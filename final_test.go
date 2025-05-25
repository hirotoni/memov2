package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/utils"
)

func main() {
	// Create content with table followed by code block
	originalContent := `---
category: ["test"]
---

# Test Document

## Data Table

| Name | Value |
|------|-------|
| Item1 | 100   |
| Item2 | 200   |

## Code Example

` + "```" + `go
func main() {
    fmt.Println("Hello World")
}
` + "```" + `
`

	fmt.Println("=== ORIGINAL CONTENT ===")
	fmt.Println(originalContent)

	// Parse and render back
	handler := utils.NewMarkdownHandler()

	// Get heading blocks
	headingBlocks, err := handler.HeadingBlocksByLevel([]byte(originalContent), 2)
	if err != nil {
		panic(err)
	}

	// Create memo file
	memo, err := domain.NewMemoFile(
		time.Now(),
		"Test Document",
		[]string{"test"},
	)
	if err != nil {
		panic(err)
	}

	memo.SetHeadingBlocks(headingBlocks)

	// Get rendered content
	renderedContent := memo.ContentString()

	fmt.Println("\n=== RENDERED CONTENT ===")
	fmt.Println(renderedContent)

	fmt.Println("\n=== COMPARISON ===")
	originalLines := strings.Split(originalContent, "\n")
	renderedLines := strings.Split(renderedContent, "\n")

	// Focus on the area around the table and code block
	fmt.Println("Lines around table and code block:")
	for i := 8; i < 16 && i < len(originalLines) && i < len(renderedLines); i++ {
		orig := originalLines[i]
		rend := renderedLines[i]
		status := "✓"
		if orig != rend {
			status = "✗"
		}
		fmt.Printf("Line %d: %s orig=%q rend=%q\n", i+1, status, orig, rend)
	}
}
