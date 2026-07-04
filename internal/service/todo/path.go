package todo

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hirotoni/memov2/internal/common"
)

// resolveToTodosDir resolves a user-provided path to an absolute path under
// todosDir. It handles three cases:
//  1. Absolute path → cleaned as-is
//  2. Relative path that resolves (via CWD) to a location under todosDir → use the resolved absolute path
//  3. Relative path from list command output (e.g. "20250114Mon_todos.md") → join with todosDir
//
// Regardless of which case applies, the returned absolute path is guaranteed to be
// todosDir itself or a location under it. Any path that escapes todosDir (e.g. via
// ".." segments or an absolute path pointing elsewhere) is rejected with an
// ErrorTypeValidation error. This mirrors resolveToMemosDir in the memo package.
func resolveToTodosDir(todosDir, path string) (string, error) {
	cleanRoot := filepath.Clean(todosDir)

	var resolved string
	if filepath.IsAbs(path) {
		resolved = filepath.Clean(path)
	} else if absPath, err := filepath.Abs(path); err == nil && isUnderRoot(filepath.Clean(absPath), cleanRoot) {
		resolved = filepath.Clean(absPath)
	} else {
		resolved = filepath.Clean(filepath.Join(todosDir, path))
	}

	// Containment is purely lexical: it defends against ".." text traversal and
	// absolute paths pointing elsewhere, but NOT against a symlink located inside
	// todosDir that points outside the tree (we do not call filepath.EvalSymlinks).
	// That symlink escape is an accepted, out-of-scope threat model, matching the
	// residual risk documented for the memo package.
	if !isUnderRoot(resolved, cleanRoot) {
		return "", common.New(common.ErrorTypeValidation, fmt.Sprintf("path escapes todos directory: %s", path))
	}

	return resolved, nil
}

// isUnderRoot reports whether cleanPath is cleanRoot itself or a descendant of it.
// Both arguments must already be cleaned. The separator suffix on the prefix check
// prevents sibling directories with a shared prefix (e.g. "/a/todos-evil" vs
// "/a/todos") from being falsely accepted.
func isUnderRoot(cleanPath, cleanRoot string) bool {
	return cleanPath == cleanRoot || strings.HasPrefix(cleanPath, cleanRoot+string(filepath.Separator))
}
