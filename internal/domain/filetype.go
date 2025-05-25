package domain

import "github.com/hirotoni/memov2/internal/interfaces"

// FileType is an alias for interfaces.FileType to maintain backward compatibility
type FileType = interfaces.FileType

const (
	FileTypeTodos    = interfaces.FileTypeTodos
	FileTypeMemo     = interfaces.FileTypeMemo
	FileTypeWeekly   = interfaces.FileTypeWeekly
	FileTypeTemplate = interfaces.FileTypeTemplate
)
