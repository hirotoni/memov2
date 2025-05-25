package memo

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/ui/tui/memos/picker"
	"github.com/hirotoni/memov2/internal/ui/tui/memos/search"
)

// SearchInteractive launches the standalone search TUI (romaji-aware). Selecting
// a result opens it in the configured editor. Backs the `memos search` command.
func (uc memo) SearchInteractive() error {
	cfg := uc.config.GetTomlConfig().(*toml.Config)
	if err := cfg.EnsureDirectories(); err != nil {
		return common.Wrap(err, common.ErrorTypeService, "failed to ensure directories")
	}

	m, err := search.NewStandalone(cfg, uc.editor)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error creating search model")
	}
	final, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error running search TUI")
	}

	// The TUI has exited (alt screen torn down). Open the chosen file now so the
	// editor runs in the restored terminal and control returns to the shell when
	// it closes.
	sm, ok := final.(search.Model)
	if !ok || sm.SelectedPath() == "" {
		return nil // nothing selected / cancelled
	}
	if err := uc.editor.Open(cfg.BaseDir(), sm.SelectedPath()); err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error opening editor")
	}
	return nil
}

// RenameInteractive lets the user pick a memo and enter a new title in a TUI,
// then renames it. Backs the `memos rename` command.
func (uc memo) RenameInteractive() error {
	relPath, newTitle, ok, err := picker.SelectMemoForRename(uc.repos.Memo())
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error selecting memo to rename")
	}
	if !ok {
		return nil // cancelled
	}
	return uc.Rename(relPath, newTitle)
}

// NewInteractive lets the user pick a category and enter a title in a TUI, then
// creates the memo. Backs the `memos new` command.
func (uc memo) NewInteractive() error {
	tree, title, ok, err := picker.SelectCategoryForNew(uc.repos.Memo())
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error selecting category")
	}
	if !ok {
		return nil // cancelled
	}
	return uc.GenerateMemoFile(title, tree)
}
