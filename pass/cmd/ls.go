package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mandu/tools/pass/pkg/filesystem"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls [path]",
	Short: "List passwords",
	Long:  `List all passwords in the password store, optionally filtered by a subpath.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var subpath string
		if len(args) > 0 {
			subpath = args[0]
		}
		return listPasswords(subpath)
	},
}

// Flags for ls command
var (
	recursiveFlag bool
	dirsOnlyFlag  bool
	filesOnlyFlag bool
)

func addLsCmd() {
	lsCmd.Flags().BoolVarP(&recursiveFlag, "recursive", "r", true, "Show full paths")
	lsCmd.Flags().BoolVarP(&dirsOnlyFlag, "dirs-only", "d", false, "List only directories")
	lsCmd.Flags().BoolVarP(&filesOnlyFlag, "files-only", "f", false, "List only files")
	rootCmd.AddCommand(lsCmd)
}

// listPasswords lists all passwords in the store
func listPasswords(subpath string) error {
	storeDir := GetPasswordStoreDir()
	
	// Check if store exists
	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		// Create the store directory
		if err := os.MkdirAll(storeDir, 0700); err != nil {
			return fmt.Errorf("pass: failed to create password store: %v", err)
		}
		// Initialize git repo if not exists
		if _, err := os.Stat(filepath.Join(storeDir, ".git")); os.IsNotExist(err) {
			// Silently initialize git repo
			_ = filesystem.RunCommand("git", "init", storeDir)
		}
		return nil
	}
	
	// Build the base path
	basePath := storeDir
	if subpath != "" {
		basePath = filepath.Join(storeDir, filesystem.NormalizePath(subpath))
		// Check if subpath exists
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			return fmt.Errorf("pass: %s: No such file or directory", subpath)
		}
	}
	
	// Walk the directory tree
	var results []string
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		
		// Get relative path from store directory
		relPath, err := filepath.Rel(storeDir, path)
		if err != nil {
			return err
		}
		
		// Normalize path separators for display
		relPath = filesystem.NormalizePathForDisplay(relPath)
		
		// Skip the base path itself
		if relPath == "." {
			return nil
		}
		
		// Handle directories
		if info.IsDir() {
			// Skip all directories unless --dirs-only flag is set
			if dirsOnlyFlag {
				// For recursive listing, we want full paths
				if recursiveFlag {
					results = append(results, relPath)
				} else if subpath == "" {
					// At root level, show directory names only
					results = append(results, info.Name())
				}
			}
			return nil
		}
		
		// Handle files - only include .gpg files (unless --dirs-only is set)
		if !dirsOnlyFlag {
			// Only include .gpg files
			if strings.HasSuffix(info.Name(), ".gpg") {
				// Strip .gpg extension
				passwordPath := strings.TrimSuffix(relPath, ".gpg")
				results = append(results, passwordPath)
			}
		}
		
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("pass: failed to walk directory: %v", err)
	}
	
	// Sort results
	// Simple bubble sort for now (small datasets expected)
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
