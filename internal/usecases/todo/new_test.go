package todo

import (
	"testing"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/repositories"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	"github.com/stretchr/testify/require"
)

func TestGenerateTodoFile(t *testing.T) {
	c, err := config.LoadTomlConfig()
	require.NoError(t, err)

	// Create ConfigProvider from TomlConfig
	configProvider := interfaces.NewConfigProvider(c)

	// Test with truncate = false
	r := repositories.NewRepositories(configProvider)
	mockEditor := mock.NewMockEditor()
	uc := NewTodo(configProvider, r, mockEditor)

	err = uc.GenerateTodoFile(false)
	if err != nil {
		t.Errorf("Error generating todo file with truncate=false: %v", err)
	}
}
