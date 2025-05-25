package domain

import (
	"reflect"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoFileInterface(t *testing.T) {
	var _ MemoFileInterface = &MemoFile{}
	var _ FileInterface = &MemoFile{}
}

func TestMemoFile_Location(t *testing.T) {
	tests := []struct {
		name         string // description of this test case
		categoryTree []string
		want         string
	}{
		{
			name:         "basic pattern",
			categoryTree: []string{"top", "middle", "bottom"},
			want:         "top/middle/bottom",
		},
		{
			name:         "no category tree",
			categoryTree: []string{},
			want:         "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f MemoFile
			f.SetCategoryTree(tt.categoryTree)

			got := f.Location()
			assert.Equal(t, tt.want, got, "Location() = %v, want %v", got, tt.want)
		})
	}
}

func TestMemoTitle(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		filename string
		want     string
	}{
		{
			name:     "basic pattern",
			filename: "20250612Thu111111_memo_title.md",
			want:     "title",
		},
		{
			name:     "unmatch pattern",
			filename: "unmatch",
			want:     "unmatch",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MemoTitle(tt.filename)
			assert.Equal(t, tt.want, got, "MemoTitle() = %v, want %v", got, tt.want)
		})
	}
}

func TestNewMemoFile(t *testing.T) {
	now := time.Now()
	type args struct {
		date         time.Time
		title        string
		categoryTree []string
	}
	tests := []struct {
		name    string
		args    args
		want    MemoFileInterface
		wantErr bool
	}{
		{
			name: "basic pattern",
			args: args{
				date:         now,
				title:        "test title",
				categoryTree: []string{"some category"},
			},
			want: &MemoFile{
				file: file{
					date:     now,
					title:    "test title",
					fileType: FileTypeMemo,
				},
				categoryTree: []string{"some category"},
			},
			wantErr: false,
		},
		{
			name: "date is zero",
			args: args{
				date:         time.Time{},
				title:        "test title",
				categoryTree: []string{"some category"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMemoFile(tt.args.date, tt.args.title, tt.args.categoryTree)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMemoFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemoFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoFile_FileName(t *testing.T) {
	fixedTime, err := time.Parse(time.DateTime, "2025-10-01 12:34:56")
	require.NoError(t, err)

	tests := []struct {
		name string
		f    *MemoFile
		want string
	}{
		{
			name: "basic pattern",
			f: &MemoFile{
				file: file{
					date:     fixedTime,
					title:    "test title",
					fileType: FileTypeMemo,
				},
				categoryTree: []string{"some category"},
			},
			want: "20251001Wed123456_memo_test-title.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.FileName(); got != tt.want {
				t.Errorf("MemoFile.FileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoFile_CategoryTree(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		f    *MemoFile
		want []string
	}{
		{
			name: "basic pattern",
			f: &MemoFile{
				file: file{
					date:     now,
					title:    "test title",
					fileType: FileTypeMemo,
				},
				categoryTree: []string{"some category"},
			},
			want: []string{"some category"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.CategoryTree(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoFile.CategoryTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoFile_SetCategoryTree(t *testing.T) {
	type args struct {
		tree []string
	}
	tests := []struct {
		name string
		f    *MemoFile
		args args
		want []string
	}{
		{
			name: "basic pattern",
			f:    &MemoFile{categoryTree: []string{"some category"}},
			args: args{tree: []string{"edited tree"}},
			want: []string{"edited tree"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.SetCategoryTree(tt.args.tree)
			assert.Equal(t, tt.want, tt.f.categoryTree)
		})
	}
}

func TestMemoFile_ContentString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		f    *MemoFile
		want string
	}{
		{
			name: "basic pattern",
			f: &MemoFile{
				file: file{
					date:     time.Now(),
					fileType: FileTypeMemo,
					title:    "test title",
					topLevelBodyContent: &markdown.HeadingBlock{
						ContentText: "content",
					},
					headingBlocks: []*markdown.HeadingBlock{
						{
							HeadingText: "heading 1",
							Level:       2,
							ContentText: "context 1",
						},
						{
							HeadingText: "heading 2",
							Level:       2,
							ContentText: "context 2",
						},
					},
				},
				categoryTree: []string{"some category", "other category"},
			},
			want: `---
category: ["some category", "other category"]
---

# test title

content

## heading 1

context 1
## heading 2

context 2
`,
		},
		{
			name: "no category pattern",
			f: &MemoFile{
				file: file{
					date:     time.Now(),
					fileType: FileTypeMemo,
					title:    "test title",
					topLevelBodyContent: &markdown.HeadingBlock{
						ContentText: "content",
					},
					headingBlocks: []*markdown.HeadingBlock{
						{
							HeadingText: "heading 1",
							Level:       2,
							ContentText: "context 1",
						},
						{
							HeadingText: "heading 2",
							Level:       2,
							ContentText: "context 2",
						},
					},
				},
				categoryTree: []string{},
			},
			want: `---
category: []
---

# test title

content

## heading 1

context 1
## heading 2

context 2
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.f.ContentString()
			assert.Equal(t, tt.want, got)
		})
	}
}
