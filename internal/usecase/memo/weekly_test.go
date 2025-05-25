package memo

import (
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/repository"
	"github.com/hirotoni/memov2/internal/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildWeeklyReport_Success tests the success case
func TestBuildWeeklyReport_Success(t *testing.T) {
	m := mock.NewMock()

	// Create test memos
	date, err := time.Parse(domain.FileNameDateLayoutTodo, "20230101Mon")
	require.NoError(t, err)
	memo1, err := domain.NewMemoFile(date, "Test Memo 1", []string{})
	memo1.SetHeadingBlocks([]*markdown.HeadingBlock{
		{HeadingText: "Test Heading 1"},
		{HeadingText: "Test Heading 2"},
	})
	require.NoError(t, err)
	memo2, err := domain.NewMemoFile(date, "Test Memo 2", []string{})
	require.NoError(t, err)
	m.MockMemoEntries = func() ([]domain.MemoFileInterface, error) {
		return []domain.MemoFileInterface{memo1, memo2}, nil
	}

	o := config.TomlConfigOption{BaseDir: t.TempDir()}
	cfg, err := config.NewTomlConfig(o)
	require.NoError(t, err)
	// Execute
	r := repository.NewRepositories(*cfg)
	uc := NewMemo(*cfg, r)

	err = uc.BuildWeeklyReportMemos()

	// Assert
	assert.NoError(t, err)
}
