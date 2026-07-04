package memo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/domain"
)

// Rename renames the memo at path (non-interactive mode, backs
// `memos rename <path> <new-title>`). On success it writes the renamed memo's new
// relative path (one line) to stdout, or a renameResultDTO when asJSON is true.
// The interactive TUI path (RenameInteractive) uses renameMemo directly and must
// NOT emit this output.
func (uc memo) Rename(path string, newTitle string, asJSON bool) error {
	dto, err := uc.renameMemo(path, newTitle)
	if err != nil {
		return err
	}

	if asJSON {
		return writeJSON(dto)
	}

	fmt.Fprintln(os.Stdout, dto.NewPath)
	return nil
}

// renameMemo performs the path-safe rename and returns the resulting DTO without
// writing anything to stdout. It is shared by the non-interactive Rename (which
// adds output) and the interactive TUI path (which does not).
func (uc memo) renameMemo(path string, newTitle string) (renameResultDTO, error) {
	// The new title is untrusted input, symmetric with create: reject unsafe
	// titles (empty/whitespace-only, path separators, traversal) before they reach
	// domain.FileName() and store the trimmed form.
	newTitle, err := validateAndTrimTitle(newTitle)
	if err != nil {
		return renameResultDTO{}, err
	}

	memosDir := uc.config.MemosDir()

	// Resolve path
	path, err = resolveToMemosDir(memosDir, path)
	if err != nil {
		return renameResultDTO{}, err
	}

	// Extract filename and category from path
	fileName := filepath.Base(path)

	// Validate filename matches memo pattern
	if domain.MemoTitle(fileName) == fileName {
		return renameResultDTO{}, common.New(common.ErrorTypeValidation, fmt.Sprintf("invalid memo file path: %s", path))
	}

	// Get category tree from the path relative to memosDir
	dir := filepath.Dir(path)
	relDir, err := filepath.Rel(memosDir, dir)
	if err != nil {
		return renameResultDTO{}, common.Wrap(err, common.ErrorTypeFileSystem, "failed to compute relative path")
	}

	var categoryTree []string
	if relDir != "." {
		categoryTree = strings.Split(relDir, string(filepath.Separator))
	}

	// Find the memo by scanning entries
	entries, err := uc.repos.Memo().MemoEntries()
	if err != nil {
		return renameResultDTO{}, err
	}

	location := filepath.Join(categoryTree...)
	for _, m := range entries {
		if m.FileName() == fileName && m.Location() == location {
			oldTitle := m.Title()
			oldRel := filepath.Join(m.Location(), m.FileName())

			if err := uc.repos.Memo().Rename(m, newTitle); err != nil {
				return renameResultDTO{}, err
			}

			// Compute the new relative path from the new title (same date/category).
			renamed, err := domain.NewMemoFile(m.Date(), newTitle, m.CategoryTree())
			if err != nil {
				return renameResultDTO{}, common.Wrap(err, common.ErrorTypeValidation, "failed to compute renamed path")
			}
			newRel := filepath.Join(renamed.Location(), renamed.FileName())

			return renameResultDTO{
				Path:     oldRel,
				OldTitle: oldTitle,
				NewTitle: newTitle,
				NewPath:  newRel,
			}, nil
		}
	}

	return renameResultDTO{}, common.New(common.ErrorTypeValidation, fmt.Sprintf("memo not found: %s", path))
}
