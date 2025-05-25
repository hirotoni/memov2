package service

import (
	"log/slog"

	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/repositories"
	"github.com/hirotoni/memov2/internal/service/config"
	"github.com/hirotoni/memov2/internal/service/memo"
	"github.com/hirotoni/memov2/internal/service/todo"
)

// services implements the Services interface
type services struct {
	memo   interfaces.MemoService
	todo   interfaces.TodoService
	config interfaces.ConfigService
}

// NewServices creates a new Services instance with all dependencies
func NewServices(c interfaces.ConfigProvider, e interfaces.Editor, logger *slog.Logger) interfaces.Services {
	// Create repositories
	r := repositories.NewRepositories(c, logger)

	// Create and return services
	return services{
		memo:   memo.NewMemo(c, r, e, logger),
		todo:   todo.NewTodo(c, r, e, logger),
		config: config.NewConfig(c, r, e, logger),
	}
}

func (r services) Memo() interfaces.MemoService     { return r.memo }
func (r services) Todo() interfaces.TodoService     { return r.todo }
func (r services) Config() interfaces.ConfigService { return r.config }
