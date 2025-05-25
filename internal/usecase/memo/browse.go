package memo

import (
	"fmt"

	"github.com/hirotoni/memov2/internal/ui/tui/memos"
)

func (uc memo) Browse() error {
	err := memos.IntegratedMemos(&uc.config)
	if err != nil {
		return fmt.Errorf("error browsing memos: %v", err)
	}
	return nil
}
