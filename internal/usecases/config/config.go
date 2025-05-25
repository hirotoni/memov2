package config

import (
	"github.com/hirotoni/memov2/internal/interfaces"
)

type config struct {
	config interfaces.ConfigProvider
	repos  interfaces.Repositories
	editor interfaces.Editor
}

func NewConfig(c interfaces.ConfigProvider, r interfaces.Repositories, e interfaces.Editor) interfaces.ConfigUsecase {
	return config{
		config: c,
		repos:  r,
		editor: e,
	}
}
