package memo

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hirotoni/memov2/internal/search"
	"github.com/mattn/go-runewidth"
)

// matchedEntry holds a matched memo with its valid matches for context output
type matchedEntry struct {
	title   string
	path    string
	matches []search.Match
}

func (uc memo) Search(query string, showFullPath bool, showContext bool) error {
	rc, err := search.NewRomajiConverter()
	if err != nil {
		return fmt.Errorf("failed to load SKK dictionary: %w", err)
	}

	entries, err := uc.repos.Memo().MemoEntries()
	if err != nil {
		return err
	}

	memosDir := uc.config.MemosDir()

	// Split query into words and convert each word
	queryWords := strings.Fields(query)
	wordQueries := make([][]string, len(queryWords))
	for i, word := range queryWords {
		wordQueries[i] = rc.Convert(word)
	}

	// containsAllWords checks if all query words appear in a text
	containsAllWords := func(text string) bool {
		textLower := strings.ToLower(text)
		for _, wordVariations := range wordQueries {
			wordFound := false
			for _, word := range wordVariations {
				if strings.Contains(textLower, strings.ToLower(word)) {
					wordFound = true
					break
				}
			}
			if !wordFound {
				return false
			}
		}
		return true
	}

	// Search all memos
	var matched []matchedEntry

	for _, m := range entries {
		// Track matches by type
		matchesByType := make(map[search.MatchType][]search.Match)

		for _, wordVariations := range wordQueries {
			for _, q := range wordVariations {
				if matches := search.SearchMemo(m, q, search.MatchTitle); len(matches.Matches) > 0 {
					matchesByType[search.MatchTitle] = append(matchesByType[search.MatchTitle], matches.Matches...)
				}
				if matches := search.SearchMemo(m, q, search.MatchCategory); len(matches.Matches) > 0 {
					matchesByType[search.MatchCategory] = append(matchesByType[search.MatchCategory], matches.Matches...)
				}
				if matches := search.SearchMemo(m, q, search.MatchHeading); len(matches.Matches) > 0 {
					matchesByType[search.MatchHeading] = append(matchesByType[search.MatchHeading], matches.Matches...)
				}
				if matches := search.SearchMemo(m, q, search.MatchContent); len(matches.Matches) > 0 {
					matchesByType[search.MatchContent] = append(matchesByType[search.MatchContent], matches.Matches...)
				}
			}
		}

		// Check each match type and collect valid matches
		var validMatches []search.Match
		hasValidMatch := false

		for _, matches := range matchesByType {
			uniqueMatches := make(map[string]search.Match)
			for _, match := range matches {
				uniqueMatches[match.Content] = match
			}
			for _, match := range uniqueMatches {
				if containsAllWords(match.Content) {
					validMatches = append(validMatches, match)
					hasValidMatch = true
				}
			}
		}

		if hasValidMatch {
			title := m.Title()
			var path string
			if showFullPath {
				path = filepath.Join(memosDir, m.Location(), m.FileName())
			} else {
				path = filepath.Join(m.Location(), m.FileName())
			}
			matched = append(matched, matchedEntry{title: title, path: path, matches: validMatches})
		}
	}

	// Compute max title width for alignment
	maxWidth := 0
	for _, item := range matched {
		if w := runewidth.StringWidth(item.title); w > maxWidth {
			maxWidth = w
		}
	}

	// Sort by path (newest first)
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].path > matched[j].path
	})

	if showContext {
		for _, item := range matched {
			padding := maxWidth - runewidth.StringWidth(item.title)
			titlePadded := item.title + strings.Repeat(" ", padding)
			for _, m := range item.matches {
				label := matchTypeLabel(m.Type)
				content := formatMatchContent(m)
				fmt.Fprintf(os.Stdout, "%s\t%s\t%s\t%s\n", titlePadded, item.path, label, content)
			}
		}
	} else {
		for _, item := range matched {
			padding := maxWidth - runewidth.StringWidth(item.title)
			fmt.Fprintf(os.Stdout, "%s%s\t%s\n", item.title, strings.Repeat(" ", padding), item.path)
		}
	}

	return nil
}

func matchTypeLabel(t search.MatchType) string {
	switch t {
	case search.MatchTitle:
		return "[Title]"
	case search.MatchCategory:
		return "[Category]"
	case search.MatchHeading:
		return "[Heading]"
	case search.MatchContent:
		return "[Content]"
	default:
		return "[Unknown]"
	}
}

func formatMatchContent(m search.Match) string {
	content := strings.ReplaceAll(m.Content, "\t", " ")
	if m.Type == search.MatchContent && m.Heading != "" {
		heading := strings.ReplaceAll(m.Heading, "\t", " ")
		return heading + " > " + content
	}
	return content
}
