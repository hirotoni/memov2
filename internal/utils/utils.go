package utils

import (
	"fmt"
	"time"
)

// WeekSplitter splits date into year and week
func WeekSplitter(date time.Time) string {
	year, week := date.ISOWeek()
	return fmt.Sprint(year) + " | Week " + fmt.Sprint(week)
}

