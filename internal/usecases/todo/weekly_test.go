package todo

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/repositories"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	todoRepo "github.com/hirotoni/memov2/internal/repositories/todo"
	"github.com/hirotoni/memov2/utils"
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

	r := todoRepo.NewTodo(todoDir)
	err = r.Save(todo1, true)
	require.NoError(t, err)
	err = r.Save(todo2, true)
	require.NoError(t, err)

	// Execute
	o := config.TomlConfigOption{BaseDir: baseDir}
	cfg, err := config.NewTomlConfig(o)
	require.NoError(t, err)

	configProvider := interfaces.NewConfigProvider(cfg)
	rs := repositories.NewRepositories(configProvider)
	mockEditor := mock.NewMockEditor()
	uc := NewTodo(configProvider, rs, mockEditor)
	err = uc.BuildWeeklyReportTodos()

	// Assert
	assert.NoError(t, err)
}
