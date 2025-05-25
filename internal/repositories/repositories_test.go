package repositories

import (
	"log/slog"
	"os"
	"testing"

	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepositories(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Execute
	repos := NewRepositories(configProvider, logger)

	// Assert
	assert.NotNil(t, repos)
	assert.NotNil(t, repos.Memo())
	assert.NotNil(t, repos.Todo())
	assert.NotNil(t, repos.MemoWeekly())
	assert.NotNil(t, repos.TodoWeekly())
}

func TestRepositories_Memo(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repos := NewRepositories(configProvider, logger)

	// Execute
	memoRepo := repos.Memo()

	// Assert
	assert.NotNil(t, memoRepo)
}

func TestRepositories_Todo(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repos := NewRepositories(configProvider, logger)

	// Execute
	todoRepo := repos.Todo()

	// Assert
	assert.NotNil(t, todoRepo)
}

func TestRepositories_MemoWeekly(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repos := NewRepositories(configProvider, logger)

	// Execute
	weeklyRepo := repos.MemoWeekly()

	// Assert
	assert.NotNil(t, weeklyRepo)
}

func TestRepositories_TodoWeekly(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repos := NewRepositories(configProvider, logger)

	// Execute
	weeklyRepo := repos.TodoWeekly()

	// Assert
	assert.NotNil(t, weeklyRepo)
}
