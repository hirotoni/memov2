package memo

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hirotoni/memov2/internal/common"
)

// resolveToMemosDir resolves a user-provided path to an absolute path under memosDir.
// It handles three cases:
//  1. Absolute path → cleaned as-is
//  2. Relative path that resolves (via CWD) to a location under memosDir → use the resolved absolute path
//  3. Relative path from list command output (e.g. "category/file.md") → join with memosDir
//
// Regardless of which case applies, the returned absolute path is guaranteed to be
// memosDir itself or a location under it. Any path that escapes memosDir (e.g. via
// ".." segments or an absolute path pointing elsewhere) is rejected with an
// ErrorTypeValidation error.
func resolveToMemosDir(memosDir, path string) (string, error) {
	cleanRoot := filepath.Clean(memosDir)

	var resolved string
	if filepath.IsAbs(path) {
		resolved = filepath.Clean(path)
	} else if absPath, err := filepath.Abs(path); err == nil && isUnderRoot(filepath.Clean(absPath), cleanRoot) {
		resolved = filepath.Clean(absPath)
	} else {
		resolved = filepath.Clean(filepath.Join(memosDir, path))
	}

	// Containment is purely lexical: it defends against ".." text traversal and
	// absolute paths pointing elsewhere, but NOT against a symlink located inside
	// memosDir that points outside the tree (we do not call filepath.EvalSymlinks).
	// That symlink escape is an accepted, out-of-scope threat model, documented as a
	// residual risk in docs/plans/plan_20260704081955.md.
	if !isUnderRoot(resolved, cleanRoot) {
		return "", common.New(common.ErrorTypeValidation, fmt.Sprintf("path escapes memos directory: %s", path))
	}

	return resolved, nil
}

// isUnderRoot reports whether cleanPath is cleanRoot itself or a descendant of it.
// Both arguments must already be cleaned. The separator suffix on the prefix check
// prevents sibling directories with a shared prefix (e.g. "/a/memos-evil" vs
// "/a/memos") from being falsely accepted.
func isUnderRoot(cleanPath, cleanRoot string) bool {
	return cleanPath == cleanRoot || strings.HasPrefix(cleanPath, cleanRoot+string(filepath.Separator))
}
