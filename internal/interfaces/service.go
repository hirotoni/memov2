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
	GenerateMemoFile(title string) error
	GenerateMemoIndex() error
	Browse() error
	TidyMemos() error
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
