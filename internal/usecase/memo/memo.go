package memo

import (
	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/repository"
)

type memo struct {
	config config.TomlConfig
	repos  repository.Repositories
}

func NewMemo(c config.TomlConfig, r repository.Repositories) memo {
	return memo{
		config: c,
		repos:  r,
	}
}
