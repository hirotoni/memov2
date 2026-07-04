package memo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hirotoni/memov2/internal/common"
	memsearch "github.com/hirotoni/memov2/internal/search"
)

// Search runs a non-interactive query over all memos and reports the matches.
// Unlike the interactive TUI (which additionally applies SKK romaji expansion
// and multi-word AND semantics), this performs a single case-insensitive
// substring pass: each memo is searched across all four match types (title,
// category, heading, content) and produces at most one result entry with its
// matches aggregated. In text mode it prints "title<TAB>path" per matched memo
// (deduplicated); in json mode it writes a []searchItemDTO.
func (uc memo) Search(query string, asJSON bool) error {
	if strings.TrimSpace(query) == "" {
		return common.New(common.ErrorTypeValidation, "search query must not be empty")
	}

	entries, err := uc.repos.Memo().MemoEntries()
	if err != nil {
		return err
	}

	matchTypes := []memsearch.MatchType{
		memsearch.MatchTitle,
		memsearch.MatchCategory,
		memsearch.MatchHeading,
		memsearch.MatchContent,
	}

	items := make([]searchItemDTO, 0)
	for _, m := range entries {
		var matches []searchMatchDTO
		for _, mt := range matchTypes {
			res := memsearch.SearchMemo(m, query, mt)
			for _, match := range res.Matches {
				// match.Line is only a 1-based line number for content matches;
				// for other match types the underlying engine reuses the field
				// as an ordinal, so expose 0 there to avoid a misleading contract.
				line := 0
				if match.Type == memsearch.MatchContent {
					line = match.Line
				}
				matches = append(matches, searchMatchDTO{
					Type:    matchTypeString(match.Type),
					Text:    match.Content,
					Heading: match.Heading,
					Line:    line,
				})
			}
		}
		if len(matches) == 0 {
			continue
		}
		items = append(items, searchItemDTO{
			Title:    m.Title(),
			Path:     filepath.Join(m.Location(), m.FileName()),
			Category: categoryTree(m),
			Matches:  matches,
		})
	}

	if asJSON {
		return writeJSON(items)
	}

	for _, item := range items {
		fmt.Fprintf(os.Stdout, "%s\t%s\n", item.Title, item.Path)
	}
	return nil
}
