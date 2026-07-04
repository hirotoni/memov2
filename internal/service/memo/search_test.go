package memo

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch_JSON(t *testing.T) {
	uc, _ := newTestMemoService(t)
	require.NoError(t, uc.GenerateMemoFile("GolangNotes", []string{"work"}))
	require.NoError(t, uc.GenerateMemoFile("CookingRecipes", []string{"home"}))

	out := captureStdout(t, func() {
		require.NoError(t, uc.Search("golang", true))
	})

	var items []struct {
		Title    string   `json:"title"`
		Path     string   `json:"path"`
		Category []string `json:"category"`
		Matches  []struct {
			Type    string `json:"type"`
			Text    string `json:"text"`
			Heading string `json:"heading"`
			Line    int    `json:"line"`
		} `json:"matches"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &items), "output must be valid JSON: %q", out)
	require.Len(t, items, 1)
	assert.Equal(t, "GolangNotes", items[0].Title)
	require.NotEmpty(t, items[0].Matches)
	assert.Equal(t, "title", items[0].Matches[0].Type)
}

func TestSearch_Text(t *testing.T) {
	uc, _ := newTestMemoService(t)
	require.NoError(t, uc.GenerateMemoFile("GolangNotes", []string{"work"}))

	out := captureStdout(t, func() {
		require.NoError(t, uc.Search("golang", false))
	})

	assert.Contains(t, out, "GolangNotes")
	assert.Contains(t, out, "\t")
	assert.NotContains(t, out, "\"title\"")
}

func TestSearch_NoMatch_JSON(t *testing.T) {
	uc, _ := newTestMemoService(t)
	require.NoError(t, uc.GenerateMemoFile("GolangNotes", []string{"work"}))

	out := captureStdout(t, func() {
		require.NoError(t, uc.Search("zzznomatch", true))
	})

	var items []interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &items))
	assert.Empty(t, items)
}

func TestSearch_Empty_Dir_JSON(t *testing.T) {
	uc, _ := newTestMemoService(t)

	out := captureStdout(t, func() {
		require.NoError(t, uc.Search("anything", true))
	})

	var items []interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &items))
	assert.Empty(t, items)
}

// An empty (or whitespace-only) query must be rejected as a validation error
// and must not enumerate the corpus (guards against dumping every memo).
func TestSearch_EmptyQuery_Rejected(t *testing.T) {
	uc, _ := newTestMemoService(t)
	require.NoError(t, uc.GenerateMemoFile("GolangNotes", []string{"work"}))

	for _, q := range []string{"", "   ", "\t"} {
		var err error
		out := captureStdout(t, func() {
			err = uc.Search(q, true)
		})

		require.Error(t, err, "query %q must be rejected", q)
		var appErr *common.AppError
		require.True(t, errors.As(err, &appErr), "expected *common.AppError, got %T", err)
		assert.Equal(t, common.ErrorTypeValidation, appErr.Type)
		assert.Empty(t, strings.TrimSpace(out), "nothing must be written to stdout for %q", q)
	}
}

// A memo matching under several match types (title + category + content) must
// still appear exactly once in text output (one line per matched memo).
func TestSearch_DeduplicatesMemoAcrossMatchTypes(t *testing.T) {
	uc, _ := newTestMemoService(t)
	// "golang" appears in the title and in the category, so this memo matches
	// under multiple match types simultaneously.
	require.NoError(t, uc.GenerateMemoFile("GolangNotes", []string{"golang"}))

	out := captureStdout(t, func() {
		require.NoError(t, uc.Search("golang", false))
	})

	lines := 0
	for _, l := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if strings.TrimSpace(l) != "" {
			lines++
		}
	}
	assert.Equal(t, 1, lines, "matched memo must be deduplicated to a single line, got:\n%s", out)
}
