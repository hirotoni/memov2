package interfaces

import "github.com/hirotoni/memov2/internal/config"

// ConfigProviderImpl wraps TomlConfig to implement ConfigProvider
type ConfigProviderImpl struct {
	config *config.TomlConfig
}

// NewConfigProvider creates a new ConfigProvider from TomlConfig
func NewConfigProvider(config *config.TomlConfig) ConfigProvider {
	return &ConfigProviderImpl{config: config}
}

// BaseDir returns the base directory path
func (cp *ConfigProviderImpl) BaseDir() string {
	return cp.config.BaseDir()
}

// TodosDir returns the todos directory path
func (cp *ConfigProviderImpl) TodosDir() string {
	return cp.config.TodosDir()
}

// MemosDir returns the memos directory path
func (cp *ConfigProviderImpl) MemosDir() string {
	return cp.config.MemosDir()
}

// TodosDaysToSeek returns the number of days to seek for todos
func (cp *ConfigProviderImpl) TodosDaysToSeek() int {
	return cp.config.TodosDaysToSeek()
}

// ConfigDirPath returns the config directory and file path
func (cp *ConfigProviderImpl) ConfigDirPath() (string, string, error) {
	return cp.config.ConfigDirPath()
}

// GetTomlConfig returns the underlying TomlConfig for cases where it's needed
// This method should be used sparingly and only when absolutely necessary
func (cp *ConfigProviderImpl) GetTomlConfig() interface{} {
	return cp.config
}
