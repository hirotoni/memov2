package interfaces

import "os"

// Editor defines the interface for editor operations
type Editor interface {
	Open(basedir, path string) error
}

// FileSystem defines the interface for file system operations
type FileSystem interface {
	ReadDir(dirname string) ([]os.DirEntry, error)
}
