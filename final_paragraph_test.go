package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/utils"
)

func main() {
	// Create content with multiple paragraphs
	originalContent := `---
category: ["test"]
---

# Test Document

This is the first paragraph.

This is the second paragraph.

This is the third paragraph.
`

	fmt.Println("=== ORIGINAL CONTENT ===")
	fmt.Println(originalContent)

	// Parse and render back
	handler := utils.NewMarkdownHandler()

	// Get top level content
	topLevel := handler.TopLevelBodyContent([]byte(originalContent))

	// Create memo file
	memo, err := domain.NewMemoFile(
		time.Now(),
		"Test Document",
		[]string{"test"},
	)
	if err != nil {
		panic(err)
	}

	if topLevel != nil {
		memo.SetTopLevelBodyContent(topLevel)
	}

	// Get rendered content
	renderedContent := memo.ContentString()

	fmt.Println("=== RENDERED CONTENT ===")
	fmt.Println(renderedContent)

	fmt.Println("\n=== COMPARISON ===")
	originalLines := strings.Split(originalContent, "\n")
	renderedLines := strings.Split(renderedContent, "\n")

	// Focus on the area around the paragraphs
	fmt.Println("Lines around paragraphs:")
	for i := 4; i < 12 && i < len(originalLines) && i < len(renderedLines); i++ {
		orig := originalLines[i]
		rend := renderedLines[i]
		status := "✓"
		if orig != rend {
			status = "✗"
		}
		fmt.Printf("Line %d: %s orig=%q rend=%q\n", i+1, status, orig, rend)
	}
}
