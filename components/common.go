package components

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

// OpenEditor opens editor
func OpenEditor(basedir, path string) {
	cmd := exec.Command("code", path, "--folder-uri", basedir)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// WeekSplitter splits date into year and week
func WeekSplitter(date time.Time) string {
	year, week := date.ISOWeek()
	return fmt.Sprint(year) + " | Week " + fmt.Sprint(week)
}
