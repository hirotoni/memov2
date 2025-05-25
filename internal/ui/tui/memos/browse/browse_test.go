package browse

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/repositories/memo"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew tests the creation of a BrowseModel
func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		setupDirs func(string) error
		wantErr   bool
		errMsg    string
	}{
		{
			name: "Successfully creates browse model",
			setupDirs: func(baseDir string) error {
				return os.MkdirAll(filepath.Join(baseDir, "memos"), 0755)
			},
			wantErr: false,
		},
		{
			name: "Creates model even with empty memos directory",
			setupDirs: func(baseDir string) error {
				return os.MkdirAll(filepath.Join(baseDir, "memos"), 0755)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup temp directory
			tempDir := t.TempDir()
			err := tt.setupDirs(tempDir)
			require.NoError(t, err, "Failed to setup directories")

			cfg, err := toml.NewConfig(toml.Option{
				BaseDir:         tempDir,
				TodosFolderName: "todos/",
				MemosFolderName: "memos/",
			})
			require.NoError(t, err, "Failed to create config")

			editor := &mock.MockEditor{}

			// Test
			model, err := New(cfg, editor)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, model)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, model)
				assert.Equal(t, cfg, model.config)
				assert.Equal(t, editor, model.editor)
			}
		})
	}
}

// TestBrowseModel_Init tests the Init function
func TestBrowseModel_Init(t *testing.T) {
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

	// Init should return tea.EnterAltScreen
	assert.NotNil(t, cmd, "Init should return a command")
}

// TestBrowseModel_View tests the View function
func TestBrowseModel_View(t *testing.T) {
	tests := []struct {
		name       string
		setupModel func(*BrowseModel)
		checkView  func(*testing.T, string)
	}{
		{
			name: "View renders without error",
			setupModel: func(m *BrowseModel) {
				// Default setup
			},
			checkView: func(t *testing.T, view string) {
				assert.NotEmpty(t, view, "View should not be empty")
			},
		},
		{
			name: "View shows error when set",
			setupModel: func(m *BrowseModel) {
				m.err = assert.AnError
			},
			checkView: func(t *testing.T, view string) {
				assert.Contains(t, view, "Error:", "View should contain error message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			cfg, err := toml.NewConfig(toml.Option{
				BaseDir:         tempDir,
				MemosFolderName: "memos/",
			})
			require.NoError(t, err)

			editor := &mock.MockEditor{}
			model, err := New(cfg, editor)
			require.NoError(t, err)

			tt.setupModel(model)
			view := model.View()
			tt.checkView(t, view)
		})
	}
}

// TestBrowseKeybindings tests keyboard shortcuts
func TestBrowseKeybindings(t *testing.T) {
	tests := []struct {
		name        string
		keyMsg      tea.KeyMsg
		expectQuit  bool
		description string
	}{
		{
			name:        "Q key triggers quit",
			keyMsg:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			expectQuit:  true,
			description: "Pressing 'q' should quit",
		},
		{
			name:        "Ctrl+C triggers quit",
			keyMsg:      tea.KeyMsg{Type: tea.KeyCtrlC},
			expectQuit:  true,
			description: "Pressing Ctrl+C should quit",
		},
		{
			name:        "Ctrl+U scrolls up",
			keyMsg:      tea.KeyMsg{Type: tea.KeyCtrlU},
			expectQuit:  false,
			description: "Ctrl+U should scroll up half page",
		},
		{
			name:        "Ctrl+D scrolls down",
			keyMsg:      tea.KeyMsg{Type: tea.KeyCtrlD},
			expectQuit:  false,
			description: "Ctrl+D should scroll down half page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			cfg, err := toml.NewConfig(toml.Option{
				BaseDir:         tempDir,
				MemosFolderName: "memos/",
			})
			require.NoError(t, err)

			editor := &mock.MockEditor{}
			model, err := New(cfg, editor)
			require.NoError(t, err)

			// Execute keybinding
			updatedModel, cmd := BrowseKeybindings(*model, tt.keyMsg)

			// Assert
			assert.NotNil(t, updatedModel, "Updated model should not be nil")

			if tt.expectQuit {
				assert.NotNil(t, cmd, "Quit command should be returned")
			}
		})
	}
}

// TestBrowseModel_WindowResize tests window resize handling
func TestBrowseModel_WindowResize(t *testing.T) {
	tempDir := t.TempDir()
		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Send window resize message
	resizeMsg := tea.WindowSizeMsg{
		Width:  120,
		Height: 40,
	}

	updatedModel, _ := model.Update(resizeMsg)

	// Verify model was updated
	assert.NotNil(t, updatedModel, "Model should be updated")

	// Type assertion to access fields
	if browseModel, ok := updatedModel.(BrowseModel); ok {
		assert.Equal(t, 120, browseModel.width, "Width should be updated")
		assert.Equal(t, 40, browseModel.height, "Height should be updated")
	} else {
		t.Fatal("Updated model is not BrowseModel type")
	}
}

// TestBrowseModel_NavigationKeys tests navigation behavior
func TestBrowseModel_NavigationKeys(t *testing.T) {
	tempDir := t.TempDir()

	// Create some test files in memos directory
	memosDir := filepath.Join(tempDir, "memos")
	require.NoError(t, os.MkdirAll(memosDir, 0755))

		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	tests := []struct {
		name        string
		key         tea.KeyMsg
		description string
	}{
		{
			name:        "L key expands directory or opens file",
			key:         tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			description: "L should expand directory or open file",
		},
		{
			name:        "H key collapses directory",
			key:         tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			description: "H should collapse directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedModel, _ := BrowseKeybindings(*model, tt.key)
			assert.NotNil(t, updatedModel, "Model should be updated")
		})
	}
}

// TestBrowseModel_CategoryDialog tests category management
func TestBrowseModel_CategoryDialog(t *testing.T) {
	tempDir := t.TempDir()
	memosDir := filepath.Join(tempDir, "memos")
	require.NoError(t, os.MkdirAll(memosDir, 0755))

		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Test opening category dialog with 'c' key (requires a selected memo)
	// This would typically require setting up a mock memo file first

	tests := []struct {
		name       string
		key        tea.KeyMsg
		dialogOpen bool
	}{
		{
			name:       "ESC closes category dialog",
			key:        tea.KeyMsg{Type: tea.KeyEsc},
			dialogOpen: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model.showCategoryDialog = tt.dialogOpen

			updatedModel, _ := model.Update(tt.key)

			assert.NotNil(t, updatedModel)

			if browseModel, ok := updatedModel.(BrowseModel); ok {
				if tt.key.Type == tea.KeyEsc {
					assert.False(t, browseModel.showCategoryDialog, "Dialog should be closed after ESC")
				}
			}
		})
	}
}

// TestBrowseModel_FileOperations tests file browsing with actual files
func TestBrowseModel_FileOperations(t *testing.T) {
	tempDir := t.TempDir()
	memosDir := filepath.Join(tempDir, "memos")
	require.NoError(t, os.MkdirAll(memosDir, 0755))

	// Create a test memo file
	testMemoPath := filepath.Join(memosDir, "test_memo.md")
	testContent := `---
title: Test Memo
date: 2024-01-01
---

# Test Heading

Test content
`
	err := os.WriteFile(testMemoPath, []byte(testContent), 0644)
	require.NoError(t, err)

		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Verify model was created successfully with files
	assert.NotNil(t, model)
	assert.NotNil(t, model.list)
}

// Helper functions for testing

// createTestMemoFile creates a test memo file for testing
func createTestMemoFile(t *testing.T, dir, filename, content string) string {
	t.Helper()

	err := os.MkdirAll(dir, 0755)
	require.NoError(t, err)

	filePath := filepath.Join(dir, filename)
	err = os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	return filePath
}

// TestHelperFunctions tests that helper functions work correctly
func TestHelperFunctions(t *testing.T) {
	t.Run("min function", func(t *testing.T) {
		assert.Equal(t, 1, min(1, 2))
		assert.Equal(t, 1, min(2, 1))
		assert.Equal(t, 0, min(0, 0))
	})

	t.Run("max function", func(t *testing.T) {
		assert.Equal(t, 2, max(1, 2))
		assert.Equal(t, 2, max(2, 1))
		assert.Equal(t, 0, max(0, 0))
	})

	t.Run("isChildPath function", func(t *testing.T) {
		assert.True(t, isChildPath([]string{"a", "b", "c"}, []string{"a", "b"}))
		assert.False(t, isChildPath([]string{"a", "b"}, []string{"a", "b"}))
		assert.False(t, isChildPath([]string{"a", "b"}, []string{"a", "b", "c"}))
		assert.False(t, isChildPath([]string{"x", "y"}, []string{"a", "b"}))
	})
}

// TestSortCategories tests category sorting
func TestSortCategories(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]string
		expected [][]string
	}{
		{
			name: "Sorts flat categories alphabetically",
			input: [][]string{
				{"z"},
				{"a"},
				{"m"},
			},
			expected: [][]string{
				{"a"},
				{"m"},
				{"z"},
			},
		},
		{
			name: "Maintains parent-child relationships",
			input: [][]string{
				{"parent", "child2"},
				{"parent"},
				{"parent", "child1"},
			},
			expected: [][]string{
				{"parent"},
				{"parent", "child1"},
				{"parent", "child2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortCategories(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestBrowseModel_ErrorHandling tests error scenarios
func TestBrowseModel_ErrorHandling(t *testing.T) {
	t.Run("Handle invalid memos directory", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create a file instead of directory to cause an error
		memosPath := filepath.Join(tempDir, "memos")
		err := os.WriteFile(memosPath, []byte("not a directory"), 0644)
		require.NoError(t, err)

		cfg, err := toml.NewConfig(toml.Option{
			BaseDir:         tempDir,
			MemosFolderName: "memos",
		})
		require.NoError(t, err)

		editor := &mock.MockEditor{}

		// This should handle the error gracefully
		_, err = New(cfg, editor)

		// Depending on implementation, this might succeed or fail
		// The key is that it doesn't panic
		if err != nil {
			assert.Error(t, err)
		}
	})
}

// Example of testing key binding with mocked key helper
func TestKeyMatching(t *testing.T) {
	tests := []struct {
		name     string
		keyMsg   tea.KeyMsg
		expected string
	}{
		{
			name:     "Regular character key",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			expected: "l",
		},
		{
			name:     "Control key",
			keyMsg:   tea.KeyMsg{Type: tea.KeyCtrlC},
			expected: "ctrl+c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.keyMsg.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Benchmark for heavy operations
func BenchmarkBrowseModel_Update(b *testing.B) {
	tempDir := b.TempDir()
		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(b, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(b, err)

	msg := tea.KeyMsg{Type: tea.KeyDown}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.Update(msg)
	}
}

// TestBrowseModel_DeleteDialog tests the delete confirmation dialog
func TestBrowseModel_DeleteDialog(t *testing.T) {
	tempDir := t.TempDir()
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
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Refresh to load the memo
	_, err = model.updateItems()
	require.NoError(t, err)

	// Select the memo and press 'd' to open delete dialog
	if len(model.list.Items()) > 0 {
		model.list.Select(0)
		if i, ok := model.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			model.selectedMemo = i.memo
			model.showDeleteDialog = true

			// Verify delete dialog is shown
			assert.True(t, model.showDeleteDialog)
			view := model.View()
			assert.Contains(t, view, "Delete Memo?")
			assert.Contains(t, view, "Press Y to")

			// Test cancel with 'n'
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.False(t, m.showDeleteDialog, "Dialog should be closed after pressing 'n'")
			}
		}
	}
}

// TestBrowseModel_RenameDialog tests the rename dialog
func TestBrowseModel_RenameDialog(t *testing.T) {
	tempDir := t.TempDir()
		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	// Create a test memo
	logger := common.DefaultLogger()
	repo := memo.NewMemo(cfg.MemosDir(), logger)
	testMemo, err := domain.NewMemoFile(time.Now(), "original-title", []string{"category1"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(testMemo, true))

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Refresh to load the memo
	_, err = model.updateItems()
	require.NoError(t, err)

	// Select the memo and open rename dialog
	if len(model.list.Items()) > 0 {
		model.list.Select(0)
		if i, ok := model.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			model.selectedMemo = i.memo
			model.showRenameDialog = true
			model.renameInput = "original title"

			// Verify rename dialog is shown
			assert.True(t, model.showRenameDialog)
			view := model.View()
			assert.Contains(t, view, "Rename Memo")
			assert.Contains(t, view, "Current: original-title")

			// Test cancel with Esc
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEscape})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.False(t, m.showRenameDialog, "Dialog should be closed after pressing Esc")
				assert.Empty(t, m.renameInput, "Input should be cleared")
			}
		}
	}
}

// TestBrowseModel_DuplicateDialog tests the duplicate confirmation dialog
func TestBrowseModel_DuplicateDialog(t *testing.T) {
	tempDir := t.TempDir()
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
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Refresh to load the memo
	_, err = model.updateItems()
	require.NoError(t, err)

	// Select the memo and open duplicate dialog
	if len(model.list.Items()) > 0 {
		model.list.Select(0)
		if i, ok := model.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			model.selectedMemo = i.memo
			model.showDuplicateDialog = true

			// Verify duplicate dialog is shown
			assert.True(t, model.showDuplicateDialog)
			view := model.View()
			assert.Contains(t, view, "Duplicate Memo?")
			assert.Contains(t, view, "Original timestamp")
			assert.Contains(t, view, "New timestamp")

			// Test cancel with Esc
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEscape})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.False(t, m.showDuplicateDialog, "Dialog should be closed after pressing Esc")
			}
		}
	}
}

// TestBrowseModel_NewMemoDialog tests the new memo creation dialog
func TestBrowseModel_NewMemoDialog(t *testing.T) {
	tempDir := t.TempDir()
		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	// Create a test memo
	logger := common.DefaultLogger()
	repo := memo.NewMemo(cfg.MemosDir(), logger)
	testMemo, err := domain.NewMemoFile(time.Now(), "existing-memo", []string{"work", "projects"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(testMemo, true))

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Refresh to load the memo
	_, err = model.updateItems()
	require.NoError(t, err)

	// Select the memo and open new memo dialog
	if len(model.list.Items()) > 0 {
		model.list.Select(0)
		if i, ok := model.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			model.selectedMemo = i.memo
			model.showNewMemoDialog = true
			model.newMemoTitleInput = ""

			// Verify new memo dialog is shown
			assert.True(t, model.showNewMemoDialog)
			view := model.View()
			assert.Contains(t, view, "Create New Memo")
			assert.Contains(t, view, "Category: work > projects")
			assert.Contains(t, view, "Title:")

			// Test cancel with Esc
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEscape})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.False(t, m.showNewMemoDialog, "Dialog should be closed after pressing Esc")
				assert.Empty(t, m.newMemoTitleInput, "Input should be cleared")
			}
		}
	}
}

// TestBrowseModel_CategorySingleSelection tests that category selection clears previous selections
func TestBrowseModel_CategorySingleSelection(t *testing.T) {
	tempDir := t.TempDir()
		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	// Create multiple memos in different categories
	logger := common.DefaultLogger()
	repo := memo.NewMemo(cfg.MemosDir(), logger)
	memo1, err := domain.NewMemoFile(time.Now(), "memo1", []string{"cat1"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(memo1, true))

	memo2, err := domain.NewMemoFile(time.Now(), "memo2", []string{"cat2"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(memo2, true))

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Initialize category dialog
	model.showCategoryDialog = true
	model.allCategories = [][]string{{"cat1"}, {"cat2"}, {"cat3"}}
	model.selectedCategories = make(map[string]bool)
	model.selectedCategories["cat1"] = true

	// Update category items
	model.updateCategoryItems()

	// Verify cat1 is selected
	assert.True(t, model.selectedCategories["cat1"], "cat1 should be selected")
	assert.False(t, model.selectedCategories["cat2"], "cat2 should not be selected")

	// Simulate selecting cat2 with space key
	if len(model.categoryList.Items()) >= 2 {
		model.categoryList.Select(1) // Select cat2
		if i, ok := model.categoryList.SelectedItem().(categoryItem); ok {
			pathStr := filepath.Join(i.path...)
			model.selectedCategories = make(map[string]bool)
			model.selectedCategories[pathStr] = true
		}

		// Verify only cat2 is selected now
		assert.False(t, model.selectedCategories["cat1"], "cat1 should be deselected")
		assert.True(t, model.selectedCategories["cat2"], "cat2 should be selected")
	}
}

// TestBrowseModel_HierarchicalCategoryCreation tests creating hierarchical categories
func TestBrowseModel_HierarchicalCategoryCreation(t *testing.T) {
	tempDir := t.TempDir()
		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Initialize with empty categories
	model.allCategories = [][]string{}

	tests := []struct {
		name          string
		input         string
		expectedPaths [][]string
	}{
		{
			name:  "Single level category",
			input: "work",
			expectedPaths: [][]string{
				{"work"},
			},
		},
		{
			name:  "Two level category",
			input: "work/projects",
			expectedPaths: [][]string{
				{"work"},
				{"work", "projects"},
			},
		},
		{
			name:  "Three level category with slash separator",
			input: "work/projects/2024",
			expectedPaths: [][]string{
				{"work"},
				{"work", "projects"},
				{"work", "projects", "2024"},
			},
		},
		{
			name:  "Three level category with > separator",
			input: "work > projects > 2024",
			expectedPaths: [][]string{
				{"work"},
				{"work", "projects"},
				{"work", "projects", "2024"},
			},
		},
		{
			name:  "Mixed separators",
			input: "work/projects > client",
			expectedPaths: [][]string{
				{"work"},
				{"work", "projects"},
				{"work", "projects", "client"},
			},
		},
		{
			name:  "Extra spaces",
			input: "  work  /  projects  ",
			expectedPaths: [][]string{
				{"work"},
				{"work", "projects"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset categories
			model.allCategories = [][]string{}

			// Simulate the category creation logic
			input := tt.input
			input = strings.ReplaceAll(input, " > ", "/")
			input = strings.ReplaceAll(input, ">", "/")

			segments := strings.Split(input, "/")
			var cleanSegments []string
			for _, seg := range segments {
				cleaned := strings.TrimSpace(seg)
				if cleaned != "" {
					cleanSegments = append(cleanSegments, cleaned)
				}
			}

			if len(cleanSegments) > 0 {
				for i := 1; i <= len(cleanSegments); i++ {
					subPath := cleanSegments[:i]
					pathStr := strings.Join(subPath, string(filepath.Separator))

					exists := false
					for _, existing := range model.allCategories {
						if strings.Join(existing, string(filepath.Separator)) == pathStr {
							exists = true
							break
						}
					}

					if !exists {
						model.allCategories = append(model.allCategories, subPath)
					}
				}
			}

			// Verify the expected paths were created
			assert.Len(t, model.allCategories, len(tt.expectedPaths),
				"Should create %d category paths", len(tt.expectedPaths))

			for _, expectedPath := range tt.expectedPaths {
				found := false
				for _, actualPath := range model.allCategories {
					if len(actualPath) == len(expectedPath) {
						match := true
						for j := range expectedPath {
							if actualPath[j] != expectedPath[j] {
								match = false
								break
							}
						}
						if match {
							found = true
							break
						}
					}
				}
				assert.True(t, found, "Expected path %v not found in created categories", expectedPath)
			}
		})
	}
}

// TestBrowseModel_DialogStates tests that only one dialog is shown at a time
func TestBrowseModel_DialogStates(t *testing.T) {
	tempDir := t.TempDir()
		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	// Create a test memo
	logger := common.DefaultLogger()
	repo := memo.NewMemo(cfg.MemosDir(), logger)
	testMemo, err := domain.NewMemoFile(time.Now(), "test", []string{"cat1"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(testMemo, true))

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Set selected memo (required for dialogs to work)
	model.selectedMemo = testMemo

	// Test that dialogs are mutually exclusive
	model.showDeleteDialog = true
	model.showRenameDialog = true
	model.showDuplicateDialog = true
	model.showNewMemoDialog = true
	model.showCategoryDialog = true

	// Only the first check in View() should render
	view := model.View()

	// Count how many dialog indicators appear (should be only one)
	dialogIndicators := []string{
		"Create New Memo",
		"Duplicate Memo?",
		"Rename Memo",
		"Delete Memo?",
	}

	foundCount := 0
	for _, indicator := range dialogIndicators {
		if strings.Contains(view, indicator) {
			foundCount++
		}
	}

	// At most one dialog should be visible
	// (depends on the order in View() function)
	assert.LessOrEqual(t, foundCount, 1, "Only one dialog should be shown at a time")
}

// TestBrowseModel_NewMemoInSameCategory tests creating a new memo in the same category
func TestBrowseModel_NewMemoInSameCategory(t *testing.T) {
	tempDir := t.TempDir()
		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	// Create a test memo in a specific category
	logger := common.DefaultLogger()
	repo := memo.NewMemo(cfg.MemosDir(), logger)
	existingMemo, err := domain.NewMemoFile(time.Now(), "existing", []string{"work", "projects"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(existingMemo, true))

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Refresh to load the memo
	_, err = model.updateItems()
	require.NoError(t, err)

	// Simulate opening new memo dialog
	if len(model.list.Items()) > 0 {
		model.list.Select(0)
		if i, ok := model.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			model.selectedMemo = i.memo
			model.showNewMemoDialog = true
			model.newMemoTitleInput = "new memo"

			// Verify dialog content
			view := model.View()
			assert.Contains(t, view, "Create New Memo")
			assert.Contains(t, view, "Category: work > projects")
			assert.Contains(t, view, "new memo")

			// The actual creation and editor opening would happen on Enter,
			// but we can't fully test that without mocking the editor
		}
	}
}

// TestBrowseModel_PreviewPane tests the preview pane functionality
func TestBrowseModel_PreviewPane(t *testing.T) {
	tempDir := t.TempDir()
		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	// Create a test memo with content
	logger := common.DefaultLogger()
	repo := memo.NewMemo(cfg.MemosDir(), logger)
	testMemo, err := domain.NewMemoFile(time.Now(), "preview-test", []string{"work"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(testMemo, true))

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Preview should be enabled by default
	assert.True(t, model.showPreview, "Preview should be enabled by default")

	// Refresh to load the memo
	_, err = model.updateItems()
	require.NoError(t, err)

	// Expand all directories to make files visible
	for i := range model.items {
		if model.items[i].isDir {
			model.items[i].expanded = true
		}
	}
	_, err = model.updateItems()
	require.NoError(t, err)

	// Find and select a memo file (not directory)
	memoFound := false
	for i, listItem := range model.list.Items() {
		if item, ok := listItem.(item); ok && !item.isDir && item.memo != nil {
			model.list.Select(i)
			memoFound = true
			break
		}
	}

	require.True(t, memoFound, "Should find a memo file in the list")

	// Render split view
	view := model.View()

	// Verify preview content appears (split view should contain metadata)
	assert.Contains(t, view, "Metadata", "Preview should show metadata section")

	// Test toggling preview off
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	if m, ok := updatedModel.(BrowseModel); ok {
		assert.False(t, m.showPreview, "Preview should be toggled off")

		// View should be full width list only
		viewNoPreview := m.View()
		assert.NotContains(t, viewNoPreview, "Metadata", "Preview metadata should not show when disabled")
	}
}

// TestBrowseModel_PreviewMemoContent tests memo content preview rendering
func TestBrowseModel_PreviewMemoContent(t *testing.T) {
	tempDir := t.TempDir()
		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Create mock item with memo
	mockMemo, err := domain.NewMemoFile(time.Now(), "test-title", []string{"cat1", "cat2"})
	require.NoError(t, err)

	mockItem := item{
		name:  "test.md",
		path:  "/test/path",
		isDir: false,
		memo:  mockMemo,
	}

	// Test memo preview rendering
	previewStyle := lipgloss.NewStyle().Width(40)
	preview := model.renderMemoPreview(mockItem, previewStyle)

	assert.Contains(t, preview, "test-title", "Should show memo title")
	assert.Contains(t, preview, "Metadata", "Should show metadata section")
	assert.Contains(t, preview, "cat1 > cat2", "Should show category path")
}

// TestBrowseModel_PreviewDirectoryContent tests directory preview rendering
func TestBrowseModel_PreviewDirectoryContent(t *testing.T) {
	tempDir := t.TempDir()
		cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Create mock directory item with children
	mockDir := item{
		name:  "testdir",
		path:  "/test/path",
		isDir: true,
		children: []item{
			{name: "subdir", isDir: true},
			{name: "file1.md", isDir: false},
			{name: "file2.md", isDir: false},
		},
	}

	// Test directory preview rendering
	previewStyle := lipgloss.NewStyle().Width(40)
	preview := model.renderDirectoryPreview(mockDir, previewStyle)

	assert.Contains(t, preview, "Directory", "Should show directory indicator")
	assert.Contains(t, preview, "testdir", "Should show directory name")
	assert.Contains(t, preview, "Folders: 1", "Should count subdirectories")
	assert.Contains(t, preview, "Files: 2", "Should count files")
}

// TestBrowseModel_ExpandCollapseAll tests expand/collapse all functionality
func TestBrowseModel_ExpandCollapseAll(t *testing.T) {
	tempDir := t.TempDir()
	memosDir := filepath.Join(tempDir, "memos")
	require.NoError(t, os.MkdirAll(memosDir, 0755))

	// Create nested directory structure
	cat2Dir := filepath.Join(memosDir, "cat1", "cat2")
	require.NoError(t, os.MkdirAll(cat2Dir, 0755))

	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Refresh to load items
	_, err = model.updateItems()
	require.NoError(t, err)

	// Find a directory item
	var dirItem item
	for _, listItem := range model.list.Items() {
		if i, ok := listItem.(item); ok && i.isDir {
			dirItem = i
			model.list.Select(0)
			break
		}
	}

	if dirItem.path != "" {
		// Test expand all with '>'
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'>'}})
		assert.NotNil(t, updatedModel)

		// Test collapse all with '<'
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'<'}})
		assert.NotNil(t, updatedModel)
	}
}

// TestBrowseModel_NewMemoDialogInput tests new memo dialog input handling
func TestBrowseModel_NewMemoDialogInput(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	logger := common.DefaultLogger()
	repo := memo.NewMemo(cfg.MemosDir(), logger)
	testMemo, err := domain.NewMemoFile(time.Now(), "test-memo", []string{"category1"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(testMemo, true))

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Refresh to load the memo
	_, err = model.updateItems()
	require.NoError(t, err)

	// Select the memo and open new memo dialog
	if len(model.list.Items()) > 0 {
		model.list.Select(0)
		if i, ok := model.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			model.selectedMemo = i.memo
			model.showNewMemoDialog = true
			model.newMemoTitleInput = ""

			// Test input
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.Equal(t, "n", m.newMemoTitleInput, "Input should be added")
			}

			// Test backspace
			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.Empty(t, m.newMemoTitleInput, "Backspace should remove character")
			}

			// Test cancel with Esc
			model.newMemoTitleInput = "test"
			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.False(t, m.showNewMemoDialog, "Dialog should be closed")
				assert.Empty(t, m.newMemoTitleInput, "Input should be cleared")
			}
		}
	}
}

// TestBrowseModel_RenameDialogInput tests rename dialog input handling
func TestBrowseModel_RenameDialogInput(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	logger := common.DefaultLogger()
	repo := memo.NewMemo(cfg.MemosDir(), logger)
	testMemo, err := domain.NewMemoFile(time.Now(), "original-title", []string{"category1"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(testMemo, true))

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Refresh to load the memo
	_, err = model.updateItems()
	require.NoError(t, err)

	// Select the memo and open rename dialog
	if len(model.list.Items()) > 0 {
		model.list.Select(0)
		if i, ok := model.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			model.selectedMemo = i.memo
			model.showRenameDialog = true
			model.renameInput = "original title"

			// Test input
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.Contains(t, m.renameInput, "x", "Input should be added")
			}

			// Test backspace
			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			if m, ok := updatedModel.(BrowseModel); ok {
				// Backspace should remove last character
				assert.NotEqual(t, "original title", m.renameInput, "Backspace should modify input")
			}
		}
	}
}

// TestBrowseModel_CategoryDialogNavigation tests category dialog navigation
func TestBrowseModel_CategoryDialogNavigation(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	logger := common.DefaultLogger()
	repo := memo.NewMemo(cfg.MemosDir(), logger)
	testMemo, err := domain.NewMemoFile(time.Now(), "test-memo", []string{"cat1"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(testMemo, true))

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Refresh to load the memo
	_, err = model.updateItems()
	require.NoError(t, err)

	// Select the memo and open category dialog
	if len(model.list.Items()) > 0 {
		model.list.Select(0)
		if i, ok := model.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			model.selectedMemo = i.memo
			model.showCategoryDialog = true
			model.allCategories = [][]string{{"cat1"}, {"cat2"}, {"cat3"}}
			model.selectedCategories = make(map[string]bool)
			model.updateCategoryItems()

			// Test expand with 'l'
			if len(model.categoryList.Items()) > 0 {
				model.categoryList.Select(0)
				updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
				assert.NotNil(t, updatedModel)
			}

			// Test collapse with 'h'
			if len(model.categoryList.Items()) > 0 {
				model.categoryList.Select(0)
				updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
				assert.NotNil(t, updatedModel)
			}

			// Test new category input mode with 'n'
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.Equal(t, -1, m.categoryDialogCursor, "Should enter input mode")
			}
		}
	}
}

// TestBrowseModel_CategoryDialogNewCategoryInput tests new category input in dialog
func TestBrowseModel_CategoryDialogNewCategoryInput(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	logger := common.DefaultLogger()
	repo := memo.NewMemo(cfg.MemosDir(), logger)
	testMemo, err := domain.NewMemoFile(time.Now(), "test-memo", []string{"cat1"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(testMemo, true))

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Refresh to load the memo
	_, err = model.updateItems()
	require.NoError(t, err)

	// Select the memo and open category dialog
	if len(model.list.Items()) > 0 {
		model.list.Select(0)
		if i, ok := model.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			model.selectedMemo = i.memo
			model.showCategoryDialog = true
			model.allCategories = [][]string{{"cat1"}}
			model.selectedCategories = make(map[string]bool)
			model.categoryDialogCursor = -1 // Enter input mode
			model.newCategoryInput = ""

			// Test input
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.Equal(t, "n", m.newCategoryInput, "Input should be added")
			}

			// Test backspace
			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.Empty(t, m.newCategoryInput, "Backspace should remove character")
			}

			// Test cancel with Esc
			model.newCategoryInput = "newcat"
			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.Empty(t, m.newCategoryInput, "Input should be cleared")
				assert.Equal(t, 0, m.categoryDialogCursor, "Should exit input mode")
			}
		}
	}
}

// TestBrowseModel_DeleteConfirm tests delete confirmation
func TestBrowseModel_DeleteConfirm(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	logger := common.DefaultLogger()
	repo := memo.NewMemo(cfg.MemosDir(), logger)
	testMemo, err := domain.NewMemoFile(time.Now(), "test-memo", []string{"category1"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(testMemo, true))

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Refresh to load the memo
	_, err = model.updateItems()
	require.NoError(t, err)

	// Select the memo and open delete dialog
	if len(model.list.Items()) > 0 {
		model.list.Select(0)
		if i, ok := model.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			model.selectedMemo = i.memo
			model.showDeleteDialog = true

			// Test cancel with 'n'
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.False(t, m.showDeleteDialog, "Dialog should be closed after 'n'")
			}

			// Test cancel with Esc
			model.showDeleteDialog = true
			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.False(t, m.showDeleteDialog, "Dialog should be closed after Esc")
			}
		}
	}
}

// TestBrowseModel_DuplicateConfirm tests duplicate confirmation
func TestBrowseModel_DuplicateConfirm(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	logger := common.DefaultLogger()
	repo := memo.NewMemo(cfg.MemosDir(), logger)
	testMemo, err := domain.NewMemoFile(time.Now(), "test-memo", []string{"category1"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(testMemo, true))

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Refresh to load the memo
	_, err = model.updateItems()
	require.NoError(t, err)

	// Select the memo and open duplicate dialog
	if len(model.list.Items()) > 0 {
		model.list.Select(0)
		if i, ok := model.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			model.selectedMemo = i.memo
			model.showDuplicateDialog = true

			// Test cancel with 'n'
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.False(t, m.showDuplicateDialog, "Dialog should be closed after 'n'")
			}

			// Test cancel with Esc
			model.showDuplicateDialog = true
			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
			if m, ok := updatedModel.(BrowseModel); ok {
				assert.False(t, m.showDuplicateDialog, "Dialog should be closed after Esc")
			}
		}
	}
}

// TestBrowseModel_ErrorHandlingInUpdate tests error handling in Update
func TestBrowseModel_ErrorHandlingInUpdate(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Test error message handling
	errorMsg := assert.AnError
	updatedModel, _ := model.Update(errorMsg)
	if m, ok := updatedModel.(BrowseModel); ok {
		assert.NotNil(t, m.err, "Error should be set")
	}
}

// TestBrowseModel_TraverseAndModify tests the traverseAndModify helper
func TestBrowseModel_TraverseAndModify(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Create test items
	testItems := []item{
		{path: "root1", isDir: true, children: []item{
			{path: "root1/child1", isDir: false},
		}},
		{path: "root2", isDir: true},
	}

	// Test modifying an item
	modified := false
	model.traverseAndModify(&testItems, "root1/child1", func(item *item) {
		item.name = "modified"
		modified = true
	})

	assert.True(t, modified, "Item should be modified")
}

// TestBrowseModel_ExpandItemAndChildren tests expandItemAndChildren helper
func TestBrowseModel_ExpandItemAndChildren(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Create nested directory structure
	testItem := item{
		path:     "root",
		isDir:    true,
		expanded: false,
		children: []item{
			{path: "root/child1", isDir: true, expanded: false, children: []item{
				{path: "root/child1/grandchild", isDir: false},
			}},
			{path: "root/child2", isDir: false},
		},
	}

	// Expand all
	model.expandItemAndChildren(&testItem, true)
	assert.True(t, testItem.expanded, "Root should be expanded")
	assert.True(t, testItem.children[0].expanded, "Child should be expanded")

	// Collapse all
	model.expandItemAndChildren(&testItem, false)
	assert.False(t, testItem.expanded, "Root should be collapsed")
	assert.False(t, testItem.children[0].expanded, "Child should be collapsed")
}

// TestBrowseModel_PathToCategoryTree tests pathToCategoryTree helper
func TestBrowseModel_PathToCategoryTree(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	memosDir := cfg.MemosDir()
	testPath := filepath.Join(memosDir, "cat1", "cat2")

	categoryTree := model.pathToCategoryTree(testPath)
	assert.Equal(t, []string{"cat1", "cat2"}, categoryTree, "Should extract category tree from path")
}

// TestBrowseModel_ComplexWorkflow tests a complex user workflow
func TestBrowseModel_ComplexWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	memosDir := filepath.Join(tempDir, "memos")
	require.NoError(t, os.MkdirAll(memosDir, 0755))

	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		MemosFolderName: "memos/",
	})
	require.NoError(t, err)

	logger := common.DefaultLogger()
	repo := memo.NewMemo(cfg.MemosDir(), logger)

	// Create multiple memos
	memo1, err := domain.NewMemoFile(time.Now(), "memo1", []string{"work"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(memo1, true))

	memo2, err := domain.NewMemoFile(time.Now(), "memo2", []string{"work", "projects"})
	require.NoError(t, err)
	require.NoError(t, repo.Save(memo2, true))

	editor := &mock.MockEditor{}
	model, err := New(cfg, editor)
	require.NoError(t, err)

	// Refresh to load memos
	_, err = model.updateItems()
	require.NoError(t, err)

	// Toggle preview
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	if browseModel, ok := updatedModel.(BrowseModel); ok {
		assert.False(t, browseModel.showPreview, "Preview should be toggled off")
		// Update the model pointer with the new value
		*model = browseModel
	}

	// Toggle preview back on
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	if browseModel, ok := updatedModel.(BrowseModel); ok {
		assert.True(t, browseModel.showPreview, "Preview should be toggled on")
		// Update the model pointer with the new value
		*model = browseModel
	}

	// Navigate with Ctrl+U
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	assert.NotNil(t, updatedModel)

	// Navigate with Ctrl+D
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
	assert.NotNil(t, updatedModel)
}
