/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/hirotoni/memov2/components"
	"github.com/hirotoni/memov2/config"
	"github.com/spf13/cobra"
)

// todoCmd represents the todo command
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "generate todo file",
	Long:  `generate todo file`,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.LoadTomlConfig()

		err := components.GenerateTodoFile(*c, truncateFlag)
		if err != nil {
			cmd.PrintErrf("Error generating todo file: %v\n", err)
			return
		}
	},
}

var truncateFlag bool

func init() {
	newCmd.AddCommand(todoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// todoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	todoCmd.Flags().BoolVarP(&truncateFlag, "truncate", "t", false, "Truncate the file if it exists")
}
