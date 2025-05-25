/*
Copyright © 2025 hirotoni
*/
package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

// renameCmd launches an interactive TUI to pick a memo and enter a new title,
// then renames it (updates both title and filename).
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "interactively pick a memo and rename it",
	Long:  `Pick a memo in an embedded TUI, type a new title, and rename it (updates both title and filename).`,
	Run: func(cmd *cobra.Command, args []string) {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			cmd.PrintErrf("Error initializing app: %v\n", err)
			return
		}
		if err := ap.Services().Memo().RenameInteractive(); err != nil {
			cmd.PrintErrf("Error renaming memo: %v\n", err)
			return
		}
	},
}
