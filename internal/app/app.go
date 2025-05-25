package app

import (
	"fmt"
	"log/slog"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/repository"
	"github.com/hirotoni/memov2/internal/usecase/memo"
	"github.com/hirotoni/memov2/internal/usecase/todo"
)

// App represents the main application container
type App struct {
	config      *config.TomlConfig
	repos       repository.Repositories
	memoUseCase memo.Memo
	todoUseCase todo.Todo
	logger      *slog.Logger
}

// NewApp creates a new application instance with all dependencies
func NewApp(cfg *config.TomlConfig, logger *slog.Logger) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	// Initialize repositories
	repos := repository.NewRepositories(*cfg)

	// Initialize use cases
	memoUseCase := memo.NewMemo(*cfg, repos)
	todoUseCase := todo.NewTodo(*cfg, repos)

	app := &App{
		config:      cfg,
		repos:       repos,
		memoUseCase: memoUseCase,
		todoUseCase: todoUseCase,
		logger:      logger,
	}

	return app, nil
}

// Config returns the application configuration
func (a *App) Config() *config.TomlConfig {
	return a.config
}

// Repositories returns the application repositories
func (a *App) Repositories() repository.Repositories {
	return a.repos
}

// MemoUseCase returns the memo use case
func (a *App) MemoUseCase() memo.Memo {
	return a.memoUseCase
}

// TodoUseCase returns the todo use case
func (a *App) TodoUseCase() todo.Todo {
	return a.todoUseCase
}

// Logger returns the application logger
func (a *App) Logger() *slog.Logger {
	return a.logger
}
