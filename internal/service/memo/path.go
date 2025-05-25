package memo

import (
	"path/filepath"
	"strings"
)

// resolveToMemosDir resolves a user-provided path to an absolute path under memosDir.
// It handles three cases:
//  1. Absolute path → returned as-is
//  2. Relative path that resolves (via CWD) to a location under memosDir → use the resolved absolute path
//  3. Relative path from list command output (e.g. "category/file.md") → join with memosDir
func resolveToMemosDir(memosDir, path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	absPath, err := filepath.Abs(path)
	if err == nil {
		cleanAbs := filepath.Clean(absPath)
		cleanMemosDir := filepath.Clean(memosDir) + string(filepath.Separator)
		if strings.HasPrefix(cleanAbs, cleanMemosDir) {
			return cleanAbs
		}
	}

	return filepath.Join(memosDir, path)
}
