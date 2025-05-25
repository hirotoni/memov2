package repos

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/hirotoni/memov2/models"
	"github.com/hirotoni/memov2/utils"
)

type MemoRepo interface {
	MemoEntries() ([]models.MemoFileInterface, error)
	Metadata(file models.MemoFileInterface) (map[string]interface{}, error)
	Save(file models.MemoFileInterface, truncate bool) error
	TidyMemos() error
}

type memoRepoImpl struct {
	dir string
}

func NewMemoRepo(dir string) MemoRepo {
	return &memoRepoImpl{
		dir: dir,
	}
}

func (r *memoRepoImpl) MemoEntries() ([]models.MemoFileInterface, error) {
	reg, err := regexp.Compile(models.FileNameRegexMemo)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}
	var files []models.MemoFileInterface
	err = filepath.Walk(r.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}
		if path == r.dir {
			return nil // Skip the root directory itself
		}

		if reg.MatchString(info.Name()) {
			mm, err := memofilefrominfo(path, info)
			if err != nil {
				// Log the error but continue processing other files
				fmt.Printf("Warning: skipping file %s due to error: %v\n", path, err)
				return nil
			}

			files = append(files, mm)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	slices.SortFunc(files, func(a, b models.MemoFileInterface) int {
		if a.Date().Before(b.Date()) {
			return -1
		}
		return 1
	})

	return files, nil
}

func memofilefrominfo(path string, info os.FileInfo) (models.MemoFileInterface, error) {
	datetimeReg, err := regexp.Compile(models.FileNameDateTimeRegexMemo)
	if err != nil {
		return nil, fmt.Errorf("invalid date regex pattern: %w", err)
	}
	datetimestring := datetimeReg.FindString(info.Name())
	if datetimestring == "" {
		return nil, fmt.Errorf("invalid date format in filename: %s", info.Name())
	}

	// 日付
	date, err := time.Parse(models.FileNameDateLayoutMemo, datetimestring)
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

	mm, err := models.NewMemoFile(date, title, category)
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

func (r *memoRepoImpl) Metadata(f models.MemoFileInterface) (map[string]interface{}, error) {
	fpath := filepath.Join(r.dir, filepath.Join(f.CategoryTree()...), f.FileName())

	b, err := os.ReadFile(fpath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", fpath, err)
	}

	h := utils.NewMarkdownHandler()
	metadata := h.Metadata(b)

	return metadata, nil
}

func (r *memoRepoImpl) Save(file models.MemoFileInterface, truncate bool) error {
	// initialize location
	locationPath := filepath.Join(r.dir, file.Location())
	if !utils.Exists(locationPath) {
		err := os.MkdirAll(locationPath, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %v", locationPath, err)
		}
		fmt.Println("Created directory:", locationPath)
	}

	path := filepath.Join(locationPath, file.FileName())

	if !utils.Exists(path) || truncate {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// write metadata
		f.WriteString(metadataString(file))

		// write title
		f.WriteString("# " + file.Title() + "\n\n")

		// write heading blocks
		for _, v := range file.HeadingBlocks() {
			f.WriteString(v.String())
		}
		fmt.Printf("File saved: %s\n", path)
	}

	return nil
}

func metadataString(file models.MemoFileInterface) string {
	wrap := func(s string) string {
		return "\"" + s + "\""
	}

	sb := strings.Builder{}
	sb.WriteString("---\n")
	sb.WriteString("category: ")
	if len(file.CategoryTree()) > 0 {
		sb.WriteString("[")
		for i, v := range file.CategoryTree() {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(wrap(v))
		}
		sb.WriteString("]")
	} else {
		sb.WriteString("[]")
	}
	sb.WriteString("\n")
	sb.WriteString("---\n\n")

	return sb.String()
}

func (r *memoRepoImpl) TidyMemos() error {
	err := r.moveFilesToCorrectLocation()
	if err != nil {
		return fmt.Errorf("error moving files to correct location: %w", err)
	}
	err = r.removeEmptyDirectories()
	if err != nil {
		return fmt.Errorf("error removing empty directories: %w", err)
	}
	fmt.Println("Memo files tidied successfully.")
	return nil
}

// traverse directory and if path does not match with the location computed from metadata, move it to the correct location
func (r *memoRepoImpl) moveFilesToCorrectLocation() error {
	return filepath.Walk(r.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil // Skip directories and non-md files
		}
		if info.Name() == "weekly_report.md" {
			return nil
		}
		if info.Name() == "index.md" {
			return nil
		}

		mm, err := memofilefrominfo(path, info)
		if err != nil {
			return fmt.Errorf("error creating MemoFile from info: %w", err)
		}

		category := mm.CategoryTree()

		location := filepath.Join(r.dir, filepath.Join(category...), info.Name())

		if path != location {
			err = os.MkdirAll(filepath.Dir(location), 0755)
			if err != nil {
				return fmt.Errorf("failed to create directory for %s: %w", location, err)
			}
			fmt.Printf("Moving file %s to %s\n", path, location)
			err := os.Rename(path, location)
			if err != nil {
				return fmt.Errorf("failed to move file %s to %s: %w", path, location, err)
			}
			fmt.Printf("Moved file %s to %s\n", path, location)
		}

		return nil
	})
}

// traverse directory and if empty derectory, remove it
func (r *memoRepoImpl) removeEmptyDirectories() error {
	var limit = 10 // limit to prevent infinite loop

	removeFunc := func() (int, error) {
		count := 0
		err := filepath.Walk(r.dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && path != r.dir {
				entries, err := os.ReadDir(path)
				if err != nil {
					return fmt.Errorf("failed to read directory %s: %w", path, err)
				}
				if len(entries) == 0 {
					err = os.Remove(path)
					if err != nil {
						return fmt.Errorf("failed to remove empty directory %s: %w", path, err)
					}
					count++
					fmt.Printf("Removed empty directory: %s\n", path)
				}
			}
			return nil
		})

		if err != nil {
			return 0, fmt.Errorf("error walking through directories: %w", err)
		}

		return count, nil
	}

	for range limit {
		c, err := removeFunc()
		if err != nil {
			return fmt.Errorf("error removing empty directories: %w", err)
		} else if c == 0 {
			break // No more empty directories to remove
		}
	}

	return nil
}
