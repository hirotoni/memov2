package weekly

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
)

// createTestWeekly creates a weekly file for testing
func createTestWeekly(t *testing.T, headingBlocks []*markdown.HeadingBlock) domain.WeeklyFileInterface {
	t.Helper()

	weekly, err := domain.NewWeekly()
	if err != nil {
		t.Fatalf("failed to create weekly file: %v", err)
	}

	if headingBlocks != nil {
		weekly.SetHeadingBlocks(headingBlocks)
	}

	return weekly
}

func TestWeeklyRepoImpl_Save(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewWeekly(tmpDir, logger)

	testCases := []struct {
		name          string
		headingBlocks []*markdown.HeadingBlock
		truncate      bool
		goldenFile    string
	}{
		{
			name: "Save new weekly file",
			headingBlocks: []*markdown.HeadingBlock{
				createTestHeadingBlock(2, "Weekly Summary", "This week was productive\n"),
				createTestHeadingBlock(2, "Goals for Next Week", "Complete project milestone\nLearn new technology\n"),
			},
			truncate:   false,
			goldenFile: "weekly_save_new",
		},
		{
			name: "Override existing weekly file with truncate = true",
			headingBlocks: []*markdown.HeadingBlock{
				createTestHeadingBlock(2, "Updated Summary", "Updated weekly content\n"),
			},
			truncate:   true,
			goldenFile: "weekly_override_existing",
		},
		{
			name: "Don't override when truncate is false",
			headingBlocks: []*markdown.HeadingBlock{
				createTestHeadingBlock(2, "Should not appear", "This should not be saved\n"),
			},
			truncate:   false,
			goldenFile: "weekly_dont_override",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			weekly := createTestWeekly(t, tc.headingBlocks)

			err := repo.Save(weekly, tc.truncate)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check if file exists
			path := filepath.Join(tmpDir, weekly.FileName())
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

// TestWeeklyRepoImpl_Save_EdgeCases tests edge cases for the Save method
func TestWeeklyRepoImpl_Save_EdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewWeekly(tmpDir, logger)

	testCases := []struct {
		name          string
		headingBlocks []*markdown.HeadingBlock
		goldenFile    string
	}{
		{
			name:          "Save weekly with no heading blocks",
			headingBlocks: nil,
			goldenFile:    "weekly_no_heading_blocks",
		},
		{
			name: "Save weekly with empty heading blocks",
			headingBlocks: []*markdown.HeadingBlock{
				createTestHeadingBlock(2, "Empty Section", ""),
			},
			goldenFile: "weekly_empty_heading_blocks",
		},
		{
			name: "Save weekly with complex heading structure",
			headingBlocks: []*markdown.HeadingBlock{
				createTestHeadingBlock(1, "Main Section", "Main content\n"),
				createTestHeadingBlock(2, "Subsection 1", "Subsection content\n"),
				createTestHeadingBlock(3, "Sub-subsection", "Detailed content\n"),
				createTestHeadingBlock(2, "Subsection 2", "Another subsection\n"),
			},
			goldenFile: "weekly_complex_structure",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			weekly := createTestWeekly(t, tc.headingBlocks)

			err := repo.Save(weekly, false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			filePath := filepath.Join(tmpDir, weekly.FileName())
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

// TestWeeklyRepoImpl_Save_FileSystem tests file system related aspects
func TestWeeklyRepoImpl_Save_FileSystem(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewWeekly(tmpDir, logger)

	t.Run("Create weekly file in correct location", func(t *testing.T) {
		weekly := createTestWeekly(t, []*markdown.HeadingBlock{
			createTestHeadingBlock(2, "Test Section", "Test content\n"),
		})

		err := repo.Save(weekly, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify file exists with correct name
		expectedFileName := "weekly_report.md"
		filePath := filepath.Join(tmpDir, expectedFileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("expected file to exist at %s", filePath)
		}

		// Verify file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		expectedContent := readGoldenFile(t, "weekly_file_location")
		if string(content) != expectedContent {
			t.Errorf("file content mismatch.\nExpected: %q\nGot: %q", expectedContent, string(content))
		}
	})

	t.Run("Multiple saves with different content", func(t *testing.T) {
		// First save
		weekly1 := createTestWeekly(t, []*markdown.HeadingBlock{
			createTestHeadingBlock(2, "First Save", "First content\n"),
		})

		err := repo.Save(weekly1, false)
		if err != nil {
			t.Fatalf("failed to save first weekly: %v", err)
		}

		// Second save with truncate = true
		weekly2 := createTestWeekly(t, []*markdown.HeadingBlock{
			createTestHeadingBlock(2, "Second Save", "Second content\n"),
		})

		err = repo.Save(weekly2, true)
		if err != nil {
			t.Fatalf("failed to save second weekly: %v", err)
		}

		// Verify second content was saved
		filePath := filepath.Join(tmpDir, weekly2.FileName())
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		expectedContent := readGoldenFile(t, "weekly_multiple_saves")
		if string(content) != expectedContent {
			t.Errorf("file content mismatch.\nExpected: %q\nGot: %q", expectedContent, string(content))
		}
	})
}
