package picker

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func runes(s string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }

var (
	enter = tea.KeyMsg{Type: tea.KeyEnter}
	esc   = tea.KeyMsg{Type: tea.KeyEsc}
	down  = tea.KeyMsg{Type: tea.KeyDown}
	ctrlC = tea.KeyMsg{Type: tea.KeyCtrlC}
)

func send(m *Model, msgs ...tea.Msg) *Model {
	for _, msg := range msgs {
		updated, _ := m.Update(msg)
		m = updated.(*Model)
	}
	return m
}

func sampleItems() []Item {
	return []Item{
		{Display: "alpha", FilterBy: "alpha", Payload: "a"},
		{Display: "beta", FilterBy: "beta", Payload: "b"},
		{Display: "gamma", FilterBy: "gamma", Payload: "c"},
	}
}

func TestFilterNarrowsAndSelects(t *testing.T) {
	m := New(Config{Items: sampleItems()})
	if len(m.rows) != 3 {
		t.Fatalf("expected 3 rows before filtering, got %d", len(m.rows))
	}

	m = send(m, runes("bet"))
	if len(m.rows) != 1 {
		t.Fatalf("expected 1 row after filtering 'bet', got %d", len(m.rows))
	}

	m = send(m, enter)
	if m.chosen == nil {
		t.Fatal("expected an item to be chosen")
	}
	if m.chosen.Payload.(string) != "b" {
		t.Errorf("expected payload 'b', got %v", m.chosen.Payload)
	}
	if m.cancelled {
		t.Error("did not expect cancelled")
	}
}

func TestNavigationDown(t *testing.T) {
	m := New(Config{Items: sampleItems()})
	m = send(m, down, enter)
	if m.chosen == nil || m.chosen.Payload.(string) != "b" {
		t.Fatalf("expected second item 'b' selected, got %v", m.chosen)
	}
}

func TestWithInputPhase(t *testing.T) {
	m := New(Config{Items: sampleItems(), WithInput: true, InputPrompt: "Title: "})

	m = send(m, runes("alpha"), enter)
	if m.phase != phaseInput {
		t.Fatalf("expected to advance to input phase, got %v", m.phase)
	}

	// Empty input must not confirm.
	m = send(m, enter)
	if m.phase != phaseInput {
		t.Fatal("empty input should keep the input phase")
	}

	m = send(m, runes("New Title"), enter)
	if m.inputText != "New Title" {
		t.Errorf("expected input 'New Title', got %q", m.inputText)
	}
	if m.chosen == nil || m.chosen.Payload.(string) != "a" {
		t.Errorf("expected chosen 'a', got %v", m.chosen)
	}
}

func TestEscInInputReturnsToFilter(t *testing.T) {
	m := New(Config{Items: sampleItems(), WithInput: true, InputPrompt: "Title: "})
	m = send(m, runes("alpha"), enter, esc)
	if m.phase != phaseFilter {
		t.Fatalf("expected to return to filter phase, got %v", m.phase)
	}
	if m.chosen != nil {
		t.Error("expected chosen cleared after esc in input")
	}
}

func TestCancelInFilter(t *testing.T) {
	m := New(Config{Items: sampleItems()})
	m = send(m, esc)
	if !m.cancelled {
		t.Error("expected cancelled after esc")
	}

	m2 := New(Config{Items: sampleItems()})
	m2 = send(m2, ctrlC)
	if !m2.cancelled {
		t.Error("expected cancelled after ctrl+c")
	}
}

func TestEnterWithNoMatchesDoesNothing(t *testing.T) {
	m := New(Config{Items: sampleItems()})
	m = send(m, runes("zzz"), enter)
	if m.chosen != nil {
		t.Error("expected no selection when there are no matches")
	}
	if m.cancelled {
		t.Error("expected not cancelled")
	}
}

func TestFreeTextRowAppearsAndSelects(t *testing.T) {
	m := New(Config{Items: sampleItems(), AllowFreeText: true})

	// Typing something with no exact existing match adds a free-text row.
	m = send(m, runes("work/new"))
	if len(m.rows) != 1 || !m.rows[0].freeText {
		t.Fatalf("expected a single free-text row, got %+v", m.rows)
	}

	m = send(m, enter)
	if m.freeText != "work/new" {
		t.Errorf("expected freeText 'work/new', got %q", m.freeText)
	}
	if m.chosen != nil {
		t.Error("free-text selection should not set chosen")
	}
}

func TestFreeTextSuppressedOnExactMatch(t *testing.T) {
	m := New(Config{Items: sampleItems(), AllowFreeText: true})
	// "alpha" exactly matches an existing item: no free-text row, just the item.
	m = send(m, runes("alpha"))
	for _, r := range m.rows {
		if r.freeText {
			t.Fatal("did not expect a free-text row on an exact match")
		}
	}
	if len(m.rows) != 1 {
		t.Fatalf("expected 1 matching row, got %d", len(m.rows))
	}
}

func TestFreeTextDisabledByDefault(t *testing.T) {
	m := New(Config{Items: sampleItems()}) // AllowFreeText not set
	m = send(m, runes("zzz"), enter)
	if m.freeText != "" {
		t.Errorf("expected no free text when disabled, got %q", m.freeText)
	}
	if m.chosen != nil {
		t.Error("expected no selection")
	}
}

func TestRomajiMatcherFallbackSubstring(t *testing.T) {
	match := romajiMatcher(nil) // nil converter => plain per-word substring
	it := Item{FilterBy: "Weekly Meeting Notes"}
	if !match("meeting notes", it) {
		t.Error("expected multi-word substring match")
	}
	if match("absent", it) {
		t.Error("did not expect a match for absent word")
	}
	if !match("", it) {
		t.Error("empty query should match")
	}
}
