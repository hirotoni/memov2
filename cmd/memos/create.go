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

// createCmd creates a memo in one of two modes, mirroring the dual-mode `rename`
// command. With no --title it launches the interactive TUI + editor flow (the
// former `memos new`, retained as an alias). With --title it creates the memo
// non-interactively (no TUI, no editor): category from --category, body from
// stdin; by default it prints the new memo's relative path, or a JSON object with
// --json. It uses RunE + SilenceUsage so invalid input exits non-zero with nothing
// on stdout.
var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "create a memo (interactively, or non-interactively with --title)",
	Long: `Create a memo. With no --title, pick a category in an embedded TUI, type a title, and open the new memo in the editor (the same flow as the "new" alias). With --title, create non-interactively for AI/automation: --category is an optional path (a/b/c) and the body is read from stdin (a TTY or empty stdin creates a memo with no body). The non-interactive mode prints the new memo's relative path, or a JSON object with --json. Titles that escape the memos directory are rejected.`,
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			return err
		}

		// Rule: a blank --title (empty or whitespace-only) always means interactive.
		// We never fall through to a non-interactive create with a blank title, so an
		// untitled memo can never be created silently. Stdin is irrelevant here — the
		// TUI owns the terminal in this mode (matching the former `memos new`).
		if strings.TrimSpace(createTitleFlag) == "" {
			return ap.Services().Memo().NewInteractive()
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
	createCmd.Flags().StringVar(&createTitleFlag, "title", "", "memo title; when omitted, the interactive TUI is launched")
	createCmd.Flags().StringVar(&createCategoryFlag, "category", "", "category path, e.g. a/b/c (non-interactive mode)")
	createCmd.Flags().BoolVar(&createJSONFlag, "json", false, "output as JSON to stdout (non-interactive mode)")
}
