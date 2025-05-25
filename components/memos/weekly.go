package memos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hirotoni/memov2/components"
	"github.com/hirotoni/memov2/config"
	"github.com/hirotoni/memov2/models"
	"github.com/hirotoni/memov2/repos"
	"github.com/hirotoni/memov2/utils"
)

func isSameWithPrevDate(e models.TodoFileInterface, date string) bool {
	if le := e.LastHeadingBlock(); le != nil && le.HeadingText == date {
		return true
	}
	return false
}

// BuildWeeklyReportMemos builds weekly report
func BuildWeeklyReportMemos(c config.TomlConfig) error {
	if !utils.Exists(c.MemosDir()) {
		err := os.MkdirAll(c.MemosDir(), 0755)
		if err != nil {
			return fmt.Errorf("error creating memos directory: %v", err)
		}
		fmt.Println("Created memos directory:", c.MemosDir())
	}

	fmt.Print("Building weekly report...\n")

	r := repos.NewMemoRepo(c.MemosDir())
	err := r.TidyMemos()
	if err != nil {
		fmt.Print("Error tidying memos: ", err, "\n")
	}

	memos, err := r.MemoEntries()
	if err != nil {
		fmt.Print("Error fetching memo entries: ", err, "\n")
		return err
	}

	w, err := models.NewWeeklyFile()
	if err != nil {
		fmt.Print("Error creating weekly file: ", err, "\n")
		return err
	}

	var prevWeekNum int
	var order int
	var b = utils.NewMarkdownBuilder()
	for _, memo := range memos {
		// Add week header if week number changes
		if _, week := memo.Date().ISOWeek(); week != prevWeekNum {
			weekHeader := components.WeekSplitter(memo.Date())
			e := &models.HeadingBlock{HeadingText: weekHeader, Level: 2}
			w.SetHeadingBlocks(append(w.HeadingBlocks(), e))
			prevWeekNum = week
		}

		// determine if new date or same as previous
		var e *models.HeadingBlock
		date := memo.Date().Format(models.FileNameDateLayoutTodo) // date
		sameWithPrevDate := isSameWithPrevDate(w, date)

		if sameWithPrevDate {
			order++
			e = w.LastHeadingBlock() // reuse last entity if date is the same
		} else {
			order = 1
			e = &models.HeadingBlock{HeadingText: date, Level: 3} // new entity for new date
		}

		// memo title
		var tt string
		link := b.BuildLink(memo.Title(), filepath.Join(memo.Location(), memo.FileName()), memo.Title())
		tt += b.BuildOrderedList(order, link, 1, 1)

		// memo headings
		var innerOrder int
		for _, entity := range memo.HeadingBlocks() {
			innerOrder++
			link := b.BuildLink(entity.HeadingText, filepath.Join(memo.Location(), memo.FileName()), entity.HeadingText)
			tt += b.BuildOrderedList(innerOrder, link, 2, order)
		}
		e.ContentText = e.ContentText + tt // append content

		if sameWithPrevDate {
			etts := w.HeadingBlocks()
			newEtts := append(etts[:len(etts)-1], e) // replace last entity
			w.SetHeadingBlocks(newEtts)              // set updated entities
		} else {
			w.SetHeadingBlocks(append(w.HeadingBlocks(), e))
		}
	}

	wfr := repos.NewWeeklyFileRepo(c.MemosDir())

	err = wfr.Save(w, true)
	if err != nil {
		fmt.Print("Error saving weekly report: ", err, "\n")
		return err
	}

	fpath := filepath.Join(c.MemosDir(), w.FileName())
	components.OpenEditor(c.BaseDir, fpath)

	return nil
}
