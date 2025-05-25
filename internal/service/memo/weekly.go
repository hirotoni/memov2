package memo

import (
	"fmt"
	"path/filepath"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/utils"
)

func isSameWithPrevDate(e domain.TodoFileInterface, date string) bool {
	if le := e.LastHeadingBlock(); le != nil && le.HeadingText == date {
		return true
	}
	return false
}

func (uc memo) BuildWeeklyReportMemos() error {
	fmt.Print("Building weekly report...\n")

	err := uc.TidyMemos()
	if err != nil {
		// continue even if error
		fmt.Print("Error tidying memos: ", err, "\n")
	}

	memos, err := uc.repos.Memo().MemoEntries()
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error fetching memo entries")
	}

	w, err := domain.NewWeekly()
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error creating weekly file")
	}

	var prevWeekNum int
	var order int
	var b = utils.NewMarkdownBuilder()
	for _, memo := range memos {
		// Add week header if week number changes
		if _, week := memo.Date().ISOWeek(); week != prevWeekNum {
			weekHeader := utils.WeekSplitter(memo.Date())
			e := &markdown.HeadingBlock{HeadingText: weekHeader, Level: 2}
			w.SetHeadingBlocks(append(w.HeadingBlocks(), e))
			prevWeekNum = week
		}

		// determine if new date or same as previous
		var e *markdown.HeadingBlock
		date := memo.Date().Format(domain.FileNameDateLayoutTodo) // date
		sameWithPrevDate := isSameWithPrevDate(w, date)

		if sameWithPrevDate {
			order++
			e = w.LastHeadingBlock() // reuse last entity if date is the same
		} else {
			order = 1
			e = &markdown.HeadingBlock{HeadingText: date, Level: 3} // new entity for new date
		}

		// memo title
		var tt string
		path := filepath.ToSlash(filepath.Join(memo.Location(), memo.FileName()))
		link := b.BuildLink(memo.Title(), path, memo.Title())
		tt += b.BuildOrderedList(order, link, 1, 1)

		// memo headings
		var innerOrder int
		for _, entity := range memo.HeadingBlocks() {
			innerOrder++
			path := filepath.ToSlash(filepath.Join(memo.Location(), memo.FileName()))
			link := b.BuildLink(entity.HeadingText, path, entity.HeadingText)
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

	err = uc.repos.MemoWeekly().Save(w, true)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error saving weekly report")
	}

	fpath := filepath.Join(uc.config.MemosDir(), w.FileName())
	err = uc.editor.Open(uc.config.BaseDir(), fpath)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error opening editor")
	}

	return nil
}
