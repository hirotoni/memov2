package platform

import (
	"testing"

	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestNewEditor(t *testing.T) {
	editor := NewEditor("vi", []string{"{path}"})

	assert.NotNil(t, editor)
	assert.Implements(t, (*interfaces.Editor)(nil), editor)
}

func TestDefaultEditor_BuildArgs(t *testing.T) {
	tests := []struct {
		name     string
		template []string
		basedir  string
		path     string
		want     []string
	}{
		{
			name:     "vi style",
			template: []string{"{path}"},
			basedir:  "/tmp",
			path:     "/tmp/test.txt",
			want:     []string{"/tmp/test.txt"},
		},
		{
			name:     "vscode style",
			template: []string{"--folder-uri", "{basedir}", "--goto", "{path}:7"},
			basedir:  "/tmp",
			path:     "/tmp/test.txt",
			want:     []string{"--folder-uri", "/tmp", "--goto", "/tmp/test.txt:7"},
		},
		{
			name:     "neovim with neo-tree",
			template: []string{"-c", "Neotree reveal", "{path}"},
			basedir:  "/tmp",
			path:     "/tmp/test.txt",
			want:     []string{"-c", "Neotree reveal", "/tmp/test.txt"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := DefaultEditor{command: "editor", argsTemplate: tt.template}
			got := e.buildArgs(tt.basedir, tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDefaultEditor_Type(t *testing.T) {
	editor := DefaultEditor{}
	var _ interfaces.Editor = editor
}

