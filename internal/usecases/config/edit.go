package config

import (
	"fmt"
)

func (uc config) Edit() error {
	dir, path, err := uc.config.ConfigDirPath()
	if err != nil {
		return fmt.Errorf("error getting config path: %v", err)
	}

	return uc.editor.Open(dir, path)
}
