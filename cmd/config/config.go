/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package config

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "config",
	Long:  `config`,
}

func init() {
	ConfigCmd.AddCommand(showCmd)
	ConfigCmd.AddCommand(editCmd)
}
