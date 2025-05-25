package utils

import (
	"bytes"
	"strings"

	"github.com/hirotoni/memov2/models"
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

func (h *MarkdownHandler) findHeadingEntity(source []byte, heading Heading) (*models.HeadingBlock, error) {
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

	markdownEntity := &models.HeadingBlock{
		HeadingText: string(foundHeading.Lines().Value(source)),
		Level:       foundHeading.(*ast.Heading).Level,
		ContentText: contents.String(),
	}

	return markdownEntity, nil
}

func (h *MarkdownHandler) HeadingBlocksByLevel(source []byte, level int) ([]*models.HeadingBlock, error) {
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

	res := make([]*models.HeadingBlock, len(foundNodes))
	for i, node := range foundNodes {
		target := Heading{
			Level: node.(*ast.Heading).Level,
			Text:  string(node.Lines().Value(source)),
		}

		markdownEntity, err := h.findHeadingEntity(source, target)
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

func (h *MarkdownHandler) Metadata(source []byte) map[string]interface{} {
	reader := text.NewReader(source)
	doc := h.md.Parser().Parse(reader)
	return doc.OwnerDocument().Meta()
}
