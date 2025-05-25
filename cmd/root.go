/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/hirotoni/memov2/cmd/config"
	"github.com/hirotoni/memov2/cmd/memos"
	"github.com/hirotoni/memov2/cmd/todos"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "memov2",
	Short: "memo v2",
	Long:  `memo v2`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(memos.MemosCmd)
	RootCmd.AddCommand(todos.TodosCmd)
	RootCmd.AddCommand(config.ConfigCmd)

	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
