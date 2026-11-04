package todo

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateTodoFile: with truncate=false, today's daily todo file is created
// in the configured todos directory and opened in the editor. The service is
// backed by a temp base dir so the test never reads or writes the runner's real
// config (~/.config/memov2) or memospace.
func TestGenerateTodoFile(t *testing.T) {
	uc, _, todosDir, editor := newTestTodoServiceE(t)

	require.NoError(t, uc.GenerateTodoFile(false))

	name, _ := readTodayTodoFile(t, todosDir)
	require.Len(t, editor.Calls, 1)
	assert.Equal(t, filepath.Join(todosDir, name), editor.Calls[0].Path)
}
