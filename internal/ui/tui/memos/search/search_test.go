package search

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/repositories/memo"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	memsearch "github.com/hirotoni/memov2/internal/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew tests the creation of a search Model
func TestNew(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}

	model, err := New(cfg, editor)
	require.NoError(t, err)
	assert.NotNil(t, model)
	assert.Equal(t, cfg, model.config)
	assert.Equal(t, editor, model.editor)
	assert.NotNil(t, model.romajiConv)
	assert.Equal(t, focusInput, model.focus)
}

// TestModel_Init tests the Init function
func TestModel_Init(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	cmd := model.Init()
	assert.NotNil(t, cmd, "Init should return a command")
}

// TestModel_View tests the View function
func TestModel_View(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Test view with focus on input
	model.focus = focusInput
	view := model.View()
	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Search memos", "View should contain search input")

	// Test view with focus on list
	model.focus = focusList
	view = model.View()
	assert.NotEmpty(t, view)
}

// TestModel_FocusManagement tests focus switching
func TestModel_FocusManagement(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	modelPtr, err := New(cfg, editor)
	require.NoError(t, err)
	model := *modelPtr

	// Start with focus on input
	assert.Equal(t, focusInput, model.focus)

	// Switch to list with Ctrl+J
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlJ})
	if updatedModel, ok := updatedModel.(Model); ok {
		model = updatedModel
		assert.Equal(t, focusList, model.focus, "Focus should switch to list")
	}

	// Switch back to input with Ctrl+K
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlK})
	if updatedModel, ok := updatedModel.(Model); ok {
		model = updatedModel
		assert.Equal(t, focusInput, model.focus, "Focus should switch back to input")
	}

	// Switch to list with Enter when there are results
	testMemo, err := domain.NewMemoFile(time.Now(), "test-memo", []string{"cat1"})
	require.NoError(t, err)
	model.results = []memsearch.SearchResult{
		{Memo: testMemo, Matches: []memsearch.Match{}},
	}
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if updatedModel, ok := updatedModel.(Model); ok {
		model = updatedModel
		assert.Equal(t, focusList, model.focus, "Focus should switch to list on Enter with results")
	}
}

// TestModel_EscKeyBehavior tests Esc key behavior
func TestModel_EscKeyBehavior(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	modelPtr, err := New(cfg, editor)
	require.NoError(t, err)
	model := *modelPtr

	// Esc from input should quit
	model.focus = focusInput
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	assert.NotNil(t, cmd, "Esc from input should trigger quit")

	// Esc from list should return to input
	model.focus = focusList
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if updatedModel, ok := updatedModel.(Model); ok {
		model = updatedModel
		assert.Nil(t, cmd, "Esc from list should not quit")
		assert.Equal(t, focusInput, model.focus, "Esc from list should return focus to input")
	}
}

// TestModel_NavigationKeys tests navigation in list mode
func TestModel_NavigationKeys(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	modelPtr, err := New(cfg, editor)
	require.NoError(t, err)
	model := *modelPtr

	// Create test memos for results (renderResults needs valid Memo objects)
	memo1, err := domain.NewMemoFile(time.Now(), "test-memo-1", []string{"cat1"})
	require.NoError(t, err)
	memo2, err := domain.NewMemoFile(time.Now(), "test-memo-2", []string{"cat2"})
	require.NoError(t, err)
	memo3, err := domain.NewMemoFile(time.Now(), "test-memo-3", []string{"cat3"})
	require.NoError(t, err)

	// Set up some results
	model.results = []memsearch.SearchResult{
		{Memo: memo1, Matches: []memsearch.Match{}},
		{Memo: memo2, Matches: []memsearch.Match{}},
		{Memo: memo3, Matches: []memsearch.Match{}},
	}
	model.focus = focusList
	model.selected = 1

	// Test up arrow
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m, ok := updatedModel.(Model); ok {
		model = m
		assert.Equal(t, 0, model.selected, "Up arrow should move selection up")
	}

	// Test down arrow
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m, ok := updatedModel.(Model); ok {
		model = m
		assert.Equal(t, 1, model.selected, "Down arrow should move selection down")
	}

	// Test 'k' key
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m, ok := updatedModel.(Model); ok {
		model = m
		assert.Equal(t, 0, model.selected, "'k' should move selection up")
	}

	// Test 'j' key
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m, ok := updatedModel.(Model); ok {
		model = m
		assert.Equal(t, 1, model.selected, "'j' should move selection down")
	}

	// Test Ctrl+U (move up by 5)
	model.selected = 5
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	if m, ok := updatedModel.(Model); ok {
		model = m
		assert.Equal(t, 0, model.selected, "Ctrl+U should move up by 5 (clamped to 0)")
	}

	// Test Ctrl+D (move down by 5)
	model.selected = 0
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
	if m, ok := updatedModel.(Model); ok {
		model = m
		assert.Equal(t, 2, model.selected, "Ctrl+D should move down by 5 (clamped to max)")
	}

	// Test 'G' (go to bottom) - use KeyRunes with 'G'
	model.selected = 0
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	if m, ok := updatedModel.(Model); ok {
		model = m
		assert.Equal(t, 2, model.selected, "'G' should go to bottom")
	}

	// Test 'gg' (go to top)
	model.selected = 2
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	if m, ok := updatedModel.(Model); ok {
		model = m
	}
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	if m, ok := updatedModel.(Model); ok {
		model = m
		assert.Equal(t, 0, model.selected, "'gg' should go to top")
	}
}

// TestModel_OpenMemo tests opening a memo from search results
func TestModel_OpenMemo(t *testing.T) {
	tempDir := t.TempDir()
	memosDir := filepath.Join(tempDir, "memos")
	require.NoError(t, os.MkdirAll(memosDir, 0755))

	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	// Create a test memo
	logger := common.DefaultLogger()
	repo := memo.NewMemo(cfg.MemosDir(), logger)
	testMemo, err := domain.NewMemoFile(time.Now(), "test-memo", []string{"category1"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(testMemo, true))

	editor := &mock.MockEditor{}
	modelPtr, err := New(cfg, editor)
	require.NoError(t, err)
	model := *modelPtr

	// Set up results with the test memo
	model.results = []memsearch.SearchResult{
		{Memo: testMemo, Matches: []memsearch.Match{}},
	}
	model.focus = focusList
	model.selected = 0

	// Test opening with Enter
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.NotNil(t, updatedModel)
	// Editor should have been called (checked via mock if needed)

	// Test opening with 'l'
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	assert.NotNil(t, updatedModel)
}

// TestModel_WindowResize tests window resize handling
func TestModel_WindowResize(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	modelPtr, err := New(cfg, editor)
	require.NoError(t, err)
	model := *modelPtr

	// Resize window
	resizeMsg := tea.WindowSizeMsg{
		Width:  120,
		Height: 40,
	}

	updatedModel, _ := model.Update(resizeMsg)
	if updatedModel, ok := updatedModel.(Model); ok {
		model = updatedModel
		assert.Equal(t, 120, model.width)
		assert.Equal(t, 40, model.height)
	}
}

// TestModel_SearchTrigger tests that search is triggered on input change
func TestModel_SearchTrigger(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	modelPtr, err := New(cfg, editor)
	require.NoError(t, err)
	model := *modelPtr

	// Type in search input
	model.focus = focusInput
	model.searchInput.Focus()

	// Simulate typing
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	assert.NotNil(t, updatedModel)
	// Command should be returned (for search)
	if cmd != nil {
		// Execute the command to get search results
		msg := cmd()
		if searchMsg, ok := msg.(searchResultMsg); ok {
			updatedModel, _ = model.Update(searchMsg)
			assert.NotNil(t, updatedModel)
		}
	}
}

// TestModel_EmptySearch tests behavior with empty search
func TestModel_EmptySearch(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	modelPtr, err := New(cfg, editor)
	require.NoError(t, err)
	model := *modelPtr

	// Empty search should return empty results
	results := model.search()
	searchMsg, ok := results.(searchResultMsg)
	require.True(t, ok)
	assert.Empty(t, searchMsg.results, "Empty search should return no results")
}

// TestModel_ViewWithResults tests view rendering with search results
func TestModel_ViewWithResults(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	modelPtr, err := New(cfg, editor)
	require.NoError(t, err)
	model := *modelPtr

	// Create test memo for results
	testMemo, err := domain.NewMemoFile(time.Now(), "test-memo", []string{"cat1"})
	require.NoError(t, err)

	// Set up results
	model.results = []memsearch.SearchResult{
		{
			Memo: testMemo,
			Matches: []memsearch.Match{
				{Type: memsearch.MatchTitle, Content: "Test Memo"},
			},
		},
	}
	model.focus = focusList

	view := model.View()
	assert.NotEmpty(t, view)
}

// TestModel_Quit tests quit functionality
func TestModel_Quit(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Ctrl+C should quit
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.NotNil(t, cmd, "Ctrl+C should trigger quit")
}

// TestModel_ErrorHandling tests error handling
func TestModel_ErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Set an error
	model.err = assert.AnError

	// View should still render (error handling depends on implementation)
	view := model.View()
	assert.NotEmpty(t, view)
}

// TestModel_LastKeyGTracking tests the 'gg' command tracking
func TestModel_LastKeyGTracking(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	modelPtr, err := New(cfg, editor)
	require.NoError(t, err)
	model := *modelPtr

	// Create test memos for results (renderResults needs valid Memo objects)
	memo1, err := domain.NewMemoFile(time.Now(), "test-memo-1", []string{"cat1"})
	require.NoError(t, err)
	memo2, err := domain.NewMemoFile(time.Now(), "test-memo-2", []string{"cat2"})
	require.NoError(t, err)

	model.results = []memsearch.SearchResult{
		{Memo: memo1, Matches: []memsearch.Match{}},
		{Memo: memo2, Matches: []memsearch.Match{}},
	}
	model.focus = focusList
	model.selected = 1

	// First 'g' should set lastKeyG
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	if updatedModel, ok := updatedModel.(Model); ok {
		model = updatedModel
		assert.True(t, model.lastKeyG, "First 'g' should set lastKeyG")
	}

	// Second 'g' should go to top
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	if updatedModel, ok := updatedModel.(Model); ok {
		model = updatedModel
		assert.False(t, model.lastKeyG, "Second 'g' should clear lastKeyG")
		assert.Equal(t, 0, model.selected, "Second 'g' should go to top")
	}

	// Other keys should clear lastKeyG
	model.lastKeyG = true
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	if updatedModel, ok := updatedModel.(Model); ok {
		model = updatedModel
		assert.False(t, model.lastKeyG, "Other keys should clear lastKeyG")
	}
}
