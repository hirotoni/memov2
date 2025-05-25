package domain

type FileType string

const (
	FileTypeTodos    FileType = "todos"
	FileTypeMemo     FileType = "memo"
	FileTypeWeekly   FileType = "weekly"
	FileTypeTemplate FileType = "template"
)

func (ft FileType) String() string { return string(ft) }
