package search

import (
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/stretchr/testify/assert"
)

func TestSearchMemo(t *testing.T) {
	// Create a test memo
	date := time.Now()
	memo, err := domain.NewMemoFile(date, "Test Memo", []string{"Category1", "Category2"})
	if err != nil {
		t.Fatalf("Failed to create test memo: %v", err)
	}

	// Add heading blocks with content
	memo.SetHeadingBlocks([]*markdown.HeadingBlock{
		{
			Level:       2,
			HeadingText: "Meeting Notes 会議メモ",
			ContentText: "Discussed project timeline\nNext steps: implementation\n予定を確認した",
		},
		{
			Level:       2,
			HeadingText: "Tasks",
			ContentText: "1. Review code\n2. Write tests\n3. Update docs\n",
		},
		{
			Level:       2,
			HeadingText: "References リファレンス",
			ContentText: "See documentation at:\nhttps://example.com\n参考資料を確認",
		},
	})

	// Test cases
	tests := []struct {
		name          string
		query         string
		matchType     MatchType
		wantMatch     bool
		wantMatchText string // For specific content verification
	}{
		// Title matches
		{
			name:      "Title match",
			query:     "Test",
			matchType: MatchTitle,
			wantMatch: true,
		},
		{
			name:      "Title no match",
			query:     "NonExistent",
			matchType: MatchTitle,
			wantMatch: false,
		},

		// Category matches
		{
			name:      "Category match",
			query:     "Category1",
			matchType: MatchCategory,
			wantMatch: true,
		},
		{
			name:      "Category no match",
			query:     "NonExistent",
			matchType: MatchCategory,
			wantMatch: false,
		},

		// Heading matches
		{
			name:          "Heading match English",
			query:         "Meeting",
			matchType:     MatchHeading,
			wantMatch:     true,
			wantMatchText: "Meeting Notes 会議メモ",
		},
		{
			name:          "Heading match Japanese",
			query:         "会議",
			matchType:     MatchHeading,
			wantMatch:     true,
			wantMatchText: "Meeting Notes 会議メモ",
		},
		{
			name:      "Heading no match",
			query:     "NonExistent",
			matchType: MatchHeading,
			wantMatch: false,
		},
		{
			name:          "Heading case insensitive",
			query:         "TASKS",
			matchType:     MatchHeading,
			wantMatch:     true,
			wantMatchText: "Tasks",
		},

		// Content matches
		{
			name:          "Content match English",
			query:         "timeline",
			matchType:     MatchContent,
			wantMatch:     true,
			wantMatchText: "Discussed project timeline",
		},
		{
			name:          "Content match Japanese",
			query:         "予定",
			matchType:     MatchContent,
			wantMatch:     true,
			wantMatchText: "予定を確認した",
		},
		{
			name:      "Content no match",
			query:     "NonExistent",
			matchType: MatchContent,
			wantMatch: false,
		},
		{
			name:          "Content case insensitive",
			query:         "REVIEW",
			matchType:     MatchContent,
			wantMatch:     true,
			wantMatchText: "1. Review code",
		},
		{
			name:          "Content with whitespace",
			query:         "  timeline  ",
			matchType:     MatchContent,
			wantMatch:     true,
			wantMatchText: "Discussed project timeline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchMemo(memo, tt.query, tt.matchType)

			if tt.wantMatch {
				assert.Greater(t, len(result.Matches), 0, "Expected matches but got none")
				if tt.wantMatchText != "" {
					found := false
					for _, match := range result.Matches {
						if match.Content == tt.wantMatchText {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected to find content %q but didn't", tt.wantMatchText)
				}
			} else {
				assert.Equal(t, 0, len(result.Matches), "Expected no matches but got some")
			}
		})
	}
}
