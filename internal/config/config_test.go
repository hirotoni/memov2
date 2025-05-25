package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTomlConfig(t *testing.T) {
	tests := []struct {
		name string
		opts TomlConfigOption
		want *TomlConfig
	}{
		{
			name: "default values when empty options",
			opts: TomlConfigOption{},
			want: func() *TomlConfig {
				c, err := NewDefaultConfig()
				require.NoError(t, err)
				return c
			}(),
		},
		{
			name: "custom base dir",
			opts: TomlConfigOption{BaseDir: "/tmp/test"},
			want: func() *TomlConfig {
				c, err := NewDefaultConfig()
				require.NoError(t, err)
				c.baseDir = "/tmp/test"
				return c
			}(),
		},
		{
			name: "custom todos folder",
			opts: TomlConfigOption{TodosFolderName: "custom_todos"},
			want: func() *TomlConfig {
				c, err := NewDefaultConfig()
				require.NoError(t, err)
				c.todosFolderName = "custom_todos"
				return c
			}(),
		},
		{
			name: "custom memos folder",
			opts: TomlConfigOption{MemosFolderName: "custom_memos"},
			want: func() *TomlConfig {
				c, err := NewDefaultConfig()
				require.NoError(t, err)
				c.memosFolderName = "custom_memos"
				return c
			}(),
		},
		{
			name: "custom days to seek",
			opts: TomlConfigOption{TodosDaysToSeek: 20},
			want: func() *TomlConfig {
				c, err := NewDefaultConfig()
				require.NoError(t, err)
				c.todosDaysToSeek = 20
				return c
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTomlConfig(tt.opts)
			require.NoError(t, err)
			if got.BaseDir() != tt.want.BaseDir() {
				t.Errorf("BaseDir = %v, want %v", got.BaseDir(), tt.want.BaseDir())
			}
			if got.TodosDir() != tt.want.TodosDir() {
				t.Errorf("TodosDir = %v, want %v", got.TodosDir(), tt.want.TodosDir())
			}
			if got.MemosDir() != tt.want.MemosDir() {
				t.Errorf("MemosDir = %v, want %v", got.MemosDir(), tt.want.MemosDir())
			}
			if got.TodosDaysToSeek() != tt.want.TodosDaysToSeek() {
				t.Errorf("TodosDaysToSeek = %v, want %v", got.TodosDaysToSeek(), tt.want.TodosDaysToSeek())
			}
		})
	}
}

func TestTomlConfig_Dirs(t *testing.T) {
	cfg := &TomlConfig{
		baseDir:         "/base",
		todosFolderName: "todos",
		memosFolderName: "memos",
	}

	if got := cfg.TodosDir(); got != "/base/todos" {
		t.Errorf("TodosDir() = %v, want %v", got, "/base/todos")
	}

	if got := cfg.MemosDir(); got != "/base/memos" {
		t.Errorf("MemosDir() = %v, want %v", got, "/base/memos")
	}
}
