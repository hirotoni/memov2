/*
Copyright © 2025 hirotoni
*/
package memos

import (
	"github.com/hirotoni/memov2/cmd/app"
	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "generate memo index file",
	Long:  `generate a memo index file that lists all memos.`,
	// RunE + SilenceUsage so a failure exits non-zero with the error on stderr.
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ap, err := app.InitializeApp(cmd)
		if err != nil {
			return err
		}
		return ap.Services().Memo().GenerateMemoIndex()
	},
}

func init() {}
