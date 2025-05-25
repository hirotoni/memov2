package domain

import (
	"bufio"
	"errors"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type MemoFileInterface interface {
	fileInterface

	CategoryTree() []string
	SetCategoryTree(tree []string)
	Location() string
}

var (
	FileNameDateLayoutMemo    = "20060102Mon150405"
	FileNameDateTimeRegexMemo = `^\d{8}\S{3}\d{6}`
	FileNameRegexMemo         = `^\d{8}\S{3}\d{6}_memo_.*\.md$`
	FileNameExtractRegexMemo  = `^\d{8}\S{3}\d{6}_memo_(.*)\.md$`
)

type MemoFile struct {
	file
	categoryTree []string // tree structure for memo files
}

func NewMemoFile(date time.Time, title string, categoryTree []string) (MemoFileInterface, error) {
	if date.IsZero() {
		return nil, errors.New("invalid date")
	}

	mf := &MemoFile{
		file: file{
			date:     date,
			fileType: FileTypeMemo,
			title:    title,
		},
		categoryTree: categoryTree,
	}

	return mf, nil
}

func (f *MemoFile) FileName() string {
	var filename string

	datetimestring := f.date.Format(FileNameDateLayoutMemo)
	title := strings.ReplaceAll(f.title, " ", filler)
	filename = datetimestring + Sep + f.fileType.String() + Sep + title

	return filename + Ext
}

func (f *MemoFile) CategoryTree() []string {
	return f.categoryTree
}

func (f *MemoFile) SetCategoryTree(tree []string) {
	f.categoryTree = tree
}

func (f *MemoFile) Location() string {
	if len(f.categoryTree) == 0 {
		return ""
	}

	return strings.Join(f.categoryTree, string(filepath.Separator))
}

func MemoTitle(filename string) string {
	re := regexp.MustCompile(FileNameExtractRegexMemo)
	matches := re.FindStringSubmatch(filename)

	if len(matches) > 1 {
		return matches[1]
	} else {
		return filename
	}
}

// Save writes memo file content including metadata, title, body, and headings.
// Order:
//  1. YAML frontmatter with category
//  2. Title
//  3. Top-level body content
//  4. Heading blocks
func (f *MemoFile) Save(w *bufio.Writer) error {
	if w == nil {
		return errors.New("writer is nil")
	}

	// 1) Metadata
	if _, err := w.WriteString(f.metadataString()); err != nil {
		return err
	}

	// Delegate the rest to the common serializer (title, body, headings)
	return f.file.Save(w)
}

func (f *MemoFile) metadataString() string {
	wrap := func(s string) string { return "\"" + s + "\"" }

	sb := strings.Builder{}
	sb.WriteString("---\n")
	sb.WriteString("category: ")
	if len(f.CategoryTree()) > 0 {
		sb.WriteString("[")
		for i, v := range f.CategoryTree() {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(wrap(v))
		}
		sb.WriteString("]")
	} else {
		sb.WriteString("[]")
	}
	sb.WriteString("\n")
	sb.WriteString("---\n\n")

	return sb.String()
}
