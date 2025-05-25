package memo

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/domain"
)

func (uc memo) Rename(path string, newTitle string) error {
	memosDir := uc.config.MemosDir()

	// Resolve path
	path = resolveToMemosDir(memosDir, path)

	// Extract filename and category from path
	fileName := filepath.Base(path)

	// Validate filename matches memo pattern
	if domain.MemoTitle(fileName) == fileName {
		return common.New(common.ErrorTypeValidation, fmt.Sprintf("invalid memo file path: %s", path))
	}

	// Get category tree from the path relative to memosDir
	dir := filepath.Dir(path)
	relDir, err := filepath.Rel(memosDir, dir)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to compute relative path")
	}

	var categoryTree []string
	if relDir != "." {
		categoryTree = strings.Split(relDir, string(filepath.Separator))
	}

	// Find the memo by scanning entries
	entries, err := uc.repos.Memo().MemoEntries()
	if err != nil {
		return err
	}

	location := filepath.Join(categoryTree...)
	for _, m := range entries {
		if m.FileName() == fileName && m.Location() == location {
			return uc.repos.Memo().Rename(m, newTitle)
		}
	}

	return common.New(common.ErrorTypeValidation, fmt.Sprintf("memo not found: %s", path))
}
