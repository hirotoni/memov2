package todo

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateTodoFileNonInteractive_Path(t *testing.T) {
	uc, _, todosDir, editor := newTestTodoServiceE(t)
	// Seed a recent previous file so inheritance succeeds.
	seedTodoFile(t, todosDir, time.Now().AddDate(0, 0, -1), "- carried over")

	out := captureStdout(t, func() {
		require.NoError(t, uc.GenerateTodoFileNonInteractive(false, false))
	})

	rel := strings.TrimSpace(out)
	require.NotEmpty(t, rel)
	assert.Equal(t, 1, len(strings.Split(rel, "\n")), "default output must be a single line")
	assert.True(t, strings.HasSuffix(rel, "_todos.md"), "path %q must be a todos file", rel)

	// The editor must NOT be opened in non-interactive mode.
	assert.Empty(t, editor.Calls, "editor must not be opened in non-interactive mode")
}

func TestGenerateTodoFileNonInteractive_JSON(t *testing.T) {
	uc, _, todosDir, editor := newTestTodoServiceE(t)
	seedTodoFile(t, todosDir, time.Now().AddDate(0, 0, -1), "- carried over")

	out := captureStdout(t, func() {
		require.NoError(t, uc.GenerateTodoFileNonInteractive(false, true))
	})

	var dto struct {
		Date string `json:"date"`
		Path string `json:"path"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &dto), "output must be valid JSON: %q", out)
	assert.True(t, strings.HasSuffix(dto.Path, "_todos.md"))
	assert.NotEmpty(t, dto.Date)
	assert.Empty(t, editor.Calls, "editor must not be opened in --json mode")
}
