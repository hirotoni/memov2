package memos

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExplorerModel_Update tests the Update function with various messages
func TestExplorerModel_Update(t *testing.T) {
	tests := []struct {
		name         string
		initialMode  Mode
		msg          tea.Msg
		expectedMode Mode
		expectQuit   bool
		description  string
	}{
		{
			name:         "Tab key toggles from Browse to Search",
			initialMode:  BrowseMode,
			msg:          tea.KeyMsg{Type: tea.KeyTab},
			expectedMode: SearchMode,
			expectQuit:   false,
			description:  "Pressing Tab should switch from Browse to Search mode",
		},
		{
			name:         "Tab key toggles from Search to Browse",
			initialMode:  SearchMode,
			msg:          tea.KeyMsg{Type: tea.KeyTab},
			expectedMode: BrowseMode,
			expectQuit:   false,
			description:  "Pressing Tab should switch from Search to Browse mode",
		},
		{
			name:         "Q key triggers quit",
			initialMode:  BrowseMode,
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			expectedMode: BrowseMode,
			expectQuit:   true,
			description:  "Pressing 'q' should trigger quit command",
		},
		{
			name:         "Ctrl+C triggers quit",
			initialMode:  BrowseMode,
			msg:          tea.KeyMsg{Type: tea.KeyCtrlC},
			expectedMode: BrowseMode,
			expectQuit:   true,
			description:  "Pressing Ctrl+C should trigger quit command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			cfg, err := toml.NewConfig(toml.Option{
				BaseDir:         t.TempDir(), // Use temp dir for tests
				TodosFolderName: "todos/",
				MemosFolderName: "memos/",
			})
			require.NoError(t, err, "Failed to create config")

			editor := &mock.MockEditor{}

			m, err := NewIntegratedModel(cfg, editor)
			require.NoError(t, err, "Failed to create model")

			// Set initial mode
			m.SetMode(tt.initialMode)

			// Execute
			updatedModel, cmd := m.Update(tt.msg)

			// Assert
			explorerModel, ok := updatedModel.(*ExplorerModel)
			require.True(t, ok, "Updated model should be ExplorerModel type")

			assert.Equal(t, tt.expectedMode, explorerModel.CurrentMode(),
				"Mode should be %v but got %v", tt.expectedMode, explorerModel.CurrentMode())

			// Check if quit command was issued
			if tt.expectQuit {
				assert.NotNil(t, cmd, "Quit command should be returned")
				// Note: We can't directly compare cmd to tea.Quit, but we can check it's not nil
			} else if !tt.expectQuit && cmd != nil {
				// If we don't expect quit, cmd might still be not nil (batch commands, etc)
				// This is fine, we just check we got the expected mode
			}
		})
	}
}

// TestExplorerModel_Init tests the Init function
func TestExplorerModel_Init(t *testing.T) {
		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         t.TempDir(),
		TodosFolderName: "todos/",
		MemosFolderName: "memos/",
	})
	require.NoError(t, err, "Failed to create config")

	editor := &mock.MockEditor{}

	m, err := NewIntegratedModel(cfg, editor)
	require.NoError(t, err)

	cmd := m.Init()

	// Init should return nil or a valid command
	// For this model, Init returns nil
	assert.Nil(t, cmd, "Init should return nil for this model")
}

// TestExplorerModel_View tests the View rendering
func TestExplorerModel_View(t *testing.T) {
	tests := []struct {
		name        string
		mode        Mode
		setupModel  func(*ExplorerModel)
		checkView   func(*testing.T, string)
		description string
	}{
		{
			name: "View renders without error in BrowseMode",
			mode: BrowseMode,
			setupModel: func(m *ExplorerModel) {
				// Default setup, no modifications needed
			},
			checkView: func(t *testing.T, view string) {
				assert.NotEmpty(t, view, "View should not be empty")
				// The view should contain browse mode content
				// Note: Exact content depends on browseModel's View implementation
			},
			description: "View should render in Browse mode",
		},
		{
			name: "View renders without error in SearchMode",
			mode: SearchMode,
			setupModel: func(m *ExplorerModel) {
				// Default setup
			},
			checkView: func(t *testing.T, view string) {
				assert.NotEmpty(t, view, "View should not be empty")
				// The view should contain search mode content
			},
			description: "View should render in Search mode",
		},
		{
			name: "View shows error message when error is set",
			mode: BrowseMode,
			setupModel: func(m *ExplorerModel) {
				m.err = assert.AnError
			},
			checkView: func(t *testing.T, view string) {
				assert.Contains(t, view, "Error:", "View should contain error message")
				assert.Contains(t, view, "Press any key to continue", "View should show continuation prompt")
			},
			description: "View should display error when present",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			cfg, err := toml.NewConfig(toml.Option{
				BaseDir:         t.TempDir(),
				TodosFolderName: "todos/",
				MemosFolderName: "memos/",
			})
			require.NoError(t, err, "Failed to create config")

			editor := &mock.MockEditor{}

			m, err := NewIntegratedModel(cfg, editor)
			require.NoError(t, err)

			m.SetMode(tt.mode)
			tt.setupModel(m)

			// Execute
			view := m.View()

			// Assert
			tt.checkView(t, view)
		})
	}
}

// TestExplorerModel_ModeToggling tests mode toggling behavior
func TestExplorerModel_ModeToggling(t *testing.T) {
		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         t.TempDir(),
		TodosFolderName: "todos/",
		MemosFolderName: "memos/",
	})
	require.NoError(t, err, "Failed to create config")

	editor := &mock.MockEditor{}

	m, err := NewIntegratedModel(cfg, editor)
	require.NoError(t, err)

	// Start in Browse mode
	assert.Equal(t, BrowseMode, m.CurrentMode(), "Should start in Browse mode")

	// Toggle to Search
	m.toggleMode()
	assert.Equal(t, SearchMode, m.CurrentMode(), "Should switch to Search mode")

	// Toggle back to Browse
	m.toggleMode()
	assert.Equal(t, BrowseMode, m.CurrentMode(), "Should switch back to Browse mode")
}

// TestNewIntegratedModel_NilConfig tests error handling
func TestNewIntegratedModel_NilConfig(t *testing.T) {
	editor := &mock.MockEditor{}

	m, err := NewIntegratedModel(nil, editor)

	assert.Error(t, err, "Should return error with nil config")
	assert.Nil(t, m, "Model should be nil when error occurs")
	assert.Contains(t, err.Error(), "config cannot be nil", "Error message should be descriptive")
}

// Example of testing window resize
func TestExplorerModel_WindowResize(t *testing.T) {
		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         t.TempDir(),
		TodosFolderName: "todos/",
		MemosFolderName: "memos/",
	})
	require.NoError(t, err, "Failed to create config")

	editor := &mock.MockEditor{}

	m, err := NewIntegratedModel(cfg, editor)
	require.NoError(t, err)

	// Send window resize message
	resizeMsg := tea.WindowSizeMsg{
		Width:  100,
		Height: 50,
	}

	updatedModel, _ := m.Update(resizeMsg)

	// Verify the model was updated (exact behavior depends on implementation)
	assert.NotNil(t, updatedModel, "Model should be updated after resize")
}
