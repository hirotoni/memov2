/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package memos

import (
	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/ui/tui/memos"
	"github.com/spf13/cobra"
)

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse memos in a terminal UI",
	Long:  `Browse memos in an interactive terminal UI where you can navigate through folders and files.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.LoadTomlConfig()
		if err != nil {
			cmd.PrintErrf("Error loading config: %v\n", err)
			return
		}
		err = memos.IntegratedMemos(c)
		if err != nil {
			cmd.PrintErrf("Error browsing memos: %v\n", err)
			return
		}
	},
}

func init() {}
