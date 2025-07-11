package models

import "time"

type WeeklyFileInterface interface {
	fileInterface
}

type WeeklyFile struct {
	file
}

func NewWeeklyFile() (WeeklyFileInterface, error) {
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
	return filename + Ext
}
