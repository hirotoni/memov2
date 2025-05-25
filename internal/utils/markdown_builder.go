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

func (mb *MarkdownBuilder) BuildParagraph(text string) string {
	if text == "" {
		return ""
	}
	return text + "\n\n"
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

func (mb *MarkdownBuilder) BuildImage(altText string, url string) string {
	if altText == "" || url == "" {
		return ""
	}
	return "![ " + altText + "](" + url + ")\n"
}

func (mb *MarkdownBuilder) BuildBlockquote(text string) string {
	if text == "" {
		return ""
	}
	lines := strings.Split(text, "\n")
	var sb strings.Builder
	for _, line := range lines {
		sb.WriteString("> " + line + "\n")
	}
	return sb.String() + "\n"
}

func (mb *MarkdownBuilder) BuildHorizontalRule() string {
	return "---\n\n"
}

const (
	tableCellSeparator = " | "
	tableRowPrefix     = "| "
	tableRowSuffix     = " |\n"
	tableSeparatorCell = "---"
)

func (mb *MarkdownBuilder) BuildTable(headers []string, rows [][]string) string {
	if len(headers) == 0 || len(rows) == 0 {
		return ""
	}
	var sb strings.Builder
	mb.writeTableHeader(&sb, headers)
	mb.writeTableSeparator(&sb, len(headers))
	mb.writeTableRows(&sb, rows)

	return sb.String() + "\n"
}

func (mb *MarkdownBuilder) writeTableHeader(sb *strings.Builder, headers []string) {
	sb.WriteString(tableRowPrefix)
	sb.WriteString(strings.Join(headers, tableCellSeparator))
	sb.WriteString(tableRowSuffix)
}

func (mb *MarkdownBuilder) writeTableSeparator(sb *strings.Builder, headerCount int) {
	sb.WriteString(tableRowPrefix)
	sb.WriteString(strings.Repeat(tableSeparatorCell+tableCellSeparator, headerCount-1))
	sb.WriteString(tableSeparatorCell)
	sb.WriteString(tableRowSuffix)
}

func (mb *MarkdownBuilder) writeTableRows(sb *strings.Builder, rows [][]string) {
	for _, row := range rows {
		sb.WriteString(tableRowPrefix)
		sb.WriteString(strings.Join(row, tableCellSeparator))
		sb.WriteString(tableRowSuffix)
	}
}

func (mb *MarkdownBuilder) BuildTaskList(items []string, completed []bool) string {
	if len(items) == 0 || len(completed) != len(items) {
		return ""
	}
	var sb strings.Builder
	for i, item := range items {
		if completed[i] {
			sb.WriteString("- [x] " + item + "\n")
		} else {
			sb.WriteString("- [ ] " + item + "\n")
		}
	}
	return sb.String() + "\n"
}

func (mb *MarkdownBuilder) BuildDefinitionList(definitions map[string]string) string {
	if len(definitions) == 0 {
		return ""
	}
	var sb strings.Builder
	for term, definition := range definitions {
		sb.WriteString("**" + term + "**: " + definition + "\n")
	}
	return sb.String() + "\n"
}

func (mb *MarkdownBuilder) BuildFootnote(text string, footnote string) string {
	if text == "" || footnote == "" {
		return ""
	}
	return text + " [^1]\n\n[^1]: " + footnote + "\n"
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
