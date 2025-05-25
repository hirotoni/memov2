package memo

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
)

// createTestMemo creates a memo file for testing
func createTestMemo(t *testing.T, title string, categories []string, headingBlocks []*markdown.HeadingBlock) domain.MemoFileInterface {
	t.Helper()

	memo, err := domain.NewMemoFile(
		time.Now(),
		title,
		categories,
	)
	if err != nil {
		t.Fatalf("failed to create memo file: %v", err)
	}

	if headingBlocks != nil {
		memo.SetHeadingBlocks(headingBlocks)
	}

	return memo
}

// createTestHeadingBlock creates a heading block for testing
func createTestHeadingBlock(level int, heading, content string) *markdown.HeadingBlock {
	return &markdown.HeadingBlock{
		Level:       level,
		HeadingText: heading,
		ContentText: content,
	}
}

// readGoldenFile reads the expected content from a golden file
func readGoldenFile(t *testing.T, testName string) string {
	t.Helper()

	goldenPath := filepath.Join("testdata", testName+".golden")
	content, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file %s: %v", goldenPath, err)
	}
	return string(content)
}

// writeGoldenFile writes content to a golden file (useful for updating expected results)
func writeGoldenFile(t *testing.T, testName, content string) {
	t.Helper()

	// Create testdata directory if it doesn't exist
	testDataDir := "testdata"
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		t.Fatalf("failed to create testdata directory: %v", err)
	}

	goldenPath := filepath.Join(testDataDir, testName+".golden")
	if err := os.WriteFile(goldenPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write golden file %s: %v", goldenPath, err)
	}
}

func TestMemoRepoImpl_Save(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tmpDir, logger)

	testCases := []struct {
		name          string
		title         string
		categories    []string
		headingBlocks []*markdown.HeadingBlock
		truncate      bool
		goldenFile    string
	}{
		{
			name:       "Save new file",
			title:      "Test Memo",
			categories: []string{"TestCategory", "SubCategory"},
			headingBlocks: []*markdown.HeadingBlock{
				createTestHeadingBlock(2, "Section 1", "This is section 1 content\n"),
				createTestHeadingBlock(2, "Section 2", "This is section 2 content\n"),
			},
			truncate:   false,
			goldenFile: "memo_save_new_file",
		},
		{
			name:       "Override existing file with truncate = true",
			title:      "Test Memo",
			categories: []string{"TestCategory", "SubCategory"},
			headingBlocks: []*markdown.HeadingBlock{
				createTestHeadingBlock(2, "Updated Section", "Updated content\n"),
			},
			truncate:   true,
			goldenFile: "memo_override_existing_file",
		},
		{
			name:       "Don't override when truncate is false",
			title:      "Test Memo",
			categories: []string{"TestCategory", "SubCategory"},
			headingBlocks: []*markdown.HeadingBlock{
				createTestHeadingBlock(2, "Should not be saved", "This should not appear\n"),
			},
			truncate:   false,
			goldenFile: "memo_dont_override_file",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			memo := createTestMemo(t, tc.title, tc.categories, tc.headingBlocks)

			err := repo.Save(memo, tc.truncate)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check if file exists
			path := filepath.Join(tmpDir, memo.Location(), memo.FileName())
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

// TestMemoRepoImpl_Save_EdgeCases tests edge cases for the Save method
func TestMemoRepoImpl_Save_EdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tmpDir, logger)

	testCases := []struct {
		name          string
		title         string
		categories    []string
		headingBlocks []*markdown.HeadingBlock
		goldenFile    string
	}{
		{
			name:       "Save memo with empty categories",
			title:      "Empty Categories Memo",
			categories: []string{},
			headingBlocks: []*markdown.HeadingBlock{
				createTestHeadingBlock(1, "Main Heading", "Main content\n"),
			},
			goldenFile: "memo_empty_categories",
		},
		{
			name:          "Save memo with no heading blocks",
			title:         "No Headings Memo",
			categories:    []string{"Category"},
			headingBlocks: nil,
			goldenFile:    "memo_no_heading_blocks",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			memo := createTestMemo(t, tc.title, tc.categories, tc.headingBlocks)

			err := repo.Save(memo, false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			filePath := filepath.Join(tmpDir, memo.Location(), memo.FileName())
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("expected file to exist at %s", filePath)
			}

			content, err := os.ReadFile(filePath)
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

func TestMemoEntries_Success(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tmpDir, logger)

	// Create multiple memos
	date1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)

	memo1, _ := domain.NewMemoFile(date1, "Memo 1", []string{})
	memo2, _ := domain.NewMemoFile(date2, "Memo 2", []string{"cat1"})
	memo3, _ := domain.NewMemoFile(date3, "Memo 3", []string{"cat1", "cat2"})

	if err := repo.Save(memo1, false); err != nil {
		t.Fatalf("failed to save memo1: %v", err)
	}
	if err := repo.Save(memo2, false); err != nil {
		t.Fatalf("failed to save memo2: %v", err)
	}
	if err := repo.Save(memo3, false); err != nil {
		t.Fatalf("failed to save memo3: %v", err)
	}

	// Execute
	entries, err := repo.MemoEntries()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
	// Verify they're sorted by date
	if !entries[0].Date().Before(entries[1].Date()) {
		t.Error("entries not sorted by date")
	}
	if !entries[1].Date().Before(entries[2].Date()) {
		t.Error("entries not sorted by date")
	}
}

func TestMemoEntries_EmptyDirectory(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tmpDir, logger)

	// Execute
	entries, err := repo.MemoEntries()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestMetadata(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tmpDir, logger)

	// Create and save a memo
	memo := createTestMemo(t, "Test Memo", []string{"cat1", "cat2"}, []*markdown.HeadingBlock{
		createTestHeadingBlock(2, "Heading 1", "Content 1"),
	})
	if err := repo.Save(memo, false); err != nil {
		t.Fatalf("failed to save memo: %v", err)
	}

	// Execute
	metadata, err := repo.Metadata(memo)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if metadata == nil {
		t.Fatal("metadata is nil")
	}
	// The actual Metadata method returns { "category": []string }
	if _, ok := metadata["category"]; !ok {
		t.Error("metadata does not contain 'category'")
	}
	categories, ok := metadata["category"].([]interface{})
	if !ok {
		t.Error("category should be a slice")
	}
	if len(categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(categories))
	}
}

// TestTidyMemos has been moved to Service layer
// These tests are no longer valid as TidyMemos is now a Service layer responsibility
// The tests should be moved to internal/service/memo/tidy_test.go

func TestCategories_Success(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tmpDir, logger)

	// Create memos with different categories
	memo1, _ := domain.NewMemoFile(time.Now(), "Memo 1", []string{"category1"})
	memo2, _ := domain.NewMemoFile(time.Now(), "Memo 2", []string{"category1", "subcategory1"})
	memo3, _ := domain.NewMemoFile(time.Now(), "Memo 3", []string{"category2"})

	if err := repo.Save(memo1, false); err != nil {
		t.Fatalf("failed to save memo1: %v", err)
	}
	if err := repo.Save(memo2, false); err != nil {
		t.Fatalf("failed to save memo2: %v", err)
	}
	if err := repo.Save(memo3, false); err != nil {
		t.Fatalf("failed to save memo3: %v", err)
	}

	// Execute
	categories, err := repo.Categories()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(categories) == 0 {
		t.Error("expected categories, got none")
	}

	// Check that we have at least the categories we created
	categoryStrings := make([]string, 0)
	for _, cat := range categories {
		categoryStrings = append(categoryStrings, cat[len(cat)-1])
	}
	hasCategory1 := false
	hasCategory2 := false
	for _, c := range categoryStrings {
		if c == "category1" {
			hasCategory1 = true
		}
		if c == "category2" {
			hasCategory2 = true
		}
	}
	if !hasCategory1 {
		t.Error("expected to find category1")
	}
	if !hasCategory2 {
		t.Error("expected to find category2")
	}
}

func TestCategories_EmptyDirectory(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tmpDir, logger)

	// Execute
	categories, err := repo.Categories()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(categories) != 0 {
		t.Errorf("expected empty categories, got %d", len(categories))
	}
}

func TestMove_Success(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tmpDir, logger)

	// Create a memo
	memo := createTestMemo(t, "Move Test", []string{"old", "category"}, nil)
	if err := repo.Save(memo, false); err != nil {
		t.Fatalf("failed to save memo: %v", err)
	}

	oldPath := filepath.Join(tmpDir, memo.Location(), memo.FileName())
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		t.Fatalf("file should exist at %s before move", oldPath)
	}

	// Execute - move to new category
	newCategoryTree := []string{"new", "category", "tree"}
	err := repo.Move(memo, newCategoryTree)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check old path doesn't exist
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Errorf("file should not exist at old path %s", oldPath)
	}

	// Check new path exists
	newPath := filepath.Join(tmpDir, filepath.Join(newCategoryTree...), memo.FileName())
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Errorf("file should exist at new path %s", newPath)
	}
}

func TestMove_SameLocation(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tmpDir, logger)

	// Create a memo
	categories := []string{"category1", "subcategory1"}
	memo := createTestMemo(t, "Same Location", categories, nil)
	if err := repo.Save(memo, false); err != nil {
		t.Fatalf("failed to save memo: %v", err)
	}

	originalPath := filepath.Join(tmpDir, memo.Location(), memo.FileName())
	if _, err := os.Stat(originalPath); os.IsNotExist(err) {
		t.Fatalf("file should exist before move: %s", originalPath)
	}

	// Execute - move to same location
	err := repo.Move(memo, categories)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(originalPath); os.IsNotExist(err) {
		t.Errorf("file should still exist in original location: %s", originalPath)
	}
}

func TestMove_ToRoot(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tmpDir, logger)

	// Create a memo with categories
	memo := createTestMemo(t, "To Root", []string{"deep", "category"}, nil)
	if err := repo.Save(memo, false); err != nil {
		t.Fatalf("failed to save memo: %v", err)
	}

	oldPath := filepath.Join(tmpDir, memo.Location(), memo.FileName())
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		t.Fatalf("file should exist before move: %s", oldPath)
	}

	// Execute - move to root (empty category tree)
	err := repo.Move(memo, []string{})

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check old path doesn't exist
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Errorf("file should not exist at old path %s", oldPath)
	}

	// Check new path exists at root
	newPath := filepath.Join(tmpDir, memo.FileName())
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Errorf("file should exist at new path %s", newPath)
	}
}

func TestMove_FileNotExist(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tmpDir, logger)

	// Create a memo but don't save it
	memo, _ := domain.NewMemoFile(time.Now(), "Not Saved", []string{})

	// Execute
	err := repo.Move(memo, []string{"new", "location"})

	// Assert
	if err == nil {
		t.Error("expected error when moving non-existent file")
	}
}

func TestScanDirectory_WithSubdirectories(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tmpDir, logger).(*memo)

	// Create memos in different subdirectories
	memo1, _ := domain.NewMemoFile(time.Now(), "Root Memo", []string{})
	memo2, _ := domain.NewMemoFile(time.Now(), "Cat1 Memo", []string{"cat1"})
	memo3, _ := domain.NewMemoFile(time.Now(), "Deep Memo", []string{"cat1", "cat2", "cat3"})

	if err := repo.Save(memo1, false); err != nil {
		t.Fatalf("failed to save memo1: %v", err)
	}
	if err := repo.Save(memo2, false); err != nil {
		t.Fatalf("failed to save memo2: %v", err)
	}
	if err := repo.Save(memo3, false); err != nil {
		t.Fatalf("failed to save memo3: %v", err)
	}

	// Execute
	files, err := repo.scanDirectory()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 3 {
		t.Errorf("expected 3 files, got %d", len(files))
	}
}

func TestDelete_Success(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	testTrashDir := filepath.Join(tempDir, ".trash")

	// Set test trash directory environment variable
	os.Setenv("TEST_TRASH_DIR", testTrashDir)
	defer os.Unsetenv("TEST_TRASH_DIR")

	// Create a memo repository
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tempDir, logger)

	// Create and save a test memo
	memo := createTestMemo(t, "test", []string{"category1"}, nil)
	if err := repo.Save(memo, true); err != nil {
		t.Fatalf("failed to save memo: %v", err)
	}

	// Verify the file exists
	memoPath := filepath.Join(tempDir, memo.Location(), memo.FileName())
	if _, err := os.Stat(memoPath); os.IsNotExist(err) {
		t.Fatal("Memo file should exist before deletion")
	}

	// Delete the memo
	if err := repo.Delete(memo); err != nil {
		t.Fatalf("failed to delete memo: %v", err)
	}

	// Verify the file no longer exists
	if _, err := os.Stat(memoPath); !os.IsNotExist(err) {
		t.Error("Memo file should not exist after deletion")
	}
}

func TestDelete_FileNotExist(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a memo repository
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tempDir, logger)

	// Create a test memo but don't save it
	memo := createTestMemo(t, "nonexistent", []string{"category1"}, nil)

	// Try to delete the non-existent memo
	err := repo.Delete(memo)
	if err == nil {
		t.Error("Expected error when deleting non-existent file")
	}
}

func TestDelete_CleansUpEmptyDirectories(t *testing.T) {
	// Note: Empty directory cleanup has been moved to Service layer (TidyMemos)
	// This test verifies that Delete only removes the file, not the directories
	// Directory cleanup is now handled by Service layer's TidyMemos method

	// Create a temporary directory
	tempDir := t.TempDir()
	testTrashDir := filepath.Join(tempDir, ".trash")

	// Set test trash directory environment variable
	os.Setenv("TEST_TRASH_DIR", testTrashDir)
	defer os.Unsetenv("TEST_TRASH_DIR")

	// Create a memo repository
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tempDir, logger)

	// Create and save a test memo in a nested category
	memo := createTestMemo(t, "test", []string{"parent", "child"}, nil)
	if err := repo.Save(memo, true); err != nil {
		t.Fatalf("failed to save memo: %v", err)
	}

	// Verify the directory structure exists
	categoryDir := filepath.Join(tempDir, "parent", "child")
	if _, err := os.Stat(categoryDir); os.IsNotExist(err) {
		t.Fatal("Category directory should exist before deletion")
	}

	// Delete the memo
	if err := repo.Delete(memo); err != nil {
		t.Fatalf("failed to delete memo: %v", err)
	}

	// Verify the file is deleted
	memoPath := filepath.Join(tempDir, memo.Location(), memo.FileName())
	if _, err := os.Stat(memoPath); !os.IsNotExist(err) {
		t.Error("Memo file should not exist after deletion")
	}

	// Note: Empty directories are NOT cleaned up by Delete method
	// This is now the responsibility of Service layer's TidyMemos method
	// The directory may still exist after deletion
	if _, err := os.Stat(categoryDir); os.IsNotExist(err) {
		t.Log("Directory was removed (this is fine, but cleanup is now handled by Service layer)")
	}
}

func TestRename_Success(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a memo repository
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tempDir, logger)

	// Create and save a test memo
	memo := createTestMemo(t, "old-title", []string{"category1"}, nil)
	if err := repo.Save(memo, true); err != nil {
		t.Fatalf("failed to save memo: %v", err)
	}

	// Verify the original file exists
	oldPath := filepath.Join(tempDir, memo.Location(), memo.FileName())
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		t.Fatal("Original memo file should exist before rename")
	}

	// Rename the memo
	newTitle := "new-title"
	if err := repo.Rename(memo, newTitle); err != nil {
		t.Fatalf("failed to rename memo: %v", err)
	}

	// Verify the old file no longer exists
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Error("Old memo file should not exist after rename")
	}

	// Load the renamed memo and verify title
	renamedMemo, err := repo.Memo(memo)
	if err != nil {
		// Try loading with new title
		newMemo := createTestMemo(t, newTitle, []string{"category1"}, nil)
		renamedMemo, err = repo.Memo(newMemo)
		if err != nil {
			t.Fatalf("failed to load renamed memo: %v", err)
		}
	}

	if renamedMemo.Title() != newTitle {
		t.Errorf("expected title %s, got %s", newTitle, renamedMemo.Title())
	}

	// Verify the new file exists
	newPath := filepath.Join(tempDir, renamedMemo.Location(), renamedMemo.FileName())
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Error("New memo file should exist after rename")
	}
}

func TestRename_FileNotExist(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a memo repository
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tempDir, logger)

	// Create a test memo but don't save it
	memo := createTestMemo(t, "nonexistent", []string{"category1"}, nil)

	// Try to rename the non-existent memo
	err := repo.Rename(memo, "new-title")
	if err == nil {
		t.Error("Expected error when renaming non-existent file")
	}
}

func TestDuplicate_Success(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a memo repository
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tempDir, logger)

	// Create and save a test memo with content
	headingBlocks := []*markdown.HeadingBlock{
		createTestHeadingBlock(2, "Section 1", "Content 1"),
		createTestHeadingBlock(2, "Section 2", "Content 2"),
	}
	memo := createTestMemo(t, "original", []string{"category1"}, headingBlocks)
	if err := repo.Save(memo, true); err != nil {
		t.Fatalf("failed to save memo: %v", err)
	}

	// Duplicate the memo
	duplicated, err := repo.Duplicate(memo)
	if err != nil {
		t.Fatalf("failed to duplicate memo: %v", err)
	}

	// Verify the duplicated memo exists
	dupPath := filepath.Join(tempDir, duplicated.Location(), duplicated.FileName())
	if _, err := os.Stat(dupPath); os.IsNotExist(err) {
		t.Error("Duplicated memo file should exist")
	}

	// Verify the title has " copied" appended
	expectedTitle := "original copied"
	if duplicated.Title() != expectedTitle {
		t.Errorf("expected title %s, got %s", expectedTitle, duplicated.Title())
	}

	// Verify the category is the same
	if len(duplicated.CategoryTree()) != len(memo.CategoryTree()) {
		t.Errorf("duplicated memo should have same category tree")
	}
	for i := range memo.CategoryTree() {
		if duplicated.CategoryTree()[i] != memo.CategoryTree()[i] {
			t.Errorf("category mismatch at index %d: expected %s, got %s",
				i, memo.CategoryTree()[i], duplicated.CategoryTree()[i])
		}
	}

	// Verify content is copied
	if len(duplicated.HeadingBlocks()) != len(memo.HeadingBlocks()) {
		t.Errorf("expected %d heading blocks, got %d",
			len(memo.HeadingBlocks()), len(duplicated.HeadingBlocks()))
	}

	// Verify timestamp is different (newer)
	if !duplicated.Date().After(memo.Date()) && !duplicated.Date().Equal(memo.Date()) {
		t.Error("duplicated memo should have same or newer timestamp")
	}

	// Verify original still exists
	origPath := filepath.Join(tempDir, memo.Location(), memo.FileName())
	if _, err := os.Stat(origPath); os.IsNotExist(err) {
		t.Error("Original memo should still exist after duplication")
	}
}

func TestDuplicate_FileNotExist(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a memo repository
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tempDir, logger)

	// Create a test memo but don't save it
	memo := createTestMemo(t, "nonexistent", []string{"category1"}, nil)

	// Try to duplicate the non-existent memo
	_, err := repo.Duplicate(memo)
	if err == nil {
		t.Error("Expected error when duplicating non-existent file")
	}
}

func TestDuplicate_PreservesContent(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a memo repository
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewMemo(tempDir, logger)

	// Create a memo with content (only level 2 headings as per repo logic)
	headingBlocks := []*markdown.HeadingBlock{
		createTestHeadingBlock(2, "Introduction", "This is the intro"),
		createTestHeadingBlock(2, "Details", "These are the details"),
	}
	memo := createTestMemo(t, "complex-memo", []string{"work", "projects"}, headingBlocks)

	// Set top-level content
	topLevel := createTestHeadingBlock(0, "", "Top level content here")
	memo.SetTopLevelBodyContent(topLevel)

	if err := repo.Save(memo, true); err != nil {
		t.Fatalf("failed to save memo: %v", err)
	}

	// Duplicate the memo
	duplicated, err := repo.Duplicate(memo)
	if err != nil {
		t.Fatalf("failed to duplicate memo: %v", err)
	}

	// Reload the duplicated memo to verify saved content
	reloaded, err := repo.Memo(duplicated)
	if err != nil {
		t.Fatalf("failed to reload duplicated memo: %v", err)
	}

	// Verify heading blocks are preserved
	if len(reloaded.HeadingBlocks()) != len(headingBlocks) {
		t.Errorf("expected %d heading blocks, got %d",
			len(headingBlocks), len(reloaded.HeadingBlocks()))
	}

	// Verify top-level content is preserved (with possible formatting differences)
	if !strings.Contains(reloaded.TopLevelBodyContent().ContentText, topLevel.ContentText) {
		t.Errorf("top-level content not preserved: expected to contain %q, got %q",
			topLevel.ContentText, reloaded.TopLevelBodyContent().ContentText)
	}
}
