package memos

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hirotoni/memov2/cmd/app"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/platform"
	"github.com/spf13/cobra"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename <path> [new-title]",
	Short: "rename a memo file",
	Long: `Rename a memo file. The path is relative to the memos directory (same format as list/search output).

If new-title is omitted, prompts interactively for the new title.

Example usage with fzf:
  memov2 memos list | fzf | cut -f2 | xargs memov2 memos rename`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			cmd.PrintErrf("Error initializing app: %v\n", err)
			return
		}

		path := args[0]
		var newTitle string

		if len(args) >= 2 {
			// Direct mode: new title provided as argument
			newTitle = strings.Join(args[1:], " ")
		} else {
			// Interactive mode: prompt for new title
			fileName := filepath.Base(path)
			currentTitle := strings.ReplaceAll(domain.MemoTitle(fileName), "-", " ")
			fmt.Fprintf(cmd.ErrOrStderr(), "Current title: %s\n", currentTitle)

			newTitle, err = platform.ReadLine("New title: ")
			if err != nil {
				cmd.PrintErrf("Error reading input: %v\n", err)
				return
			}
			newTitle = strings.TrimSpace(newTitle)
			if newTitle == "" {
				cmd.PrintErrln("Cancelled: empty title")
				return
			}
		}

		err = ap.Services().Memo().Rename(path, newTitle)
		if err != nil {
			cmd.PrintErrf("Error renaming memo: %v\n", err)
			return
		}
	},
}
