/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse memos in a terminal UI",
	Long:  `Browse memos in an interactive terminal UI where you can navigate through folders and files.`,
	Run: func(cmd *cobra.Command, args []string) {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			cmd.PrintErrf("Error initializing app: %v\n", err)
			return
		}
		err = ap.Services().Memo().Browse()
		if err != nil {
			cmd.PrintErrf("Error browsing memos: %v\n", err)
			return
		}
	},
}

func init() {}
