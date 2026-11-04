package domain

import (
	"regexp"
	"strings"
)

// The functions in this file rewrite the raw text of an existing memo file in
// place. Rename, category move and duplicate must change only the H1 title
// line or the frontmatter `category` entry; everything else the user wrote
// (other frontmatter keys, sub-headings before the first `##`, formatting)
// has to survive byte for byte. Re-serializing the parsed MemoFile model
// would drop whatever the model does not hold, so these operate on lines.

var (
	atxH1Regex           = regexp.MustCompile(`^ {0,3}#(?:[ \t]|$)`)
	setextH1Underline    = regexp.MustCompile(`^ {0,3}=+[ \t]*$`)
	frontmatterCategory  = regexp.MustCompile(`^category[ \t]*:`)
	yamlContinuationLine = regexp.MustCompile(`^(?:[ \t]|-(?:[ \t]|$))`)
)

// splitLines splits s after each "\n" so that strings.Join(lines, "") == s.
// Each element keeps its terminator ("\n" or "\r\n"); the last element may be
// an unterminated line or "".
func splitLines(s string) []string {
	return strings.SplitAfter(s, "\n")
}

// lineEnding returns the line terminator used by line ("\r\n", "\n" or "").
func lineEnding(line string) string {
	switch {
	case strings.HasSuffix(line, "\r\n"):
		return "\r\n"
	case strings.HasSuffix(line, "\n"):
		return "\n"
	default:
		return ""
	}
}

func trimLineEnding(line string) string {
	return strings.TrimRight(line, "\r\n")
}

// frontmatterBounds returns the index of the closing "---" line when lines
// begin with a YAML frontmatter block, or -1 when there is none.
func frontmatterBounds(lines []string) int {
	if len(lines) == 0 || trimLineEnding(lines[0]) != "---" {
		return -1
	}
	for i := 1; i < len(lines); i++ {
		if trimLineEnding(lines[i]) == "---" {
			return i
		}
	}
	return -1
}

// CategoryLine returns the frontmatter line encoding the given category tree,
// in the same format MetadataString writes (without the line terminator).
func CategoryLine(categoryTree []string) string {
	m := &MemoFile{categoryTree: categoryTree}
	meta := m.MetadataString()
	for _, l := range strings.Split(meta, "\n") {
		if frontmatterCategory.MatchString(l) {
			return l
		}
	}
	return "category: []"
}

// RewriteMemoTitle returns content with its H1 title replaced by "# title".
// The H1 is the first non-blank line after the frontmatter (memo files must
// start with an H1 to be parseable); both ATX ("# x") and setext ("x\n===")
// forms are replaced. If no H1 is found there, "# title" is inserted right
// after the frontmatter. All other bytes are left untouched.
func RewriteMemoTitle(content string, title string) string {
	lines := splitLines(content)
	start := frontmatterBounds(lines) + 1

	idx := start
	for idx < len(lines) && lines[idx] != "" && strings.TrimSpace(lines[idx]) == "" {
		idx++
	}

	if idx < len(lines) {
		if atxH1Regex.MatchString(trimLineEnding(lines[idx])) {
			lines[idx] = "# " + title + lineEnding(lines[idx])
			return strings.Join(lines, "")
		}
		if idx+1 < len(lines) && strings.TrimSpace(lines[idx]) != "" &&
			setextH1Underline.MatchString(trimLineEnding(lines[idx+1])) {
			out := append([]string{}, lines[:idx]...)
			out = append(out, "# "+title+lineEnding(lines[idx+1]))
			out = append(out, lines[idx+2:]...)
			return strings.Join(out, "")
		}
	}

	// No H1 found: insert one right after the frontmatter.
	out := append([]string{}, lines[:start]...)
	out = append(out, "# "+title+"\n\n")
	out = append(out, lines[start:]...)
	return strings.Join(out, "")
}

// RewriteMemoCategory returns content with the frontmatter `category` entry
// replaced by the given category tree. Other frontmatter keys and the body
// are left untouched. A block-style value (indented or "- item" lines after
// "category:") is replaced as a whole. If the frontmatter has no category key
// it is added as the first key; if there is no frontmatter, one is prepended.
func RewriteMemoCategory(content string, categoryTree []string) string {
	lines := splitLines(content)
	end := frontmatterBounds(lines)
	newLine := CategoryLine(categoryTree)

	if end < 0 {
		return "---\n" + newLine + "\n---\n\n" + content
	}

	for i := 1; i < end; i++ {
		if !frontmatterCategory.MatchString(trimLineEnding(lines[i])) {
			continue
		}
		j := i + 1
		for j < end && yamlContinuationLine.MatchString(lines[j]) {
			j++
		}
		out := append([]string{}, lines[:i]...)
		out = append(out, newLine+lineEnding(lines[i]))
		out = append(out, lines[j:]...)
		return strings.Join(out, "")
	}

	out := append([]string{}, lines[:1]...)
	out = append(out, newLine+lineEnding(lines[0]))
	out = append(out, lines[1:]...)
	return strings.Join(out, "")
}
