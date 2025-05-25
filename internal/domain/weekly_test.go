package domain

import (
	"testing"
	"time"
)

func TestNewWeeklyFile(t *testing.T) {
	weeklyFile, err := NewWeekly()

	if err != nil {
		t.Fatalf("NewWeeklyFile() returned an error: %v", err)
	}

	if weeklyFile == nil {
		t.Fatal("NewWeeklyFile() returned nil")
	}

	// Test the underlying file properties
	wf := weeklyFile.(*WeeklyFile)

	if wf.fileType != FileTypeWeekly {
		t.Errorf("Expected fileType to be FileTypeWeekly, got %v", wf.fileType)
	}

	if wf.title != "weekly_report" {
		t.Errorf("Expected title to be 'weekly_report', got %s", wf.title)
	}

	// Check that date is set to current time (within reasonable bounds)
	now := time.Now()
	timeDiff := now.Sub(wf.date)
	if timeDiff < 0 || timeDiff > time.Second {
		t.Errorf("Expected date to be close to current time, got %v", wf.date)
	}
}

func TestWeeklyFile_FileName(t *testing.T) {
	weeklyFile, err := NewWeekly()
	if err != nil {
		t.Fatalf("NewWeeklyFile() returned an error: %v", err)
	}

	fileName := weeklyFile.FileName()
	expected := "weekly_report" + FileExtension

	if fileName != expected {
		t.Errorf("FileName() = %s, want %s", fileName, expected)
	}
}

func TestWeeklyFileInterface(t *testing.T) {
	// Test that WeeklyFile implements both WeeklyFileInterface and FileInterface
	var _ WeeklyFileInterface = &WeeklyFile{}
	var _ FileInterface = &WeeklyFile{}
}
