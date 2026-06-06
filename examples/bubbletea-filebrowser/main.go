// Package main implements a minimal file browser TUI using Bubble Tea.
// This demonstrates cross-platform terminal application development.
//
// Features:
// - Lists files in current directory
// - Navigate with j/k or arrow keys
// - Open selected file with Enter (uses $EDITOR or platform default)
// - Quit with Escape or q
//
// Run with: go run main.go
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
