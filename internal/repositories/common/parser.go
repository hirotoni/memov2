package common

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/utils"
)

// DateParserConfig holds configuration for date parsing
type DateParserConfig struct {
	DateTimeRegex string
	DateLayout    string
}

// ParseDateFromFilename extracts and parses a date from a filename using the provided regex and layout
func ParseDateFromFilename(filename string, config DateParserConfig) (time.Time, error) {
	dateReg, err := regexp.Compile(config.DateTimeRegex)
	if err != nil {
		return time.Time{}, common.Wrap(err, common.ErrorTypeRepository, "invalid date regex pattern")
	}

	datestring := dateReg.FindString(filename)
	if datestring == "" {
		return time.Time{}, common.New(common.ErrorTypeRepository, fmt.Sprintf("no date found in filename: %s", filename))
	}

	date, err := time.Parse(config.DateLayout, datestring)
	if err != nil {
		return time.Time{}, common.Wrap(err, common.ErrorTypeRepository, fmt.Sprintf("error parsing date from filename: %s", filename))
	}

	return date, nil
}

// ReadMarkdownFile reads a markdown file and returns its content as bytes
func ReadMarkdownFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, common.Wrap(err, common.ErrorTypeFileSystem, fmt.Sprintf("error reading file: %s", path))
	}
	return b, nil
}

// MarkdownParser provides markdown parsing functionality
type MarkdownParser struct {
	handler *utils.MarkdownHandler
}

// NewMarkdownParser creates a new markdown parser
func NewMarkdownParser() *MarkdownParser {
	return &MarkdownParser{
		handler: utils.NewMarkdownHandler(),
	}
}

// Metadata extracts metadata from markdown content
func (p *MarkdownParser) Metadata(content []byte) map[string]interface{} {
	return p.handler.Metadata(content)
}

// HeadingBlocksByLevel extracts heading blocks at the specified level
func (p *MarkdownParser) HeadingBlocksByLevel(content []byte, level int) ([]*markdown.HeadingBlock, error) {
	return p.handler.HeadingBlocksByLevel(content, level)
}

// TopLevelBodyContent extracts top-level body content
func (p *MarkdownParser) TopLevelBodyContent(content []byte) *markdown.HeadingBlock {
	return p.handler.TopLevelBodyContent(content)
}

