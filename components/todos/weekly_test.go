package todos

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hirotoni/memov2/config"
	"github.com/hirotoni/memov2/models"
	"github.com/hirotoni/memov2/repos"
	"github.com/hirotoni/memov2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildWeeklyReportTodos(t *testing.T) {
	// Create test todos
	date1, err := time.Parse(models.FileNameDateLayoutTodo, "20230101Mon")
	require.NoError(t, err)
	date2, err := time.Parse(models.FileNameDateLayoutTodo, "20230102Tue")
	require.NoError(t, err)
	todo1, err := models.NewTodosFile(date1)
	require.NoError(t, err)
	todo1.SetHeadingBlocks([]*models.HeadingBlock{
		{HeadingText: "Test Heading 1", Level: 2},
		{HeadingText: "Test Heading 2", Level: 2},
		{HeadingText: utils.HeadingTodos.Text, ContentText: "This is a test content for todos.", Level: 2},
	})
	todo2, err := models.NewTodosFile(date2)
	require.NoError(t, err)

	baseDir := t.TempDir()
	todoDir := filepath.Join(baseDir, config.DefaultFolderNameTodos)
	err = os.MkdirAll(todoDir, 0755)
	require.NoError(t, err)

	r := repos.NewTodoFileRepo(todoDir)
	err = r.Save(todo1, true)
	require.NoError(t, err)
	err = r.Save(todo2, true)
	require.NoError(t, err)

	// Execute
	o := config.TomlConfigOption{BaseDir: baseDir}
	cfg, err := config.NewTomlConfig(o)
	require.NoError(t, err)
	err = BuildWeeklyReportTodos(*cfg)

	// Assert
	assert.NoError(t, err)
}
