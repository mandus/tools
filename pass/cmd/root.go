// Package cmd contains the command-line interface for the pass tool.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pass [path]",
	Short: "A Windows-compatible password store manager",
	Long: `pass - A Windows-compatible replacement for the Unix password-store tool.

pass stores passwords as GPG-encrypted files in ~/.password-store/
and integrates with git for version control.

If called with a path argument and no command, it shows the password (same as 'pass show').

Usage:
  pass [command] [options] [path]
  pass [options] <path>              # Show password (default command)

Examples:
  pass insert email/gmail.com/user    Insert a new password
  pass email/gmail.com/user           Show a password
  pass -c email/gmail.com/user       Copy password to clipboard
  pass ls                            List all passwords
  pass find gmail                    Search for passwords`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no args and no command flags, show help
		if len(args) == 0 {
			// Check if any flags were set
			if cmd.Flags().Changed("help") || cmd.Flags().Changed("version") || cmd.Flags().Changed("clip") {
				return nil // Let cobra handle these flags
			}
			return cmd.Help()
		}
		// If args provided without explicit command, treat as show command
		return showPassword(args[0])
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
