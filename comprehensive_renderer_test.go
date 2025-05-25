package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/utils"
)

func main() {
	// Create comprehensive test content with various markdown elements
	originalContent := `---
category: ["test", "comprehensive"]
---

# Main Heading

This is a paragraph with **bold text** and *italic text*.

## Lists Section

### Unordered List
- First item
- Second item
  - Nested item
  - Another nested item
- Third item

### Ordered List
1. First numbered item
2. Second numbered item
3. Third numbered item

## Code Section

### Inline Code
Here is some `inline code` in a sentence.

### Code Block
` + "```" + `go
func main() {
    fmt.Println("Hello, World!")
    return 0
}
` + "```" + `

### Fenced Code Block with Language
` + "```" + `python
def hello():
    print("Hello from Python!")
    return True
` + "```" + `

## Table Section

| Name | Age | City |
|------|-----|------|
| Alice | 25 | New York |
| Bob | 30 | London |
| Charlie | 35 | Tokyo |

## Links and References

[Google](https://google.com)

[Internal Link](#main-heading)

## Blockquotes

> This is a blockquote.
> 
> It can span multiple lines.

## Horizontal Rule

---

## Task Lists

- [x] Completed task
- [ ] Incomplete task
- [x] Another completed task

## Emphasis and Styling

**Bold text** and *italic text* and ***bold italic***

` + "`" + `code span` + "`" + ` and [link text](https://example.com)

## Nested Lists

1. First level
   - Second level
     - Third level
       - Fourth level
2. Back to first level

## Mixed Content

This paragraph has **bold** and *italic* text with `inline code`.

> Blockquote with **bold** and `code`

- List item with **bold** text
- List item with `code` text
- List item with [link](https://example.com)
`

	fmt.Println("=== ORIGINAL CONTENT ===")
	fmt.Println(originalContent)

	// Parse and render back
	handler := utils.NewMarkdownHandler()
	
	// Get top level content
	topLevel := handler.TopLevelBodyContent([]byte(originalContent))
	
	// Get heading blocks
	headingBlocks, err := handler.HeadingBlocksByLevel([]byte(originalContent), 2)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n=== TOP LEVEL CONTENT ANALYSIS ===")
	if topLevel != nil {
		fmt.Printf("Top Level Content Length: %d characters\n", len(topLevel.ContentText))
		fmt.Printf("Top Level Content Preview: %q\n", topLevel.ContentText[:min(100, len(topLevel.ContentText))])
	}

	fmt.Println("\n=== HEADING BLOCKS ANALYSIS ===")
	for i, hb := range headingBlocks {
		fmt.Printf("Heading Block %d: %q\n", i+1, hb.HeadingText)
		fmt.Printf("  Content Length: %d characters\n", len(hb.ContentText))
		fmt.Printf("  Content Preview: %q\n", hb.ContentText[:min(100, len(hb.ContentText))])
	}

	// Create memo file
	memo, err := domain.NewMemoFile(
		time.Now(),
		"Main Heading",
		[]string{"test", "comprehensive"},
	)
	if err != nil {
		panic(err)
	}

	if topLevel != nil {
		memo.SetTopLevelBodyContent(topLevel)
	}
	memo.SetHeadingBlocks(headingBlocks)

	// Get rendered content
	renderedContent := memo.ContentString()

	fmt.Println("\n=== RENDERED CONTENT ===")
	fmt.Println(renderedContent)

	fmt.Println("\n=== COMPARISON ANALYSIS ===")
	originalLines := strings.Split(originalContent, "\n")
	renderedLines := strings.Split(renderedContent, "\n")
	
	fmt.Printf("Original lines: %d\n", len(originalLines))
	fmt.Printf("Rendered lines: %d\n", len(renderedLines))
	
	// Check for specific elements
	checkElement(originalLines, renderedLines, "**bold text**", "Bold text")
	checkElement(originalLines, renderedLines, "*italic text*", "Italic text")
	checkElement(originalLines, renderedLines, "`inline code`", "Inline code")
	checkElement(originalLines, renderedLines, "```go", "Go code block")
	checkElement(originalLines, renderedLines, "```python", "Python code block")
	checkElement(originalLines, renderedLines, "| Name | Age |", "Table")
	checkElement(originalLines, renderedLines, "[Google]", "Link")
	checkElement(originalLines, renderedLines, "> This is a blockquote", "Blockquote")
	checkElement(originalLines, renderedLines, "---", "Horizontal rule")
	checkElement(originalLines, renderedLines, "- [x] Completed task", "Task list")
	checkElement(originalLines, renderedLines, "1. First numbered item", "Ordered list")
	checkElement(originalLines, renderedLines, "- First item", "Unordered list")
}

func checkElement(original, rendered []string, searchText, elementName string) {
	origFound := false
	rendFound := false
	
	for _, line := range original {
		if strings.Contains(line, searchText) {
			origFound = true
			break
		}
	}
	
	for _, line := range rendered {
		if strings.Contains(line, searchText) {
			rendFound = true
			break
		}
	}
	
	status := "✓"
	if origFound != rendFound {
		status = "✗"
	}
	
	fmt.Printf("%s %s: orig=%t rend=%t\n", status, elementName, origFound, rendFound)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
