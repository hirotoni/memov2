package usecases

import (
	"testing"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUsecases(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := config.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := config.NewTomlConfig(opts)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)

	// Execute
	ucs := NewUsecases(configProvider)

	// Assert
	assert.NotNil(t, ucs)
	assert.NotNil(t, ucs.Memo())
	assert.NotNil(t, ucs.Todo())
	assert.NotNil(t, ucs.Config())
}

func TestUsecases_Memo(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := config.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := config.NewTomlConfig(opts)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	ucs := NewUsecases(configProvider)

	// Execute
	memoUsecase := ucs.Memo()

	// Assert
	assert.NotNil(t, memoUsecase)
}

func TestUsecases_Todo(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := config.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := config.NewTomlConfig(opts)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	ucs := NewUsecases(configProvider)

	// Execute
	todoUsecase := ucs.Todo()

	// Assert
	assert.NotNil(t, todoUsecase)
}

func TestUsecases_Config(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := config.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := config.NewTomlConfig(opts)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	ucs := NewUsecases(configProvider)

	// Execute
	configUsecase := ucs.Config()

	// Assert
	assert.NotNil(t, configUsecase)
}
