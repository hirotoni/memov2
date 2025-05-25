package toml

import (
	"path/filepath"

	"github.com/hirotoni/memov2/internal/config"
)

// Config represents the TOML-based configuration
type Config struct {
	baseDir         string
	todosFolderName string
	memosFolderName string
	todosDaysToSeek int
}

// Option holds configuration options for creating a new Config
type Option struct {
	BaseDir         string
	TodosFolderName string
	MemosFolderName string
	TodosDaysToSeek int
}

// DTO is a decode/encode surrogate with exported fields for the TOML library
type DTO struct {
	BaseDir         string `toml:"base_dir"`
	TodosFolderName string `toml:"todos_foldername"`
	MemosFolderName string `toml:"memos_foldername"`
	TodosDaysToSeek int    `toml:"todos_daystoseek"`
}

// toDTO converts Config to DTO for TOML encoding
func (c *Config) toDTO() DTO {
	return DTO{
		BaseDir:         c.baseDir,
		TodosFolderName: c.todosFolderName,
		MemosFolderName: c.memosFolderName,
		TodosDaysToSeek: c.todosDaysToSeek,
	}
}

// fromDTO creates Config from DTO after TOML decoding
func fromDTO(d DTO) *Config {
	return &Config{
		baseDir:         d.BaseDir,
		todosFolderName: d.TodosFolderName,
		memosFolderName: d.MemosFolderName,
		todosDaysToSeek: d.TodosDaysToSeek,
	}
}

// BaseDir returns the base directory path
func (c *Config) BaseDir() string {
	return c.baseDir
}

// TodosDir returns the todos directory path
func (c *Config) TodosDir() string {
	return filepath.Join(c.baseDir, c.todosFolderName)
}

// MemosDir returns the memos directory path
func (c *Config) MemosDir() string {
	return filepath.Join(c.baseDir, c.memosFolderName)
}

// TodosDaysToSeek returns the number of days to seek for todos
func (c *Config) TodosDaysToSeek() int {
	return c.todosDaysToSeek
}

// ConfigDirPath returns the config directory and file path
func (c *Config) ConfigDirPath() (string, string, error) {
	dir, err := config.ConfigDir()
	if err != nil {
		return "", "", err
	}
	path := filepath.Join(dir, "config.toml")
	return dir, path, nil
}

// EnsureDirectories creates all necessary directories for the configuration
func (c *Config) EnsureDirectories() error {
	dirs := []string{
		c.BaseDir(),
		c.TodosDir(),
		c.MemosDir(),
	}

	for _, dir := range dirs {
		if err := ensureDir(dir); err != nil {
			return err
		}
	}

	return nil
}
