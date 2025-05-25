package trash

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// MoveToTrash moves a file to the system's trash/recycle bin
func MoveToTrash(path string) error {
	if path == "" {
		return fmt.Errorf("path is empty")
	}

	// Check if file exists
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	switch runtime.GOOS {
	case "darwin":
		return moveToTrashMacOS(path)
	case "linux":
		return moveToTrashLinux(path)
	case "windows":
		return moveToTrashWindows(path)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// moveToTrashMacOS moves a file to macOS Trash
func moveToTrashMacOS(path string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	trashDir := filepath.Join(homeDir, ".Trash")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return fmt.Errorf("failed to create trash directory: %w", err)
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
		return fmt.Errorf("failed to move file to trash: %w", err)
	}

	fmt.Printf("Moved to trash: %s -> %s\n", path, trashPath)
	return nil
}

// moveToTrashLinux moves a file to Linux trash following XDG specification
func moveToTrashLinux(path string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// XDG trash directory
	trashDir := filepath.Join(homeDir, ".local", "share", "Trash", "files")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return fmt.Errorf("failed to create trash directory: %w", err)
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
		return fmt.Errorf("failed to move file to trash: %w", err)
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
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	trashDir := filepath.Join(homeDir, ".Trash")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return fmt.Errorf("failed to create trash directory: %w", err)
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
		return fmt.Errorf("failed to move file to trash: %w", err)
	}

	fmt.Printf("Moved to trash: %s -> %s\n", path, trashPath)
	return nil
}
