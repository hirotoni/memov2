package memo

import (
	"time"

	"github.com/hirotoni/memov2/internal/interfaces"
	memsearch "github.com/hirotoni/memov2/internal/search"
)

// This file defines the JSON DTOs used by the machine-readable (--json) output of
// the read-only memo commands. Domain types are intentionally NOT marshalled
// directly: these DTOs form a stable, documented contract for AI/automation
// consumers. Fields may be added in a backward-compatible way, but existing
// fields and json tags must not change meaning.

// listItemDTO is one memo in `memos list --json`.
type listItemDTO struct {
	Title    string    `json:"title"`
	Path     string    `json:"path"`
	Category []string  `json:"category"`
	Date     time.Time `json:"date"`
}

// categoryDTO is one category in `memos categories --json`.
type categoryDTO struct {
	Category []string `json:"category"`
	Path     string   `json:"path"`
}

// headingDTO is one heading block in `memos show --json`.
type headingDTO struct {
	Level   int    `json:"level"`
	Heading string `json:"heading"`
	Content string `json:"content"`
}

// showDTO is the full memo returned by `memos show --json`.
type showDTO struct {
	Title           string       `json:"title"`
	Path            string       `json:"path"`
	Category        []string     `json:"category"`
	Date            time.Time    `json:"date"`
	TopLevelContent string       `json:"top_level_content"`
	Headings        []headingDTO `json:"headings"`
}

// searchMatchDTO is one match within a searched memo in `memos search --json`.
type searchMatchDTO struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	Heading string `json:"heading"`
	Line    int    `json:"line"`
}

// searchItemDTO is one matched memo (with its aggregated matches) in
// `memos search --json`.
type searchItemDTO struct {
	Title    string           `json:"title"`
	Path     string           `json:"path"`
	Category []string         `json:"category"`
	Matches  []searchMatchDTO `json:"matches"`
}

// categoryTree normalizes a memo's category tree so it always marshals to a JSON
// array (never null) even when the memo is at the root.
func categoryTree(m interfaces.MemoFileInterface) []string {
	tree := m.CategoryTree()
	if tree == nil {
		return []string{}
	}
	return tree
}

// newShowDTO builds a showDTO from a parsed memo entry and its resolved path.
func newShowDTO(m interfaces.MemoFileInterface, path string) showDTO {
	headings := make([]headingDTO, 0, len(m.HeadingBlocks()))
	for _, hb := range m.HeadingBlocks() {
		headings = append(headings, headingDTO{
			Level:   hb.Level,
			Heading: hb.HeadingText,
			Content: hb.ContentText,
		})
	}

	var topLevel string
	if tl := m.TopLevelBodyContent(); tl != nil {
		topLevel = tl.ContentText
	}

	return showDTO{
		Title:           m.Title(),
		Path:            path,
		Category:        categoryTree(m),
		Date:            m.Date(),
		TopLevelContent: topLevel,
		Headings:        headings,
	}
}

// matchTypeString maps a search MatchType to its stable JSON string.
func matchTypeString(t memsearch.MatchType) string {
	switch t {
	case memsearch.MatchTitle:
		return "title"
	case memsearch.MatchCategory:
		return "category"
	case memsearch.MatchHeading:
		return "heading"
	case memsearch.MatchContent:
		return "content"
	default:
		return "unknown"
	}
}
