package domain

import (
	"reflect"
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

	if wf.FileName() != "weekly_report.md" {
		t.Errorf("expected weekly filename %v, got %v", "weekly_report.md", wf.FileName())
	}

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

func TestNewWeekly(t *testing.T) {
	tests := []struct {
		name    string
		want    WeeklyFileInterface
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWeekly()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWeekly() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWeekly() = %v, want %v", got, tt.want)
			}
		})
	}
}
