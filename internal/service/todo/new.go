package todo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/utils"
)

func (uc todo) GenerateTodoFile(truncate bool) error {
	now := time.Now()
	repo := uc.r.Todo()

	md, err := uc.inheritTodos(now, uc.c.TodosDaysToSeek())
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error inheriting todos")
	}

	err = repo.Save(md, truncate)
	if err != nil {
		return err
	}

	fpath := filepath.Join(uc.c.TodosDir(), md.FileName())
	err = uc.e.Open(uc.c.BaseDir(), fpath)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error opening editor")
	}

	return nil
}

// inheritTodos inherits information of the specified heading from previous day's memo
func (uc todo) inheritTodos(today time.Time, daysToSeek int) (domain.TodoFileInterface, error) {
	repo := uc.r.Todo()

	// templateファイルから雛形生成
	f, err := repo.TodosTemplate(time.Now())
	if err != nil {
		return nil, common.Wrap(err, common.ErrorTypeService, "failed to load todos template")
	}

	// 過去のファイルからtodosを継承
	found, err := uc.findPrevTodosFile(today, daysToSeek)
	if err != nil {
		return nil, common.Wrap(err, common.ErrorTypeService, "failed to find previous todos file")
	}

	if found != nil {
		// ファイルに必要な情報を設定
		for _, entity := range found.HeadingBlocks() {
			switch entity.HeadingText {
			// todos, wanttodos のものだけ継承する
			case utils.HeadingTodos.Text, utils.HeadingWantTodos.Text:
				f.OverrideHeadingBlockMatched(entity)
			}
		}
	}

	return f, nil
}

func (uc todo) findPrevTodosFile(today time.Time, daysToSeek int) (domain.TodoFileInterface, error) {
	repo := uc.r.Todo()

	var found domain.TodoFileInterface
	for i := range daysToSeek {
		prevDay := today.AddDate(0, 0, -1*(i+1))
		md, err := repo.FindTodosFileByDate(prevDay)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if i+1 == daysToSeek {
					return nil, common.New(common.ErrorTypeValidation, fmt.Sprintf("previous todos were not found in previous %d days", daysToSeek))
				}
				continue
			}
			return nil, common.Wrap(err, common.ErrorTypeService, "error finding previous todos file")
		}

		if md != nil {
			found = md
			break
		}
	}

	if found == nil {
		return nil, nil
	}

	return found, nil
}
