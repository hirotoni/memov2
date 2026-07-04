/*
Copyright © 2025 hirotoni
*/
package todos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

var showJSONFlag bool

// showCmd prints a single todo file identified by its path (relative to the todos
// directory, e.g. from `todos list --short`, or absolute). By default it emits the
// raw markdown content; with --json it emits a structured JSON object (daily todos
// only). It uses RunE so that a rejected path (e.g. one escaping the todos
// directory) produces a non-zero exit code with nothing written to stdout.
var showCmd = &cobra.Command{
	Use:           "show <path>",
	Short:         "show a single todo",
	Long:          `Show a single todo by path. By default prints the raw markdown content; use --json for a machine-readable object (daily todos only). Paths escaping the todos directory are rejected.`,
	Args:          cobra.ExactArgs(1),
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			return err
		}
		return ap.Services().Todo().Show(args[0], showJSONFlag)
	},
}

func init() {
	showCmd.Flags().BoolVar(&showJSONFlag, "json", false, "output as JSON to stdout")
}
