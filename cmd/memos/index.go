/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "generate memo index file",
	Long:  `generate a memo index file that lists all memos.`,
	Run: func(cmd *cobra.Command, args []string) {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			cmd.PrintErrf("Error initializing app: %v\n", err)
			return
		}

		err = ap.Services().Memo().GenerateMemoIndex()
		if err != nil {
			cmd.PrintErrf("Error generating memo file: %v\n", err)
			return
		}
	},
}

func init() {}
