package config

import (
	"encoding/json"
	"os"

	"github.com/hirotoni/memov2/internal/common"
)

func (uc config) Show(asJSON bool) error {
	if asJSON {
		dto := configDTO{
			BaseDir:         uc.config.BaseDir(),
			TodosDir:        uc.config.TodosDir(),
			MemosDir:        uc.config.MemosDir(),
			TodosDaysToSeek: uc.config.TodosDaysToSeek(),
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		// Return encode errors (e.g. a closed stdout) so the command exits
		// non-zero instead of silently reporting success to a JSON consumer.
		if err := enc.Encode(dto); err != nil {
			return common.Wrap(err, common.ErrorTypeService, "failed to encode config")
		}
		return nil
	}

	uc.logger.Info("Configuration", "base_dir", uc.config.BaseDir())
	uc.logger.Info("Configuration", "todos_dir", uc.config.TodosDir())
	uc.logger.Info("Configuration", "memos_dir", uc.config.MemosDir())
	uc.logger.Info("Configuration", "todos_daystoseek", uc.config.TodosDaysToSeek())
	return nil
}
