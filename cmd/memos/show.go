package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

var showJSONFlag bool

// showCmd prints a single memo identified by its path (relative to the memos
// directory, e.g. from `memos list --short`, or absolute). By default it emits
// the raw markdown content; with --json it emits a structured JSON object.
// It uses RunE so that a rejected path (e.g. one escaping the memos directory)
// produces a non-zero exit code with nothing written to stdout.
var showCmd = &cobra.Command{
	Use:           "show <path>",
	Short:         "show a single memo",
	Long:          `Show a single memo by path. By default prints the raw markdown content; use --json for a machine-readable object. Paths escaping the memos directory are rejected.`,
	Args:          cobra.ExactArgs(1),
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			return err
		}
		return ap.Services().Memo().Show(args[0], showJSONFlag)
	},
}

func init() {
	showCmd.Flags().BoolVar(&showJSONFlag, "json", false, "output as JSON to stdout")
}
