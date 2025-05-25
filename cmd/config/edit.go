package config

import (
	"os"
	"os/exec"

	"github.com/hirotoni/memov2/config"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit config",
	Long:  `edit config`,
	Run: func(cmd *cobra.Command, args []string) {
		_, path, err := config.ConfigDirPath()
		if err != nil {
			cmd.PrintErrf("Error getting config path: %v\n", err)
			return
		}
		editorCmd := exec.Command("vim", path)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr
		err = editorCmd.Run()
		if err != nil {
			cmd.PrintErrf("Error running editor: %v\n", err)
			return
		}
	},
}

func init() {}
