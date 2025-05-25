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

// TestExplorerModel_MessagePropagation tests that messages are properly propagated to sub-models
func TestExplorerModel_MessagePropagation(t *testing.T) {
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         t.TempDir(),
		TodosFolderName: "todos/",
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}

	m, err := NewIntegratedModel(cfg, editor)
	require.NoError(t, err)

	// Test that messages are passed to browse model when in BrowseMode
	m.SetMode(BrowseMode)
	updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.NotNil(t, updatedModel)
	// Command might be nil or not, depending on implementation
	_ = cmd

	// Test that messages are passed to search model when in SearchMode
	m.SetMode(SearchMode)
	updatedModel, cmd = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.NotNil(t, updatedModel)
	_ = cmd
}

// TestExplorerModel_ErrorHandling tests error handling in the integrated model
func TestExplorerModel_ErrorHandling(t *testing.T) {
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         t.TempDir(),
		TodosFolderName: "todos/",
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}

	m, err := NewIntegratedModel(cfg, editor)
	require.NoError(t, err)

	// Set an error
	m.err = assert.AnError

	// View should show error
	view := m.View()
	assert.Contains(t, view, "Error:", "View should show error message")

	// Update should still work
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	assert.NotNil(t, updatedModel)
}

// TestExplorerModel_ComplexInteractionFlow tests a complex interaction flow
func TestExplorerModel_ComplexInteractionFlow(t *testing.T) {
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         t.TempDir(),
		TodosFolderName: "todos/",
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}

	m, err := NewIntegratedModel(cfg, editor)
	require.NoError(t, err)

	// Start in Browse mode
	assert.Equal(t, BrowseMode, m.CurrentMode())

	// Switch to Search mode
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updatedModel.(*ExplorerModel)
	assert.Equal(t, SearchMode, m.CurrentMode())

	// Switch back to Browse mode
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updatedModel.(*ExplorerModel)
	assert.Equal(t, BrowseMode, m.CurrentMode())

	// Resize window
	updatedModel, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	assert.NotNil(t, updatedModel)

	// Quit
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.NotNil(t, cmd)
}

// TestExplorerModel_CurrentMode tests CurrentMode and SetMode
func TestExplorerModel_CurrentMode(t *testing.T) {
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         t.TempDir(),
		TodosFolderName: "todos/",
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}

	m, err := NewIntegratedModel(cfg, editor)
	require.NoError(t, err)

	// Default mode should be BrowseMode
	assert.Equal(t, BrowseMode, m.CurrentMode())

	// Set to SearchMode
	m.SetMode(SearchMode)
	assert.Equal(t, SearchMode, m.CurrentMode())

	// Set back to BrowseMode
	m.SetMode(BrowseMode)
	assert.Equal(t, BrowseMode, m.CurrentMode())
}

// TestExplorerModel_ViewModeSwitching tests that View correctly switches between modes
func TestExplorerModel_ViewModeSwitching(t *testing.T) {
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         t.TempDir(),
		TodosFolderName: "todos/",
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}

	m, err := NewIntegratedModel(cfg, editor)
	require.NoError(t, err)

	// View in BrowseMode
	m.SetMode(BrowseMode)
	viewBrowse := m.View()
	assert.NotEmpty(t, viewBrowse)

	// View in SearchMode
	m.SetMode(SearchMode)
	viewSearch := m.View()
	assert.NotEmpty(t, viewSearch)

	// Views should be different (though we can't guarantee exact content)
	// At minimum, they should both be non-empty
	assert.NotEqual(t, viewBrowse, viewSearch, "Views should differ between modes")
}

// TestExplorerModel_QuitFromBothModes tests quit functionality from both modes
func TestExplorerModel_QuitFromBothModes(t *testing.T) {
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         t.TempDir(),
		TodosFolderName: "todos/",
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}

	tests := []struct {
		name      string
		mode      Mode
		keyMsg    tea.KeyMsg
		expectQuit bool
	}{
		{
			name:       "Quit with 'q' from BrowseMode",
			mode:       BrowseMode,
			keyMsg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			expectQuit: true,
		},
		{
			name:       "Quit with Ctrl+C from BrowseMode",
			mode:       BrowseMode,
			keyMsg:     tea.KeyMsg{Type: tea.KeyCtrlC},
			expectQuit: true,
		},
		{
			name:       "Quit with 'q' from SearchMode",
			mode:       SearchMode,
			keyMsg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			expectQuit: true,
		},
		{
			name:       "Quit with Ctrl+C from SearchMode",
			mode:       SearchMode,
			keyMsg:     tea.KeyMsg{Type: tea.KeyCtrlC},
			expectQuit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewIntegratedModel(cfg, editor)
			require.NoError(t, err)

			m.SetMode(tt.mode)
			_, cmd := m.Update(tt.keyMsg)

			if tt.expectQuit {
				assert.NotNil(t, cmd, "Quit command should be returned")
			}
		})
	}
}

// TestExplorerModel_MultipleModeToggles tests rapid mode toggling
func TestExplorerModel_MultipleModeToggles(t *testing.T) {
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         t.TempDir(),
		TodosFolderName: "todos/",
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}

	m, err := NewIntegratedModel(cfg, editor)
	require.NoError(t, err)

	// Start in BrowseMode (default)
	currentMode := BrowseMode

	// Toggle multiple times
	for i := 0; i < 10; i++ {
		// Toggle mode
		updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
		m = updatedModel.(*ExplorerModel)
		
		// Expected mode after toggle
		if currentMode == BrowseMode {
			currentMode = SearchMode
		} else {
			currentMode = BrowseMode
		}
		
		assert.Equal(t, currentMode, m.CurrentMode(), "Mode should toggle correctly on iteration %d", i)
	}
}

// TestExplorerModel_NonKeyMessages tests handling of non-key messages
func TestExplorerModel_NonKeyMessages(t *testing.T) {
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         t.TempDir(),
		TodosFolderName: "todos/",
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}

	m, err := NewIntegratedModel(cfg, editor)
	require.NoError(t, err)

	// Test with WindowSizeMsg
	resizeMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := m.Update(resizeMsg)
	assert.NotNil(t, updatedModel)

	// Test with error message
	errorMsg := assert.AnError
	updatedModel, _ = m.Update(errorMsg)
	assert.NotNil(t, updatedModel)
}
