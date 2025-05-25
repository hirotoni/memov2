package repository

import "github.com/hirotoni/memov2/internal/config"

type Repostiories struct {
	Memo       Memo
	MemoWeekly Weekly
	Todo       Todo
	TodoWeekly Weekly
}

func NewRepositories(c config.TomlConfig) Repostiories {
	r := Repostiories{
		Memo:       NewMemo(c.MemosDir()),
		MemoWeekly: NewWeekly(c.MemosDir()),
		Todo:       NewTodo(c.TodosDir()),
		TodoWeekly: NewWeekly(c.TodosDir()),
	}
	return r
}
