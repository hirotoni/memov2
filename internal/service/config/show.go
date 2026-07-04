package config

import (
	"encoding/json"
	"os"
)

func (uc config) Show(asJSON bool) {
	if asJSON {
		dto := configDTO{
			BaseDir:         uc.config.BaseDir(),
			TodosDir:        uc.config.TodosDir(),
			MemosDir:        uc.config.MemosDir(),
			TodosDaysToSeek: uc.config.TodosDaysToSeek(),
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		// Errors here are unexpected (encoding a plain struct to stdout); log to
		// stderr to keep stdout clean for the JSON consumer.
		if err := enc.Encode(dto); err != nil {
			uc.logger.Error("failed to encode config", "error", err)
		}
		return
	}

	uc.logger.Info("Configuration", "base_dir", uc.config.BaseDir())
	uc.logger.Info("Configuration", "todos_dir", uc.config.TodosDir())
	uc.logger.Info("Configuration", "memos_dir", uc.config.MemosDir())
	uc.logger.Info("Configuration", "todos_daystoseek", uc.config.TodosDaysToSeek())
}
