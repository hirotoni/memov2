package repositories

import (
	"log/slog"

	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/repositories/memo"
	"github.com/hirotoni/memov2/internal/repositories/todo"
	"github.com/hirotoni/memov2/internal/repositories/weekly"
)

type repositories struct {
	memo       interfaces.MemoRepo
	memoWeekly interfaces.WeeklyRepo
	todo       interfaces.TodoRepo
	todoWeekly interfaces.WeeklyRepo
}

func NewRepositories(c interfaces.ConfigProvider, logger *slog.Logger) interfaces.Repositories {
	r := repositories{
		memo:       memo.NewMemo(c.MemosDir(), logger),
		memoWeekly: weekly.NewWeekly(c.MemosDir(), logger),
		todo:       todo.NewTodo(c.TodosDir(), logger),
		todoWeekly: weekly.NewWeekly(c.TodosDir(), logger),
	}
	return r
}

func (r repositories) Memo() interfaces.MemoRepo         { return r.memo }
func (r repositories) Todo() interfaces.TodoRepo         { return r.todo }
func (r repositories) MemoWeekly() interfaces.WeeklyRepo { return r.memoWeekly }
func (r repositories) TodoWeekly() interfaces.WeeklyRepo { return r.todoWeekly }
