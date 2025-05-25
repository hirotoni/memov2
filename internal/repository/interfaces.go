package repository

import (
	"time"

	"github.com/hirotoni/memov2/internal/domain"
)

type Repositories interface {
	Memo() Memo
	Todo() Todo
	MemoWeekly() Weekly
	TodoWeekly() Weekly
}

type Memo interface {
	MemoEntries() ([]domain.MemoFileInterface, error)
	Metadata(file domain.MemoFileInterface) (map[string]interface{}, error)
	Save(file domain.MemoFileInterface, truncate bool) error
	TidyMemos() error
	Categories() ([][]string, error)
	Move(file domain.MemoFileInterface, newCategoryTree []string) error
}

type Todo interface {
	TodoEntries() ([]domain.TodoFileInterface, error)
	Save(file domain.TodoFileInterface, truncate bool) error
	TodosTemplate(date time.Time) (domain.TodoFileInterface, error)
	FindTodosFileByDate(date time.Time) (domain.TodoFileInterface, error)
}

type Weekly interface {
	Save(file domain.WeeklyFileInterface, truncate bool) error
}
