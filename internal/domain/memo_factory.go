package domain

import (
	"errors"
	"time"

	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/interfaces"
)

// MemoFileFromParsedData creates a MemoFile from already parsed data
// This allows Repository layer to parse and Domain layer to construct
// This function constructs the entity directly without using setters to improve encapsulation
func MemoFileFromParsedData(
	date time.Time,
	title string,
	category []string,
	topLevelBodyContent *markdown.HeadingBlock,
	headingBlocks []*markdown.HeadingBlock,
) (interfaces.MemoFileInterface, error) {
	if date.IsZero() {
		return nil, errors.New("invalid date")
	}

	// Construct the entity directly without using setters
	mf := &MemoFile{
		file: file{
			date:                date,
			fileType:            FileTypeMemo,
			title:               title,
			topLevelBodyContent: topLevelBodyContent,
			headingBlocks:       headingBlocks,
		},
		categoryTree: category,
	}

	return mf, nil
}

