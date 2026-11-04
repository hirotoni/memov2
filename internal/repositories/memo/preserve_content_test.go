package memo

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// richMemoBody is a memo whose content goes beyond what the parsed MemoFile
// model holds: extra frontmatter keys and a sub-heading with body text
// between the H1 and the first `##`.
const richMemoBody = `---
category: ["work"]
tags: [alpha, beta]
created: 2026-01-01
---

# Orig

intro line

### sub before h2

subtext

## Section

#### deep heading

deep text
`

func writeRichMemo(t *testing.T, dir string) domain.MemoFileInterface {
	t.Helper()
	date := time.Date(2026, 1, 2, 3, 4, 5, 0, time.Local)
	m, err := domain.NewMemoFile(date, "Orig", []string{"work"})
	require.NoError(t, err)
	path := filepath.Join(dir, m.Location(), m.FileName())
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(richMemoBody), 0o644))
	return m
}

func newQuietRepo(dir string) *memo {
	return &memo{dir: dir, logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(b)
}

func assertKeepsUserContent(t *testing.T, got string) {
	t.Helper()
	for _, want := range []string{
		"tags: [alpha, beta]\n",
		"created: 2026-01-01\n",
		"intro line\n",
		"### sub before h2\n\nsubtext\n",
		"#### deep heading\n\ndeep text\n",
	} {
		assert.Contains(t, got, want)
	}
}

func TestRename_PreservesContentExceptTitle(t *testing.T) {
	dir := t.TempDir()
	repo := newQuietRepo(dir)
	m := writeRichMemo(t, dir)
	oldPath := filepath.Join(dir, m.Location(), m.FileName())

	require.NoError(t, repo.Rename(m, "Renamed"))

	newPath := filepath.Join(dir, "work", strings.Replace(m.FileName(), "_Orig.md", "_Renamed.md", 1))
	assert.NoFileExists(t, oldPath)
	got := readFile(t, newPath)

	// The only difference from the original is the H1 title line.
	want := strings.Replace(richMemoBody, "# Orig\n", "# Renamed\n", 1)
	assert.Equal(t, want, got)
	assertKeepsUserContent(t, got)
}

func TestMove_PreservesContentExceptCategory(t *testing.T) {
	dir := t.TempDir()
	repo := newQuietRepo(dir)
	m := writeRichMemo(t, dir)
	oldPath := filepath.Join(dir, m.Location(), m.FileName())

	require.NoError(t, repo.Move(m, []string{"personal", "notes"}))

	newPath := filepath.Join(dir, "personal", "notes", m.FileName())
	assert.NoFileExists(t, oldPath)
	got := readFile(t, newPath)

	want := strings.Replace(richMemoBody, `category: ["work"]`, `category: ["personal", "notes"]`, 1)
	assert.Equal(t, want, got)
	assertKeepsUserContent(t, got)

	// The moved file must still parse with the new category.
	moved, err := repo.Memo(&movedRef{m, []string{"personal", "notes"}})
	require.NoError(t, err)
	assert.Equal(t, []string{"personal", "notes"}, moved.CategoryTree())
}

func TestDuplicate_PreservesContentExceptTitle(t *testing.T) {
	dir := t.TempDir()
	repo := newQuietRepo(dir)
	m := writeRichMemo(t, dir)
	origPath := filepath.Join(dir, m.Location(), m.FileName())

	dup, err := repo.Duplicate(m)
	require.NoError(t, err)

	assert.Equal(t, richMemoBody, readFile(t, origPath), "original must be untouched")
	assert.Equal(t, "Orig copied", dup.Title())
	assert.Equal(t, []string{"work"}, dup.CategoryTree())

	got := readFile(t, filepath.Join(dir, dup.Location(), dup.FileName()))
	want := strings.Replace(richMemoBody, "# Orig\n", "# Orig copied\n", 1)
	assert.Equal(t, want, got)
	assertKeepsUserContent(t, got)
}

func TestRename_RejectsInvalidTitleWithoutTouchingFile(t *testing.T) {
	dir := t.TempDir()
	repo := newQuietRepo(dir)
	m := writeRichMemo(t, dir)
	oldPath := filepath.Join(dir, m.Location(), m.FileName())

	assert.Error(t, repo.Rename(m, "../escape"))
	assert.Equal(t, richMemoBody, readFile(t, oldPath))
}

func TestRename_SameTitleRewritesInPlace(t *testing.T) {
	dir := t.TempDir()
	repo := newQuietRepo(dir)
	m := writeRichMemo(t, dir)
	path := filepath.Join(dir, m.Location(), m.FileName())

	require.NoError(t, repo.Rename(m, "Orig"))

	assert.Equal(t, richMemoBody, readFile(t, path))
	assertNoTempFiles(t, filepath.Dir(path))
}

func TestMove_SameLocationKeepsContent(t *testing.T) {
	dir := t.TempDir()
	repo := newQuietRepo(dir)
	m := writeRichMemo(t, dir)
	path := filepath.Join(dir, m.Location(), m.FileName())

	require.NoError(t, repo.Move(m, []string{"work"}))

	assert.Equal(t, richMemoBody, readFile(t, path))
	assertNoTempFiles(t, filepath.Dir(path))
}

func TestMove_InvalidCategoryLeavesFileUntouched(t *testing.T) {
	dir := t.TempDir()
	repo := newQuietRepo(dir)
	m := writeRichMemo(t, dir)
	path := filepath.Join(dir, m.Location(), m.FileName())

	assert.Error(t, repo.Move(m, []string{"..", "escape"}))
	assert.Equal(t, richMemoBody, readFile(t, path))
}

func TestRename_CaseOnlyChangeKeepsMemo(t *testing.T) {
	dir := t.TempDir()
	repo := newQuietRepo(dir)
	m := writeRichMemo(t, dir)
	oldName := m.FileName()

	require.NoError(t, repo.Rename(m, "orig"))

	// Works on both case-sensitive and case-insensitive filesystems: the memo
	// must survive under the new name with only the H1 changed.
	newName := strings.Replace(oldName, "_Orig.md", "_orig.md", 1)
	newPath := filepath.Join(dir, "work", newName)
	want := strings.Replace(richMemoBody, "# Orig\n", "# orig\n", 1)
	assert.Equal(t, want, readFile(t, newPath))

	// The on-disk name must carry the new casing (a path lookup alone would
	// also succeed via the old-cased entry on a case-insensitive filesystem).
	entries, err := os.ReadDir(filepath.Join(dir, "work"))
	require.NoError(t, err)
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	assert.Equal(t, []string{newName}, names)
}

func TestMove_CaseOnlyCategoryChangeKeepsMemo(t *testing.T) {
	dir := t.TempDir()
	repo := newQuietRepo(dir)
	m := writeRichMemo(t, dir)

	require.NoError(t, repo.Move(m, []string{"Work"}))

	newPath := filepath.Join(dir, "Work", m.FileName())
	want := strings.Replace(richMemoBody, `category: ["work"]`, `category: ["Work"]`, 1)
	assert.Equal(t, want, readFile(t, newPath))
}

func TestRename_KeepsFileModeAndWritesThroughSymlink(t *testing.T) {
	dir := t.TempDir()
	repo := newQuietRepo(dir)
	m := writeRichMemo(t, dir)
	path := filepath.Join(dir, m.Location(), m.FileName())
	require.NoError(t, os.Chmod(path, 0o600))

	// Move the real file elsewhere and leave a symlink at the memo path.
	realPath := filepath.Join(t.TempDir(), "real.md")
	require.NoError(t, os.Rename(path, realPath))
	require.NoError(t, os.Symlink(realPath, path))

	// Same-title rename rewrites in place: the symlink must remain and the
	// real file must keep its mode.
	require.NoError(t, repo.Rename(m, "Orig"))

	fi, err := os.Lstat(path)
	require.NoError(t, err)
	assert.NotZero(t, fi.Mode()&os.ModeSymlink, "symlink must be kept")
	real, err := os.Stat(realPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), real.Mode().Perm())
	assert.Equal(t, richMemoBody, readFile(t, realPath))
}

func assertNoTempFiles(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	for _, e := range entries {
		assert.False(t, strings.HasSuffix(e.Name(), ".tmp"), "temp file left behind: %s", e.Name())
	}
}

// movedRef points at a memo by its original filename under another category.
type movedRef struct {
	domain.MemoFileInterface
	tree []string
}

func (m *movedRef) CategoryTree() []string { return m.tree }
func (m *movedRef) Location() string       { return filepath.Join(m.tree...) }
