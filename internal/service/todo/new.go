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

// generateTodoFile builds today's todo file (inheriting from previous days) and
// persists it, WITHOUT launching the editor. It returns the domain file so
// callers can decide what to do next (open the editor, or report its path). This
// is the shared core of the interactive GenerateTodoFile and the non-interactive
// GenerateTodoFileNonInteractive. Existing "today's file" semantics are preserved:
// Save is a no-op when the file already exists unless truncate is true.
func (uc todo) generateTodoFile(truncate bool) (domain.TodoFileInterface, error) {
	now := time.Now()
	repo := uc.r.Todo()

	md, err := uc.inheritTodos(now, uc.c.TodosDaysToSeek())
	if err != nil {
		return nil, common.Wrap(err, common.ErrorTypeService, "error inheriting todos")
	}

	if err := repo.Save(md, truncate); err != nil {
		return nil, err
	}

	return md, nil
}

func (uc todo) GenerateTodoFile(truncate bool) error {
	md, err := uc.generateTodoFile(truncate)
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

// GenerateTodoFileNonInteractive creates today's todo file without launching the
// editor. It is the write counterpart of the read-only `--json` commands and backs
// `todos create --no-open` / `todos create --json`. It keeps the existing "today's file"
// semantics: if the file already exists it is left intact unless truncate is true,
// and the existing file's path is still reported (this is not a collision error).
//
// On success it writes the todo file's TodosDir-relative path (one line) to
// stdout, or a todoNewResultDTO when asJSON is true. Logs and errors go to stderr.
func (uc todo) GenerateTodoFileNonInteractive(truncate, asJSON bool) error {
	md, err := uc.generateTodoFile(truncate)
	if err != nil {
		return err
	}

	relPath := md.FileName()

	if asJSON {
		return writeJSON(todoNewResultDTO{
			Date: md.Date(),
			Path: relPath,
		})
	}

	fmt.Fprintln(os.Stdout, relPath)
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
			// todos, backlog のものだけ継承する
			case utils.HeadingTodos.Text, utils.HeadingBacklog.Text:
				f.OverrideHeadingBlockMatched(entity)
			// rename 前の legacy backlog 見出し ("wanttodos") は backlog に正規化してから継承する
			// (未 migrate のまま todos create された場合の carryover 消失を防ぐ)
			case utils.HeadingBacklogLegacy.Text:
				entity.HeadingText = utils.HeadingBacklog.Text
				f.OverrideHeadingBlockMatched(entity)
			}
		}
	}

	return f, nil
}

// findPrevTodosFile returns the most recent daily todo file within the previous
// daysToSeek days. When none exists in that window (first run, or after a long
// break) it returns (nil, nil) so today's file is built from the template alone.
func (uc todo) findPrevTodosFile(today time.Time, daysToSeek int) (domain.TodoFileInterface, error) {
	repo := uc.r.Todo()

	for i := range daysToSeek {
		prevDay := today.AddDate(0, 0, -1*(i+1))
		md, err := repo.FindTodosFileByDate(prevDay)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, common.Wrap(err, common.ErrorTypeService, "error finding previous todos file")
		}

		if md != nil {
			return md, nil
		}
	}

	return nil, nil
}
