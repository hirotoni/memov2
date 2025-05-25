package config

import (
	"fmt"
)

func (uc config) Show() {
	fmt.Printf("base_dir: %s\n", uc.config.BaseDir())
	fmt.Printf("todos_dir: %s\n", uc.config.TodosDir())
	fmt.Printf("memos_dir: %s\n", uc.config.MemosDir())
	fmt.Printf("todos_daystoseek: %d\n", uc.config.TodosDaysToSeek())
}
