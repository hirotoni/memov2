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

func TestListCategories_Success(t *testing.T) {
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

	// Create memos with categories to populate category list
	err = uc.GenerateMemoFile("Memo A", []string{"work"})
	require.NoError(t, err)
	err = uc.GenerateMemoFile("Memo B", []string{"work", "projects"})
	require.NoError(t, err)
	err = uc.GenerateMemoFile("Memo C", []string{"personal"})
	require.NoError(t, err)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	// Execute
	err = uc.ListCategories()
	require.NoError(t, err)

	// Read captured output
	w.Close()
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	os.Stdout = oldStdout
	output := string(buf[:n])

	// Assert
	assert.Contains(t, output, "work\n")
	assert.Contains(t, output, "work/projects\n")
	assert.Contains(t, output, "personal\n")
}

func TestListCategories_Empty(t *testing.T) {
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

	// Create a memo at root (no category) so memos dir exists
	err = uc.GenerateMemoFile("Root Memo", []string{})
	require.NoError(t, err)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	// Execute
	err = uc.ListCategories()
	require.NoError(t, err)

	// Read captured output
	w.Close()
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	os.Stdout = oldStdout
	output := string(buf[:n])

	// Assert - no category output when all memos are at root
	assert.Empty(t, output)
}
