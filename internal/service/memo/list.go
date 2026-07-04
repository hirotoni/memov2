package memo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-runewidth"
)

func (uc memo) List(showFullPath, asJSON bool) error {
	entries, err := uc.repos.Memo().MemoEntries()
	if err != nil {
		return err
	}

	memosDir := uc.config.MemosDir()

	if asJSON {
		items := make([]listItemDTO, 0, len(entries))
		for _, m := range entries {
			var path string
			if showFullPath {
				path = filepath.Join(memosDir, m.Location(), m.FileName())
			} else {
				path = filepath.Join(m.Location(), m.FileName())
			}
			items = append(items, listItemDTO{
				Title:    m.Title(),
				Path:     path,
				Category: categoryTree(m),
				Date:     m.Date(),
			})
		}
		return writeJSON(items)
	}

	// Build titles and paths, and compute max title width
	type entry struct {
		title string
		path  string
	}
	items := make([]entry, len(entries))
	maxWidth := 0
	for i, m := range entries {
		title := m.Title()
		var path string
		if showFullPath {
			path = filepath.Join(memosDir, m.Location(), m.FileName())
		} else {
			path = filepath.Join(m.Location(), m.FileName())
		}
		items[i] = entry{title: title, path: path}
		if w := runewidth.StringWidth(title); w > maxWidth {
			maxWidth = w
		}
	}

	for _, item := range items {
		padding := maxWidth - runewidth.StringWidth(item.title)
		fmt.Fprintf(os.Stdout, "%s%s\t%s\n", item.title, strings.Repeat(" ", padding), item.path)
	}

	return nil
}

// writeJSON encodes v as a single JSON value to stdout with a trailing newline.
// It is the sole stdout writer for the machine-readable (--json) memo commands;
// logs and errors continue to go to stderr via slog.
func writeJSON(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
