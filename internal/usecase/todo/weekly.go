package todo

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/platform/fs"
	"github.com/hirotoni/memov2/internal/repository"
	"github.com/hirotoni/memov2/internal/usecase"
	"github.com/hirotoni/memov2/utils"
)

func (uc todo) BuildWeeklyReportTodos() error {
	if err := fs.EnsureDir(uc.c.TodosDir()); err != nil {
		return fmt.Errorf("error ensuring todos directory: %v", err)
	}

	fmt.Print("Building weekly report...\n")

	r := repository.NewTodo(uc.c.TodosDir())

	todos, err := r.TodoEntries()
	if err != nil {
		return fmt.Errorf("error fetching todo entries: %v", err)
	}
	dw, err := domain.NewWeekly()
	if err != nil {
		return fmt.Errorf("error creating weekly file: %v", err)
	}

	var prevWeekNum int
	for i, todo := range todos {
		if i == 0 {
			continue // skip the first iteration
		}

		prev := todos[i-1]
		curr := todo

		if _, week := todo.Date().ISOWeek(); week != prevWeekNum {
			weekHeader := usecase.WeekSplitter(todo.Date())
			e := &markdown.HeadingBlock{HeadingText: weekHeader, Level: 2}
			dw.SetHeadingBlocks(append(dw.HeadingBlocks(), e))
			prevWeekNum = week
		}

		// core logic
		s := generateTodoDiff(prev, curr)
		b := utils.NewMarkdownBuilder()
		l := b.BuildLink(curr.FileName(), filepath.ToSlash(curr.FileName()), "")
		if s != "" {
			s = b.BuildCodeBlock(s, "diff")
		}

		e := markdown.HeadingBlock{
			HeadingText: l,
			Level:       3,
			ContentText: s,
		}

		dw.SetHeadingBlocks(append(dw.HeadingBlocks(), &e))
	}

	err = uc.r.TodoWeekly().Save(dw, true)
	if err != nil {
		return fmt.Errorf("error saving weekly report: %v", err)
	}

	fpath := filepath.Join(uc.c.TodosDir(), dw.FileName())
	err = uc.e.Open(uc.c.BaseDir(), fpath)
	if err != nil {
		return fmt.Errorf("error opening editor: %v", err)
	}

	return nil
}

func generateTodoDiff(prev, curr domain.TodoFileInterface) string {
	var (
		prevTodos string
		currTodos string
	)

	for _, hb := range prev.HeadingBlocks() {
		if hb.HeadingText == utils.HeadingTodos.Text {
			prevTodos = hb.ContentText
			break
		}
	}

	for _, hb := range curr.HeadingBlocks() {
		if hb.HeadingText == utils.HeadingTodos.Text {
			currTodos = hb.ContentText
			break
		}
	}

	s := todoDiff(prev.FileName(), prevTodos, curr.FileName(), currTodos)
	s = strings.Trim(s, "\n")

	return s
}

func todoDiff(aname, atext, bname, btext string) string {
	edits := myers.ComputeEdits(span.URIFromPath(aname), atext, btext)
	diff := fmt.Sprint(gotextdiff.ToUnified(aname, bname, atext, edits))
	return fmt.Sprintln(diff)
}
