package picker

import (
	"path/filepath"
	"strings"

	"github.com/hirotoni/memov2/internal/interfaces"
	memsearch "github.com/hirotoni/memov2/internal/search"
)

// romajiMatcher returns a MatchFunc that expands each query word through the
// SKK romaji converter, so romaji input (e.g. "kaigi") matches Japanese text
// (e.g. "会議"). Every query word must match some variation.
func romajiMatcher(rc *memsearch.RomajiConverter) MatchFunc {
	return func(query string, item Item) bool {
		if query == "" {
			return true
		}
		text := strings.ToLower(item.FilterBy)
		for _, word := range strings.Fields(query) {
			variations := []string{word}
			if rc != nil {
				variations = rc.Convert(word)
			}
			found := false
			for _, v := range variations {
				if strings.Contains(text, strings.ToLower(v)) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
}

// SelectMemoForRename shows a romaji-aware memo picker, then prompts for a new
// title. It returns the selected memo's path (relative to the memos directory,
// the same format the rename command accepts) and the entered title. ok is
// false if the user cancelled.
func SelectMemoForRename(repo interfaces.MemoRepo) (relPath string, newTitle string, ok bool, err error) {
	entries, err := repo.MemoEntries()
	if err != nil {
		return "", "", false, err
	}

	items := make([]Item, 0, len(entries))
	for _, m := range entries {
		rel := filepath.ToSlash(filepath.Join(m.Location(), m.FileName()))
		items = append(items, Item{
			Display:   m.Title(),
			Secondary: rel,
			FilterBy:  m.Title() + " " + rel,
			Payload:   rel,
		})
	}

	// Romaji converter is best-effort; on failure fall back to plain substring.
	rc, _ := memsearch.NewRomajiConverter()

	res, err := Run(Config{
		Title:       "Rename memo — pick a file",
		Items:       items,
		Match:       romajiMatcher(rc),
		WithInput:   true,
		InputPrompt: "New title: ",
	})
	if err != nil {
		return "", "", false, err
	}
	if res.Cancelled || res.Item == nil {
		return "", "", false, nil
	}
	return res.Item.Payload.(string), res.InputText, true, nil
}
