package usecase

import (
	"github.com/hirotoni/memov2/internal/domain"
)

// MemoUseCase defines the contract for memo business logic
type MemoUseCase interface {
	// BuildWeeklyReportMemos builds a weekly report of memos
	BuildWeeklyReportMemos() error
	// GenerateMemoFile generates a new memo file
	GenerateMemoFile(title string) error
	// GenerateMemoIndex generates an index of all memos
	GenerateMemoIndex() error
	// SearchMemos searches for memos by query
	SearchMemos(query string) ([]domain.MemoFileInterface, error)
	// GetMemoByTitle gets a memo by its title
	GetMemoByTitle(title string) (domain.MemoFileInterface, error)
}

// TodoUseCase defines the contract for todo business logic
type TodoUseCase interface {
	// GenerateTodoFile generates a new todo file
	GenerateTodoFile(truncate bool) error
	// BuildWeeklyReportTodos builds a weekly report of todos
	BuildWeeklyReportTodos() error
	// GetTodosByDate gets todos for a specific date
	GetTodosByDate(date string) ([]domain.TodoFileInterface, error)
	// MarkTodoComplete marks a todo as complete
	MarkTodoComplete(todoID string) error
}

// WeeklyUseCase defines the contract for weekly report business logic
type WeeklyUseCase interface {
	// GenerateWeeklyReport generates a weekly report
	GenerateWeeklyReport(week int, year int) error
	// GetWeeklyReport gets a weekly report by week and year
	GetWeeklyReport(week int, year int) (domain.WeeklyFileInterface, error)
	// ListWeeklyReports lists all weekly reports
	ListWeeklyReports() ([]domain.WeeklyFileInterface, error)
}
