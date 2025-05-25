package domain

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	FileNameDateLayoutMemo    = "20060102Mon150405"
	FileNameDateTimeRegexMemo = `^\d{8}\S{3}\d{6}`
	FileNameRegexMemo         = `^\d{8}\S{3}\d{6}_memo_.*\.md$`
	FileNameExtractRegexMemo  = `^\d{8}\S{3}\d{6}_memo_(.*)\.md$`
)

type MemoFileInterface interface {
	FileInterface

	CategoryTree() []string
	SetCategoryTree(tree []string)
	Location() string
}

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
	title := strings.ReplaceAll(f.title, " ", FileFiller)
	filename = datetimestring + FileSeparator + f.fileType.String() + FileSeparator + title

	return filename + FileExtension
}

func (f *MemoFile) CategoryTree() []string        { return f.categoryTree }
func (f *MemoFile) SetCategoryTree(tree []string) { f.categoryTree = tree }

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

// ContentString returns memo file content including metadata, title, body, and headings.
// Order:
//  1. YAML frontmatter with category
//  2. Title
//  3. Top-level body content
//  4. Heading blocks
func (f *MemoFile) ContentString() string {
	meta := f.MetadataString()
	rest := f.file.ContentString()
	// Concatenate metadata and common body
	return meta + rest
}

func (f *MemoFile) MetadataString() string {
	const (
		header         = "---\n"
		categoryPrefix = "category: "
		footer         = "\n---\n\n"
	)

	sb := strings.Builder{}
	sb.WriteString(header)
	sb.WriteString(categoryPrefix)

	sb.WriteString("[")
	for i, v := range f.CategoryTree() {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%q", v))
	}
	sb.WriteString("]")

	sb.WriteString(footer)
	return sb.String()
}
