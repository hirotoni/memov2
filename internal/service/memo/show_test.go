package memo

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// firstMemoRelPath generates a memo and returns its relative path (as produced
// by `list --short --json`) for use with Show.
func firstMemoRelPath(t *testing.T, uc interface {
	GenerateMemoFile(title string, categoryTree []string) error
	List(showFullPath, asJSON bool) error
}, title string, cat []string) string {
	t.Helper()
	require.NoError(t, uc.GenerateMemoFile(title, cat))
	out := captureStdout(t, func() {
		require.NoError(t, uc.List(false, true))
	})
	var items []struct {
		Title string `json:"title"`
		Path  string `json:"path"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &items))
	for _, it := range items {
		if it.Title == title {
			return it.Path
		}
	}
	t.Fatalf("memo %q not found in list output", title)
	return ""
}

func TestShow_JSON(t *testing.T) {
	uc, _ := newTestMemoService(t)
	rel := firstMemoRelPath(t, uc, "ShowMe", []string{"work"})

	out := captureStdout(t, func() {
		require.NoError(t, uc.Show(rel, true))
	})

	var dto struct {
		Title           string   `json:"title"`
		Path            string   `json:"path"`
		Category        []string `json:"category"`
		Date            string   `json:"date"`
		TopLevelContent string   `json:"top_level_content"`
		Headings        []struct {
			Level   int    `json:"level"`
			Heading string `json:"heading"`
			Content string `json:"content"`
		} `json:"headings"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &dto), "output must be valid JSON: %q", out)
	assert.Equal(t, "ShowMe", dto.Title)
	assert.Equal(t, []string{"work"}, dto.Category)
	assert.NotEmpty(t, dto.Date)
}

func TestShow_Text_RawContent(t *testing.T) {
	uc, _ := newTestMemoService(t)
	rel := firstMemoRelPath(t, uc, "RawContent", []string{})

	out := captureStdout(t, func() {
		require.NoError(t, uc.Show(rel, false))
	})

	// Text mode is the raw markdown file: frontmatter + title heading, not JSON.
	assert.Contains(t, out, "---")
	assert.Contains(t, out, "RawContent")
	assert.NotContains(t, out, "\"top_level_content\"")
}

func TestShow_PathEscape_Rejected(t *testing.T) {
	uc, _ := newTestMemoService(t)

	out := captureStdout(t, func() {
		err := uc.Show("../../../etc/passwd", true)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "escapes")
	})
	// Nothing must leak to stdout on rejection.
	assert.Empty(t, out)
}

func TestShow_NotFound_Rejected(t *testing.T) {
	uc, _ := newTestMemoService(t)
	// Create the memos dir so it exists, then request a non-existent memo.
	require.NoError(t, uc.GenerateMemoFile("Existing", []string{}))

	out := captureStdout(t, func() {
		err := uc.Show("nonexistent.md", true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "memo not found")
	})
	assert.Empty(t, out)
}
