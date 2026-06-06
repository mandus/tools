package tui

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"

	"github.com/mandu/tools/pass/cmd"
)

// item represents a password entry in the list
type item struct {
	path         string
	matchScore   int
	matchIndices []int
}

// Title returns the title of the item (the password path)
func (i item) Title() string { return i.path }

// Description returns the description of the item (empty for passwords)
func (i item) Description() string { return "" }

// FilterValue returns the value used for filtering
func (i item) FilterValue() string { return i.path }

// passwordDelegate is a custom delegate for rendering password items
type passwordDelegate struct {
	list.DefaultDelegate
}

// NewPasswordDelegate creates a new password delegate
func NewPasswordDelegate() passwordDelegate {
	return passwordDelegate{
		DefaultDelegate: list.NewDefaultDelegate(),
	}
}

// Height returns the height of the item
func (d passwordDelegate) Height() int { return 1 }

// Spacing returns the spacing between items
func (d passwordDelegate) Spacing() int { return 0 }

// Render renders a password item with match highlighting
func (d passwordDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	// Get the current selection
	isSelected := index == m.Index()

	// Create the base string
	var prefix string
	if isSelected {
		prefix = "> "
	} else {
		prefix = "  "
	}

	// Apply match highlighting if we have match indices
	displayPath := i.path
	if len(i.matchIndices) > 0 {
		displayPath = highlightMatches(i.path, i.matchIndices)
	}

	// Apply selection style if selected
	if isSelected {
		// Remove ANSI codes for width calculation
		plainPrefix := removeANSI(prefix)
		plainPath := removeANSI(displayPath)
		maxWidth := m.Width() - len(plainPrefix)
		
		// Truncate if necessary
		if len(plainPath) > maxWidth {
			plainPath = truncate(plainPath, maxWidth)
		}
		
		// Apply selected style to the whole line
		line := selectedStyle.Render(prefix + plainPath)
		fmt.Fprint(w, line)
	} else {
		// For non-selected items
		plainPrefix := removeANSI(prefix)
		plainPath := removeANSI(displayPath)
		maxWidth := m.Width() - len(plainPrefix)
		
		if len(plainPath) > maxWidth {
			plainPath = truncate(plainPath, maxWidth)
		}
		
		fmt.Fprint(w, prefix+plainPath)
	}
}

// highlightMatches highlights the matching characters in a path
func highlightMatches(path string, indices []int) string {
	if len(indices) == 0 {
		return path
	}

	// Sort indices
	sortedIndices := make([]int, len(indices))
	copy(sortedIndices, indices)
	sort.Ints(sortedIndices)

	var result strings.Builder
	prevIdx := 0

	for _, idx := range sortedIndices {
		// Add the non-matching part
		if idx > prevIdx {
			result.WriteString(path[prevIdx:idx])
		}
		// Add the matching character with highlighting
		if idx < len(path) {
			result.WriteString(matchStyle.Render(string(path[idx])))
			prevIdx = idx + 1
		}
	}

	// Add the remaining part
	if prevIdx < len(path) {
		result.WriteString(path[prevIdx:])
	}

	return result.String()
}

// removeANSI removes ANSI escape codes from a string
func removeANSI(s string) string {
	// Simple ANSI removal - this is a simplified version
	var result strings.Builder
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\u001b' {
			inEscape = true
			continue
		}
		if inEscape {
			if s[i] == 'm' {
				inEscape = false
			}
			continue
		}
		result.WriteByte(s[i])
	}
	return result.String()
}

// truncate truncates a string to the specified length
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	if length <= 0 {
		return ""
	}
	if length <= 3 {
		return s[:length]
	}
	return s[:length-3] + "..."
}

// Model is the main model for the fuzzy search TUI
type Model struct {
	// Bubble Tea list component
	list list.Model
	
	// Search input
	input textinput.Model
	
	// Current mode (using the cmd package type)
	mode cmd.FuzzySearchMode
	
	// All passwords (for filtering)
	allPasswords []string
	
	// State
	loading   bool
	error     error
	quitting  bool
	selected  string
	width     int
	height    int
}

// getTitle returns the title for the given mode
func getTitle(mode cmd.FuzzySearchMode) string {
	switch mode {
	case cmd.FuzzyModeShow:
		return "Select password (Enter to show, Esc to cancel)"
	case cmd.FuzzyModeClip:
		return "Select password to copy (Enter to copy, Esc to cancel)"
	case cmd.FuzzyModeRm:
		return "Select password to remove (Enter to delete, Esc to cancel)"
	default:
		return "Select password"
	}
}

// getPrompt returns the prompt for the given mode
func getPrompt(mode cmd.FuzzySearchMode) string {
	switch mode {
	case cmd.FuzzyModeShow:
		return "Search: "
	case cmd.FuzzyModeClip:
		return "Copy: "
	case cmd.FuzzyModeRm:
		return "Remove: "
	default:
		return "Search: "
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle window resize
		m.width = msg.Width
		m.height = msg.Height
		
		// Set list dimensions
		listWidth := msg.Width - 4 // Leave some padding
		if listWidth < 20 {
			listWidth = 20
		}
		listHeight := msg.Height - 8 // Account for header, input, help
		if listHeight < 5 {
			listHeight = 5
		}
		
		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)
		
		// Recreate input with new width
		m.input = recreateInput(m.input, msg.Width-10)
		
		return m, nil
	
	case tea.KeyMsg:
		// Handle keyboard input
		return m.handleKeyMsg(msg)
	}
	
	// Handle list-specific messages
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	
	return m, tea.Batch(cmds...)
}

// View renders the current state
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	
	if m.loading {
		return "Loading passwords...\n"
	}
	
	if m.error != nil {
		return "Error: " + m.error.Error() + "\n"
	}
	
	// Build the view
	var view string
	
	// Header
	view += getTitle(m.mode) + "\n\n"
	
	// Input prompt and value
	prompt := getPrompt(m.mode)
	view += prompt + m.input.View() + "\n\n"
	
	// List
	view += m.list.View() + "\n\n"
	
	// Help
	view += helpView()
	
	return view
}

// handleKeyMsg handles keyboard messages
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	// Exit keys
	case msg.String() == "esc" || msg.String() == "ctrl+c" || msg.String() == "ctrl+d" || msg.String() == "ctrl+q":
		m.quitting = true
		return m, tea.Quit
		
	// Enter - select
	case msg.String() == "enter":
		if m.list.SelectedItem() != nil {
			selectedItem := m.list.SelectedItem().(item)
			m.selected = selectedItem.path
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil
		
	// Handle input for search
	default:
		// Update the input
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		
		// Filter the list based on input
		m.filterList()
		
		return m, cmd
	}
}

// filterList filters the list based on the current input value
func (m *Model) filterList() {
	query := m.input.Value()
	
	// Use simple filtering for now
	var filtered []list.Item
	for _, path := range m.allPasswords {
		// Simple case-insensitive contains for now
		if strings.Contains(strings.ToLower(path), strings.ToLower(query)) {
			filtered = append(filtered, item{path: path})
		}
	}
	
	// If query is empty, show all passwords
	if query == "" {
		filtered = make([]list.Item, len(m.allPasswords))
		for i, path := range m.allPasswords {
			filtered[i] = item{path: path}
		}
	}
	
	m.list.SetItems(filtered)
}

// helpView returns the help text
func helpView() string {
	return "↑/↓: Navigate | Enter: Select | Esc/Ctrl+C: Cancel | Ctrl+Q: Quit"
}

// recreateInput recreates a textinput with a new width while preserving its value
func recreateInput(old textinput.Model, width int) textinput.Model {
	newInput := textinput.New()
	newInput.Prompt = old.Prompt
	newInput.Placeholder = old.Placeholder
	newInput.EchoMode = old.EchoMode
	newInput.EchoCharacter = old.EchoCharacter
	newInput.CharLimit = old.CharLimit
	newInput.Width = width
	newInput.SetValue(old.Value())
	newInput.Focus()
	return newInput
}

// NewModel creates a new fuzzy search model
func NewModel(passwords []string, mode cmd.FuzzySearchMode) *Model {
	// Create input
	input := textinput.New()
	input.Placeholder = "Type to search..."
	input.Focus()
	input.CharLimit = 100
	input.Width = 50
	
	// Create list
	items := make([]list.Item, len(passwords))
	for i, path := range passwords {
		items[i] = item{path: path}
	}
	
	delegate := NewPasswordDelegate()
	listModel := list.New(items, delegate, 0, 0)
	listModel.SetShowStatusBar(false)
	listModel.SetFilteringEnabled(false) // We handle filtering ourselves
	listModel.DisableQuitKeybindings()
	listModel.SetShowHelp(false)
	listModel.SetShowTitle(false)
	
	// Set initial dimensions
	listModel.SetWidth(76)
	listModel.SetHeight(15)
	
	// Create model
	model := &Model{
		list:         listModel,
		input:        input,
		mode:         mode,
		allPasswords: passwords,
		loading:      false,
		quitting:     false,
		width:        80,
		height:       24,
	}
	
	return model
}
