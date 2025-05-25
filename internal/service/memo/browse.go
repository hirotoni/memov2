package memo

import (
	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/ui/tui/memos"
)

func (uc memo) Browse() error {
	// Get the underlying TomlConfig for UI layer compatibility
	tomlConfig := uc.config.GetTomlConfig().(*toml.Config)
	err := memos.IntegratedMemos(tomlConfig, uc.editor)
	if err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error browsing memos")
	}
	return nil
}
