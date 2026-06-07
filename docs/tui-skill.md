---
name: tui-skill
description: Create cross-platform terminal applications using Bubble Tea. Use when you need to build interactive TUIs with keyboard navigation, file listing, and editor integration that work consistently on Windows, Linux, and macOS.
license: MIT
compatibility: Go 1.20+, Windows/Linux/macOS
metadata:
  framework: charmbracelet/bubbletea
  language: Go
  tags: [tui, terminal, cross-platform, file-browser]
---

# TUI Development with Bubble Tea

This skill provides patterns and examples for building cross-platform terminal user interfaces using [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea). Bubble Tea is a Go framework for building terminal applications that works consistently across Windows, Linux, and macOS - solving the terminal portability issues common with other approaches.

## When to Use

Use this skill when:
- You need a cross-platform TUI that works on Windows, Linux, and macOS
- You want consistent keyboard handling across different terminals
- You need to build interactive file browsers, forms, or dashboards
- You're experiencing terminal compatibility issues with your current approach (like in the `pass` program)

## Quick Start

### 1. Install Bubble Tea

```bash
# Install the Bubble Tea library
go get github.com/charmbracelet/bubbletea

# Install optional dependencies for enhanced features
go get github.com/charmbracelet/bubbles  # Pre-built components (list, textinput, etc.)
go get github.com/charmbracelet/lipgloss  # Style library
```

### 2. Run the Test Program

A working file browser example is included in the repository:

```bash
cd examples/bubbletea-filebrowser
go run main.go
```

Or build and run:
```bash
go build -o filebrowser main.go
./filebrowser
```

On Windows:
```powershell
cd examples\bubbletea-filebrowser
go run main.go
```

### 3. Minimal File Browser Example

Here's a complete, working example of a simplified `ls` with keyboard navigation:

```go
package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
)

// item represents a file in the list
type item struct {
	name string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.name }

// model holds the application state
type model struct {
	list     list.Model
	quitting bool
}

// Init is called when the program starts
func (m model) Init() tea.Cmd {
	return nil
}

// initialModel creates the initial application state
func initialModel() model {
	// Read files from current directory
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
		os.Exit(1)
	}

	// Create list items from files (excluding directories)
	var items []list.Item
	for _, f := range files {
		if !f.IsDir() {
			items = append(items, item{name: f.Name()})
		}
	}

	// Create the list with default delegate
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Files (j/↓ k/↑ Enter Esc)"
	l.SetShowHelp(false)

	return model{list: l}
}

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard input
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			// Quit the application
			m.quitting = true
			return m, tea.Quit

		case "enter":
			// Open selected file in editor
			if m.list.SelectedItem() != nil {
				selected := m.list.SelectedItem().(item)
				editor := getEditor()
				cmd := exec.Command(editor, selected.name)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
					return nil
				})
			}

		case "j", "down":
			// Move cursor down
			newListModel, cmd := m.list.Update(msg)
			m.list = newListModel
			return m, cmd

		case "k", "up":
			// Move cursor up
			newListModel, cmd := m.list.Update(msg)
			m.list = newListModel
			return m, cmd
		}

	case tea.WindowSizeMsg:
		// Handle terminal resize
		m.list.SetWidth(msg.Width)
		return m, nil
	}

	// Handle list-specific messages
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the current state
func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}
	return m.list.View()
}

// getEditor returns the appropriate editor for the current platform
func getEditor() string {
	// Try EDITOR environment variable first
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	// Try VISUAL
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	// Platform defaults
	if os.PathSeparator == '\\' {
		// Windows
		return "notepad"
	}
	// Unix-like systems
	return "nano"
}

func main() {
	// Initialize and run the Bubble Tea program
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
```

## How to Use This Example

1. Save the code above to a file named `main.go`
2. Initialize a Go module: `go mod init filebrowser`
3. Add dependencies: `go get github.com/charmbracelet/bubbletea github.com/charmbracelet/bubbles`
4. Run it: `go run main.go`

**Controls:**
- `j` or `↓` (Down Arrow) - Move down
- `k` or `↑` (Up Arrow) - Move up
- `Enter` - Open file in editor
- `Esc` or `q` - Quit

## Key Concepts

### The Tea Program

Bubble Tea applications revolve around the `tea.Program`:

```go
p := tea.NewProgram(initialModel())
if _, err := p.Run(); err != nil {
    // handle error
}
```

### The Model

Your model holds application state and must implement the `tea.Model` interface:

```go
type model struct {
    list     list.Model
    selected string
}

// Required: Init method
func (m model) Init() tea.Cmd {
    return nil
}
```

### Update Function

Handles messages and returns updated model + command:

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "j": // handle j key
        case "k": // handle k key
        case "enter": // handle enter
        case "esc", "q": // handle quit
        }
    case tea.WindowSizeMsg:
        // handle window resize
    }
    return m, nil
}
```

### View Function

Renders the current state:

```go
func (m model) View() string {
    return "Your view here"
}
```

## Keyboard Handling

Bubble Tea normalizes keyboard input across platforms:

| Key | Windows | Linux/macOS | msg.String() |
|-----|---------|-------------|--------------|
| Down Arrow | ↓ | ↓ | "down" |
| Up Arrow | ↑ | ↑ | "up" |
| Enter | Enter | Enter | "enter" |
| Escape | Esc | Esc | "esc" |
| j | j | j | "j" |
| k | k | k | "k" |
| Ctrl+C | Ctrl+C | Ctrl+C | "ctrl+c" |

**This means your key handling code works the same on all platforms!** This is the key feature that solves the Windows/Linux terminal compatibility issues you're experiencing with the `pass` program.

## Cross-Platform Considerations

### Editor Detection

```go
func getEditor() string {
    // Try EDITOR first
    if editor := os.Getenv("EDITOR"); editor != "" {
        return editor
    }
    // Try VISUAL
    if editor := os.Getenv("VISUAL"); editor != "" {
        return editor
    }
    // Platform defaults
    if runtime.GOOS == "windows" {
        return "notepad"
    }
    return "nano"
}
```

### File Paths

Use `os.PathSeparator` and `filepath` package for cross-platform paths:

```go
import "path/filepath"

// Join paths
path := filepath.Join("dir", "subdir", "file.txt")

// Get absolute path
absPath, _ := filepath.Abs("relative/path")
```

## Enhanced Version with Styling

For a more polished look, use the `lipgloss` library:

```go
package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.Padding(1, 0)
)

type item struct {
	name string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.name }

type model struct {
	list     list.Model
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
		os.Exit(1)
	}

	var items []list.Item
	for _, f := range files {
		if !f.IsDir() {
			items = append(items, item{name: f.Name()})
		}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "File Browser"
	l.Styles.Title = titleStyle
	l.Styles.Pagination = paginationStyle
	l.Styles.Help = helpStyle

	return model{list: l}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if m.list.SelectedItem() != nil {
				selected := m.list.SelectedItem().(item)
				editor := getEditor()
				cmd := exec.Command(editor, selected.name)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
					return nil
				})
			}

		case "j", "down":
			newListModel, cmd := m.list.Update(msg)
			m.list = newListModel
			return m, cmd

		case "k", "up":
			newListModel, cmd := m.list.Update(msg)
			m.list = newListModel
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func getEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	if os.PathSeparator == '\\' {
		return "notepad"
	}
	return "nano"
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}
	return m.list.View()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
```

## Custom Delegate for Better File Display

For more control over how files are displayed, create a custom delegate:

```go
// Custom delegate for file items
type fileDelegate struct{}

func (d fileDelegate) Height() int                             { return 1 }
func (d fileDelegate) Spacing() int                            { return 0 }
func (d fileDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d fileDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedStyle.Render("> " + s)
		}
	}

	// Add file extension color
	if strings.HasSuffix(i.name, ".go") {
		fn = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ADD8")).Render
	} else if strings.HasSuffix(i.name, ".md") {
		fn = lipgloss.NewStyle().Foreground(lipgloss.Color("#50E3C2")).Render
	}

	fmt.Fprint(w, fn(str))
}
```

## Testing Your TUI

### Local Testing

```bash
# Run directly
go run main.go

# Build and run
go build -o filebrowser main.go
./filebrowser
```

### Windows-Specific Testing

```powershell
# In PowerShell
go run main.go

# Build for Windows
go build -o filebrowser.exe main.go
```

## Common Patterns

### Adding Help Text

```go
func (m model) View() string {
    help := "j/↓: Down | k/↑: Up | Enter: Open | Esc: Quit"
    return m.list.View() + "\n" + helpStyle.Render(help)
}
```

### Filtering Files

```go
func initialModel() model {
    files, _ := os.ReadDir(".")
    var items []list.Item
    
    // Filter by extension
    for _, f := range files {
        if !f.IsDir() && strings.HasSuffix(f.Name(), ".go") {
            items = append(items, item{name: f.Name()})
        }
    }
    
    return model{list: list.New(items, delegate, 0, 0)}
}
```

### Async File Loading

```go
func loadFiles() tea.Cmd {
    return func() tea.Msg {
        files, err := os.ReadDir(".")
        if err != nil {
            return errorMsg{err}
        }
        var items []list.Item
        for _, f := range files {
            items = append(items, item{name: f.Name()})
        }
        return filesLoadedMsg{items: items}
    }
}

type filesLoadedMsg struct {
    items []list.Item
}

type errorMsg struct {
    err error
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case filesLoadedMsg:
        m.list.SetItems(msg.items)
        return m, nil
    case errorMsg:
        m.error = msg.err
        return m, nil
    }
    return m, nil
}

func initialModel() model {
    return model{
        list: list.New([]list.Item{}, delegate, 0, 0),
        loading: true,
    }
}

func main() {
    p := tea.NewProgram(initialModel())
    p.Send(loadFiles())
    p.Run()
}
```

## Project Structure for Bubble Tea Apps

```
my-tui-app/
├── main.go           # Entry point
├── go.mod            # Go module
├── go.sum            # Dependencies
├── models/           # Application models
│   └── filebrowser.go
├── components/       # Reusable components
│   └── list.go
├── styles/           # Styling
│   └── theme.go
└── utils/            # Utilities
    └── filesystem.go
```

## Troubleshooting

### Windows Terminal Issues

1. **Keys not working**: Ensure you're using a modern terminal (Windows Terminal, not cmd.exe)
2. **Colors not showing**: Check if terminal supports ANSI colors
3. **Slow rendering**: Reduce the number of items or use pagination

### Common Errors

**"not a terminal" error**:
```bash
# Force terminal mode (Windows)
go run -exec "winpty -X allow-non-tty" main.go
```

**Dependency issues**:
```bash
go mod tidy
go get -u github.com/charmbracelet/bubbletea
```

## Tree View Rendering

For displaying hierarchical data like file systems or password stores, use box-drawing characters to create a tree structure.

### Box-Drawing Characters

| Character | Unicode | Description | Usage |
|-----------|---------|-------------|-------|
| `├` | U+251C | Box drawings light vertical and horizontal | Branch with siblings below |
| `└` | U+2514 | Box drawings light up and horizontal | Last branch in group |
| `│` | U+2502 | Box drawings light vertical | Vertical connector |
| `─` | U+2500 | Box drawings light horizontal | Horizontal line |

**Combined characters:**
- `├── ` : Item with siblings below
- `└── ` : Last item in group
- `│   ` : Vertical connector with indentation
- `    ` : Pure indentation (4 spaces)

### Tree Node Structure

```go
package tree

import (
	"fmt"
	"strings"
)

// TreeNode represents a node in a hierarchical tree
type TreeNode struct {
	Name     string
	IsDir    bool
	Children []*TreeNode
}

// NewTreeNode creates a new tree node
func NewTreeNode(name string, isDir bool) *TreeNode {
	return &TreeNode{
		Name:     name,
		IsDir:    isDir,
		Children: []*TreeNode{},
	}
}

// AddChild adds a child and maintains alphabetical sort
func (n *TreeNode) AddChild(child *TreeNode) {
	n.Children = append(n.Children, child)
	// Sort children alphabetically (bubble sort for simplicity)
	for i := 0; i < len(n.Children)-1; i++ {
		for j := 0; j < len(n.Children)-i-1; j++ {
			if n.Children[j].Name > n.Children[j+1].Name {
				n.Children[j], n.Children[j+1] = n.Children[j+1], n.Children[j]
			}
		}
	}
}

// FindOrCreateChild finds a child by name or creates it
func (n *TreeNode) FindOrCreateChild(name string, isDir bool) *TreeNode {
	for _, child := range n.Children {
		if child.Name == name {
			return child
		}
	}
	child := NewTreeNode(name, isDir)
	n.AddChild(child)
	return child
}
```

### Tree Rendering with Box-Drawing Characters

```go
// Render renders the tree with proper indentation and connectors
func (n *TreeNode) Render(prefix string) string {
	var sb strings.Builder

	// Determine connector based on position
	connector := "└── "
	if len(prefix) > 0 && !strings.HasSuffix(prefix, "    ") {
		connector = "├── "
	}

	// Format name
	displayName := n.Name
	if n.IsDir {
		displayName += "/"
	}

	sb.WriteString(prefix + connector + displayName + "\n")

	// Render children
	for i, child := range n.Children {
		isLast := i == len(n.Children)-1
		childPrefix := prefix
		if !isLast {
			childPrefix += "│   "
		} else {
			childPrefix += "    "
		}
		sb.WriteString(child.Render(childPrefix))
	}

	return sb.String()
}

// BuildTreeFromPaths creates a tree from flat path list
func BuildTreeFromPaths(paths []string) *TreeNode {
	root := NewTreeNode("", false)

	for _, path := range paths {
		// Split path into components
		components := strings.Split(path, "/")

		// Build tree structure
		current := root
		for i, component := range components {
			isDir := i < len(components)-1
			current = current.FindOrCreateChild(component, isDir)
		}
	}

	return root
}
```

### Complete Tree View Example

```go
package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	// Build a sample tree
	root := tree.NewTreeNode("", false)
	
	// Add structure: dev/hafslund/mistral-vibe-key
	dev := root.FindOrCreateChild("dev", true)
	hafslund := dev.FindOrCreateChild("hafslund", true)
	hafslund.AddChild(tree.NewTreeNode("mistral-vibe-key", false))
	
	// Add structure: dev/mistral.ai/api-key
	mistralAI := dev.FindOrCreateChild("mistral.ai", true)
	mistralAI.AddChild(tree.NewTreeNode("api-key", false))
	
	// Render with optional styling
	output := ""
	for i, child := range root.Children {
		output += child.Render("")
	}
	
	fmt.Print(output)
	// Output:
	// ├── dev/
	// │   ├── hafslund/
	// │   │   └── mistral-vibe-key
	// │   └── mistral.ai/
	// │       └── api-key
}
```

### Tree View with Lip Gloss Styling

For enhanced visual appearance, apply styles to different tree elements:

```go
var (
	// Style for directory names
	dirStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5DADE2")) // Blue

	// Style for file names
	fileStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")) // White

	// Style for tree characters (├──, └──, │)
	treeCharStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D3C98")) // Purple
)

// StyledRender renders tree with colors
func (n *TreeNode) StyledRender(prefix string) string {
	var sb strings.Builder

	connector := "└── "
	if len(prefix) > 0 && !strings.HasSuffix(prefix, "    ") {
		connector = "├── "
	}

	displayName := n.Name
	if n.IsDir {
		displayName += "/"
	}

	// Style the connector and name differently
	styledConnector := treeCharStyle.Render(connector)
	var styledName string
	if n.IsDir {
		styledName = dirStyle.Render(displayName)
	} else {
		styledName = fileStyle.Render(displayName)
	}

	sb.WriteString(prefix + styledConnector + styledName + "\n")

	for i, child := range n.Children {
		isLast := i == len(n.Children)-1
		childPrefix := prefix
		if !isLast {
			childPrefix += treeCharStyle.Render("│   ")
		} else {
			childPrefix += "    "
		}
		sb.WriteString(child.StyledRender(childPrefix))
	}

	return sb.String()
}
```

### Handling Deep Nesting

For very deep directory structures, consider:

1. **Limiting depth**: Stop rendering after a certain level
2. **Collapsing**: Show "..." for deep branches
3. **Horizontal scrolling**: For terminals that support it
4. **Truncation**: Use ellipsis for long names

```go
const maxDepth = 10

func (n *TreeNode) RenderWithDepth(prefix string, depth int) string {
	if depth > maxDepth {
		return prefix + "└── ... (truncated)\n"
	}
	// ... normal rendering ...
}
```

## Resources

- [Bubble Tea GitHub](https://github.com/charmbracelet/bubbletea)
- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/examples/tutorials)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Lip Gloss Styling](https://github.com/charmbracelet/lipgloss)
- [Charmbracelet Awesome](https://github.com/charmbracelet/awesome)

## Why This Solves Your Problem

The `pass` program likely has terminal portability issues between Windows and Linux because:

1. **Different key codes**: Terminals on different OSes send different escape sequences for the same keys
2. **Terminal capability differences**: Not all terminals support the same features
3. **Manual terminal handling**: If you're using low-level terminal libraries, they may not abstract these differences

Bubble Tea solves all of these:
- Normalizes keyboard input across platforms (j/k/arrow keys work the same everywhere)
- Automatically detects terminal capabilities
- Provides a consistent API that works on Windows Terminal, Linux terminals, and macOS Terminal
- Handles window resizing and other terminal events gracefully

## Next Steps for Your Pass Program

To integrate Bubble Tea into your `pass` program:

1. Create a new directory for TUI components: `mkdir -p pass/pkg/tui`
2. Add Bubble Tea to your go.mod: `go get github.com/charmbracelet/bubbletea`
3. Start with a simple TUI for one feature (e.g., password selection)
4. Gradually migrate existing CLI functionality to TUI
5. Test on both Windows and Linux early and often

The key advantage is that the same code will work on both platforms without modification.
