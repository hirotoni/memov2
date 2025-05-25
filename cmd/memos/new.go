/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

// newCmd represents the memo command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "generate a new memo file",
	Long:  `generate a new memo file",`,
	Run: func(cmd *cobra.Command, args []string) {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			cmd.PrintErrf("Error initializing app: %v\n", err)
			return
		}

		var title string
		if len(args) > 0 {
			title = args[0]
		}

		err = ap.Services().Memo().GenerateMemoFile(title)
		if err != nil {
			cmd.PrintErrf("Error generating memo file: %v\n", err)
			return
		}
	},
}

func init() {}
