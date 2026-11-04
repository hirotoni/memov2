package todo

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/repositories"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	todoRepo "github.com/hirotoni/memov2/internal/repositories/todo"
	"github.com/hirotoni/memov2/internal/utils"
	"github.com/stretchr/testify/require"
)

// newTestTodoService builds a todo service backed by a fresh temp base dir. It
// returns the service, the base dir, and the todos dir (created on disk). Logs go
// to a discard logger so stderr stays clean while capturing stdout.
func newTestTodoService(t *testing.T) (interfaces.TodoService, string, string) {
	t.Helper()
	baseDir := t.TempDir()
	todosDir := filepath.Join(baseDir, config.DefaultFolderNameTodos)
	require.NoError(t, os.MkdirAll(todosDir, 0o755))

	cfg, err := toml.NewConfig(toml.Option{BaseDir: baseDir})
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repos := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()
	uc := NewTodo(configProvider, repos, mockEditor, logger)
	return uc, baseDir, todosDir
}

// newTestTodoServiceE is like newTestTodoService but also returns the mock editor
// so tests can assert whether the editor was (not) opened.
func newTestTodoServiceE(t *testing.T) (interfaces.TodoService, string, string, *mock.MockEditor) {
	t.Helper()
	baseDir := t.TempDir()
	todosDir := filepath.Join(baseDir, config.DefaultFolderNameTodos)
	require.NoError(t, os.MkdirAll(todosDir, 0o755))

	cfg, err := toml.NewConfig(toml.Option{BaseDir: baseDir})
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repos := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()
	uc := NewTodo(configProvider, repos, mockEditor, logger)
	return uc, baseDir, todosDir, mockEditor
}

// seedTodoFile writes a daily todo file for the given date directly into todosDir,
// with a todos heading carrying the given content. It is used to exercise
// inheritance from a previous day's file in GenerateTodoFile*.
func seedTodoFile(t *testing.T, todosDir string, date time.Time, todosContent string) {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	r := todoRepo.NewTodo(todosDir, logger)
	f, err := domain.NewTodosFile(date)
	require.NoError(t, err)
	f.SetHeadingBlocks([]*markdown.HeadingBlock{
		{HeadingText: utils.HeadingTodos.Text, ContentText: todosContent, Level: 2},
		{HeadingText: utils.HeadingBacklog.Text, Level: 2},
	})
	require.NoError(t, r.Save(f, true))
}

// seedTodoFileWithBacklog is like seedTodoFile but also fills the backlog heading.
func seedTodoFileWithBacklog(t *testing.T, todosDir string, date time.Time, todosContent, backlogContent string) {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	r := todoRepo.NewTodo(todosDir, logger)
	f, err := domain.NewTodosFile(date)
	require.NoError(t, err)
	f.SetHeadingBlocks([]*markdown.HeadingBlock{
		{HeadingText: utils.HeadingTodos.Text, ContentText: todosContent, Level: 2},
		{HeadingText: utils.HeadingBacklog.Text, ContentText: backlogContent, Level: 2},
	})
	require.NoError(t, r.Save(f, true))
}

// seedLegacyTodoFile is like seedTodoFile but writes the backlog section using the
// pre-rename heading text ("wanttodos") instead of "backlog". It simulates a daily
// todo file created before the heading rename (and before running `todos migrate`),
// so inheritance can be exercised against a legacy file.
func seedLegacyTodoFile(t *testing.T, todosDir string, date time.Time, backlogContent string) {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	r := todoRepo.NewTodo(todosDir, logger)
	f, err := domain.NewTodosFile(date)
	require.NoError(t, err)
	f.SetHeadingBlocks([]*markdown.HeadingBlock{
		{HeadingText: utils.HeadingTodos.Text, Level: 2},
		{HeadingText: utils.HeadingBacklogLegacy.Text, ContentText: backlogContent, Level: 2},
	})
	require.NoError(t, r.Save(f, true))
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
