/*
Copyright © 2025 hirotoni
*/
package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

var searchJSONFlag bool

// searchCmd searches memos. With no argument it launches an interactive,
// romaji-aware search TUI (selecting a result opens it in the configured
// editor). With a single query argument it runs a non-interactive search and
// prints the matches to stdout ("title\tpath" per memo, or JSON with --json),
// which is suitable for scripting and AI/automation.
var searchCmd = &cobra.Command{
	Use:           "search [query]",
	Short:         "search memos (interactive TUI, or non-interactive with a query)",
	Long:          `Search memos. With no argument, incrementally search in an embedded TUI (romaji-aware): type to filter, ctrl+n/ctrl+p to move, enter to open the highlighted memo. With a query argument, run a non-interactive search and print matches ("title\tpath", or JSON with --json).`,
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			// Preserve the interactive path's legacy behavior (print, exit 0);
			// the non-interactive path surfaces a non-zero exit for scripting/AI.
			if len(args) == 0 {
				cmd.PrintErrf("Error searching memos: %v\n", err)
				return nil
			}
			return err
		}
		if len(args) == 0 {
			if err := ap.Services().Memo().SearchInteractive(); err != nil {
				cmd.PrintErrf("Error searching memos: %v\n", err)
				return nil
			}
			return nil
		}
		return ap.Services().Memo().Search(args[0], searchJSONFlag)
	},
}

func init() {
	searchCmd.Flags().BoolVar(&searchJSONFlag, "json", false, "output as JSON to stdout (only with a query argument)")
}
