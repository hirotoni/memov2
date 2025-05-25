package todo

import (
	"testing"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestGenerateTodoFile(t *testing.T) {
	c, err := config.LoadTomlConfig()
	require.NoError(t, err)

	// Test with truncate = false
	r := repository.NewRepositories(*c)
	uc := NewTodo(*c, r)

	err = uc.GenerateTodoFile(false)
	if err != nil {
		t.Errorf("Error generating todo file with truncate=false: %v", err)
	}
}
