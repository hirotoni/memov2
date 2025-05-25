package config

import (
	"log"
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

func NewTomlConfig(tco TomlConfigOption) *TomlConfig {
	c := NewDefaultConfig()

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

	return c
}

func NewDefaultConfig() *TomlConfig {
	dir, _, err := ConfigDirPath()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &TomlConfig{
		BaseDir:         filepath.Join(dir, DefaultFolderNameBase),
		TodosFolderName: DefaultFolderNameTodos,
		MemosFolderName: DefaultFolderNameMemos,
		TodosDaysToSeek: DefaultTodosDaysToSeek,
	}
}

func (tc *TomlConfig) TodosDir() string {
	return filepath.Join(tc.BaseDir, tc.TodosFolderName)
}
func (tc *TomlConfig) MemosDir() string {
	return filepath.Join(tc.BaseDir, tc.MemosFolderName)
}

func LoadTomlConfig() *TomlConfig {
	dir, path, err := ConfigDirPath()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	if !utils.Exists(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal(err)
			return nil
		}
	}

	if !utils.Exists(path) {
		f, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		defer f.Close()

		tomlConfig := NewDefaultConfig()

		enc := toml.NewEncoder(f)
		if err := enc.Encode(tomlConfig); err != nil {
			log.Fatal(err)
			return nil
		}

		os.MkdirAll(tomlConfig.BaseDir, 0755)
		os.MkdirAll(tomlConfig.TodosDir(), 0755)
		os.MkdirAll(tomlConfig.MemosDir(), 0755)

		return tomlConfig
	}

	tomlConfig := &TomlConfig{}
	_, err = toml.DecodeFile(path, tomlConfig)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return tomlConfig
}
