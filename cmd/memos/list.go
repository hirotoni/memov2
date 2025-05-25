package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

var shortFlag bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all memos",
	Long:  `List all memos in "title\tpath" format. Useful with peco for interactive selection.`,
	Run: func(cmd *cobra.Command, args []string) {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			cmd.PrintErrf("Error initializing app: %v\n", err)
			return
		}
		err = ap.Services().Memo().List(!shortFlag)
		if err != nil {
			cmd.PrintErrf("Error listing memos: %v\n", err)
			return
		}
	},
}

func init() {
	listCmd.Flags().BoolVarP(&shortFlag, "short", "s", false, "show relative paths instead of full paths")
}
