package interfaces

// Editor defines the interface for editor operations
type Editor interface {
	Open(basedir, path string) error
}
