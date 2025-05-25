package memo

import (
	"path/filepath"
	"time"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/platform"
)

func (uc memo) GenerateMemoFile(title string) error {
	if title == "" {
		var err error
		title, err = platform.ReadLine("Title: ")
		if err != nil {
			return err
		}
	}

	// Create a new memo file with the given title
	today := time.Now()
	memoFile, err := domain.NewMemoFile(today, title, []string{})
	if err != nil {
		return err
	}

	// Save the memo file to the base directory
	err = uc.repos.Memo().Save(memoFile, false)
	if err != nil {
		return err
	}

	fpath := filepath.Join(uc.config.MemosDir(), memoFile.Location(), memoFile.FileName())

	err = uc.editor.Open(uc.config.BaseDir(), fpath)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error opening editor")
	}

	return nil
}
