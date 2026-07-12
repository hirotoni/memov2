/*
Copyright © 2025 hirotoni
*/
package todos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

// migrateCmd renames the legacy "wanttodos" heading to "backlog" across existing
// todo files (the daily files and the template). New files already use "backlog";
// this brings older files in line. Use --dry-run to preview the affected files
// without modifying them.
var migrateCmd = &cobra.Command{
	Use:           "migrate",
	Short:         `migrate the legacy "wanttodos" heading to "backlog"`,
	Long:          `Rewrite the legacy "## wanttodos" heading to "## backlog" in all existing todo files (daily files and the template). New todo files already use "backlog"; run this once to migrate files created before the rename. Use --dry-run to preview the affected files without modifying them.`,
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			return err
		}

		return ap.Services().Todo().MigrateBacklogHeading(migrateDryRunFlag)
	},
}

var migrateDryRunFlag bool

func init() {
	migrateCmd.Flags().BoolVar(&migrateDryRunFlag, "dry-run", false, "preview the files that would change without modifying them")
}
