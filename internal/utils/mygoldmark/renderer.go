package mygoldmark

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func NewMarkdownRenderer() renderer.Renderer {
	return renderer.NewRenderer(
		renderer.WithNodeRenderers(
			util.Prioritized(NewRenderer(), 1),
		),
	)
}

type Renderer struct{}

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks
	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, r.renderHTMLBlock)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindTextBlock, r.renderTextBlock)
	reg.Register(ast.KindThematicBreak, r.renderThematicBreak)
	reg.Register(extast.KindTable, r.renderTable)
	reg.Register(extast.KindTableHeader, r.renderTableHeader)
	reg.Register(extast.KindTableRow, r.renderTableRow)
	reg.Register(extast.KindTableCell, r.renderTableCell)

	// // inlines
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	// reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindLink, r.renderLink)
	// reg.Register(ast.KindRawHTML, r.renderRawHTML)
	reg.Register(ast.KindText, r.renderText)
	// reg.Register(ast.KindString, r.renderString)
	reg.Register(extast.KindTaskCheckBox, r.renderTaskCheckBox)
}

func (r *Renderer) renderDocument(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// nothing to do
	return ast.WalkContinue, nil
}

func (r *Renderer) renderHeading(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		if n.PreviousSibling() != nil && n.HasBlankPreviousLines() {
			_, _ = w.WriteString("\n\n")
		}
		_, _ = w.WriteString(strings.Repeat("#", n.Level) + " ")
	} else {
		// Always add newline after heading to separate from content
		_, _ = w.WriteString("\n\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderParagraph(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Paragraph)
	if entering {
		if n.HasBlankPreviousLines() {
			_, _ = w.WriteString("\n")
		}
	} else {
		// Add newline after paragraph if there's a next sibling
		if n.NextSibling() != nil {
			_, _ = w.WriteString("\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderText(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Text)
	if entering {
		p := n.Parent()
		if p == nil {
			return ast.WalkContinue, nil
		}
		if p.Kind() == ast.KindLink {
			// r.renderLink() renders text in advance. no rendering needed here.
			return ast.WalkContinue, nil
		}

		w.WriteString(string(n.Text(source)))

		if n.SoftLineBreak() {
			w.WriteString("\n")

			pp := p.Parent()
			if pp, ok := pp.(*ast.ListItem); ok {
				// ListItem - TextBlock - Text(SoftLineBreak)
				w.WriteString(strings.Repeat(" ", pp.Offset))
			}
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderList(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.List)
	if entering {
		if n.HasBlankPreviousLines() {
			_, _ = w.WriteString("\n\n")
		}
	} else {
		// Add newline after list if there's a next sibling
		if n.NextSibling() != nil {
			_, _ = w.WriteString("\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderListItem(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.ListItem)

	// at least there must be a parent and a grandparent
	// e.g. Document - List - ListItem
	p := n.Parent()
	if p == nil {
		return ast.WalkContinue, nil
	}
	pp := p.Parent()
	if pp == nil {
		return ast.WalkContinue, nil
	}

	if entering {
		// If it is not the first element of the list or it is a nested listitems, add a line break
		if n.PreviousSibling() != nil || pp.Kind() == ast.KindListItem {
			w.WriteString("\n")
		}

		// Increase the indent for nested lists
		curpp := pp
		for curpp != nil {
			li, ok := curpp.(*ast.ListItem)
			if ok {
				w.WriteString(strings.Repeat(" ", li.Offset))
			}
			curpp = safeGroundParent(curpp)
		}

		if p, ok := p.(*ast.List); ok {
			if p.IsOrdered() {
				order := p.Start
				for node.PreviousSibling() != nil {
					order++
					node = node.PreviousSibling()
				}
				w.WriteString(fmt.Sprintf("%d", order))
			}

			w.WriteString(string(p.Marker) + " ")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTextBlock(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// n := node.(*ast.TextBlock)
	// nothing to do
	return ast.WalkContinue, nil
}

func (r *Renderer) renderThematicBreak(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("\n\n---")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderBlockquote(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Blockquote)
	if entering {
		if n.HasBlankPreviousLines() {
			_, _ = w.WriteString("\n\n")
		}
		_, _ = w.WriteString("> ")
	} else {
		// Add newline after blockquote if there's a next sibling
		if n.NextSibling() != nil {
			_, _ = w.WriteString("\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeBlock(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.CodeBlock)
	if entering {
		if n.HasBlankPreviousLines() {
			_, _ = w.WriteString("\n\n")
		}

		// Process each line of the code block
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			_, _ = w.WriteString("    ") // 4 spaces for indentation
			_, _ = w.Write(line.Value(source))
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderFencedCodeBlock(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.FencedCodeBlock)
	if entering {
		if n.HasBlankPreviousLines() {
			_, _ = w.WriteString("\n")
		}

		// Write opening fence
		_, _ = w.WriteString("```")

		// Write language info if present
		if n.Info != nil {
			_, _ = w.Write(n.Info.Value(source))
		}
		_, _ = w.WriteString("\n")

		// Process each line of the code block
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			_, _ = w.Write(line.Value(source))
		}

		// Write closing fence
		_, _ = w.WriteString("```")

		if n.NextSibling() != nil {
			_, _ = w.WriteString("\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderHTMLBlock(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.HTMLBlock)
	if entering {
		if n.HasBlankPreviousLines() {
			_, _ = w.WriteString("\n\n")
		}

		// Process each line of the HTML block
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			_, _ = w.Write(line.Value(source))
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderEmphasis(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)
	if entering {
		w.WriteString(strings.Repeat("*", n.Level))
	} else {
		w.WriteString(strings.Repeat("*", n.Level))
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTaskCheckBox(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*extast.TaskCheckBox)
	if entering {
		if n.IsChecked {
			_, _ = w.WriteString("[x] ")
		} else {
			_, _ = w.WriteString("[ ] ")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderLink(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		// NOTE As of goldmark v1.7.1, ast.Link.Title is not set by default markdown parser, so that n.Text(source) is used here instead.
		// n.Text(source) retrieves text from n's child node (ast.Text) in advance to node-walking operation.
		w.WriteString(fmt.Sprintf("[%s](%s)", n.Text(source), n.Destination))
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderAutoLink(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.AutoLink)
	if entering {
		w.WriteString(fmt.Sprint(string(n.URL(source))))
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeSpan(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("`")
	} else {
		w.WriteString("`")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTable(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*extast.Table)
	if entering {
		if n.HasBlankPreviousLines() {
			_, _ = w.WriteString("\n\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTableHeader(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("| ")
	} else {
		_, _ = w.WriteString(" |\n")

		// Render separator row after header
		cellCount := 0
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			cellCount++
		}

		_, _ = w.WriteString("| ")
		for i := 0; i < cellCount; i++ {
			if i > 0 {
				_, _ = w.WriteString(" | ")
			}
			// Get alignment from parent table
			table := node.Parent()
			if table != nil && table.Kind() == extast.KindTable {
				tableNode := table.(*extast.Table)
				if i < len(tableNode.Alignments) {
					alignment := tableNode.Alignments[i]
					switch alignment {
					case extast.AlignLeft:
						_, _ = w.WriteString(":---")
					case extast.AlignCenter:
						_, _ = w.WriteString(":---:")
					case extast.AlignRight:
						_, _ = w.WriteString("---:")
					default:
						_, _ = w.WriteString("---")
					}
				} else {
					_, _ = w.WriteString("---")
				}
			} else {
				_, _ = w.WriteString("---")
			}
		}
		_, _ = w.WriteString(" |\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTableRow(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("| ")
	} else {
		_, _ = w.WriteString(" |\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTableCell(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		// Add cell separator after each cell (except the last one in the row)
		// Check if this is the last cell in the row
		if node.NextSibling() != nil {
			_, _ = w.WriteString(" | ")
		}
	}
	return ast.WalkContinue, nil
}
