package interfaces

import "time"

// Repositories is the main repository interface that aggregates all repository types
type Repositories interface {
	Memo() MemoRepo
	Todo() TodoRepo
	MemoWeekly() WeeklyRepo
	TodoWeekly() WeeklyRepo
}

// MemoRepo defines the interface for memo repository operations
type MemoRepo interface {
	MemoEntries() ([]MemoFileInterface, error)
	Metadata(file MemoFileInterface) (map[string]interface{}, error)
	Save(file MemoFileInterface, truncate bool) error
	Categories() ([][]string, error)
	Move(file MemoFileInterface, newCategoryTree []string) error
	Delete(file MemoFileInterface) error
	Rename(file MemoFileInterface, newTitle string) error
	Duplicate(file MemoFileInterface) (MemoFileInterface, error)
	Memo(file MemoFileInterface) (MemoFileInterface, error)
}

// TodoRepo defines the interface for todo repository operations
type TodoRepo interface {
	TodoEntries() ([]TodoFileInterface, error)
	Save(file TodoFileInterface, truncate bool) error
	TodosTemplate(date time.Time) (TodoFileInterface, error)
	FindTodosFileByDate(date time.Time) (TodoFileInterface, error)
}

// WeeklyRepo defines the interface for weekly report repository operations
type WeeklyRepo interface {
	Save(file WeeklyFileInterface, truncate bool) error
}
