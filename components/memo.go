package components

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hirotoni/memov2/config"
	"github.com/hirotoni/memov2/models"
	"github.com/hirotoni/memov2/repos"
	"github.com/hirotoni/memov2/utils"
)

func GenerateMemoFile(c config.TomlConfig, title string) error {
	if !utils.Exists(c.MemosDir()) {
		err := os.MkdirAll(c.MemosDir(), 0755)
		if err != nil {
			return fmt.Errorf("error creating memos directory: %v", err)
		}
		fmt.Println("Created memos directory:", c.MemosDir())
	}

	if title == "" {
		fmt.Print("Title: ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return fmt.Errorf("canceled")
		}
		if scanner.Err() != nil {
			return fmt.Errorf("error reading input: %v", scanner.Err())
		}
		title = scanner.Text()
	}

	// Create a new memo file with the given title
	today := time.Now()
	memoFile, err := models.NewMemoFile(today, title, []string{})
	if err != nil {
		return err
	}

	// Save the memo file to the base directory
	repo := repos.NewMemoRepo(c.MemosDir())
	err = repo.Save(memoFile, false)
	if err != nil {
		return err
	}

	fpath := filepath.Join(c.MemosDir(), memoFile.Location(), memoFile.FileName())

	OpenEditor(c.BaseDir, fpath)

	return nil
}

func GenerateMemoIndex(c config.TomlConfig) error {
	if !utils.Exists(c.MemosDir()) {
		err := os.MkdirAll(c.MemosDir(), 0755)
		if err != nil {
			return fmt.Errorf("error creating memos directory: %v", err)
		}
		fmt.Println("Created memos directory:", c.MemosDir())
	}

	// reg, err := regexp.Compile(models.MemoFileNameRegex)
	// if err != nil {
	// 	return err // Invalid regex pattern
	// }

	s, err := tree(c, c.MemosDir(), 2)
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
	OpenEditor(c.BaseDir, filepath.Join(c.MemosDir(), "index.md"))

	return nil
}

func tree(c config.TomlConfig, path string, level int) (string, error) {
	sb := strings.Builder{}
	b := utils.NewMarkdownBuilder()
	reg, err := regexp.Compile(models.MemoFileNameRegex)
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
			h := b.BuildHeading(level, name)
			sb.WriteString(h + "\n")

			childPath := filepath.Join(path, name)
			s, err := tree(c, childPath, level+1)
			if err != nil {
				return "", err
			}
			sb.WriteString(s)
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
			l := b.BuildList(link, 0)
			sb.WriteString(l)
		}
	}

	return sb.String(), nil
}
