package config

import (
	"fmt"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show config",
	Long:  `show config`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.LoadTomlConfig()
		if err != nil {
			cmd.PrintErrf("Error loading config: %v\n", err)
			return
		}

		fmt.Printf("base_dir: %s\n", c.BaseDir())
		fmt.Printf("todos_dir: %s\n", c.TodosDir())
		fmt.Printf("memos_dir: %s\n", c.MemosDir())
		fmt.Printf("todos_daystoseek: %d\n", c.TodosDaysToSeek())
	},
}

func init() {}
