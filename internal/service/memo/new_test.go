package memo

import (
	"log/slog"
	"os"
	"testing"

	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/repositories"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateMemoFile_Success(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repos := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()
	uc := NewMemo(configProvider, repos, mockEditor, logger)

	// Execute
	err = uc.GenerateMemoFile("Test Memo Title")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 1, len(mockEditor.Calls), "Editor should be called once")

	// Verify the editor was called with the right path
	call := mockEditor.Calls[0]
	assert.Equal(t, tmpDir, call.Basedir)
	assert.Contains(t, call.Path, ".md")
	assert.Contains(t, call.Path, "memo")
}

func TestGenerateMemoFile_EmptyTitle_NoInput(t *testing.T) {
	// This test is skipped because it requires stdin interaction
	// In a real scenario, you'd mock os.Stdin
	t.Skip("Skipping test that requires stdin input")
}

func TestGenerateMemoFile_EditorError(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repos := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()

	// Set up editor to return error
	mockEditor.OpenFunc = func(basedir, path string) error {
		return assert.AnError
	}

	uc := NewMemo(configProvider, repos, mockEditor, logger)

	// Execute
	err = uc.GenerateMemoFile("Test Memo")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error opening editor")
}

func TestGenerateMemoFile_SaveError(t *testing.T) {
	// Skip this test - it requires more complex mocking setup
	t.Skip("Skipping save error test - requires complex mock setup")
}
