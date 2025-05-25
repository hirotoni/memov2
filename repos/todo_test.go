package repos

import (
	"testing"
	"time"

	"github.com/hirotoni/memov2/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_todoFileRepoImpl_TodoEntries(t *testing.T) {
	tmpDir := t.TempDir()
	r := NewTodoFileRepo(tmpDir)

	date, err := time.Parse(time.DateOnly, "2023-10-01")
	require.NoError(t, err)
	todo, err := models.NewTodosFile(date)
	require.NoError(t, err)

	r.Save(todo, true)

	tests := []struct {
		name string
		want []models.TodoFileInterface
	}{
		{
			name: "Get Todo Entries",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.TodoEntries()
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		})
	}
}
