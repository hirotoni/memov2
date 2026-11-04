package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupHome points HOME at a fresh temp directory and writes a config.toml
// whose editor is the no-op `true`, so commands that open an editor on
// success return immediately and never touch the user's real data.
// It returns the memos directory.
func setupHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)

	base := filepath.Join(home, "base")
	memosDir := filepath.Join(base, "memos")
	require.NoError(t, os.MkdirAll(memosDir, 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(base, "todos"), 0o755))

	cfgDir := filepath.Join(home, ".config", "memov2")
	require.NoError(t, os.MkdirAll(cfgDir, 0o755))
	cfg := "base_dir = \"" + base + "\"\n" +
		"memos_foldername = \"memos/\"\n" +
		"todos_foldername = \"todos/\"\n" +
		"todos_daystoseek = 10\n" +
		"editor = \"true\"\n" +
		"editor_args = [\"{path}\"]\n"
	require.NoError(t, os.WriteFile(filepath.Join(cfgDir, "config.toml"), []byte(cfg), 0o644))
	return memosDir
}

// runRoot executes RootCmd with args and returns what was written to the
// process stdout (services write to os.Stdout directly) and the error that
// Execute returned; a non-nil error makes main exit with status 1.
// These tests share the package-global RootCmd and swap os.Stdout, so they
// must not call t.Parallel.
func runRoot(t *testing.T, args ...string) (string, error) {
	t.Helper()

	// Flag variables are package globals that survive across Execute calls;
	// reset the target command's flags so earlier runs do not leak in.
	c, _, err := RootCmd.Find(args)
	require.NoError(t, err)
	c.Flags().VisitAll(func(f *pflag.Flag) {
		require.NoError(t, f.Value.Set(f.DefValue))
		f.Changed = false
	})

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	var errBuf bytes.Buffer
	RootCmd.SetOut(w)
	RootCmd.SetErr(&errBuf)
	RootCmd.SetArgs(args)
	defer func() {
		RootCmd.SetOut(nil)
		RootCmd.SetErr(nil)
		RootCmd.SetArgs(nil)
	}()
	execErr := RootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout
	return <-done, execErr
}

// makeUnreadable chmods dir to 000 and restores it on cleanup so t.TempDir
// can remove it.
func makeUnreadable(t *testing.T, dir string) {
	t.Helper()
	if os.Geteuid() == 0 {
		t.Skip("permission checks do not apply to root")
	}
	require.NoError(t, os.Chmod(dir, 0o000))
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
}

func TestMemosOpen_PathEscape_FailsWithEmptyStdout(t *testing.T) {
	setupHome(t)
	stdout, err := runRoot(t, "memos", "open", "../../etc/passwd")
	require.Error(t, err)
	assert.Empty(t, stdout)
}

func TestCommands_InternalError_ReturnsError(t *testing.T) {
	tests := []struct {
		args        []string
		emptyStdout bool
	}{
		{[]string{"memos", "categories"}, true},
		{[]string{"memos", "categories", "--json"}, true},
		{[]string{"memos", "index"}, true},
		// memos weekly prints a progress line to stdout before failing; only
		// the exit status is asserted here.
		{[]string{"memos", "weekly"}, false},
	}
	for _, tt := range tests {
		args := tt.args
		t.Run(filepath.Join(args...), func(t *testing.T) {
			memosDir := setupHome(t)
			// Give the memos directory a category so Categories has to read it.
			require.NoError(t, os.MkdirAll(filepath.Join(memosDir, "work"), 0o755))
			makeUnreadable(t, memosDir)

			stdout, err := runRoot(t, args...)
			require.Error(t, err)
			if tt.emptyStdout {
				assert.Empty(t, stdout)
			}
		})
	}
}

func TestConfigShow_BrokenConfig_ReturnsError(t *testing.T) {
	setupHome(t)
	home := os.Getenv("HOME")
	cfgPath := filepath.Join(home, ".config", "memov2", "config.toml")
	require.NoError(t, os.WriteFile(cfgPath, []byte("base_dir = [unterminated"), 0o644))

	for _, args := range [][]string{{"config", "show"}, {"config", "show", "--json"}} {
		stdout, err := runRoot(t, args...)
		require.Error(t, err)
		assert.Empty(t, stdout)
	}
}

func TestCommands_Success_ReturnsNilAndKeepsStdout(t *testing.T) {
	memosDir := setupHome(t)
	require.NoError(t, os.MkdirAll(filepath.Join(memosDir, "work", "daily"), 0o755))

	stdout, err := runRoot(t, "memos", "categories")
	require.NoError(t, err)
	assert.Equal(t, "work\nwork/daily\n", stdout)

	stdout, err = runRoot(t, "memos", "categories", "--json")
	require.NoError(t, err)
	assert.JSONEq(t, `[{"category":["work"],"path":"work"},{"category":["work","daily"],"path":"work/daily"}]`, stdout)

	stdout, err = runRoot(t, "config", "show", "--json")
	require.NoError(t, err)
	assert.Contains(t, stdout, `"memos_dir": "`+memosDir)

	for _, tt := range []struct {
		args   []string
		stdout string
	}{
		{[]string{"memos", "index"}, ""},
		{[]string{"memos", "weekly"}, "Building weekly report...\n"},
		{[]string{"config", "edit"}, ""},
	} {
		stdout, err = runRoot(t, tt.args...)
		require.NoError(t, err, "%v", tt.args)
		assert.Equal(t, tt.stdout, stdout, "%v", tt.args)
	}
}
