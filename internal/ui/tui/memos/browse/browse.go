package browse

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/platform/editor"
	"github.com/hirotoni/memov2/internal/platform/fs"
	"github.com/hirotoni/memov2/internal/repository"
	"golang.org/x/term"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("63")).
			Padding(0, 1)

	ItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			MarginLeft(1).
			Padding(0, 1)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("230")).
				Background(lipgloss.Color("63")).
				MarginLeft(1).
				Padding(0, 1)
)

type item struct {
	name     string
	path     string
	isDir    bool
	memo     domain.MemoFileInterface
	depth    int
	expanded bool
	parent   *item
	children []item
}

func (i item) Title() string {
	prefix := strings.Repeat("  ", i.depth)
	if i.isDir {
		expandChar := "▶"
		if i.expanded {
			expandChar = "▼"
		}
		return prefix + expandChar + " " + i.name + "/"
	}
	return prefix + "  " + domain.MemoTitle(i.name)
}

func (i item) Description() string {
	return ""
}

func (i item) FilterValue() string {
	return i.name
}

type BrowseModel struct {
	list                 list.Model
	config               *config.TomlConfig
	currPath             string
	err                  error
	width                int
	height               int
	items                []item // Root level items
	showCategoryDialog   bool
	selectedCategories   map[string]bool
	newCategoryInput     string
	allCategories        [][]string // Change to [][]string to store hierarchical paths
	selectedMemo         domain.MemoFileInterface
	categoryDialogCursor int
	collapsedCategories  map[string]bool // Track which categories are collapsed
	categoryList         list.Model      // Add list model for category dialog
}

// Helper type for category items
type categoryItem struct {
	path      []string
	indent    string
	indicator string
	selected  bool
}

func (i categoryItem) FilterValue() string {
	return strings.Join(i.path, "/")
}

func (i categoryItem) Title() string {
	displayName := i.path[len(i.path)-1]
	checkmark := " "
	if i.selected {
		checkmark = "✓"
	}
	return fmt.Sprintf("%s%s[%s] %s", i.indent, i.indicator, checkmark, displayName)
}

func (i categoryItem) Description() string {
	return ""
}

func New(c *config.TomlConfig) (*BrowseModel, error) {
	if err := fs.EnsureDir(c.MemosDir()); err != nil {
		return nil, fmt.Errorf("failed to ensure memos directory %s: %w", c.MemosDir(), err)
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = SelectedItemStyle
	delegate.Styles.NormalTitle = ItemStyle
	delegate.Styles.NormalDesc = ItemStyle

	// Make the delegate more compact
	delegate.SetSpacing(0) // Remove spacing between items
	delegate.ShowDescription = false

	// Get initial terminal dimensions
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w, h = 80, 20 // fallback dimensions
	}

	l := list.New([]list.Item{}, delegate, w, h-2) // reduce space for title
	l.SetShowTitle(true)
	l.SetShowHelp(true)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	l.KeyMap.NextPage = key.NewBinding() // Disable page down
	l.KeyMap.PrevPage = key.NewBinding() // Disable page up

	l.Title = "Memos Browser"
	l.Styles.Title = TitleStyle

	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("ctrl+c/q", "quit")),
			key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "expand/open")),
			key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "collapse")),
			key.NewBinding(key.WithKeys("ctrl+u"), key.WithHelp("ctrl+u", "half page up")),
			key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("ctrl+d", "half page down")),
			key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "manage categories")),
		}
	}
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("ctrl+c/q", "quit")),
			key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "expand/open")),
			key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "collapse")),
			key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "manage categories")),
		}
	}

	cl := list.New([]list.Item{}, delegate, w, h)
	cl.Title = "Category Management"
	cl.Styles.Title = TitleStyle
	cl.SetShowTitle(true)
	cl.SetShowHelp(true)
	cl.SetShowStatusBar(false)
	cl.SetShowPagination(false)
	cl.SetFilteringEnabled(false)
	cl.DisableQuitKeybindings()

	m := &BrowseModel{
		list:                 l,
		config:               c,
		currPath:             c.MemosDir(),
		width:                w,
		height:               h,
		selectedCategories:   make(map[string]bool),
		categoryDialogCursor: 0,
		collapsedCategories:  make(map[string]bool),
		categoryList:         cl,
	}

	repo := repository.NewMemo(m.config.MemosDir())
	err = repo.TidyMemos()
	if err != nil {
		fmt.Print("Error tidying memos: ", err, "\n")
	}

	if _, err := m.updateItems(); err != nil {
		return nil, fmt.Errorf("error loading initial items: %w", err)
	}

	return m, nil
}

func (m BrowseModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m BrowseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 2)
		cmd, err := m.updateItems()
		if err != nil {
			return m, tea.Quit
		}
		return m, tea.Batch(cmd)

	case tea.KeyMsg:
		if m.showCategoryDialog {
			if m.newCategoryInput != "" {
				// When in category input mode, handle all keys as input except for special keys
				switch msg.String() {
				case "esc":
					m.newCategoryInput = ""
					return m, nil
				case "enter":
					// Add new category as a subcategory of the currently selected path
					newCat := strings.TrimSpace(m.newCategoryInput)
					if newCat != "" {
						var newPath []string
						if i, ok := m.categoryList.SelectedItem().(categoryItem); ok {
							newPath = append([]string{}, i.path...)
						}
						newPath = append(newPath, newCat)
						pathStr := strings.Join(newPath, string(filepath.Separator))
						m.selectedCategories[pathStr] = true
						m.allCategories = append(m.allCategories, newPath)
						m.newCategoryInput = ""
						m.updateCategoryItems()
					}
					return m, nil
				case "backspace":
					if len(m.newCategoryInput) > 0 {
						m.newCategoryInput = m.newCategoryInput[:len(m.newCategoryInput)-1]
					}
					return m, nil
				default:
					m.newCategoryInput += string(msg.Runes)
					return m, nil
				}
			}

			// Normal category dialog mode
			switch msg.String() {
			case "esc":
				m.showCategoryDialog = false
				return m, nil
			case "l":
				if i, ok := m.categoryList.SelectedItem().(categoryItem); ok {
					pathStr := strings.Join(i.path, string(filepath.Separator))
					m.collapsedCategories[pathStr] = false
					m.updateCategoryItems()
				}
				return m, nil
			case "h":
				if i, ok := m.categoryList.SelectedItem().(categoryItem); ok {
					pathStr := strings.Join(i.path, string(filepath.Separator))
					m.collapsedCategories[pathStr] = true
					m.updateCategoryItems()
				}
				return m, nil
			case " ":
				if i, ok := m.categoryList.SelectedItem().(categoryItem); ok {
					pathStr := strings.Join(i.path, string(filepath.Separator))
					m.selectedCategories[pathStr] = !m.selectedCategories[pathStr]
					m.updateCategoryItems()
				}
				return m, nil
			case "enter":
				// Save categories
				var selectedPath []string
				for pathStr, selected := range m.selectedCategories {
					if selected {
						// Only use the last selected path
						path := strings.Split(pathStr, string(filepath.Separator))
						if len(path) > len(selectedPath) {
							selectedPath = path
						}
					}
				}
				repo := repository.NewMemo(m.config.MemosDir())
				if err := repo.Move(m.selectedMemo, selectedPath); err != nil {
					m.err = fmt.Errorf("failed to move memo: %w", err)
					return m, nil
				}
				m.showCategoryDialog = false
				cmd, err := m.updateItems()
				if err != nil {
					return m, tea.Quit
				}
				return m, tea.Batch(cmd)
			case "n":
				m.newCategoryInput = " " // Start with a space to indicate input mode
				return m, nil
			default:
				var cmd tea.Cmd
				m.categoryList, cmd = m.categoryList.Update(msg)
				return m, cmd
			}
		} else {
			m, cmd = BrowseKeybindings(m, msg)
			cmds = append(cmds, cmd)
		}
	case error:
		m.err = msg
		return m, nil
	}

	if !m.showCategoryDialog {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func BrowseKeybindings(m BrowseModel, msg tea.KeyMsg) (BrowseModel, tea.Cmd) {
	switch {
	case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "q"))):
		return m, tea.Quit
	case key.Matches(msg, key.NewBinding(key.WithKeys("c"))):
		if i, ok := m.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			m.selectedMemo = i.memo
			m.showCategoryDialog = true
			m.selectedCategories = make(map[string]bool)

			// Store the current memo's category path
			if len(i.memo.CategoryTree()) > 0 {
				pathStr := strings.Join(i.memo.CategoryTree(), string(filepath.Separator))
				m.selectedCategories[pathStr] = true
			}

			// Get all unique categories
			repo := repository.NewMemo(m.config.MemosDir())
			categories, err := repo.Categories()
			if err != nil {
				m.err = fmt.Errorf("failed to get categories: %w", err)
				return m, nil
			}

			m.allCategories = categories
			m.updateCategoryItems()
			return m, nil
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("l"))):
		if i, ok := m.list.SelectedItem().(item); ok {
			if i.isDir {
				// Find and expand the directory
				found := false
				m.traverseAndModify(&m.items, i.path, func(item *item) {
					item.expanded = true
					found = true
				})
				if found {
					cmd, err := m.updateItems()
					if err != nil {
						return m, tea.Quit
					}
					return m, tea.Batch(cmd)
				}
			} else {
				err := editor.DEO.Open(m.config.BaseDir, i.path)
				if err != nil {
					m.err = fmt.Errorf("failed to open editor: %w", err)
					return m, nil
				}
				return m, nil
			}
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("h"))):
		if i, ok := m.list.SelectedItem().(item); ok {
			if i.isDir {
				// If directory has a parent, collapse parent and move cursor there
				if i.parent != nil {
					parentPath := i.parent.path
					found := false
					m.traverseAndModify(&m.items, parentPath, func(item *item) {
						item.expanded = false
						found = true
					})
					if found {
						cmd, err := m.updateItems()
						if err != nil {
							return m, tea.Quit
						}

						// Move cursor to parent
						items := m.list.Items()
						for idx, listItem := range items {
							if it, ok := listItem.(item); ok && it.path == parentPath {
								m.list.Select(idx)
								break
							}
						}

						return m, tea.Batch(cmd)
					}
				} else {
					// If it's a root directory, just collapse it
					found := false
					m.traverseAndModify(&m.items, i.path, func(item *item) {
						item.expanded = false
						found = true
					})
					if found {
						cmd, err := m.updateItems()
						if err != nil {
							return m, tea.Quit
						}
						return m, tea.Batch(cmd)
					}
				}
			} else if i.parent != nil {
				// For files, collapse their parent directory and move cursor to parent
				parentPath := i.parent.path
				found := false
				m.traverseAndModify(&m.items, parentPath, func(item *item) {
					item.expanded = false
					found = true
				})
				if found {
					cmd, err := m.updateItems()
					if err != nil {
						return m, tea.Quit
					}

					// Move cursor to parent
					items := m.list.Items()
					for idx, listItem := range items {
						if it, ok := listItem.(item); ok && it.path == parentPath {
							m.list.Select(idx)
							break
						}
					}

					return m, tea.Batch(cmd)
				}
			}
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+u"))):
		// Move cursor up half a page
		currentIndex := m.list.Index()
		pageSize := len(m.list.VisibleItems()) / 2
		newIndex := max(0, currentIndex-pageSize)
		m.list.Select(newIndex)
		return m, nil
	case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+d"))):
		// Move cursor down half a page
		currentIndex := m.list.Index()
		pageSize := len(m.list.VisibleItems()) / 2
		newIndex := min(len(m.list.Items())-1, currentIndex+pageSize)
		m.list.Select(newIndex)
		return m, nil
	}
	return m, nil
}

// traverseAndModify traverses the tree and applies the modifier function when the target path is found
func (m *BrowseModel) traverseAndModify(items *[]item, targetPath string, modifier func(*item)) {
	for i := range *items {
		if (*items)[i].path == targetPath {
			modifier(&(*items)[i])
			return
		}
		if (*items)[i].isDir {
			m.traverseAndModify(&(*items)[i].children, targetPath, modifier)
		}
	}
}

func (m BrowseModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	if m.showCategoryDialog {
		if m.newCategoryInput != "" {
			var sb strings.Builder
			currentPath := ""
			visibleCats := m.visibleCategories()
			if m.categoryDialogCursor < len(visibleCats) {
				currentPath = strings.Join(visibleCats[m.categoryDialogCursor], " > ")
			}
			sb.WriteString(fmt.Sprintf("Adding new category under: %s\n", currentPath))
			sb.WriteString(fmt.Sprintf("New category: %s\n", strings.TrimSpace(m.newCategoryInput)))
			sb.WriteString("\nPress Enter to add, Esc to cancel\n")
			return sb.String()
		}

		return m.categoryList.View()
	}

	return m.list.View()
}

func (m *BrowseModel) updateItems() (tea.Cmd, error) {
	// Build the tree structure starting from the memos directory
	rootItems, err := m.buildTree(m.config.MemosDir(), 0, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build tree: %w", err)
	}

	// Preserve expansion state from previous items
	m.preserveExpansionState(&rootItems, &m.items)
	m.items = rootItems

	// Flatten the tree for display
	flatItems := m.flattenTree(rootItems)
	listItems := make([]list.Item, len(flatItems))
	for i, item := range flatItems {
		listItems[i] = item
	}

	m.list.SetItems(listItems)
	return nil, nil
}

func (m *BrowseModel) buildTree(path string, depth int, parent *item) ([]item, error) {
	var items []item

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	// First process directories
	for _, entry := range entries {
		if entry.IsDir() {
			newItem := item{
				name:     entry.Name(),
				path:     filepath.Join(path, entry.Name()),
				isDir:    true,
				depth:    depth,
				expanded: false, // Start collapsed
				parent:   parent,
			}

			// Recursively build children
			children, err := m.buildTree(newItem.path, depth+1, &newItem)
			if err != nil {
				return nil, err
			}
			newItem.children = children
			items = append(items, newItem)
		}
	}

	// Then process files

	repo := repository.NewMemo(m.config.MemosDir())
	memos, err := repo.MemoEntries()
	if err != nil {
		return nil, fmt.Errorf("failed to get memo entries: %w", err)
	}

	memoMap := make(map[string]domain.MemoFileInterface)
	for _, memo := range memos {
		memoMap[filepath.Join(memo.Location(), memo.FileName())] = memo
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			relPath := strings.TrimPrefix(filepath.Join(path, entry.Name()), m.config.MemosDir())
			relPath = strings.TrimPrefix(relPath, string(filepath.Separator))

			newItem := item{
				name:   entry.Name(),
				path:   filepath.Join(path, entry.Name()),
				isDir:  false,
				depth:  depth,
				parent: parent,
			}

			if memo, ok := memoMap[relPath]; ok {
				newItem.memo = memo
			}

			items = append(items, newItem)
		}
	}

	return items, nil
}

func (m *BrowseModel) flattenTree(items []item) []item {
	var result []item
	for _, item := range items {
		result = append(result, item)
		if item.isDir && item.expanded {
			result = append(result, m.flattenTree(item.children)...)
		}
	}
	return result
}

// preserveExpansionState copies the expanded state from old items to new items
func (m *BrowseModel) preserveExpansionState(newItems *[]item, oldItems *[]item) {
	oldMap := make(map[string]bool)
	m.buildExpansionMap(oldItems, oldMap)
	m.applyExpansionMap(newItems, oldMap)
}

func (m *BrowseModel) buildExpansionMap(items *[]item, expansionMap map[string]bool) {
	for _, item := range *items {
		if item.isDir {
			expansionMap[item.path] = item.expanded
			m.buildExpansionMap(&item.children, expansionMap)
		}
	}
}

func (m *BrowseModel) applyExpansionMap(items *[]item, expansionMap map[string]bool) {
	for i := range *items {
		if (*items)[i].isDir {
			if expanded, ok := expansionMap[(*items)[i].path]; ok {
				(*items)[i].expanded = expanded
			}
			m.applyExpansionMap(&(*items)[i].children, expansionMap)
		}
	}
}

// Helper functions for min/max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Helper function to check if a path is a child of another path
func isChildPath(child, parent []string) bool {
	if len(child) <= len(parent) {
		return false
	}
	for i := range parent {
		if parent[i] != child[i] {
			return false
		}
	}
	return true
}

// Helper function to get visible categories based on collapsed state
func (m *BrowseModel) visibleCategories() [][]string {
	if len(m.allCategories) == 0 {
		return nil
	}

	visible := make([][]string, 0)
	for _, path := range m.allCategories {
		// Check if any parent is collapsed
		isHidden := false
		for i := 1; i < len(path); i++ {
			parentPath := path[:i]
			parentStr := strings.Join(parentPath, string(filepath.Separator))
			if m.collapsedCategories[parentStr] {
				isHidden = true
				break
			}
		}
		if !isHidden {
			visible = append(visible, path)
		}
	}
	return visible
}

// Helper function to sort categories maintaining parent-child relationships
func sortCategories(categories [][]string) [][]string {
	// Create a map for quick parent lookup
	parentMap := make(map[string][][]string)
	rootCategories := [][]string{}

	// First, organize categories by their parent path
	for _, path := range categories {
		if len(path) == 1 {
			rootCategories = append(rootCategories, path)
		} else {
			parentPath := strings.Join(path[:len(path)-1], string(filepath.Separator))
			parentMap[parentPath] = append(parentMap[parentPath], path)
		}
	}

	// Sort root categories
	sort.Slice(rootCategories, func(i, j int) bool {
		return rootCategories[i][0] < rootCategories[j][0]
	})

	// Helper function to sort children recursively
	var sortedPaths [][]string
	var sortChildren func(path []string)
	sortChildren = func(path []string) {
		sortedPaths = append(sortedPaths, path)

		// Get and sort children
		pathStr := strings.Join(path, string(filepath.Separator))
		children := parentMap[pathStr]
		sort.Slice(children, func(i, j int) bool {
			return children[i][len(children[i])-1] < children[j][len(children[j])-1]
		})

		// Recursively process children
		for _, child := range children {
			sortChildren(child)
		}
	}

	// Process all root categories
	for _, root := range rootCategories {
		sortChildren(root)
	}

	return sortedPaths
}

func (m *BrowseModel) updateCategoryItems() {
	var items []list.Item
	lastParentPath := []string{}

	// Sort categories before creating items
	sortedCategories := sortCategories(m.allCategories)
	visibleCats := [][]string{}

	// Filter visible categories while maintaining order
	for _, path := range sortedCategories {
		isVisible := true
		for i := 1; i < len(path); i++ {
			parentPath := path[:i]
			if m.collapsedCategories[strings.Join(parentPath, string(filepath.Separator))] {
				isVisible = false
				break
			}
		}
		if isVisible {
			visibleCats = append(visibleCats, path)
		}
	}

	for _, path := range visibleCats {
		// Calculate indent
		indent := ""
		if len(path) > 1 {
			isChild := len(path) > len(lastParentPath) &&
				strings.HasPrefix(strings.Join(path, string(filepath.Separator)),
					strings.Join(lastParentPath, string(filepath.Separator))+string(filepath.Separator))

			if isChild {
				indent = strings.Repeat("  ", len(path)-1)
			} else {
				lastParentPath = path[:len(path)-1]
				indent = strings.Repeat("  ", len(path)-1)
			}
		}

		// Add expand/collapse indicator
		indicator := "  "
		pathStr := strings.Join(path, string(filepath.Separator))
		hasChildren := false
		isCollapsed := m.collapsedCategories[pathStr]
		for _, otherPath := range m.allCategories {
			if isChildPath(otherPath, path) {
				hasChildren = true
				break
			}
		}
		if hasChildren {
			if isCollapsed {
				indicator = "▶ "
			} else {
				indicator = "▼ "
			}
		}

		items = append(items, categoryItem{
			path:      path,
			indent:    indent,
			indicator: indicator,
			selected:  m.selectedCategories[pathStr],
		})

		// Update lastParentPath if this is a potential parent
		if len(path) == 1 || !strings.HasPrefix(pathStr, strings.Join(lastParentPath, string(filepath.Separator))) {
			lastParentPath = path
		}
	}

	m.categoryList.SetItems(items)
}
