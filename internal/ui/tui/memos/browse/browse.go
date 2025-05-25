package browse

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/domain"
	"github.com/hirotoni/memov2/internal/domain/markdown"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/platform"
	"github.com/hirotoni/memov2/internal/repositories/memo"
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
	config               *toml.Config
	editor               interfaces.Editor
	currPath             string
	err                  error
	width                int
	height               int
	items                []item // Root level items
	showCategoryDialog   bool
	showDeleteDialog     bool
	showRenameDialog     bool
	showDuplicateDialog  bool
	showNewMemoDialog    bool
	showPreview          bool // Toggle preview pane on/off
	selectedCategories   map[string]bool
	newCategoryInput     string
	renameInput          string
	newMemoTitleInput    string
	allCategories        [][]string // Change to [][]string to store hierarchical paths
	selectedMemo         domain.MemoFileInterface
	selectedCategoryTree []string // Category tree when a directory is selected for new memo
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

func New(c *toml.Config, e interfaces.Editor) (*BrowseModel, error) {
	if err := platform.EnsureDir(c.MemosDir()); err != nil {
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
			// Navigation
			key.NewBinding(key.WithKeys("ctrl+u"), key.WithHelp("ctrl+u", "up by 10 items")),
			key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("ctrl+d", "down by 10 items")),
			// Hierarchy
			key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "expand/open")),
			key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "collapse")),
			key.NewBinding(key.WithKeys(">"), key.WithHelp(">", "expand all under")),
			key.NewBinding(key.WithKeys("<"), key.WithHelp("<", "collapse all under")),
			// View
			key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "toggle preview")),
			// Memo operations
			key.NewBinding(key.WithKeys("N"), key.WithHelp("N", "new memo in category")),
			key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rename")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "duplicate (with new timestamp)")),
			key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "delete (move to trash)")),
			key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "manage categories")),
			// Quit
			key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("ctrl+c/q", "quit")),
		}
	}
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			// Navigation
			key.NewBinding(key.WithKeys("ctrl+u"), key.WithHelp("ctrl+u", "up 10")),
			key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("ctrl+d", "down 10")),
			// Hierarchy
			key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "expand/open")),
			key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "collapse")),
			key.NewBinding(key.WithKeys(">"), key.WithHelp(">", "expand all")),
			key.NewBinding(key.WithKeys("<"), key.WithHelp("<", "collapse all")),
			// View
			key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "preview")),
			// Memo operations
			key.NewBinding(key.WithKeys("N"), key.WithHelp("N", "new")),
			key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rename")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "duplicate")),
			key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "delete")),
			key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "categories")),
			// Quit
			key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("ctrl+c/q", "quit")),
		}
	}

	cl := list.New([]list.Item{}, delegate, w, h)
	cl.Title = "Category Management"
	cl.Styles.Title = TitleStyle
	cl.SetShowTitle(true)
	cl.SetShowHelp(true)
	cl.SetShowStatusBar(false)
	cl.SetShowPagination(false)
	cl.SetFilteringEnabled(false) // Disable filtering to keep UI simple and consistent
	cl.DisableQuitKeybindings()

	// Add custom help keys for category dialog
	cl.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "select")),
			key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "expand")),
			key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "collapse")),
			key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "save")),
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
		}
	}

	m := &BrowseModel{
		list:                 l,
		config:               c,
		editor:               e,
		currPath:             c.MemosDir(),
		width:                w,
		height:               h,
		showPreview:          true, // Preview enabled by default
		selectedCategories:   make(map[string]bool),
		categoryDialogCursor: 0,
		collapsedCategories:  make(map[string]bool),
		categoryList:         cl,
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

		// Adjust list width based on preview state
		if m.showPreview {
			listWidth := int(float64(msg.Width) * 0.5)
			m.list.SetWidth(listWidth)
		} else {
			m.list.SetWidth(msg.Width)
		}

		m.list.SetHeight(msg.Height - 2)
		m.categoryList.SetWidth(msg.Width)
		m.categoryList.SetHeight(msg.Height)

		cmd, err := m.updateItems()
		if err != nil {
			return m, tea.Quit
		}
		return m, tea.Batch(cmd)

	case tea.KeyMsg:
		if m.showNewMemoDialog {
			// Handle new memo dialog
			switch msg.String() {
			case "esc":
				// Cancel new memo creation
				m.showNewMemoDialog = false
				m.newMemoTitleInput = ""
				return m, nil
			case "enter":
				// Create new memo
				newTitle := strings.TrimSpace(m.newMemoTitleInput)
				if newTitle != "" {
					// Get category tree from either the selected memo or the selected directory
					var categoryTree []string
					if m.selectedMemo != nil {
						categoryTree = m.selectedMemo.CategoryTree()
					} else {
						categoryTree = m.selectedCategoryTree
					}

					// Create new memo with the determined category
					logger := common.DefaultLogger()
					repo := memo.NewMemo(m.config.MemosDir(), logger)
					newMemo, err := domain.NewMemoFile(
						time.Now(),
						newTitle,
						categoryTree,
					)
					if err != nil {
						m.err = fmt.Errorf("failed to create new memo: %w", err)
						m.showNewMemoDialog = false
						m.newMemoTitleInput = ""
						return m, nil
					}

					// Initialize with empty content structure (proper memo format)
					emptyContent := &markdown.HeadingBlock{
						Level:       0,
						HeadingText: "",
						ContentText: "",
					}
					newMemo.SetTopLevelBodyContent(emptyContent)

					// Save the new memo
					if err := repo.Save(newMemo, true); err != nil {
						m.err = fmt.Errorf("failed to save new memo: %w", err)
						m.showNewMemoDialog = false
						m.newMemoTitleInput = ""
						return m, nil
					}

					m.showNewMemoDialog = false
					m.newMemoTitleInput = ""

					// Refresh the list
					cmd, err := m.updateItems()
					if err != nil {
						return m, tea.Quit
					}

					// Open the new memo in the editor
					newPath := filepath.Join(m.config.MemosDir(), newMemo.Location(), newMemo.FileName())
					if err := m.editor.Open(m.config.BaseDir(), newPath); err != nil {
						m.err = fmt.Errorf("failed to open editor: %w", err)
						return m, nil
					}

					return m, tea.Batch(cmd)
				}
				// If empty, just cancel
				m.showNewMemoDialog = false
				m.newMemoTitleInput = ""
				return m, nil
			case "backspace":
				if len(m.newMemoTitleInput) > 0 {
					// Properly handle multi-byte UTF-8 characters (like Japanese)
					_, size := utf8.DecodeLastRuneInString(m.newMemoTitleInput)
					m.newMemoTitleInput = m.newMemoTitleInput[:len(m.newMemoTitleInput)-size]
				}
				return m, nil
			default:
				m.newMemoTitleInput += string(msg.Runes)
				return m, nil
			}
		} else if m.showDuplicateDialog {
			// Handle duplicate confirmation dialog
			switch msg.String() {
			case "y", "Y":
				// Confirm duplication
				logger := common.DefaultLogger()
				repo := memo.NewMemo(m.config.MemosDir(), logger)
				_, err := repo.Duplicate(m.selectedMemo)
				if err != nil {
					m.err = fmt.Errorf("failed to duplicate memo: %w", err)
					m.showDuplicateDialog = false
					return m, nil
				}
				m.showDuplicateDialog = false
				cmd, err := m.updateItems()
				if err != nil {
					return m, tea.Quit
				}
				return m, tea.Batch(cmd)
			case "n", "N", "esc":
				// Cancel duplication
				m.showDuplicateDialog = false
				return m, nil
			}
			return m, nil
		} else if m.showRenameDialog {
			// Handle rename dialog
			switch msg.String() {
			case "esc":
				// Cancel rename
				m.showRenameDialog = false
				m.renameInput = ""
				return m, nil
			case "enter":
				// Confirm rename
				newTitle := strings.TrimSpace(m.renameInput)
				if newTitle != "" && newTitle != m.selectedMemo.Title() {
					logger := common.DefaultLogger()
					repo := memo.NewMemo(m.config.MemosDir(), logger)
					if err := repo.Rename(m.selectedMemo, newTitle); err != nil {
						m.err = fmt.Errorf("failed to rename memo: %w", err)
						m.showRenameDialog = false
						m.renameInput = ""
						return m, nil
					}
					m.showRenameDialog = false
					m.renameInput = ""
					cmd, err := m.updateItems()
					if err != nil {
						return m, tea.Quit
					}
					return m, tea.Batch(cmd)
				}
				// If same title or empty, just cancel
				m.showRenameDialog = false
				m.renameInput = ""
				return m, nil
			case "backspace":
				if len(m.renameInput) > 0 {
					// Properly handle multi-byte UTF-8 characters (like Japanese)
					_, size := utf8.DecodeLastRuneInString(m.renameInput)
					m.renameInput = m.renameInput[:len(m.renameInput)-size]
				}
				return m, nil
			default:
				m.renameInput += string(msg.Runes)
				return m, nil
			}
		} else if m.showDeleteDialog {
			// Handle delete confirmation dialog
			switch msg.String() {
			case "y", "Y":
				// Confirm deletion
				logger := common.DefaultLogger()
				repo := memo.NewMemo(m.config.MemosDir(), logger)
				if err := repo.Delete(m.selectedMemo); err != nil {
					m.err = fmt.Errorf("failed to delete memo: %w", err)
					m.showDeleteDialog = false
					return m, nil
				}
				m.showDeleteDialog = false
				cmd, err := m.updateItems()
				if err != nil {
					return m, tea.Quit
				}
				return m, tea.Batch(cmd)
			case "n", "N", "esc":
				// Cancel deletion
				m.showDeleteDialog = false
				return m, nil
			}
			return m, nil
		} else if m.showCategoryDialog {
			// Check if we're in new category input mode
			if m.categoryDialogCursor == -1 {
				// When in category input mode, handle all keys as input except for special keys
				switch msg.String() {
				case "esc":
					m.newCategoryInput = ""
					m.categoryDialogCursor = 0
					return m, nil
				case "enter":
					// Add new hierarchical category path
					newCat := strings.TrimSpace(m.newCategoryInput)
					if newCat != "" {
						// Parse hierarchical path - support both / and > as separators
						newCat = strings.ReplaceAll(newCat, " > ", "/")
						newCat = strings.ReplaceAll(newCat, ">", "/")

						// Split by / to get path segments
						segments := strings.Split(newCat, "/")

						// Clean each segment (trim spaces)
						var cleanSegments []string
						for _, seg := range segments {
							cleaned := strings.TrimSpace(seg)
							if cleaned != "" {
								cleanSegments = append(cleanSegments, cleaned)
							}
						}

						if len(cleanSegments) > 0 {
							// Add all parent paths if they don't exist
							for i := 1; i <= len(cleanSegments); i++ {
								subPath := cleanSegments[:i]
								pathStr := strings.Join(subPath, string(filepath.Separator))

								// Check if this path already exists
								exists := false
								for _, existing := range m.allCategories {
									if strings.Join(existing, string(filepath.Separator)) == pathStr {
										exists = true
										break
									}
								}

								// Add if it doesn't exist
								if !exists {
									m.allCategories = append(m.allCategories, subPath)
								}
							}

							// Don't auto-select, let user manually select the category
							m.newCategoryInput = ""
							m.categoryDialogCursor = 0
							m.updateCategoryItems()
						}
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
					// Clear all previous selections and select only this one
					m.selectedCategories = make(map[string]bool)
					m.selectedCategories[pathStr] = true
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
				logger := common.DefaultLogger()
				repo := memo.NewMemo(m.config.MemosDir(), logger)
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
				m.categoryDialogCursor = -1 // Use -1 as sentinel for input mode
				m.newCategoryInput = ""     // Start with empty input
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
	case key.Matches(msg, key.NewBinding(key.WithKeys("p"))):
		// Toggle preview pane
		m.showPreview = !m.showPreview

		// Adjust list width
		if m.showPreview {
			listWidth := int(float64(m.width) * 0.5)
			m.list.SetWidth(listWidth)
		} else {
			m.list.SetWidth(m.width)
		}

		return m, nil
	case key.Matches(msg, key.NewBinding(key.WithKeys("N"))):
		if i, ok := m.list.SelectedItem().(item); ok {
			if i.isDir {
				// Directory selected - extract category tree from path
				m.selectedMemo = nil
				m.selectedCategoryTree = m.pathToCategoryTree(i.path)
				m.showNewMemoDialog = true
				m.newMemoTitleInput = ""
				return m, nil
			} else if i.memo != nil {
				// File selected - use memo's category tree
				m.selectedMemo = i.memo
				m.selectedCategoryTree = nil
				m.showNewMemoDialog = true
				m.newMemoTitleInput = ""
				return m, nil
			}
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
		if i, ok := m.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			m.selectedMemo = i.memo
			m.showRenameDialog = true
			// Initialize with current title (replace hyphens with spaces for editing)
			m.renameInput = strings.ReplaceAll(i.memo.Title(), "-", " ")
			return m, nil
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("d"))):
		if i, ok := m.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			m.selectedMemo = i.memo
			m.showDuplicateDialog = true
			return m, nil
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys("D"))):
		if i, ok := m.list.SelectedItem().(item); ok && !i.isDir && i.memo != nil {
			m.selectedMemo = i.memo
			m.showDeleteDialog = true
			return m, nil
		}
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
			logger := common.DefaultLogger()
			repo := memo.NewMemo(m.config.MemosDir(), logger)
			categories, err := repo.Categories()
			if err != nil {
				m.err = fmt.Errorf("failed to get categories: %w", err)
				return m, nil
			}

			m.allCategories = categories

			// Update dialog title with memo context
			memoTitle := domain.MemoTitle(i.memo.FileName())
			currentCategory := "None"
			if len(i.memo.CategoryTree()) > 0 {
				currentCategory = strings.Join(i.memo.CategoryTree(), " > ")
			}
			m.categoryList.Title = fmt.Sprintf("Move '%s' (Currently in: %s)", memoTitle, currentCategory)

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
				err := m.editor.Open(m.config.BaseDir(), i.path)
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
		// Move cursor up by 10 items
		currentIndex := m.list.Index()
		newIndex := max(0, currentIndex-10)
		m.list.Select(newIndex)
		return m, nil
	case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+d"))):
		// Move cursor down by 10 items
		currentIndex := m.list.Index()
		newIndex := min(len(m.list.Items())-1, currentIndex+10)
		m.list.Select(newIndex)
		return m, nil
	case key.Matches(msg, key.NewBinding(key.WithKeys("<"))):
		// Collapse all directories under selected item
		if i, ok := m.list.SelectedItem().(item); ok {
			if i.isDir {
				// Collapse the selected directory and all its children
				found := false
				m.traverseAndModify(&m.items, i.path, func(item *item) {
					m.expandItemAndChildren(item, false)
					found = true
				})
				if found {
					cmd, err := m.updateItems()
					if err != nil {
						return m, tea.Quit
					}
					return m, cmd
				}
			} else if i.parent != nil {
				// For files, collapse their parent directory and all its children
				parentPath := i.parent.path
				found := false
				m.traverseAndModify(&m.items, parentPath, func(item *item) {
					m.expandItemAndChildren(item, false)
					found = true
				})
				if found {
					cmd, err := m.updateItems()
					if err != nil {
						return m, tea.Quit
					}
					return m, cmd
				}
			}
		}
	case key.Matches(msg, key.NewBinding(key.WithKeys(">"))):
		// Expand all directories under selected item
		if i, ok := m.list.SelectedItem().(item); ok {
			if i.isDir {
				// Expand the selected directory and all its children
				found := false
				m.traverseAndModify(&m.items, i.path, func(item *item) {
					m.expandItemAndChildren(item, true)
					found = true
				})
				if found {
					cmd, err := m.updateItems()
					if err != nil {
						return m, tea.Quit
					}
					return m, cmd
				}
			} else if i.parent != nil {
				// For files, expand their parent directory and all its children
				parentPath := i.parent.path
				found := false
				m.traverseAndModify(&m.items, parentPath, func(item *item) {
					m.expandItemAndChildren(item, true)
					found = true
				})
				if found {
					cmd, err := m.updateItems()
					if err != nil {
						return m, tea.Quit
					}
					return m, cmd
				}
			}
		}
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

// expandItemAndChildren recursively sets the expanded state for a directory and all its children
func (m *BrowseModel) expandItemAndChildren(item *item, expanded bool) {
	if item.isDir {
		item.expanded = expanded
		for i := range item.children {
			m.expandItemAndChildren(&item.children[i], expanded)
		}
	}
}

func (m BrowseModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	if m.showNewMemoDialog {
		var sb strings.Builder

		// Create a styled new memo dialog
		dialogStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Width(60)

		// Get category tree from either the selected memo or the selected directory
		var categoryTree []string
		if m.selectedMemo != nil {
			categoryTree = m.selectedMemo.CategoryTree()
		} else {
			categoryTree = m.selectedCategoryTree
		}

		categoryPath := "root"
		if len(categoryTree) > 0 {
			categoryPath = strings.Join(categoryTree, " > ")
		}

		content := lipgloss.NewStyle().Bold(true).Render("Create New Memo") + "\n\n"
		content += lipgloss.NewStyle().Faint(true).Render("Category: ") + categoryPath + "\n\n"
		content += "Title: " + lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Render(m.newMemoTitleInput+"│") + "\n\n"
		content += lipgloss.NewStyle().Faint(true).Render("This will create a new memo in the selected category\nand open it in your editor.") + "\n\n"
		content += lipgloss.NewStyle().Italic(true).Render("Press Enter to create, Esc to cancel")

		sb.WriteString("\n")
		sb.WriteString(dialogStyle.Render(content))

		return sb.String()
	}

	if m.showDuplicateDialog {
		var sb strings.Builder

		// Create a styled duplicate confirmation dialog
		dialogStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Width(60)

		memoTitle := domain.MemoTitle(m.selectedMemo.FileName())
		currentTimestamp := m.selectedMemo.Date().Format("2006-01-02 15:04:05")
		newTimestamp := time.Now().Format("2006-01-02 15:04:05")

		content := lipgloss.NewStyle().Bold(true).Render("Duplicate Memo?") + "\n\n"
		content += fmt.Sprintf("Memo: %s\n", memoTitle)
		content += fmt.Sprintf("Original timestamp: %s\n", currentTimestamp)
		content += fmt.Sprintf("New timestamp:      %s\n\n", newTimestamp)
		content += lipgloss.NewStyle().Faint(true).Render("This will create a copy with the current timestamp.\nAll content and category will be preserved.\nTitle will have ' copied' appended.") + "\n\n"
		content += lipgloss.NewStyle().Bold(true).Render("Press Y to duplicate, N or Esc to cancel")

		sb.WriteString("\n")
		sb.WriteString(dialogStyle.Render(content))

		return sb.String()
	}

	if m.showRenameDialog {
		var sb strings.Builder

		// Create a styled rename dialog
		dialogStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Width(60)

		oldTitle := m.selectedMemo.Title()

		content := lipgloss.NewStyle().Bold(true).Render("Rename Memo") + "\n\n"
		content += lipgloss.NewStyle().Faint(true).Render("Current: ") + oldTitle + "\n\n"
		content += "New title: " + lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Render(m.renameInput+"│") + "\n\n"
		content += lipgloss.NewStyle().Faint(true).Render("This will update both the filename and title in content.") + "\n\n"
		content += lipgloss.NewStyle().Italic(true).Render("Press Enter to rename, Esc to cancel")

		sb.WriteString("\n")
		sb.WriteString(dialogStyle.Render(content))

		return sb.String()
	}

	if m.showDeleteDialog {
		var sb strings.Builder

		// Create a styled delete confirmation dialog
		dialogStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("196")).
			Padding(1, 2).
			Width(60)

		memoTitle := domain.MemoTitle(m.selectedMemo.FileName())

		content := "Delete Memo?\n\n"
		content += fmt.Sprintf("File: %s\n", memoTitle)
		content += fmt.Sprintf("Path: %s\n\n", filepath.Join(m.selectedMemo.Location(), m.selectedMemo.FileName()))
		content += "This will move the file to trash.\n\n"
		content += lipgloss.NewStyle().Bold(true).Render("Press Y to confirm, N or Esc to cancel")

		sb.WriteString("\n")
		sb.WriteString(dialogStyle.Render(content))

		return sb.String()
	}

	// Return split view if preview is enabled and not showing dialogs
	if !m.showNewMemoDialog && !m.showDuplicateDialog && !m.showRenameDialog &&
		!m.showDeleteDialog && !m.showCategoryDialog && m.showPreview {
		return m.renderSplitView()
	}

	if m.showCategoryDialog {
		// Show new category input prompt when in input mode
		if m.categoryDialogCursor == -1 {
			var sb strings.Builder

			// Create a styled input box
			inputStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(1, 2).
				Width(60)

			content := lipgloss.NewStyle().Bold(true).Render("Add New Category (Hierarchical)") + "\n\n"
			content += "Category path: " + lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Render(m.newCategoryInput+"│") + "\n\n"

			// Show examples
			examplesTitle := lipgloss.NewStyle().Underline(true).Render("Examples:")
			content += lipgloss.NewStyle().Faint(true).Render(examplesTitle + "\n")
			content += lipgloss.NewStyle().Faint(true).Render("  • work              → Creates 'work'\n")
			content += lipgloss.NewStyle().Faint(true).Render("  • work/projects     → Creates 'work' > 'projects'\n")
			content += lipgloss.NewStyle().Faint(true).Render("  • work/projects/2024 → Creates full tree\n\n")

			content += lipgloss.NewStyle().Faint(true).Render("Use / or > to separate levels. Existing folders are reused.") + "\n\n"
			content += lipgloss.NewStyle().Italic(true).Render("Press Enter to create, Esc to cancel")

			sb.WriteString("\n")
			sb.WriteString(inputStyle.Render(content))

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

	logger := common.DefaultLogger()
	repo := memo.NewMemo(m.config.MemosDir(), logger)
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

// pathToCategoryTree converts a directory path to a category tree
func (m *BrowseModel) pathToCategoryTree(dirPath string) []string {
	// Strip the memos directory prefix
	relPath := strings.TrimPrefix(dirPath, m.config.MemosDir())
	relPath = strings.TrimPrefix(relPath, string(filepath.Separator))

	// If empty (root directory), return empty slice
	if relPath == "" {
		return []string{}
	}

	// Split by path separator to get category tree
	return strings.Split(relPath, string(filepath.Separator))
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

// renderSplitView renders the browse list on the left and preview on the right
func (m BrowseModel) renderSplitView() string {
	// Calculate split widths (60/40 split)
	listWidth := int(float64(m.width) * 0.5)
	previewWidth := m.width - listWidth - 2 // -2 for padding

	// Update list width for split view
	m.list.SetWidth(listWidth)

	// Render the list
	listView := m.list.View()

	// Render the preview
	previewView := m.renderPreview(previewWidth, m.height)

	// Combine side by side
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		listView,
		previewView,
	)
}

// renderPreview generates the preview pane content
func (m BrowseModel) renderPreview(width, height int) string {
	previewStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		Width(width).
		Height(height - 2)

	// Get selected item
	selectedItem := m.list.SelectedItem()
	if selectedItem == nil {
		return previewStyle.Render("No item selected")
	}

	item, ok := selectedItem.(item)
	if !ok {
		return previewStyle.Render("Invalid item")
	}

	// If directory, show directory info
	if item.isDir {
		return m.renderDirectoryPreview(item, previewStyle)
	}

	// If memo, show memo preview
	if item.memo != nil {
		return m.renderMemoPreview(item, previewStyle)
	}

	return previewStyle.Render("No preview available")
}

// renderDirectoryPreview shows directory information
func (m BrowseModel) renderDirectoryPreview(item item, style lipgloss.Style) string {
	var content strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	content.WriteString(titleStyle.Render("📁 Directory") + "\n\n")
	content.WriteString(labelStyle.Render("Name: ") + item.name + "\n")
	content.WriteString(labelStyle.Render("Path: ") + item.path + "\n\n")

	// Count children
	fileCount := 0
	dirCount := 0
	for _, child := range item.children {
		if child.isDir {
			dirCount++
		} else {
			fileCount++
		}
	}

	content.WriteString(labelStyle.Render("Contents:") + "\n")
	content.WriteString(fmt.Sprintf("  Folders: %d\n", dirCount))
	content.WriteString(fmt.Sprintf("  Files: %d\n", fileCount))

	return style.Render(content.String())
}

// renderMemoPreview shows memo metadata and content preview
func (m BrowseModel) renderMemoPreview(item item, style lipgloss.Style) string {
	var content strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	headingStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230"))
	faintStyle := lipgloss.NewStyle().Faint(true)

	memo := item.memo

	// Title
	content.WriteString(titleStyle.Render("📝 "+domain.MemoTitle(memo.FileName())) + "\n\n")

	// Metadata section
	content.WriteString(headingStyle.Render("Metadata") + "\n")
	content.WriteString(labelStyle.Render("Date: ") + memo.Date().Format("2006-01-02 15:04") + "\n")

	if len(memo.CategoryTree()) > 0 {
		content.WriteString(labelStyle.Render("Category: ") + strings.Join(memo.CategoryTree(), " > ") + "\n")
	} else {
		content.WriteString(labelStyle.Render("Category: ") + faintStyle.Render("(root)") + "\n")
	}

	content.WriteString(labelStyle.Render("File: ") + memo.FileName() + "\n")

	// Content preview section - smart 8-line preview
	content.WriteString("\n" + headingStyle.Render("Content Preview") + "\n")

	var allLines []string

	// Start with top-level body content
	topLevel := memo.TopLevelBodyContent()
	if topLevel != nil && topLevel.ContentText != "" {
		topLevelText := strings.TrimSpace(topLevel.ContentText)
		topLines := strings.Split(topLevelText, "\n")
		allLines = append(allLines, topLines...)
	}

	// If we need more lines, peek into heading blocks
	if len(allLines) < 8 {
		headingBlocks := memo.HeadingBlocks()
		for _, hb := range headingBlocks {
			if len(allLines) >= 8 {
				break
			}

			// Add heading
			indent := strings.Repeat("  ", hb.Level-2)
			allLines = append(allLines, "")
			allLines = append(allLines, indent+"## "+hb.HeadingText)

			// Add content from this heading
			if hb.ContentText != "" {
				contentText := strings.TrimSpace(hb.ContentText)
				contentLines := strings.Split(contentText, "\n")
				for _, line := range contentLines {
					if len(allLines) >= 8 {
						break
					}
					allLines = append(allLines, indent+line)
				}
			}
		}
	}

	// Render the collected lines
	if len(allLines) > 0 {
		displayLines := allLines
		if len(displayLines) > 8 {
			displayLines = displayLines[:8]
			content.WriteString(strings.Join(displayLines, "\n") + "\n" + faintStyle.Render("..."))
		} else {
			content.WriteString(strings.Join(displayLines, "\n"))
		}
	} else {
		content.WriteString(faintStyle.Render("(no content)"))
	}

	// Show heading blocks section
	headingBlocks := memo.HeadingBlocks()
	if len(headingBlocks) > 0 {
		content.WriteString("\n\n" + headingStyle.Render("Sections") + "\n")
		maxHeadings := 5
		for i, hb := range headingBlocks {
			if i >= maxHeadings {
				content.WriteString(faintStyle.Render(fmt.Sprintf("  ... and %d more", len(headingBlocks)-maxHeadings)) + "\n")
				break
			}
			indent := strings.Repeat("  ", hb.Level-2)
			content.WriteString(indent + "• " + hb.HeadingText + "\n")
		}
	}

	return style.Render(content.String())
}
