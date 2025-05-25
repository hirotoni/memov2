package todo

import (
	"log/slog"

	"github.com/hirotoni/memov2/internal/interfaces"
)

type todo struct {
	c      interfaces.ConfigProvider
	r      interfaces.Repositories
	e      interfaces.Editor
	logger *slog.Logger
}

func NewTodo(c interfaces.ConfigProvider, r interfaces.Repositories, e interfaces.Editor, logger *slog.Logger) interfaces.TodoService {
	return todo{
		c:      c,
		r:      r,
		e:      e,
		logger: logger,
	}
}
