package toml

import (
	"github.com/hirotoni/memov2/internal/interfaces"
)

// Provider wraps Config to implement ConfigProvider interface
type Provider struct {
	config *Config
}

// NewProvider creates a new ConfigProvider from Config
func NewProvider(cfg *Config) interfaces.ConfigProvider {
	return &Provider{config: cfg}
}

// BaseDir returns the base directory path
func (p *Provider) BaseDir() string {
	return p.config.BaseDir()
}

// TodosDir returns the todos directory path
func (p *Provider) TodosDir() string {
	return p.config.TodosDir()
}

// MemosDir returns the memos directory path
func (p *Provider) MemosDir() string {
	return p.config.MemosDir()
}

// TodosDaysToSeek returns the number of days to seek for todos
func (p *Provider) TodosDaysToSeek() int {
	return p.config.TodosDaysToSeek()
}

// ConfigDirPath returns the config directory and file path
func (p *Provider) ConfigDirPath() (string, string, error) {
	return p.config.ConfigDirPath()
}

// GetTomlConfig returns the underlying Config for cases where it's needed
// This method should be used sparingly and only when absolutely necessary
func (p *Provider) GetTomlConfig() interface{} {
	return p.config
}
