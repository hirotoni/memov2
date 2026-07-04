package memo

import (
	"fmt"
	"os"
	"strings"
)

func (uc memo) ListCategories(asJSON bool) error {
	categories, err := uc.repos.Memo().Categories()
	if err != nil {
		return err
	}

	if asJSON {
		items := make([]categoryDTO, 0, len(categories))
		for _, cat := range categories {
			items = append(items, categoryDTO{
				Category: cat,
				Path:     strings.Join(cat, "/"),
			})
		}
		return writeJSON(items)
	}

	for _, cat := range categories {
		fmt.Fprintln(os.Stdout, strings.Join(cat, "/"))
	}

	return nil
}
