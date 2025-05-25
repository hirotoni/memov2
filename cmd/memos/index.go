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

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "generate memo index file",
	Long:  `generate a memo index file that lists all memos.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.LoadTomlConfig()
		if err != nil {
			cmd.PrintErrf("Error loading config: %v\n", err)
			return
		}

		r := repository.NewRepositories(*c)
		uc := memo.NewMemo(*c, r)

		err = uc.GenerateMemoIndex()
		if err != nil {
			cmd.PrintErrf("Error generating memo index: %v\n", err)
			return
		}
	},
}

func init() {}
