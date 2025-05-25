package weekly

import (
	"bufio"
	"log/slog"
	"path/filepath"

	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/platform"
)

type weekly struct {
	dir    string
	logger *slog.Logger
}

func NewWeekly(dir string, logger *slog.Logger) interfaces.WeeklyRepo {
	return &weekly{dir: dir, logger: logger}
}

func (r *weekly) Save(file interfaces.WeeklyFileInterface, truncate bool) error {
	path := filepath.Join(r.dir, file.FileName())
	if !platform.Exists(path) || truncate {
		if err := platform.WriteFileStream(path, truncate, func(w *bufio.Writer) error {
			_, err := w.WriteString(file.ContentString())
			return err
		}); err != nil {
			return err
		}
		r.logger.Info("File saved", "path", path)
	}
	return nil
}
