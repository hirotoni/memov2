package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hirotoni/memov2/internal/config/toml"
)

// TestConfig creates a test configuration with a temporary directory
func TestConfig(t *testing.T) *toml.Config {
	t.Helper()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "memov2_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Clean up the temp directory after the test
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	// Create test configuration
	cfg, err := toml.NewConfig(toml.Option{
		BaseDir:         tempDir,
		TodosFolderName: "todos",
		MemosFolderName: "memos",
		TodosDaysToSeek: 7,
	})
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	return cfg
}

// CreateTestFile creates a test file with the given content
func CreateTestFile(t *testing.T, dir, filename, content string) string {
	t.Helper()

	filepath := filepath.Join(dir, filename)
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file %s: %v", filepath, err)
	}

	return filepath
}

// CreateTestDirectory creates a test directory
func CreateTestDirectory(t *testing.T, parentDir, dirName string) string {
	t.Helper()

	dirPath := filepath.Join(parentDir, dirName)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		t.Fatalf("Failed to create test directory %s: %v", dirPath, err)
	}

	return dirPath
}

// AssertFileExists checks if a file exists
func AssertFileExists(t *testing.T, filepath string) {
	t.Helper()

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		t.Errorf("Expected file to exist: %s", filepath)
	}
}

// AssertFileNotExists checks if a file does not exist
func AssertFileNotExists(t *testing.T, filepath string) {
	t.Helper()

	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		t.Errorf("Expected file to not exist: %s", filepath)
	}
}

// AssertDirectoryExists checks if a directory exists
func AssertDirectoryExists(t *testing.T, dirpath string) {
	t.Helper()

	info, err := os.Stat(dirpath)
	if err != nil {
		t.Errorf("Expected directory to exist: %s", dirpath)
		return
	}

	if !info.IsDir() {
		t.Errorf("Expected path to be a directory: %s", dirpath)
	}
}
