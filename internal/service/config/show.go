package config

func (uc config) Show() {
	uc.logger.Info("Configuration", "base_dir", uc.config.BaseDir())
	uc.logger.Info("Configuration", "todos_dir", uc.config.TodosDir())
	uc.logger.Info("Configuration", "memos_dir", uc.config.MemosDir())
	uc.logger.Info("Configuration", "todos_daystoseek", uc.config.TodosDaysToSeek())
}
