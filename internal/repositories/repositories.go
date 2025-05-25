package repositories

import (
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

func NewRepositories(c interfaces.ConfigProvider) interfaces.Repositories {
	r := repositories{
		memo:       memo.NewMemo(c.MemosDir()),
		memoWeekly: weekly.NewWeekly(c.MemosDir()),
		todo:       todo.NewTodo(c.TodosDir()),
		todoWeekly: weekly.NewWeekly(c.TodosDir()),
	}
	return r
}

func (r repositories) Memo() interfaces.MemoRepo         { return r.memo }
func (r repositories) Todo() interfaces.TodoRepo         { return r.todo }
func (r repositories) MemoWeekly() interfaces.WeeklyRepo { return r.memoWeekly }
func (r repositories) TodoWeekly() interfaces.WeeklyRepo { return r.todoWeekly }
