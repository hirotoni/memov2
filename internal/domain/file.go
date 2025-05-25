package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/interfaces"
)

const (
	FileSeparator = "_"
	FileFiller    = "-"
	FileExtension = ".md"
)

// FileInterface is an alias for interfaces.FileInterface to maintain backward compatibility
type FileInterface = interfaces.FileInterface

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

func (f *file) SetTitle(title string) {
	f.title = title
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
		sb.WriteString(tl.ContentText + "\n")
	}

	// 3) Heading blocks
	hasTopLevelContent := f.TopLevelBodyContent() != nil && f.TopLevelBodyContent().ContentText != ""
	for i, hb := range f.HeadingBlocks() {
		// Add spacing before each heading block only if there was top-level content or it's not the first heading
		// But skip if the previous heading block was empty (ends with \n\n) to avoid extra newlines
		if hasTopLevelContent || i > 0 {
			// Check if previous heading block was empty
			shouldAddNewline := true
			if i > 0 {
				prevHb := f.HeadingBlocks()[i-1]
				// If previous heading block is empty, it ends with \n\n, so we don't need extra newline
				if prevHb.ContentText == "" {
					shouldAddNewline = false
				}
			}
			if shouldAddNewline {
				sb.WriteString("\n")
			}
		}

		headingStr := hb.String()
		// For the last heading block, don't add trailing newlines
		if i == len(f.HeadingBlocks())-1 {
			headingStr = strings.TrimRight(headingStr, "\n")
		}
		sb.WriteString(headingStr)
	}

	// Add final newline only if there are heading blocks
	if len(f.HeadingBlocks()) > 0 {
		sb.WriteString("\n")
	}
	return sb.String()
}
