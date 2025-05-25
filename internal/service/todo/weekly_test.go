package todo

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/repositories"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	todoRepo "github.com/hirotoni/memov2/internal/repositories/todo"
	"github.com/hirotoni/memov2/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildWeeklyReportTodos(t *testing.T) {
	// Create test todos
	date1, err := time.Parse(domain.FileNameDateLayoutTodo, "20230101Mon")
	require.NoError(t, err)
	date2, err := time.Parse(domain.FileNameDateLayoutTodo, "20230102Tue")
	require.NoError(t, err)
	todo1, err := domain.NewTodosFile(date1)
	require.NoError(t, err)
	todo1.SetHeadingBlocks([]*markdown.HeadingBlock{
		{HeadingText: "Test Heading 1", Level: 2},
		{HeadingText: "Test Heading 2", Level: 2},
		{HeadingText: utils.HeadingTodos.Text, ContentText: "This is a test content for todos.", Level: 2},
	})
	todo2, err := domain.NewTodosFile(date2)
	require.NoError(t, err)

	baseDir := t.TempDir()
	todoDir := filepath.Join(baseDir, config.DefaultFolderNameTodos)
	err = os.MkdirAll(todoDir, 0o755)
	require.NoError(t, err)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	r := todoRepo.NewTodo(todoDir, logger)
	err = r.Save(todo1, true)
	require.NoError(t, err)
	err = r.Save(todo2, true)
	require.NoError(t, err)

	// Execute
	o := toml.Option{BaseDir: baseDir}
	cfg, err := toml.NewConfig(o)
	require.NoError(t, err)

	configProvider := toml.NewProvider(cfg)
	rs := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()
	uc := NewTodo(configProvider, rs, mockEditor, logger)
	err = uc.BuildWeeklyReportTodos()

	// Assert
	assert.NoError(t, err)
}
