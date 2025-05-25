package interfaces

import (
	"time"

	domainmarkdown "github.com/hirotoni/memov2/internal/domain/markdown"
)

// FileType represents the type of file
type FileType string

const (
	FileTypeTodos    FileType = "todos"
	FileTypeMemo     FileType = "memo"
	FileTypeWeekly   FileType = "weekly"
	FileTypeTemplate FileType = "template"
)

func (ft FileType) String() string { return string(ft) }

// HeadingBlock is an alias for domain.markdown.HeadingBlock
// The concrete implementation is in the domain layer
type HeadingBlock = domainmarkdown.HeadingBlock

// FileInterface defines the interface for file operations
type FileInterface interface {
	// Meta data
	FileName() string
	Date() time.Time
	Title() string
	FileType() FileType

	// Meta data modification
	SetDate(date time.Time)
	SetTitle(title string)

	// Heading blocks
	TopLevelBodyContent() *HeadingBlock
	HeadingBlocks() []*HeadingBlock
	LastHeadingBlock() *HeadingBlock

	// Heading blocks modification
	SetTopLevelBodyContent(content *HeadingBlock)
	SetHeadingBlocks(hbs []*HeadingBlock)
	OverrideHeadingBlockMatched(input *HeadingBlock) error
	OverrideHeadingBlocksMatched(hbs []*HeadingBlock) error

	// Data operation for serialization
	ContentString() string
}

// MemoFileInterface defines the interface for memo file operations
type MemoFileInterface interface {
	FileInterface

	CategoryTree() []string
	SetCategoryTree(tree []string)
	Location() string
}

// TodoFileInterface defines the interface for todo file operations
type TodoFileInterface interface {
	FileInterface
}

// WeeklyFileInterface defines the interface for weekly file operations
type WeeklyFileInterface interface {
	FileInterface
}
