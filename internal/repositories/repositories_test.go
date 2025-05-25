package repositories

import (
	"testing"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepositories(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := config.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := config.NewTomlConfig(opts)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)

	// Execute
	repos := NewRepositories(configProvider)

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
	opts := config.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := config.NewTomlConfig(opts)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	repos := NewRepositories(configProvider)

	// Execute
	memoRepo := repos.Memo()

	// Assert
	assert.NotNil(t, memoRepo)
}

func TestRepositories_Todo(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := config.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := config.NewTomlConfig(opts)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	repos := NewRepositories(configProvider)

	// Execute
	todoRepo := repos.Todo()

	// Assert
	assert.NotNil(t, todoRepo)
}

func TestRepositories_MemoWeekly(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := config.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := config.NewTomlConfig(opts)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	repos := NewRepositories(configProvider)

	// Execute
	weeklyRepo := repos.MemoWeekly()

	// Assert
	assert.NotNil(t, weeklyRepo)
}

func TestRepositories_TodoWeekly(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := config.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := config.NewTomlConfig(opts)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	repos := NewRepositories(configProvider)

	// Execute
	weeklyRepo := repos.TodoWeekly()

	// Assert
	assert.NotNil(t, weeklyRepo)
}
