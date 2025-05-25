package memo

import (
	"github.com/hirotoni/memov2/internal/interfaces"
)

type memo struct {
	config interfaces.ConfigProvider
	repos  interfaces.Repositories
	editor interfaces.Editor
}

func NewMemo(c interfaces.ConfigProvider, r interfaces.Repositories, e interfaces.Editor) interfaces.MemoUsecase {
	return memo{
		config: c,
		repos:  r,
		editor: e,
	}
}
