package repository

import "github.com/hirotoni/memov2/internal/config"

type Repositories struct {
	Memo       Memo
	MemoWeekly Weekly
	Todo       Todo
	TodoWeekly Weekly
}

func NewRepositories(c config.TomlConfig) Repositories {
	r := Repositories{
		Memo:       NewMemo(c.MemosDir()),
		MemoWeekly: NewWeekly(c.MemosDir()),
		Todo:       NewTodo(c.TodosDir()),
		TodoWeekly: NewWeekly(c.TodosDir()),
	}
	return r
}
