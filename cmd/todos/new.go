/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package todos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

// newCmd represents the todo command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "generate a todo file",
	Long:  `generate a todo file`,
	Run: func(cmd *cobra.Command, args []string) {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			cmd.PrintErrf("Error initializing app: %v\n", err)
			return
		}

		// generate todo
		err = ap.Services().Todo().GenerateTodoFile(truncateFlag)
		if err != nil {
			cmd.PrintErrf("Error generating todo file: %v\n", err)
			return
		}
	},
}

var truncateFlag bool

func init() {
	newCmd.Flags().BoolVarP(&truncateFlag, "truncate", "t", false, "Truncate the file if it exists")
}
