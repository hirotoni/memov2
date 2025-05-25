package todo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/platform/fs"
	todoRepo "github.com/hirotoni/memov2/internal/repositories/todo"
	"github.com/hirotoni/memov2/utils"
)

func (uc todo) GenerateTodoFile(truncate bool) error {
	if err := fs.EnsureDir(uc.c.TodosDir()); err != nil {
		return fmt.Errorf("error ensuring todos directory: %v", err)
	}

	now := time.Now()
	repo := todoRepo.NewTodo(uc.c.TodosDir())

	md, err := inheritTodos(uc.c.TodosDir(), now, uc.c.TodosDaysToSeek())
	if err != nil {
		return fmt.Errorf("error inheriting todos: %v", err)
	}

	err = repo.Save(md, truncate)
	if err != nil {
		return err
	}

	fpath := filepath.Join(uc.c.TodosDir(), md.FileName())
	err = uc.e.Open(uc.c.BaseDir(), fpath)
	if err != nil {
		return fmt.Errorf("error opening editor: %v", err)
	}

	return nil
}

// inheritTodos inherits information of the specified heading from previous day's memo
func inheritTodos(dir string, today time.Time, daysToSeek int) (domain.TodoFileInterface, error) {
	repo := todoRepo.NewTodo(dir)

	// templateファイルから雛形生成
	f, err := repo.TodosTemplate(time.Now())
	if err != nil {
		return nil, errors.New("failed to load todos template")
	}

	// 過去のファイルからtodosを継承
	found, err := findPrevTodosFile(dir, today, daysToSeek)
	if err != nil {
		return nil, errors.New("failed to find previous todos file")
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

func findPrevTodosFile(baseDir string, today time.Time, daysToSeek int) (domain.TodoFileInterface, error) {
	repo := todoRepo.NewTodo(baseDir)

	var found domain.TodoFileInterface
	for i := range daysToSeek {
		prevDay := today.AddDate(0, 0, -1*(i+1))
		md, err := repo.FindTodosFileByDate(prevDay)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if i+1 == daysToSeek {
					return nil, fmt.Errorf("previous todos were not found in previous %d days", daysToSeek)
				}
				continue
			}
			return nil, fmt.Errorf("error finding previous todos file: %v", err)
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
