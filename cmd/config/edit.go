package config

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit config",
	Long:  `edit config`,
	Run: func(cmd *cobra.Command, args []string) {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			cmd.PrintErrf("Error initializing app: %v\n", err)
			return
		}

		err = ap.Services().Config().Edit()
		if err != nil {
			cmd.PrintErrf("Error editing config: %v\n", err)
			return
		}
	},
}

func init() {}
