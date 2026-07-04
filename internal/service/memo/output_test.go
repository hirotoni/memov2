package memo

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/repositories"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	"github.com/stretchr/testify/require"
)

// newTestMemoService builds a memo service backed by a fresh temp memos dir.
// Logs go to a discard logger so stderr stays clean while capturing stdout.
func newTestMemoService(t *testing.T) (interfaces.MemoService, string) {
	t.Helper()
	tmpDir := t.TempDir()
	opts := toml.Option{BaseDir: tmpDir}
	cfg, err := toml.NewConfig(opts)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repos := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()
	uc := NewMemo(configProvider, repos, mockEditor, logger)
	return uc, tmpDir
}

// captureStdout redirects os.Stdout for the duration of fn and returns whatever
// was written to it. It fully drains the pipe so large outputs are not truncated.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()

	w.Close()
	os.Stdout = oldStdout
	return <-done
}
