package memos

import (
	"strings"

	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

var searchShortFlag bool
var searchContextFlag bool

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "search memos with romaji-to-Japanese conversion",
	Long: `Search memos using romaji-to-Japanese conversion via SKK dictionary.
Output format is "title\tpath", same as the list command.

With --context flag, outputs match details in grep-like format:
  title\tpath\t[Type]\tmatch_content

Example usage with fzf:
  memov2 memos list | fzf \
    --disabled \
    --bind "change:reload:memov2 memos search {q}" \
  | cut -f2 | xargs memov2 memos open

Example usage with fzf and --context:
  memov2 memos list | fzf \
    --disabled \
    --delimiter=$'\t' \
    --with-nth='1,3,4' \
    --bind "change:reload:memov2 memos search --context {q}" \
  | cut -f2 | xargs memov2 memos open`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			cmd.PrintErrf("Error initializing app: %v\n", err)
			return
		}
		query := strings.Join(args, " ")
		err = ap.Services().Memo().Search(query, !searchShortFlag, searchContextFlag)
		if err != nil {
			cmd.PrintErrf("Error searching memos: %v\n", err)
			return
		}
	},
}

func init() {
	searchCmd.Flags().BoolVarP(&searchShortFlag, "short", "s", false, "show relative paths instead of full paths")
	searchCmd.Flags().BoolVarP(&searchContextFlag, "context", "c", false, "show match context (type and content) for each match")
}
