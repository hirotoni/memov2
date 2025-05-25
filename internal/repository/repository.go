package repository

import "github.com/hirotoni/memov2/internal/config"

type repositories struct {
	memo       Memo
	memoWeekly Weekly
	todo       Todo
	todoWeekly Weekly
}

func NewRepositories(c config.TomlConfig) Repositories {
	r := repositories{
		memo:       NewMemo(c.MemosDir()),
		memoWeekly: NewWeekly(c.MemosDir()),
		todo:       NewTodo(c.TodosDir()),
		todoWeekly: NewWeekly(c.TodosDir()),
	}
	return r
}

func (r repositories) Memo() Memo         { return r.memo }
func (r repositories) Todo() Todo         { return r.todo }
func (r repositories) MemoWeekly() Weekly { return r.memoWeekly }
func (r repositories) TodoWeekly() Weekly { return r.todoWeekly }
