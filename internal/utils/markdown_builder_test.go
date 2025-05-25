package utils

import (
	"strings"
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

func TestBuildParagraph(t *testing.T) {
	mb := NewMarkdownBuilder()
	tests := []struct {
		text     string
		expected string
	}{
		{"This is a paragraph", "This is a paragraph\n\n"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := mb.BuildParagraph(tt.text)
			if result != tt.expected {
				t.Errorf("BuildParagraph(%q) = %q, want %q", tt.text, result, tt.expected)
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

func TestBuildImage(t *testing.T) {
	mb := NewMarkdownBuilder()
	tests := []struct {
		altText  string
		url      string
		expected string
	}{
		{"Image", "https://example.com/image.jpg", "![ Image](https://example.com/image.jpg)\n"},
		{"", "https://example.com/image.jpg", ""},
		{"Image", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.altText, func(t *testing.T) {
			result := mb.BuildImage(tt.altText, tt.url)
			if result != tt.expected {
				t.Errorf("BuildImage(%q, %q) = %q, want %q", tt.altText, tt.url, result, tt.expected)
			}
		})
	}
}

func TestBuildBlockquote(t *testing.T) {
	mb := NewMarkdownBuilder()
	tests := []struct {
		text     string
		expected string
	}{
		{"Single line", "> Single line\n\n"},
		{"Line 1\nLine 2", "> Line 1\n> Line 2\n\n"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := mb.BuildBlockquote(tt.text)
			if result != tt.expected {
				t.Errorf("BuildBlockquote(%q) = %q, want %q", tt.text, result, tt.expected)
			}
		})
	}
}

func TestBuildHorizontalRule(t *testing.T) {
	mb := NewMarkdownBuilder()
	expected := "---\n\n"
	result := mb.BuildHorizontalRule()
	if result != expected {
		t.Errorf("BuildHorizontalRule() = %q, want %q", result, expected)
	}
}

func TestBuildTable(t *testing.T) {
	mb := NewMarkdownBuilder()
	tests := []struct {
		headers  []string
		rows     [][]string
		expected string
	}{
		{
			[]string{"Name", "Age"},
			[][]string{{"Alice", "30"}, {"Bob", "25"}},
			"| Name | Age |\n| --- | --- |\n| Alice | 30 |\n| Bob | 25 |\n\n",
		},
		{[]string{}, [][]string{{"Alice", "30"}}, ""},
		{[]string{"Name"}, [][]string{}, ""},
	}

	for i, tt := range tests {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			result := mb.BuildTable(tt.headers, tt.rows)
			if result != tt.expected {
				t.Errorf("BuildTable(%v, %v) = %q, want %q", tt.headers, tt.rows, result, tt.expected)
			}
		})
	}
}

func TestBuildTaskList(t *testing.T) {
	mb := NewMarkdownBuilder()
	tests := []struct {
		items     []string
		completed []bool
		expected  string
	}{
		{
			[]string{"Task 1", "Task 2", "Task 3"},
			[]bool{true, false, true},
			"- [x] Task 1\n- [ ] Task 2\n- [x] Task 3\n\n",
		},
		{[]string{}, []bool{}, ""},
		{[]string{"Task"}, []bool{}, ""},
	}

	for i, tt := range tests {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			result := mb.BuildTaskList(tt.items, tt.completed)
			if result != tt.expected {
				t.Errorf("BuildTaskList(%v, %v) = %q, want %q", tt.items, tt.completed, result, tt.expected)
			}
		})
	}
}

func TestBuildDefinitionList(t *testing.T) {
	mb := NewMarkdownBuilder()
	tests := []struct {
		definitions map[string]string
		expected    string
	}{
		{
			map[string]string{"Term 1": "Definition 1", "Term 2": "Definition 2"},
			"**Term 1**: Definition 1\n**Term 2**: Definition 2\n\n",
		},
		{map[string]string{}, ""},
	}

	for i, tt := range tests {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			result := mb.BuildDefinitionList(tt.definitions)

			// Since map iteration order is not guaranteed, we need to check if all terms and definitions are present
			if result == "" && tt.expected == "" {
				return
			}

			for term, definition := range tt.definitions {
				expected := "**" + term + "**: " + definition + "\n"
				if !strings.Contains(result, expected) {
					t.Errorf("BuildDefinitionList result doesn't contain %q", expected)
				}
			}

			if len(result) != len(tt.expected) {
				t.Errorf("BuildDefinitionList result length %d doesn't match expected length %d", len(result), len(tt.expected))
			}
		})
	}
}

func TestBuildFootnote(t *testing.T) {
	mb := NewMarkdownBuilder()
	tests := []struct {
		text     string
		footnote string
		expected string
	}{
		{"Main text", "Footnote text", "Main text [^1]\n\n[^1]: Footnote text\n"},
		{"", "Footnote", ""},
		{"Text", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := mb.BuildFootnote(tt.text, tt.footnote)
			if result != tt.expected {
				t.Errorf("BuildFootnote(%q, %q) = %q, want %q", tt.text, tt.footnote, result, tt.expected)
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
