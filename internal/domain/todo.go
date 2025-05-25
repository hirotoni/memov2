package domain

import (
	"errors"
	"time"
)

type TodoFileInterface interface {
	fileInterface
}

type TodoFile struct {
	file
}

var (
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

	f.SetHeadingBlocks([]*HeadingBlock{
		{HeadingText: "todos", Level: 2},
		{HeadingText: "wanttodos", Level: 2},
	})

	return f, nil
}

func (f *TodoFile) FileName() string {
	if f.fileType == FileTypeTemplate {
		filename := "todos_template"
		return filename + Ext
	}

	datestring := f.date.Format(FileNameDateLayoutTodo)
	filename := datestring + Sep + f.fileType.String()

	return filename + Ext
}
