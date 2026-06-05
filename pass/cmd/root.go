// Package cmd contains the command-line interface for the pass tool.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pass",
	Short: "A Windows-compatible password store manager",
	Long: `pass - A Windows-compatible replacement for the Unix password-store tool.

pass stores passwords as GPG-encrypted files in ~/.password-store/
and integrates with git for version control.

Usage:
  pass [command] [options] [path]

Examples:
  pass insert email/gmail.com/user    Insert a new password
  pass email/gmail.com/user           Show a password
  pass -c email/gmail.com/user       Copy password to clipboard
  pass ls                            List all passwords
  pass find gmail                    Search for passwords`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no command and no args, show help
		if len(args) == 0 {
			return cmd.Help()
		}
		// If args provided without command, treat as show command
		// This will be handled by the show command registration
		return nil
	},
}

// clipFlag is a global flag for copying to clipboard
var clipFlag bool

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	// Add global flags
	rootCmd.PersistentFlags().BoolVarP(&clipFlag, "clip", "c", false, "Copy password to clipboard instead of stdout")

	// Add version flag
	rootCmd.Version = "0.1.0"

	// Add all commands
	addInsertCmd()
	addShowCmd()
	addLsCmd()
	addFindCmd()

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

// GetPasswordStoreDir returns the password store directory path
func GetPasswordStoreDir() string {
	dir := os.Getenv("PASSWORD_STORE_DIR")
	if dir != "" {
		return dir
	}
	// Default to %USERPROFILE%\.password-store on Windows
	// or ~/.password-store on Unix
	return fmt.Sprintf("%s/.password-store", os.Getenv("USERPROFILE"))
}

// IsClipboardFlagSet returns whether the -c/--clip flag is set
func IsClipboardFlagSet() bool {
	return clipFlag
}
