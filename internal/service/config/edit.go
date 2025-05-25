package config

import (
	"github.com/hirotoni/memov2/internal/common"
)

func (uc config) Edit() error {
	dir, path, err := uc.config.ConfigDirPath()
	if err != nil {
		return common.Wrap(err, common.ErrorTypeConfig, "error getting config path")
	}

	return uc.editor.Open(dir, path)
}
