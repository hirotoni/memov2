package memo

import (
	"fmt"
	"os"
	"strings"
)

func (uc memo) ListCategories() error {
	categories, err := uc.repos.Memo().Categories()
	if err != nil {
		return err
	}

	for _, cat := range categories {
		fmt.Fprintln(os.Stdout, strings.Join(cat, "/"))
	}

	return nil
}
