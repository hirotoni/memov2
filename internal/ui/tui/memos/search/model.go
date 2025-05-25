package search

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/config/toml"
	"github.com/hirotoni/memov2/internal/interfaces"
	"github.com/hirotoni/memov2/internal/repositories/memo"
	memsearch "github.com/hirotoni/memov2/internal/search"
	"golang.org/x/term"
)

type focusState int

const (
	focusInput focusState = iota
	focusList
)

type Model struct {
	searchInput textinput.Model
	viewport    viewport.Model
	config      *toml.Config
	editor      interfaces.Editor
	results     []memsearch.SearchResult
	selected    int
	width       int
	height      int
	focus       focusState
	romajiConv  *memsearch.RomajiConverter
	lastKeyG    bool // Track if last key was 'g' for 'gg' command
	err         error
}

var (
	titleStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	matchStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	dimStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("247"))
	pageStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	typeStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("110"))
	focusedStyle   = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62")).Padding(0, 1)
	unfocusedStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("241")).Padding(0, 1)
)

func New(c *toml.Config, e interfaces.Editor) (*Model, error) {
	// Get initial terminal dimensions
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w, h = 80, 20 // fallback dimensions
	}

	ti := textinput.New()
	ti.Placeholder = "Search memos..."
	ti.Focus()
	ti.Width = w

	vp := viewport.New(w, h-6)
	vp.Style = unfocusedStyle

	// Initialize romaji converter
	rc, err := memsearch.NewRomajiConverter()
	if err != nil {
		return nil, fmt.Errorf("failed to load SKK dictionary: %w", err)
	}

	m := &Model{
		searchInput: ti,
		viewport:    vp,
		config:      c,
		editor:      e,
		results:     []memsearch.SearchResult{},
		selected:    0,
		focus:       focusInput,
		romajiConv:  rc,
		lastKeyG:    false,
	}

	return m, nil
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Esc key behavior depends on focus state
			if m.focus == focusList {
				// Move focus back to input from list
				m.focus = focusInput
				m.searchInput.Focus()
				m.viewport.Style = unfocusedStyle
				return m, nil
			} else {
				// Exit app when Esc is pressed in input
				return m, tea.Quit
			}
		case "ctrl+j":
			// Toggle focus from input to list
			if m.focus == focusInput {
				m.focus = focusList
				m.searchInput.Blur()
				m.viewport.Style = focusedStyle
				return m, nil
			}
		case "ctrl+k":
			// Move focus back to input from list
			if m.focus == focusList {
				m.focus = focusInput
				m.searchInput.Focus()
				m.viewport.Style = unfocusedStyle
				return m, nil
			}
		case "enter":
			if m.focus == focusInput {
				// Move focus to list when Enter is pressed in input
				if len(m.results) > 0 {
					m.focus = focusList
					m.searchInput.Blur()
					m.viewport.Style = focusedStyle
					return m, nil
				}
			} else if m.focus == focusList && m.selected < len(m.results) {
				// Open selected memo when Enter is pressed in list
				result := m.results[m.selected]
				filePath := filepath.Join(m.config.MemosDir(), result.Memo.Location(), result.Memo.FileName())
				err := m.editor.Open(m.config.BaseDir(), filePath)
				if err != nil {
					m.err = fmt.Errorf("failed to open editor: %w", err)
					return m, nil
				}
				return m, nil
			}
		case "l":
			if m.focus == focusList && m.selected < len(m.results) {
				result := m.results[m.selected]
				filePath := filepath.Join(m.config.MemosDir(), result.Memo.Location(), result.Memo.FileName())
				err := m.editor.Open(m.config.BaseDir(), filePath)
				if err != nil {
					m.err = fmt.Errorf("failed to open editor: %w", err)
					return m, nil
				}
				return m, nil
			}
		}

		// Handle navigation only when list is focused
		if m.focus == focusList {
			switch msg.String() {
			case "g":
				if m.lastKeyG {
					// Second 'g' pressed - go to top
					m.selected = 0
					content := m.renderResults()
					m.viewport.SetContent(content)
					m.scrollToSelected(content)
					m.lastKeyG = false
				} else {
					// First 'g' pressed - wait for second 'g'
					m.lastKeyG = true
				}
			case "G", "shift+g":
				// Go to bottom
				if len(m.results) > 0 {
					m.selected = len(m.results) - 1
					content := m.renderResults()
					m.viewport.SetContent(content)
					m.scrollToSelected(content)
				}
				m.lastKeyG = false
			case "up", "k":
				if m.selected > 0 {
					m.selected--
					content := m.renderResults()
					m.viewport.SetContent(content)
					m.scrollToSelected(content)
				}
				m.lastKeyG = false
			case "down", "j":
				if m.selected < len(m.results)-1 {
					m.selected++
					content := m.renderResults()
					m.viewport.SetContent(content)
					m.scrollToSelected(content)
				}
				m.lastKeyG = false
			case "ctrl+u":
				// Move cursor up by 5 items
				if m.selected > 0 {
					m.selected = max(0, m.selected-5)
					content := m.renderResults()
					m.viewport.SetContent(content)
					m.scrollToSelected(content)
				}
				m.lastKeyG = false
			case "ctrl+d":
				// Move cursor down by 5 items
				if m.selected < len(m.results)-1 {
					m.selected = min(len(m.results)-1, m.selected+5)
					content := m.renderResults()
					m.viewport.SetContent(content)
					m.scrollToSelected(content)
				}
				m.lastKeyG = false
			default:
				m.lastKeyG = false
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 6
		m.searchInput.Width = msg.Width - 11

		// Update viewport content
		if len(m.results) > 0 {
			content := m.renderResults()
			m.viewport.SetContent(content)
			m.scrollToSelected(content)
		}

	case searchResultMsg:
		m.results = msg.results
		m.selected = 0
		content := m.renderResults()
		m.viewport.SetContent(content)
		m.viewport.GotoTop() // Reset scroll position when new results arrive
		return m, nil
	}

	if m.focus == focusInput {
		m.searchInput, tiCmd = m.searchInput.Update(msg)
		if m.searchInput.Value() != "" {
			return m, tea.Batch(tiCmd, m.search)
		}
		return m, tiCmd
	}

	m.viewport, vpCmd = m.viewport.Update(msg)
	return m, vpCmd
}

func (m Model) View() string {
	var s strings.Builder

	// Search input with focus indicator
	inputStyle := unfocusedStyle
	if m.focus == focusInput {
		inputStyle = focusedStyle
	}
	s.WriteString("\n" + inputStyle.Render(m.searchInput.View()) + "\n\n")

	// Results viewport
	s.WriteString(m.viewport.View())

	// Help
	var help string
	if m.focus == focusInput {
		help = "\n  tab/ctrl+j: switch to list | type to search | esc: quit\n"
	} else {
		help = "\n  ↑/k,↓/j: navigate | ctrl+u/ctrl+d: move by 5 | gg/G: top/bottom | enter/l: open | tab/ctrl+k/esc: switch to search | ctrl+c: quit\n"
	}
	s.WriteString(dimStyle.Render(help))

	return s.String()
}

type searchResultMsg struct {
	results []memsearch.SearchResult
}

func (m *Model) search() tea.Msg {
	query := m.searchInput.Value()
	if query == "" {
		return searchResultMsg{results: []memsearch.SearchResult{}}
	}

	logger := common.DefaultLogger()
	memos, err := memo.NewMemo(m.config.MemosDir(), logger).MemoEntries()
	if err != nil {
		return searchResultMsg{results: []memsearch.SearchResult{}}
	}

	// Map to store results by file path to merge matches from same file
	resultsByFile := make(map[string]*memsearch.SearchResult)

	// Split query into words and convert each word
	queryWords := strings.Fields(query)
	wordQueries := make([][]string, len(queryWords))

	// Convert each word and collect all possible variations
	for i, word := range queryWords {
		if m.romajiConv != nil {
			wordQueries[i] = m.romajiConv.Convert(word)
		} else {
			wordQueries[i] = []string{word}
		}
	}

	// Search with all query candidates
	for _, memo := range memos {
		filePath := filepath.Join(memo.Location(), memo.FileName())

		// Initialize result for this file
		result := &memsearch.SearchResult{
			Memo:    memo,
			Matches: []memsearch.Match{},
		}

		// Track matches by type
		matchesByType := make(map[memsearch.MatchType][]memsearch.Match)

		// For each word in the query
		for _, wordVariations := range wordQueries {
			// Try each variation of the current word
			for _, q := range wordVariations {
				// Search by title
				if matches := memsearch.SearchMemo(memo, q, memsearch.MatchTitle); len(matches.Matches) > 0 {
					matchesByType[memsearch.MatchTitle] = append(matchesByType[memsearch.MatchTitle], matches.Matches...)
				}

				// Search by category
				if matches := memsearch.SearchMemo(memo, q, memsearch.MatchCategory); len(matches.Matches) > 0 {
					matchesByType[memsearch.MatchCategory] = append(matchesByType[memsearch.MatchCategory], matches.Matches...)
				}

				// Search by heading
				if matches := memsearch.SearchMemo(memo, q, memsearch.MatchHeading); len(matches.Matches) > 0 {
					matchesByType[memsearch.MatchHeading] = append(matchesByType[memsearch.MatchHeading], matches.Matches...)
				}

				// Search by content
				if matches := memsearch.SearchMemo(memo, q, memsearch.MatchContent); len(matches.Matches) > 0 {
					matchesByType[memsearch.MatchContent] = append(matchesByType[memsearch.MatchContent], matches.Matches...)
				}
			}
		}

		// Check each match type independently
		var validMatches []memsearch.Match
		hasValidMatches := false

		// Helper function to check if all query words appear in a text
		containsAllWords := func(text string) bool {
			textLower := strings.ToLower(text)
			for _, wordVariations := range wordQueries {
				wordFound := false
				for _, word := range wordVariations {
					if strings.Contains(textLower, strings.ToLower(word)) {
						wordFound = true
						break
					}
				}
				if !wordFound {
					return false
				}
			}
			return true
		}

		// Process each match type
		for _, matches := range matchesByType {
			// Group matches by their unique content
			uniqueMatches := make(map[string]memsearch.Match)
			for _, match := range matches {
				uniqueMatches[match.Content] = match
			}

			// Check each unique match
			for _, match := range uniqueMatches {
				if containsAllWords(match.Content) {
					validMatches = append(validMatches, match)
					hasValidMatches = true
				}
			}
		}

		// Only keep results that have valid matches
		if hasValidMatches {
			result.Matches = validMatches
			resultsByFile[filePath] = result
		}
	}

	// Convert map to slice and sort results
	var allResults []memsearch.SearchResult
	for _, result := range resultsByFile {
		// Sort matches within each result by type and position
		sort.Slice(result.Matches, func(i, j int) bool {
			// First sort by match type
			if result.Matches[i].Type != result.Matches[j].Type {
				return result.Matches[i].Type < result.Matches[j].Type
			}

			// Then sort by heading order for all match types
			if result.Matches[i].HeadingOrder != result.Matches[j].HeadingOrder {
				return result.Matches[i].HeadingOrder < result.Matches[j].HeadingOrder
			}

			// For content matches, sort by line number within heading
			if result.Matches[i].Type == memsearch.MatchContent {
				return result.Matches[i].Line < result.Matches[j].Line
			}

			// For heading matches at the same position (should be rare), sort by content for stability
			return result.Matches[i].Content < result.Matches[j].Content
		})
		allResults = append(allResults, *result)
	}

	// Sort results by date (newest first) and then by file path for stability
	sort.Slice(allResults, func(i, j int) bool {
		dateI := allResults[i].Memo.Date()
		dateJ := allResults[j].Memo.Date()
		if !dateI.Equal(dateJ) {
			return dateI.After(dateJ)
		}
		// If dates are equal, sort by file path for stability
		pathI := filepath.Join(allResults[i].Memo.Location(), allResults[i].Memo.FileName())
		pathJ := filepath.Join(allResults[j].Memo.Location(), allResults[j].Memo.FileName())
		return pathI < pathJ
	})

	return searchResultMsg{results: allResults}
}

// matchInfo stores information about where matches occur in text
type matchInfo struct {
	start int
	end   int
}

// getQueryVariations returns all variations of the search query words
func (m *Model) getQueryVariations() []string {
	if m.searchInput.Value() == "" {
		return nil
	}

	queryWords := strings.Fields(m.searchInput.Value())
	var allQueries []string
	seen := make(map[string]bool) // For deduplication

	for _, word := range queryWords {
		if m.romajiConv != nil {
			variations := m.romajiConv.Convert(word)
			for _, v := range variations {
				if !seen[v] {
					allQueries = append(allQueries, v)
					seen[v] = true
				}
			}
		} else if !seen[word] {
			allQueries = append(allQueries, word)
			seen[word] = true
		}
	}
	return allQueries
}

// findAllMatches finds all match positions in the text for given queries
func findAllMatches(text string, queries []string) []matchInfo {
	if len(queries) == 0 {
		return nil
	}

	textLower := strings.ToLower(text)
	var matches []matchInfo

	// Find matches for each query
	for _, query := range queries {
		queryLower := strings.ToLower(query)
		pos := 0
		for {
			idx := strings.Index(textLower[pos:], queryLower)
			if idx == -1 {
				break
			}
			idx += pos // Adjust index for the full string
			matches = append(matches, matchInfo{
				start: idx,
				end:   idx + len(queryLower),
			})
			pos = idx + len(queryLower)
		}
	}

	// Sort matches by start position
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].start < matches[j].start
	})

	// Merge overlapping matches
	if len(matches) <= 1 {
		return matches
	}

	merged := make([]matchInfo, 0, len(matches))
	current := matches[0]

	for i := 1; i < len(matches); i++ {
		if matches[i].start <= current.end {
			// Merge overlapping matches
			if matches[i].end > current.end {
				current.end = matches[i].end
			}
		} else {
			merged = append(merged, current)
			current = matches[i]
		}
	}
	merged = append(merged, current)

	return merged
}

// highlightMatches highlights the matched portions of text based on the search query
func (m *Model) highlightMatches(text string, queryVariations []string, style lipgloss.Style) string {
	if len(queryVariations) == 0 {
		return style.Render(text)
	}

	// Find all matches
	matches := findAllMatches(text, queryVariations)
	if len(matches) == 0 {
		return style.Render(text)
	}

	// Build the result string with highlights
	var result strings.Builder
	result.Grow(len(text) * 2) // Preallocate space for efficiency

	lastPos := 0
	for _, match := range matches {
		// Add text before match with original style
		result.WriteString(style.Render(text[lastPos:match.start]))
		// Add highlighted match
		result.WriteString(matchStyle.Render(text[match.start:match.end]))
		lastPos = match.end
	}
	// Add remaining text with original style
	result.WriteString(style.Render(text[lastPos:]))

	return result.String()
}

func (m Model) renderResults() string {
	var s strings.Builder
	s.Grow(4096) // Preallocate initial capacity

	if len(m.results) == 0 {
		s.WriteString("No results found")
		return s.String()
	}

	// Cache query variations
	queryVariations := m.getQueryVariations()

	// Results count
	resultInfo := fmt.Sprintf("Found %d results", len(m.results))
	s.WriteString(pageStyle.Render(resultInfo))
	s.WriteString("\n\n")

	// Render all results
	for i, result := range m.results {
		// Highlight selected result
		prefix := " "
		if i == m.selected {
			prefix = "▸"
		}

		// Memo title and category
		title := m.highlightMatches(result.Memo.Title(), queryVariations, titleStyle)
		fmt.Fprintf(&s, "%s %s\n", prefix, title)
		datetime := result.Memo.Date().Format("2006-01-02 15:04")
		datetime = dimStyle.Render(datetime)
		category := ""
		if len(result.Memo.CategoryTree()) > 0 {
			categoryText := strings.Join(result.Memo.CategoryTree(), " > ")
			category = m.highlightMatches(categoryText, queryVariations, dimStyle)
		}
		fmt.Fprintf(&s, "   %s | %s\n", datetime, category)

		// Group matches by type and deduplicate based on content
		matchesByType := make(map[memsearch.MatchType]map[string]memsearch.Match)
		for _, match := range result.Matches {
			if matchesByType[match.Type] == nil {
				matchesByType[match.Type] = make(map[string]memsearch.Match)
			}
			// Use content as key for deduplication
			matchesByType[match.Type][match.Content] = match
		}

		// Display matches grouped by type
		matchTypes := []struct {
			Type  memsearch.MatchType
			Label string
		}{
			{memsearch.MatchTitle, "Title"},
			{memsearch.MatchCategory, "Category"},
			{memsearch.MatchHeading, "Heading"},
			{memsearch.MatchContent, "Content"},
		}

		for _, mt := range matchTypes {
			matches := matchesByType[mt.Type]
			if len(matches) > 0 {
				fmt.Fprintf(&s, "   %s matches:\n", typeStyle.Render(mt.Label))

				// Convert matches map to slice and sort by line number for stable order
				var sortedMatches []memsearch.Match
				for _, match := range matches {
					sortedMatches = append(sortedMatches, match)
				}
				sort.Slice(sortedMatches, func(i, j int) bool {
					if sortedMatches[i].HeadingOrder != sortedMatches[j].HeadingOrder {
						return sortedMatches[i].HeadingOrder < sortedMatches[j].HeadingOrder
					}
					return sortedMatches[i].Line < sortedMatches[j].Line
				})

				if mt.Type == memsearch.MatchContent {
					// Group content matches by heading
					currentHeading := ""
					for i, match := range sortedMatches {
						if match.Heading != currentHeading {
							currentHeading = match.Heading
							if i > 0 {
								s.WriteString("\n") // Add extra newline between different headings
							}
							fmt.Fprintf(&s, "     In heading: %s\n", m.highlightMatches(currentHeading, queryVariations, dimStyle))
						}
						// Show context and matched line
						if match.PrevLineContext != "" {
							fmt.Fprintf(&s, "           %d: %s\n", match.Line-1, dimStyle.Render(match.PrevLineContext))
						}
						fmt.Fprintf(&s, "       >>> %d: %s\n", match.Line, m.highlightMatches(match.Content, queryVariations, lipgloss.NewStyle()))
						if match.NextLineContext != "" {
							fmt.Fprintf(&s, "           %d: %s\n", match.Line+1, dimStyle.Render(match.NextLineContext))
						}

						// Add separator between matches in the same heading, but not after the last match
						if i < len(sortedMatches)-1 && match.Heading == sortedMatches[i+1].Heading {
							fmt.Fprintf(&s, "           %s\n", dimStyle.Render("---"))
						}
					}
				} else {
					for _, match := range sortedMatches {
						switch match.Type {
						case memsearch.MatchTitle, memsearch.MatchCategory, memsearch.MatchHeading:
							fmt.Fprintf(&s, "     %s\n", m.highlightMatches(match.Content, queryVariations, lipgloss.NewStyle()))
						}
					}
				}
			}
		}
		s.WriteString("\n")
	}

	return s.String()
}

// scrollToSelected ensures the selected item is visible in the viewport
func (m *Model) scrollToSelected(content string) {
	lines := strings.Split(content, "\n")

	// Find the selected item's position
	selectedLine := -1
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "▸") {
			selectedLine = i
			break
		}
	}

	if selectedLine == -1 {
		return // No selected item found
	}

	// Calculate the ideal position for the selected item
	// We want it to be roughly in the middle of the viewport when possible
	idealOffset := max(0, selectedLine-m.viewport.Height/2)

	// Adjust viewport position
	m.viewport.YOffset = idealOffset

	// Ensure we don't scroll past the bottom
	maxOffset := max(0, len(lines)-m.viewport.Height)
	m.viewport.YOffset = min(m.viewport.YOffset, maxOffset)
}

type SearchResultMsg struct {
	Results []memsearch.SearchResult
}

// RenderResults renders the search results in a formatted string
func RenderResults(results []memsearch.SearchResult, selected int, query string, romajiConv *memsearch.RomajiConverter) string {
	var s strings.Builder
	s.Grow(4096) // Preallocate initial capacity

	if len(results) == 0 {
		s.WriteString("No results found")
		return s.String()
	}

	// Results count
	resultInfo := fmt.Sprintf("Found %d results", len(results))
	s.WriteString(pageStyle.Render(resultInfo))
	s.WriteString("\n\n")

	// Cache query variations
	var queryVariations []string
	if query != "" {
		if romajiConv != nil {
			queryWords := strings.Fields(query)
			for _, word := range queryWords {
				variations := romajiConv.Convert(word)
				queryVariations = append(queryVariations, variations...)
			}
		} else {
			queryVariations = []string{query}
		}
	}

	// Render all results
	for i, result := range results {
		// Highlight selected result
		prefix := " "
		if i == selected {
			prefix = "▸"
		}

		// Memo title and category
		title := highlightMatches(result.Memo.Title(), queryVariations, titleStyle)
		fmt.Fprintf(&s, "%s %s\n", prefix, title)
		datetime := result.Memo.Date().Format("2006-01-02 15:04")
		datetime = dimStyle.Render(datetime)
		category := ""
		if len(result.Memo.CategoryTree()) > 0 {
			categoryText := strings.Join(result.Memo.CategoryTree(), " > ")
			category = highlightMatches(categoryText, queryVariations, dimStyle)
		}
		fmt.Fprintf(&s, "   %s | %s\n", datetime, category)

		// Group matches by type and deduplicate based on content
		matchesByType := make(map[memsearch.MatchType]map[string]memsearch.Match)
		for _, match := range result.Matches {
			if matchesByType[match.Type] == nil {
				matchesByType[match.Type] = make(map[string]memsearch.Match)
			}
			// Use content as key for deduplication
			matchesByType[match.Type][match.Content] = match
		}

		// Display matches grouped by type
		matchTypes := []struct {
			Type  memsearch.MatchType
			Label string
		}{
			{memsearch.MatchTitle, "Title"},
			{memsearch.MatchCategory, "Category"},
			{memsearch.MatchHeading, "Heading"},
			{memsearch.MatchContent, "Content"},
		}

		for _, mt := range matchTypes {
			matches := matchesByType[mt.Type]
			if len(matches) > 0 {
				fmt.Fprintf(&s, "   %s matches:\n", typeStyle.Render(mt.Label))

				// Convert matches map to slice and sort by line number for stable order
				var sortedMatches []memsearch.Match
				for _, match := range matches {
					sortedMatches = append(sortedMatches, match)
				}
				sort.Slice(sortedMatches, func(i, j int) bool {
					if sortedMatches[i].HeadingOrder != sortedMatches[j].HeadingOrder {
						return sortedMatches[i].HeadingOrder < sortedMatches[j].HeadingOrder
					}
					return sortedMatches[i].Line < sortedMatches[j].Line
				})

				if mt.Type == memsearch.MatchContent {
					// Group content matches by heading
					currentHeading := ""
					for i, match := range sortedMatches {
						if match.Heading != currentHeading {
							currentHeading = match.Heading
							if i > 0 {
								s.WriteString("\n") // Add extra newline between different headings
							}
							fmt.Fprintf(&s, "     In heading: %s\n", highlightMatches(currentHeading, queryVariations, dimStyle))
						}
						// Show context and matched line
						if match.PrevLineContext != "" {
							fmt.Fprintf(&s, "           %d: %s\n", match.Line-1, dimStyle.Render(match.PrevLineContext))
						}
						fmt.Fprintf(&s, "       >>> %d: %s\n", match.Line, highlightMatches(match.Content, queryVariations, lipgloss.NewStyle()))
						if match.NextLineContext != "" {
							fmt.Fprintf(&s, "           %d: %s\n", match.Line+1, dimStyle.Render(match.NextLineContext))
						}

						// Add separator between matches in the same heading, but not after the last match
						if i < len(sortedMatches)-1 && match.Heading == sortedMatches[i+1].Heading {
							fmt.Fprintf(&s, "           %s\n", dimStyle.Render("---"))
						}
					}
				} else {
					for _, match := range sortedMatches {
						switch match.Type {
						case memsearch.MatchTitle, memsearch.MatchCategory, memsearch.MatchHeading:
							fmt.Fprintf(&s, "     %s\n", highlightMatches(match.Content, queryVariations, lipgloss.NewStyle()))
						}
					}
				}
			}
		}
		s.WriteString("\n")
	}

	return s.String()
}

// highlightMatches highlights the matched portions of text based on the search query
func highlightMatches(text string, queryVariations []string, style lipgloss.Style) string {
	if len(queryVariations) == 0 {
		return style.Render(text)
	}

	// Find all matches
	matches := findAllMatches(text, queryVariations)
	if len(matches) == 0 {
		return style.Render(text)
	}

	// Build the result string with highlights
	var result strings.Builder
	result.Grow(len(text) * 2) // Preallocate space for efficiency

	lastPos := 0
	for _, match := range matches {
		// Add text before match with original style
		result.WriteString(style.Render(text[lastPos:match.start]))
		// Add highlighted match
		result.WriteString(matchStyle.Render(text[match.start:match.end]))
		lastPos = match.end
	}
	// Add remaining text with original style
	result.WriteString(style.Render(text[lastPos:]))

	return result.String()
}
