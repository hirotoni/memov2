package memo

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/platform"
	repoCommon "github.com/hirotoni/memov2/internal/repositories/common"
)

// TidyMemos organizes memo files by moving them to correct locations based on metadata
// and removes empty directories
func (uc memo) TidyMemos() error {
	err := uc.moveFilesToCorrectLocation()
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error moving files to correct location")
	}
	err = uc.removeEmptyDirectories()
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error removing empty directories")
	}
	return nil
}

// moveFilesToCorrectLocation traverses directory and moves files to correct location
// based on metadata category
func (uc memo) moveFilesToCorrectLocation() error {
	memosDir := uc.config.MemosDir()
	// Check if directory exists before walking
	if _, err := os.Stat(memosDir); os.IsNotExist(err) {
		// Directory doesn't exist, nothing to tidy
		return nil
	}

	return filepath.Walk(memosDir, func(sourcePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !uc.shouldProcess(info) {
			return nil
		}

		return uc.moveFilesIfNeeded(sourcePath, info)
	})
}

func (uc memo) shouldProcess(info os.FileInfo) bool {
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

func (uc memo) moveFilesIfNeeded(sourcePath string, info os.FileInfo) error {
	// Parse date and title from filename to create a minimal memo file interface
	date, err := parseDateFromMemoFilename(info.Name())
	if err != nil {
		// Skip files we can't parse
		return nil
	}
	title := extractTitleFromMemoFilename(info.Name())
	tempMemo, err := domain.NewMemoFile(date, title, []string{})
	if err != nil {
		return nil
	}

	// Load the actual memo to get its category
	loadedMemo, err := uc.repos.Memo().Memo(tempMemo)
	if err != nil {
		// Skip files we can't load
		return nil
	}

	category := loadedMemo.CategoryTree()
	targetPath := filepath.Join(uc.config.MemosDir(), filepath.Join(category...), info.Name())

	if sourcePath != targetPath {
		if err := platform.EnsureDir(filepath.Dir(targetPath)); err != nil {
			return common.Wrap(err, common.ErrorTypeFileSystem, fmt.Sprintf("failed to create directory: %s", targetPath))
		}
		uc.logger.Info("Moving file", "source", sourcePath, "target", targetPath)
		err := os.Rename(sourcePath, targetPath)
		if err != nil {
			return common.Wrap(err, common.ErrorTypeFileSystem, fmt.Sprintf("failed to move file from %s to %s", sourcePath, targetPath))
		}
		uc.logger.Info("Moved file", "source", sourcePath, "target", targetPath)
	}

	return nil
}

// removeEmptyDirectories removes empty directories in multiple passes
func (uc memo) removeEmptyDirectories() error {
	var limit = 10 // limit to prevent infinite loop

	for range limit {
		c, err := uc.removeEmptyDirectoriesInOnePass()
		if err != nil {
			return common.Wrap(err, common.ErrorTypeService, "error removing empty directories")
		} else if c == 0 {
			break // No more empty directories to remove
		}
	}

	return nil
}

func (uc memo) removeEmptyDirectoriesInOnePass() (int, error) {
	memosDir := uc.config.MemosDir()
	// Check if directory exists before walking
	if _, err := os.Stat(memosDir); os.IsNotExist(err) {
		// Directory doesn't exist, nothing to remove
		return 0, nil
	}

	count := 0
	err := filepath.Walk(memosDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !uc.shouldCheckDirectory(info, path) {
			return nil
		}

		if uc.isDirectoryEmpty(path) {
			err = os.Remove(path)
			if err != nil {
				return common.Wrap(err, common.ErrorTypeFileSystem, fmt.Sprintf("failed to remove empty directory: %s", path))
			}
			count++
			uc.logger.Info("Removed empty directory", "path", path)
		}

		return nil
	})

	if err != nil {
		return 0, common.Wrap(err, common.ErrorTypeFileSystem, "error walking through directories")
	}

	return count, nil
}

func (uc memo) shouldCheckDirectory(info os.FileInfo, path string) bool {
	return info.IsDir() && path != uc.config.MemosDir()
}

func (uc memo) isDirectoryEmpty(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	return len(entries) == 0
}

// parseDateFromMemoFilename extracts date from memo filename
func parseDateFromMemoFilename(filename string) (time.Time, error) {
	date, err := repoCommon.ParseDateFromFilename(filename, repoCommon.DateParserConfig{
		DateTimeRegex: domain.FileNameDateTimeRegexMemo,
		DateLayout:    domain.FileNameDateLayoutMemo,
	})
	return date, err
}

// extractTitleFromMemoFilename extracts title from memo filename
func extractTitleFromMemoFilename(filename string) string {
	datetimestring := regexp.MustCompile(domain.FileNameDateTimeRegexMemo).FindString(filename)
	title := strings.TrimPrefix(filename, datetimestring+"_memo_")
	title = strings.TrimSuffix(title, ".md")
	return title
}

