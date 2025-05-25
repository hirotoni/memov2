package app

import (
	"log/slog"
	"os"
	"testing"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewApp_Success(t *testing.T) {
	// Setup
	cfg := &config.TomlConfig{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Execute
	app, err := NewApp(cfg, logger)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, app)
	assert.NotNil(t, app.Config())
	assert.NotNil(t, app.Usecases())
	assert.NotNil(t, app.Logger())
	assert.Equal(t, logger, app.Logger())
}

func TestNewApp_NilConfig(t *testing.T) {
	// Setup
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Execute
	app, err := NewApp(nil, logger)

	// Assert
	require.Error(t, err)
	assert.Nil(t, app)
	assert.Contains(t, err.Error(), "config cannot be nil")
}

func TestNewApp_NilLogger(t *testing.T) {
	// Setup
	cfg := &config.TomlConfig{}

	// Execute
	app, err := NewApp(cfg, nil)

	// Assert
	require.Error(t, err)
	assert.Nil(t, app)
	assert.Contains(t, err.Error(), "logger cannot be nil")
}

func TestApp_Config(t *testing.T) {
	// Setup
	cfg := &config.TomlConfig{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app, err := NewApp(cfg, logger)
	require.NoError(t, err)

	// Execute
	config := app.Config()

	// Assert
	assert.NotNil(t, config)
}

func TestApp_Usecases(t *testing.T) {
	// Setup
	cfg := &config.TomlConfig{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app, err := NewApp(cfg, logger)
	require.NoError(t, err)

	// Execute
	usecases := app.Usecases()

	// Assert
	assert.NotNil(t, usecases)
	assert.NotNil(t, usecases.Memo())
	assert.NotNil(t, usecases.Todo())
	assert.NotNil(t, usecases.Config())
}

func TestApp_Logger(t *testing.T) {
	// Setup
	cfg := &config.TomlConfig{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app, err := NewApp(cfg, logger)
	require.NoError(t, err)

	// Execute
	retrievedLogger := app.Logger()

	// Assert
	assert.NotNil(t, retrievedLogger)
	assert.Equal(t, logger, retrievedLogger)
}
