package interfaces

// ConfigProvider abstracts configuration access to reduce coupling
type ConfigProvider interface {
	BaseDir() string
	TodosDir() string
	MemosDir() string
	TodosDaysToSeek() int
	ConfigDirPath() (string, string, error)
	// GetTomlConfig returns the underlying TomlConfig for cases where it's needed
	// This method should be used sparingly and only when absolutely necessary
	GetTomlConfig() interface{}
}
