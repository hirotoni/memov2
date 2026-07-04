package todo

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShow_JSON(t *testing.T) {
	uc, _, todosDir := newTestTodoService(t)
	date := time.Date(2023, 1, 2, 0, 0, 0, 0, time.Local)
	seedTodoFile(t, todosDir, date, "- inherited item")

	out := captureStdout(t, func() {
		require.NoError(t, uc.Show("20230102Mon_todos.md", true))
	})

	var dto struct {
		Date            string `json:"date"`
		Path            string `json:"path"`
		TopLevelContent string `json:"top_level_content"`
		Headings        []struct {
			Level   int    `json:"level"`
			Heading string `json:"heading"`
			Content string `json:"content"`
		} `json:"headings"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &dto), "output must be valid JSON: %q", out)
	assert.Equal(t, "20230102Mon_todos.md", dto.Path)
	assert.NotEmpty(t, dto.Date)
	assert.NotEmpty(t, dto.Headings)
}

func TestShow_Text_RawContent(t *testing.T) {
	uc, _, todosDir := newTestTodoService(t)
	seedTodoFile(t, todosDir, time.Date(2023, 1, 2, 0, 0, 0, 0, time.Local), "- raw item")

	out := captureStdout(t, func() {
		require.NoError(t, uc.Show("20230102Mon_todos.md", false))
	})

	// Text mode is the raw markdown file, not JSON.
	assert.Contains(t, out, "raw item")
	assert.NotContains(t, out, "\"top_level_content\"")
}

func TestShow_PathEscape_Rejected(t *testing.T) {
	uc, _, _ := newTestTodoService(t)

	out := captureStdout(t, func() {
		err := uc.Show("../../../etc/passwd", true)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "escapes")
	})
	// Nothing must leak to stdout on rejection.
	assert.Empty(t, out)
}

func TestShow_NotFound_Rejected(t *testing.T) {
	uc, _, todosDir := newTestTodoService(t)
	seedTodoFile(t, todosDir, time.Date(2023, 1, 2, 0, 0, 0, 0, time.Local), "- x")

	out := captureStdout(t, func() {
		err := uc.Show("20991231Fri_todos.md", true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "todo not found")
	})
	assert.Empty(t, out)
}
