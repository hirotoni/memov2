package todo

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTodoShowDTO_Marshals(t *testing.T) {
	date := time.Date(2023, 1, 2, 0, 0, 0, 0, time.Local)
	f, err := domain.NewTodosFile(date)
	require.NoError(t, err)
	f.SetHeadingBlocks([]*markdown.HeadingBlock{
		{HeadingText: utils.HeadingTodos.Text, ContentText: "- item", Level: 2},
	})

	dto := newTodoShowDTO(f, "20230102Mon_todos.md")
	assert.Equal(t, "20230102Mon_todos.md", dto.Path)
	assert.Equal(t, date, dto.Date)
	require.Len(t, dto.Headings, 1)
	assert.Equal(t, utils.HeadingTodos.Text, dto.Headings[0].Heading)
	assert.Equal(t, "- item", dto.Headings[0].Content)

	// Round-trips through JSON without error.
	b, err := json.Marshal(dto)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"path":"20230102Mon_todos.md"`)
}

func TestNewTodoShowDTO_EmptyHeadingsIsArray(t *testing.T) {
	date := time.Date(2023, 1, 2, 0, 0, 0, 0, time.Local)
	f, err := domain.NewTodosFile(date)
	require.NoError(t, err)
	f.SetHeadingBlocks(nil)

	dto := newTodoShowDTO(f, "20230102Mon_todos.md")
	b, err := json.Marshal(dto)
	require.NoError(t, err)
	// headings must marshal to [] (never null).
	assert.Contains(t, string(b), `"headings":[]`)
}
