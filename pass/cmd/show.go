package cmd

import (
	"fmt"
	"os"

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
	fullPath, err := getPasswordFullPath(path)
	if err != nil {
		return err
	}
	return decryptAndDisplay(fullPath, path)
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
		// Explicitly flush stdout to avoid buffering delay
		// This is especially important for passwords without trailing newlines
		// that are piped to other commands
		_ = os.Stdout.Sync()
	}
	
	return nil
}
