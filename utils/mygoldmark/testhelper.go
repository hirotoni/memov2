package mygoldmark

import (
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
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
