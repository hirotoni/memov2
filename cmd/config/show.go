package config

import (
	"encoding/json"
	"fmt"

	"github.com/hirotoni/memov2/config"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show config",
	Long:  `show config`,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.LoadTomlConfig()

		// pretty print the config
		prettyConfig, err := json.MarshalIndent(c, "", "  ")
		if err != nil {
			cmd.PrintErrf("Error pretty printing config: %v\n", err)
			return
		}
		fmt.Println(string(prettyConfig))
	},
}

func init() {}
