package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/hirotoni/memov2/utils"
)

const (
	DefaultFolderNameConfig = ".config/memov2/"
	DefaultFolderNameBase   = "dailymemo/"
	DefaultFolderNameTodos  = "todos/"
	DefaultFolderNameMemos  = "memos/"
	DefaultTodosDaysToSeek  = 10 // Default number of days to seek back for todos
)

type Config interface {
	TodosDir() string
	MemosDir() string
}

type TomlConfig struct {
	BaseDir         string `toml:"base_dir"`         // base directory for storing data
	TodosFolderName string `toml:"todos_foldername"` // folder name for todos
	MemosFolderName string `toml:"memos_foldername"` // folder name for memos
	TodosDaysToSeek int    `toml:"todos_daystoseek"` // days to seek back
}

type TomlConfigOption struct {
	BaseDir         string
	TodosFolderName string
	MemosFolderName string
	TodosDaysToSeek int
}

func ConfigDirPath() (string, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}

	configDir := filepath.Join(home, DefaultFolderNameConfig)
	fpath := filepath.Join(configDir, "config.toml")

	return configDir, fpath, nil
}

func NewTomlConfig(tco TomlConfigOption) (*TomlConfig, error) {
	c, err := NewDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create default config: %w", err)
	}

	if tco.BaseDir != "" {
		c.BaseDir = tco.BaseDir
	}
	if tco.TodosFolderName != "" {
		c.TodosFolderName = tco.TodosFolderName
	}
	if tco.MemosFolderName != "" {
		c.MemosFolderName = tco.MemosFolderName
	}
	if tco.TodosDaysToSeek > 0 {
		c.TodosDaysToSeek = tco.TodosDaysToSeek
	}

	return c, nil
}

func NewDefaultConfig() (*TomlConfig, error) {
	dir, _, err := ConfigDirPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory path: %w", err)
	}

	return &TomlConfig{
		BaseDir:         filepath.Join(dir, DefaultFolderNameBase),
		TodosFolderName: DefaultFolderNameTodos,
		MemosFolderName: DefaultFolderNameMemos,
		TodosDaysToSeek: DefaultTodosDaysToSeek,
	}, nil
}

func (tc *TomlConfig) TodosDir() string {
	return filepath.Join(tc.BaseDir, tc.TodosFolderName)
}
func (tc *TomlConfig) MemosDir() string {
	return filepath.Join(tc.BaseDir, tc.MemosFolderName)
}

func LoadTomlConfig() (*TomlConfig, error) {
	dir, path, err := ConfigDirPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory path: %w", err)
	}

	if !utils.Exists(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	if !utils.Exists(path) {
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
		if err := enc.Encode(tomlConfig); err != nil {
			return nil, fmt.Errorf("failed to encode default config: %w", err)
		}

		if err := os.MkdirAll(tomlConfig.BaseDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create base directory: %w", err)
		}
		if err := os.MkdirAll(tomlConfig.TodosDir(), 0755); err != nil {
			return nil, fmt.Errorf("failed to create todos directory: %w", err)
		}
		if err := os.MkdirAll(tomlConfig.MemosDir(), 0755); err != nil {
			return nil, fmt.Errorf("failed to create memos directory: %w", err)
		}

		return tomlConfig, nil
	}

	tomlConfig := &TomlConfig{}
	_, err = toml.DecodeFile(path, tomlConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return tomlConfig, nil
}
