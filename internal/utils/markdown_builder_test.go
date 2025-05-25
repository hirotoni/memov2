package utils

import (
	"testing"
)

func TestNewMarkdownBuilder(t *testing.T) {
	mb := NewMarkdownBuilder()
	if mb == nil {
		t.Fatal("Expected non-nil MarkdownBuilder")
	}
	if mb.tabSize != 2 {
		t.Errorf("Expected tabSize to be 2, got %d", mb.tabSize)
	}
}

func TestText2Tag(t *testing.T) {
	mb := NewMarkdownBuilder()
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "Hello-World"},
		{"Test #1", "Test-1"},
		{"Full-width　chars！", "Full-widthchars"},
		{"No Change", "No-Change"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := mb.text2tag(tt.input)
			if result != tt.expected {
				t.Errorf("text2tag(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBuildHeading(t *testing.T) {
	mb := NewMarkdownBuilder()
	tests := []struct {
		level    int
		text     string
		expected string
	}{
		{1, "Heading 1", "# Heading 1\n"},
		{2, "Heading 2", "## Heading 2\n"},
		{6, "Heading 6", "###### Heading 6\n"},
		{0, "Invalid", ""},
		{7, "Invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := mb.BuildHeading(tt.level, tt.text)
			if result != tt.expected {
				t.Errorf("BuildHeading(%d, %q) = %q, want %q", tt.level, tt.text, result, tt.expected)
			}
		})
	}
}

func TestBuildList(t *testing.T) {
	mb := NewMarkdownBuilder()
	tests := []struct {
		item     string
		level    int
		expected string
	}{
		{"Item 1", 1, "- Item 1\n"},
		{"Nested item", 2, "  - Nested item\n"},
		{"Deeply nested", 3, "    - Deeply nested\n"},
	}

	for _, tt := range tests {
		t.Run(tt.item, func(t *testing.T) {
			result := mb.BuildList(tt.item, tt.level)
			if result != tt.expected {
				t.Errorf("BuildList(%q, %d) = %q, want %q", tt.item, tt.level, result, tt.expected)
			}
		})
	}
}

func TestBuildOrderedList(t *testing.T) {
	mb := NewMarkdownBuilder()
	tests := []struct {
		order    int
		text     string
		expected string
	}{
		{1, "First item", "1. First item\n"},
		{2, "Second item", "2. Second item\n"},
		{0, "Auto-numbered", "1. Auto-numbered\n"},
		{-1, "Negative", "1. Negative\n"},
		{10, "Tenth", "10. Tenth\n"},
		{1, "", "1. \n"},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := mb.BuildOrderedList(tt.order, tt.text, 1, 1)
			if result != tt.expected {
				t.Errorf("BuildOrderedList(%d, %q) = %q, want %q", tt.order, tt.text, result, tt.expected)
			}
		})
	}
}

func TestBuildCodeBlock(t *testing.T) {
	mb := NewMarkdownBuilder()
	tests := []struct {
		code     string
		language string
		expected string
	}{
		{"func main() {}", "go", "```go\nfunc main() {}\n```\n"},
		{"console.log('hello')", "js", "```js\nconsole.log('hello')\n```\n"},
		{"print('hello')", "", "```\nprint('hello')\n```\n"},
		{"", "go", ""},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result := mb.BuildCodeBlock(tt.code, tt.language)
			if result != tt.expected {
				t.Errorf("BuildCodeBlock(%q, %q) = %q, want %q", tt.code, tt.language, result, tt.expected)
			}
		})
	}
}

func TestBuildLink(t *testing.T) {
	mb := NewMarkdownBuilder()
	tests := []struct {
		text     string
		url      string
		tag      string
		expected string
	}{
		{"Link text", "https://example.com", "", "[Link text](https://example.com)"},
		{"Link text", "https://example.com", "Section", "[Link text](https://example.com#Section)"},
		{"Link text", "https://example.com", "Special # @", "[Link text](https://example.com#Special--@)"},
		{"", "https://example.com", "", ""},
		{"Link text", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := mb.BuildLink(tt.text, tt.url, tt.tag)
			if result != tt.expected {
				t.Errorf("BuildLink(%q, %q, %q) = %q, want %q", tt.text, tt.url, tt.tag, result, tt.expected)
			}
		})
	}
}
