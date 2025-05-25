package todo

import (
	"github.com/hirotoni/memov2/internal/interfaces"
)

type todo struct {
	c interfaces.ConfigProvider
	r interfaces.Repositories
	e interfaces.Editor
}

func NewTodo(c interfaces.ConfigProvider, r interfaces.Repositories, e interfaces.Editor) interfaces.TodoUsecase {
	return todo{
		c: c,
		r: r,
		e: e,
	}
}
