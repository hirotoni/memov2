package interfaces

import (
	ccc "github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/platform/editor"
	"github.com/hirotoni/memov2/internal/repository"

	"github.com/hirotoni/memov2/internal/usecase/config"
	"github.com/hirotoni/memov2/internal/usecase/memo"
	"github.com/hirotoni/memov2/internal/usecase/todo"
)

type Usecases interface {
	Memo() Memo
	Todo() Todo
	Config() Config
}

type Memo interface {
	BuildWeeklyReportMemos() error
	GenerateMemoFile(title string) error
	GenerateMemoIndex() error
	Browse() error
}

type Todo interface {
	GenerateTodoFile(truncate bool) error
	BuildWeeklyReportTodos() error
}

type Config interface {
	Show()
	Edit() error
}

type usecases struct {
	memo   Memo
	todo   Todo
	config Config
}

func NewUsecases(c ccc.TomlConfig) Usecases {
	// Initialize editor
	e := editor.New(&c)

	// Initialize repositories
	r := repository.NewRepositories(c)

	uc := usecases{
		memo:   memo.NewMemo(c, r),
		todo:   todo.NewTodo(c, r, e),
		config: config.NewConfig(c, r, e),
	}
	return uc
}

func (r usecases) Memo() Memo     { return r.memo }
func (r usecases) Todo() Todo     { return r.todo }
func (r usecases) Config() Config { return r.config }
