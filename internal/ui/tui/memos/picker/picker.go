// Package picker provides a small, reusable two-phase selection TUI.
//
// Phase 1 (filter): the user types into a filter input and navigates a list of
// candidate items. Phase 2 (input, optional): after selecting an item, the user
// types free text (e.g. a new title) to confirm the action.
//
// It is the embedded-TUI building block behind the `memov2 alt` commands and is
// intentionally decoupled from any memo-specific logic: callers supply the items
// and a match function. See memo_picker.go and category_picker.go for the
// memo/category-specific wiring.
package picker

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hirotoni/memov2/internal/ui/tui/styles"
	"golang.org/x/term"
)

// Item is a single selectable candidate.
type Item struct {
	// Display is the primary text shown for this item.
	Display string
	// Secondary is optional dim text shown after Display (e.g. a path).
	Secondary string
	// FilterBy is the text matched against the query (case handling is up to Match).
	FilterBy string
	// Payload is the value returned to the caller when this item is chosen.
	Payload any
}

// MatchFunc reports whether item matches query. An empty query should match all.
type MatchFunc func(query string, item Item) bool

// Config configures a picker run.
type Config struct {
	Title       string    // header shown above the filter input
	Items       []Item    // candidates
	Match       MatchFunc // filter predicate; nil => case-insensitive substring on FilterBy
	WithInput   bool      // if true, a text-input phase follows selection
	InputPrompt string    // label for the input phase (e.g. "New title: ")

	// AllowFreeText offers the typed query itself as a selectable row when it
	// does not exactly match an existing item, so the caller can accept a new
	// value (e.g. a brand-new category). FreeTextLabel renders that row; if nil,
	// a default `+ create "<query>"` label is used.
	AllowFreeText bool
	FreeTextLabel func(query string) string
}

// Result is the outcome of a picker run.
type Result struct {
	Item      *Item  // the chosen existing item (nil if cancelled or free text)
	InputText string // text entered in the input phase (empty if WithInput is false)
	FreeText  string // the typed query, when a free-text row was chosen
	Cancelled bool   // true if the user aborted (esc/ctrl+c)
}

type phase int

const (
	phaseFilter phase = iota
	phaseInput
)

// row is a single visible line in the filter list: either an existing item or
// the synthetic free-text row.
type row struct {
	freeText  bool
	itemIndex int // index into cfg.Items, valid when !freeText
}

// Model is the picker Bubbletea model. After Run, read Result.
type Model struct {
	cfg         Config
	filterInput textinput.Model
	textInput   textinput.Model
	match       MatchFunc

	rows     []row // visible rows for the current query, in display order
	selected int   // index into rows
	phase    phase
	width    int
	height   int

	// result
	chosen    *Item
	freeText  string
	inputText string
	cancelled bool
}

func defaultMatch(query string, item Item) bool {
	if query == "" {
		return true
	}
	text := strings.ToLower(item.FilterBy)
	for _, word := range strings.Fields(query) {
		if !strings.Contains(text, strings.ToLower(word)) {
			return false
		}
	}
	return true
}

// New builds a picker model from cfg.
func New(cfg Config) *Model {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w, h = 80, 24
	}

	fi := textinput.New()
	fi.Placeholder = "Filter..."
	fi.Focus()
	fi.Width = w

	ti := textinput.New()
	ti.Placeholder = strings.TrimSpace(cfg.InputPrompt)
	ti.Width = w

	match := cfg.Match
	if match == nil {
		match = defaultMatch
	}

	m := &Model{
		cfg:         cfg,
		filterInput: fi,
		textInput:   ti,
		match:       match,
		phase:       phaseFilter,
		width:       w,
		height:      h,
	}
	m.recompute()
	return m
}

// recompute rebuilds the visible rows for the current query and clamps the
// selection. When free text is offered, its row is appended at the bottom.
func (m *Model) recompute() {
	q := m.filterInput.Value()
	m.rows = m.rows[:0]

	exactMatch := false
	for i, it := range m.cfg.Items {
		if m.match(q, it) {
			m.rows = append(m.rows, row{itemIndex: i})
		}
		if it.Display == strings.TrimSpace(q) {
			exactMatch = true
		}
	}

	if m.cfg.AllowFreeText && strings.TrimSpace(q) != "" && !exactMatch {
		m.rows = append(m.rows, row{freeText: true})
	}

	if m.selected >= len(m.rows) {
		m.selected = len(m.rows) - 1
	}
	if m.selected < 0 {
		m.selected = 0
	}
}

// freeTextLabel renders the synthetic free-text row for the current query.
func (m *Model) freeTextLabel() string {
	q := strings.TrimSpace(m.filterInput.Value())
	if m.cfg.FreeTextLabel != nil {
		return m.cfg.FreeTextLabel(q)
	}
	return "+ create \"" + q + "\""
}

func (m *Model) Init() tea.Cmd { return textinput.Blink }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.phase == phaseInput {
			return m.updateInput(msg)
		}
		return m.updateFilter(msg)
	}
	return m, nil
}

func (m *Model) updateFilter(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.cancelled = true
		return m, tea.Quit
	case "up", "ctrl+p":
		if m.selected > 0 {
			m.selected--
		}
		return m, nil
	case "down", "ctrl+n":
		if m.selected < len(m.rows)-1 {
			m.selected++
		}
		return m, nil
	case "enter":
		if len(m.rows) == 0 {
			return m, nil
		}
		r := m.rows[m.selected]
		if r.freeText {
			m.freeText = strings.TrimSpace(m.filterInput.Value())
		} else {
			it := m.cfg.Items[r.itemIndex]
			m.chosen = &it
		}
		if m.cfg.WithInput {
			m.phase = phaseInput
			m.filterInput.Blur()
			m.textInput.Focus()
			return m, textinput.Blink
		}
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.filterInput, cmd = m.filterInput.Update(msg)
	m.recompute()
	return m, cmd
}

func (m *Model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.chosen = nil
		m.freeText = ""
		m.cancelled = true
		return m, tea.Quit
	case "esc":
		// Go back to selection rather than aborting outright.
		m.chosen = nil
		m.freeText = ""
		m.phase = phaseFilter
		m.textInput.SetValue("")
		m.textInput.Blur()
		m.filterInput.Focus()
		return m, textinput.Blink
	case "enter":
		text := strings.TrimSpace(m.textInput.Value())
		if text == "" {
			return m, nil // require non-empty input
		}
		m.inputText = text
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	if m.phase == phaseInput {
		return m.viewInput()
	}
	return m.viewFilter()
}

func (m *Model) viewFilter() string {
	var b strings.Builder
	if m.cfg.Title != "" {
		b.WriteString(styles.Header.Render(m.cfg.Title) + "\n")
	}
	b.WriteString(m.filterInput.View() + "\n\n")

	// Window the visible list around the selection.
	listHeight := m.height - 6
	if listHeight < 1 {
		listHeight = 1
	}
	start := 0
	if m.selected >= listHeight {
		start = m.selected - listHeight + 1
	}
	end := start + listHeight
	if end > len(m.rows) {
		end = len(m.rows)
	}

	for i := start; i < end; i++ {
		r := m.rows[i]
		var display, secondary string
		if r.freeText {
			display = m.freeTextLabel()
		} else {
			it := m.cfg.Items[r.itemIndex]
			display, secondary = it.Display, it.Secondary
		}

		if i == m.selected {
			text := "▸ " + display
			if secondary != "" {
				text += "  " + secondary
			}
			b.WriteString(styles.SelectedBar(m.width, text) + "\n")
		} else {
			line := "  " + display
			if secondary != "" {
				line += "  " + styles.Dim.Render(secondary)
			}
			b.WriteString(line + "\n")
		}
	}
	if len(m.rows) == 0 {
		b.WriteString(styles.Dim.Render("  (no matches)") + "\n")
	}

	b.WriteString("\n" + styles.Dim.Render("↑/ctrl+p, ↓/ctrl+n: navigate | type to filter | enter: select | esc: cancel"))
	return b.String()
}

func (m *Model) viewInput() string {
	var b strings.Builder
	if m.cfg.Title != "" {
		b.WriteString(styles.Header.Render(m.cfg.Title) + "\n")
	}
	if m.chosen != nil {
		b.WriteString(styles.Dim.Render("Selected: "+m.chosen.Display) + "\n\n")
	} else if m.freeText != "" {
		b.WriteString(styles.Dim.Render("New: "+m.freeText) + "\n\n")
	}
	b.WriteString(styles.Type.Render(m.cfg.InputPrompt) + "\n")
	b.WriteString(m.textInput.View() + "\n\n")
	b.WriteString(styles.Dim.Render("enter: confirm | esc: back | ctrl+c: cancel"))
	return b.String()
}

// Run launches the picker and returns its Result.
func Run(cfg Config) (Result, error) {
	m := New(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return Result{Cancelled: true}, err
	}
	fm, ok := final.(*Model)
	if !ok {
		return Result{Cancelled: true}, nil
	}
	return Result{
		Item:      fm.chosen,
		InputText: fm.inputText,
		FreeText:  fm.freeText,
		Cancelled: fm.cancelled || (fm.chosen == nil && fm.freeText == ""),
	}, nil
}
