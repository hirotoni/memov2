package picker

import (
	"strings"

	"github.com/hirotoni/memov2/internal/interfaces"
)

// SelectCategoryForNew shows a category picker (with an explicit "no category"
// option), then prompts for a title. It returns the chosen category tree (empty
// slice for "no category") and the entered title. ok is false if cancelled.
func SelectCategoryForNew(repo interfaces.MemoRepo) (tree []string, title string, ok bool, err error) {
	categories, err := repo.Categories()
	if err != nil {
		return nil, "", false, err
	}

	items := make([]Item, 0, len(categories)+1)
	// First option: create without a category.
	items = append(items, Item{
		Display:  "(no category)",
		FilterBy: "no category",
		Payload:  []string{},
	})
	for _, cat := range categories {
		joined := strings.Join(cat, "/")
		items = append(items, Item{
			Display:  joined,
			FilterBy: joined,
			Payload:  cat,
		})
	}

	res, err := Run(Config{
		Title:         "New memo — pick a category",
		Items:         items,
		WithInput:     true,
		InputPrompt:   "Title: ",
		AllowFreeText: true,
		FreeTextLabel: func(q string) string { return "+ new category \"" + q + "\"" },
	})
	if err != nil {
		return nil, "", false, err
	}
	if res.Cancelled {
		return nil, "", false, nil
	}
	// A typed-but-unlisted category: split the path like the old --category flag.
	if res.FreeText != "" {
		return strings.Split(res.FreeText, "/"), res.InputText, true, nil
	}
	if res.Item == nil {
		return nil, "", false, nil
	}
	return res.Item.Payload.([]string), res.InputText, true, nil
}
