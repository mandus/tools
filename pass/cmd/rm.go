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
If no path is provided, will enter interactive fuzzy search mode to select a password.

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
		force, _ := cmd.Flags().GetBool("force")
		clip, _ := cmd.Flags().GetBool("clip")
		
		// --force is an alias for --no-commit
		if force {
			noCommit = true
		}
		
		if path != "" {
			return removePassword(path, noCommit, clip)
		}
		
		// No path provided - enter fuzzy search mode
		return runRmFuzzySearch(noCommit, clip)
	},
}

// Flags for rm command
var (
	noCommitFlagRm bool
	forceFlagRm    bool
	clipFlagRm     bool
)

func addRmCmd() {
	rmCmd.Flags().BoolVarP(&noCommitFlagRm, "no-commit", "n", false, "Skip git commit after removal")
	rmCmd.Flags().BoolVarP(&forceFlagRm, "force", "f", false, "Alias for --no-commit")
	rmCmd.Flags().BoolVarP(&clipFlagRm, "clip", "c", false, "Copy password to clipboard before deleting")
	rootCmd.AddCommand(rmCmd)
}

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

// runRmFuzzySearch enters interactive fuzzy search mode for removing a password.
// When user selects a password, it will be removed.
func runRmFuzzySearch(noCommit, clip bool) error {
	// If clip flag is set, we need special handling
	// Since fuzzy search for rm with clip needs to copy before deleting
	if clip {
		// For clip mode in rm, we'll handle it in the RunInteractiveFuzzySearch
		// by passing a special mode or handling it separately
		// For now, we'll use a simple approach: do fuzzy search, then copy and delete
		selected, err := InteractiveFuzzySearch(FuzzyModeShow)
		if err != nil {
			return err
		}
		if selected == "" {
			return nil
		}
		// Get the full path
		fullPath := getRmFullPath(selected)
		
		// Copy to clipboard first
		password, err := gpg.DecryptFile(fullPath)
		if err != nil {
			return err
		}
		if err := filesystem.CopyToClipboard(password); err != nil {
			return fmt.Errorf("pass: failed to copy to clipboard: %v", err)
		}
		fmt.Printf("Copied %s to clipboard.\n", selected)
		go filesystem.StartClipboardClearTimer()
		
		// Now remove it
		return removePasswordInternal(fullPath, selected, noCommit)
	}
	
	// Normal rm without clip
	selected, err := InteractiveFuzzySearch(FuzzyModeRm)
	if err != nil {
		return err
	}
	if selected == "" {
		return nil
	}
	
	fullPath := getRmFullPath(selected)
	return removePasswordInternal(fullPath, selected, noCommit)
}

// getRmFullPath converts a password path to its full filesystem path
func getRmFullPath(path string) string {
	storeDir := GetPasswordStoreDir()
	normalized := filesystem.NormalizePath(path)
	if !strings.HasSuffix(normalized, ".gpg") {
		normalized += ".gpg"
	}
	return filepath.Join(storeDir, filepath.FromSlash(normalized))
}

// removePasswordInternal removes a password file without doing path validation
// (path is already validated in fuzzy search)
func removePasswordInternal(fullPath, displayPath string, noCommit bool) error {
	// Remove the file
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("pass: failed to remove %s: %v", displayPath, err)
	}

	// Git integration - remove and commit
	if !noCommit {
		// Check if this is a git repo
		storeDir := GetPasswordStoreDir()
		gitDir := filepath.Join(storeDir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			if err := git.RemoveAndCommit(fullPath, "Remove "+displayPath); err != nil {
				fmt.Fprintf(os.Stderr, "pass: warning: git operations failed: %v\n", err)
			}
		}
	}

	fmt.Printf("Password removed successfully.\n")
	return nil
}
