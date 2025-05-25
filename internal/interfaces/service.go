package interfaces

// Services is the main service interface that aggregates all service types
type Services interface {
	Memo() MemoService
	Todo() TodoService
	Config() ConfigService
}

// MemoService defines the interface for memo service operations
type MemoService interface {
	BuildWeeklyReportMemos() error
	GenerateMemoFile(title string, categoryTree []string) error
	ListCategories() error
	GenerateMemoIndex() error
	Browse() error
	List(showFullPath bool) error
	Open(path string) error
	Rename(path string, newTitle string) error
	TidyMemos() error

	// Interactive embedded-TUI commands (memos search/rename/new).
	SearchInteractive() error
	RenameInteractive() error
	NewInteractive() error
}

// TodoService defines the interface for todo service operations
type TodoService interface {
	GenerateTodoFile(truncate bool) error
	BuildWeeklyReportTodos() error
}

// ConfigService defines the interface for config service operations
type ConfigService interface {
	Show()
	Edit() error
}
