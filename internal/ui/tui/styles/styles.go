// Package styles centralizes the lipgloss palette and styles shared across the
// memo TUIs (browse, search, picker) so they present a consistent look.
//
// The palette is anchored on the browse view's existing accent (blue 63) and
// on-accent text (230); search and picker align to it.
package styles

import "github.com/charmbracelet/lipgloss"

// Palette colors (256-color ANSI).
const (
	ColorAccent   = lipgloss.Color("63")  // borders, headers, selection background
	ColorOnAccent = lipgloss.Color("230") // text drawn on the accent background
	ColorTitle    = lipgloss.Color("205") // memo titles / primary emphasis (pink)
	ColorMatch    = lipgloss.Color("86")  // matched substrings (cyan)
	ColorType     = lipgloss.Color("110") // labels / prompts (light blue)
	ColorDim      = lipgloss.Color("247") // secondary text
	ColorFaint    = lipgloss.Color("243") // tertiary text (counts, hints)
	ColorBorder   = lipgloss.Color("241") // inactive borders
)

// Shared styles.
var (
	// Selected is the base for the full-width selection bar.
	Selected = lipgloss.NewStyle().Background(ColorAccent).Foreground(ColorOnAccent).Bold(true)
	// Header titles a view or section.
	Header = lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
	// Title styles a memo title.
	Title = lipgloss.NewStyle().Foreground(ColorTitle)
	// Match styles a matched substring.
	Match = lipgloss.NewStyle().Foreground(ColorMatch)
	// Type styles a label or prompt.
	Type = lipgloss.NewStyle().Foreground(ColorType)
	// Dim styles secondary text.
	Dim = lipgloss.NewStyle().Foreground(ColorDim)
	// Faint styles tertiary text.
	Faint = lipgloss.NewStyle().Foreground(ColorFaint)

	// FocusedBorder / UnfocusedBorder frame focusable panes.
	FocusedBorder   = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(ColorAccent).Padding(0, 1)
	UnfocusedBorder = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(ColorBorder).Padding(0, 1)
)

// SelectedBar renders text as a full-width selection bar. Width values <= 0 fall
// back to a sensible default.
func SelectedBar(width int, text string) string {
	if width <= 0 {
		width = 80
	}
	return Selected.Width(width).Render(text)
}
