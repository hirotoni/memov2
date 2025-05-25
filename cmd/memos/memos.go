/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package memos

import (
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var MemosCmd = &cobra.Command{
	Use:   "memos",
	Short: "commands about memos",
	Long:  `commands about memos`,
}

func init() {
	MemosCmd.AddCommand(newCmd)
	MemosCmd.AddCommand(weeklyCmd)
	MemosCmd.AddCommand(indexCmd)
	MemosCmd.AddCommand(browseCmd)
	MemosCmd.AddCommand(listCmd)
	MemosCmd.AddCommand(openCmd)
	MemosCmd.AddCommand(searchCmd)
	MemosCmd.AddCommand(renameCmd)
}
