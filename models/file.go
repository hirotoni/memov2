package models

import (
	"errors"
	"time"
)

type FileInterface interface {
	FileName() string
	Date() time.Time
	FileType() string
	Title() string
	HeadingBlocks() []*HeadingBlock
	LastHeadingBlock() *HeadingBlock

	SetDate(date time.Time)
	SetHeadingBlocks(hbs []*HeadingBlock)
	OverrideHeadingBlockMatched(input *HeadingBlock) error
	OverrideHeadingBlocksMatched(hbs []*HeadingBlock) error
}

type file struct {
	date     time.Time
	fileType FileType
	title    string
	entities []*HeadingBlock
}

func NewFile(date time.Time, fileType FileType, title string) (*file, error) {
	if date.IsZero() {
		return nil, errors.New("invalid date")
	}

	return &file{
		date:     date,
		fileType: fileType,
		title:    title,
	}, nil
}

func NewTodosFile(date time.Time) (FileInterface, error) {
	f, err := NewFile(date, FileTypeTodos, date.Format(FileNameDateLayoutTodo))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func NewWeeklyFile() (FileInterface, error) {
	date := time.Now()

	// set the current date but wont use it in the filename
	f, err := NewFile(date, FileTypeWeekly, "weekly_report")
	if err != nil {
		return nil, err
	}
	return f, nil
}

func NewTodoTemplateFile() (FileInterface, error) {
	date := time.Now()

	// set the current date but wont use it in the filename
	f, err := NewFile(date, FileTypeTemplate, "todos_template")
	if err != nil {
		return nil, err
	}

	f.SetHeadingBlocks([]*HeadingBlock{
		{HeadingText: "todos", Level: 2},
		{HeadingText: "wanttodos", Level: 2},
	})

	return f, nil
}

var (
	FileNameDateLayoutTodo   = "20060102Mon"
	FileNameDateLayoutMemo   = "20060102Mon150405"
	Sep                      = "_"
	filler                   = "-"
	Ext                      = ".md"
	MemoFileNameRegex        = `^\d{8}\S{3}\d{6}_memo_.*\.md$`
	MemoFileNameExtractRegex = `^\d{8}\S{3}\d{6}_memo_(.*)\.md$`
	DateTimeRegex            = `^\d{8}\S{3}\d{6}`
)

func (f *file) FileName() string {
	var filename string

	if f.fileType == FileTypeWeekly {
		filename = "weekly_report"
		return filename + Ext
	}

	if f.fileType == FileTypeTemplate {
		filename = "todos_template"
		return filename + Ext
	}

	if f.fileType == FileTypeTodos {
		datestring := f.date.Format(FileNameDateLayoutTodo)
		filename = datestring + Sep + f.fileType.String()
	}

	return filename + ".md"
}

func (f *file) Date() time.Time {
	return f.date
}
func (f *file) FileType() string {
	return f.fileType.String()
}
func (f *file) Title() string {
	return f.title
}
func (f *file) HeadingBlocks() []*HeadingBlock {
	return f.entities
}
func (f *file) LastHeadingBlock() *HeadingBlock {
	if len(f.entities) == 0 {
		return nil
	}
	return f.entities[len(f.entities)-1]
}

func (f *file) SetDate(date time.Time) {
	if date.IsZero() {
		return
	}
	f.date = date
}

func (f *file) SetHeadingBlocks(entities []*HeadingBlock) {
	f.entities = entities
}

func (f *file) OverrideHeadingBlockMatched(input *HeadingBlock) error {
	found := false
	for i, e := range f.entities {
		if e.Level == input.Level && e.HeadingText == input.HeadingText {
			f.entities[i] = input
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
