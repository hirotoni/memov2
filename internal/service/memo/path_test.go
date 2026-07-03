package memo

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveToMemosDir(t *testing.T) {
	memosDir := filepath.Join(string(filepath.Separator)+"a", "memos")

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name: "accept relative list-output form",
			path: filepath.Join("category", "file.md"),
			want: filepath.Join(memosDir, "category", "file.md"),
		},
		{
			name: "accept absolute path under memosDir",
			path: filepath.Join(memosDir, "category", "file.md"),
			want: filepath.Join(memosDir, "category", "file.md"),
		},
		{
			name: "accept memosDir root itself",
			path: memosDir,
			want: filepath.Clean(memosDir),
		},
		{
			name: "accept empty path as root",
			path: "",
			want: filepath.Clean(memosDir),
		},
		{
			name: "accept dot path as root",
			path: ".",
			want: filepath.Clean(memosDir),
		},
		{
			name:    "reject relative traversal",
			path:    filepath.Join("..", "..", "etc", "passwd"),
			wantErr: true,
		},
		{
			name:    "reject absolute path outside memosDir",
			path:    filepath.Join(string(filepath.Separator)+"etc", "passwd"),
			wantErr: true,
		},
		{
			name:    "reject escape after normalization",
			path:    filepath.Join("category", "..", "..", ".."),
			wantErr: true,
		},
		{
			name:    "reject sibling dir with shared prefix",
			path:    filepath.Join(string(filepath.Separator)+"a", "memos-evil", "x"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveToMemosDir(memosDir, tt.path)
			if tt.wantErr {
				require.Error(t, err)
				var appErr *common.AppError
				require.True(t, errors.As(err, &appErr), "expected *common.AppError, got %T", err)
				assert.Equal(t, common.ErrorTypeValidation, appErr.Type)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestResolveToMemosDir_CWDRelativeBranch exercises case 2: a bare relative path
// that resolves via the current working directory to a location under memosDir.
// It requires a real memosDir on disk so filepath.Abs (which joins with CWD) yields
// a path under it. t.Chdir (Go 1.24+) sets and restores the working directory.
func TestResolveToMemosDir_CWDRelativeBranch(t *testing.T) {
	// Canonicalize so this matches os.Getwd() on platforms where TempDir is a
	// symlink (e.g. macOS /var → /private/var), ensuring the CWD-relative branch
	// (filepath.Abs under memosDir) is exercised rather than the join fallback.
	memosDir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	t.Chdir(memosDir)

	got, err := resolveToMemosDir(memosDir, "file.md")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(filepath.Clean(memosDir), "file.md"), got)
}
