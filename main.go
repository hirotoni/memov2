/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log/slog"
	"os"

	"github.com/hirotoni/memov2/cmd"
	"github.com/hirotoni/memov2/internal/app"
	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/logger"
)

func main() {
	// Initialize logger
	log := logger.DefaultLogger()
	slog.SetDefault(log)

	// Load configuration
	cfg, err := config.LoadTomlConfig()
	if err != nil {
		log.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Create application
	_, err = app.NewApp(cfg, log)
	if err != nil {
		log.Error("Failed to create application", "error", err)
		os.Exit(1)
	}

	// Execute the application
	cmd.Execute()
}
