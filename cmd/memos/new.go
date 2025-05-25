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

// newCmd represents the memo command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "generate a new memo file",
	Long:  `generate a new memo file",`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.LoadTomlConfig()
		if err != nil {
			cmd.PrintErrf("Error loading config: %v\n", err)
			return
		}

		var title string
		if len(args) > 0 {
			title = args[0]
		}

		r := repository.NewRepositories(*c)
		uc := memo.NewMemo(*c, r)

		err = uc.GenerateMemoFile(title)
		if err != nil {
			cmd.PrintErrf("Error generating memo file: %v\n", err)
			return
		}
	},
}

func init() {}
