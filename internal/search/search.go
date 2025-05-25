package search

import (
	"strings"

	"github.com/hirotoni/memov2/internal/domain"
)

type MatchType int

const (
	MatchTitle MatchType = iota
	MatchCategory
	MatchContent
	MatchHeading
)

type SearchResult struct {
	Memo    domain.MemoFileInterface
	Matches []Match
}

type Match struct {
	Type            MatchType
	HeadingOrder    int // Order index in HeadingBlocks slice
	Line            int // Line number within heading block (for content matches)
	Content         string
	PrevLineContext string
	NextLineContext string
	Heading         string // If match is under a heading
}

func SearchMemo(memo domain.MemoFileInterface, query string, matchType MatchType) SearchResult {
	result := SearchResult{
		Memo:    memo,
		Matches: []Match{},
	}

	queryLower := strings.ToLower(strings.TrimSpace(query))

	// For title matches
	if matchType == MatchTitle && strings.Contains(strings.ToLower(memo.Title()), queryLower) {
		result.Matches = append(result.Matches, Match{
			Type:    MatchTitle,
			Content: memo.Title(),
		})
	}

	// For category matches
	if matchType == MatchCategory {
		for _, category := range memo.CategoryTree() {
			if strings.Contains(strings.ToLower(category), queryLower) {
				result.Matches = append(result.Matches, Match{
					Type:    MatchCategory,
					Content: strings.Join(memo.CategoryTree(), "/"),
				})
				break
			}
		}
	}

	// For content and heading matches
	if matchType == MatchContent || matchType == MatchHeading {
		// Top level body content
		lines := strings.Split(memo.TopLevelBodyContent().ContentText, "\n")
		for i, line := range lines {
			if strings.Contains(strings.ToLower(line), queryLower) {
				var prevLineContext string
				var nextLineContext string

				if i > 0 {
					prevLineContext = lines[i-1]
				}

				if i < len(lines)-1 {
					nextLineContext = lines[i+1]
				}

				result.Matches = append(result.Matches, Match{
					Type:            MatchContent,
					Content:         line,
					HeadingOrder:    -1,
					Line:            i + 1,
					PrevLineContext: prevLineContext,
					NextLineContext: nextLineContext,
					Heading:         memo.TopLevelBodyContent().HeadingText,
				})
			}
		}

		// Heading blocks
		for i, block := range memo.HeadingBlocks() {
			if matchType == MatchHeading && strings.Contains(strings.ToLower(block.HeadingText), queryLower) {
				result.Matches = append(result.Matches, Match{
					Type:         MatchHeading,
					Content:      block.HeadingText,
					HeadingOrder: i,
					Line:         i,
				})
			} else if matchType == MatchContent {
				// Split content into lines for line number tracking
				lines := strings.Split(block.ContentText, "\n")
				for j, line := range lines {
					if strings.Contains(strings.ToLower(line), queryLower) {
						// Get context lines
						var prevLineContext string
						var nextLineContext string

						// Add previous line if exists
						if j > 0 {
							prevLineContext = lines[j-1]
						}

						// Add next line if exists
						if j < len(lines)-1 {
							nextLineContext = lines[j+1]
						}

						result.Matches = append(result.Matches, Match{
							Type:            MatchContent,
							Content:         line,
							HeadingOrder:    i,
							Line:            j + 1,
							PrevLineContext: prevLineContext,
							NextLineContext: nextLineContext,
							Heading:         block.HeadingText,
						})
					}
				}
			}
		}
	}

	return result
}
