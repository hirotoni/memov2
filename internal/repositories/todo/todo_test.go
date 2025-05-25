package todo

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
)

// createTestTodo creates a todo file for testing
func createTestTodo(t *testing.T, date time.Time, headingBlocks []*markdown.HeadingBlock) domain.TodoFileInterface {
	t.Helper()

	todo, err := domain.NewTodosFile(date)
	if err != nil {
		t.Fatalf("failed to create todo file: %v", err)
	}

	if headingBlocks != nil {
		todo.SetHeadingBlocks(headingBlocks)
	}

	return todo
}

func TestTodoRepoImpl_Save(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewTodo(tmpDir, logger)

	date, err := time.Parse(time.DateOnly, "2023-10-01")
	if err != nil {
		t.Fatalf("failed to parse date: %v", err)
	}

	testCases := []struct {
		name          string
		date          time.Time
		headingBlocks []*markdown.HeadingBlock
		truncate      bool
		goldenFile    string
	}{
		{
			name: "Save new todo file",
			date: date,
			headingBlocks: []*markdown.HeadingBlock{
				createTestHeadingBlock(2, "todos", "Buy groceries\nClean house\n"),
				createTestHeadingBlock(2, "wanttodos", "Learn Go\nRead book\n"),
			},
			truncate:   false,
			goldenFile: "todo_save_new_todo",
		},
		{
			name: "Override existing todo file with truncate = true",
			date: date,
			headingBlocks: []*markdown.HeadingBlock{
				createTestHeadingBlock(2, "todos", "Updated task list\n"),
			},
			truncate:   true,
			goldenFile: "todo_override_existing_todo",
		},
		{
			name: "Don't override when truncate is false",
			date: date,
			headingBlocks: []*markdown.HeadingBlock{
				createTestHeadingBlock(2, "should not appear", "This should not be saved\n"),
			},
			truncate:   false,
			goldenFile: "todo_dont_override_todo",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			todo := createTestTodo(t, tc.date, tc.headingBlocks)

			err := repo.Save(todo, tc.truncate)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check if file exists
			path := filepath.Join(tmpDir, todo.FileName())
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("expected file to be created at %s", path)
			}

			// Verify content against golden file
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			expectedContent := readGoldenFile(t, tc.goldenFile)
			if string(content) != expectedContent {
				t.Errorf("file content mismatch.\nExpected: %q\nGot: %q", expectedContent, string(content))
			}
		})
	}
}

func TestTodoRepoImpl_TodoEntries(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewTodo(tmpDir, logger)

	// Create multiple todo files
	dates := []string{"2023-10-01", "2023-10-02", "2023-10-03"}
	for _, dateStr := range dates {
		date, err := time.Parse(time.DateOnly, dateStr)
		if err != nil {
			t.Fatalf("failed to parse date %s: %v", dateStr, err)
		}

		todo := createTestTodo(t, date, []*markdown.HeadingBlock{
			createTestHeadingBlock(2, "todos", "Task for "+dateStr+"\n"),
		})

		err = repo.Save(todo, false)
		if err != nil {
			t.Fatalf("failed to save todo file: %v", err)
		}
	}

	t.Run("Get Todo Entries", func(t *testing.T) {
		entries, err := repo.TodoEntries()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(entries) != 3 {
			t.Errorf("expected 3 todo entries, got %d", len(entries))
		}

		// Verify entries are sorted by date (oldest first)
		for i := 0; i < len(entries)-1; i++ {
			if entries[i].Date().After(entries[i+1].Date()) {
				t.Errorf("entries not sorted by date: %s comes after %s",
					entries[i].Date().Format(time.DateOnly),
					entries[i+1].Date().Format(time.DateOnly))
			}
		}
	})
}

func TestTodoRepoImpl_TodosTemplate(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewTodo(tmpDir, logger)

	date, err := time.Parse(time.DateOnly, "2023-10-01")
	if err != nil {
		t.Fatalf("failed to parse date: %v", err)
	}

	t.Run("Create template file when it doesn't exist", func(t *testing.T) {
		template, err := repo.TodosTemplate(date)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Check if template file was created
		templatePath := filepath.Join(tmpDir, "todos_template.md")
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			t.Errorf("expected template file to be created at %s", templatePath)
		}

		// Verify template content
		content, err := os.ReadFile(templatePath)
		if err != nil {
			t.Fatalf("failed to read template file: %v", err)
		}

		expectedContent := readGoldenFile(t, "todo_todos_template")
		if string(content) != expectedContent {
			t.Errorf("template content mismatch.\nExpected: %q\nGot: %q", expectedContent, string(content))
		}

		// Verify returned template has correct heading blocks
		headingBlocks := template.HeadingBlocks()
		if len(headingBlocks) != 2 {
			t.Errorf("expected 2 heading blocks, got %d", len(headingBlocks))
		}

		expectedHeadings := []string{"todos", "wanttodos"}
		for i, expected := range expectedHeadings {
			if headingBlocks[i].HeadingText != expected {
				t.Errorf("expected heading %s, got %s", expected, headingBlocks[i].HeadingText)
			}
		}
	})

	t.Run("Reuse existing template file", func(t *testing.T) {
		// Call again - should reuse existing template
		template, err := repo.TodosTemplate(date)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should still have the same heading blocks
		headingBlocks := template.HeadingBlocks()
		if len(headingBlocks) != 2 {
			t.Errorf("expected 2 heading blocks, got %d", len(headingBlocks))
		}
	})
}

func TestTodoRepoImpl_FindTodosFileByDate(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewTodo(tmpDir, logger)

	date, err := time.Parse(time.DateOnly, "2023-10-01")
	if err != nil {
		t.Fatalf("failed to parse date: %v", err)
	}

	t.Run("Find existing todo file", func(t *testing.T) {
		// Create a todo file first
		todo := createTestTodo(t, date, []*markdown.HeadingBlock{
			createTestHeadingBlock(2, "todos", "Existing task\n"),
		})

		err := repo.Save(todo, false)
		if err != nil {
			t.Fatalf("failed to save todo file: %v", err)
		}

		// Try to find it
		found, err := repo.FindTodosFileByDate(date)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if found == nil {
			t.Fatal("expected to find todo file, got nil")
		}

		if !found.Date().Equal(date) {
			t.Errorf("expected date %s, got %s", date.Format(time.DateOnly), found.Date().Format(time.DateOnly))
		}
	})

	t.Run("Return error for non-existent file", func(t *testing.T) {
		nonExistentDate, err := time.Parse(time.DateOnly, "2023-12-31")
		if err != nil {
			t.Fatalf("failed to parse date: %v", err)
		}

		_, err = repo.FindTodosFileByDate(nonExistentDate)
		if err == nil {
			t.Error("expected error for non-existent file, got nil")
		}

		if !os.IsNotExist(err) {
			t.Errorf("expected os.ErrNotExist, got %v", err)
		}
	})
}
