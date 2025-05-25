/*
Copyright © 2025 hirotoni
*/
package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

// searchCmd launches an interactive, romaji-aware search TUI. Selecting a result
// opens it in the configured editor.
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "interactively search memos and open the selection",
	Long:  `Incrementally search memos in an embedded TUI (romaji-aware): type to filter, ctrl+n/ctrl+p to move, enter to open the highlighted memo in the configured editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			cmd.PrintErrf("Error initializing app: %v\n", err)
			return
		}
		if err := ap.Services().Memo().SearchInteractive(); err != nil {
			cmd.PrintErrf("Error searching memos: %v\n", err)
			return
		}
	},
}
