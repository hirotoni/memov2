package memo

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/repositories"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateMemoIndex_Success(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := config.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := config.NewTomlConfig(opts)
	require.NoError(t, err)

	// Create some test memos
	memosDir := cfg.MemosDir()
	err = os.MkdirAll(memosDir, 0o755)
	require.NoError(t, err)

	// Create a test memo file
	date := time.Now()
	memoFile, err := domain.NewMemoFile(date, "Test Memo", []string{})
	require.NoError(t, err)

	repos := repositories.NewRepositories(interfaces.NewConfigProvider(cfg))
	err = repos.Memo().Save(memoFile, false)
	require.NoError(t, err)

	mockEditor := mock.NewMockEditor()
	uc := NewMemo(interfaces.NewConfigProvider(cfg), repos, mockEditor)

	// Execute
	err = uc.GenerateMemoIndex()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 1, len(mockEditor.Calls), "Editor should be called once")

	// Verify index file was created
	indexPath := filepath.Join(memosDir, "index.md")
	assert.FileExists(t, indexPath)

	// Verify index content
	content, err := os.ReadFile(indexPath)
	require.NoError(t, err)
	// The index contains the filename which has the title in a transformed format
	assert.Contains(t, string(content), "Test")
}

func TestGenerateMemoIndex_EmptyMemos(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := config.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := config.NewTomlConfig(opts)
	require.NoError(t, err)

	// Create empty memos directory
	err = os.MkdirAll(cfg.MemosDir(), 0o755)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	repos := repositories.NewRepositories(configProvider)
	mockEditor := mock.NewMockEditor()
	uc := NewMemo(configProvider, repos, mockEditor)

	// Execute
	err = uc.GenerateMemoIndex()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 1, len(mockEditor.Calls), "Editor should be called once")

	// Verify index file was created
	indexPath := filepath.Join(cfg.MemosDir(), "index.md")
	assert.FileExists(t, indexPath)
}

func TestGenerateMemoIndex_EditorError(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := config.TomlConfigOption{BaseDir: tmpDir}
	cfg, err := config.NewTomlConfig(opts)
	require.NoError(t, err)

	// Create empty memos directory
	err = os.MkdirAll(cfg.MemosDir(), 0o755)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	repos := repositories.NewRepositories(configProvider)
	mockEditor := mock.NewMockEditor()

	// Set up editor to return error
	mockEditor.OpenFunc = func(basedir, path string) error {
		return assert.AnError
	}

	uc := NewMemo(configProvider, repos, mockEditor)

	// Execute
	err = uc.GenerateMemoIndex()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error opening editor")
}
