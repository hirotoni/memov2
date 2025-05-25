package domain

import (
	"reflect"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/stretchr/testify/assert"
)

func Test_file_Date(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		f    *file
		want time.Time
	}{
		{
			name: "basic pattern",
			f:    &file{date: now},
			want: now,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Date(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("file.Date() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_file_FileType(t *testing.T) {
	tests := []struct {
		name string
		f    *file
		want FileType
	}{
		{
			name: "basic pattern",
			f:    &file{fileType: FileTypeMemo},
			want: FileTypeMemo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.FileType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("file.FileType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_file_Title(t *testing.T) {
	tests := []struct {
		name string
		f    *file
		want string
	}{
		{
			name: "basic pattern",
			f:    &file{title: "test title"},
			want: "test title",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Title(); got != tt.want {
				t.Errorf("file.Title() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_file_TopLevelBodyContent(t *testing.T) {
	tlbc := &markdown.HeadingBlock{
		HeadingText: "test heading",
		Level:       1,
		ContentText: "test content",
	}
	tests := []struct {
		name string
		f    *file
		want *markdown.HeadingBlock
	}{
		{
			name: "basic pattern",
			f:    &file{topLevelBodyContent: tlbc},
			want: tlbc,
		},
		{
			name: "nil pattern",
			f:    &file{},
			want: &markdown.HeadingBlock{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.TopLevelBodyContent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("file.TopLevelBodyContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_file_HeadingBlocks(t *testing.T) {
	hbs := []*markdown.HeadingBlock{
		{
			HeadingText: "test heading",
			Level:       2,
			ContentText: "test content",
		},
		{
			HeadingText: "next heading",
			Level:       2,
			ContentText: "next content",
		},
		{
			HeadingText: "last heading",
			Level:       2,
			ContentText: "last content",
		},
	}

	tests := []struct {
		name string
		f    *file
		want []*markdown.HeadingBlock
	}{
		{
			name: "basic pattern",
			f:    &file{headingBlocks: hbs},
			want: hbs,
		},
		{
			name: "length zero pattern",
			f:    &file{},
			want: []*markdown.HeadingBlock{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.HeadingBlocks(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("file.HeadingBlocks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_file_LastHeadingBlock(t *testing.T) {
	hbs := []*markdown.HeadingBlock{
		{
			HeadingText: "test heading",
			Level:       2,
			ContentText: "test content",
		},
		{
			HeadingText: "next heading",
			Level:       2,
			ContentText: "next content",
		},
		{
			HeadingText: "last heading",
			Level:       2,
			ContentText: "last content",
		},
	}
	tests := []struct {
		name string
		f    *file
		want *markdown.HeadingBlock
	}{
		{
			name: "basic pattern",
			f:    &file{headingBlocks: hbs},
			want: hbs[2],
		},
		{
			name: "length zero pattern",
			f:    &file{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.LastHeadingBlock(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("file.LastHeadingBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_file_SetDate(t *testing.T) {
	now := time.Now()
	zero := time.Time{}
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		f    *file
		args args
		want time.Time
	}{
		{
			name: "basic pattern",
			f:    &file{},
			args: args{date: now},
			want: now,
		},
		{
			name: "zero value pattern",
			f:    &file{},
			args: args{date: zero},
			want: zero,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.SetDate(tt.args.date)
			assert.Equal(t, tt.f.date, tt.want)
		})
	}
}

func Test_file_SetTopLevelBodyContent(t *testing.T) {
	tlbc := &markdown.HeadingBlock{
		HeadingText: "test heading",
		Level:       1,
		ContentText: "test content",
	}
	type args struct {
		content *markdown.HeadingBlock
	}
	tests := []struct {
		name string
		f    *file
		args args
		want *markdown.HeadingBlock
	}{
		{
			name: "basic pattern",
			f:    &file{},
			args: args{content: tlbc},
			want: tlbc,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.SetTopLevelBodyContent(tt.args.content)
			if got := tt.f.TopLevelBodyContent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("file.TopLevelBodyContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_file_SetHeadingBlocks(t *testing.T) {
	hbs := []*markdown.HeadingBlock{
		{
			HeadingText: "test heading",
			Level:       2,
			ContentText: "test content",
		},
		{
			HeadingText: "next heading",
			Level:       2,
			ContentText: "next content",
		},
		{
			HeadingText: "last heading",
			Level:       2,
			ContentText: "last content",
		},
	}
	type args struct {
		entities []*markdown.HeadingBlock
	}
	tests := []struct {
		name string
		f    *file
		args args
		want []*markdown.HeadingBlock
	}{
		{
			name: "basic pattern",
			f:    &file{},
			args: args{entities: hbs},
			want: hbs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.SetHeadingBlocks(tt.args.entities)
			if got := tt.f.HeadingBlocks(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("file.HeadingBlocks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_file_OverrideHeadingBlockMatched(t *testing.T) {
	hbs := []*markdown.HeadingBlock{
		{
			HeadingText: "test heading",
			Level:       2,
			ContentText: "test content",
		},
		{
			HeadingText: "next heading",
			Level:       2,
			ContentText: "next content",
		},
		{
			HeadingText: "last heading",
			Level:       2,
			ContentText: "last content",
		},
	}
	target := &markdown.HeadingBlock{
		HeadingText: "next heading",
		Level:       2,
		ContentText: "edited content",
	}
	expected := []*markdown.HeadingBlock{
		{
			HeadingText: "test heading",
			Level:       2,
			ContentText: "test content",
		},
		{
			HeadingText: "next heading",
			Level:       2,
			ContentText: "edited content",
		},
		{
			HeadingText: "last heading",
			Level:       2,
			ContentText: "last content",
		},
	}

	type args struct {
		input *markdown.HeadingBlock
	}
	tests := []struct {
		name    string
		f       *file
		args    args
		wantErr bool
		want    []*markdown.HeadingBlock
	}{
		{
			name:    "basic pattern",
			f:       &file{headingBlocks: hbs},
			args:    args{input: target},
			wantErr: false,
			want:    expected,
		},
		{
			name:    "not found pattern",
			f:       &file{headingBlocks: []*markdown.HeadingBlock{}},
			args:    args{input: target},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.OverrideHeadingBlockMatched(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("file.OverrideHeadingBlockMatched() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if got := tt.f.HeadingBlocks(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("file.HeadingBlocks() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_file_OverrideHeadingBlocksMatched(t *testing.T) {
	hbs := []*markdown.HeadingBlock{
		{
			HeadingText: "test heading",
			Level:       2,
			ContentText: "test content",
		},
		{
			HeadingText: "next heading",
			Level:       2,
			ContentText: "next content",
		},
		{
			HeadingText: "last heading",
			Level:       2,
			ContentText: "last content",
		},
	}
	target := &markdown.HeadingBlock{
		HeadingText: "next heading",
		Level:       2,
		ContentText: "edited content",
	}
	type args struct {
		entities []*markdown.HeadingBlock
	}
	tests := []struct {
		name    string
		f       *file
		args    args
		wantErr bool
	}{
		{
			name:    "basic pattern",
			f:       &file{headingBlocks: hbs},
			args:    args{entities: []*markdown.HeadingBlock{target}},
			wantErr: false,
		},
		{
			name:    "not found pattern",
			f:       &file{headingBlocks: []*markdown.HeadingBlock{}},
			args:    args{entities: []*markdown.HeadingBlock{target}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.OverrideHeadingBlocksMatched(tt.args.entities); (err != nil) != tt.wantErr {
				t.Errorf("file.OverrideHeadingBlocksMatched() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
