package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRewriteMemoTitle(t *testing.T) {
	tests := []struct {
		name    string
		content string
		title   string
		want    string
	}{
		{
			name:    "replaces ATX H1 and keeps everything else",
			content: "---\ncategory: [\"work\"]\ntags: [a, b]\ncreated: 2026-01-01\n---\n\n# Orig\n\nintro line\n\n### sub before h2\n\nsubtext\n\n## H2\n\n#### deep\n",
			title:   "Renamed",
			want:    "---\ncategory: [\"work\"]\ntags: [a, b]\ncreated: 2026-01-01\n---\n\n# Renamed\n\nintro line\n\n### sub before h2\n\nsubtext\n\n## H2\n\n#### deep\n",
		},
		{
			name:    "keeps missing trailing newline",
			content: "---\ncategory: []\n---\n\n# Orig",
			title:   "New",
			want:    "---\ncategory: []\n---\n\n# New",
		},
		{
			name:    "keeps CRLF line endings",
			content: "---\r\ncategory: []\r\n---\r\n\r\n# Orig\r\nbody\r\n",
			title:   "New",
			want:    "---\r\ncategory: []\r\n---\r\n\r\n# New\r\nbody\r\n",
		},
		{
			name:    "replaces setext H1",
			content: "---\ncategory: []\n---\n\nOrig\n====\n\nbody\n",
			title:   "New",
			want:    "---\ncategory: []\n---\n\n# New\n\nbody\n",
		},
		{
			name:    "replaces multi-line setext H1",
			content: "Orig first\nOrig second\n===\n\nbody\n",
			title:   "New",
			want:    "# New\n\nbody\n",
		},
		{
			name:    "frontmatter after BOM",
			content: "\ufeff---\ncategory: []\n---\n\n# Orig\n",
			title:   "New",
			want:    "\ufeff---\ncategory: []\n---\n\n# New\n",
		},
		{
			name:    "no frontmatter",
			content: "# Orig\n\nbody\n",
			title:   "New",
			want:    "# New\n\nbody\n",
		},
		{
			name:    "only first H1 is replaced",
			content: "# Orig\n\n```\n# not a heading\n```\n\n# second\n",
			title:   "New",
			want:    "# New\n\n```\n# not a heading\n```\n\n# second\n",
		},
		{
			name:    "inserts H1 when missing",
			content: "---\ncategory: []\n---\n## H2\n",
			title:   "New",
			want:    "---\ncategory: []\n---\n# New\n\n## H2\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, RewriteMemoTitle(tt.content, tt.title))
		})
	}
}

func TestRewriteMemoCategory(t *testing.T) {
	tests := []struct {
		name    string
		content string
		tree    []string
		want    string
	}{
		{
			name:    "replaces flow-style category and keeps other keys",
			content: "---\ncategory: [\"old\"]\ntags: [a, b]\ncreated: 2026-01-01\n---\n\n# T\n\n### sub\n\ntext\n",
			tree:    []string{"new", "sub"},
			want:    "---\ncategory: [\"new\", \"sub\"]\ntags: [a, b]\ncreated: 2026-01-01\n---\n\n# T\n\n### sub\n\ntext\n",
		},
		{
			name:    "replaces block-style category",
			content: "---\ntags: x\ncategory:\n  - old\n  - deeper\ncreated: 2026-01-01\n---\n\n# T\n",
			tree:    []string{"new"},
			want:    "---\ntags: x\ncategory: [\"new\"]\ncreated: 2026-01-01\n---\n\n# T\n",
		},
		{
			name:    "replaces zero-indented block sequence",
			content: "---\ncategory:\n- old\ntags: x\n---\n\n# T\n",
			tree:    []string{},
			want:    "---\ncategory: []\ntags: x\n---\n\n# T\n",
		},
		{
			name:    "replaces multi-line flow list closed at column 0",
			content: "---\ncategory: [\n  \"old\",\n  \"deeper\"\n]\ntags: x\n---\n\n# T\n",
			tree:    []string{"new"},
			want:    "---\ncategory: [\"new\"]\ntags: x\n---\n\n# T\n",
		},
		{
			name:    "replaces block value containing blank and comment lines",
			content: "---\ncategory:\n  - old\n\n# note\n  - deeper\n\n# kept comment\ntags: x\n---\n\n# T\n",
			tree:    []string{"new"},
			want:    "---\ncategory: [\"new\"]\n\n# kept comment\ntags: x\n---\n\n# T\n",
		},
		{
			name:    "block value with CRLF",
			content: "---\r\ncategory:\r\n- old\r\n-\r\ntags: x\r\n---\r\n\r\n# T\r\n",
			tree:    []string{"new"},
			want:    "---\r\ncategory: [\"new\"]\r\ntags: x\r\n---\r\n\r\n# T\r\n",
		},
		{
			name:    "replaces quoted category key",
			content: "---\n\"category\": [\"old\"]\ntags: x\n---\n\n# T\n",
			tree:    []string{"new"},
			want:    "---\ncategory: [\"new\"]\ntags: x\n---\n\n# T\n",
		},
		{
			name:    "frontmatter after BOM",
			content: "\ufeff---\ncategory: [\"old\"]\n---\n\n# T\n",
			tree:    []string{"new"},
			want:    "\ufeff---\ncategory: [\"new\"]\n---\n\n# T\n",
		},
		{
			name:    "adds category when key missing",
			content: "---\ntags: x\n---\n\n# T\n",
			tree:    []string{"a"},
			want:    "---\ncategory: [\"a\"]\ntags: x\n---\n\n# T\n",
		},
		{
			name:    "prepends frontmatter when missing",
			content: "# T\n",
			tree:    []string{"a"},
			want:    "---\ncategory: [\"a\"]\n---\n\n# T\n",
		},
		{
			name:    "does not touch category-like text in body",
			content: "---\ncategory: [\"old\"]\n---\n\n# T\n\ncategory: keep\n",
			tree:    []string{"new"},
			want:    "---\ncategory: [\"new\"]\n---\n\n# T\n\ncategory: keep\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, RewriteMemoCategory(tt.content, tt.tree))
		})
	}
}
