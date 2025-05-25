/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package todos

import (
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var TodosCmd = &cobra.Command{
	Use:   "todos",
	Short: "commands about todos",
	Long:  `commands about todos`,
}

func init() {
	TodosCmd.AddCommand(newCmd)
	TodosCmd.AddCommand(weeklyCmd)
}
