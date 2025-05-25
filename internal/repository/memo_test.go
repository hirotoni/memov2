package repository

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
)

func TestMemoRepoImpl_Save(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create a mock MemoFileInterface
	mockMemo, err := domain.NewMemoFile(
		time.Now(),
		"Test Memo",
		[]string{"TestCategory", "SubCategory"},
	)
	if err != nil {
		t.Fatalf("failed to create memo file: %v", err)
	}

	// Add some heading blocks to the mock
	mockMemo.SetHeadingBlocks([]*markdown.HeadingBlock{
		{
			Level:       2,
			HeadingText: "Section 1",
			ContentText: "This is section 1 content\n",
		},
		{
			Level:       2,
			HeadingText: "Section 2",
			ContentText: "This is section 2 content\n",
		},
	})

	// Import the repos package
	repo := NewMemo(tmpDir)

	// Test case 1: Save a new file
	t.Run("Save new file", func(t *testing.T) {
		err := repo.Save(mockMemo, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Check if file exists
		path := filepath.Join(tmpDir, mockMemo.Location(), mockMemo.FileName())
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file to be created at %s", path)
		}

		// Verify content
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		expectedContent := "---\ncategory: [\"TestCategory\", \"SubCategory\"]\n---\n\n# Test Memo\n\n## Section 1\n\nThis is section 1 content\n\n## Section 2\n\nThis is section 2 content\n\n"
		if string(content) != expectedContent {
			t.Errorf("file content mismatch.\nExpected: %q\nGot: %q", expectedContent, string(content))
		}
	})

	// Test case 2: Override existing file with truncate = true
	t.Run("Override existing file", func(t *testing.T) {
		// Modify the mock heading blocks
		mockMemo.SetHeadingBlocks([]*markdown.HeadingBlock{
			{
				Level:       2,
				HeadingText: "Updated Section",
				ContentText: "Updated content\n",
			},
		})

		err := repo.Save(mockMemo, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify updated content
		path := filepath.Join(tmpDir, mockMemo.Location(), mockMemo.FileName())
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		expectedContent := "---\ncategory: [\"TestCategory\", \"SubCategory\"]\n---\n\n# Test Memo\n\n## Updated Section\n\nUpdated content\n\n"
		if string(content) != expectedContent {
			t.Errorf("file content mismatch after update.\nExpected: %q\nGot: %q", expectedContent, string(content))
		}
	})

	// Test case 3: Don't override when truncate = false
	t.Run("Don't override when truncate is false", func(t *testing.T) {
		// Set different heading blocks
		mockMemo.SetHeadingBlocks([]*markdown.HeadingBlock{
			{
				Level:       2,
				HeadingText: "Should not be saved",
				ContentText: "This should not appear\n",
			},
		})

		err := repo.Save(mockMemo, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Content should remain the same as previous test
		path := filepath.Join(tmpDir, mockMemo.Location(), mockMemo.FileName())
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		expectedContent := "---\ncategory: [\"TestCategory\", \"SubCategory\"]\n---\n\n# Test Memo\n\n## Updated Section\n\nUpdated content\n\n"
		if string(content) != expectedContent {
			t.Errorf("file should not have been updated.\nExpected: %q\nGot: %q", expectedContent, string(content))
		}
	})
}
