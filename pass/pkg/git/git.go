// Package git provides git integration for the pass tool.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// InitRepo initializes a git repository in the given directory.
func InitRepo(dir string) error {
	// Check if .git already exists
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); !os.IsNotExist(err) {
		return nil // Already initialized
	}
	
	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git init failed: %v", err)
	}
	
	// Configure git user if not already set
	if err := configureGitUser(dir); err != nil {
		// Non-fatal warning
		fmt.Fprintf(os.Stderr, "pass: warning: failed to configure git user: %v\n", err)
	}
	
	return nil
}

// configureGitUser configures git user.name and user.email if not already set.
func configureGitUser(dir string) error {
	// Check if user.name is set in the repo (local config)
	cmd := exec.Command("git", "config", "--local", "user.name")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		// Try to set from environment or default
		name := os.Getenv("PASS_GIT_NAME")
		if name == "" {
			name = os.Getenv("USERNAME")
			if name == "" {
				name = "Password Store User"
			}
		}
		
		// Set user.name
		if err := runGitConfig(dir, "user.name", name); err != nil {
			return err
		}
	}
	
	// Check if user.email is set in the repo (local config)
	cmd = exec.Command("git", "config", "--local", "user.email")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		// Try to set from environment or default
		email := os.Getenv("PASS_GIT_EMAIL")
		if email == "" {
			email = os.Getenv("USERNAME") + "@localhost"
		}
		
		// Set user.email
		if err := runGitConfig(dir, "user.email", email); err != nil {
			return err
		}
	}
	
	return nil
}

// runGitConfig runs git config to set a value.
func runGitConfig(dir, key, value string) error {
	cmd := exec.Command("git", "config", key, value)
	cmd.Dir = dir
	return cmd.Run()
}

// AddAndCommit adds a file to git and commits it.
func AddAndCommit(filePath, message string) error {
	// Get directory of the file
	dir := filepath.Dir(filePath)
	
	// Add file
	cmd := exec.Command("git", "add", filepath.Base(filePath))
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git add failed: %v", err)
	}
	
	// Commit
	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %v", err)
	}
	
	return nil
}

// RemoveAndCommit removes a file from git and commits the removal.
func RemoveAndCommit(filePath, message string) error {
	// Get directory of the file
	dir := filepath.Dir(filePath)
	
	// Remove file from git
	cmd := exec.Command("git", "rm", filepath.Base(filePath))
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		// Non-fatal: file is already removed, just warn
		fmt.Fprintf(os.Stderr, "pass: warning: git rm failed: %v\n", err)
	}
	
	// Commit the removal
	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		// Non-fatal: warn but continue
		fmt.Fprintf(os.Stderr, "pass: warning: git commit failed: %v\n", err)
	}
	
	return nil
}

// CheckGit checks if git is installed and available.
func CheckGit() error {
	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git: command not found. Please install Git for Windows")
	}
	return nil
}
