package config

import (
	"log/slog"

	"github.com/hirotoni/memov2/internal/interfaces"
)

type config struct {
	config interfaces.ConfigProvider
	repos  interfaces.Repositories
	editor interfaces.Editor
	logger *slog.Logger
}

func NewConfig(c interfaces.ConfigProvider, r interfaces.Repositories, e interfaces.Editor, logger *slog.Logger) interfaces.ConfigService {
	return config{
		config: c,
		repos:  r,
		editor: e,
		logger: logger,
	}
}
