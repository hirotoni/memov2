/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package todos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

var weeklyCmd = &cobra.Command{
	Use:   "weekly",
	Short: "generate weekly report for todos",
	Long:  `generate weekly report for todos`,
	Run: func(cmd *cobra.Command, args []string) {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			cmd.PrintErrf("Error initializing app: %v\n", err)
			return
		}

		// generate weekly report
		err = ap.Services().Todo().BuildWeeklyReportTodos()
		if err != nil {
			cmd.PrintErrf("Error generating weekly report: %v\n", err)
			return
		}
	},
}

func init() {}
