/*
Copyright © 2025 hirotoni
*/
package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

// newCmd launches an interactive TUI to pick a category and enter a title, then
// creates the memo and opens it in the configured editor.
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "interactively pick a category and create a memo",
	Long:  `Pick a category in an embedded TUI (or "no category"), type a title, and create the memo.`,
	Run: func(cmd *cobra.Command, args []string) {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			cmd.PrintErrf("Error initializing app: %v\n", err)
			return
		}
		if err := ap.Services().Memo().NewInteractive(); err != nil {
			cmd.PrintErrf("Error creating memo: %v\n", err)
			return
		}
	},
}
