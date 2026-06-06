package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mandu/tools/pass/cmd"
	"github.com/mandu/tools/pass/pkg/filesystem"
)

// styles for the TUI
var (
	// General styles
	appStyle = lipgloss.NewStyle().Padding(1, 2)
	
	// Header styles
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(100)
	
	// Input styles
	inputStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00"))
	
	// List styles
	listStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#888888"))
	
	// Selected item style
	selectedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF00"))
	
	// Help styles
	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)
	
	// Error styles
	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF0000"))
	
	// Match highlight style
	matchStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFF00"))
)

// RunFuzzySearch runs the fuzzy search TUI and returns the selected password
func RunFuzzySearch(passwords []string, mode cmd.FuzzySearchMode) (string, error) {
	// Create the model
	model := NewModel(passwords, mode)
	
	// Create the tea program
	p := tea.NewProgram(model, tea.WithAltScreen())
	
	// Run the program
	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("tui: %v", err)
	}
	
	// Get the final model
	if finalModel == nil {
		return "", nil
	}
	
	final, ok := finalModel.(*Model)
	if !ok {
		return "", fmt.Errorf("tui: unexpected model type")
	}
	
	if final.quitting && final.selected == "" {
		// User cancelled
		return "", nil
	}
	
	return final.selected, nil
}

// CollectAllPasswords walks the password store and collects all password paths
func CollectAllPasswords(storeDir string) ([]string, error) {
	var passwords []string
	
	err := filepath.Walk(storeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		
		// Only process .gpg files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".gpg") {
			relPath, err := filepath.Rel(storeDir, path)
			if err != nil {
				return err
			}
			relPath = filesystem.NormalizePathForDisplay(relPath)
			passwordPath := strings.TrimSuffix(relPath, ".gpg")
			passwords = append(passwords, passwordPath)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to list passwords: %v", err)
	}
	
	return passwords, nil
}

// GetPasswordStoreDir returns the password store directory path
func GetPasswordStoreDir() string {
	return cmd.GetPasswordStoreDir()
}

// RunInteractiveFuzzySearch runs interactive search and returns the selected password path
// This is the main entry point for the TUI from the existing code
func RunInteractiveFuzzySearch(mode cmd.FuzzySearchMode) (string, error) {
	storeDir := GetPasswordStoreDir()
	
	// Check if store exists
	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		return "", fmt.Errorf("pass: password store does not exist. Use 'pass insert' to create it.")
	}
	
	// Collect all passwords
	passwords, err := CollectAllPasswords(storeDir)
	if err != nil {
		return "", err
	}
	
	if len(passwords) == 0 {
		return "", fmt.Errorf("pass: no passwords found in store")
	}
	
	// Run the TUI
	selected, err := RunFuzzySearch(passwords, mode)
	if err != nil {
		return "", err
	}
	
	return selected, nil
}
