package todo

import (
	"log/slog"
	"os"
	"testing"

	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/repositories"
	"github.com/hirotoni/memov2/internal/repositories/mock"
	"github.com/stretchr/testify/require"
)

func TestGenerateTodoFile(t *testing.T) {
	c, err := toml.LoadConfig()
	require.NoError(t, err)

	// Create ConfigProvider from Config
	configProvider := toml.NewProvider(c)

	// Test with truncate = false
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	r := repositories.NewRepositories(configProvider, logger)
	mockEditor := mock.NewMockEditor()
	uc := NewTodo(configProvider, r, mockEditor, logger)

	err = uc.GenerateTodoFile(false)
	if err != nil {
		t.Errorf("Error generating todo file with truncate=false: %v", err)
	}
}
