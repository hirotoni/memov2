package todo

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildWeeklyReportTodosNonInteractive_Path(t *testing.T) {
	uc, _, todosDir, editor := newTestTodoServiceE(t)
	seedTodoFile(t, todosDir, time.Date(2023, 1, 2, 0, 0, 0, 0, time.Local), "- a")
	seedTodoFile(t, todosDir, time.Date(2023, 1, 3, 0, 0, 0, 0, time.Local), "- b")

	out := captureStdout(t, func() {
		require.NoError(t, uc.BuildWeeklyReportTodosNonInteractive(false))
	})

	rel := strings.TrimSpace(out)
	require.NotEmpty(t, rel)
	// Only the report path may appear on stdout; the "Building weekly report..."
	// progress line must NOT pollute it.
	assert.Equal(t, 1, len(strings.Split(rel, "\n")), "default output must be a single line")
	assert.NotContains(t, out, "Building weekly report")
	assert.Empty(t, editor.Calls, "editor must not be opened in non-interactive mode")
}

func TestBuildWeeklyReportTodosNonInteractive_JSON(t *testing.T) {
	uc, _, todosDir, editor := newTestTodoServiceE(t)
	seedTodoFile(t, todosDir, time.Date(2023, 1, 2, 0, 0, 0, 0, time.Local), "- a")
	seedTodoFile(t, todosDir, time.Date(2023, 1, 3, 0, 0, 0, 0, time.Local), "- b")

	out := captureStdout(t, func() {
		require.NoError(t, uc.BuildWeeklyReportTodosNonInteractive(true))
	})

	var dto struct {
		Path string `json:"path"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &dto), "output must be valid JSON: %q", out)
	assert.NotEmpty(t, dto.Path)
	assert.NotContains(t, out, "Building weekly report")
	assert.Empty(t, editor.Calls, "editor must not be opened in --json mode")
}
