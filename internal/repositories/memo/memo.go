package memo

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/platform"
	repoCommon "github.com/hirotoni/memov2/internal/repositories/common"
)

type memo struct {
	dir    string
	logger *slog.Logger
}

func NewMemo(dir string, logger *slog.Logger) interfaces.MemoRepo {
	return &memo{dir: dir, logger: logger}
}

func (r *memo) Memo(file interfaces.MemoFileInterface) (interfaces.MemoFileInterface, error) {
	path := filepath.Join(r.dir, file.Location(), file.FileName())

	info, err := os.Stat(path)
	if err != nil {
		return nil, common.Wrap(err, common.ErrorTypeFileSystem, "error getting file info")
	}

	return memofilefromosfileinfo(path, info, r.logger)
}

func (r *memo) MemoEntries() ([]interfaces.MemoFileInterface, error) {
	files, err := r.scanDirectory()
	if err != nil {
		return nil, err
	}

	slices.SortFunc(files, func(a, b domain.MemoFileInterface) int {
		if a.Date().Before(b.Date()) {
			return -1
		}
		return 1
	})

	return files, nil
}

func (r *memo) scanDirectory() ([]interfaces.MemoFileInterface, error) {
	// Check if directory exists before walking
	if _, err := os.Stat(r.dir); os.IsNotExist(err) {
		// Directory doesn't exist, return empty list
		return []interfaces.MemoFileInterface{}, nil
	}

	reg, err := regexp.Compile(domain.FileNameRegexMemo)
	if err != nil {
		return nil, common.Wrap(err, common.ErrorTypeRepository, "invalid regex pattern")
	}
	var files []interfaces.MemoFileInterface
	err = filepath.Walk(r.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return common.Wrap(err, common.ErrorTypeFileSystem, "error walking path")
		}
		if path == r.dir {
			return nil // Skip the root directory itself
		}

		if reg.MatchString(info.Name()) {
			mm, err := memofilefromosfileinfo(path, info, r.logger)
			if err != nil {
				return common.Wrap(err, common.ErrorTypeRepository, "error creating MemoFile from info")
			}

			files = append(files, mm)
		}

		return nil
	})
	if err != nil {
		return nil, common.Wrap(err, common.ErrorTypeFileSystem, "error walking directory")
	}
	return files, nil
}

func memofilefromosfileinfo(path string, info os.FileInfo, logger *slog.Logger) (interfaces.MemoFileInterface, error) {
	// 日付抽出（共通パーサーを使用）
	date, err := repoCommon.ParseDateFromFilename(info.Name(), repoCommon.DateParserConfig{
		DateTimeRegex: domain.FileNameDateTimeRegexMemo,
		DateLayout:    domain.FileNameDateLayoutMemo,
	})
	if err != nil {
		return nil, err
	}

	// タイトル抽出（Memo固有）
	datetimestring := regexp.MustCompile(domain.FileNameDateTimeRegexMemo).FindString(info.Name())
	title := strings.TrimPrefix(info.Name(), datetimestring+"_memo_")
	title = strings.TrimSuffix(title, ".md")

	// Markdownファイル読み込み（共通パーサーを使用）
	b, err := repoCommon.ReadMarkdownFile(path)
	if err != nil {
		return nil, err
	}

	// Markdown解析（共通パーサーを使用）
	parser := repoCommon.NewMarkdownParser()
	meta := parser.Metadata(b)

	var category []string
	if v, ok := meta["category"]; ok {
		switch val := v.(type) {
		case []string:
			category = val
		case string:
			category = []string{val}
		case []interface{}:
			for _, item := range val {
				if str, ok := item.(string); ok {
					category = append(category, str)
				}
			}
		default:
			logger.Warn("Unexpected category type, using empty category", "type", fmt.Sprintf("%T", v), "path", path)
		}
	}

	// TopLevelBodyContent抽出（Memo固有、共通パーサーを使用）
	tlbc := parser.TopLevelBodyContent(b)
	if tlbc == nil {
		return nil, common.New(common.ErrorTypeRepository, fmt.Sprintf("failed to parse top level body content in file: %s", path))
	}

	// HeadingBlocks抽出（共通パーサーを使用）
	hbs, err := parser.HeadingBlocksByLevel(b, 2)
	if err != nil {
		return nil, common.Wrap(err, common.ErrorTypeRepository, fmt.Sprintf("failed to parse heading blocks in file: %s", path))
	}

	// Domain層のファクトリを使用してMemoFileを構築
	return domain.MemoFileFromParsedData(date, title, category, tlbc, hbs)
}

func (r *memo) Metadata(f interfaces.MemoFileInterface) (map[string]interface{}, error) {
	fpath := filepath.Join(r.dir, filepath.Join(f.CategoryTree()...), f.FileName())

	b, err := repoCommon.ReadMarkdownFile(fpath)
	if err != nil {
		return nil, err
	}

	parser := repoCommon.NewMarkdownParser()
	metadata := parser.Metadata(b)

	return metadata, nil
}

func (r *memo) Save(file interfaces.MemoFileInterface, truncate bool) error {
	// initialize location
	locationPath := filepath.Join(r.dir, file.Location())
	if err := platform.EnsureDir(locationPath); err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, fmt.Sprintf("error ensuring directory: %s", locationPath))
	}

	path := filepath.Join(locationPath, file.FileName())

	if !platform.Exists(path) || truncate {
		if err := platform.WriteFileStream(path, truncate, func(w *bufio.Writer) error {
			_, err := w.WriteString(file.ContentString())
			return err
		}); err != nil {
			return err
		}
		r.logger.Info("File saved", "path", path)
	}

	return nil
}


type CategoryCollector struct {
	memorepo      interfaces.MemoRepo
	dir           string
	seen          map[string]bool
	allCategories [][]string
}

func NewCategoryCollector(memorepo interfaces.MemoRepo, dir string) *CategoryCollector {
	return &CategoryCollector{
		memorepo:      memorepo,
		dir:           dir,
		allCategories: make([][]string, 0),
		seen:          make(map[string]bool),
	}
}

func (cc *CategoryCollector) collectFromDirectory() error {
	err := filepath.Walk(cc.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == cc.dir {
			return nil // Skip the root directory itself
		}

		// Get relative path from memos directory
		relPath, err := filepath.Rel(cc.dir, path)
		if err != nil {
			return err
		}

		// Split path into parts
		parts := strings.Split(relPath, string(filepath.Separator))

		// For each part of the path, create a category path up to that part
		for i := 0; i < len(parts); i++ {
			subPath := parts[:i+1]
			pathStr := strings.Join(subPath, string(filepath.Separator))
			if !cc.seen[pathStr] && info.IsDir() {
				cc.seen[pathStr] = true
				cc.allCategories = append(cc.allCategories, subPath)
			}
		}
		return nil
	})
	if err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to walk directory")
	}
	return nil
}

func (cc *CategoryCollector) collectFromFiles() error {
	files, err := cc.memorepo.MemoEntries()
	if err != nil {
		return common.Wrap(err, common.ErrorTypeRepository, "error getting memo entries")
	}
	for _, memo := range files {
		categoryTree := memo.CategoryTree()
		if len(categoryTree) > 0 {
			// Add each subpath of the category tree
			for i := 1; i <= len(categoryTree); i++ {
				subPath := categoryTree[:i]
				pathStr := strings.Join(subPath, string(filepath.Separator))
				if !cc.seen[pathStr] {
					cc.seen[pathStr] = true
					cc.allCategories = append(cc.allCategories, subPath)
				}
			}
		}
	}
	return nil
}

func (cc *CategoryCollector) sort() {
	// Sort the categories:
	// 1. By parent path (keep children with their parents)
	// 2. By name within the same level
	slices.SortFunc(cc.allCategories, func(a, b []string) int {
		// Compare common parent path components
		minLen := len(a)
		if len(b) < minLen {
			minLen = len(b)
		}
		for i := 0; i < minLen; i++ {
			if a[i] != b[i] {
				return strings.Compare(a[i], b[i])
			}
		}
		// If one is a prefix of the other, shorter path comes first
		if len(a) != len(b) {
			return len(a) - len(b)
		}
		// If same length and same prefix, compare the last component
		return strings.Compare(a[len(a)-1], b[len(b)-1])
	})
}

func (r *memo) Categories() ([][]string, error) {
	cc := NewCategoryCollector(r, r.dir)
	// First, get categories from directory structure
	if err := cc.collectFromDirectory(); err != nil {
		return nil, err
	}
	// Then, get categories from memo metadata
	if err := cc.collectFromFiles(); err != nil {
		return nil, err
	}
	// Sort the categories
	cc.sort()

	return cc.allCategories, nil
}

func (r *memo) Move(file interfaces.MemoFileInterface, newCategoryTree []string) error {
	// Get current file path
	currentPath := filepath.Join(r.dir, file.Location(), file.FileName())

	// retrieve again
	mm, err := r.Memo(file)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeRepository, "error getting memo file")
	}

	// Update category tree
	mm.SetCategoryTree(newCategoryTree)

	// Create new directory if needed
	newLocation := filepath.Join(r.dir, mm.Location())
	if err := platform.EnsureDir(newLocation); err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, fmt.Sprintf("error ensuring directory: %s", newLocation))
	}

	// Create new file path
	newPath := filepath.Join(newLocation, mm.FileName())

	// Write content to new location
	err = r.Save(mm, true)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeRepository, fmt.Sprintf("failed to save file: %s", newPath))
	}

	// Remove old file only if it's different from the new location
	if currentPath != newPath {
		err = os.Remove(currentPath)
		if err != nil {
			return common.Wrap(err, common.ErrorTypeFileSystem, fmt.Sprintf("failed to remove old file: %s", currentPath))
		}
		r.logger.Info("Moved file", "from", currentPath, "to", newPath)
	}

	return nil
}

func (r *memo) Delete(file interfaces.MemoFileInterface) error {
	// Get file path
	filePath := filepath.Join(r.dir, file.Location(), file.FileName())

	// Check if file exists
	if !platform.Exists(filePath) {
		return common.New(common.ErrorTypeRepository, fmt.Sprintf("file does not exist: %s", filePath))
	}

	// Move the file to trash instead of permanently deleting
	err := platform.MoveToTrash(filePath)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, fmt.Sprintf("failed to move file to trash: %s", filePath))
	}

	return nil
}

func (r *memo) Rename(file interfaces.MemoFileInterface, newTitle string) error {
	// Get current file path
	oldPath := filepath.Join(r.dir, file.Location(), file.FileName())

	// Check if file exists
	if !platform.Exists(oldPath) {
		return common.New(common.ErrorTypeRepository, fmt.Sprintf("file does not exist: %s", oldPath))
	}

	// Retrieve the full memo with content
	mm, err := r.Memo(file)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeRepository, "error getting memo file")
	}

	// Update the title
	mm.SetTitle(newTitle)

	// Create new file path with new title
	newPath := filepath.Join(r.dir, mm.Location(), mm.FileName())

	// If the path is the same (shouldn't happen but check anyway), just update content
	if oldPath == newPath {
		// Just save with updated title in content
		if err := r.Save(mm, true); err != nil {
			return common.Wrap(err, common.ErrorTypeRepository, "failed to save renamed file")
		}
		return nil
	}

	// Save the file with new name
	if err := r.Save(mm, true); err != nil {
		return common.Wrap(err, common.ErrorTypeRepository, "failed to save renamed file")
	}

	// Remove old file
	if err := os.Remove(oldPath); err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, fmt.Sprintf("failed to remove old file: %s", oldPath))
	}

	r.logger.Info("Renamed file", "from", oldPath, "to", newPath)

	return nil
}

func (r *memo) Duplicate(file interfaces.MemoFileInterface) (interfaces.MemoFileInterface, error) {
	// Get original file path
	origPath := filepath.Join(r.dir, file.Location(), file.FileName())

	// Check if file exists
	if !platform.Exists(origPath) {
		return nil, common.New(common.ErrorTypeRepository, fmt.Sprintf("file does not exist: %s", origPath))
	}

	// Retrieve the full memo with content
	origMemo, err := r.Memo(file)
	if err != nil {
		return nil, common.Wrap(err, common.ErrorTypeRepository, "error getting original memo")
	}

	// Create new memo with current timestamp and "copied" suffix
	newTitle := origMemo.Title() + " copied"
	newMemo, err := domain.NewMemoFile(
		time.Now(),
		newTitle,
		origMemo.CategoryTree(),
	)
	if err != nil {
		return nil, common.Wrap(err, common.ErrorTypeRepository, "error creating duplicate memo")
	}

	// Copy content
	newMemo.SetTopLevelBodyContent(origMemo.TopLevelBodyContent())
	newMemo.SetHeadingBlocks(origMemo.HeadingBlocks())

	// Save the duplicate
	if err := r.Save(newMemo, true); err != nil {
		return nil, common.Wrap(err, common.ErrorTypeRepository, "failed to save duplicate")
	}

	newPath := filepath.Join(r.dir, newMemo.Location(), newMemo.FileName())
	r.logger.Info("Duplicated memo", "from", origPath, "to", newPath)

	return newMemo, nil
}
