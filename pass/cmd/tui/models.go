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
	"github.com/mandu/tools/pass/cmd/tree"
	"github.com/mandu/tools/pass/pkg/fuzzy"
	"github.com/mandu/tools/pass/pkg/git"
)

// item represents a password entry in the list (flat view)
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



// treeFormattedItem extends item to support tree view display
type treeFormattedItem struct {
	item
	displayName string // Tree-formatted display name
}

// Title returns the tree-formatted display name
func (t treeFormattedItem) Title() string { return t.displayName }

// passwordDelegate is a custom delegate for rendering password items
type passwordDelegate struct {
	list.DefaultDelegate
}

// CreateTreeFormattedItems creates list items with tree formatting
func CreateTreeFormattedItems(passwords []string) []list.Item {
	if len(passwords) == 0 {
		return []list.Item{}
	}
	
	// Build tree from all passwords
	treeRoot := tree.BuildTreeFromPaths(passwords)
	
	// Flatten tree to items with tree formatting
	var items []list.Item
	if treeRoot != nil && treeRoot.Name == "" {
		// Root is empty container, always process its children
		for i := range treeRoot.Children {
			isLast := i == len(treeRoot.Children)-1
			flattenTreeToFormattedItems(treeRoot.Children[i], "", isLast, "", &items)
		}
	} else if treeRoot != nil {
		// Root has a name (shouldn't happen with BuildTreeFromPaths, but handle it)
		flattenTreeToFormattedItems(treeRoot, "", false, "", &items)
	}
	
	return items
}

// flattenTreeToFormattedItems recursively flattens the tree to formatted items
func flattenTreeToFormattedItems(node *tree.TreeNode, prefix string, isLast bool, parentPath string, items *[]list.Item) {
	// Build the full path for this node
	nodePath := parentPath
	if parentPath != "" {
		nodePath += "/"
	}
	nodePath += node.Name
	
	// Determine the connector for this node
	connector := "\u2514\u2500\u2500 " // └──
	if !isLast {
		connector = "\u251C\u2500\u2500 " // ├──
	}
	
	// Format the display name with tree structure
	displayName := prefix + connector + node.Name
	
	// Add / suffix for directories
	if node.IsDir {
		displayName += "/"
	}
	
	// Create item for all nodes (both directories and passwords)
	// Note: A node can be both a password (IsPassword=true) and a directory (IsDir=true)
	formattedItem := treeFormattedItem{
		item: item{
			path: nodePath,
		},
		displayName: displayName,
	}
	*items = append(*items, formattedItem)
	
	// Process children for directories or nodes with children
	if len(node.Children) > 0 {
		// Build child prefix
		childPrefix := prefix
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "\u2502   " // │   
		}
		
		// Process children
		for i, child := range node.Children {
			childIsLast := i == len(node.Children)-1
			flattenTreeToFormattedItems(child, childPrefix, childIsLast, nodePath, items)
		}
	}
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
	// Handle both item and treeFormattedItem types
	var matchIndices []int
	var displayName string
	
	switch v := listItem.(type) {
	case item:
		matchIndices = v.matchIndices
		displayName = v.path // For flat view, display the path
	case treeFormattedItem:
		matchIndices = v.matchIndices
		displayName = v.displayName // For tree view, display the formatted name
	default:
		return // Unknown type, skip rendering
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
	displayPath := displayName
	if len(matchIndices) > 0 {
		displayPath = highlightMatches(displayName, matchIndices)
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
	error          error
	quitting       bool
	selected       string
	width          int
	height         int
	gitStatusReady bool // Indicates if git status has been loaded
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
	// Start async git status check
	return tea.Batch(
		textinput.Blink,
		loadGitStatusCmd(getPasswordStoreDir()),
	)
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
			// Refresh git status and password list after successful push
			m.gitStatusReady = false
			cmds = append(cmds, loadGitStatusCmd(getPasswordStoreDir()))
			m.refreshPasswordList()
		}
		return m, tea.Batch(cmds...)
		
	case gitUpdateResult:
		// Handle git update result
		if msg.err != nil {
			m.error = fmt.Errorf("git update failed: %v", msg.err)
		} else {
			// Refresh git status and password list after successful update
			m.gitStatusReady = false
			cmds = append(cmds, loadGitStatusCmd(getPasswordStoreDir()))
			m.refreshPasswordList()
		}
		return m, tea.Batch(cmds...)
		
	case gitStatusResult:
		// Handle async git status result
		m.gitStatus = msg.status
		m.gitStatusErr = msg.err
		m.gitStatusReady = true
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
	
	if m.error != nil {
		return "Error: " + m.error.Error() + "\n"
	}
	
	// Build the view
	var view string
	
	// Header
	view += getTitle(m.mode) + "\n"
	
	// Git status line (if available)
	if !m.gitStatusReady {
		view += "Git: loading...\n"
	} else if m.gitStatusErr != nil {
		view += fmt.Sprintf("Git: %v\n", m.gitStatusErr)
	} else if m.gitStatus.IsGitRepo {
		gitStatusLine := getGitStatusLine(m.gitStatus)
		if gitStatusLine != "" {
			view += gitStatusLine + "\n"
		}
	} else {
		view += "Git: not a repository\n"
	}
	
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
		// Refresh git status asynchronously
		m.gitStatusReady = false
		return m, loadGitStatusCmd(getPasswordStoreDir())
		
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
	
	// Use fuzzy matching on all passwords
	var matchedPaths []string
	if query == "" {
		// Empty query: show all passwords
		matchedPaths = m.allPasswords
	} else {
		// Filter using fuzzy matching
		results := fuzzy.Filter(query, m.allPasswords)
		matchedPaths = make([]string, len(results))
		for i, result := range results {
			matchedPaths[i] = result.Path
		}
	}
	
	// Create filtered items with tree formatting
	filtered := CreateTreeFormattedItems(matchedPaths)
	
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

// gitStatusResult represents the result of async git status loading
type gitStatusResult struct {
	status git.GitStatus
	err    error
}

// loadGitStatusCmd creates a command to load git status asynchronously
func loadGitStatusCmd(dir string) tea.Cmd {
	return func() tea.Msg {
		status := git.GetGitStatus(dir)
		return gitStatusResult{status: status, err: nil}
	}
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
	
	// Create list items with tree formatting
	items := CreateTreeFormattedItems(passwords)
	
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
	
	// Create model with empty git status (will be loaded async)
	model := &Model{
		list:          listModel,
		input:         input,
		mode:          mode,
		allPasswords:  passwords,

		gitStatus:     git.GitStatus{IsGitRepo: false},
		gitStatusErr:  nil,
		quitting:      false,
		gitStatusReady: false,
		width:         80,
		height:        24,
	}
	
	return model
}

// refreshPasswordList refreshes the list of passwords from the store
func (m *Model) refreshPasswordList() {
	passwords, err := CollectAllPasswords(getPasswordStoreDir())
	if err != nil {
		// Keep existing list if there's an error
		return
	}
	m.allPasswords = passwords
	// Re-filter with current query
	m.filterList()
}