package config

import (
	ccc "github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/platform/editor"
	"github.com/hirotoni/memov2/internal/repository"
)

type config struct {
	config ccc.TomlConfig
	repos  repository.Repositories
	editor editor.Editor
}

func NewConfig(c ccc.TomlConfig, r repository.Repositories, e editor.Editor) config {
	return config{
		config: c,
		repos:  r,
		editor: e,
	}
}
