package todo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// readTodayTodoFile returns the content of today's daily todo file in todosDir.
func readTodayTodoFile(t *testing.T, todosDir string) (string, string) {
	t.Helper()
	today, err := domain.NewTodosFile(time.Now())
	require.NoError(t, err)
	data, err := os.ReadFile(filepath.Join(todosDir, today.FileName()))
	require.NoError(t, err, "today's todo file must be created")
	return today.FileName(), string(data)
}

// TestGenerateTodoFileNonInteractive_NoPreviousFile: with no daily file at all
// (first run), today's file is still created from the template and its path is
// printed.
func TestGenerateTodoFileNonInteractive_NoPreviousFile(t *testing.T) {
	uc, _, todosDir := newTestTodoService(t)

	out := captureStdout(t, func() {
		require.NoError(t, uc.GenerateTodoFileNonInteractive(false, false))
	})

	name, content := readTodayTodoFile(t, todosDir)
	assert.Equal(t, name, strings.TrimSpace(out))
	assert.Contains(t, content, "## "+utils.HeadingTodos.Text)
	assert.Contains(t, content, "## "+utils.HeadingBacklog.Text)
}

// TestGenerateTodoFileNonInteractive_OnlyOlderThanDaysToSeek: when the only
// daily file is older than todos_daystoseek, today's file is created without
// inheriting from it.
func TestGenerateTodoFileNonInteractive_OnlyOlderThanDaysToSeek(t *testing.T) {
	uc, _, todosDir := newTestTodoService(t)

	const stale = "- stale item beyond seek window"
	seedTodoFile(t, todosDir, time.Now().AddDate(0, 0, -(config.DefaultTodosDaysToSeek+1)), stale)

	out := captureStdout(t, func() {
		require.NoError(t, uc.GenerateTodoFileNonInteractive(false, false))
	})

	name, content := readTodayTodoFile(t, todosDir)
	assert.Equal(t, name, strings.TrimSpace(out))
	assert.NotContains(t, content, stale)
}

// TestGenerateTodoFileNonInteractive_InheritsWithinDaysToSeek: a daily file on
// the last day of the seek window is still inherited (todos and backlog).
func TestGenerateTodoFileNonInteractive_InheritsWithinDaysToSeek(t *testing.T) {
	uc, _, todosDir := newTestTodoService(t)

	const carriedTodo = "- carried todo"
	const carriedBacklog = "- carried backlog"
	seedTodoFileWithBacklog(t, todosDir, time.Now().AddDate(0, 0, -config.DefaultTodosDaysToSeek), carriedTodo, carriedBacklog)

	captureStdout(t, func() {
		require.NoError(t, uc.GenerateTodoFileNonInteractive(false, false))
	})

	_, content := readTodayTodoFile(t, todosDir)
	assert.Contains(t, content, carriedTodo)
	assert.Contains(t, content, carriedBacklog)
}
