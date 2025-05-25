package memo

import (
	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/repository"
)

type Memo interface {
	BuildWeeklyReportMemos(c config.TomlConfig) error
	GenerateMemoFile(c config.TomlConfig, title string) error
	GenerateMemoIndex(c config.TomlConfig) error
}

type memo struct {
	c config.TomlConfig
	r repository.Repostiories
}

func NewMemo(c config.TomlConfig, r repository.Repostiories) memo {
	return memo{
		c: c,
		r: r,
	}
}
