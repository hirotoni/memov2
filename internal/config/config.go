package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/hirotoni/memov2/internal/platform/fs"
)

const (
	DefaultFolderNameConfig = ".config/memov2/"
	DefaultFolderNameBase   = "dailymemo/"
	DefaultFolderNameTodos  = "todos/"
	DefaultFolderNameMemos  = "memos/"
	DefaultTodosDaysToSeek  = 10 // Default number of days to seek back for todos
)

type Config interface {
	BaseDir() string
	TodosDir() string
	MemosDir() string
	TodosDaysToSeek() string
}

type TomlConfig struct {
	baseDir         string
	todosFolderName string
	memosFolderName string
	todosDaysToSeek int
}

type TomlConfigOption struct {
	BaseDir         string
	TodosFolderName string
	MemosFolderName string
	TodosDaysToSeek int
}

// tomlDTO is a decode/encode surrogate with exported fields for the TOML library
type tomlDTO struct {
	BaseDir         string `toml:"base_dir"`
	TodosFolderName string `toml:"todos_foldername"`
	MemosFolderName string `toml:"memos_foldername"`
	TodosDaysToSeek int    `toml:"todos_daystoseek"`
}

func (tc *TomlConfig) toDTO() tomlDTO {
	return tomlDTO{
		BaseDir:         tc.baseDir,
		TodosFolderName: tc.todosFolderName,
		MemosFolderName: tc.memosFolderName,
		TodosDaysToSeek: tc.todosDaysToSeek,
	}
}

func fromDTO(d tomlDTO) *TomlConfig {
	return &TomlConfig{
		baseDir:         d.BaseDir,
		todosFolderName: d.TodosFolderName,
		memosFolderName: d.MemosFolderName,
		todosDaysToSeek: d.TodosDaysToSeek,
	}
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
		c.baseDir = tco.BaseDir
	}
	if tco.TodosFolderName != "" {
		c.todosFolderName = tco.TodosFolderName
	}
	if tco.MemosFolderName != "" {
		c.memosFolderName = tco.MemosFolderName
	}
	if tco.TodosDaysToSeek > 0 {
		c.todosDaysToSeek = tco.TodosDaysToSeek
	}

	return c, nil
}

func NewDefaultConfig() (*TomlConfig, error) {
	dir, _, err := ConfigDirPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory path: %w", err)
	}

	return &TomlConfig{
		baseDir:         filepath.Join(dir, DefaultFolderNameBase),
		todosFolderName: DefaultFolderNameTodos,
		memosFolderName: DefaultFolderNameMemos,
		todosDaysToSeek: DefaultTodosDaysToSeek,
	}, nil
}

// EnsureDirectories creates all necessary directories for the configuration
func (tc *TomlConfig) EnsureDirectories() error {
	dirs := []string{
		tc.BaseDir(),
		tc.TodosDir(),
		tc.MemosDir(),
	}

	for _, dir := range dirs {
		if err := fs.EnsureDir(dir); err != nil {
			return fmt.Errorf("failed to ensure directory %s: %w", dir, err)
		}
	}

	return nil
}

// getters
func (tc *TomlConfig) TodosDir() string     { return filepath.Join(tc.baseDir, tc.todosFolderName) }
func (tc *TomlConfig) MemosDir() string     { return filepath.Join(tc.baseDir, tc.memosFolderName) }
func (tc *TomlConfig) BaseDir() string      { return tc.baseDir }
func (tc *TomlConfig) TodosDaysToSeek() int { return tc.todosDaysToSeek }

func LoadTomlConfig() (*TomlConfig, error) {
	dir, path, err := ConfigDirPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory path: %w", err)
	}

	if err := fs.EnsureDir(dir); err != nil {
		return nil, fmt.Errorf("failed to ensure config directory: %w", err)
	}

	if !fs.Exists(path) {
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

		if err := fs.EnsureDir(tomlConfig.BaseDir()); err != nil {
			return nil, fmt.Errorf("failed to ensure base directory: %w", err)
		}
		if err := fs.EnsureDir(tomlConfig.TodosDir()); err != nil {
			return nil, fmt.Errorf("failed to ensure todos directory: %w", err)
		}
		if err := fs.EnsureDir(tomlConfig.MemosDir()); err != nil {
			return nil, fmt.Errorf("failed to ensure memos directory: %w", err)
		}

		return tomlConfig, nil
	}

	var dto tomlDTO
	if _, err = toml.DecodeFile(path, &dto); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return fromDTO(dto), nil
}
