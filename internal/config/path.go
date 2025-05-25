package config

import (
	"os"
	"path/filepath"
)

// ConfigDir returns the configuration directory path
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(home, DefaultFolderNameConfig)
	return configDir, nil
}
