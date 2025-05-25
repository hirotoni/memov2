package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/hirotoni/memov2/internal/domain/markdown"
)

const (
	FileSeparator = "_"
	FileFiller    = "-"
	FileExtension = ".md"
)

type FileInterface interface {
	// Meta data
	FileName() string
	Date() time.Time
	Title() string
	FileType() FileType

	// Meta data modification
	SetDate(date time.Time)

	// Heading blocks
	TopLevelBodyContent() *markdown.HeadingBlock
	HeadingBlocks() []*markdown.HeadingBlock
	LastHeadingBlock() *markdown.HeadingBlock

	// Heading blocks modification
	SetTopLevelBodyContent(content *markdown.HeadingBlock)
	SetHeadingBlocks(hbs []*markdown.HeadingBlock)
	OverrideHeadingBlockMatched(input *markdown.HeadingBlock) error
	OverrideHeadingBlocksMatched(hbs []*markdown.HeadingBlock) error

	// Data operation for serialization
	ContentString() string
}

type file struct {
	date                time.Time
	fileType            FileType
	title               string
	topLevelBodyContent *markdown.HeadingBlock
	headingBlocks       []*markdown.HeadingBlock
}

func (f *file) Date() time.Time    { return f.date }
func (f *file) FileType() FileType { return f.fileType }
func (f *file) Title() string      { return f.title }

func (f *file) TopLevelBodyContent() *markdown.HeadingBlock {
	if f.topLevelBodyContent == nil {
		return &markdown.HeadingBlock{}
	}
	return f.topLevelBodyContent
}
func (f *file) HeadingBlocks() []*markdown.HeadingBlock {
	if len(f.headingBlocks) == 0 {
		return []*markdown.HeadingBlock{}
	}
	return f.headingBlocks
}
func (f *file) LastHeadingBlock() *markdown.HeadingBlock {
	if len(f.headingBlocks) == 0 {
		return nil
	}
	return f.headingBlocks[len(f.headingBlocks)-1]
}

func (f *file) SetDate(date time.Time) {
	if date.IsZero() {
		return
	}
	f.date = date
}

func (f *file) SetTopLevelBodyContent(content *markdown.HeadingBlock) {
	f.topLevelBodyContent = content
}

func (f *file) SetHeadingBlocks(entities []*markdown.HeadingBlock) {
	f.headingBlocks = entities
}

func (f *file) OverrideHeadingBlockMatched(input *markdown.HeadingBlock) error {
	found := false
	for i, e := range f.headingBlocks {
		if e.Level == input.Level && e.HeadingText == input.HeadingText {
			f.headingBlocks[i] = input
			found = true
			break
		}
	}

	if !found {
		return errors.New("target entity not found")
	}
	return nil
}

func (f *file) OverrideHeadingBlocksMatched(entities []*markdown.HeadingBlock) error {
	for _, input := range entities {
		err := f.OverrideHeadingBlockMatched(input)
		if err != nil {
			return err
		}
	}
	return nil
}

// Data implements the common markdown serialization for files.
// File-type specific structs embedding file may override to prepend/append sections via composition.
func (f *file) ContentString() string {
	var sb strings.Builder
	// 1) Title
	sb.WriteString("# " + f.Title() + "\n\n")

	// 2) Top level body content (if any)
	if tl := f.TopLevelBodyContent(); tl != nil && tl.ContentText != "" {
		sb.WriteString(tl.ContentText + "\n\n")
	}

	// 3) Heading blocks
	for _, hb := range f.HeadingBlocks() {
		sb.WriteString(hb.String())
	}
	return sb.String()
}
