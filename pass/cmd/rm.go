package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mandu/tools/pass/pkg/filesystem"
	"github.com/mandu/tools/pass/pkg/git"
	"github.com/mandu/tools/pass/pkg/gpg"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm [OPTIONS] [<path>]",
	Short: "Remove a password",
	Long: `Remove a password from the store.

If a path is provided, removes that specific password.
If no path is provided, will list available passwords for selection.

The password file is deleted and the removal is committed to git
(by default). The password remains in git history for recovery if needed.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) > 0 {
			path = args[0]
		}
		
		// Get flags
		noCommit, _ := cmd.Flags().GetBool("no-commit")
		clip, _ := cmd.Flags().GetBool("clip")
		
		if path != "" {
			return removePassword(path, noCommit, clip)
		}
		
		// No path provided - list passwords and let user select
		// For now, just show error. Fuzzy search will be added later.
		return listAndRemove(noCommit, clip)
	},
}

// Flags for rm command
var (
	noCommitFlagRm bool
	forceFlagRm    bool
)

func addRmCmd() {
	rmCmd.Flags().BoolVarP(&noCommitFlagRm, "no-commit", "n", false, "Skip git commit after removal")
	rmCmd.Flags().BoolVarP(&forceFlagRm, "force", "f", false, "Alias for --no-commit")
	rmCmd.Flags().BoolVarP(&clipFlagRm, "clip", "c", false, "Copy password to clipboard before deleting")
	rootCmd.AddCommand(rmCmd)
}

var clipFlagRm bool

// removePassword removes a password at the specified path.
func removePassword(path string, noCommit, clip bool) error {
	storeDir := GetPasswordStoreDir()

	// Normalize path and add .gpg extension if needed
	filePath := filesystem.NormalizePath(path)
	if !strings.HasSuffix(filePath, ".gpg") {
		filePath += ".gpg"
	}

	// Construct full path
	fullPath := filepath.Join(storeDir, filepath.FromSlash(filePath))

	// Strategy 2: Try with forward slashes directly
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		storeDirForward := strings.ReplaceAll(storeDir, "\\", "/")
		fullPath = storeDirForward + "/" + filePath
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("pass: %s: No such file or directory", path)
		}
	}

	// If clip flag is set, decrypt and copy to clipboard first
	if clip {
		password, err := gpg.DecryptFile(fullPath)
		if err != nil {
			return err
		}
		if err := filesystem.CopyToClipboard(password); err != nil {
			return fmt.Errorf("pass: failed to copy to clipboard: %v", err)
		}
		fmt.Printf("Copied %s to clipboard.\n", path)
		// Start auto-clear timer
		go filesystem.StartClipboardClearTimer()
	}

	// Remove the file
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("pass: failed to remove %s: %v", path, err)
	}

	// Git integration - remove and commit
	if !noCommit {
		// Check if this is a git repo
		gitDir := filepath.Join(storeDir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			// Run git rm and commit
			// We need to run git from the store directory
			if err := git.RemoveAndCommit(fullPath, "Remove "+path); err != nil {
				fmt.Fprintf(os.Stderr, "pass: warning: git operations failed: %v\n", err)
			}
		}
	}

	fmt.Printf("Password removed successfully.\n")
	return nil
}

// listAndRemove lists all passwords and lets user select one to remove.
func listAndRemove(noCommit, clip bool) error {
	storeDir := GetPasswordStoreDir()

	// Check if store exists
	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		return fmt.Errorf("pass: password store does not exist. Use 'pass insert' to create it.")
	}

	// Collect all password files
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
		return fmt.Errorf("pass: failed to list passwords: %v", err)
	}

	if len(passwords) == 0 {
		return fmt.Errorf("pass: no passwords found in store")
	}

	// For now, just print the list and return
	// TODO: Implement fuzzy search selection
	fmt.Println("Passwords:")
	for i, p := range passwords {
		fmt.Printf("  %d. %s\n", i+1, p)
	}
	fmt.Println("\nNote: Fuzzy search selection coming soon. Specify path directly for now.")

	return nil
}
