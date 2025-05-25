package app

import (
	"fmt"
	"log/slog"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/usecases"
)

// App represents the main application container
type App struct {
	config   interfaces.ConfigProvider
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
	configProvider := interfaces.NewConfigProvider(cfg)
	ucs := usecases.NewUsecases(configProvider)

	app := &App{
		config:   configProvider,
		usecases: ucs,
		logger:   logger,
	}

	return app, nil
}

// Config returns the application configuration
func (a *App) Config() interfaces.ConfigProvider {
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
