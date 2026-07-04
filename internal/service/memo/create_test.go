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

// findMemoFile returns the absolute path of the single created memo file whose
// name matches the given title fragment, or fails the test.
func findMemoFile(t *testing.T, memosDir, titleFrag string) string {
	t.Helper()
	var found string
	err := filepath.Walk(memosDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") && strings.Contains(info.Name(), titleFrag) {
			found = path
		}
		return nil
	})
	require.NoError(t, err)
	require.NotEmpty(t, found, "expected a memo file containing %q under %s", titleFrag, memosDir)
	return found
}

func TestCreateNonInteractive_Success(t *testing.T) {
	uc, baseDir := newTestMemoService(t)
	memosDir := filepath.Join(baseDir, "memos")

	out := captureStdout(t, func() {
		require.NoError(t, uc.CreateNonInteractive("My Note", []string{"work", "proj"}, "# body line\nsecond", false))
	})

	// Default output: single relative-path line.
	rel := strings.TrimSpace(out)
	require.NotEmpty(t, rel)
	assert.Equal(t, 1, len(strings.Split(rel, "\n")), "default output must be a single line")
	assert.True(t, strings.HasPrefix(rel, filepath.Join("work", "proj")+string(filepath.Separator)), "path %q must be under category dir", rel)
	assert.Contains(t, rel, "_memo_My-Note.md")

	// File exists on disk under the category and contains title + body.
	fpath := findMemoFile(t, memosDir, "My-Note")
	assert.Equal(t, filepath.Join(memosDir, rel), fpath)
	content, err := os.ReadFile(fpath)
	require.NoError(t, err)
	s := string(content)
	assert.Contains(t, s, `category: ["work", "proj"]`)
	assert.Contains(t, s, "# My Note")
	assert.Contains(t, s, "# body line\nsecond")
}

func TestCreateNonInteractive_EmptyTitle_Error(t *testing.T) {
	uc, _ := newTestMemoService(t)

	out := captureStdout(t, func() {
		err := uc.CreateNonInteractive("   ", []string{}, "body", false)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "title is empty")
	})
	// Nothing must leak to stdout on rejection.
	assert.Empty(t, out)
}

func TestCreateNonInteractive_TraversalTitle_Rejected(t *testing.T) {
	uc, baseDir := newTestMemoService(t)
	memosDir := filepath.Join(baseDir, "memos")

	out := captureStdout(t, func() {
		err := uc.CreateNonInteractive("../../../pwned", nil, "leak", false)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "title")
	})
	assert.Empty(t, out, "nothing must leak to stdout on rejection")

	// No file must be written anywhere under (or outside) MemosDir.
	assertNoMarkdownWritten(t, baseDir)
	if _, err := os.Stat(memosDir); err == nil {
		assertNoMarkdownWritten(t, memosDir)
	}
	// The traversal target (baseDir's parent) must not have gained a pwned.md.
	_, err := os.Stat(filepath.Join(filepath.Dir(memosDir), "pwned.md"))
	assert.True(t, os.IsNotExist(err), "traversal title must not write outside MemosDir")
}

func TestCreateNonInteractive_InvalidCategory_Rejected(t *testing.T) {
	uc, baseDir := newTestMemoService(t)

	out := captureStdout(t, func() {
		err := uc.CreateNonInteractive("Note", []string{"..", "etc"}, "body", false)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "category")
	})
	assert.Empty(t, out, "nothing must leak to stdout on rejection")
	assertNoMarkdownWritten(t, baseDir)
}

func TestCreateNonInteractive_Collision_Rejected(t *testing.T) {
	uc, baseDir := newTestMemoService(t)
	memosDir := filepath.Join(baseDir, "memos")

	// First create succeeds.
	first := captureStdout(t, func() {
		require.NoError(t, uc.CreateNonInteractive("Dup", []string{"c"}, "first body", false))
	})
	rel := strings.TrimSpace(first)
	require.NotEmpty(t, rel)

	// Second create with the same title+category (same-second filename) must be
	// rejected, not silently reported as success while dropping the new body.
	out := captureStdout(t, func() {
		err := uc.CreateNonInteractive("Dup", []string{"c"}, "second body", false)
		require.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "already exists")
	})
	assert.Empty(t, out, "nothing must leak to stdout on collision")

	// The original file's body must be intact (the second body must not overwrite).
	content, err := os.ReadFile(filepath.Join(memosDir, rel))
	require.NoError(t, err)
	assert.Contains(t, string(content), "first body")
	assert.NotContains(t, string(content), "second body")
}

func TestCreateNonInteractive_TitleTrimmed(t *testing.T) {
	uc, baseDir := newTestMemoService(t)
	memosDir := filepath.Join(baseDir, "memos")

	out := captureStdout(t, func() {
		require.NoError(t, uc.CreateNonInteractive("  Padded Title  ", []string{}, "", false))
	})
	rel := strings.TrimSpace(out)
	// Filename must encode the trimmed title (spaces -> "-"), no leading/trailing fill.
	assert.Contains(t, rel, "_memo_Padded-Title.md")
	assert.NotContains(t, rel, "_memo_-Padded")

	content, err := os.ReadFile(filepath.Join(memosDir, rel))
	require.NoError(t, err)
	assert.Contains(t, string(content), "# Padded Title\n")
}

// assertNoMarkdownWritten fails if any .md file exists under root.
func assertNoMarkdownWritten(t *testing.T, root string) {
	t.Helper()
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			t.Errorf("unexpected markdown file written: %s", path)
		}
		return nil
	})
}

func TestCreateNonInteractive_JSON(t *testing.T) {
	uc, _ := newTestMemoService(t)

	out := captureStdout(t, func() {
		require.NoError(t, uc.CreateNonInteractive("JsonNote", []string{"cat"}, "hi", true))
	})

	var dto struct {
		Title    string   `json:"title"`
		Path     string   `json:"path"`
		Category []string `json:"category"`
		Date     string   `json:"date"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &dto), "output must be valid JSON: %q", out)
	assert.Equal(t, "JsonNote", dto.Title)
	assert.Equal(t, []string{"cat"}, dto.Category)
	assert.Contains(t, dto.Path, "_memo_JsonNote.md")
	assert.NotEmpty(t, dto.Date)
}

func TestCreateNonInteractive_JSON_RootCategoryIsArray(t *testing.T) {
	uc, _ := newTestMemoService(t)

	out := captureStdout(t, func() {
		require.NoError(t, uc.CreateNonInteractive("RootNote", nil, "", true))
	})
	// category must marshal to [] (never null) at the root.
	assert.Contains(t, out, `"category": []`)
}

func TestCreateNonInteractive_EmptyBody_NoBodyCreated(t *testing.T) {
	uc, baseDir := newTestMemoService(t)
	memosDir := filepath.Join(baseDir, "memos")

	require.NoError(t, captureStdoutErr(t, func() error {
		return uc.CreateNonInteractive("Bare", []string{}, "", false)
	}))

	fpath := findMemoFile(t, memosDir, "Bare")
	content, err := os.ReadFile(fpath)
	require.NoError(t, err)
	s := string(content)
	assert.Contains(t, s, "# Bare")
	// No body: the file ends right after the title heading.
	assert.Equal(t, "---\ncategory: []\n---\n\n# Bare\n", s)
}

// captureStdoutErr is a small wrapper that runs fn (which may error) while
// discarding its stdout, returning the error.
func captureStdoutErr(t *testing.T, fn func() error) error {
	t.Helper()
	var retErr error
	captureStdout(t, func() {
		retErr = fn()
	})
	return retErr
}
