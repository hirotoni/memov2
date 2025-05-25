package memo

import (
	"fmt"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/ui/tui/memos"
)

func (uc memo) Browse() error {
	// Get the underlying TomlConfig for UI layer compatibility
	tomlConfig := uc.config.GetTomlConfig().(*config.TomlConfig)
	err := memos.IntegratedMemos(tomlConfig, uc.editor)
	if err != nil {
		return fmt.Errorf("error browsing memos: %v", err)
	}
	return nil
}
