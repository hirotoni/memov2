package repository

import (
	"bufio"
	"fmt"
	"path/filepath"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/platform/fs"
)

type Weekly interface {
	Save(file domain.WeeklyFileInterface, truncate bool) error
}

type weekly struct {
	dir string
}

func NewWeekly(dir string) Weekly {
	return &weekly{dir: dir}
}

func (r *weekly) Save(file domain.WeeklyFileInterface, truncate bool) error {
	path := filepath.Join(r.dir, file.FileName())
	if !fs.Exists(path) || truncate {
		if err := fs.WriteFileStream(path, truncate, func(w *bufio.Writer) error {
			_, err := w.Write(file.Data())
			return err
		}); err != nil {
			return err
		}
		fmt.Printf("File saved: %s\n", path)
	}
	return nil
}
