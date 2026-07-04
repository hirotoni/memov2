package memos

import (
	"io"
	"os"
	"strings"

	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	createTitleFlag    string
	createCategoryFlag string
	createJSONFlag     bool
)

// createCmd creates a memo non-interactively (no TUI, no editor). It is the write
// counterpart of the read-only --json commands: the title comes from --title, the
// category from --category (a path like "a/b/c"), and the body is read from stdin.
// By default it prints the new memo's relative path (one line); with --json it
// prints a structured object. It uses RunE + SilenceUsage so an invalid input
// yields a non-zero exit with nothing on stdout.
var createCmd = &cobra.Command{
	Use:           "create",
	Short:         "create a memo non-interactively (for AI/automation)",
	Long:          `Create a memo without launching the TUI or the editor. Title is required (--title); category is an optional path (--category a/b/c). The body is read from stdin (a TTY or empty stdin creates a memo with no body). Prints the new memo's relative path, or a JSON object with --json.`,
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			return err
		}

		var categoryTree []string
		if trimmed := strings.Trim(createCategoryFlag, "/"); trimmed != "" {
			categoryTree = strings.Split(trimmed, "/")
		}

		body := readStdinBody()

		return ap.Services().Memo().CreateNonInteractive(createTitleFlag, categoryTree, body, createJSONFlag)
	},
}

// readStdinBody reads the memo body from stdin. If stdin is an interactive
// terminal (no piped input) it returns an empty string instead of blocking.
func readStdinBody() string {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		return ""
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return ""
	}
	return string(data)
}

func init() {
	createCmd.Flags().StringVar(&createTitleFlag, "title", "", "memo title (required)")
	createCmd.Flags().StringVar(&createCategoryFlag, "category", "", "category path, e.g. a/b/c")
	createCmd.Flags().BoolVar(&createJSONFlag, "json", false, "output as JSON to stdout")
}
