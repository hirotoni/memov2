package todos

import (
	"testing"

	"github.com/hirotoni/memov2/config"
)

func TestGenerateTodoFile(t *testing.T) {
	c := config.LoadTomlConfig()

	// Test with truncate = false
	err := GenerateTodoFile(*c, false)
	if err != nil {
		t.Errorf("Error generating todo file with truncate=false: %v", err)
	}
}
