package todos

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/hirotoni/memov2/components"
	"github.com/hirotoni/memov2/config"
	"github.com/hirotoni/memov2/models"
	"github.com/hirotoni/memov2/repos"
	"github.com/hirotoni/memov2/utils"
)

func BuildWeeklyReportTodos(c config.TomlConfig) error {
	if !utils.Exists(c.TodosDir()) {
		err := os.MkdirAll(c.TodosDir(), 0755)
		if err != nil {
			return fmt.Errorf("error creating todos directory: %v", err)
		}
		fmt.Println("Created memos directory:", c.TodosDir())
	}

	fmt.Print("Building weekly report...\n")

	r := repos.NewTodoFileRepo(c.TodosDir())

	todos, err := r.TodoEntries()
	if err != nil {
		fmt.Print("Error fetching todo entries: ", err, "\n")
		return err
	}
	w, err := models.NewWeeklyFile()
	if err != nil {
		fmt.Print("Error creating weekly file: ", err, "\n")
		return err
	}

	var prevWeekNum int
	for i, todo := range todos {
		if i == 0 {
			continue // skip the first iteration
		}

		prev := todos[i-1]
		curr := todo

		if _, week := todo.Date().ISOWeek(); week != prevWeekNum {
			weekHeader := components.WeekSplitter(todo.Date())
			e := &models.HeadingBlock{HeadingText: weekHeader, Level: 2}
			w.SetHeadingBlocks(append(w.HeadingBlocks(), e))
			prevWeekNum = week
		}

		// core logic
		s := generateTodoDiff(prev, curr)
		b := utils.NewMarkdownBuilder()
		l := b.BuildLink(curr.FileName(), curr.FileName(), "")
		if s != "" {
			s = b.BuildCodeBlock(s, "diff")
		}

		e := models.HeadingBlock{
			HeadingText: l,
			Level:       3,
			ContentText: s,
		}

		w.SetHeadingBlocks(append(w.HeadingBlocks(), &e))
	}

	wfr := repos.NewWeeklyFileRepo(c.TodosDir())

	err = wfr.Save(w, true)
	if err != nil {
		fmt.Print("Error saving weekly report: ", err, "\n")
		return err
	}

	fpath := filepath.Join(c.TodosDir(), w.FileName())
	components.OpenEditor(c.BaseDir, fpath)

	return nil

}

func generateTodoDiff(prev, curr models.TodoFileInterface) string {
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
