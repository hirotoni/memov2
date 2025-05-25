package todos

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hirotoni/memov2/components"
	"github.com/hirotoni/memov2/config"
	"github.com/hirotoni/memov2/models"
	"github.com/hirotoni/memov2/repos"
	"github.com/hirotoni/memov2/utils"
)

func GenerateTodoFile(c config.TomlConfig, truncate bool) error {
	if !utils.Exists(c.TodosDir()) {
		err := os.MkdirAll(c.TodosDir(), 0755)
		if err != nil {
			log.Printf("Error creating todos directory: %v\n", err)
			return err
		}
		log.Println("Created todos directory:", c.TodosDir())
	}

	now := time.Now()
	repo := repos.NewTodoFileRepo(c.TodosDir())

	md, err := inheritTodos(c.TodosDir(), now, c.TodosDaysToSeek)
	if err != nil {
		log.Println("No previous todos found, creating a new file.")
		return err
	}

	err = repo.Save(md, truncate)
	if err != nil {
		return err
	}

	fpath := filepath.Join(c.TodosDir(), md.FileName())
	components.OpenEditor(c.BaseDir, fpath)

	return nil
}

// inheritTodos inherits information of the specified heading from previous day's memo
func inheritTodos(dir string, today time.Time, daysToSeek int) (models.TodoFileInterface, error) {
	repo := repos.NewTodoFileRepo(dir)

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

func findPrevTodosFile(baseDir string, today time.Time, daysToSeek int) (models.TodoFileInterface, error) {
	repo := repos.NewTodoFileRepo(baseDir)

	var found models.TodoFileInterface
	for i := range daysToSeek {
		prevDay := today.AddDate(0, 0, -1*(i+1))
		md, err := repo.FindTodosFileByDate(prevDay)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if i+1 == daysToSeek {
					log.Printf("previous todos were not found in previous %d days.", daysToSeek)
				}
				continue
			}
			log.Fatal(err)
		}

		if md != nil {
			log.Println("Found previous todos for", prevDay.Format(models.FileNameDateLayoutTodo))
			found = md
			break
		}
	}

	if found == nil {
		log.Printf("No todos found in the previous %d days.", daysToSeek)
		return nil, nil
	}

	return found, nil
}
