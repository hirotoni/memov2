package usecases

import (
	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/platform/editor"
	"github.com/hirotoni/memov2/internal/repositories"
	usecaseConfig "github.com/hirotoni/memov2/internal/usecases/config"
	"github.com/hirotoni/memov2/internal/usecases/memo"
	"github.com/hirotoni/memov2/internal/usecases/todo"
)

// usecases implements the Usecases interface
type usecases struct {
	memo   interfaces.MemoUsecase
	todo   interfaces.TodoUsecase
	config interfaces.ConfigUsecase
}

// NewUsecases creates a new Usecases instance with all dependencies
func NewUsecases(c interfaces.ConfigProvider) interfaces.Usecases {
	// Create editor - it only needs the config for now
	e := editor.New(c.GetTomlConfig().(*config.TomlConfig))

	// Create repositories
	r := repositories.NewRepositories(c)

	// Create and return usecases
	return usecases{
		memo:   memo.NewMemo(c, r, e),
		todo:   todo.NewTodo(c, r, e),
		config: usecaseConfig.NewConfig(c, r, e),
	}
}

func (r usecases) Memo() interfaces.MemoUsecase     { return r.memo }
func (r usecases) Todo() interfaces.TodoUsecase     { return r.todo }
func (r usecases) Config() interfaces.ConfigUsecase { return r.config }
