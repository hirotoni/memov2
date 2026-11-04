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
	CreateNonInteractive(title string, categoryTree []string, body string, asJSON bool) error
	ListCategories(asJSON bool) error
	GenerateMemoIndex() error
	Browse() error
	List(showFullPath, asJSON bool) error
	Show(path string, asJSON bool) error
	Search(query string, asJSON bool) error
	Open(path string) error
	Rename(path string, newTitle string, asJSON bool) error
	TidyMemos() error

	// Interactive embedded-TUI commands (memos search/rename/new).
	SearchInteractive() error
	RenameInteractive() error
	NewInteractive() error
}

// TodoService defines the interface for todo service operations
type TodoService interface {
	GenerateTodoFile(truncate bool) error
	GenerateTodoFileNonInteractive(truncate, asJSON bool) error
	BuildWeeklyReportTodos() error
	BuildWeeklyReportTodosNonInteractive(asJSON bool) error
	List(showFullPath, asJSON bool) error
	Show(path string, asJSON bool) error
	MigrateBacklogHeading(dryRun bool) error
}

// ConfigService defines the interface for config service operations
type ConfigService interface {
	Show(asJSON bool) error
	Edit() error
}
