package browse

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hirotoni/memov2/components"
	"github.com/hirotoni/memov2/config"
	"github.com/hirotoni/memov2/models"
	"github.com/hirotoni/memov2/repos"
	"github.com/hirotoni/memov2/utils"
	"golang.org/x/term"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("63")).
			Padding(0, 1)

	ItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("230")).
				Background(lipgloss.Color("63"))

	// Remove padding from pagination and help styles
	PaginationStyle = list.DefaultStyles().PaginationStyle
	HelpStyle       = list.DefaultStyles().HelpStyle
)

type item struct {
	name     string
	path     string
	isDir    bool
	memo     models.MemoFileInterface
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
	return prefix + "  " + models.MemoTitle(i.name)
}

func (i item) Description() string {
	return ""
}

func (i item) FilterValue() string {
	return i.name
}

type BrowseModel struct {
	list         list.Model
	config       *config.TomlConfig
	currPath     string
	err          error
	width        int
	height       int
	items        []item // Root level items
	lastKeyPress string // Track last key press for 'gg' command
}

func New(c *config.TomlConfig) (*BrowseModel, error) {
	if !utils.Exists(c.MemosDir()) {
		if err := os.MkdirAll(c.MemosDir(), 0755); err != nil {
			return nil, fmt.Errorf("failed to create memos directory %s: %w", c.MemosDir(), err)
		}
		fmt.Printf("Created memos directory: %s\n", c.MemosDir())
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
	l.SetShowStatusBar(false)                      // Hide status bar
	l.SetShowHelp(true)                            // Show help footer
	l.SetFilteringEnabled(false)                   // Disable filtering
	l.DisableQuitKeybindings()                     // Disable default quit keybindings
	l.Title = ""                                   // Remove list title
	l.Styles.Title = TitleStyle
	l.Styles.PaginationStyle = PaginationStyle
	l.Styles.HelpStyle = HelpStyle
	l.KeyMap.NextPage = key.NewBinding() // Disable page down
	l.KeyMap.PrevPage = key.NewBinding() // Disable page up
	l.SetShowPagination(false)           // Hide pagination

	m := &BrowseModel{
		list:     l,
		config:   c,
		currPath: c.MemosDir(),
		width:    w,
		height:   h,
	}

	repo := repos.NewMemoRepo(m.config.MemosDir())
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
		m, cmd = BrowseKeybindings(m, msg)
		cmds = append(cmds, cmd)
	case error:
		m.err = msg
		return m, nil
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func BrowseKeybindings(m BrowseModel, msg tea.KeyMsg) (BrowseModel, tea.Cmd) {
	switch {
	case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "q"))):
		return m, tea.Quit
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
				components.OpenEditor(m.config.BaseDir, i.path)
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
		pageSize := m.height / 2
		newIndex := max(0, currentIndex-pageSize)
		m.list.Select(newIndex)
		return m, nil
	case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+d"))):
		// Move cursor down half a page
		currentIndex := m.list.Index()
		pageSize := m.height / 2
		newIndex := min(len(m.list.Items())-1, currentIndex+pageSize)
		m.list.Select(newIndex)
		return m, nil
	case key.Matches(msg, key.NewBinding(key.WithKeys("G"))):
		// Move to bottom (Shift+G)
		m.list.Select(len(m.list.Items()) - 1)
		return m, nil
	case key.Matches(msg, key.NewBinding(key.WithKeys("g"))):
		// Double 'g' to go to top
		if m.lastKeyPress == "g" {
			m.list.Select(0)
			m.lastKeyPress = ""
		} else {
			m.lastKeyPress = "g"
		}
		return m, nil
	}
	// Reset last key press if any other key is pressed
	if msg.String() != "g" {
		m.lastKeyPress = ""
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

	title := TitleStyle.Render("File Browser (h: collapse, l: expand)")
	return title + "\n" + m.list.View()
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

	repo := repos.NewMemoRepo(m.config.MemosDir())
	memos, err := repo.MemoEntries()
	if err != nil {
		return nil, fmt.Errorf("failed to get memo entries: %w", err)
	}

	memoMap := make(map[string]models.MemoFileInterface)
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
