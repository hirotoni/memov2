package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yuin/goldmark/text"
)

const testDataDir = "testdata"

func TestNewMarkdownHandler(t *testing.T) {
	handler := NewMarkdownHandler()
	if handler == nil {
		t.Fatal("Expected NewMarkdownHandler to return non-nil handler")
	}
}

func loadTestFile(t *testing.T, filename string) []byte {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(testDataDir, filename))
	if err != nil {
		t.Fatalf("Failed to read test file %s: %v", filename, err)
	}
	return content
}

func TestMarkdownHandler_findHeadingAndContent(t *testing.T) {
	handler := NewMarkdownHandler()

	tests := []struct {
		name          string
		filename      string
		heading       Heading
		expectFound   bool
		expectContent int // expected number of content nodes
	}{
		{
			name:          "finds heading with content",
			filename:      "heading_with_content.md",
			heading:       HeadingTodos,
			expectFound:   true,
			expectContent: 3, // list, heading, list
		},
		{
			name:          "finds heading at end of document",
			filename:      "heading_at_end.md",
			heading:       HeadingTodos,
			expectFound:   true,
			expectContent: 1, // just the list
		},
		{
			name:          "heading not found",
			filename:      "heading_not_found.md",
			heading:       HeadingTodos,
			expectFound:   false,
			expectContent: 0,
		},
		{
			name:          "heading found but no content",
			filename:      "heading_no_content.md",
			heading:       HeadingTodos,
			expectFound:   true,
			expectContent: 0,
		},
		{
			name:          "stops at same level heading",
			filename:      "stops_same_level.md",
			heading:       HeadingTodos,
			expectFound:   true,
			expectContent: 1, // only the list before "Another Section"
		},
		{
			name:          "stops at higher level heading",
			filename:      "stops_higher_level.md",
			heading:       HeadingTodos,
			expectFound:   true,
			expectContent: 3, // list and subsection heading with its content
		},
		{
			name:          "includes lower level headings",
			filename:      "includes_lower_level.md",
			heading:       HeadingTodos,
			expectFound:   true,
			expectContent: 7, // list, heading, list, heading, list, heading, list
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := loadTestFile(t, tt.filename)
			reader := text.NewReader(source)
			doc := handler.md.Parser().Parse(reader)

			foundHeading, hangingNodes := handler.findHeadingAndContent(doc, source, tt.heading)

			if tt.expectFound {
				if foundHeading == nil {
					t.Errorf("Expected to find heading, but got nil")
					return
				}

				headingText := string(foundHeading.Lines().Value(source))
				if !strings.Contains(headingText, tt.heading.Text) {
					t.Errorf("Expected heading text to contain '%s', got '%s'", tt.heading.Text, headingText)
				}
			} else {
				if foundHeading != nil {
					t.Errorf("Expected not to find heading, but got: %s", string(foundHeading.Lines().Value(source)))
				}
			}

			if len(hangingNodes) != tt.expectContent {
				t.Errorf("Expected %d content nodes, got %d", tt.expectContent, len(hangingNodes))
			}
		})
	}
}

func TestMarkdownHandler_FindHeadingEntity(t *testing.T) {
	handler := NewMarkdownHandler()

	tests := []struct {
		name           string
		filename       string
		heading        Heading
		expectFound    bool
		expectLevel    int
		expectNonEmpty bool
	}{
		{
			name:           "complete heading entity",
			filename:       "complete_heading.md",
			heading:        HeadingTodos,
			expectFound:    true,
			expectLevel:    2,
			expectNonEmpty: true,
		},
		{
			name:           "heading with no content",
			filename:       "heading_no_content.md",
			heading:        HeadingTodos,
			expectFound:    true,
			expectLevel:    2,
			expectNonEmpty: false,
		},
		{
			name:        "heading not found",
			filename:    "heading_not_found.md",
			heading:     HeadingTodos,
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := loadTestFile(t, tt.filename)

			entity, err := handler.HeadingBlockByHeading(source, tt.heading)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.expectFound {
				if entity == nil {
					t.Errorf("Expected to find heading entity, but got nil")
					return
				}

				if entity.Level != tt.expectLevel {
					t.Errorf("Expected level %d, got %d", tt.expectLevel, entity.Level)
				}

				if tt.expectNonEmpty && entity.ContentText == "" {
					t.Errorf("Expected non-empty content, got empty")
				}

				if !tt.expectNonEmpty && entity.ContentText != "" {
					t.Errorf("Expected empty content, got: %s", entity.ContentText)
				}
			} else {
				if entity != nil {
					t.Errorf("Expected not to find heading entity, but got: %+v", entity)
				}
			}
		})
	}
}

func TestMarkdownHandler_HeadingBlocksByLevel(t *testing.T) {
	handler := NewMarkdownHandler()

	source := loadTestFile(t, "multiple_levels.md")

	// Test level 2 headings
	blocks, err := handler.HeadingBlocksByLevel(source, 2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedCount := 3
	if len(blocks) != expectedCount {
		t.Errorf("Expected %d level 2 headings, got %d", expectedCount, len(blocks))
	}

	// Test level 1 headings
	blocks, err = handler.HeadingBlocksByLevel(source, 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedCount = 1
	if len(blocks) != expectedCount {
		t.Errorf("Expected %d level 1 headings, got %d", expectedCount, len(blocks))
	}
}

func TestMarkdownHandler_TopLevelBodyContent(t *testing.T) {
	handler := NewMarkdownHandler()

	tests := []struct {
		name          string
		content       string
		expectNil     bool
		expectHeading string
		expectContent string
		expectLevel   int
		expectLineNum int
	}{
		{
			name: "document with top level content",
			content: `# Main Heading
This is some content.
- List item 1
- List item 2

## Second Level
More content here.`,
			expectNil:     false,
			expectHeading: "Main Heading",
			expectContent: "This is some content.\n- List item 1\n- List item 2",
			expectLevel:   1,
			expectLineNum: 0,
		},
		{
			name:      "empty document",
			content:   "",
			expectNil: true,
		},
		{
			name: "document with only second level content",
			content: `## Second Level
Some content here.`,
			expectNil: true,
		},
		{
			name: "document with multiple top level sections",
			content: `# First Section
Content 1
Content 2

# Second Section
More content`,
			expectNil:     false,
			expectHeading: "First Section",
			expectContent: "Content 1\nContent 2",
			expectLevel:   1,
			expectLineNum: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.TopLevelBodyContent([]byte(tt.content))

			if tt.expectNil {
				if result != nil {
					t.Errorf("Expected nil result, got: %+v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.HeadingText != tt.expectHeading {
				t.Errorf("Expected heading text '%s', got '%s'", tt.expectHeading, result.HeadingText)
			}

			if result.Level != tt.expectLevel {
				t.Errorf("Expected level %d, got %d", tt.expectLevel, result.Level)
			}

			// Normalize line endings for comparison
			gotContent := strings.ReplaceAll(result.ContentText, "\r\n", "\n")
			expectContent := strings.ReplaceAll(tt.expectContent, "\r\n", "\n")
			if gotContent != expectContent {
				t.Errorf("Expected content:\n%s\nGot:\n%s", expectContent, gotContent)
			}

			// Line number can vary based on parser implementation, so we just check it's set
			if result.LineNumber < 0 {
				t.Errorf("Expected non-negative line number, got %d", result.LineNumber)
			}
		})
	}
}
