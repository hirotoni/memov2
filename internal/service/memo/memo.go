package memo

import (
	"log/slog"

	"github.com/hirotoni/memov2/internal/interfaces"
)

type memo struct {
	config interfaces.ConfigProvider
	repos  interfaces.Repositories
	editor interfaces.Editor
	logger *slog.Logger
}

func NewMemo(c interfaces.ConfigProvider, r interfaces.Repositories, e interfaces.Editor, logger *slog.Logger) interfaces.MemoService {
	return memo{
		config: c,
		repos:  r,
		editor: e,
		logger: logger,
	}
}
