package utils

import (
	"fmt"
	"strings"
)

type MarkdownBuilder struct {
	tabSize int
}

func NewMarkdownBuilder() *MarkdownBuilder {
	return &MarkdownBuilder{
		tabSize: 2,
	}
}

const (
	halfWidthChars = "#."
	fullWidthChars = "　！＠＃＄％＾＆＊（）＋｜〜＝￥｀「」｛｝；’：”、。・＜＞？【】『』《》〔〕［］‹›«»〘〙〚〛"
)

func (mb *MarkdownBuilder) text2tag(text string) string {
	if text == "" {
		return ""
	}

	var tag = text
	tag = mb.removeCharacters(tag, halfWidthChars)
	tag = mb.removeCharacters(tag, fullWidthChars)
	tag = strings.ReplaceAll(tag, " ", "-")

	return tag
}

func (mb *MarkdownBuilder) removeCharacters(text, charsToRemove string) string {
	for _, char := range charsToRemove {
		text = strings.ReplaceAll(text, string(char), "")
	}
	return text
}

// MARK: block

func (mb *MarkdownBuilder) BuildHeading(level int, text string) string {
	if level < 1 || level > 6 {
		return ""
	}
	return strings.Repeat("#", level) + " " + text + "\n"
}

func (mb *MarkdownBuilder) BuildList(item string, level int) string {
	if level < 1 {
		level = 1
	}

	var sb strings.Builder
	indent := strings.Repeat(" ", mb.tabSize*(level-1))
	sb.WriteString(indent + "- " + item + "\n")
	return sb.String()
}

func (mb *MarkdownBuilder) BuildOrderedList(order int, text string, level int, parentOrder int) string {
	if order < 1 {
		order = 1
	}
	if parentOrder < 1 {
		parentOrder = 1
	}

	var sb strings.Builder
	repeat := level - 1
	parentOrderStrLen := len(fmt.Sprintf("%d", parentOrder))

	// Adjust indent for ordered lists based on number of digits in the order
	indent := strings.Repeat(" ", (mb.tabSize+1+parentOrderStrLen-1)*repeat)
	sb.WriteString(indent + fmt.Sprintf("%d. ", order) + text + "\n")

	return sb.String()
}

func (mb *MarkdownBuilder) BuildCodeBlock(code string, language string) string {
	if code == "" {
		return ""
	}
	if language != "" {
		return "```" + language + "\n" + code + "\n```\n"
	}
	return "```\n" + code + "\n```\n"
}

// MARK: inline

func (mb *MarkdownBuilder) BuildLink(text string, url string, tag string) string {
	if text == "" || url == "" {
		return ""
	}

	if tag == "" {
		return "[" + text + "](" + url + ")"
	}

	tag = mb.text2tag(tag)
	return "[" + text + "](" + url + "#" + tag + ")"
}
