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
	
	// Normalize path: replace / with OS separator, add .gpg extension
	filePath := filesystem.NormalizePath(path)
	if !strings.HasSuffix(filePath, ".gpg") {
		filePath += ".gpg"
	}
	fullPath := filepath.Join(storeDir, filePath)
	
	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("pass: %s: No such file or directory", path)
	}
	
	// Decrypt the file
	password, err := gpg.DecryptFile(fullPath)
	if err != nil {
		return fmt.Errorf("pass: GPG decryption failed: %v", err)
	}
	
	// Output based on flags
	if IsClipboardFlagSet() {
		// Copy to clipboard
		if err := filesystem.CopyToClipboard(password); err != nil {
			return fmt.Errorf("pass: failed to copy to clipboard: %v", err)
		}
		fmt.Printf("Copied %s to clipboard.\n", path)
		// Start auto-clear timer in background
		go filesystem.StartClipboardClearTimer()
	} else {
		// Print to stdout
		fmt.Print(password)
	}
	
	return nil
}
