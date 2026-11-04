package config

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

var showJSONFlag bool

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show config",
	Long:  `show config`,
	// RunE + SilenceUsage so a failure exits non-zero with the error on stderr.
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			return err
		}
		return ap.Services().Config().Show(showJSONFlag)
	},
}

func init() {
	showCmd.Flags().BoolVar(&showJSONFlag, "json", false, "output as JSON to stdout")
}
