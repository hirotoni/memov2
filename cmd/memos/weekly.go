/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package memos

import (
	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/repository"
	"github.com/hirotoni/memov2/internal/usecase/memo"
	"github.com/spf13/cobra"
)

var weeklyCmd = &cobra.Command{
	Use:   "weekly",
	Short: "generate weekly report for memos",
	Long:  `generate weekly report for memos`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.LoadTomlConfig()
		if err != nil {
			cmd.PrintErrf("Error loading config: %v\n", err)
			return
		}

		r := repository.NewRepositories(*c)
		uc := memo.NewMemo(*c, r)

		err = uc.BuildWeeklyReportMemos()
		if err != nil {
			cmd.PrintErrf("Error generating weekly report: %v\n", err)
			return
		}
	},
}

func init() {}
