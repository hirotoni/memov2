package memo

import (
	"testing"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/repositories"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	"github.com/stretchr/testify/require"
)

func TestBrowse_Integration(t *testing.T) {
	// Skip integration test as it would launch the TUI
	t.Skip("Skipping TUI integration test")

	// Setup
	tmpDir := t.TempDir()
	opts := config.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := config.NewTomlConfig(opts)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	repos := repositories.NewRepositories(configProvider)
	mockEditor := mock.NewMockEditor()
	uc := NewMemo(configProvider, repos, mockEditor)

	// Execute
	err = uc.Browse()

	// Assert - would depend on TUI interaction
	_ = err
}
