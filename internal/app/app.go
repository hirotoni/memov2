package app

import (
	"log/slog"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/platform"
	"github.com/hirotoni/memov2/internal/service"
)

// App represents the main application container
type App struct {
	config   interfaces.ConfigProvider
	services interfaces.Services
	logger   *slog.Logger
}

// NewApp creates a new application instance with all dependencies
func NewApp(cfg *toml.Config, logger *slog.Logger) (*App, error) {
	if cfg == nil {
		return nil, common.New(common.ErrorTypeConfig, "config cannot be nil")
	}
	if logger == nil {
		return nil, common.New(common.ErrorTypeConfig, "logger cannot be nil")
	}

	// Initialize services
	configProvider := toml.NewProvider(cfg)
	editor := platform.NewEditor()
	ucs := service.NewServices(configProvider, editor, logger)

	app := &App{
		config:   configProvider,
		services: ucs,
		logger:   logger,
	}

	return app, nil
}

// Config returns the application configuration
func (a *App) Config() interfaces.ConfigProvider {
	return a.config
}

// Services returns the application Services
func (a *App) Services() interfaces.Services {
	return a.services
}

// Logger returns the application logger
func (a *App) Logger() *slog.Logger {
	return a.logger
}
