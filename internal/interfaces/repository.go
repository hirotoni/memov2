package interfaces

import (
	"time"

	"github.com/hirotoni/memov2/internal/domain"
)

// Repositories is the main repository interface that aggregates all repository types
type Repositories interface {
	Memo() MemoRepo
	Todo() TodoRepo
	MemoWeekly() WeeklyRepo
	TodoWeekly() WeeklyRepo
}

// MemoRepo defines the interface for memo repository operations
type MemoRepo interface {
	MemoEntries() ([]domain.MemoFileInterface, error)
	Metadata(file domain.MemoFileInterface) (map[string]interface{}, error)
	Save(file domain.MemoFileInterface, truncate bool) error
	TidyMemos() error
	Categories() ([][]string, error)
	Move(file domain.MemoFileInterface, newCategoryTree []string) error
	Delete(file domain.MemoFileInterface) error
	Rename(file domain.MemoFileInterface, newTitle string) error
	Duplicate(file domain.MemoFileInterface) (domain.MemoFileInterface, error)
	Memo(file domain.MemoFileInterface) (domain.MemoFileInterface, error)
}

// TodoRepo defines the interface for todo repository operations
type TodoRepo interface {
	TodoEntries() ([]domain.TodoFileInterface, error)
	Save(file domain.TodoFileInterface, truncate bool) error
	TodosTemplate(date time.Time) (domain.TodoFileInterface, error)
	FindTodosFileByDate(date time.Time) (domain.TodoFileInterface, error)
}

// WeeklyRepo defines the interface for weekly report repository operations
type WeeklyRepo interface {
	Save(file domain.WeeklyFileInterface, truncate bool) error
}
