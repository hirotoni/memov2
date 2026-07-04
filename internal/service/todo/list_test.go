package todo

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList_JSON(t *testing.T) {
	uc, _, todosDir := newTestTodoService(t)
	d1 := time.Date(2023, 1, 2, 0, 0, 0, 0, time.Local)
	d2 := time.Date(2023, 1, 3, 0, 0, 0, 0, time.Local)
	seedTodoFile(t, todosDir, d1, "- a")
	seedTodoFile(t, todosDir, d2, "- b")

	out := captureStdout(t, func() {
		require.NoError(t, uc.List(true, true))
	})

	var items []struct {
		Date time.Time `json:"date"`
		Path string    `json:"path"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &items), "output must be valid JSON: %q", out)
	require.Len(t, items, 2)
	for _, it := range items {
		assert.NotEmpty(t, it.Path)
		assert.False(t, it.Date.IsZero())
	}
}

func TestList_JSON_Empty(t *testing.T) {
	uc, _, _ := newTestTodoService(t)

	out := captureStdout(t, func() {
		require.NoError(t, uc.List(true, true))
	})

	// Must be a valid empty JSON array, not null and not a crash.
	var items []interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &items))
	assert.Empty(t, items)
}

func TestList_Text_Unchanged(t *testing.T) {
	uc, _, todosDir := newTestTodoService(t)
	seedTodoFile(t, todosDir, time.Date(2023, 1, 2, 0, 0, 0, 0, time.Local), "- a")

	out := captureStdout(t, func() {
		require.NoError(t, uc.List(false, false))
	})

	// Text mode: "date\tpath\n"; not JSON.
	assert.Contains(t, out, "20230102Mon")
	assert.Contains(t, out, "\t")
	assert.NotContains(t, out, "\"path\"")
}
