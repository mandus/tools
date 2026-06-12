package tui

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mandu/tools/pass/pkg/fuzzy"
	"github.com/mandu/tools/pass/pkg/git"
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
	
	// Current mode
	mode FuzzySearchMode
	
	// All passwords (for filtering)
	allPasswords []string
	
	// Git status
	gitStatus    git.GitStatus
	gitStatusErr error
	
	// State
	loading   bool
	error     error
	quitting  bool
	selected  string
	width     int
	height    int
}

// getTitle returns the title for the given mode
func getTitle(mode FuzzySearchMode) string {
	switch mode {
	case FuzzyModeShow:
		return "Select password (Enter to show, Esc to cancel)"
	case FuzzyModeClip:
		return "Select password to copy (Enter to copy, Esc to cancel)"
	case FuzzyModeRm:
		return "Select password to remove (Enter to delete, Esc to cancel)"
	case FuzzyModeEdit:
		return "Select password to edit (Enter to edit, Esc to cancel)"
	default:
		return "Select password"
	}
}

// getPrompt returns the prompt for the given mode
func getPrompt(mode FuzzySearchMode) string {
	switch mode {
	case FuzzyModeShow:
		return "Search: "
	case FuzzyModeClip:
		return "Copy: "
	case FuzzyModeRm:
		return "Remove: "
	case FuzzyModeEdit:
		return "Edit: "
	default:
		return "Search: "
	}
}

// getGitStatusLine returns a formatted git status line for display with colors
func getGitStatusLine(status git.GitStatus) string {
	if !status.IsGitRepo {
		return ""
	}
	
	// Use the String() method which already has the simple symbols
	// But we'll add colors to make it fancy
	base := status.String()
	
	if base == "" || !status.IsGitRepo {
		return ""
	}
	
	// Color code the status based on state
	// Green for up to date, Yellow for ahead/behind, Red for conflicts
	if status.HasMergeConflict {
		// Red for conflicts
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("Git: " + base)
	} else if status.Ahead > 0 || status.Behind > 0 {
		// Yellow for not in sync
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render("Git: " + base)
	} else if status.HasUncommitted {
		// Yellow for uncommitted changes
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render("Git: " + base)
	} else {
		// Green for up to date
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("Git: " + base)
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		
	case gitPushResult:
		// Handle git push result
		if msg.err != nil {
			m.error = fmt.Errorf("git push failed: %v", msg.err)
		} else {
			// Refresh git status after successful push
			m.gitStatus = git.GetGitStatus(getPasswordStoreDir())
		}
		return m, nil
		
	case gitUpdateResult:
		// Handle git update result
		if msg.err != nil {
			m.error = fmt.Errorf("git update failed: %v", msg.err)
		} else {
			// Refresh git status after successful update
			m.gitStatus = git.GetGitStatus(getPasswordStoreDir())
		}
		return m, nil
	}
	
	// Handle list-specific messages
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	
	return m, tea.Batch(cmds...)
}

// View renders the current state
func (m *Model) View() string {
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
	view += getTitle(m.mode) + "\n"
	
	// Git status line (if available)
	if m.gitStatus.IsGitRepo {
		gitStatusLine := getGitStatusLine(m.gitStatus)
		if gitStatusLine != "" {
			view += gitStatusLine + "\n"
		}
	} else if m.gitStatusErr != nil {
		view += fmt.Sprintf("Git: %v\n", m.gitStatusErr)
	} else {
		view += "Git: not a repository\n"
	}
	view += ""
	
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
func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		
	// Tab - cycle through results
	case msg.String() == "tab":
		// Cycle to next item in the list
		if len(m.list.Items()) > 1 {
			nextIdx := (m.list.Index() + 1) % len(m.list.Items())
			m.list.Select(nextIdx)
		}
		return m, nil
		
	// Git operations
	case msg.String() == "ctrl+p":
		// Push changes to remote
		return m, handleGitPush
		
	case msg.String() == "ctrl+u":
		// Update (pull) from remote
		return m, handleGitUpdate
		
	case msg.String() == "ctrl+r":
		// Refresh git status
		m.gitStatus = git.GetGitStatus(getPasswordStoreDir())
		return m, nil
		
	// List navigation keys - pass to list
	case msg.String() == "up" || msg.String() == "down" || msg.String() == "pageup" || msg.String() == "pagedown":
		// Pass vertical navigation keys to the list
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
		
	// Input cursor keys - pass to input field
	case msg.String() == "left" || msg.String() == "right" || msg.String() == "home" || msg.String() == "end":
		// Pass cursor movement keys to the input field
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
		
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

// filterList filters the list based on the current input value using fuzzy matching
func (m *Model) filterList() {
	query := m.input.Value()
	
	// Use fuzzy matching
	var filtered []list.Item
	if query == "" {
		// Empty query: show all passwords
		filtered = make([]list.Item, len(m.allPasswords))
		for i, path := range m.allPasswords {
			filtered[i] = item{path: path}
		}
	} else {
		// Filter using fuzzy matching
		results := fuzzy.Filter(query, m.allPasswords)
		filtered = make([]list.Item, len(results))
		for i, result := range results {
			filtered[i] = item{
				path:         result.Path,
				matchIndices: result.MatchIndices,
			}
		}
	}
	
	m.list.SetItems(filtered)
}

// helpView returns the help text
func helpView() string {
	return "↑/↓: Navigate | Enter: Select | Esc/Ctrl+C: Cancel | Ctrl+Q: Quit | Ctrl+P: Push | Ctrl+U: Update | Ctrl+R: Refresh"
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

// Git operation commands

// gitPushResult represents the result of a git push operation
type gitPushResult struct {
	err error
}

// gitUpdateResult represents the result of a git update operation
type gitUpdateResult struct {
	err error
}

// handleGitPush attempts to push changes to the remote
func handleGitPush() tea.Msg {
	storeDir := getPasswordStoreDir()
	err := git.Push(storeDir)
	return gitPushResult{err: err}
}

// handleGitUpdate attempts to update (pull) from the remote
func handleGitUpdate() tea.Msg {
	storeDir := getPasswordStoreDir()
	err := git.Update(storeDir)
	return gitUpdateResult{err: err}
}

// NewModel creates a new fuzzy search model
func NewModel(passwords []string, mode FuzzySearchMode) *Model {
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
		gitStatus:    git.GetGitStatus(getPasswordStoreDir()),
		loading:      false,
		quitting:     false,
		width:        80,
		height:       24,
	}
	
	return model
}
