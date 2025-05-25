package markdown

import "strings"

// HeadingBlock represents a markdown heading block
type HeadingBlock struct {
	Level       int
	HeadingText string
	ContentText string
	LineNumber  int // Line number of the heading in the original file
}

// String returns the string representation of the heading block
func (me HeadingBlock) String() string {
	if me.ContentText == "" {
		return strings.Repeat("#", me.Level) + " " + me.HeadingText + "\n\n"
	}

	return strings.Repeat("#", me.Level) + " " + me.HeadingText + "\n\n" + me.ContentText
}
