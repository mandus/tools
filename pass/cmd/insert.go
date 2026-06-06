package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mandu/tools/pass/pkg/gpg"
	"github.com/mandu/tools/pass/pkg/filesystem"
	"github.com/mandu/tools/pass/pkg/git"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert [path]",
	Short: "Insert a new password",
	Long:  `Insert a new password at the given path. Prompts for password twice for verification.`,
	Args:  cobra.ExactArgs(1),
	SilenceUsage: true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		return insertPassword(path)
	},
}

// Flags for insert command
var (
	echoFlag      bool
	multilineFlag bool
	noCommitFlag  bool
)

func addInsertCmd() {
	insertCmd.Flags().BoolVarP(&echoFlag, "echo", "e", false, "Echo password while typing")
	insertCmd.Flags().BoolVarP(&multilineFlag, "multiline", "m", false, "Allow multi-line password")
	insertCmd.Flags().BoolVar(&noCommitFlag, "no-commit", false, "Skip git commit after insert")
	rootCmd.AddCommand(insertCmd)
}

// insertPassword inserts a new password
func insertPassword(path string) error {
	// Normalize path and create file path
	storeDir := GetPasswordStoreDir()
	filePath := filesystem.NormalizePath(path)
	if !strings.HasSuffix(filePath, ".gpg") {
		filePath += ".gpg"
	}
	
	// Try multiple path strategies for cross-platform compatibility
	// Strategy 1: Use filepath.Join with FromSlash
	fullPath := filepath.Join(storeDir, filepath.FromSlash(filePath))
	
	// Check if file exists with strategy 1
	if _, err := os.Stat(fullPath); err == nil {
		return fmt.Errorf("pass: %s: Already exists", path)
	}
	
	// Strategy 2: Try with forward slashes directly
	storeDirForward := strings.ReplaceAll(storeDir, "\\", "/")
	fullPath = storeDirForward + "/" + filePath
	
	// Check if file exists with strategy 2
	if _, err := os.Stat(fullPath); err == nil {
		return fmt.Errorf("pass: %s: Already exists", path)
	}
	
	// Get password from user
	password, err := promptForPassword(path)
	if err != nil {
		return err
	}
	
	// Validate password is not empty
	if password == "" {
		return fmt.Errorf("pass: password cannot be empty")
	}
	
	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
		return fmt.Errorf("pass: failed to create directory: %v", err)
	}
	
	// Create temp file
	tempFile, err := os.CreateTemp("", "pass-*.tmp")
	if err != nil {
		return fmt.Errorf("pass: failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	// Write password to temp file
	if _, err := tempFile.WriteString(password); err != nil {
		tempFile.Close()
		return fmt.Errorf("pass: failed to write temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("pass: failed to close temp file: %v", err)
	}
	
	// Encrypt the temp file
	if err := gpg.EncryptFile(tempFile.Name(), fullPath); err != nil {
		return fmt.Errorf("pass: GPG encryption failed: %v", err)
	}
	
	// Securely delete temp file
	if err := filesystem.SecureDelete(tempFile.Name()); err != nil {
		// Non-fatal error
		fmt.Fprintf(os.Stderr, "pass: warning: failed to securely delete temp file: %v\n", err)
	}
	
	// Git integration
	if !noCommitFlag {
		if err := git.AddAndCommit(fullPath, "Add "+path); err != nil {
			fmt.Fprintf(os.Stderr, "pass: warning: git commit failed: %v\n", err)
		}
	}
	
	fmt.Println("Password inserted successfully.")
	return nil
}

// promptForPassword prompts the user for a password twice and verifies they match
func promptForPassword(path string) (string, error) {
	var password1, password2 string
	var err error
	
	// Get terminal state
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		// Fall back to regular input if terminal raw mode fails
		oldState = nil
	}
	defer func() {
		if oldState != nil {
			term.Restore(int(os.Stdin.Fd()), oldState)
		}
	}()
	
	// Prompt for password
	fmt.Printf("Enter password for %s: ", path)
	if echoFlag {
		reader := bufio.NewReader(os.Stdin)
		password1, err = reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("pass: failed to read password: %v", err)
		}
		password1 = strings.TrimSpace(password1)
	} else {
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", fmt.Errorf("pass: failed to read password: %v", err)
		}
		password1 = string(passwordBytes)
		fmt.Println() // Print newline after hidden input
	}
	
	// Prompt for confirmation
	fmt.Printf("Retype password for %s: ", path)
	if echoFlag {
		reader := bufio.NewReader(os.Stdin)
		password2, err = reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("pass: failed to read password confirmation: %v", err)
		}
		password2 = strings.TrimSpace(password2)
	} else {
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", fmt.Errorf("pass: failed to read password confirmation: %v", err)
		}
		password2 = string(passwordBytes)
		fmt.Println() // Print newline after hidden input
	}
	
	// Verify passwords match
	if password1 != password2 {
		return "", fmt.Errorf("pass: password verification failed")
	}
	
	return password1, nil
}
