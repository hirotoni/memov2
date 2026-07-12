package todo

import (
	"fmt"
	"os"

	"github.com/hirotoni/memov2/internal/utils"
)

// MigrateBacklogHeading renames the legacy "## wanttodos" heading to "## backlog"
// across all existing todo files (daily files and the template) in the todo
// directory. New todo files already use "backlog"; this brings older files in
// line so that backlog contents keep being inherited after the rename.
//
// When dryRun is true no file is modified; the files that would change are
// reported instead. A human-readable summary is written to stdout.
func (uc todo) MigrateBacklogHeading(dryRun bool) error {
	changed, err := uc.r.Todo().RewriteHeadingText(
		utils.HeadingBacklogLegacy.Text,
		utils.HeadingBacklog.Text,
		dryRun,
	)
	if err != nil {
		return err
	}

	if len(changed) == 0 {
		fmt.Fprintf(os.Stdout, "No files to migrate: no %q heading found.\n", utils.HeadingBacklogLegacy.Text)
		return nil
	}

	if dryRun {
		fmt.Fprintf(os.Stdout, "Would rewrite %q -> %q in %d file(s):\n", utils.HeadingBacklogLegacy.Text, utils.HeadingBacklog.Text, len(changed))
	} else {
		fmt.Fprintf(os.Stdout, "Rewrote %q -> %q in %d file(s):\n", utils.HeadingBacklogLegacy.Text, utils.HeadingBacklog.Text, len(changed))
	}
	for _, p := range changed {
		fmt.Fprintf(os.Stdout, "  %s\n", p)
	}

	return nil
}
