package markdown

import "testing"

func TestHeadingBlockString(t *testing.T) {
	testCases := []struct {
		name     string
		heading  HeadingBlock
		expected string
	}{
		{
			name: "h1 with content",
			heading: HeadingBlock{
				HeadingText: "Title",
				Level:       1,
				ContentText: "Some content",
			},
			expected: "# Title\n\nSome content",
		},
		{
			name: "h2 with content",
			heading: HeadingBlock{
				HeadingText: "Subtitle",
				Level:       2,
				ContentText: "More content",
			},
			expected: "## Subtitle\n\nMore content",
		},
		{
			name: "h3 with empty content",
			heading: HeadingBlock{
				HeadingText: "Section",
				Level:       3,
				ContentText: "",
			},
			expected: "### Section\n\n",
		},
		{
			name: "h4 with multiline content",
			heading: HeadingBlock{
				HeadingText: "Subsection",
				Level:       4,
				ContentText: "Line 1\nLine 2",
			},
			expected: "#### Subsection\n\nLine 1\nLine 2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.heading.String()
			if result != tc.expected {
				t.Errorf("expected:\n%q\ngot:\n%q", tc.expected, result)
			}
		})
	}
}
