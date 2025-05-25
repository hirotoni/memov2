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
	Run: func(cmd *cobra.Command, args []string) {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			cmd.PrintErrf("Error initializing app: %v\n", err)
			return
		}
		err = ap.Services().Memo().Open(args[0])
		if err != nil {
			cmd.PrintErrf("Error opening memo: %v\n", err)
			return
		}
	},
}
