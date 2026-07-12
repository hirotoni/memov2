package todo

import (
	"bufio"
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/platform"
	repoCommon "github.com/hirotoni/memov2/internal/repositories/common"
)

type todo struct {
	dir    string
	logger *slog.Logger
}

func NewTodo(dir string, logger *slog.Logger) interfaces.TodoRepo {
	return &todo{dir: dir, logger: logger}
}

func (r *todo) TodoEntries() ([]interfaces.TodoFileInterface, error) {
	reg, err := regexp.Compile(domain.FileNameRegexTodo)
	if err != nil {
		return nil, common.Wrap(err, common.ErrorTypeRepository, "invalid regex pattern")
	}

	var files []interfaces.TodoFileInterface
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
				return common.Wrap(err, common.ErrorTypeRepository, "error creating TodoFile from info")
			}

			files = append(files, todoFile)
		}

		return nil
	})
	if err != nil {
		return nil, common.Wrap(err, common.ErrorTypeFileSystem, "error walking through directory")
	}

	return files, nil
}

func todofilefrominfo(path string, info os.FileInfo) (interfaces.TodoFileInterface, error) {
	// 日付抽出（共通パーサーを使用）
	date, err := repoCommon.ParseDateFromFilename(info.Name(), repoCommon.DateParserConfig{
		DateTimeRegex: domain.FileNameDateTimeRegexTodo,
		DateLayout:    domain.FileNameDateLayoutTodo,
	})
	if err != nil {
		return nil, err
	}

	// Markdownファイル読み込み（共通パーサーを使用）
	b, err := repoCommon.ReadMarkdownFile(path)
	if err != nil {
		return nil, err
	}

	// Markdown解析（共通パーサーを使用）
	parser := repoCommon.NewMarkdownParser()
	entities, _ := parser.HeadingBlocksByLevel(b, 2)

	// Domain層のファクトリを使用してTodoFileを構築
	return domain.TodoFileFromParsedData(date, entities)
}

func (r *todo) Save(file interfaces.TodoFileInterface, truncate bool) error {
	path := filepath.Join(r.dir, file.FileName())
	if !platform.Exists(path) || truncate {
		err := platform.WriteFileStream(path, truncate, func(w *bufio.Writer) error {
			_, err := w.WriteString(file.ContentString())
			return err
		})
		if err != nil {
			return err
		}
		r.logger.Info("File saved", "path", path)
	}
	return nil
}

func (r *todo) TodosTemplate(date time.Time) (interfaces.TodoFileInterface, error) {
	fpath := filepath.Join(r.dir, "todos_template.md")

	if !platform.Exists(fpath) {
		t, err := domain.NewTodoTemplateFile()
		if err != nil {
			return nil, common.Wrap(err, common.ErrorTypeRepository, "failed to create todo template file")
		}
		err = r.Save(t, false) // Save the template file if it does not exist
		if err != nil {
			return nil, common.Wrap(err, common.ErrorTypeRepository, "failed to save todo template file")
		}
		r.logger.Info("Template file created", "path", fpath)
	}

	b, err := repoCommon.ReadMarkdownFile(fpath)
	if err != nil {
		return nil, err // Error reading template file
	}

	parser := repoCommon.NewMarkdownParser()
	hbs, err := parser.HeadingBlocksByLevel(b, 2)
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

// RewriteHeadingText rewrites every level-2 heading whose text is oldText to
// newText across the todo files (daily files matching FileNameRegexTodo and the
// todos_template.md) in the todo directory. The rewrite is a targeted,
// line-level replacement so the rest of each file's content is preserved
// byte-for-byte. When dryRun is true, files are inspected but not modified.
func (r *todo) RewriteHeadingText(oldText, newText string, dryRun bool) ([]string, error) {
	reg, err := regexp.Compile(domain.FileNameRegexTodo)
	if err != nil {
		return nil, common.Wrap(err, common.ErrorTypeRepository, "invalid regex pattern")
	}

	// Matches a level-2 heading line whose text is exactly oldText, e.g. "## wanttodos".
	headingRe := regexp.MustCompile(`(?m)^##[ \t]+` + regexp.QuoteMeta(oldText) + `[ \t]*$`)
	replacement := "## " + newText

	var changed []string
	err = filepath.Walk(r.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		name := info.Name()
		if !reg.MatchString(name) && name != "todos_template.md" {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return common.Wrap(err, common.ErrorTypeFileSystem, "error reading todo file")
		}

		out := headingRe.ReplaceAll(b, []byte(replacement))
		if bytes.Equal(out, b) {
			return nil // heading not present; nothing to rewrite
		}

		rel, relErr := filepath.Rel(r.dir, path)
		if relErr != nil {
			rel = path
		}
		changed = append(changed, rel)

		if dryRun {
			return nil
		}

		if err := os.WriteFile(path, out, info.Mode().Perm()); err != nil {
			return common.Wrap(err, common.ErrorTypeFileSystem, "error writing todo file")
		}
		r.logger.Info("Heading rewritten", "path", path, "from", oldText, "to", newText)
		return nil
	})
	if err != nil {
		return nil, common.Wrap(err, common.ErrorTypeFileSystem, "error walking through directory")
	}

	return changed, nil
}

func (r *todo) FindTodosFileByDate(date time.Time) (interfaces.TodoFileInterface, error) {
	f, err := domain.NewTodosFile(date)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(r.dir, f.FileName())

	if !platform.Exists(path) {
		return nil, os.ErrNotExist // File does not exist
	}

	// set entities
	b, err := repoCommon.ReadMarkdownFile(path)
	if err != nil {
		return nil, err // Error reading file
	}

	parser := repoCommon.NewMarkdownParser()
	entities, _ := parser.HeadingBlocksByLevel(b, 2)

	f.SetDate(date)
	f.SetHeadingBlocks(entities)

	return f, nil
}
