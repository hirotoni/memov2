package memo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/platform"
)

func (uc memo) GenerateMemoFile(title string, categoryTree []string) error {
	if title == "" {
		var err error
		title, err = platform.ReadLine("Title: ")
		if err != nil {
			return err
		}
	}

	// Create a new memo file with the given title
	today := time.Now()
	memoFile, err := domain.NewMemoFile(today, title, categoryTree)
	if err != nil {
		return err
	}

	// Save the memo file to the base directory
	err = uc.repos.Memo().Save(memoFile, false)
	if err != nil {
		return err
	}

	fpath := filepath.Join(uc.config.MemosDir(), memoFile.Location(), memoFile.FileName())

	err = uc.editor.Open(uc.config.BaseDir(), fpath)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error opening editor")
	}

	return nil
}

// CreateNonInteractive creates a memo without launching the TUI or the editor.
// It is the write counterpart of the read-only `--json` commands and backs
// `memos create`. The title is REQUIRED: an empty title is a validation error
// (it must NOT fall through to the interactive prompt). The body (typically read
// from stdin) is attached to the memo as its top-level body content; an empty
// body creates a memo with just its title heading.
//
// On success it writes the new memo's relative path (one line) to stdout, or a
// createResultDTO when asJSON is true. Logs and errors go to stderr.
func (uc memo) CreateNonInteractive(title string, categoryTree []string, body string, asJSON bool) error {
	// Title is untrusted AI/automation input: reject unsafe titles (path
	// separators, traversal) up front and store the trimmed form.
	title, err := validateAndTrimTitle(title)
	if err != nil {
		return err
	}

	// Reject an unsafe category tree before it is joined into a filesystem path,
	// so the collision check below always operates on a path under MemosDir.
	if err := domain.ValidateCategoryTree(categoryTree); err != nil {
		return common.Wrap(err, common.ErrorTypeValidation, "invalid category tree")
	}

	today := time.Now()
	memoFile, err := domain.NewMemoFile(today, title, categoryTree)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeValidation, "failed to create memo file")
	}

	// Attach the stdin body as top-level body content. Trailing newlines are
	// trimmed because the serializer appends its own single trailing newline.
	if trimmed := strings.TrimRight(body, "\n"); trimmed != "" {
		memoFile.SetTopLevelBodyContent(&markdown.HeadingBlock{ContentText: trimmed})
	}

	relPath := filepath.Join(memoFile.Location(), memoFile.FileName())

	// Save(_, false) is a no-op when the target already exists, which for scripted
	// bulk-create would silently drop this memo's body and falsely report success.
	// Detect the collision and error instead of reporting a misleading success.
	targetPath := filepath.Join(uc.config.MemosDir(), memoFile.Location(), memoFile.FileName())
	if platform.Exists(targetPath) {
		return common.New(common.ErrorTypeValidation, fmt.Sprintf("memo already exists: %s", relPath))
	}

	// Persist without truncating an existing file and WITHOUT opening the editor.
	if err := uc.repos.Memo().Save(memoFile, false); err != nil {
		return err
	}

	if asJSON {
		// Normalize the category so it always marshals to a JSON array, never null.
		cat := categoryTree
		if cat == nil {
			cat = []string{}
		}
		return writeJSON(createResultDTO{
			Title:    memoFile.Title(),
			Path:     relPath,
			Category: cat,
			Date:     memoFile.Date(),
		})
	}

	fmt.Fprintln(os.Stdout, relPath)
	return nil
}
