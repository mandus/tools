package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mandu/tools/pass/pkg/filesystem"
	"github.com/mandu/tools/pass/pkg/fuzzy"
)

// fuzzySearchMode enters interactive fuzzy search mode.
// It lists all passwords and allows the user to search and select one.
// When a password is selected (Enter), it shows the password.
func fuzzySearchMode() error {
	storeDir := GetPasswordStoreDir()

	// Check if store exists
	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		return fmt.Errorf("pass: password store does not exist. Use 'pass insert' to create it.")
	}

	// Collect all password files
	passwords, err := collectAllPasswords(storeDir)
	if err != nil {
		return err
	}

	if len(passwords) == 0 {
		return fmt.Errorf("pass: no passwords found in store")
	}

	// For now, implement a simple non-interactive fuzzy search
	// The user can type a query and we'll show matching passwords
	// This is a simplified version that doesn't require complex terminal handling
	
	fmt.Println("Fuzzy Search Mode (type to search, Enter to select, Ctrl+C to exit)")
	fmt.Println("Passwords:")
	for i, p := range passwords {
		fmt.Printf("  %d. %s\n", i+1, p)
	}
	
	// Read query from user
	fmt.Print("\nSearch: ")
	reader := bufio.NewReader(os.Stdin)
	query, err := reader.ReadString('\n')
	if err != nil {
		return nil // User pressed Ctrl+C or similar
	}
	
	// Trim newline
	query = strings.TrimSuffix(query, "\n")
	query = strings.TrimSuffix(query, "\r")

	// If empty query, exit
	if query == "" {
		return nil
	}

	// Filter passwords using fuzzy matching
	matches := fuzzy.Filter(query, passwords)

	if len(matches) == 0 {
		fmt.Println("No matches found.")
		return nil
	}

	// Show matches
	fmt.Println("\nMatching passwords:")
	for i, match := range matches {
		// Highlight matching characters (simple version without colors)
		fmt.Printf("  %d. %s\n", i+1, match.Path)
	}

	// If exactly one match, select it automatically
	if len(matches) == 1 {
		return showSelectedPassword(matches[0].Path)
	}

	// Let user select from matches
	fmt.Print("\nSelect (number or Enter for first): ")
	selection, err := reader.ReadString('\n')
	if err != nil {
		return nil
	}

	selection = strings.TrimSuffix(selection, "\n")
	selection = strings.TrimSuffix(selection, "\r")

	// Parse selection
	var selectedIdx int
	if selection == "" {
		selectedIdx = 0 // Default to first
	} else {
		_, err := fmt.Sscanf(selection, "%d", &selectedIdx)
		if err != nil || selectedIdx < 1 || selectedIdx > len(matches) {
			fmt.Println("Invalid selection.")
			return nil
		}
		selectedIdx-- // Convert to 0-indexed
	}

	return showSelectedPassword(matches[selectedIdx].Path)
}

// collectAllPasswords walks the password store and collects all password paths.
func collectAllPasswords(storeDir string) ([]string, error) {
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
			// Get relative path from store directory
			relPath, err := filepath.Rel(storeDir, path)
			if err != nil {
				return err
			}
			// Normalize path separators for display
			relPath = filesystem.NormalizePathForDisplay(relPath)
			// Strip .gpg extension
			passwordPath := strings.TrimSuffix(relPath, ".gpg")
			passwords = append(passwords, passwordPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("pass: failed to list passwords: %v", err)
	}

	return passwords, nil
}

// showSelectedPassword shows the password at the selected path.
func showSelectedPassword(path string) error {
	// Reuse the showPassword function from show.go
	return showPassword(path)
}

// InteractiveFuzzySearch provides a more advanced fuzzy search with real-time filtering.
// This is the full implementation that will replace fuzzySearchMode once terminal
// handling is more robust.
func InteractiveFuzzySearch() error {
	storeDir := GetPasswordStoreDir()

	// Check if store exists
	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		return fmt.Errorf("pass: password store does not exist. Use 'pass insert' to create it.")
	}

	// Collect all password files
	passwords, err := collectAllPasswords(storeDir)
	if err != nil {
		return err
	}

	if len(passwords) == 0 {
		return fmt.Errorf("pass: no passwords found in store")
	}

	// For now, fall back to simple search
	// TODO: Implement full interactive mode with terminal package
	return fuzzySearchMode()
}

// RmFuzzySearch is called by rm command when no path is provided.
// It allows user to select a password to remove using fuzzy search.
func RmFuzzySearch(noCommit, clip bool) error {
	storeDir := GetPasswordStoreDir()

	// Collect all password files
	passwords, err := collectAllPasswords(storeDir)
	if err != nil {
		return err
	}

	if len(passwords) == 0 {
		return fmt.Errorf("pass: no passwords found in store")
	}

	// Simple non-interactive selection
	fmt.Println("Select password to remove:")
	for i, p := range passwords {
		fmt.Printf("  %d. %s\n", i+1, p)
	}
	
	fmt.Print("\nEnter number to remove (or Ctrl+C to cancel): ")
	reader := bufio.NewReader(os.Stdin)
	selection, err := reader.ReadString('\n')
	if err != nil {
		return nil // User cancelled
	}

	selection = strings.TrimSuffix(selection, "\n")
	selection = strings.TrimSuffix(selection, "\r")

	// Parse selection
	var selectedIdx int
	if selection == "" {
		return nil // No selection
	}
	
	_, err = fmt.Sscanf(selection, "%d", &selectedIdx)
	if err != nil || selectedIdx < 1 || selectedIdx > len(passwords) {
		fmt.Println("Invalid selection.")
		return nil
	}
	selectedIdx-- // Convert to 0-indexed

	// Remove the selected password
	return removePassword(passwords[selectedIdx], noCommit, clip)
}

// UseFuzzyPath is a helper that uses fuzzy matching to find the best match
// for a partial path and returns the full path.
func UseFuzzyPath(query string) (string, error) {
	storeDir := GetPasswordStoreDir()

	// Collect all password files
	passwords, err := collectAllPasswords(storeDir)
	if err != nil {
		return "", err
	}

	if len(passwords) == 0 {
		return "", fmt.Errorf("pass: no passwords found in store")
	}

	// Find best match
	bestMatch := fuzzy.FindBestMatch(query, passwords)
	if bestMatch == "" {
		return "", fmt.Errorf("pass: no match found for %q", query)
	}

	return bestMatch, nil
}
