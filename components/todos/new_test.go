package todos

import (
	"testing"

	"github.com/hirotoni/memov2/config"
	"github.com/stretchr/testify/require"
)

func TestGenerateTodoFile(t *testing.T) {
	c, err := config.LoadTomlConfig()
	require.NoError(t, err)

	// Test with truncate = false
	err = GenerateTodoFile(*c, false)
	if err != nil {
		t.Errorf("Error generating todo file with truncate=false: %v", err)
	}
}
