package memo

import (
	"bufio"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/platform"
	"github.com/hirotoni/memov2/internal/utils"
)

func (uc memo) GenerateMemoIndex() error {
	err := uc.TidyMemos()
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error tidying memos")
	}

	s, err := tree(uc.config, uc.config.MemosDir(), 0, true)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error generating memo tree")
	}

	indexPath := filepath.Join(uc.config.MemosDir(), "index.md")
	if err := platform.WriteFileStream(indexPath, true, func(w *bufio.Writer) error {
		if _, err := w.WriteString(s); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error writing to index file")
	}
	uc.logger.Info("Memo index generated successfully", "path", indexPath)
	err = uc.editor.Open(uc.config.BaseDir(), indexPath)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error opening editor")
	}

	return nil
}

func tree(c interfaces.ConfigProvider, path string, level int, root bool) (string, error) {
	sb := strings.Builder{}
	b := utils.NewMarkdownBuilder()
	reg, err := regexp.Compile(domain.FileNameRegexMemo)
	if err != nil {
		return "", common.Wrap(err, common.ErrorTypeService, "invalid regex pattern")
	}

	entries, err := platform.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		name := entry.Name()

		if entry.IsDir() {
			if root {
				h := b.BuildHeading(2, name)
				sb.WriteString("\n")
				sb.WriteString(h + "\n")

				childPath := filepath.Join(path, name)
				s, err := tree(c, childPath, level+1, false)
				if err != nil {
					return "", err
				}
				sb.WriteString(s)
			} else {
				l := b.BuildList(name, level)
				sb.WriteString(l)

				childPath := filepath.Join(path, name)
				s, err := tree(c, childPath, level+1, false)
				if err != nil {
					return "", err
				}
				sb.WriteString(s)
			}
		} else {
			if !reg.MatchString(name) {
				continue
			}

			title := domain.MemoTitle(name)
			fullpath := filepath.Join(path, name)
			rel, err := filepath.Rel(c.MemosDir(), fullpath)
			if err != nil {
				return "", common.Wrap(err, common.ErrorTypeService, "error getting relative path")
			}
			link := b.BuildLink(title, filepath.ToSlash(rel), "")
			l := b.BuildList(link, level)
			sb.WriteString(l)
		}
	}

	return sb.String(), nil
}
