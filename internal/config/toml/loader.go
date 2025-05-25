package toml

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/hirotoni/memov2/internal/config"
)

// NewConfig creates a new Config from Option
func NewConfig(opt Option) (*Config, error) {
	c, err := NewDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create default config: %w", err)
	}

	if opt.BaseDir != "" {
		c.baseDir = opt.BaseDir
	}
	if opt.TodosFolderName != "" {
		c.todosFolderName = opt.TodosFolderName
	}
	if opt.MemosFolderName != "" {
		c.memosFolderName = opt.MemosFolderName
	}
	if opt.TodosDaysToSeek > 0 {
		c.todosDaysToSeek = opt.TodosDaysToSeek
	}

	return c, nil
}

// NewDefaultConfig creates a Config with default values
func NewDefaultConfig() (*Config, error) {
	dir, err := config.ConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory path: %w", err)
	}

	return &Config{
		baseDir:         filepath.Join(dir, config.DefaultFolderNameBase),
		todosFolderName: config.DefaultFolderNameTodos,
		memosFolderName: config.DefaultFolderNameMemos,
		todosDaysToSeek: config.DefaultTodosDaysToSeek,
	}, nil
}

// ConfigFilePath returns the full path to the TOML configuration file
func ConfigFilePath() (string, error) {
	dir, err := config.ConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}
	return filepath.Join(dir, "config.toml"), nil
}

// LoadConfig loads configuration from TOML file
func LoadConfig() (*Config, error) {
	dir, err := config.ConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}
	path, err := ConfigFilePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config file path: %w", err)
	}

	if err := ensureDir(dir); err != nil {
		return nil, fmt.Errorf("failed to ensure config directory: %w", err)
	}

	if !fileExists(path) {
		f, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("failed to create config file: %w", err)
		}
		defer f.Close()

		tomlConfig, err := NewDefaultConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}

		enc := toml.NewEncoder(f)
		if err := enc.Encode(tomlConfig.toDTO()); err != nil {
			return nil, fmt.Errorf("failed to encode default config: %w", err)
		}

		if err := ensureDir(tomlConfig.BaseDir()); err != nil {
			return nil, fmt.Errorf("failed to ensure base directory: %w", err)
		}
		if err := ensureDir(tomlConfig.TodosDir()); err != nil {
			return nil, fmt.Errorf("failed to ensure todos directory: %w", err)
		}
		if err := ensureDir(tomlConfig.MemosDir()); err != nil {
			return nil, fmt.Errorf("failed to ensure memos directory: %w", err)
		}

		return tomlConfig, nil
	}

	var dto DTO
	if _, err = toml.DecodeFile(path, &dto); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return fromDTO(dto), nil
}

// ensureDir ensures the directory exists, creating it (and parents) if necessary.
func ensureDir(dir string) error {
	if dir == "" {
		return nil
	}
	if fileExists(dir) {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	return nil
}

// fileExists returns true if the file or directory at path exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
