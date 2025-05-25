package repos

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/hirotoni/memov2/models"
	"github.com/hirotoni/memov2/utils"
)

type todoFileRepo interface {
	TodoEntries() ([]models.TodoFileInterface, error)
	Save(file models.TodoFileInterface, truncate bool) error
	TodosTemplate(date time.Time) (models.TodoFileInterface, error)
	FindTodosFileByDate(date time.Time) (models.TodoFileInterface, error)
}

type todoFileRepoImpl struct {
	dir string
}

func NewTodoFileRepo(dir string) todoFileRepo {
	return &todoFileRepoImpl{
		dir: dir,
	}
}

func (r *todoFileRepoImpl) TodoEntries() ([]models.TodoFileInterface, error) {
	reg, err := regexp.Compile(models.FileNameRegexTodo)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err) // Invalid regex pattern
	}

	var files []models.TodoFileInterface
	err = filepath.Walk(r.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == r.dir {
			return nil // Skip the root directory itself
		}

		if reg.MatchString(info.Name()) {
			todoFile, err := todofilefrominfo(path, info)
			if err != nil {
				return fmt.Errorf("error creating TodoFile from info: %w", err)
			}

			files = append(files, todoFile)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking through directory: %w", err)
	}

	return files, nil
}

func todofilefrominfo(path string, info os.FileInfo) (models.TodoFileInterface, error) {
	dateReg, err := regexp.Compile(models.FileNameDateTimeRegexTodo)
	if err != nil {
		return nil, fmt.Errorf("invalid date regex pattern: %w", err) // Invalid regex pattern
	}
	datestring := dateReg.FindString(info.Name())
	if datestring == "" {
		return nil, fmt.Errorf("no date found in filename %s", info.Name())
	}

	// 日付
	date, err := time.Parse(models.FileNameDateLayoutTodo, datestring)
	if err != nil {
		return nil, fmt.Errorf("error parsing date from filename %s: %w", info.Name(), err)
	}

	t, err := models.NewTodosFile(date)
	if err != nil {
		return nil, fmt.Errorf("error creating new TodosFile: %w", err)
	}

	h := utils.NewMarkdownHandler()
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err) // Error reading file
	}

	entities, _ := h.HeadingBlocksByLevel(b, 2)

	t.SetDate(date)
	t.SetHeadingBlocks(entities)

	return t, nil
}

func (r *todoFileRepoImpl) Save(file models.TodoFileInterface, truncate bool) error {
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

func (r *todoFileRepoImpl) TodosTemplate(date time.Time) (models.TodoFileInterface, error) {
	fpath := filepath.Join(r.dir, "todos_template.md")

	if !utils.Exists(fpath) {
		t, err := models.NewTodoTemplateFile()
		if err != nil {
			return nil, fmt.Errorf("failed to create todo template file: %w", err)
		}
		err = r.Save(t, false) // Save the template file if it does not exist
		if err != nil {
			return nil, fmt.Errorf("failed to save todo template file: %w", err)
		}
		fmt.Printf("Template file created: %s\n", fpath)
	}

	b, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err // Error reading template file
	}

	h := utils.NewMarkdownHandler()
	hbs, err := h.HeadingBlocksByLevel(b, 2)
	if err != nil {
		return nil, err // Error parsing markdown
	}

	f, err := models.NewTodosFile(date)
	if err != nil {
		return nil, err // Error creating new TodosFile
	}

	f.SetHeadingBlocks(hbs)

	return f, nil
}

func (r *todoFileRepoImpl) FindTodosFileByDate(date time.Time) (models.TodoFileInterface, error) {
	f, err := models.NewTodosFile(date)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(r.dir, f.FileName())

	if !utils.Exists(path) {
		return nil, os.ErrNotExist // File does not exist
	}

	// set entities
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err // Error reading file
	}

	h := utils.NewMarkdownHandler()

	entities, _ := h.HeadingBlocksByLevel(b, 2)

	f.SetDate(date)
	f.SetHeadingBlocks(entities)

	return f, nil
}
