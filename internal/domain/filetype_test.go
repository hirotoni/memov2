package domain

import "testing"

func TestFileType_String(t *testing.T) {
	if FileTypeTodos.String() != "todos" {
		t.Errorf("expected 'todos', got %v", FileTypeTodos.String())
	}
	if FileTypeMemo.String() != "memo" {
		t.Errorf("expected 'memo', got %v", FileTypeMemo.String())
	}
}
