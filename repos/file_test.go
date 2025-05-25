package repos

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hirotoni/memov2/models"
	"github.com/hirotoni/memov2/utils"
)

func TestFileRepoImpl_Save(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup a mock file
	date := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	filetype := models.FileTypeTodos
	title := "savetest"
	mock, err := models.NewFile(date, filetype, title)
	if err != nil {
		t.Fatalf("failed to create mock file: %v", err)
	}
	repo := NewFileRepo(tmpDir)
	t.Run("saves file correctly", func(t *testing.T) {
		err := repo.Save(mock, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Check if file exists
		expectedPath := filepath.Join(tmpDir, date.Format(models.FileNameDateLayoutTodo)+models.Sep+filetype.String()+models.Ext)
		if !utils.Exists(expectedPath) {
			t.Errorf("expected file to be saved at %s, but it does not exist", expectedPath)
		}

	})
}
func TestNewFileRepo(t *testing.T) {
	t.Run("creates a new file repo with the specified base directory", func(t *testing.T) {
		baseDir := "/test/path"
		repo := NewFileRepo(baseDir)

		// Since fileRepo is an interface, we need to check the underlying type
		impl, ok := repo.(*fileRepoImpl)
		if !ok {
			t.Fatalf("expected repo to be of type *fileRepoImpl, got %T", repo)
		}

		if impl.dir != baseDir {
			t.Errorf("expected baseDir to be %q, got %q", baseDir, impl.dir)
		}
	})

	t.Run("returns a valid fileRepo implementation", func(t *testing.T) {
		tmpDir := t.TempDir()
		repo := NewFileRepo(tmpDir)

		// Make sure the repo implements the fileRepo interface correctly
		var _ fileRepo = repo

		// Basic functionality test
		date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		mockFile, _ := models.NewTodosFile(date)

		// This will fail if the repo doesn't properly implement the interface methods
		err := repo.Save(mockFile, true)
		if err != nil {
			t.Fatalf("basic functionality test failed: %v", err)
		}
	})
}
func TestFileRepoImpl_TodosTemplate(t *testing.T) {
	t.Run("successfully creates template file", func(t *testing.T) {
		// Setup a test directory and template file
		tmpDir := t.TempDir()
		templatePath := filepath.Join(tmpDir, "todos_template.md")

		// Create a sample template file
		templateContent := "# Template\n\n## Section 1\n\nContent 1\n\n## Section 2\n\nContent 2\n"
		err := os.WriteFile(templatePath, []byte(templateContent), 0644)
		if err != nil {
			t.Fatalf("failed to create test template file: %v", err)
		}

		repo := NewFileRepo(tmpDir)
		date := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

		file, err := repo.TodosTemplate(date)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify the file was created correctly
		if file == nil {
			t.Fatal("expected file to be non-nil")
		}

		// Check if entities were parsed correctly
		entities := file.HeadingBlocks()
		if len(entities) != 2 {
			t.Errorf("expected 2 entities, got %d", len(entities))
		}

		// Check if date was set correctly
		expectedDate := date.Format("2006-01-02")
		if file.Date().Format("2006-01-02") != expectedDate {
			t.Errorf("expected date %s, got %s", expectedDate, file.Date().Format("2006-01-02"))
		}
	})
}

func TestFileRepoImpl_FindTodosFileByDate(t *testing.T) {
	t.Run("returns file when it exists", func(t *testing.T) {
		// Setup a test directory and file
		tmpDir := t.TempDir()
		date := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

		// Create a mock todos file first
		mockFile, err := models.NewTodosFile(date)
		if err != nil {
			t.Fatalf("failed to create mock file: %v", err)
		}

		repo := NewFileRepo(tmpDir)

		// Save the file first
		if err := repo.Save(mockFile, true); err != nil {
			t.Fatalf("failed to save mock file: %v", err)
		}

		// Now try to find it
		file, err := repo.FindTodosFileByDate(date)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if file == nil {
			t.Fatal("expected file to be non-nil")
		}

		// Check if date was set correctly
		expectedDate := date.Format("2006-01-02")
		if file.Date().Format("2006-01-02") != expectedDate {
			t.Errorf("expected date %s, got %s", expectedDate, file.Date().Format("2006-01-02"))
		}
	})

	t.Run("returns error when file doesn't exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		date := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

		repo := NewFileRepo(tmpDir)

		_, err := repo.FindTodosFileByDate(date)
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected os.ErrNotExist, got %v", err)
		}
	})
}
