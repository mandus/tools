package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mandu/tools/pass/cmd/tui"
	"github.com/mandu/tools/pass/pkg/filesystem"
	"github.com/mandu/tools/pass/pkg/gpg"
	"github.com/mandu/tools/pass/pkg/git"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [OPTIONS] [<path>]",
	Short: "Edit a password",
	Long: `Edit an existing password. Decrypts the password, opens it in your editor,
and re-encrypts it when you save.

If a path is provided, edits that specific password.
If no path is provided, will enter interactive fuzzy search mode to select a password.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) > 0 {
			path = args[0]
		}

		// Get flags
		noCommit, _ := cmd.Flags().GetBool("no-commit")
		force, _ := cmd.Flags().GetBool("force")

		// --force is an alias for --no-commit
		if force {
			noCommit = true
		}

		if path != "" {
			return editPassword(path, noCommit)
		}

		// No path provided - enter fuzzy search mode
		return runEditFuzzySearch(noCommit)
	},
}

// Flags for edit command
var (
	noCommitFlagEdit bool
	forceFlagEdit    bool
)

func addEditCmd() {
	editCmd.Flags().BoolVarP(&noCommitFlagEdit, "no-commit", "n", false, "Skip git commit after editing")
	editCmd.Flags().BoolVarP(&forceFlagEdit, "force", "f", false, "Alias for --no-commit")
	rootCmd.AddCommand(editCmd)
}

// editPassword edits a password at the specified path
func editPassword(path string, noCommit bool) error {
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

	// Check if it's a directory
	info, err := os.Stat(fullPath)
	if err != nil {
		return fmt.Errorf("pass: failed to stat %s: %v", path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("pass: %s: Is a directory", path)
	}

	// Decrypt the file
	password, err := gpg.DecryptFile(fullPath)
	if err != nil {
		return err // err already has proper prefix from gpg package
	}

	// Create temp file
	tempFile, err := os.CreateTemp("", "pass-edit-*.tmp")
	if err != nil {
		return fmt.Errorf("pass: failed to create temp file: %v", err)
	}
	defer func() {
		// Always try to securely delete the temp file
		_ = filesystem.SecureDelete(tempFile.Name())
	}()

	// Set restrictive permissions on temp file (readable/writable only by owner)
	if err := tempFile.Chmod(0600); err != nil {
		// Non-fatal warning
		fmt.Fprintf(os.Stderr, "pass: warning: failed to set permissions on temp file: %v\n", err)
	}

	// Write password to temp file
	if _, err := tempFile.WriteString(password); err != nil {
		tempFile.Close()
		return fmt.Errorf("pass: failed to write temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("pass: failed to close temp file: %v", err)
	}

	// Open in editor
	editor := getEditor()
	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("pass: editor exited with status %d", exitErr.ExitCode())
		}
		return fmt.Errorf("pass: failed to open editor: %v", err)
	}

	// Read the modified content from temp file
	modifiedContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return fmt.Errorf("pass: failed to read modified content: %v", err)
	}

	// Convert to string
	modifiedPassword := string(modifiedContent)

	// Create another temp file for re-encryption
	tempFile2, err := os.CreateTemp("", "pass-edit-*.tmp")
	if err != nil {
		return fmt.Errorf("pass: failed to create temp file for encryption: %v", err)
	}
	defer func() {
		// Always try to securely delete the second temp file
		_ = filesystem.SecureDelete(tempFile2.Name())
	}()

	// Set restrictive permissions on second temp file
	if err := tempFile2.Chmod(0600); err != nil {
		// Non-fatal warning
		fmt.Fprintf(os.Stderr, "pass: warning: failed to set permissions on temp file: %v\n", err)
	}

	// Write modified content to second temp file
	if _, err := tempFile2.WriteString(modifiedPassword); err != nil {
		tempFile2.Close()
		return fmt.Errorf("pass: failed to write temp file for encryption: %v", err)
	}
	if err := tempFile2.Close(); err != nil {
		return fmt.Errorf("pass: failed to close temp file for encryption: %v", err)
	}

	// Encrypt the modified content back to the original file
	if err := gpg.EncryptFile(tempFile2.Name(), fullPath); err != nil {
		return fmt.Errorf("pass: GPG encryption failed: %v", err)
	}

	// Git integration
	if !noCommit {
		// Check if this is a git repo
		gitDir := filepath.Join(storeDir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			// Run git add and commit
			// We need to run git from the store directory
			if err := git.AddAndCommit(fullPath, "Edit "+path); err != nil {
				fmt.Fprintf(os.Stderr, "pass: warning: git operations failed: %v\n", err)
			}
		}
	}

	fmt.Println("Password updated successfully.")
	return nil
}

// getEditor returns the editor command to use
func getEditor() string {
	editor := os.Getenv("EDITOR")
	if editor != "" {
		return editor
	}
	// Platform-specific defaults
	if runtime.GOOS == "windows" {
		return "notepad"
	}
	return "vi"
}

// runEditFuzzySearch enters interactive fuzzy search mode for editing a password
func runEditFuzzySearch(noCommit bool) error {
	// Use the TUI for fuzzy search
	selected, err := tui.RunInteractiveFuzzySearch(tui.FuzzyModeEdit)
	if err != nil {
		return err
	}
	if selected == "" {
		return nil
	}

	return editPassword(selected, noCommit)
}
