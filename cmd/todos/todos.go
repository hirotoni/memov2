/*
Copyright © 2025 hirotoni
*/
package todos

import (
	"github.com/spf13/cobra"
)

// TodosCmd represents the todos command
var TodosCmd = &cobra.Command{
	Use:   "todos",
	Short: "commands about todos",
	Long:  `commands about todos`,
}

func init() {
	TodosCmd.AddCommand(createCmd)
	TodosCmd.AddCommand(weeklyCmd)
	TodosCmd.AddCommand(listCmd)
	TodosCmd.AddCommand(showCmd)
}
