/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package memos

import (
	"github.com/hirotoni/memov2/components/memos"
	"github.com/hirotoni/memov2/config"
	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "generate memo index file",
	Long:  `generate a memo index file that lists all memos.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.LoadTomlConfig()
		if err != nil {
			cmd.PrintErrf("Error loading config: %v\n", err)
			return
		}
		err = memos.GenerateMemoIndex(*c)
		if err != nil {
			cmd.PrintErrf("Error generating memo index: %v\n", err)
			return
		}
	},
}

func init() {}
