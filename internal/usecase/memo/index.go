package memo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/platform/editor"
	"github.com/hirotoni/memov2/internal/platform/fs"
	"github.com/hirotoni/memov2/utils"
)

func (uc memo) GenerateMemoIndex() error {
	if err := fs.EnsureDir(uc.config.MemosDir()); err != nil {
		return fmt.Errorf("error ensuring memos directory: %v", err)
	}

	err := uc.repos.Memo.TidyMemos()
	if err != nil {
		return fmt.Errorf("error tidying memos: %v", err)
	}

	s, err := tree(uc.config, uc.config.MemosDir(), 0, true)
	if err != nil {
		return fmt.Errorf("error generating memo tree: %v", err)
	}

	indexPath := filepath.Join(uc.config.MemosDir(), "index.md")
	if err := fs.WriteFileStream(indexPath, true, func(w *bufio.Writer) error {
		if _, err := w.WriteString(s); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("error writing to index file: %v", err)
	}
	fmt.Println("Memo index generated successfully at:", indexPath)
	err = editor.DEO.Open(uc.config.BaseDir(), indexPath)
	if err != nil {
		return fmt.Errorf("error opening editor: %v", err)
	}

	return nil
}

func tree(c config.TomlConfig, path string, level int, root bool) (string, error) {
	sb := strings.Builder{}
	b := utils.NewMarkdownBuilder()
	reg, err := regexp.Compile(domain.FileNameRegexMemo)
	if err != nil {
		return "", fmt.Errorf("invalid regex pattern: %v", err)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("error reading directory %s: %v", path, err)
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
				return "", fmt.Errorf("error getting relative path: %v", err)
			}
			link := b.BuildLink(title, filepath.ToSlash(rel), "")
			l := b.BuildList(link, level)
			sb.WriteString(l)
		}
	}

	return sb.String(), nil
}
