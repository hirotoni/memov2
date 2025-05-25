package domain

import (
	"errors"
	"time"

	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/interfaces"
)

// TodoFileFromParsedData creates a TodoFile from already parsed data
// This allows Repository layer to parse and Domain layer to construct
// This function constructs the entity directly without using setters to improve encapsulation
func TodoFileFromParsedData(
	date time.Time,
	headingBlocks []*markdown.HeadingBlock,
) (interfaces.TodoFileInterface, error) {
	if date.IsZero() {
		return nil, errors.New("invalid date")
	}

	// Construct the entity directly without using setters
	return &TodoFile{
		file: file{
			date:          date,
			fileType:      FileTypeTodos,
			title:         date.Format(FileNameDateLayoutTodo),
			headingBlocks: headingBlocks,
		},
	}, nil
}

