package service

import (
	"log/slog"
	"os"
	"testing"

	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServices(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	editor := platform.NewEditor()

	// Execute
	ucs := NewServices(configProvider, editor, logger)

	// Assert
	assert.NotNil(t, ucs)
	assert.NotNil(t, ucs.Memo())
	assert.NotNil(t, ucs.Todo())
	assert.NotNil(t, ucs.Config())
}

func TestServices_Memo(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	editor := platform.NewEditor()
	ucs := NewServices(configProvider, editor, logger)

	// Execute
	memoService := ucs.Memo()

	// Assert
	assert.NotNil(t, memoService)
}

func TestServices_Todo(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	editor := platform.NewEditor()
	ucs := NewServices(configProvider, editor, logger)

	// Execute
	todoService := ucs.Todo()

	// Assert
	assert.NotNil(t, todoService)
}

func TestServices_Config(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	editor := platform.NewEditor()
	ucs := NewServices(configProvider, editor, logger)

	// Execute
	configService := ucs.Config()

	// Assert
	assert.NotNil(t, configService)
}
