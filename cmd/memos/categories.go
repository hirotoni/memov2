package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "list all categories",
	Long:  `List all categories, one per line. Useful with peco/fzf for interactive selection.`,
	Run: func(cmd *cobra.Command, args []string) {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			cmd.PrintErrf("Error initializing app: %v\n", err)
			return
		}
		err = ap.Services().Memo().ListCategories()
		if err != nil {
			cmd.PrintErrf("Error listing categories: %v\n", err)
			return
		}
	},
}
