package fs

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Exists returns true if the file or directory at path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	if err != nil {
		// Avoid panicking/logging loudly at call sites; treat other errors as non-existent.
		return false
	}
	return true
}

// EnsureDir ensures the directory exists, creating it (and parents) if necessary.
func EnsureDir(dir string) error {
	if dir == "" {
		return nil
	}
	if Exists(dir) {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	return nil
}

// WriteFileStream creates or truncates the file at path and invokes the provided
// writer function with a buffered writer for streaming large contents.
//
// - If truncate is false and the file already exists, it returns nil without writing.
// - If truncate is true, the file is truncated if it exists.
// - The parent directory is created if missing.
func WriteFileStream(path string, truncate bool, write func(w *bufio.Writer) error) error {
	if path == "" {
		return fmt.Errorf("path is empty")
	}

	// Respect existing file if not truncating
	if !truncate && Exists(path) {
		return nil
	}

	// Ensure directory exists
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}

	flags := os.O_CREATE | os.O_WRONLY
	if truncate {
		flags |= os.O_TRUNC
	}

	f, err := os.OpenFile(path, flags, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer func() {
		_ = f.Close()
	}()

	bw := bufio.NewWriter(f)

	// Execute client write
	if err := write(bw); err != nil {
		// Attempt best-effort flush before returning
		_ = bw.Flush()
		return err
	}

	if err := bw.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffered writer for %s: %w", path, err)
	}

	return nil
}
