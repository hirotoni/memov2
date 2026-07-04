package config

// configDTO is the machine-readable representation of the configuration emitted
// by `config show --json`. It is a stable, documented contract for AI/automation
// consumers; fields may be added in a backward-compatible way only.
type configDTO struct {
	BaseDir         string `json:"base_dir"`
	TodosDir        string `json:"todos_dir"`
	MemosDir        string `json:"memos_dir"`
	TodosDaysToSeek int    `json:"todos_days_to_seek"`
}
