package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mandu/tools/pass/pkg/filesystem"
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find [string]",
	Short: "Search for passwords",
	Long:  `Search for passwords containing the given string anywhere in their path.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return findPasswords(args[0])
	},
}

// Flags for find command
var ignoreCaseFlag bool

func addFindCmd() {
	findCmd.Flags().BoolVarP(&ignoreCaseFlag, "ignore-case", "i", false, "Case-insensitive search")
	rootCmd.AddCommand(findCmd)
}

// findPasswords searches for passwords containing the search string
func findPasswords(searchString string) error {
	if searchString == "" {
		return fmt.Errorf("pass: search string cannot be empty")
	}
	
	storeDir := GetPasswordStoreDir()
	
	// Check if store exists
	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		fmt.Println() // Print nothing if store doesn't exist
		return nil
	}
	
	// Prepare search string
	target := searchString
	if ignoreCaseFlag {
		target = strings.ToLower(searchString)
	}
	
	// Walk the directory tree
	var results []string
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
			
			// Normalize path separators
			relPath = filesystem.NormalizePathForDisplay(relPath)
			
			// Strip .gpg extension
			passwordPath := strings.TrimSuffix(relPath, ".gpg")
			
			// Check if path contains search string
			checkPath := passwordPath
			if ignoreCaseFlag {
				checkPath = strings.ToLower(passwordPath)
			}
			
			if strings.Contains(checkPath, target) {
				results = append(results, passwordPath)
			}
		}
		
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("pass: failed to walk directory: %v", err)
	}
	
	// Sort results
	for i := 0; i < len(results)-1; i++ {
		for j := 0; j < len(results)-i-1; j++ {
			if results[j] > results[j+1] {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}
	
	// Print results
	for _, result := range results {
		fmt.Println(result)
	}
	
	return nil
}
