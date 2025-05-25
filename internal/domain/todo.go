package domain

import (
	"errors"
	"time"

	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/interfaces"
)

// TodoFileInterface is an alias for interfaces.TodoFileInterface to maintain backward compatibility
type TodoFileInterface = interfaces.TodoFileInterface
type TodoFile struct{ file }

const (
	FileNameDateLayoutTodo    = "20060102Mon"
	FileNameRegexTodo         = `^\d{8}\S{3}_todos\.md$`
	FileNameDateTimeRegexTodo = `^\d{8}\S{3}`
)

func NewTodosFile(date time.Time) (TodoFileInterface, error) {
	if date.IsZero() {
		return nil, errors.New("invalid date")
	}

	return &TodoFile{
		file: file{
			date:     date,
			fileType: FileTypeTodos,
			title:    date.Format(FileNameDateLayoutTodo),
		},
	}, nil
}

func NewTodoTemplateFile() (TodoFileInterface, error) {
	// set the current date but wont use it in the filename
	date := time.Now()
	f := &TodoFile{
		file: file{
			date:     date,
			fileType: FileTypeTemplate,
			title:    "todos_template",
		},
	}

	f.SetHeadingBlocks([]*markdown.HeadingBlock{
		{HeadingText: "todos", Level: 2},
		{HeadingText: "wanttodos", Level: 2},
	})

	return f, nil
}

func (f *TodoFile) FileName() string {
	if f.fileType == FileTypeTemplate {
		filename := "todos_template"
		return filename + FileExtension
	}

	datestring := f.date.Format(FileNameDateLayoutTodo)
	filename := datestring + FileSeparator + f.fileType.String()

	return filename + FileExtension
}

// ContentString overrides the base implementation to ensure trailing newline
func (f *TodoFile) ContentString() string {
	content := f.file.ContentString()
	// Ensure content ends with \n\n (two newlines) for todo files
	if len(f.HeadingBlocks()) > 0 {
		// If there are heading blocks, base ContentString() already adds one \n
		// We need to add one more to match golden files
		return content + "\n"
	}
	// If no heading blocks, base ContentString() ends with \n\n from title
	// But we still want \n\n at the end
	if content != "" && len(content) >= 2 && content[len(content)-2:] != "\n\n" {
		return content + "\n"
	}
	return content
}
