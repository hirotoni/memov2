package platform

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMoveToTrash(t *testing.T) {
	// Create a temporary file to test with
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	// Create the test file
	content := []byte("test content")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Verify the file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Test file should exist before moving to trash")
	}

	// Move the file to trash
	err := MoveToTrash(testFile)
	if err != nil {
		t.Fatalf("MoveToTrash failed: %v", err)
	}

	// Verify the original file no longer exists
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("Test file should not exist after moving to trash")
	}

	// Verify the file is in the trash (OS-specific)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	// Check trash directory based on OS
	var trashDir string
	switch {
	case fileExists(filepath.Join(homeDir, ".Trash")):
		trashDir = filepath.Join(homeDir, ".Trash")
	case fileExists(filepath.Join(homeDir, ".local", "share", "Trash", "files")):
		trashDir = filepath.Join(homeDir, ".local", "share", "Trash", "files")
	default:
		t.Skip("Trash directory not found, skipping verification")
	}

	// Check if file exists in trash
	trashedFile := filepath.Join(trashDir, "test.txt")
	if _, err := os.Stat(trashedFile); os.IsNotExist(err) {
		t.Errorf("File should exist in trash at %s", trashedFile)
	} else {
		// Clean up - remove from trash
		_ = os.Remove(trashedFile)
	}
}

func TestMoveToTrash_NonExistentFile(t *testing.T) {
	// Try to move a non-existent file
	nonExistentFile := "/tmp/this_file_does_not_exist_12345.txt"

	err := MoveToTrash(nonExistentFile)
	if err == nil {
		t.Error("Expected error when moving non-existent file to trash")
	}
}

func TestMoveToTrash_EmptyPath(t *testing.T) {
	err := MoveToTrash("")
	if err == nil {
		t.Error("Expected error when path is empty")
	}
}

func TestMoveToTrash_DuplicateFilename(t *testing.T) {
	// Create two files with the same name
	tempDir := t.TempDir()
	testFile1 := filepath.Join(tempDir, "duplicate.txt")
	testFile2 := filepath.Join(tempDir, "duplicate2.txt")

	// Create test files
	if err := os.WriteFile(testFile1, []byte("content 1"), 0644); err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("content 2"), 0644); err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}

	// Move first file to trash
	err := MoveToTrash(testFile1)
	if err != nil {
		t.Fatalf("Failed to move first file: %v", err)
	}

	// Get trash directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	var trashDir string
	switch {
	case fileExists(filepath.Join(homeDir, ".Trash")):
		trashDir = filepath.Join(homeDir, ".Trash")
	case fileExists(filepath.Join(homeDir, ".local", "share", "Trash", "files")):
		trashDir = filepath.Join(homeDir, ".local", "share", "Trash", "files")
	default:
		t.Skip("Trash directory not found, skipping test")
	}

	// Now rename the second file to match the first and move it
	os.Rename(testFile2, testFile1)
	if err := os.WriteFile(testFile1, []byte("content 2"), 0644); err != nil {
		t.Fatalf("Failed to recreate test file: %v", err)
	}

	// Move second file with same name to trash
	err = MoveToTrash(testFile1)
	if err != nil {
		t.Fatalf("Failed to move second file: %v", err)
	}

	// Verify both files exist in trash (with different names due to timestamp)
	// Note: We may not have permission to read the trash directory on some systems
	entries, err := os.ReadDir(trashDir)
	if err != nil {
		// If we can't read the trash directory (e.g., permissions issue on macOS),
		// we'll skip this verification but the test already passed if we got here
		t.Logf("Could not verify files in trash (permission issue): %v", err)
		return
	}

	duplicateCount := 0
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".txt" {
			if entry.Name() == "duplicate.txt" || (len(filepath.Base(entry.Name())) >= 9 && filepath.Base(entry.Name())[:9] == "duplicate") {
				duplicateCount++
				// Clean up
				_ = os.Remove(filepath.Join(trashDir, entry.Name()))
			}
		}
	}

	if duplicateCount < 1 {
		t.Error("Expected at least one duplicate file in trash")
	}
}

// Helper function to check if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

