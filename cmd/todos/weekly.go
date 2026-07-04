/*
Copyright © 2025 hirotoni
*/
package todos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

// weeklyCmd generates the weekly todo report. By default it opens the report in
// the editor. With --no-open (or --json, which implies it) it runs
// non-interactively: no editor is launched and the report's relative path (or a
// JSON object with --json) is printed to stdout.
var weeklyCmd = &cobra.Command{
	Use:           "weekly",
	Short:         "generate weekly report for todos",
	Long:          `Generate the weekly todo report. By default it opens the report in the editor; with --no-open (or --json) it runs non-interactively and prints the report's relative path (or a JSON object with --json).`,
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			return err
		}

		// --json implies --no-open: machine-readable output must never open the editor.
		if weeklyNoOpenFlag || weeklyJSONFlag {
			return ap.Services().Todo().BuildWeeklyReportTodosNonInteractive(weeklyJSONFlag)
		}

		return ap.Services().Todo().BuildWeeklyReportTodos()
	},
}

var (
	weeklyNoOpenFlag bool
	weeklyJSONFlag   bool
)

func init() {
	weeklyCmd.Flags().BoolVar(&weeklyNoOpenFlag, "no-open", false, "do not open the editor; print the report path to stdout")
	weeklyCmd.Flags().BoolVar(&weeklyJSONFlag, "json", false, "output as JSON to stdout (implies --no-open)")
}
