package domain

import "testing"

func TestValidateCategoryTree(t *testing.T) {
	testCases := []struct {
		name    string
		tree    []string
		wantErr bool
	}{
		{
			name:    "nil tree is valid",
			tree:    nil,
			wantErr: false,
		},
		{
			name:    "empty tree is valid",
			tree:    []string{},
			wantErr: false,
		},
		{
			name:    "normal segments are valid",
			tree:    []string{"work", "projects"},
			wantErr: false,
		},
		{
			name:    "empty segment is invalid",
			tree:    []string{"work", ""},
			wantErr: true,
		},
		{
			name:    "dot segment is invalid",
			tree:    []string{"."},
			wantErr: true,
		},
		{
			name:    "dot-dot segment is invalid",
			tree:    []string{".."},
			wantErr: true,
		},
		{
			name:    "dot-dot nested segment is invalid",
			tree:    []string{"work", ".."},
			wantErr: true,
		},
		{
			name:    "segment containing forward slash is invalid",
			tree:    []string{"work/projects"},
			wantErr: true,
		},
		{
			name:    "segment containing backslash is invalid",
			tree:    []string{"work\\projects"},
			wantErr: true,
		},
		{
			name:    "absolute path segment is invalid",
			tree:    []string{"/etc"},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCategoryTree(tc.tree)
			if tc.wantErr && err == nil {
				t.Errorf("expected an error for tree %v, got nil", tc.tree)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("expected no error for tree %v, got %v", tc.tree, err)
			}
		})
	}
}
