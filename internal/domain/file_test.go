package domain

import (
	"testing"
	"time"
)

func TestNewFile_ValidInput(t *testing.T) {
	date := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	fileType := FileTypeTodos

	f, err := NewTodosFile(date)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if f.Date() != date {
		t.Errorf("expected date %v, got %v", date, f.Date())
	}
	if f.FileType() != string(fileType) {
		t.Errorf("expected fileType %v, got %v", fileType, f.FileType())
	}
}

func TestNewFile_InvalidDate(t *testing.T) {
	title := "invalid-date"
	_, err := NewMemoFile(time.Time{}, title, nil)
	if err == nil {
		t.Fatal("expected error for zero date, got nil")
	}
}

func TestFileType_String(t *testing.T) {
	if FileTypeTodos.String() != "todos" {
		t.Errorf("expected 'todos', got %v", FileTypeTodos.String())
	}
	if FileTypeMemo.String() != "memo" {
		t.Errorf("expected 'memo', got %v", FileTypeMemo.String())
	}
}
func TestFileName_Weekly(t *testing.T) {
	f, _ := NewWeekly()
	expected := "weekly_report.md"

	if f.FileName() != expected {
		t.Errorf("expected weekly filename %v, got %v", expected, f.FileName())
	}
}

func TestFileName_Template(t *testing.T) {
	f, _ := NewTodoTemplateFile()
	expected := "todos_template.md"

	if f.FileName() != expected {
		t.Errorf("expected template filename %v, got %v", expected, f.FileName())
	}
}
