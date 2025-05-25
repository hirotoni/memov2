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
				c.BaseDir = "/tmp/test"
				return c
			}(),
		},
		{
			name: "custom todos folder",
			opts: TomlConfigOption{TodosFolderName: "custom_todos"},
			want: func() *TomlConfig {
				c, err := NewDefaultConfig()
				require.NoError(t, err)
				c.TodosFolderName = "custom_todos"
				return c
			}(),
		},
		{
			name: "custom memos folder",
			opts: TomlConfigOption{MemosFolderName: "custom_memos"},
			want: func() *TomlConfig {
				c, err := NewDefaultConfig()
				require.NoError(t, err)
				c.MemosFolderName = "custom_memos"
				return c
			}(),
		},
		{
			name: "custom days to seek",
			opts: TomlConfigOption{TodosDaysToSeek: 20},
			want: func() *TomlConfig {
				c, err := NewDefaultConfig()
				require.NoError(t, err)
				c.TodosDaysToSeek = 20
				return c
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTomlConfig(tt.opts)
			require.NoError(t, err)
			if got.BaseDir != tt.want.BaseDir {
				t.Errorf("BaseDir = %v, want %v", got.BaseDir, tt.want.BaseDir)
			}
			if got.TodosFolderName != tt.want.TodosFolderName {
				t.Errorf("TodosFolderName = %v, want %v", got.TodosFolderName, tt.want.TodosFolderName)
			}
			if got.MemosFolderName != tt.want.MemosFolderName {
				t.Errorf("MemosFolderName = %v, want %v", got.MemosFolderName, tt.want.MemosFolderName)
			}
			if got.TodosDaysToSeek != tt.want.TodosDaysToSeek {
				t.Errorf("TodosDaysToSeek = %v, want %v", got.TodosDaysToSeek, tt.want.TodosDaysToSeek)
			}
		})
	}
}

func TestTomlConfig_Dirs(t *testing.T) {
	cfg := &TomlConfig{
		BaseDir:         "/base",
		TodosFolderName: "todos",
		MemosFolderName: "memos",
	}

	if got := cfg.TodosDir(); got != "/base/todos" {
		t.Errorf("TodosDir() = %v, want %v", got, "/base/todos")
	}

	if got := cfg.MemosDir(); got != "/base/memos" {
		t.Errorf("MemosDir() = %v, want %v", got, "/base/memos")
	}
}
