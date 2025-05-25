package utils

import (
	"bytes"
	"strings"

	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/utils/mygoldmark"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

// MarkdownHandler is responsible for parsing, manipulating, and rendering markdown content
type MarkdownHandler struct {
	md goldmark.Markdown
}

// NewMarkdownHandler creates a new instance of MarkdownHandler
func NewMarkdownHandler() *MarkdownHandler {
	return &MarkdownHandler{
		md: goldmark.New(
			goldmark.WithRenderer(
				mygoldmark.NewMarkdownRenderer(),
			),
			goldmark.WithExtensions(
				extension.GFM,
				meta.New(meta.WithStoresInDocument()), // Enable meta data parsing and store in document
			),
		),
	}
}

type Heading struct {
	Level int    // Heading level (h1, h2, etc.)
	Text  string // Text of the heading to search for
}

var (
	HeadingTodos     = Heading{Text: "todos", Level: 2}
	HeadingWantTodos = Heading{Text: "wanttodos", Level: 2}
)

func (h *MarkdownHandler) findHeadingAndContent(doc ast.Node, source []byte, heading Heading) (ast.Node, []ast.Node) {
	const (
		modeSearching = iota
		modeExiting
	)

	var foundHeading ast.Node
	var hangingNodes []ast.Node
	var mode = modeSearching

	for c := doc.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Kind() == ast.KindHeading {
			switch mode {
			case modeSearching:
				levelMatched := c.(*ast.Heading).Level == heading.Level
				textMatched := strings.Contains(string(c.Lines().Value(source)), heading.Text)
				if levelMatched && textMatched {
					foundHeading = c
					mode = modeExiting
				}
			case modeExiting:
				if c.(*ast.Heading).Level <= heading.Level {
					return foundHeading, hangingNodes
				} else {
					hangingNodes = append(hangingNodes, c)
				}
			}
		} else {
			switch mode {
			case modeSearching:
				continue
			case modeExiting:
				hangingNodes = append(hangingNodes, c)
			}
		}
	}

	return foundHeading, hangingNodes
}

func (h *MarkdownHandler) HeadingBlockByHeading(source []byte, heading Heading) (*markdown.HeadingBlock, error) {
	reader := text.NewReader(source)
	doc := h.md.Parser().Parse(reader)

	foundHeading, hangingNodes := h.findHeadingAndContent(doc, source, heading)

	if foundHeading == nil {
		return nil, nil
	}

	contents := new(bytes.Buffer)
	for _, node := range hangingNodes {
		tmp := new(bytes.Buffer)
		err := h.md.Renderer().Render(tmp, source, node)
		if err != nil {
			return nil, err
		}
		contents.Write(tmp.Bytes())
	}

	markdownEntity := &markdown.HeadingBlock{
		HeadingText: string(foundHeading.Lines().Value(source)),
		Level:       foundHeading.(*ast.Heading).Level,
		ContentText: contents.String(),
		LineNumber:  foundHeading.Lines().At(0).Start,
	}

	return markdownEntity, nil
}

func (h *MarkdownHandler) HeadingBlocksByLevel(source []byte, level int) ([]*markdown.HeadingBlock, error) {
	reader := text.NewReader(source)
	doc := h.md.Parser().Parse(reader)

	var foundNodes []ast.Node
	for c := doc.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Kind() == ast.KindHeading {
			levelMatched := c.(*ast.Heading).Level == level
			if levelMatched {
				foundNodes = append(foundNodes, c)
			}
		}
	}

	res := make([]*markdown.HeadingBlock, len(foundNodes))
	for i, node := range foundNodes {
		target := Heading{
			Level: node.(*ast.Heading).Level,
			Text:  string(node.Lines().Value(source)),
		}

		markdownEntity, err := h.HeadingBlockByHeading(source, target)
		if err != nil {
			return nil, err
		}

		// remove starting newline if exists
		for strings.HasPrefix(markdownEntity.ContentText, "\n") {
			markdownEntity.ContentText = strings.TrimPrefix(markdownEntity.ContentText, "\n")
		}

		// Ensure content ends with newlines
		if !strings.HasSuffix(markdownEntity.ContentText, "\n") {
			markdownEntity.ContentText += "\n"
		}

		res[i] = markdownEntity
	}

	return res, nil
}

// TopLevelBodyContent returns the top level body content of the markdown document
// It returns a HeadingBlock with the heading text, level, content text, and line number
// If the document has no top level body content, it returns nil
func (h *MarkdownHandler) TopLevelBodyContent(source []byte) *markdown.HeadingBlock {
	reader := text.NewReader(source)
	doc := h.md.Parser().Parse(reader)

	firstChild := doc.FirstChild()
	if firstChild == nil {
		return nil
	}

	// Check if first child is a level 1 heading
	if firstChild.Kind() != ast.KindHeading || firstChild.(*ast.Heading).Level != 1 {
		return nil
	}

	contents := new(bytes.Buffer)
	for c := firstChild.NextSibling(); c != nil; c = c.NextSibling() {
		// Stop when we hit another heading of any level
		if c.Kind() == ast.KindHeading {
			break
		}

		tmp := new(bytes.Buffer)
		err := h.md.Renderer().Render(tmp, source, c)
		if err != nil {
			return nil
		}
		contents.Write(tmp.Bytes())
		contents.WriteString("\n")
	}

	markdownEntity := &markdown.HeadingBlock{
		HeadingText: string(firstChild.Lines().Value(source)),
		Level:       1,
		ContentText: contents.String(),
		LineNumber:  firstChild.Lines().At(0).Start,
	}

	return markdownEntity
}

func (h *MarkdownHandler) Metadata(source []byte) map[string]interface{} {
	reader := text.NewReader(source)
	doc := h.md.Parser().Parse(reader)
	return doc.OwnerDocument().Meta()
}
