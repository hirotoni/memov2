package memo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hirotoni/memov2/internal/common"
)

// Show prints a single memo identified by path. The path is resolved through
// resolveToMemosDir so any location escaping the memos directory is rejected.
// In text mode the memo's raw markdown file content is written to stdout (cat
// equivalent); in json mode a showDTO is written instead. A path that does not
// correspond to an existing memo file is rejected with ErrorTypeValidation.
func (uc memo) Show(path string, asJSON bool) error {
	memosDir := uc.config.MemosDir()

	resolved, err := resolveToMemosDir(memosDir, path)
	if err != nil {
		return err
	}

	fileName := filepath.Base(resolved)

	// Derive the location (category tree) relative to memosDir.
	dir := filepath.Dir(resolved)
	relDir, err := filepath.Rel(memosDir, dir)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to compute relative path")
	}
	location := ""
	if relDir != "." {
		location = filepath.Join(strings.Split(relDir, string(filepath.Separator))...)
	}

	// Find the parsed memo entry matching the resolved path.
	entries, err := uc.repos.Memo().MemoEntries()
	if err != nil {
		return err
	}

	var found = -1
	for i, m := range entries {
		if m.FileName() == fileName && m.Location() == location {
			found = i
			break
		}
	}
	if found == -1 {
		return common.New(common.ErrorTypeValidation, fmt.Sprintf("memo not found: %s", path))
	}
	m := entries[found]

	if asJSON {
		// In json mode, prefer the relative path (stable, list-compatible).
		relPath := filepath.Join(m.Location(), m.FileName())
		return writeJSON(newShowDTO(m, relPath))
	}

	// Text mode: emit the raw markdown file content.
	content, err := os.ReadFile(resolved)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to read memo file")
	}
	if _, err := os.Stdout.Write(content); err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to write memo content")
	}
	return nil
}
