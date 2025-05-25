package todo

import (
	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/repository"
)

type Todo interface {
	GenerateTodoFile(truncate bool) error
	BuildWeeklyReportTodos() error
}

type todo struct {
	c config.TomlConfig
	r repository.Repostiories
}

func NewTodo(c config.TomlConfig, r repository.Repostiories) todo {
	return todo{
		c: c,
		r: r,
	}
}
