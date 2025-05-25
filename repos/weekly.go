package repos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hirotoni/memov2/models"
	"github.com/hirotoni/memov2/utils"
)

type WeeklyFileInterface interface {
	Save(file models.WeeklyFileInterface, truncate bool) error
}

type weeklyFileRepoImpl struct {
	dir string
}

func NewWeeklyFileRepo(dir string) WeeklyFileInterface {
	return &weeklyFileRepoImpl{
		dir: dir,
	}
}

func (r *weeklyFileRepoImpl) Save(file models.WeeklyFileInterface, truncate bool) error {
	path := filepath.Join(r.dir, file.FileName())
	if !utils.Exists(path) || truncate {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		f.WriteString("# " + file.Title() + "\n\n")

		for _, v := range file.HeadingBlocks() {
			f.WriteString(v.String())
		}
		fmt.Printf("File saved: %s\n", path)
	}
	return nil
}
