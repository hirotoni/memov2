package app

import (
	"log/slog"

	"github.com/hirotoni/memov2/internal/app"
	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/spf13/cobra"
)

func InitializeApp(cmd *cobra.Command) (*app.App, error) {
	// Initialize logger
	log := common.DefaultLogger()
	slog.SetDefault(log)

	// Load configuration
	cfg, err := toml.LoadConfig()
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
