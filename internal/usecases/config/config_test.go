package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	tomlconfig "github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/repositories"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := tomlconfig.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := tomlconfig.NewTomlConfig(opts)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	repos := repositories.NewRepositories(configProvider)
	mockEditor := mock.NewMockEditor()

	// Execute
	uc := NewConfig(configProvider, repos, mockEditor)

	// Assert
	assert.NotNil(t, uc)
}

func TestConfig_Show(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := tomlconfig.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := tomlconfig.NewTomlConfig(opts)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	repos := repositories.NewRepositories(configProvider)
	mockEditor := mock.NewMockEditor()
	uc := NewConfig(configProvider, repos, mockEditor)

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute
	uc.Show()

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Assert
	assert.Contains(t, output, "base_dir")
	assert.Contains(t, output, "todos_dir")
	assert.Contains(t, output, "memos_dir")
	assert.Contains(t, output, "todos_daystoseek")
	assert.Contains(t, output, tmpDir)
}

func TestConfig_Edit_Success(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := tomlconfig.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := tomlconfig.NewTomlConfig(opts)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	repos := repositories.NewRepositories(configProvider)
	mockEditor := mock.NewMockEditor()
	uc := NewConfig(configProvider, repos, mockEditor)

	// Execute
	err = uc.Edit()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 1, len(mockEditor.Calls), "Editor should be called once")
}

func TestConfig_Edit_EditorError(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := tomlconfig.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := tomlconfig.NewTomlConfig(opts)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	repos := repositories.NewRepositories(configProvider)
	mockEditor := mock.NewMockEditor()

	// Set custom editor function that returns an error
	mockEditor.OpenFunc = func(basedir, path string) error {
		return fmt.Errorf("editor failed to open")
	}

	uc := NewConfig(configProvider, repos, mockEditor)

	// Execute
	err = uc.Edit()

	// Assert
	require.Error(t, err)
	assert.Equal(t, 1, len(mockEditor.Calls), "Editor should be called once")
}
