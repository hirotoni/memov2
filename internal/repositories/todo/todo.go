package todo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/platform/fs"
	"github.com/hirotoni/memov2/utils"
)

type todo struct {
	dir string
}

func NewTodo(dir string) interfaces.TodoRepo {
	return &todo{dir: dir}
}

func (r *todo) TodoEntries() ([]domain.TodoFileInterface, error) {
	reg, err := regexp.Compile(domain.FileNameRegexTodo)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err) // Invalid regex pattern
	}

	var files []domain.TodoFileInterface
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

func todofilefrominfo(path string, info os.FileInfo) (domain.TodoFileInterface, error) {
	dateReg, err := regexp.Compile(domain.FileNameDateTimeRegexTodo)
	if err != nil {
		return nil, fmt.Errorf("invalid date regex pattern: %w", err) // Invalid regex pattern
	}
	datestring := dateReg.FindString(info.Name())
	if datestring == "" {
		return nil, fmt.Errorf("no date found in filename %s", info.Name())
	}

	// 日付
	date, err := time.Parse(domain.FileNameDateLayoutTodo, datestring)
	if err != nil {
		return nil, fmt.Errorf("error parsing date from filename %s: %w", info.Name(), err)
	}

	t, err := domain.NewTodosFile(date)
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

func (r *todo) Save(file domain.TodoFileInterface, truncate bool) error {
	path := filepath.Join(r.dir, file.FileName())
	if !fs.Exists(path) || truncate {
		err := fs.WriteFileStream(path, truncate, func(w *bufio.Writer) error {
			_, err := w.WriteString(file.ContentString())
			return err
		})
		if err != nil {
			return err
		}
		fmt.Printf("File saved: %s\n", path)
	}
	return nil
}

func (r *todo) TodosTemplate(date time.Time) (domain.TodoFileInterface, error) {
	fpath := filepath.Join(r.dir, "todos_template.md")

	if !fs.Exists(fpath) {
		t, err := domain.NewTodoTemplateFile()
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

	f, err := domain.NewTodosFile(date)
	if err != nil {
		return nil, err // Error creating new TodosFile
	}

	f.SetHeadingBlocks(hbs)

	return f, nil
}

func (r *todo) FindTodosFileByDate(date time.Time) (domain.TodoFileInterface, error) {
	f, err := domain.NewTodosFile(date)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(r.dir, f.FileName())

	if !fs.Exists(path) {
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
