package todo

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hirotoni/memov2/internal/common"
)

// Show prints a single todo file identified by path. The path is resolved through
// resolveToTodosDir so any location escaping the todos directory is rejected. In
// text mode the file's raw markdown content is written to stdout (cat equivalent),
// which also works for the weekly report and template files. In json mode a
// todoShowDTO is built from the matching daily todo entry (weekly/template files
// are not covered by --json — read them in text mode). A path that does not
// correspond to an existing daily todo file is rejected with ErrorTypeValidation.
func (uc todo) Show(path string, asJSON bool) error {
	todosDir := uc.c.TodosDir()

	resolved, err := resolveToTodosDir(todosDir, path)
	if err != nil {
		return err
	}

	if asJSON {
		fileName := filepath.Base(resolved)

		// Find the parsed daily todo entry matching the resolved path. Todos are
		// stored flat in TodosDir, so the file name is the only key.
		entries, err := uc.r.Todo().TodoEntries()
		if err != nil {
			return common.Wrap(err, common.ErrorTypeService, "error fetching todo entries")
		}

		for _, t := range entries {
			if t.FileName() == fileName {
				return writeJSON(newTodoShowDTO(t, fileName))
			}
		}
		return common.New(common.ErrorTypeValidation, fmt.Sprintf("todo not found: %s", path))
	}

	// Text mode: emit the raw markdown file content (also valid for weekly/template).
	content, err := os.ReadFile(resolved)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to read todo file")
	}
	if _, err := os.Stdout.Write(content); err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to write todo content")
	}
	return nil
}
