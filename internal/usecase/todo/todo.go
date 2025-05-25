package todo

import (
	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/platform/editor"
	"github.com/hirotoni/memov2/internal/repository"
)

type todo struct {
	c config.TomlConfig
	r repository.Repositories
	e editor.Editor
}

func NewTodo(c config.TomlConfig, r repository.Repositories, e editor.Editor) todo {
	return todo{
		c: c,
		r: r,
		e: e,
	}
}
