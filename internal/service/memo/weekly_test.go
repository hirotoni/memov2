package memo

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/repositories"
	"github.com/hirotoni/memov2/internal/repositories/mock"
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

	o := toml.Option{BaseDir: t.TempDir()}
	cfg, err := toml.NewConfig(o)
	require.NoError(t, err)
	// Execute
	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	r := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()
	uc := NewMemo(configProvider, r, mockEditor, logger)

	err = uc.BuildWeeklyReportMemos()

	// Assert
	assert.NoError(t, err)
}

func TestBuildWeeklyReport_MultipleMemosInSameWeek(t *testing.T) {
	m := mock.NewMock()

	// Create test memos on different days in the same week
	date1, err := time.Parse(domain.FileNameDateLayoutTodo, "20230101Mon")
	require.NoError(t, err)
	date2, err := time.Parse(domain.FileNameDateLayoutTodo, "20230102Tue")
	require.NoError(t, err)
	date3, err := time.Parse(domain.FileNameDateLayoutTodo, "20230103Wed")
	require.NoError(t, err)

	memo1, err := domain.NewMemoFile(date1, "Memo 1", []string{})
	require.NoError(t, err)
	memo1.SetHeadingBlocks([]*markdown.HeadingBlock{
		{HeadingText: "Heading 1", Level: 2},
	})

	memo2, err := domain.NewMemoFile(date2, "Memo 2", []string{})
	require.NoError(t, err)
	memo2.SetHeadingBlocks([]*markdown.HeadingBlock{
		{HeadingText: "Heading 2", Level: 2},
	})

	memo3, err := domain.NewMemoFile(date3, "Memo 3", []string{})
	require.NoError(t, err)

	m.MockMemoEntries = func() ([]domain.MemoFileInterface, error) {
		return []domain.MemoFileInterface{memo1, memo2, memo3}, nil
	}

	o := toml.Option{BaseDir: t.TempDir()}
	cfg, err := toml.NewConfig(o)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	r := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()
	uc := NewMemo(configProvider, r, mockEditor, logger)

	err = uc.BuildWeeklyReportMemos()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockEditor.Calls), "Editor should be called once")
}

func TestBuildWeeklyReport_EmptyMemos(t *testing.T) {
	m := mock.NewMock()

	m.MockMemoEntries = func() ([]domain.MemoFileInterface, error) {
		return []domain.MemoFileInterface{}, nil
	}

	o := toml.Option{BaseDir: t.TempDir()}
	cfg, err := toml.NewConfig(o)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	r := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()
	uc := NewMemo(configProvider, r, mockEditor, logger)

	err = uc.BuildWeeklyReportMemos()

	// Assert
	assert.NoError(t, err)
}

func TestBuildWeeklyReport_MemoWithCategories(t *testing.T) {
	m := mock.NewMock()

	date, err := time.Parse(domain.FileNameDateLayoutTodo, "20230101Mon")
	require.NoError(t, err)

	memo, err := domain.NewMemoFile(date, "Categorized Memo", []string{"cat1", "cat2"})
	require.NoError(t, err)
	memo.SetHeadingBlocks([]*markdown.HeadingBlock{
		{HeadingText: "Section 1", ContentText: "Content 1", Level: 2},
		{HeadingText: "Section 2", ContentText: "Content 2", Level: 2},
	})

	m.MockMemoEntries = func() ([]domain.MemoFileInterface, error) {
		return []domain.MemoFileInterface{memo}, nil
	}

	o := toml.Option{BaseDir: t.TempDir()}
	cfg, err := toml.NewConfig(o)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	r := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()
	uc := NewMemo(configProvider, r, mockEditor, logger)

	err = uc.BuildWeeklyReportMemos()

	// Assert
	assert.NoError(t, err)
}

func TestBuildWeeklyReport_MultipleWeeks(t *testing.T) {
	m := mock.NewMock()

	// Create memos from different weeks
	date1, err := time.Parse(domain.FileNameDateLayoutTodo, "20230101Mon")
	require.NoError(t, err)
	date2, err := time.Parse(domain.FileNameDateLayoutTodo, "20230108Mon")
	require.NoError(t, err)

	memo1, err := domain.NewMemoFile(date1, "Week 1 Memo", []string{})
	require.NoError(t, err)
	memo2, err := domain.NewMemoFile(date2, "Week 2 Memo", []string{})
	require.NoError(t, err)

	m.MockMemoEntries = func() ([]domain.MemoFileInterface, error) {
		return []domain.MemoFileInterface{memo1, memo2}, nil
	}

	o := toml.Option{BaseDir: t.TempDir()}
	cfg, err := toml.NewConfig(o)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	r := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()
	uc := NewMemo(configProvider, r, mockEditor, logger)

	err = uc.BuildWeeklyReportMemos()

	// Assert
	assert.NoError(t, err)
}

func TestIsSameWithPrevDate(t *testing.T) {
	// Create a test weekly file
	weekly, err := domain.NewWeekly()
	require.NoError(t, err)

	// Test case 1: No headings yet
	result := isSameWithPrevDate(weekly, "20230101Mon")
	assert.False(t, result, "Should return false when no heading blocks exist")

	// Test case 2: Add a heading block with matching date
	weekly.SetHeadingBlocks([]*markdown.HeadingBlock{
		{HeadingText: "20230101Mon", Level: 3},
	})
	result = isSameWithPrevDate(weekly, "20230101Mon")
	assert.True(t, result, "Should return true when last heading matches date")

	// Test case 3: Different date
	result = isSameWithPrevDate(weekly, "20230102Tue")
	assert.False(t, result, "Should return false when date doesn't match")

	// Test case 4: Multiple headings, check against last one
	weekly.SetHeadingBlocks([]*markdown.HeadingBlock{
		{HeadingText: "20230101Mon", Level: 3},
		{HeadingText: "20230102Tue", Level: 3},
	})
	result = isSameWithPrevDate(weekly, "20230102Tue")
	assert.True(t, result, "Should check against last heading")
}
