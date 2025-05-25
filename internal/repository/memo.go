package repository

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/platform/fs"
	"github.com/hirotoni/memov2/utils"
)

type memo struct {
	dir string
}

func NewMemo(dir string) Memo {
	return &memo{dir: dir}
}

func (r *memo) Memo(file domain.MemoFileInterface) (domain.MemoFileInterface, error) {
	path := filepath.Join(r.dir, file.Location(), file.FileName())

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("error getting file info: %w", err)
	}

	return memofilefromosfileinfo(path, info)
}

func (r *memo) MemoEntries() ([]domain.MemoFileInterface, error) {
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

func (r *memo) scanDirectory() ([]domain.MemoFileInterface, error) {
	reg, err := regexp.Compile(domain.FileNameRegexMemo)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}
	var files []domain.MemoFileInterface
	err = filepath.Walk(r.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}
		if path == r.dir {
			return nil // Skip the root directory itself
		}

		if reg.MatchString(info.Name()) {
			mm, err := memofilefromosfileinfo(path, info)
			if err != nil {
				return fmt.Errorf("error creating MemoFile from info: %w", err)
			}

			files = append(files, mm)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}
	return files, nil
}

func memofilefromosfileinfo(path string, info os.FileInfo) (domain.MemoFileInterface, error) {
	datetimeReg, err := regexp.Compile(domain.FileNameDateTimeRegexMemo)
	if err != nil {
		return nil, fmt.Errorf("invalid date regex pattern: %w", err)
	}
	datetimestring := datetimeReg.FindString(info.Name())
	if datetimestring == "" {
		return nil, fmt.Errorf("invalid date format in filename: %s", info.Name())
	}

	// 日付
	date, err := time.Parse(domain.FileNameDateLayoutMemo, datetimestring)
	if err != nil {
		return nil, fmt.Errorf("error parsing date from filename: %w", err)
	}

	// タイトル
	title := strings.TrimPrefix(info.Name(), datetimestring+"_memo_")
	title = strings.TrimSuffix(title, ".md")

	// category tree
	h := utils.NewMarkdownHandler()
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	meta := h.Metadata(b)

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
			fmt.Printf("Unexpected category type: %T, using empty category\n", v)
		}
	}

	mm, err := domain.NewMemoFile(date, title, category)
	if err != nil {
		return nil, fmt.Errorf("error creating new MemoFile: %w", err)
	}

	tlbc := h.TopLevelBodyContent(b)
	if tlbc == nil {
		return nil, fmt.Errorf("failed to parse top level body content in file %s", path)
	}
	mm.SetTopLevelBodyContent(tlbc)

	hbs, err := h.HeadingBlocksByLevel(b, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to parse heading blocks in file %s: %w", path, err)
	}
	mm.SetHeadingBlocks(hbs)

	return mm, nil
}

func (r *memo) Metadata(f domain.MemoFileInterface) (map[string]interface{}, error) {
	fpath := filepath.Join(r.dir, filepath.Join(f.CategoryTree()...), f.FileName())

	b, err := os.ReadFile(fpath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", fpath, err)
	}

	h := utils.NewMarkdownHandler()
	metadata := h.Metadata(b)

	return metadata, nil
}

func (r *memo) Save(file domain.MemoFileInterface, truncate bool) error {
	// initialize location
	locationPath := filepath.Join(r.dir, file.Location())
	if err := fs.EnsureDir(locationPath); err != nil {
		return fmt.Errorf("error ensuring directory %s: %v", locationPath, err)
	}

	path := filepath.Join(locationPath, file.FileName())

	if !fs.Exists(path) || truncate {
		if err := fs.WriteFileStream(path, truncate, func(w *bufio.Writer) error {
			_, err := w.Write(file.Data())
			return err
		}); err != nil {
			return err
		}
		fmt.Printf("File saved: %s\n", path)
	}

	return nil
}

func (r *memo) TidyMemos() error {
	err := r.moveFilesToCorrectLocation()
	if err != nil {
		return fmt.Errorf("error moving files to correct location: %w", err)
	}
	err = r.removeEmptyDirectories()
	if err != nil {
		return fmt.Errorf("error removing empty directories: %w", err)
	}
	return nil
}

// traverse directory and if path does not match with the location computed from metadata, move it to the correct location
func (r *memo) moveFilesToCorrectLocation() error {
	return filepath.Walk(r.dir, func(sourcePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !r.shouldProcess(info) {
			return nil
		}

		return r.moveFilesIfNeeded(sourcePath, info)
	})
}

func (r *memo) shouldProcess(info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}
	if !strings.HasSuffix(info.Name(), domain.FileExtension) {
		return false
	}

	// skip special files
	specialFiles := map[string]bool{
		"weekly_report.md": true,
		"index.md":         true,
	}

	return !specialFiles[info.Name()]
}

func (r *memo) moveFilesIfNeeded(sourcePath string, info os.FileInfo) error {
	mm, err := memofilefromosfileinfo(sourcePath, info)
	if err != nil {
		return fmt.Errorf("error creating MemoFile from info: %w", err)
	}

	category := mm.CategoryTree()
	targetPath := filepath.Join(r.dir, filepath.Join(category...), info.Name())

	if sourcePath != targetPath {
		if err := fs.EnsureDir(filepath.Dir(targetPath)); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", targetPath, err)
		}
		fmt.Printf("Moving file %s to %s\n", sourcePath, targetPath)
		err := os.Rename(sourcePath, targetPath)
		if err != nil {
			return fmt.Errorf("failed to move file %s to %s: %w", sourcePath, targetPath, err)
		}
		fmt.Printf("Moved file %s to %s\n", sourcePath, targetPath)
	}

	return nil
}

// traverse directory and if empty derectory, remove it
func (r *memo) removeEmptyDirectories() error {
	var limit = 10 // limit to prevent infinite loop

	for range limit {
		c, err := r.removeEmptyDirectoriesInOnePass()
		if err != nil {
			return fmt.Errorf("error removing empty directories: %w", err)
		} else if c == 0 {
			break // No more empty directories to remove
		}
	}

	return nil
}

func (r *memo) removeEmptyDirectoriesInOnePass() (int, error) {
	count := 0
	err := filepath.Walk(r.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !r.shouldCheckDirectory(info, path) {
			return nil
		}

		if r.isDirectoryEmpty(path) {
			err = os.Remove(path)
			if err != nil {
				return fmt.Errorf("failed to remove empty directory %s: %w", path, err)
			}
			count++
			fmt.Printf("Removed empty directory: %s\n", path)
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("error walking through directories: %w", err)
	}

	return count, nil
}

func (r *memo) shouldCheckDirectory(info os.FileInfo, path string) bool {
	return info.IsDir() && path != r.dir
}

func (r *memo) isDirectoryEmpty(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	return len(entries) == 0
}

type CategoryCollector struct {
	memorepo      Memo
	dir           string
	seen          map[string]bool
	allCategories [][]string
}

func NewCategoryCollector(memorepo Memo, dir string) *CategoryCollector {
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
		return fmt.Errorf("failed to walk directory: %w", err)
	}
	return nil
}

func (cc *CategoryCollector) collectFromFiles() error {
	files, err := cc.memorepo.MemoEntries()
	if err != nil {
		return fmt.Errorf("error getting memo entries: %w", err)
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

func (r *memo) Move(file domain.MemoFileInterface, newCategoryTree []string) error {
	// Get current file path
	currentPath := filepath.Join(r.dir, file.Location(), file.FileName())

	// retrieve again
	mm, err := r.Memo(file)
	if err != nil {
		return fmt.Errorf("error getting memo file: %w", err)
	}

	// Update category tree
	mm.SetCategoryTree(newCategoryTree)

	// Create new directory if needed
	newLocation := filepath.Join(r.dir, mm.Location())
	if err := fs.EnsureDir(newLocation); err != nil {
		return fmt.Errorf("error ensuring directory %s: %v", newLocation, err)
	}

	// Create new file path
	newPath := filepath.Join(newLocation, mm.FileName())

	// Write content to new location
	err = r.Save(mm, true)
	if err != nil {
		return fmt.Errorf("failed to save file %s: %w", newPath, err)
	}

	// Remove old file only if it's different from the new location
	if currentPath != newPath {
		err = os.Remove(currentPath)
		if err != nil {
			return fmt.Errorf("failed to remove old file %s: %w", currentPath, err)
		}
		fmt.Printf("Moved file from %s to %s\n", currentPath, newPath)
	}

	// Clean up empty directories
	err = r.removeEmptyDirectories()
	if err != nil {
		return fmt.Errorf("failed to clean up empty directories: %w", err)
	}

	return nil
}
