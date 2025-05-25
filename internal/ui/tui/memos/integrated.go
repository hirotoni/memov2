package memos

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/ui/tui/memos/browse"
	"github.com/hirotoni/memov2/internal/ui/tui/memos/search"
)

// Mode represents the current mode of the TUI
type Mode int

const (
	BrowseMode Mode = iota
	SearchMode
)

// ExplorerModel represents the integrated memo explorer model
type ExplorerModel struct {
	// Common
	mode Mode
	// Browse mode
	browseModel browse.BrowseModel
	// Search mode
	searchModel search.Model
	// Configuration
	config *toml.Config
	// Error handling
	err error
}

// NewIntegratedModel creates a new integrated memo explorer model
func NewIntegratedModel(c *toml.Config, e interfaces.Editor) (*ExplorerModel, error) {
	if c == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	searchModel, err := search.New(c, e)
	if err != nil {
		return nil, fmt.Errorf("error creating search model: %w", err)
	}

	browseModel, err := browse.New(c, e)
	if err != nil {
		return nil, fmt.Errorf("error creating browse model: %w", err)
	}

	m := &ExplorerModel{
		mode:        BrowseMode,
		browseModel: *browseModel,
		searchModel: *searchModel,
		config:      c,
	}

	return m, nil
}

// Init initializes the model
func (m *ExplorerModel) Init() tea.Cmd {
	return nil
}

// Update handles model updates
func (m *ExplorerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle global keybindings
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "tab":
			m.toggleMode()
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	// Update the current mode's model
	if m.mode == BrowseMode {
		updatedModel, cmd := m.browseModel.Update(msg)
		if browseModel, ok := updatedModel.(browse.BrowseModel); ok {
			m.browseModel = browseModel
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	} else {
		updatedModel, cmd := m.searchModel.Update(msg)
		if searchModel, ok := updatedModel.(search.Model); ok {
			m.searchModel = searchModel
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the current view
func (m *ExplorerModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress any key to continue...", m.err)
	}

	if m.mode == BrowseMode {
		return m.browseModel.View()
	}
	return m.searchModel.View()
}

// toggleMode switches between browse and search modes
func (m *ExplorerModel) toggleMode() {
	if m.mode == BrowseMode {
		m.mode = SearchMode
	} else {
		m.mode = BrowseMode
	}
}

// CurrentMode returns the current mode
func (m *ExplorerModel) CurrentMode() Mode {
	return m.mode
}

// SetMode sets the current mode
func (m *ExplorerModel) SetMode(mode Mode) {
	m.mode = mode
}

// IntegratedMemos runs the integrated memo explorer
func IntegratedMemos(c *toml.Config, e interfaces.Editor) error {
	if c == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Ensure directories exist
	if err := c.EnsureDirectories(); err != nil {
		return fmt.Errorf("failed to ensure directories: %w", err)
	}

	// Create the integrated model
	m, err := NewIntegratedModel(c, e)
	if err != nil {
		return fmt.Errorf("error creating integrated model: %w", err)
	}

	// Create and run the program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}
