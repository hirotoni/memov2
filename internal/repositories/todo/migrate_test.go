package todo

import (
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestTodoRepoImpl_RewriteHeadingText(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewTodo(tmpDir, logger)

	// A legacy daily todo file: its "wanttodos" heading should be rewritten while
	// its content and other headings are preserved byte-for-byte.
	daily := "# 20250101Wed\n\n## todos\n\n- [ ] task A\n\n## wanttodos\n\n- old idea\n"
	dailyWant := "# 20250101Wed\n\n## todos\n\n- [ ] task A\n\n## backlog\n\n- old idea\n"
	writeFile(t, filepath.Join(tmpDir, "20250101Wed_todos.md"), daily)

	// The template must also be migrated.
	tmpl := "# todos_template\n\n## todos\n\n## wanttodos\n"
	tmplWant := "# todos_template\n\n## todos\n\n## backlog\n"
	writeFile(t, filepath.Join(tmpDir, "todos_template.md"), tmpl)

	// A non-todo file, and an inline (non-heading) mention, must be left untouched.
	other := "# weekly\n\nthis mentions wanttodos inline and should stay\n"
	writeFile(t, filepath.Join(tmpDir, "weekly_report.md"), other)

	// Dry run: reports the affected files but changes nothing on disk.
	changed, err := repo.RewriteHeadingText("wanttodos", "backlog", true)
	if err != nil {
		t.Fatalf("dry run returned error: %v", err)
	}
	assertPaths(t, changed, []string{"20250101Wed_todos.md", "todos_template.md"})
	if got := readFile(t, filepath.Join(tmpDir, "20250101Wed_todos.md")); got != daily {
		t.Errorf("dry run modified daily file:\n%q", got)
	}

	// Apply: rewrites the headings, preserving surrounding content.
	changed, err = repo.RewriteHeadingText("wanttodos", "backlog", false)
	if err != nil {
		t.Fatalf("apply returned error: %v", err)
	}
	assertPaths(t, changed, []string{"20250101Wed_todos.md", "todos_template.md"})

	if got := readFile(t, filepath.Join(tmpDir, "20250101Wed_todos.md")); got != dailyWant {
		t.Errorf("daily file not migrated as expected:\n got %q\nwant %q", got, dailyWant)
	}
	if got := readFile(t, filepath.Join(tmpDir, "todos_template.md")); got != tmplWant {
		t.Errorf("template not migrated as expected:\n got %q\nwant %q", got, tmplWant)
	}
	if got := readFile(t, filepath.Join(tmpDir, "weekly_report.md")); got != other {
		t.Errorf("non-todo file was modified:\n%q", got)
	}

	// Idempotent: a second apply finds nothing to change.
	changed, err = repo.RewriteHeadingText("wanttodos", "backlog", false)
	if err != nil {
		t.Fatalf("second apply returned error: %v", err)
	}
	if len(changed) != 0 {
		t.Errorf("expected no changes on re-run, got %v", changed)
	}
}

// TestTodoRepoImpl_RewriteHeadingText_Level3NotMatched locks in that the migration
// regex is level-2 only: a "### wanttodos" (level-3) heading must be left untouched.
func TestTodoRepoImpl_RewriteHeadingText_Level3NotMatched(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewTodo(tmpDir, logger)

	content := "# 20250101Wed\n\n## todos\n\n- [ ] task A\n\n### wanttodos\n\n- old idea\n"
	writeFile(t, filepath.Join(tmpDir, "20250101Wed_todos.md"), content)

	changed, err := repo.RewriteHeadingText("wanttodos", "backlog", false)
	if err != nil {
		t.Fatalf("apply returned error: %v", err)
	}
	if len(changed) != 0 {
		t.Errorf("level-3 heading must not be rewritten, but reported changed: %v", changed)
	}
	if got := readFile(t, filepath.Join(tmpDir, "20250101Wed_todos.md")); got != content {
		t.Errorf("level-3 heading file was modified:\n got %q\nwant %q", got, content)
	}
}

// TestTodoRepoImpl_RewriteHeadingText_InlineMentionNotMatched locks in that an inline
// mention of "wanttodos" inside a todos file's body (not on a heading line) is not
// rewritten: only heading lines are targeted.
func TestTodoRepoImpl_RewriteHeadingText_InlineMentionNotMatched(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewTodo(tmpDir, logger)

	content := "# 20250101Wed\n\n## todos\n\n- remember to move wanttodos into backlog\n"
	writeFile(t, filepath.Join(tmpDir, "20250101Wed_todos.md"), content)

	changed, err := repo.RewriteHeadingText("wanttodos", "backlog", false)
	if err != nil {
		t.Fatalf("apply returned error: %v", err)
	}
	if len(changed) != 0 {
		t.Errorf("inline body mention must not be rewritten, but reported changed: %v", changed)
	}
	if got := readFile(t, filepath.Join(tmpDir, "20250101Wed_todos.md")); got != content {
		t.Errorf("inline body mention was modified:\n got %q\nwant %q", got, content)
	}
}

// TestTodoRepoImpl_RewriteHeadingText_NoTrailingNewline locks in that a todos file
// whose content has no trailing newline is still rewritten correctly, and no spurious
// trailing newline is introduced by the rewrite.
func TestTodoRepoImpl_RewriteHeadingText_NoTrailingNewline(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewTodo(tmpDir, logger)

	// Note: no trailing "\n" after the final heading.
	content := "# 20250101Wed\n\n## todos\n\n## wanttodos"
	want := "# 20250101Wed\n\n## todos\n\n## backlog"
	writeFile(t, filepath.Join(tmpDir, "20250101Wed_todos.md"), content)

	changed, err := repo.RewriteHeadingText("wanttodos", "backlog", false)
	if err != nil {
		t.Fatalf("apply returned error: %v", err)
	}
	assertPaths(t, changed, []string{"20250101Wed_todos.md"})
	if got := readFile(t, filepath.Join(tmpDir, "20250101Wed_todos.md")); got != want {
		t.Errorf("no-trailing-newline file not migrated as expected:\n got %q\nwant %q", got, want)
	}
}

// TestTodoRepoImpl_RewriteHeadingText_MultipleSpacesNormalized locks in that a heading
// written with multiple spaces ("##  wanttodos") is matched and normalized to the
// canonical single-space form ("## backlog").
func TestTodoRepoImpl_RewriteHeadingText_MultipleSpacesNormalized(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewTodo(tmpDir, logger)

	content := "# 20250101Wed\n\n##  wanttodos\n\n- old idea\n"
	// Expected: matched despite the extra space, and normalized to a single space.
	want := "# 20250101Wed\n\n## backlog\n\n- old idea\n"
	writeFile(t, filepath.Join(tmpDir, "20250101Wed_todos.md"), content)

	changed, err := repo.RewriteHeadingText("wanttodos", "backlog", false)
	if err != nil {
		t.Fatalf("apply returned error: %v", err)
	}
	assertPaths(t, changed, []string{"20250101Wed_todos.md"})
	if got := readFile(t, filepath.Join(tmpDir, "20250101Wed_todos.md")); got != want {
		t.Errorf("multi-space heading not normalized as expected:\n got %q\nwant %q", got, want)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	return string(b)
}

func assertPaths(t *testing.T, got, want []string) {
	t.Helper()
	g := append([]string(nil), got...)
	w := append([]string(nil), want...)
	sort.Strings(g)
	sort.Strings(w)
	if len(g) != len(w) {
		t.Fatalf("changed paths mismatch:\n got %v\nwant %v", g, w)
	}
	for i := range g {
		if g[i] != w[i] {
			t.Fatalf("changed paths mismatch:\n got %v\nwant %v", g, w)
		}
	}
}
