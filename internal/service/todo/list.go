package todo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/mattn/go-runewidth"
)

// List prints all daily todo files. In text mode each line is "date\tpath"
// (aligned with runewidth); in json mode a []todoListItemDTO is written. By
// default the full path (under TodosDir) is printed; showFullPath=false emits the
// TodosDir-relative path (just the file name, since todos are stored flat).
func (uc todo) List(showFullPath, asJSON bool) error {
	entries, err := uc.r.Todo().TodoEntries()
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error fetching todo entries")
	}

	todosDir := uc.c.TodosDir()

	pathOf := func(t interface{ FileName() string }) string {
		if showFullPath {
			return filepath.Join(todosDir, t.FileName())
		}
		return t.FileName()
	}

	if asJSON {
		items := make([]todoListItemDTO, 0, len(entries))
		for _, t := range entries {
			items = append(items, todoListItemDTO{
				Date: t.Date(),
				Path: pathOf(t),
			})
		}
		return writeJSON(items)
	}

	type entry struct {
		date string
		path string
	}
	items := make([]entry, len(entries))
	maxWidth := 0
	for i, t := range entries {
		date := t.Date().Format(domain.FileNameDateLayoutTodo)
		items[i] = entry{date: date, path: pathOf(t)}
		if w := runewidth.StringWidth(date); w > maxWidth {
			maxWidth = w
		}
	}

	for _, item := range items {
		padding := maxWidth - runewidth.StringWidth(item.date)
		fmt.Fprintf(os.Stdout, "%s%s\t%s\n", item.date, strings.Repeat(" ", padding), item.path)
	}

	return nil
}
