package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

var (
	shortFlag    bool
	listJSONFlag bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all memos",
	Long:  `List all memos in "title\tpath" format. Useful with peco for interactive selection.`,
	// RunE + SilenceUsage so a failure exits non-zero with nothing on stdout,
	// upholding the machine-output contract for --json consumers.
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			return err
		}
		return ap.Services().Memo().List(!shortFlag, listJSONFlag)
	},
}

func init() {
	listCmd.Flags().BoolVarP(&shortFlag, "short", "s", false, "show relative paths instead of full paths")
	listCmd.Flags().BoolVar(&listJSONFlag, "json", false, "output as JSON to stdout")
}
