package memos

import (
	"testing"
	"time"

	"github.com/hirotoni/memov2/config"
	"github.com/hirotoni/memov2/models"
	"github.com/hirotoni/memov2/repos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildWeeklyReport_Success tests the success case
func TestBuildWeeklyReport_Success(t *testing.T) {
	m := repos.NewMockFileRepo()

	// Create test memos
	date, err := time.Parse(models.FileNameDateLayoutTodo, "20230101Mon")
	require.NoError(t, err)
	memo1, err := models.NewMemoFile(date, "Test Memo 1", []string{})
	memo1.SetHeadingBlocks([]*models.HeadingBlock{
		{HeadingText: "Test Heading 1"},
		{HeadingText: "Test Heading 2"},
	})
	require.NoError(t, err)
	memo2, err := models.NewMemoFile(date, "Test Memo 2", []string{})
	require.NoError(t, err)
	m.MockMemoEntries = func() ([]models.MemoFileInterface, error) {
		return []models.MemoFileInterface{memo1, memo2}, nil
	}

	o := config.TomlConfigOption{BaseDir: t.TempDir()}
	cfg := config.NewTomlConfig(o)
	// Execute
	err = BuildWeeklyReportMemos(*cfg)

	// Assert
	assert.NoError(t, err)
}
