package memos

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hirotoni/memov2/components"
	"github.com/hirotoni/memov2/config"
	"github.com/hirotoni/memov2/models"
	"github.com/hirotoni/memov2/repos"
	"github.com/hirotoni/memov2/utils"
)

func GenerateMemoIndex(c config.TomlConfig) error {
	if !utils.Exists(c.MemosDir()) {
		err := os.MkdirAll(c.MemosDir(), 0755)
		if err != nil {
			return fmt.Errorf("error creating memos directory: %v", err)
		}
		fmt.Println("Created memos directory:", c.MemosDir())
	}

	r := repos.NewMemoRepo(c.MemosDir())
	err := r.TidyMemos()
	if err != nil {
		return fmt.Errorf("error tidying memos: %v", err)
	}

	s, err := tree(c, c.MemosDir(), 0, true)
	if err != nil {
		return fmt.Errorf("error generating memo tree: %v", err)
	}

	f, err := os.Create(filepath.Join(c.MemosDir(), "index.md"))
	if err != nil {
		return fmt.Errorf("error creating index file: %v", err)
	}
	defer f.Close()
	_, err = f.WriteString(s)
	if err != nil {
		return fmt.Errorf("error writing to index file: %v", err)
	}
	fmt.Println("Memo index generated successfully at:", filepath.Join(c.MemosDir(), "index.md"))
	components.OpenEditor(c.BaseDir, filepath.Join(c.MemosDir(), "index.md"))

	return nil
}

func tree(c config.TomlConfig, path string, level int, root bool) (string, error) {
	sb := strings.Builder{}
	b := utils.NewMarkdownBuilder()
	reg, err := regexp.Compile(models.FileNameRegexMemo)
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

			title := models.MemoTitle(name)
			fullpath := filepath.Join(path, name)
			rel, err := filepath.Rel(c.MemosDir(), fullpath)
			if err != nil {
				return "", fmt.Errorf("error getting relative path: %v", err)
			}
			link := b.BuildLink(title, rel, "")
			l := b.BuildList(link, level)
			sb.WriteString(l)
		}
	}

	return sb.String(), nil
}
