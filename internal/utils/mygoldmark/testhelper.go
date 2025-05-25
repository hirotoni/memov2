package mygoldmark

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

// safeGroundParent returns the parent of the parent of the node ensuring nil safe.
func safeGroundParent(n ast.Node) ast.Node {
	if n.Parent() == nil {
		return nil
	}
	return n.Parent().Parent()
}

func genHeaderNode(level int, setBlankSpacePreviousLines bool, isFirstNode, isLastNode bool) ast.Node {
	h := ast.NewHeading(level)
	h.SetBlankPreviousLines(setBlankSpacePreviousLines)
	if !isFirstNode {
		h.SetPreviousSibling(ast.NewHeading(level))
	}
	if !isLastNode {
		h.SetNextSibling(ast.NewHeading(level))
	}
	return h
}

func genTextNode(text []byte, setSoftLineBreak bool, parent ast.Node) ast.Node {
	t := ast.NewText()
	t.Segment.Start = 0
	t.Segment.Stop = len(text)
	t.SetSoftLineBreak(setSoftLineBreak)

	if parent != nil {
		parent.AppendChild(parent, t)
	}
	return t
}

func genLinkNode(text, destination []byte) ast.Node {
	nl := ast.NewLink()
	nl.Destination = destination

	// segment
	t := ast.NewText()
	t.Segment.Start = 0
	t.Segment.Stop = len(text)
	nl.AppendChild(nl, t)

	return nl
}

func genAutoLinkNode(text []byte) ast.Node {
	// segment
	t := ast.NewText()
	t.Segment.Start = 0
	t.Segment.Stop = len(text)

	al := ast.NewAutoLink(ast.AutoLinkURL, t)

	return al
}

func genTaskCheckBoxNode(checked bool) ast.Node {
	return extast.NewTaskCheckBox(checked)
}

func genEnphasisNode(level int) ast.Node {
	return ast.NewEmphasis(level)
}

func genParagraphNode(setBlankSpacePreviousLines bool) ast.Node {
	p := ast.NewParagraph()
	p.SetBlankPreviousLines(setBlankSpacePreviousLines)
	return p
}

func genListNode(marker byte, setBlankSpacePreviousLines bool) ast.Node {
	l := ast.NewList(marker)
	l.SetBlankPreviousLines(setBlankSpacePreviousLines)
	return l
}

func genDocumentNode() ast.Node {
	return ast.NewDocument()
}

func genTextBlockNode() ast.Node {
	return ast.NewTextBlock()
}

func genListItemNode(num int, marker byte, offset int) []*ast.ListItem {
	lil := make([]*ast.ListItem, num)
	for i := range lil {
		lil[i] = ast.NewListItem(offset)
		if i > 0 {
			lil[i-1].SetNextSibling(lil[i])
		}
	}

	// parents
	doc := ast.NewDocument()
	l := ast.NewList(marker)
	l.Start = 1

	doc.AppendChild(doc, l)
	for _, li := range lil {
		l.AppendChild(l, li)
	}
	return lil
}

func genNestedListItemNode(num int, marker byte, offset int) []*ast.ListItem {
	lil := make([]*ast.ListItem, num)
	for i := range lil {
		lil[i] = ast.NewListItem(offset)
		if i > 0 {
			lil[i-1].SetNextSibling(lil[i])
		}
	}

	// parents
	doc := ast.NewDocument()
	middlel := ast.NewList(marker)
	middleli := ast.NewListItem(offset)
	l := ast.NewList(marker)
	l.Start = 1

	doc.AppendChild(doc, middlel)
	middlel.AppendChild(middlel, middleli)
	middleli.AppendChild(middleli, l)

	for _, li := range lil {
		l.AppendChild(l, li)
	}
	return lil
}

func genThematicBreakNode() ast.Node {
	return ast.NewThematicBreak()
}

func genCodeSpan() ast.Node {
	return ast.NewCodeSpan()
}

func genCodeBlockNode(source []byte, hasBlankPreviousLines bool) ast.Node {
	n := ast.NewCodeBlock()
	if hasBlankPreviousLines {
		n.SetBlankPreviousLines(true)
	}
	lines := n.Lines()
	// Split source into lines and add them
	sourceLines := strings.Split(string(source), "\n")
	for i, line := range sourceLines {
		if i < len(sourceLines)-1 || line != "" { // Don't add empty last line
			start := 0
			for j := 0; j < i; j++ {
				start += len(sourceLines[j]) + 1 // +1 for newline
			}
			end := start + len(line)
			if i < len(sourceLines)-1 {
				end++ // Include newline
			}
			segment := text.NewSegment(start, end)
			lines.Append(segment)
		}
	}
	return n
}

func genFencedCodeBlockNode(source []byte, language []byte, hasBlankPreviousLines bool, hasNextSibling bool) ast.Node {
	n := ast.NewFencedCodeBlock(nil)
	if hasBlankPreviousLines {
		n.SetBlankPreviousLines(true)
	}
	if language != nil {
		n.Info = ast.NewTextSegment(text.NewSegment(0, len(language)))
	}

	lines := n.Lines()
	// Split source into lines and add them, skipping first line which is language
	sourceContent := string(source)
	if strings.Contains(sourceContent, "\n") {
		parts := strings.SplitN(sourceContent, "\n", 2)
		if len(parts) > 1 {
			sourceLines := strings.Split(parts[1], "\n")
			for i, line := range sourceLines {
				if i < len(sourceLines)-1 || line != "" {
					start := len(parts[0]) + 1 // +1 for first newline
					for j := 0; j < i; j++ {
						start += len(sourceLines[j]) + 1
					}
					end := start + len(line)
					if i < len(sourceLines)-1 {
						end++
					}
					segment := text.NewSegment(start, end)
					lines.Append(segment)
				}
			}
		}
	}

	if hasNextSibling {
		// Create a dummy next sibling
		dummy := ast.NewParagraph()
		n.SetNextSibling(dummy)
	}

	return n
}

func genBlockquoteNode(hasBlankPreviousLines bool) ast.Node {
	n := ast.NewBlockquote()
	if hasBlankPreviousLines {
		n.SetBlankPreviousLines(true)
	}
	return n
}

func genHTMLBlockNode(source []byte, hasBlankPreviousLines bool) ast.Node {
	n := ast.NewHTMLBlock(ast.HTMLBlockType1)
	if hasBlankPreviousLines {
		n.SetBlankPreviousLines(true)
	}
	lines := n.Lines()

	// Split source into lines and add them
	sourceLines := strings.Split(string(source), "\n")
	for i, line := range sourceLines {
		if i < len(sourceLines)-1 || line != "" {
			start := 0
			for j := 0; j < i; j++ {
				start += len(sourceLines[j]) + 1
			}
			end := start + len(line)
			if i < len(sourceLines)-1 {
				end++
			}
			segment := text.NewSegment(start, end)
			lines.Append(segment)
		}
	}
	return n
}

func genTableNode(headers []string, rows [][]string, alignments []extast.Alignment, hasBlankPreviousLines bool) ast.Node {
	table := extast.NewTable()
	table.Alignments = alignments
	if hasBlankPreviousLines {
		table.SetBlankPreviousLines(true)
	}

	// Create header row
	if len(headers) > 0 {
		headerRow := extast.NewTableRow(alignments)
		header := extast.NewTableHeader(headerRow)

		for _, headerText := range headers {
			cell := extast.NewTableCell()
			textNode := ast.NewText()
			textNode.Segment = text.NewSegment(0, len(headerText))
			cell.AppendChild(cell, textNode)
			header.AppendChild(header, cell)
		}

		table.AppendChild(table, header)
	}

	// Create data rows
	for _, row := range rows {
		tableRow := extast.NewTableRow(alignments)
		for _, cellText := range row {
			cell := extast.NewTableCell()
			textNode := ast.NewText()
			textNode.Segment = text.NewSegment(0, len(cellText))
			cell.AppendChild(cell, textNode)
			tableRow.AppendChild(tableRow, cell)
		}
		table.AppendChild(table, tableRow)
	}

	return table
}
