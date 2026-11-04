package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open [path]",
	Short: "open a memo file in the configured editor",
	Long:  `Open a memo file in the configured editor. The path is relative to the base directory.`,
	Args:  cobra.ExactArgs(1),
	// RunE + SilenceUsage so a failure exits non-zero with the error on stderr.
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			return err
		}
		return ap.Services().Memo().Open(args[0])
	},
}
