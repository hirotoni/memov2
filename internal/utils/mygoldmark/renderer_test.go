package mygoldmark

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

func TestMarkdownRenderer_renderHeading(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering, level 1",
			args:  args{node: genHeaderNode(1, false, false, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "# ", err: false},
		},
		{
			name:  "exiting, level 1",
			args:  args{node: genHeaderNode(1, false, false, false), entering: false},
			wants: wants{status: ast.WalkContinue, str: "\n\n", err: false},
		},
		{
			name:  "entering, level 6",
			args:  args{node: genHeaderNode(6, false, false, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "###### ", err: false},
		},
		{
			name:  "entering, blank previous lines",
			args:  args{node: genHeaderNode(1, true, false, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "\n\n# ", err: false},
		},
		{
			name:  "entering, blank previous lines, first node",
			args:  args{node: genHeaderNode(1, true, true, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "# ", err: false},
		},
		{
			name:  "exiting, blank previous lines",
			args:  args{node: genHeaderNode(1, true, false, false), entering: false},
			wants: wants{status: ast.WalkContinue, str: "\n\n", err: false},
		},
		{
			name:  "exiting, blank previous lines, last node",
			args:  args{node: genHeaderNode(1, true, false, true), entering: false},
			wants: wants{status: ast.WalkContinue, str: "\n\n", err: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			source := []byte("dummy source")

			buf := new(strings.Builder)
			w := bufio.NewWriter(buf)

			got, err := r.renderHeading(w, source, tt.args.node, tt.args.entering)

			if tt.wants.err {
				assert.NotNil(err)
			}

			assert.Nil(err)
			assert.Equal(tt.wants.status, got)

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, buf.String())
		})
	}
}

func TestMarkdownRenderer_renderEmphasis(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering true",
			args:  args{node: genEnphasisNode(1), entering: true},
			wants: wants{status: ast.WalkContinue, str: "*", err: false},
		},
		{
			name:  "entering false",
			args:  args{node: genEnphasisNode(1), entering: false},
			wants: wants{status: ast.WalkContinue, str: "*", err: false},
		},
		{
			name:  "entering true, level 2",
			args:  args{node: genEnphasisNode(2), entering: true},
			wants: wants{status: ast.WalkContinue, str: "**", err: false},
		},
		{
			name:  "entering false, level 2",
			args:  args{node: genEnphasisNode(2), entering: false},
			wants: wants{status: ast.WalkContinue, str: "**", err: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)
			source := []byte("dummy source")

			got, err := r.renderEmphasis(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderEmphasis() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderEmphasis() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderTaskCheckBox(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "isChecked true, entering true",
			args:  args{node: genTaskCheckBoxNode(true), entering: true},
			wants: wants{status: ast.WalkContinue, str: "[x] ", err: false},
		},
		{
			name:  "isChecked false, entering true",
			args:  args{node: genTaskCheckBoxNode(false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "[ ] ", err: false},
		},
		{
			name:  "isChecked true, entering false",
			args:  args{node: genTaskCheckBoxNode(true), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
		{
			name:  "isChecked false, entering false",
			args:  args{node: genTaskCheckBoxNode(false), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)
			dummySource := []byte("dummy source")

			got, err := r.renderTaskCheckBox(w, dummySource, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderTaskCheckBox() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderTaskCheckBox() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderLink(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("text")

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering true",
			args:  args{node: genLinkNode([]byte("text"), []byte("destination")), entering: true},
			wants: wants{status: ast.WalkContinue, str: "[text](destination)", err: false},
		},
		{
			name:  "entering false",
			args:  args{node: genLinkNode([]byte("text"), []byte("destination")), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderLink(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderLink() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderLink() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderAutoLink(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("https://example.com")

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering true",
			args:  args{node: genAutoLinkNode(source), entering: true},
			wants: wants{status: ast.WalkContinue, str: "https://example.com", err: false},
		},
		{
			name:  "entering false",
			args:  args{node: genAutoLinkNode(source), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderAutoLink(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderAutoLink() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderAutoLink() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderText(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("test text")

	offset := 2
	li := ast.NewListItem(offset)
	tb := ast.NewTextBlock()
	tb.AppendChild(li, tb)

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering",
			args:  args{node: genTextNode(source, false, ast.NewTextBlock()), entering: true},
			wants: wants{status: ast.WalkContinue, str: string(source), err: false},
		},
		{
			name:  "entering, soft line break",
			args:  args{node: genTextNode(source, true, ast.NewTextBlock()), entering: true},
			wants: wants{status: ast.WalkContinue, str: string(source) + "\n", err: false},
		},
		{
			name:  "entering, parent parent listitem",
			args:  args{node: genTextNode(source, true, tb), entering: true},
			wants: wants{status: ast.WalkContinue, str: string(source) + "\n" + strings.Repeat(" ", offset), err: false},
		},
		{
			name:  "entering, parent link",
			args:  args{node: genTextNode(source, true, ast.NewLink()), entering: true},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
		{
			name:  "entering, parent nil",
			args:  args{node: genTextNode(source, true, nil), entering: true},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
		{
			name:  "entering false",
			args:  args{node: genTextNode(source, false, ast.NewTextBlock()), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderText(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderText() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderText() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderParagraph(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("test text")

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering",
			args:  args{node: genParagraphNode(false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
		{
			name:  "entering, blank previous lines",
			args:  args{node: genParagraphNode(true), entering: true},
			wants: wants{status: ast.WalkContinue, str: "\n", err: false},
		},
		{
			name:  "exiting",
			args:  args{node: genParagraphNode(false), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderParagraph(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderParagraph() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderParagraph() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderList(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("test text")

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering",
			args:  args{node: genListNode(byte('.'), false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
		{
			name:  "entering, blank previous lines",
			args:  args{node: genListNode(byte('.'), true), entering: true},
			wants: wants{status: ast.WalkContinue, str: "\n\n", err: false},
		},
		{
			name:  "exiting",
			args:  args{node: genListNode(byte('.'), false), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderList(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderList() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderList() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderDocument(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("test text")

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering",
			args:  args{node: genDocumentNode(), entering: true},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
		{
			name:  "exiting",
			args:  args{node: genDocumentNode(), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderDocument(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderDocument() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderDocument() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderTextBlock(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("test text")

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering",
			args:  args{node: genTextBlockNode(), entering: true},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
		{
			name:  "entering false",
			args:  args{node: genTextBlockNode(), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderTextBlock(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderTextBlock() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderTextBlock() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderListItem(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("test text")

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering",
			args:  args{node: genListItemNode(1, byte('-'), 2)[0], entering: true},
			wants: wants{status: ast.WalkContinue, str: "- ", err: false},
		},
		{
			name:  "entering, ordered listitem",
			args:  args{node: genListItemNode(1, byte('.'), 3)[0], entering: true},
			wants: wants{status: ast.WalkContinue, str: "1. ", err: false},
		},
		{
			name:  "entering, second ordered listitem",
			args:  args{node: genListItemNode(2, byte('.'), 3)[1], entering: true},
			wants: wants{status: ast.WalkContinue, str: "\n2. ", err: false},
		},
		{
			name:  "entering, nested ordered listitem",
			args:  args{node: genNestedListItemNode(1, byte('.'), 3)[0], entering: true},
			wants: wants{status: ast.WalkContinue, str: "\n   1. ", err: false},
		},
		{
			name:  "entering, second nested ordered listitem",
			args:  args{node: genNestedListItemNode(2, byte('.'), 3)[1], entering: true},
			wants: wants{status: ast.WalkContinue, str: "\n   2. ", err: false},
		},
		{
			name:  "entering false",
			args:  args{node: genListItemNode(1, byte('-'), 2)[0], entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderListItem(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderListItem() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderListItem() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderThematicBreak(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("test text")

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering",
			args:  args{node: genThematicBreakNode(), entering: true},
			wants: wants{status: ast.WalkContinue, str: "\n\n---", err: false},
		},
		{
			name:  "entering false",
			args:  args{node: genThematicBreakNode(), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := &Renderer{}
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderThematicBreak(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderThematicBreak() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderThematicBreak() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderCodeSpan(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering true",
			args:  args{node: genCodeSpan(), entering: true},
			wants: wants{status: ast.WalkContinue, str: "`", err: false},
		},
		{
			name:  "entering false",
			args:  args{node: genCodeSpan(), entering: false},
			wants: wants{status: ast.WalkContinue, str: "`", err: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := &Renderer{}
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)
			source := []byte("test text")

			got, err := r.renderCodeSpan(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderCodeSpan() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderCodeSpan() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderBlockquote(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("test text")

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering",
			args:  args{node: genBlockquoteNode(false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "> ", err: false},
		},
		{
			name:  "entering, blank previous lines",
			args:  args{node: genBlockquoteNode(true), entering: true},
			wants: wants{status: ast.WalkContinue, str: "\n\n> ", err: false},
		},
		{
			name:  "exiting",
			args:  args{node: genBlockquoteNode(false), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := &Renderer{}
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderBlockquote(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderBlockquote() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderBlockquote() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderHTMLBlock(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("<div>test content</div>\n")

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering",
			args:  args{node: genHTMLBlockNode(source, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "<div>test content</div>\n", err: false},
		},
		{
			name:  "entering, blank previous lines",
			args:  args{node: genHTMLBlockNode(source, true), entering: true},
			wants: wants{status: ast.WalkContinue, str: "\n\n<div>test content</div>\n", err: false},
		},
		{
			name:  "exiting",
			args:  args{node: genHTMLBlockNode(source, false), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := &Renderer{}
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderHTMLBlock(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderHTMLBlock() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderHTMLBlock() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderCodeBlock(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("code line 1\ncode line 2\n")

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering",
			args:  args{node: genCodeBlockNode(source, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "    code line 1\n    code line 2\n", err: false},
		},
		{
			name:  "entering, blank previous lines",
			args:  args{node: genCodeBlockNode(source, true), entering: true},
			wants: wants{status: ast.WalkContinue, str: "\n\n    code line 1\n    code line 2\n", err: false},
		},
		{
			name:  "exiting",
			args:  args{node: genCodeBlockNode(source, false), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := &Renderer{}
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderCodeBlock(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderCodeBlock() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderCodeBlock() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderFencedCodeBlock(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	source := []byte("go\nfmt.Println(\"hello\")\n")

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering",
			args:  args{node: genFencedCodeBlockNode(source, []byte("go"), false, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "```go\nfmt.Println(\"hello\")\n```", err: false},
		},
		{
			name:  "entering, blank previous lines",
			args:  args{node: genFencedCodeBlockNode(source, []byte("go"), true, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "\n```go\nfmt.Println(\"hello\")\n```", err: false},
		},
		{
			name:  "entering, with next sibling",
			args:  args{node: genFencedCodeBlockNode(source, []byte("go"), false, true), entering: true},
			wants: wants{status: ast.WalkContinue, str: "```go\nfmt.Println(\"hello\")\n```\n", err: false},
		},
		{
			name:  "entering, no language info",
			args:  args{node: genFencedCodeBlockNode(source, nil, false, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "```\nfmt.Println(\"hello\")\n```", err: false},
		},
		{
			name:  "exiting",
			args:  args{node: genFencedCodeBlockNode(source, []byte("go"), false, false), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := &Renderer{}
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)

			got, err := r.renderFencedCodeBlock(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderFencedCodeBlock() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderFencedCodeBlock() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderTable(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering, no blank previous lines",
			args:  args{node: genTableNode([]string{"Header1", "Header2"}, [][]string{{"Cell1", "Cell2"}}, []extast.Alignment{extast.AlignNone, extast.AlignNone}, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
		{
			name:  "entering, blank previous lines",
			args:  args{node: genTableNode([]string{"Header1", "Header2"}, [][]string{{"Cell1", "Cell2"}}, []extast.Alignment{extast.AlignNone, extast.AlignNone}, true), entering: true},
			wants: wants{status: ast.WalkContinue, str: "\n\n", err: false},
		},
		{
			name:  "exiting",
			args:  args{node: genTableNode([]string{"Header1", "Header2"}, [][]string{{"Cell1", "Cell2"}}, []extast.Alignment{extast.AlignNone, extast.AlignNone}, false), entering: false},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)
			source := []byte("test text")

			got, err := r.renderTable(w, source, tt.args.node, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderTable() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderTable() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderTableHeader(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering, basic table",
			args:  args{node: genTableNode([]string{"Name", "Age"}, [][]string{{"Alice", "30"}}, []extast.Alignment{extast.AlignNone, extast.AlignNone}, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "| ", err: false},
		},
		{
			name:  "entering, with alignments",
			args:  args{node: genTableNode([]string{"Left", "Center", "Right"}, [][]string{{"A", "B", "C"}}, []extast.Alignment{extast.AlignLeft, extast.AlignCenter, extast.AlignRight}, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "| ", err: false},
		},
		{
			name:  "entering, single column",
			args:  args{node: genTableNode([]string{"Single"}, [][]string{{"Value"}}, []extast.Alignment{extast.AlignNone}, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "| ", err: false},
		},
		{
			name:  "exiting",
			args:  args{node: genTableNode([]string{"Name", "Age"}, [][]string{{"Alice", "30"}}, []extast.Alignment{extast.AlignNone, extast.AlignNone}, false), entering: false},
			wants: wants{status: ast.WalkContinue, str: " |\n| --- | --- |\n", err: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)
			source := []byte("test text")

			// Get the header node from the table
			table := tt.args.node.(*extast.Table)
			header := table.FirstChild().(*extast.TableHeader)

			got, err := r.renderTableHeader(w, source, header, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderTableHeader() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderTableHeader() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderTableRow(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering, basic row",
			args:  args{node: genTableNode([]string{"Name", "Age"}, [][]string{{"Alice", "30"}}, []extast.Alignment{extast.AlignNone, extast.AlignNone}, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "| ", err: false},
		},
		{
			name:  "entering, three columns",
			args:  args{node: genTableNode([]string{"A", "B", "C"}, [][]string{{"1", "2", "3"}}, []extast.Alignment{extast.AlignNone, extast.AlignNone, extast.AlignNone}, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "| ", err: false},
		},
		{
			name:  "entering, single column row",
			args:  args{node: genTableNode([]string{"Single"}, [][]string{{"Value"}}, []extast.Alignment{extast.AlignNone}, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "| ", err: false},
		},
		{
			name:  "exiting",
			args:  args{node: genTableNode([]string{"Name", "Age"}, [][]string{{"Alice", "30"}}, []extast.Alignment{extast.AlignNone, extast.AlignNone}, false), entering: false},
			wants: wants{status: ast.WalkContinue, str: " |\n", err: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)
			source := []byte("test text")

			// Get the first data row from the table
			table := tt.args.node.(*extast.Table)
			var dataRow *extast.TableRow
			for child := table.FirstChild(); child != nil; child = child.NextSibling() {
				if child.Kind() == extast.KindTableRow {
					dataRow = child.(*extast.TableRow)
					break
				}
			}

			got, err := r.renderTableRow(w, source, dataRow, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderTableRow() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderTableRow() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}

func TestMarkdownRenderer_renderTableCell(t *testing.T) {
	type args struct {
		node     ast.Node
		entering bool
	}
	type wants struct {
		status ast.WalkStatus
		str    string
		err    bool
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "entering",
			args:  args{node: genTableNode([]string{"Name", "Age"}, [][]string{{"Alice", "30"}}, []extast.Alignment{extast.AlignNone, extast.AlignNone}, false), entering: true},
			wants: wants{status: ast.WalkContinue, str: "", err: false},
		},
		{
			name:  "exiting",
			args:  args{node: genTableNode([]string{"Name", "Age"}, [][]string{{"Alice", "30"}}, []extast.Alignment{extast.AlignNone, extast.AlignNone}, false), entering: false},
			wants: wants{status: ast.WalkContinue, str: " | ", err: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			r := NewRenderer()
			sb := new(strings.Builder)
			w := bufio.NewWriter(sb)
			source := []byte("test text")

			// Get the first cell from the first data row
			table := tt.args.node.(*extast.Table)
			var cell *extast.TableCell
			for child := table.FirstChild(); child != nil; child = child.NextSibling() {
				if child.Kind() == extast.KindTableRow {
					row := child.(*extast.TableRow)
					cell = row.FirstChild().(*extast.TableCell)
					break
				}
			}

			got, err := r.renderTableCell(w, source, cell, tt.args.entering)
			if (err != nil) != tt.wants.err {
				t.Errorf("MarkdownRenderer.renderTableCell() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if !reflect.DeepEqual(got, tt.wants.status) {
				t.Errorf("MarkdownRenderer.renderTableCell() = %v, want %v", got, tt.wants.status)
			}

			assert.NoError(w.Flush())
			assert.Equal(tt.wants.str, sb.String())
		})
	}
}
