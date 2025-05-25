package toml

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		opts Option
		want *Config
	}{
		{
			name: "default values when empty options",
			opts: Option{},
			want: func() *Config {
				c, err := NewDefaultConfig()
				require.NoError(t, err)
				return c
			}(),
		},
		{
			name: "custom base dir",
			opts: Option{BaseDir: "/tmp/test"},
			want: func() *Config {
				c, err := NewConfig(Option{BaseDir: "/tmp/test"})
				require.NoError(t, err)
				return c
			}(),
		},
		{
			name: "custom todos folder",
			opts: Option{TodosFolderName: "custom_todos"},
			want: func() *Config {
				c, err := NewConfig(Option{TodosFolderName: "custom_todos"})
				require.NoError(t, err)
				return c
			}(),
		},
		{
			name: "custom memos folder",
			opts: Option{MemosFolderName: "custom_memos"},
			want: func() *Config {
				c, err := NewConfig(Option{MemosFolderName: "custom_memos"})
				require.NoError(t, err)
				return c
			}(),
		},
		{
			name: "custom days to seek",
			opts: Option{TodosDaysToSeek: 20},
			want: func() *Config {
				c, err := NewConfig(Option{TodosDaysToSeek: 20})
				require.NoError(t, err)
				return c
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.opts)
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

func TestConfig_Dirs(t *testing.T) {
	cfg, err := NewConfig(Option{BaseDir: "/base", TodosFolderName: "todos", MemosFolderName: "memos"})
	require.NoError(t, err)

	if got := cfg.TodosDir(); got != "/base/todos" {
		t.Errorf("TodosDir() = %v, want %v", got, "/base/todos")
	}

	if got := cfg.MemosDir(); got != "/base/memos" {
		t.Errorf("MemosDir() = %v, want %v", got, "/base/memos")
	}
}

