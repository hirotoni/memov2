package domain

import (
	"fmt"
	"testing"
	"time"
)

func Test_aaaaaa(t *testing.T) {
	a, err := NewMemoFile(
		time.Now(),
		"Test Memo",
		[]string{"Category1", "Category2"},
	)
	if err != nil {
		t.Fatalf("failed to create memo file: %v", err)
	}

	fmt.Println("a.FileName():", a.FileName())
	fmt.Println("a.Date():", a.Date())
	fmt.Println("a.FileType():", a.FileType())
	fmt.Println("a.Title():", a.Title())
	fmt.Println("a.CategoryTree():", a.CategoryTree())
}
