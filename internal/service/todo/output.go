package todo

import (
	"encoding/json"
	"os"
	"time"

	"github.com/hirotoni/memov2/internal/interfaces"
)

// This file defines the JSON DTOs used by the machine-readable (--json) output of
// the read-only todo commands AND the non-interactive write commands
// (`todos create`, `todos weekly`). Domain types are intentionally NOT marshalled
// directly: these DTOs form a stable, documented contract for AI/automation
// consumers. Fields may be added in a backward-compatible way, but existing
// fields and json tags must not change meaning.

// todoListItemDTO is one daily todo file in `todos list --json`.
type todoListItemDTO struct {
	Date time.Time `json:"date"`
	Path string    `json:"path"`
}

// todoHeadingDTO is one heading block in `todos show --json`.
type todoHeadingDTO struct {
	Level   int    `json:"level"`
	Heading string `json:"heading"`
	Content string `json:"content"`
}

// todoShowDTO is the full daily todo returned by `todos show --json`.
type todoShowDTO struct {
	Date            time.Time        `json:"date"`
	Path            string           `json:"path"`
	TopLevelContent string           `json:"top_level_content"`
	Headings        []todoHeadingDTO `json:"headings"`
}

// todoNewResultDTO is emitted by `todos create --json`. It reports where the daily
// todo file was written.
type todoNewResultDTO struct {
	Date time.Time `json:"date"`
	Path string    `json:"path"`
}

// todoWeeklyResultDTO is emitted by `todos weekly --json`. It reports where the
// weekly report file was written.
type todoWeeklyResultDTO struct {
	Path string `json:"path"`
}

// newTodoShowDTO builds a todoShowDTO from a parsed todo entry and its resolved
// (TodosDir-relative) path.
func newTodoShowDTO(t interfaces.TodoFileInterface, path string) todoShowDTO {
	headings := make([]todoHeadingDTO, 0, len(t.HeadingBlocks()))
	for _, hb := range t.HeadingBlocks() {
		headings = append(headings, todoHeadingDTO{
			Level:   hb.Level,
			Heading: hb.HeadingText,
			Content: hb.ContentText,
		})
	}

	var topLevel string
	if tl := t.TopLevelBodyContent(); tl != nil {
		topLevel = tl.ContentText
	}

	return todoShowDTO{
		Date:            t.Date(),
		Path:            path,
		TopLevelContent: topLevel,
		Headings:        headings,
	}
}

// writeJSON encodes v as a single JSON value to stdout with a trailing newline.
// It is the sole stdout writer for the machine-readable (--json) todo commands;
// logs and errors continue to go to stderr via slog.
func writeJSON(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
