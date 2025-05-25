package domain

import (
	"time"

	"github.com/hirotoni/memov2/internal/interfaces"
)

// WeeklyFileInterface is an alias for interfaces.WeeklyFileInterface to maintain backward compatibility
type WeeklyFileInterface = interfaces.WeeklyFileInterface
type WeeklyFile struct{ file }

func NewWeekly() (WeeklyFileInterface, error) {
	// set the current date but wont use it in the filename
	date := time.Now()

	f := &WeeklyFile{
		file: file{
			date:     date,
			fileType: FileTypeWeekly,
			title:    "weekly_report",
		},
	}
	return f, nil
}

func (f *WeeklyFile) FileName() string {
	filename := "weekly_report"
	return filename + FileExtension
}

// ContentString overrides the base implementation to ensure trailing newline
func (f *WeeklyFile) ContentString() string {
	content := f.file.ContentString()
	// Ensure content ends with \n\n (two newlines) for weekly files
	if len(f.HeadingBlocks()) > 0 {
		// If there are heading blocks, base ContentString() already adds one \n
		// We need to add one more to match golden files
		return content + "\n"
	}
	// If no heading blocks, base ContentString() ends with \n\n from title
	// But we still want \n\n at the end
	if content != "" && len(content) >= 2 && content[len(content)-2:] != "\n\n" {
		return content + "\n"
	}
	return content
}
