package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

var categoriesJSONFlag bool

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "list all categories",
	Long:  `List all categories, one per line. Useful with peco/fzf for interactive selection.`,
	// RunE + SilenceUsage so a failure exits non-zero with the error on stderr.
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			return err
		}
		return ap.Services().Memo().ListCategories(categoriesJSONFlag)
	},
}

func init() {
	categoriesCmd.Flags().BoolVar(&categoriesJSONFlag, "json", false, "output as JSON to stdout")
}
