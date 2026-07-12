package todo

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInheritTodos_LegacyBacklogHeading verifies the rename transition is safe:
// when the previous day's file still uses the pre-rename "## wanttodos" heading
// (created before `todos migrate`), its content must be carried over into today's
// file under the new "## backlog" heading instead of being silently dropped.
func TestInheritTodos_LegacyBacklogHeading(t *testing.T) {
	uc, _, todosDir := newTestTodoService(t)

	const carried = "- legacy backlog item"
	// Seed yesterday's file with the legacy heading only.
	seedLegacyTodoFile(t, todosDir, time.Now().AddDate(0, 0, -1), carried)

	captureStdout(t, func() {
		require.NoError(t, uc.GenerateTodoFileNonInteractive(false, false))
	})

	// Read today's generated file directly from disk.
	today, err := domain.NewTodosFile(time.Now())
	require.NoError(t, err)
	data, err := os.ReadFile(filepath.Join(todosDir, today.FileName()))
	require.NoError(t, err)
	content := string(data)

	// The legacy content must be inherited under the new "backlog" heading, and
	// the legacy heading text must no longer appear.
	assert.Contains(t, content, "## "+utils.HeadingBacklog.Text)
	assert.Contains(t, content, carried)
	assert.NotContains(t, content, utils.HeadingBacklogLegacy.Text)
}
