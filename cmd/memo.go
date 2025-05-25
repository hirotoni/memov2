/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/hirotoni/memov2/components"
	"github.com/hirotoni/memov2/config"
	"github.com/spf13/cobra"
)

// memoCmd represents the memo command
var memoCmd = &cobra.Command{
	Use:   "memo",
	Short: "generate memo file",
	Long:  `generate memo file`,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.LoadTomlConfig()

		var title string
		if len(args) > 0 {
			title = args[0]
		}

		err := components.GenerateMemoFile(*c, title)
		if err != nil {
			cmd.PrintErrf("Error generating memo file: %v\n", err)
			return
		}
	},
}

func init() {
	newCmd.AddCommand(memoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// memoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// memoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
