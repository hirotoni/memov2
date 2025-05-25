/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package memos

import (
	"github.com/hirotoni/memov2/components/memos"
	"github.com/hirotoni/memov2/config"
	"github.com/spf13/cobra"
)

var weeklyCmd = &cobra.Command{
	Use:   "weekly",
	Short: "generate weekly report for memos",
	Long:  `generate weekly report for memos`,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.LoadTomlConfig()
		err := memos.BuildWeeklyReportMemos(*c)
		if err != nil {
			cmd.PrintErrf("Error generating weekly report: %v\n", err)
			return
		}
	},
}

func init() {}
