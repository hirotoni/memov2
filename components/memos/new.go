package memos

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hirotoni/memov2/components"
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

	components.OpenEditor(c.BaseDir, fpath)

	return nil
}
