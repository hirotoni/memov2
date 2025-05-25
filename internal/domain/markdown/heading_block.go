package markdown

import "strings"

// HeadingBlock represents a structural element in markdown with a heading and content
type HeadingBlock struct {
	HeadingText string // Raw text of heading
	Level       int    // Heading level (h1, h2, etc.)
	ContentText string // Raw text of content
	LineNumber  int    // Line number of the heading in the original file
}

func (me *HeadingBlock) String() string {
	if me.ContentText == "" {
		return strings.Repeat("#", me.Level) + " " + me.HeadingText + "\n\n"
	}

	return strings.Repeat("#", me.Level) + " " + me.HeadingText + "\n\n" + me.ContentText
}
