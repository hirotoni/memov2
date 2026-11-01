package memo

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/repositories"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTidyMemos_SkipsInvalidCategory verifies that a memo file whose
// frontmatter contains an invalid category (e.g. hand-edited to contain
// "..") is left in place - not moved, and definitely not moved outside
// MemosDir() - while other, validly categorized memos are still tidied
// (moved into their proper category subdirectory) in the same pass.
func TestTidyMemos_SkipsInvalidCategory(t *testing.T) {
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repos := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()
	uc := NewMemo(configProvider, repos, mockEditor, logger)

	memosDir := configProvider.MemosDir()
	require.NoError(t, os.MkdirAll(memosDir, 0o755))

	// A memo with an invalid category in its frontmatter (simulating a
	// hand-edited file). It is written directly to disk - bypassing the
	// repository's Save(), which would now reject it - and placed at the
	// memos root so we can tell whether tidy moved it.
	badDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	badMemo, err := domain.NewMemoFile(badDate, "Bad Memo", []string{".."})
	require.NoError(t, err)
	badPath := filepath.Join(memosDir, badMemo.FileName())
	require.NoError(t, os.WriteFile(badPath, []byte(badMemo.ContentString()), 0o644))

	// A validly categorized memo, also written directly to disk at the
	// memos root, but whose frontmatter category says it belongs in a
	// subdirectory. Tidy should move it there.
	goodDate := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	goodMemo, err := domain.NewMemoFile(goodDate, "Good Memo", []string{"work"})
	require.NoError(t, err)
	goodSourcePath := filepath.Join(memosDir, goodMemo.FileName())
	require.NoError(t, os.WriteFile(goodSourcePath, []byte(goodMemo.ContentString()), 0o644))

	// Execute
	err = uc.TidyMemos()

	// Assert: tidy completes without error despite the bad file.
	require.NoError(t, err)

	// The invalid-category file was left exactly where it was.
	assert.FileExists(t, badPath, "file with invalid category should not be moved")

	// Nothing escaped MemosDir(): no memo file exists in its parent
	// directory (which is what a literal ".." join would have produced).
	parentDir := filepath.Dir(memosDir)
	parentEntries, err := os.ReadDir(parentDir)
	require.NoError(t, err)
	for _, entry := range parentEntries {
		assert.NotEqual(t, badMemo.FileName(), entry.Name(), "file must not have been moved outside MemosDir()")
	}

	// The validly categorized memo was tidied into its target location.
	goodTargetPath := filepath.Join(memosDir, "work", goodMemo.FileName())
	assert.FileExists(t, goodTargetPath, "valid memo should have been moved to its category directory")
	assert.NoFileExists(t, goodSourcePath, "valid memo should no longer exist at its original location")
}

// TestTidyMemos_MovesMemoInSubdirectoryToDeeperCategory reproduces adding a
// new level to a memo's frontmatter category (work -> work/project) while the
// file still sits in work/. Tidy must move it into the new work/project/.
func TestTidyMemos_MovesMemoInSubdirectoryToDeeperCategory(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, err := toml.NewConfig(toml.Option{BaseDir: tmpDir})
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repos := repositories.NewRepositories(configProvider, logger)
	uc := NewMemo(configProvider, repos, mock.NewMockEditor(), logger)

	memosDir := configProvider.MemosDir()
	m, err := domain.NewMemoFile(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), "Deeper Memo", []string{"work", "project"})
	require.NoError(t, err)
	sourcePath := filepath.Join(memosDir, "work", m.FileName())
	require.NoError(t, os.MkdirAll(filepath.Dir(sourcePath), 0o755))
	require.NoError(t, os.WriteFile(sourcePath, []byte(m.ContentString()), 0o644))

	require.NoError(t, uc.TidyMemos())

	assert.FileExists(t, filepath.Join(memosDir, "work", "project", m.FileName()))
	assert.NoFileExists(t, sourcePath)
}
