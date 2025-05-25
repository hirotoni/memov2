package app

import (
	"fmt"
	"log/slog"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/usecase/interfaces"
)

// App represents the main application container
type App struct {
	config   *config.TomlConfig
	usecases interfaces.Usecases
	logger   *slog.Logger
}

// NewApp creates a new application instance with all dependencies
func NewApp(cfg *config.TomlConfig, logger *slog.Logger) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	// Initialize use cases
	ucs := interfaces.NewUsecases(*cfg)

	app := &App{
		config:   cfg,
		usecases: ucs,
		logger:   logger,
	}

	return app, nil
}

// Config returns the application configuration
func (a *App) Config() *config.TomlConfig {
	return a.config
}

// Usecases returns the application Usecases
func (a *App) Usecases() interfaces.Usecases {
	return a.usecases
}

// Logger returns the application logger
func (a *App) Logger() *slog.Logger {
	return a.logger
}
