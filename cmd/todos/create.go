/*
Copyright © 2025 hirotoni
*/
package todos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

// createCmd creates today's todo file. By default (interactive) it opens the file
// in the editor; the former `todos new` is retained as an alias. With --no-open (or
// --json, which implies it) it runs non-interactively: no editor is launched and the
// file's relative path (or a JSON object with --json) is printed to stdout.
var createCmd = &cobra.Command{
	Use:           "create",
	Aliases:       []string{"new"},
	Short:         "create today's todo file",
	Long:          `Create today's todo file. By default it opens the file in the editor (the same flow as the "new" alias); with --no-open (or --json) it runs non-interactively and prints the file's relative path (or a JSON object with --json).`,
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			return err
		}

		// --json implies --no-open: machine-readable output must never open the editor.
		if newNoOpenFlag || newJSONFlag {
			return ap.Services().Todo().GenerateTodoFileNonInteractive(truncateFlag, newJSONFlag)
		}

		return ap.Services().Todo().GenerateTodoFile(truncateFlag)
	},
}

var (
	truncateFlag  bool
	newNoOpenFlag bool
	newJSONFlag   bool
)

func init() {
	createCmd.Flags().BoolVarP(&truncateFlag, "truncate", "t", false, "Truncate the file if it exists")
	createCmd.Flags().BoolVar(&newNoOpenFlag, "no-open", false, "do not open the editor; print the file path to stdout")
	createCmd.Flags().BoolVar(&newJSONFlag, "json", false, "output as JSON to stdout (implies --no-open)")
}
