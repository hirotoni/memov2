package memo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createForRename creates a memo and returns its relative path.
func createForRename(t *testing.T, uc interface {
	CreateNonInteractive(title string, categoryTree []string, body string, asJSON bool) error
}, title string, cat []string) string {
	t.Helper()
	out := captureStdout(t, func() {
		require.NoError(t, uc.CreateNonInteractive(title, cat, "", false))
	})
	return strings.TrimSpace(out)
}

func TestRename_NonInteractive_Success(t *testing.T) {
	uc, baseDir := newTestMemoService(t)
	memosDir := filepath.Join(baseDir, "memos")
	rel := createForRename(t, uc, "Old Title", []string{"work"})

	out := captureStdout(t, func() {
		require.NoError(t, uc.Rename(rel, "New Title", false))
	})

	newRel := strings.TrimSpace(out)
	assert.Equal(t, 1, len(strings.Split(newRel, "\n")), "default output must be a single line")
	assert.Contains(t, newRel, "_memo_New-Title.md")

	// Old file gone, new file present with updated in-file title.
	_, err := os.Stat(filepath.Join(memosDir, rel))
	assert.True(t, os.IsNotExist(err), "old file should be removed")

	content, err := os.ReadFile(filepath.Join(memosDir, newRel))
	require.NoError(t, err)
	assert.Contains(t, string(content), "# New Title")
}

func TestRename_NonInteractive_JSON(t *testing.T) {
	uc, _ := newTestMemoService(t)
	rel := createForRename(t, uc, "Before", []string{"cat"})

	out := captureStdout(t, func() {
		require.NoError(t, uc.Rename(rel, "After", true))
	})

	var dto struct {
		Path     string `json:"path"`
		OldTitle string `json:"old_title"`
		NewTitle string `json:"new_title"`
		NewPath  string `json:"new_path"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &dto), "output must be valid JSON: %q", out)
	assert.Equal(t, rel, dto.Path)
	assert.Equal(t, "Before", dto.OldTitle)
	assert.Equal(t, "After", dto.NewTitle)
	assert.Contains(t, dto.NewPath, "_memo_After.md")
}

func TestRename_NotFound_Error(t *testing.T) {
	uc, _ := newTestMemoService(t)
	// Ensure the memos dir exists.
	_ = createForRename(t, uc, "Existing", []string{})

	out := captureStdout(t, func() {
		err := uc.Rename("20990101Mon000000_memo_Ghost.md", "X", false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "memo not found")
	})
	assert.Empty(t, out)
}

func TestRename_PathEscape_Rejected(t *testing.T) {
	uc, _ := newTestMemoService(t)

	out := captureStdout(t, func() {
		err := uc.Rename("../../../etc/passwd", "X", false)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "escapes")
	})
	assert.Empty(t, out)
}

func TestRename_TraversalTitle_Rejected(t *testing.T) {
	uc, baseDir := newTestMemoService(t)
	memosDir := filepath.Join(baseDir, "memos")
	rel := createForRename(t, uc, "Legit", []string{"work"})

	out := captureStdout(t, func() {
		err := uc.Rename(rel, "../../evil", false)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "title")
	})
	assert.Empty(t, out, "nothing must leak to stdout on rejection")

	// The original memo must be untouched and no escaped file created.
	_, err := os.Stat(filepath.Join(memosDir, rel))
	assert.NoError(t, err, "original memo must remain after a rejected rename")
	_, err = os.Stat(filepath.Join(filepath.Dir(memosDir), "evil.md"))
	assert.True(t, os.IsNotExist(err), "traversal title must not write outside MemosDir")
}

func TestRename_EmptyTitle_Rejected(t *testing.T) {
	uc, baseDir := newTestMemoService(t)
	memosDir := filepath.Join(baseDir, "memos")
	rel := createForRename(t, uc, "Keep", []string{})

	out := captureStdout(t, func() {
		err := uc.Rename(rel, "   ", false)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "title is empty")
	})
	assert.Empty(t, out)

	// The original memo must remain; no malformed "<date>_memo_.md" created.
	_, err := os.Stat(filepath.Join(memosDir, rel))
	assert.NoError(t, err, "original memo must remain after a rejected rename")
}
