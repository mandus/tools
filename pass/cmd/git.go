package cmd

import (
	"fmt"
	"os"

	"github.com/mandu/tools/pass/pkg/git"
	"github.com/spf13/cobra"
)

// GitStatus represents the git status display information
type GitStatus struct {
	IsGitRepo      bool
	Branch         string
	Remote         string
	TrackingBranch string
	Ahead          int
	Behind         int
	HasUncommitted bool
	IsClean        bool
	Error          error
}

var gitCmd = &cobra.Command{
	Use:   "git [command]",
	Short: "Manage git repository for the password store",
	Long: `Manage the git repository for your password store.

This command provides git integration for the password store, allowing you
to check the sync status and perform git operations.

If called without a subcommand, it shows the current git status.

Usage:
  pass git           Show git status
  pass git push      Push changes to remote
  pass git update    Pull changes from remote
  pass git status    Show git status (alias for 'pass git')

Examples:
  pass git              # Show git status
  pass git push         # Push changes to remote
  pass git update       # Pull changes from remote
  pass git -v           # Show detailed git status`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand, show status
		return runGitStatus(cmd, args)
	},
}

var gitStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show git status of the password store",
	Long:  "Show the git status of the password store, including branch, sync status, and uncommitted changes.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGitStatus(cmd, args)
	},
}

var gitPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push changes to the remote repository",
	Long: `Push local changes to the configured remote repository.

This command pushes all committed changes to the remote repository.
Uncommitted changes will not be pushed.

Examples:
  pass git push    # Push changes to remote`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		storeDir := GetPasswordStoreDir()
		
		// Check if store exists
		if _, err := os.Stat(storeDir); os.IsNotExist(err) {
			return fmt.Errorf("pass: password store does not exist. Use 'pass insert' to create it.")
		}
		
		// Check if it's a git repo
		if _, err := os.Stat(storeDir + "/.git"); os.IsNotExist(err) {
			return fmt.Errorf("pass: password store is not a git repository. Use 'pass git init' or initialize git manually.")
		}
		
		fmt.Println("Pushing changes to remote...")
		if err := git.Push(storeDir); err != nil {
			return fmt.Errorf("pass: git push failed: %v", err)
		}
		fmt.Println("Successfully pushed changes to remote.")
		return nil
	},
}

var gitUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull changes from the remote repository",
	Long: `Pull changes from the configured remote repository.

This command fetches changes from the remote and merges them into your
local password store.

Examples:
  pass git update    # Pull changes from remote`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		storeDir := GetPasswordStoreDir()
		
		// Check if store exists
		if _, err := os.Stat(storeDir); os.IsNotExist(err) {
			return fmt.Errorf("pass: password store does not exist. Use 'pass insert' to create it.")
		}
		
		// Check if it's a git repo
		if _, err := os.Stat(storeDir + "/.git"); os.IsNotExist(err) {
			return fmt.Errorf("pass: password store is not a git repository. Use 'pass git init' or initialize git manually.")
		}
		
		fmt.Println("Pulling changes from remote...")
		if err := git.Update(storeDir); err != nil {
			return fmt.Errorf("pass: git update failed: %v", err)
		}
		fmt.Println("Successfully updated from remote.")
		return nil
	},
}

var gitInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a git repository in the password store",
	Long: `Initialize a git repository in the password store directory.

This command creates a git repository in your password store directory
if one doesn't already exist.

Examples:
  pass git init    # Initialize git repository`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		storeDir := GetPasswordStoreDir()
		
		// Create store directory if it doesn't exist
		if err := os.MkdirAll(storeDir, 0700); err != nil {
			return fmt.Errorf("pass: failed to create password store directory: %v", err)
		}
		
		fmt.Println("Initializing git repository...")
		if err := git.InitGitRepo(storeDir); err != nil {
			return fmt.Errorf("pass: git init failed: %v", err)
		}
		fmt.Println("Successfully initialized git repository in password store.")
		return nil
	},
}

// runGitStatus runs the git status command
func runGitStatus(cmd *cobra.Command, args []string) error {
	storeDir := GetPasswordStoreDir()
	
	// Check if store exists
	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		return fmt.Errorf("pass: password store does not exist. Use 'pass insert' to create it.")
	}
	
	// Get git status
	status := git.GetGitStatus(storeDir)
	
	// Check for verbose flag
	verbose, _ := cmd.Flags().GetBool("verbose")
	
	if status.Error != nil {
		if !status.IsGitRepo {
			return fmt.Errorf("pass: password store is not a git repository.\nHint: Run 'pass git init' to initialize git.")
		}
		return fmt.Errorf("pass: failed to get git status: %v", status.Error)
	}
	
	if !status.IsGitRepo {
		return fmt.Errorf("pass: password store is not a git repository.\nHint: Run 'pass git init' to initialize git.")
	}
	
	if verbose {
		// Show detailed status
		fmt.Println(status.DetailedString())
	} else {
		// Show compact status
		fmt.Printf("Git status: %s\n", status.String())
		
		// Add additional info
		if status.Ahead > 0 {
			fmt.Printf("  %d commit(s) ahead of %s/%s\n", status.Ahead, status.Remote, status.TrackingBranch)
		}
		if status.Behind > 0 {
			fmt.Printf("  %d commit(s) behind %s/%s\n", status.Behind, status.Remote, status.TrackingBranch)
		}
		if status.HasUncommitted {
			fmt.Println("  Uncommitted changes present")
		}
		if status.IsClean && !status.HasUncommitted && status.Ahead == 0 && status.Behind == 0 {
			fmt.Println("  Up to date")
		}
	}
	
	return nil
}

// addGitCmd adds the git command and its subcommands to the root command
func addGitCmd() {
	// Add verbose flag to git command
	gitCmd.PersistentFlags().BoolP("verbose", "v", false, "Show detailed git status information")
	
	// Add subcommands
	gitCmd.AddCommand(gitStatusCmd)
	gitCmd.AddCommand(gitPushCmd)
	gitCmd.AddCommand(gitUpdateCmd)
	gitCmd.AddCommand(gitInitCmd)
	
	// Add git command to root
	rootCmd.AddCommand(gitCmd)
}
