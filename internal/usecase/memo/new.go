package memo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/platform/editor"
	"github.com/hirotoni/memov2/internal/platform/fs"
)

func (uc memo) GenerateMemoFile(title string) error {
	if err := fs.EnsureDir(uc.c.MemosDir()); err != nil {
		return fmt.Errorf("error ensuring memos directory: %v", err)
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
	memoFile, err := domain.NewMemoFile(today, title, []string{})
	if err != nil {
		return err
	}

	// Save the memo file to the base directory
	err = uc.r.Memo.Save(memoFile, false)
	if err != nil {
		return err
	}

	fpath := filepath.Join(uc.c.MemosDir(), memoFile.Location(), memoFile.FileName())

	err = editor.DEO.Open(uc.c.BaseDir, fpath)
	if err != nil {
		return fmt.Errorf("error opening editor: %v", err)
	}

	return nil
}
