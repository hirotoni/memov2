package memo

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList_JSON(t *testing.T) {
	uc, _ := newTestMemoService(t)
	require.NoError(t, uc.GenerateMemoFile("MemoA", []string{"work"}))
	require.NoError(t, uc.GenerateMemoFile("MemoB", []string{"work", "projects"}))

	out := captureStdout(t, func() {
		require.NoError(t, uc.List(true, true))
	})

	var items []struct {
		Title    string    `json:"title"`
		Path     string    `json:"path"`
		Category []string  `json:"category"`
		Date     time.Time `json:"date"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &items), "output must be valid JSON: %q", out)
	require.Len(t, items, 2)

	titles := map[string][]string{}
	for _, it := range items {
		titles[it.Title] = it.Category
		assert.NotEmpty(t, it.Path)
		assert.False(t, it.Date.IsZero())
	}
	assert.Equal(t, []string{"work"}, titles["MemoA"])
	assert.Equal(t, []string{"work", "projects"}, titles["MemoB"])
}

func TestList_JSON_Empty(t *testing.T) {
	uc, _ := newTestMemoService(t)

	out := captureStdout(t, func() {
		require.NoError(t, uc.List(true, true))
	})

	// Must be a valid empty JSON array, not null and not a crash.
	var items []interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &items))
	assert.Empty(t, items)
}

func TestList_Text_Unchanged(t *testing.T) {
	uc, _ := newTestMemoService(t)
	require.NoError(t, uc.GenerateMemoFile("MemoA", []string{"work"}))

	out := captureStdout(t, func() {
		require.NoError(t, uc.List(false, false))
	})

	// Text mode: "title\tpath\n"; not JSON.
	assert.Contains(t, out, "MemoA")
	assert.Contains(t, out, "\t")
	assert.NotContains(t, out, "\"title\"")
}
