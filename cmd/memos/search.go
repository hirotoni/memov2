package memos

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hirotoni/memov2/components/memos/search"
	"github.com/hirotoni/memov2/config"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search through memos",
	Long:  `Interactive search through memos with support for title, category, content, and heading matches.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.LoadTomlConfig()
		m := search.New(c)
		p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

		if _, err := p.Run(); err != nil {
			cmd.PrintErrf("Error running search interface: %v\n", err)
			return
		}
	},
}

func init() {
	MemosCmd.AddCommand(searchCmd)
}
