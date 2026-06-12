// Package cmd contains the command-line interface for the pass tool.
package cmd

import (
	"os"

	"github.com/mandu/tools/pass/cmd/tui"
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
If called without arguments, it enters interactive fuzzy search mode.

Usage:
  pass [command] [options] [path]
  pass [options] <path>              # Show password (default command)
  pass                              # Interactive fuzzy search mode

Examples:
  pass insert email/gmail.com/user    Insert a new password
  pass email/gmail.com/user           Show a password
  pass -c email/gmail.com/user       Copy password to clipboard
  pass                              Interactive fuzzy search
  pass ls                            List all passwords
  pass find gmail                    Search for passwords
  pass rm                            Remove password with fuzzy search
  pass rm <path>                     Remove specific password`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no args, check for flags
		if len(args) == 0 {
			// Check if help or version flags were set - let cobra handle those
			if cmd.Flags().Changed("help") || cmd.Flags().Changed("version") {
				return nil
			}
			// Check if clip flag is set - enter fuzzy search mode with clip
			clipFlagChanged, _ := cmd.Flags().GetBool("clip")
			if clipFlagChanged {
				// Set global clip flag
				clipFlag = true
				// Enter fuzzy search mode with clip using new TUI
				selected, err := tui.RunInteractiveFuzzySearch(tui.FuzzyModeClip)
				if err != nil {
					return err
				}
				if selected != "" {
					return showPassword(selected)
				}
				return nil
			}
			// Enter fuzzy search mode (default: show) using new TUI
			selected, err := tui.RunInteractiveFuzzySearch(tui.FuzzyModeShow)
			if err != nil {
				return err
			}
			if selected != "" {
				return showPassword(selected)
			}
			return nil
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
	addRmCmd()
	addEditCmd()
	addGitCmd()

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

// GetPasswordStoreDir returns the password store directory path
// Cross-platform: uses USERPROFILE on Windows, HOME on Unix
func GetPasswordStoreDir() string {
	dir := os.Getenv("PASSWORD_STORE_DIR")
	if dir != "" {
		return dir
	}
	// Use USERPROFILE on Windows, HOME on Unix
	// Always use forward slashes for consistency with pass convention
	if home := os.Getenv("USERPROFILE"); home != "" {
		return home + "/.password-store"
	}
	if home := os.Getenv("HOME"); home != "" {
		return home + "/.password-store"
	}
	return ".password-store"
}

// IsClipboardFlagSet returns whether the -c/--clip flag is set
func IsClipboardFlagSet() bool {
	return clipFlag
}
