package repository

import (
	"os"
	"path/filepath"
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
	repo := NewMemo(tmpDir)

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
	repo := NewMemo(tmpDir)

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
