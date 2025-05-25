package app

import (
	"log/slog"

	"github.com/hirotoni/memov2/internal/app"
	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/logger"
	"github.com/spf13/cobra"
)

func InitializeApp(cmd *cobra.Command) (*app.App, error) {
	// Initialize logger
	log := logger.DefaultLogger()
	slog.SetDefault(log)

	// Load configuration
	cfg, err := config.LoadTomlConfig()
	if err != nil {
		cmd.PrintErrf("Failed to load configuration: %v\n", err)
		return nil, err
	}

	// Create application
	ap, err := app.NewApp(cfg, log)
	if err != nil {
		cmd.PrintErrf("Error creating app: %v\n", err)
		return nil, err
	}

	return ap, nil
}
