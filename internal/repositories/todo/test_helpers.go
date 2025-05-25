package todo

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hirotoni/memov2/internal/domain/markdown"
)

// createTestHeadingBlock creates a heading block for testing
func createTestHeadingBlock(level int, heading, content string) *markdown.HeadingBlock {
	return &markdown.HeadingBlock{
		Level:       level,
		HeadingText: heading,
		ContentText: content,
	}
}

// readGoldenFile reads the expected content from a golden file
func readGoldenFile(t *testing.T, testName string) string {
	t.Helper()

	goldenPath := filepath.Join("testdata", testName+".golden")
	content, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file %s: %v", goldenPath, err)
	}

	return string(content)
}
