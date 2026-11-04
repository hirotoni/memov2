/*
Copyright © 2025 hirotoni
*/
package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

var weeklyCmd = &cobra.Command{
	Use:   "weekly",
	Short: "generate weekly report for memos",
	Long:  `generate weekly report for memos`,
	// RunE + SilenceUsage so a failure exits non-zero with the error on stderr.
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			return err
		}
		return ap.Services().Memo().BuildWeeklyReportMemos()
	},
}

func init() {}
