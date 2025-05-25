package platform

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExists_FileExists(t *testing.T) {
	// Setup
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	err := os.WriteFile(tmpFile, []byte("test"), 0o644)
	require.NoError(t, err)

	// Execute
	exists := Exists(tmpFile)

	// Assert
	assert.True(t, exists)
}

func TestExists_FileDoesNotExist(t *testing.T) {
	// Setup
	nonExistentPath := filepath.Join(t.TempDir(), "nonexistent.txt")

	// Execute
	exists := Exists(nonExistentPath)

	// Assert
	assert.False(t, exists)
}

func TestExists_DirectoryExists(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()

	// Execute
	exists := Exists(tmpDir)

	// Assert
	assert.True(t, exists)
}

func TestExists_EmptyPath(t *testing.T) {
	// Execute
	exists := Exists("")

	// Assert
	assert.False(t, exists)
}

func TestEnsureDir_CreatesNewDirectory(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "subdir")

	// Execute
	err := EnsureDir(newDir)

	// Assert
	require.NoError(t, err)
	assert.True(t, Exists(newDir))
}

func TestEnsureDir_CreatesNestedDirectories(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "level1", "level2", "level3")

	// Execute
	err := EnsureDir(nestedDir)

	// Assert
	require.NoError(t, err)
	assert.True(t, Exists(nestedDir))
}

func TestEnsureDir_ExistingDirectory(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()

	// Execute
	err := EnsureDir(tmpDir)

	// Assert
	require.NoError(t, err)
	assert.True(t, Exists(tmpDir))
}

func TestEnsureDir_EmptyPath(t *testing.T) {
	// Execute
	err := EnsureDir("")

	// Assert
	require.NoError(t, err)
}

func TestWriteFileStream_CreateNewFile(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "Hello, World!"

	// Execute
	err := WriteFileStream(filePath, false, func(w *bufio.Writer) error {
		_, err := w.WriteString(content)
		return err
	})

	// Assert
	require.NoError(t, err)
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, content, string(data))
}

func TestWriteFileStream_TruncateExistingFile(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	// Create initial file
	err := os.WriteFile(filePath, []byte("old content"), 0o644)
	require.NoError(t, err)

	newContent := "new content"

	// Execute
	err = WriteFileStream(filePath, true, func(w *bufio.Writer) error {
		_, err := w.WriteString(newContent)
		return err
	})

	// Assert
	require.NoError(t, err)
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, newContent, string(data))
}

func TestWriteFileStream_NoTruncateExistingFile(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	originalContent := "original content"

	// Create initial file
	err := os.WriteFile(filePath, []byte(originalContent), 0o644)
	require.NoError(t, err)

	// Execute
	err = WriteFileStream(filePath, false, func(w *bufio.Writer) error {
		_, err := w.WriteString("should not write")
		return err
	})

	// Assert
	require.NoError(t, err)
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, originalContent, string(data), "File content should remain unchanged")
}

func TestWriteFileStream_EmptyPath(t *testing.T) {
	// Execute
	err := WriteFileStream("", false, func(w *bufio.Writer) error {
		return nil
	})

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "path is empty")
}

func TestWriteFileStream_CreatesParentDirectory(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "subdir1", "subdir2", "test.txt")
	content := "test content"

	// Execute
	err := WriteFileStream(filePath, false, func(w *bufio.Writer) error {
		_, err := w.WriteString(content)
		return err
	})

	// Assert
	require.NoError(t, err)
	assert.True(t, Exists(filepath.Dir(filePath)))
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, content, string(data))
}

func TestWriteFileStream_WriteError(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	// Execute
	err := WriteFileStream(filePath, false, func(w *bufio.Writer) error {
		return assert.AnError
	})

	// Assert
	require.Error(t, err)
}

func TestWriteFileStream_LargeContent(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "large.txt")

	// Execute - write 10MB of data
	err := WriteFileStream(filePath, false, func(w *bufio.Writer) error {
		for i := 0; i < 10000; i++ {
			if _, err := w.WriteString("This is a line of test data that will be repeated many times.\n"); err != nil {
				return err
			}
		}
		return nil
	})

	// Assert
	require.NoError(t, err)
	assert.True(t, Exists(filePath))

	info, err := os.Stat(filePath)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(100000), "File should be large")
}

