package interfaces

// Usecases is the main usecase interface that aggregates all usecase types
type Usecases interface {
	Memo() MemoUsecase
	Todo() TodoUsecase
	Config() ConfigUsecase
}

// MemoUsecase defines the interface for memo usecase operations
type MemoUsecase interface {
	BuildWeeklyReportMemos() error
	GenerateMemoFile(title string) error
	GenerateMemoIndex() error
	Browse() error
}

// TodoUsecase defines the interface for todo usecase operations
type TodoUsecase interface {
	GenerateTodoFile(truncate bool) error
	BuildWeeklyReportTodos() error
}

// ConfigUsecase defines the interface for config usecase operations
type ConfigUsecase interface {
	Show()
	Edit() error
}
