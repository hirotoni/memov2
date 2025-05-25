package platform

import (
	"testing"

	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestNewEditor(t *testing.T) {
	// Execute
	editor := NewEditor()

	// Assert
	assert.NotNil(t, editor)
	assert.Implements(t, (*interfaces.Editor)(nil), editor)
}

func TestDefaultEditor_Open_Integration(t *testing.T) {
	// Skip this test in CI/CD or when editor is not available
	// This is an integration test that would actually try to open VS Code
	t.Skip("Skipping integration test for editor opening")

	// Setup
	editor := DefaultEditor{}
	basedir := "/tmp"
	path := "/tmp/test.txt"

	// Execute
	err := editor.Open(basedir, path)

	// Assert - would depend on VS Code being installed
	// In a real scenario, you'd need to mock exec.Command
	_ = err
}

func TestDefaultEditor_Type(t *testing.T) {
	// Setup
	editor := DefaultEditor{}

	// Assert - verify it implements the interface
	var _ interfaces.Editor = editor
}

