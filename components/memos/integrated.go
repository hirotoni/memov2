package memos

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hirotoni/memov2/components/memos/browse"
	"github.com/hirotoni/memov2/components/memos/search"
	"github.com/hirotoni/memov2/config"
	"github.com/hirotoni/memov2/utils"
)

// Mode represents the current mode of the TUI
type Mode int

const (
	BrowseMode Mode = iota
	SearchMode
)

type explorerModel struct {
	// Common
	mode Mode
	// Browse mode
	browseModel browse.BrowseModel
	// Search mode
	searchModel search.Model
}

func NewIntegratedModel(c *config.TomlConfig) (*explorerModel, error) {
	searchModel, err := search.New(c)
	if err != nil {
		return nil, fmt.Errorf("error creating search model: %w", err)
	}
	browseModel, err := browse.New(c)
	if err != nil {
		return nil, fmt.Errorf("error creating browse model: %w", err)
	}

	m := &explorerModel{
		mode:        BrowseMode,
		browseModel: *browseModel,
		searchModel: *searchModel,
	}

	return m, nil
}

func (m explorerModel) Init() tea.Cmd {
	return nil
}

func (m explorerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// global keybindings
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "tab":
			if m.mode == BrowseMode {
				m.mode = SearchMode
			} else {
				m.mode = BrowseMode
			}
		}
	}

	if m.mode == BrowseMode {
		mm, cmd := m.browseModel.Update(msg)
		m.browseModel = mm.(browse.BrowseModel)
		cmds = append(cmds, cmd)
	} else {
		mm, cmd := m.searchModel.Update(msg)
		m.searchModel = mm.(search.Model)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m explorerModel) View() string {
	if m.mode == BrowseMode {
		return m.browseModel.View()
	} else {
		return m.searchModel.View()
	}
}

func IntegratedMemos(c *config.TomlConfig) error {
	if !utils.Exists(c.MemosDir()) {
		if err := os.MkdirAll(c.MemosDir(), 0755); err != nil {
			return fmt.Errorf("failed to create memos directory %s: %w", c.MemosDir(), err)
		}
		fmt.Printf("Created memos directory: %s\n", c.MemosDir())
	}

	m, err := NewIntegratedModel(c)
	if err != nil {
		return fmt.Errorf("error creating integrated model: %w", err)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}
