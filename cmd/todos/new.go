/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package todos

import (
	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/repository"
	"github.com/hirotoni/memov2/internal/usecase/todo"
	"github.com/spf13/cobra"
)

// newCmd represents the todo command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "generate a todo file",
	Long:  `generate a todo file`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.LoadTomlConfig()
		if err != nil {
			cmd.PrintErrf("Error loading config: %v\n", err)
			return
		}

		r := repository.NewRepositories(*c)
		uc := todo.NewTodo(*c, r)

		err = uc.GenerateTodoFile(truncateFlag)
		if err != nil {
			cmd.PrintErrf("Error generating todo file: %v\n", err)
			return
		}
	},
}

var truncateFlag bool

func init() {
	newCmd.Flags().BoolVarP(&truncateFlag, "truncate", "t", false, "Truncate the file if it exists")
}
