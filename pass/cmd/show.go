package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mandu/tools/pass/pkg/gpg"
	"github.com/mandu/tools/pass/pkg/filesystem"
	"github.com/spf13/cobra"
)

// showCmd represents the show command (also the default command)
var showCmd = &cobra.Command{
	Use:   "show [path]",
	Short: "Show a password",
	Long:  `Show the password stored at the given path. Decrypts the GPG file and prints to stdout.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		return showPassword(path)
	},
}

func addShowCmd() {
	rootCmd.AddCommand(showCmd)
}

// showPassword retrieves and displays a password
func showPassword(path string) error {
	storeDir := GetPasswordStoreDir()
	
	// Normalize the input path: strip .gpg if present, ensure forward slashes
	filePath := strings.ReplaceAll(path, "\\", "/")
	if strings.HasSuffix(filePath, ".gpg") {
		filePath = strings.TrimSuffix(filePath, ".gpg")
	}
	filePath += ".gpg"
	
	// Construct full path - try multiple strategies for cross-platform compatibility
	// Strategy 1: Use filepath.Join with FromSlash (handles path separator conversion)
	fullPath := filepath.Join(storeDir, filepath.FromSlash(filePath))
	
	// Check if file exists
	if _, err := os.Stat(fullPath); err == nil {
		return decryptAndDisplay(fullPath, path)
	}
	
	// Strategy 2: Try with forward slashes directly (for Unix-like environments on Windows)
	fullPath = storeDir + "/" + filePath
	if _, err := os.Stat(fullPath); err == nil {
		return decryptAndDisplay(fullPath, path)
	}
	
	// Strategy 3: Try storeDir with forward slashes + path
	storeDirForward := strings.ReplaceAll(storeDir, "\\", "/")
	fullPath = storeDirForward + "/" + filePath
	if _, err := os.Stat(fullPath); err == nil {
		return decryptAndDisplay(fullPath, path)
	}
	
	return fmt.Errorf("pass: %s: No such file or directory", path)
}

// decryptAndDisplay handles the decryption and output
func decryptAndDisplay(fullPath, displayPath string) error {
	// Decrypt the file
	password, err := gpg.DecryptFile(fullPath)
	if err != nil {
		return err // err already has proper prefix from gpg package
	}
	
	// Output based on flags
	if IsClipboardFlagSet() {
		// Copy to clipboard
		if err := filesystem.CopyToClipboard(password); err != nil {
			return fmt.Errorf("pass: failed to copy to clipboard: %v", err)
		}
		fmt.Printf("Copied %s to clipboard.\n", displayPath)
		// Start auto-clear timer in background
		go filesystem.StartClipboardClearTimer()
	} else {
		// Print to stdout
		fmt.Print(password)
	}
	
	return nil
}
