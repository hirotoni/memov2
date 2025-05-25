/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/hirotoni/memov2/components"
	"github.com/hirotoni/memov2/config"
	"github.com/spf13/cobra"
)

// weeklyCmd represents the weekly command
var weeklyCmd = &cobra.Command{
	Use:   "weekly",
	Short: "generate weekly report",
	Long:  `generate weekly report`,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.LoadTomlConfig()
		err := components.BuildWeeklyReport(*c)
		if err != nil {
			cmd.PrintErrf("Error generating weekly report: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(weeklyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// weeklyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// weeklyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
