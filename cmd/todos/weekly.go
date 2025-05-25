/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package todos

import (
	"github.com/hirotoni/memov2/components/todos"
	"github.com/hirotoni/memov2/config"
	"github.com/spf13/cobra"
)

var weeklyCmd = &cobra.Command{
	Use:   "weekly",
	Short: "generate weekly report for todos",
	Long:  `generate weekly report for todos`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.LoadTomlConfig()
		if err != nil {
			cmd.PrintErrf("Error loading config: %v\n", err)
			return
		}
		err = todos.BuildWeeklyReportTodos(*c)
		if err != nil {
			cmd.PrintErrf("Error generating weekly report: %v\n", err)
			return
		}
	},
}

func init() {}
