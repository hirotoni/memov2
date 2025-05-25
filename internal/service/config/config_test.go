package config

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/repositories"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repos := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()

	// Execute
	uc := NewConfig(configProvider, repos, mockEditor, logger)

	// Assert
	assert.NotNil(t, uc)
}

func TestConfig_Show(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)

	// Capture stdout for logger
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	repos := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()
	uc := NewConfig(configProvider, repos, mockEditor, logger)

	// Execute
	uc.Show()

	// Read captured output
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
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repos := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()
	uc := NewConfig(configProvider, repos, mockEditor, logger)

	// Execute
	err = uc.Edit()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 1, len(mockEditor.Calls), "Editor should be called once")
}

func TestConfig_Edit_EditorError(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repos := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()

	// Set custom editor function that returns an error
	mockEditor.OpenFunc = func(basedir, path string) error {
		return fmt.Errorf("editor failed to open")
	}

	uc := NewConfig(configProvider, repos, mockEditor, logger)

	// Execute
	err = uc.Edit()

	// Assert
	require.Error(t, err)
	assert.Equal(t, 1, len(mockEditor.Calls), "Editor should be called once")
}
