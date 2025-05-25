package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/hirotoni/memov2/internal/common"
)

// MoveToTrash moves a file to the system's trash/recycle bin
// In test environments, it uses a test trash directory if TEST_TRASH_DIR is set
func MoveToTrash(path string) error {
	if path == "" {
		return common.New(common.ErrorTypeFileSystem, "path is empty")
	}

	// Check if file exists
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return common.New(common.ErrorTypeFileSystem, fmt.Sprintf("file does not exist: %s", path))
		}
		return common.Wrap(err, common.ErrorTypeFileSystem, fmt.Sprintf("failed to stat file: %s", path))
	}

	// Check if we're in a test environment
	if testTrashDir := os.Getenv("TEST_TRASH_DIR"); testTrashDir != "" {
		return moveToTestTrash(path, testTrashDir)
	}

	switch runtime.GOOS {
	case "darwin":
		return moveToTrashMacOS(path)
	case "linux":
		return moveToTrashLinux(path)
	case "windows":
		return moveToTrashWindows(path)
	default:
		return common.New(common.ErrorTypeFileSystem, fmt.Sprintf("unsupported operating system: %s", runtime.GOOS))
	}
}

// moveToTestTrash moves a file to a test trash directory
func moveToTestTrash(path, trashDir string) error {
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to create test trash directory")
	}

	// Get the base filename
	filename := filepath.Base(path)

	// Create a unique filename if it already exists in trash
	trashPath := filepath.Join(trashDir, filename)
	if _, err := os.Stat(trashPath); err == nil {
		// File exists, append timestamp
		ext := filepath.Ext(filename)
		nameWithoutExt := filename[:len(filename)-len(ext)]
		timestamp := time.Now().Format("20060102_150405")
		filename = fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext)
		trashPath = filepath.Join(trashDir, filename)
	}

	// Move the file
	if err := os.Rename(path, trashPath); err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to move file to test trash")
	}

	return nil
}

// moveToTrashMacOS moves a file to macOS Trash
func moveToTrashMacOS(path string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to get user home directory")
	}

	trashDir := filepath.Join(homeDir, ".Trash")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to create trash directory")
	}

	// Get the base filename
	filename := filepath.Base(path)

	// Create a unique filename if it already exists in trash
	trashPath := filepath.Join(trashDir, filename)
	if _, err := os.Stat(trashPath); err == nil {
		// File exists, append timestamp
		ext := filepath.Ext(filename)
		nameWithoutExt := filename[:len(filename)-len(ext)]
		timestamp := time.Now().Format("20060102_150405")
		filename = fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext)
		trashPath = filepath.Join(trashDir, filename)
	}

	// Move the file
	if err := os.Rename(path, trashPath); err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to move file to trash")
	}

	fmt.Printf("Moved to trash: %s -> %s\n", path, trashPath)
	return nil
}

// moveToTrashLinux moves a file to Linux trash following XDG specification
func moveToTrashLinux(path string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to get user home directory")
	}

	// XDG trash directory
	trashDir := filepath.Join(homeDir, ".local", "share", "Trash", "files")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to create trash directory")
	}

	// Get the base filename
	filename := filepath.Base(path)

	// Create a unique filename if it already exists in trash
	trashPath := filepath.Join(trashDir, filename)
	if _, err := os.Stat(trashPath); err == nil {
		// File exists, append timestamp
		ext := filepath.Ext(filename)
		nameWithoutExt := filename[:len(filename)-len(ext)]
		timestamp := time.Now().Format("20060102_150405")
		filename = fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext)
		trashPath = filepath.Join(trashDir, filename)
	}

	// Move the file
	if err := os.Rename(path, trashPath); err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to move file to trash")
	}

	fmt.Printf("Moved to trash: %s -> %s\n", path, trashPath)
	return nil
}

// moveToTrashWindows moves a file to Windows Recycle Bin
func moveToTrashWindows(path string) error {
	// For Windows, we'll use a simple approach of moving to a .Trash folder
	// A full implementation would use the Windows API to move to Recycle Bin
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to get user home directory")
	}

	trashDir := filepath.Join(homeDir, ".Trash")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to create trash directory")
	}

	// Get the base filename
	filename := filepath.Base(path)

	// Create a unique filename if it already exists in trash
	trashPath := filepath.Join(trashDir, filename)
	if _, err := os.Stat(trashPath); err == nil {
		// File exists, append timestamp
		ext := filepath.Ext(filename)
		nameWithoutExt := filename[:len(filename)-len(ext)]
		timestamp := time.Now().Format("20060102_150405")
		filename = fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext)
		trashPath = filepath.Join(trashDir, filename)
	}

	// Move the file
	if err := os.Rename(path, trashPath); err != nil {
		return common.Wrap(err, common.ErrorTypeFileSystem, "failed to move file to trash")
	}

	fmt.Printf("Moved to trash: %s -> %s\n", path, trashPath)
	return nil
}

