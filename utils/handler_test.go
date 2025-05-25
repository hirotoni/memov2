package utils

import (
	"strings"
	"testing"
)

func TestNewMarkdownHandler(t *testing.T) {
	handler := NewMarkdownHandler()
	if handler == nil {
		t.Fatal("Expected NewMarkdownHandler to return non-nil handler")
	}
}
func TestFindHeadingEntity(t *testing.T) {
	tests := []struct {
		name           string
		markdown       string
		heading        Heading
		wantFound      bool
		wantLevel      int
		wantHeadingStr string
	}{
		{
			name:           "Find H1 heading",
			markdown:       "# Heading 1\nContent under heading 1\n\n## Subheading\n\nMore content",
			heading:        Heading{Level: 1, Text: "Heading 1"},
			wantFound:      true,
			wantLevel:      1,
			wantHeadingStr: "Heading 1",
		},
		{
			name:           "Find H2 heading",
			markdown:       "# Main heading\n## Second heading\nContent under heading 2\n### Subheading\nMore content",
			heading:        Heading{Level: 2, Text: "Second"},
			wantFound:      true,
			wantLevel:      2,
			wantHeadingStr: "Second heading",
		},
		{
			name:           "Heading not found",
			markdown:       "# Main heading\n## Second heading\nContent under heading 2",
			heading:        Heading{Level: 2, Text: "Nonexistent"},
			wantFound:      false,
			wantLevel:      0,
			wantHeadingStr: "",
		},
		{
			name:           "Partial text match",
			markdown:       "# Main heading\n## Second heading with extra text\nContent",
			heading:        Heading{Level: 2, Text: "Second heading"},
			wantFound:      true,
			wantLevel:      2,
			wantHeadingStr: "Second heading with extra text",
		},
		{
			name:           "Multiple headings same level",
			markdown:       "# H1\n## First H2\nContent\n## Second H2\nMore content",
			heading:        Heading{Level: 2, Text: "Second"},
			wantFound:      true,
			wantLevel:      2,
			wantHeadingStr: "Second H2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewMarkdownHandler()
			result, err := handler.FindHeadingEntity([]byte(tt.markdown), tt.heading)
			if err != nil {
				t.Fatalf("FindHeading() error = %v", err)
			}

			if tt.wantFound && result == nil {
				t.Fatal("Expected to find heading but got nil")
			}
			if !tt.wantFound && result != nil {
				t.Fatal("Expected not to find heading but got result")
			}
			if !tt.wantFound {
				return
			}

			// Check heading level
			if result.Level != tt.wantLevel {
				t.Errorf("Expected heading level %d, got %d", tt.wantLevel, result.Level)
			}

			// Check heading text contains expected text
			headingStr := string(result.HeadingText)
			if !strings.Contains(headingStr, tt.wantHeadingStr) {
				t.Errorf("Expected heading text to contain %q, got %q", tt.wantHeadingStr, headingStr)
			}

			// Verify content is not empty when a heading is found
			if len(result.ContentText) == 0 {
				t.Error("Expected content to be non-empty")
			}
		})
	}
}
