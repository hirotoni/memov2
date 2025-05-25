package models

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type MemoFileInterface interface {
	FileInterface

	CategoryTree() []string
	SetCategoryTree(tree []string)
	Location() string
}

type memofile struct {
	file
	categoryTree []string // tree structure for memo files
}

func NewMemoFile(date time.Time, title string, categoryTree []string) (MemoFileInterface, error) {
	f, err := NewFile(date, FileTypeMemo, title)
	if err != nil {
		return nil, err
	}

	mf := &memofile{
		file:         *f,
		categoryTree: categoryTree,
	}

	return mf, nil
}

func (f *memofile) FileName() string {
	var filename string

	datetimestring := f.date.Format(FileNameDateLayoutMemo)
	title := strings.ReplaceAll(f.title, " ", filler)
	filename = datetimestring + Sep + f.fileType.String() + Sep + title

	return filename + Ext
}

func (f *memofile) CategoryTree() []string {
	return f.categoryTree
}

func (f *memofile) SetCategoryTree(tree []string) {
	f.categoryTree = tree
}

func (f *memofile) Location() string {
	if len(f.categoryTree) == 0 {
		return ""
	}

	return strings.Join(f.categoryTree, string(filepath.Separator))
}

func MemoTitle(filename string) string {
	re := regexp.MustCompile(MemoFileNameExtractRegex)
	matches := re.FindStringSubmatch(filename)

	if len(matches) > 1 {
		return matches[1]
	} else {
		return filename
	}
}
