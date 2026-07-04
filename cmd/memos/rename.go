/*
Copyright © 2025 hirotoni
*/
package memos

import (
	"fmt"

	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

var renameJSONFlag bool

// renameCmd renames a memo. With no positional arguments it launches an
// interactive TUI to pick a memo and enter a new title. With two arguments
// (`memos rename <path> <new-title>`) it renames non-interactively (for
// AI/automation): by default it prints the new relative path, or a JSON object
// with --json.
var renameCmd = &cobra.Command{
	Use:           "rename [<path> <new-title>]",
	Short:         "rename a memo (interactively, or non-interactively with <path> <new-title>)",
	Long:          `With no arguments, pick a memo in an embedded TUI and type a new title. With two arguments (path and new title), rename non-interactively (no TUI): prints the new relative path, or a JSON object with --json. Paths escaping the memos directory are rejected.`,
	SilenceUsage:  true,
	SilenceErrors: false,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 && len(args) != 2 {
			return fmt.Errorf("accepts 0 args (interactive) or 2 args (<path> <new-title>), received %d", len(args))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			return err
		}
		if len(args) == 2 {
			return ap.Services().Memo().Rename(args[0], args[1], renameJSONFlag)
		}
		return ap.Services().Memo().RenameInteractive()
	},
}

func init() {
	renameCmd.Flags().BoolVar(&renameJSONFlag, "json", false, "output as JSON to stdout (non-interactive mode)")
}
