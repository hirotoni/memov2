package domain

import (
	"bufio"
	"errors"
	"time"
)

type fileInterface interface {
	FileName() string
	Date() time.Time
	FileType() string
	Title() string

	TopLevelBodyContent() *HeadingBlock
	HeadingBlocks() []*HeadingBlock
	LastHeadingBlock() *HeadingBlock

	SetDate(date time.Time)
	SetTopLevelBodyContent(content *HeadingBlock)
	SetHeadingBlocks(hbs []*HeadingBlock)
	OverrideHeadingBlockMatched(input *HeadingBlock) error
	OverrideHeadingBlocksMatched(hbs []*HeadingBlock) error

	// Save writes the markdown representation of the file to the given writer.
	// It writes, in order:
	//   1) Title
	//   2) Top level body content (if any)
	//   3) Heading blocks
	// Note: File types may override this to add file-type specific sections
	// such as metadata for memos.
	Save(w *bufio.Writer) error
}

type file struct {
	date                time.Time
	fileType            FileType
	title               string
	topLevelBodyContent *HeadingBlock
	headingBlocks       []*HeadingBlock
}

var (
	Sep    = "_"
	filler = "-"
	Ext    = ".md"
)

func (f *file) Date() time.Time {
	return f.date
}
func (f *file) FileType() string {
	return f.fileType.String()
}
func (f *file) Title() string {
	return f.title
}
func (f *file) TopLevelBodyContent() *HeadingBlock {
	if f.topLevelBodyContent == nil {
		return &HeadingBlock{}
	}
	return f.topLevelBodyContent
}
func (f *file) HeadingBlocks() []*HeadingBlock {
	if len(f.headingBlocks) == 0 {
		return []*HeadingBlock{}
	}
	return f.headingBlocks
}
func (f *file) LastHeadingBlock() *HeadingBlock {
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

func (f *file) SetTopLevelBodyContent(content *HeadingBlock) {
	f.topLevelBodyContent = content
}

func (f *file) SetHeadingBlocks(entities []*HeadingBlock) {
	f.headingBlocks = entities
}

func (f *file) OverrideHeadingBlockMatched(input *HeadingBlock) error {
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

func (f *file) OverrideHeadingBlocksMatched(entities []*HeadingBlock) error {
	for _, input := range entities {
		err := f.OverrideHeadingBlockMatched(input)
		if err != nil {
			return err
		}
	}
	return nil
}

// Save implements the common markdown serialization for files.
// File-type specific structs embedding file may override this to prepend
// or append additional sections (e.g., metadata for memos).
func (f *file) Save(w *bufio.Writer) error {
	if w == nil {
		return errors.New("writer is nil")
	}

	// 1) Title
	if _, err := w.WriteString("# " + f.Title() + "\n\n"); err != nil {
		return err
	}

	// 2) Top level body content (if any)
	if tl := f.TopLevelBodyContent(); tl != nil && tl.ContentText != "" {
		if _, err := w.WriteString(tl.ContentText); err != nil {
			return err
		}
	}

	// 3) Heading blocks
	for _, hb := range f.HeadingBlocks() {
		if _, err := w.WriteString(hb.String()); err != nil {
			return err
		}
	}

	return nil
}
