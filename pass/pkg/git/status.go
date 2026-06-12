// Package git provides git integration for the pass tool.
// This file contains git status checking functionality.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitStatus represents the git status of a repository
type GitStatus struct {
	// IsGitRepo indicates whether the directory is a git repository
	IsGitRepo bool
	// IsClean indicates whether there are no uncommitted changes
	IsClean bool
	// HasUncommitted indicates whether there are uncommitted changes (staged or unstaged)
	HasUncommitted bool
	// Ahead is the number of commits ahead of the remote tracking branch
	Ahead int
	// Behind is the number of commits behind the remote tracking branch
	Behind int
	// Branch is the current branch name (or "HEAD" if detached)
	Branch string
	// Remote is the remote name being tracked (usually "origin")
	Remote string
	// TrackingBranch is the remote branch being tracked
	TrackingBranch string
	// Error contains any error encountered while checking status
	Error error
}

// String returns a human-readable representation of the git status
func (gs GitStatus) String() string {
	if gs.Error != nil {
		return fmt.Sprintf("Error: %v", gs.Error)
	}
	
	if !gs.IsGitRepo {
		return "Not a git repository"
	}
	
	var parts []string
	
	// Branch info
	if gs.Branch != "" {
		parts = append(parts, gs.Branch)
	} else {
		parts = append(parts, "HEAD")
	}
	
	// Sync status
	if gs.Ahead > 0 && gs.Behind > 0 {
		parts = append(parts, fmt.Sprintf("⬆%d⬇%d", gs.Ahead, gs.Behind))
	} else if gs.Ahead > 0 {
		parts = append(parts, fmt.Sprintf("⬆%d", gs.Ahead))
	} else if gs.Behind > 0 {
		parts = append(parts, fmt.Sprintf("⬇%d", gs.Behind))
	} else if gs.IsClean && !gs.HasUncommitted {
		parts = append(parts, "=")
	}
	
	// Uncommitted changes
	if gs.HasUncommitted {
		parts = append(parts, "*")
	}
	
	return strings.Join(parts, " ")
}

// DetailedString returns a detailed human-readable representation
func (gs GitStatus) DetailedString() string {
	if gs.Error != nil {
		return fmt.Sprintf("Error: %v", gs.Error)
	}
	
	if !gs.IsGitRepo {
		return "Not a git repository"
	}
	
	var lines []string
	lines = append(lines, fmt.Sprintf("Branch: %s", gs.Branch))
	
	if gs.Remote != "" && gs.TrackingBranch != "" {
		lines = append(lines, fmt.Sprintf("Tracking: %s/%s", gs.Remote, gs.TrackingBranch))
	}
	
	if gs.Ahead > 0 {
		lines = append(lines, fmt.Sprintf("Ahead: %d commit(s)", gs.Ahead))
	}
	if gs.Behind > 0 {
		lines = append(lines, fmt.Sprintf("Behind: %d commit(s)", gs.Behind))
	}
	
	if gs.HasUncommitted {
		lines = append(lines, "Uncommitted changes: yes")
	} else {
		lines = append(lines, "Uncommitted changes: no")
	}
	
	if gs.IsClean {
		lines = append(lines, "Status: clean")
	} else {
		lines = append(lines, "Status: dirty")
	}
	
	return strings.Join(lines, "\n")
}

// GetGitStatus retrieves the git status for a given directory
func GetGitStatus(dir string) GitStatus {
	status := GitStatus{
		IsGitRepo: false,
		IsClean:   true,
	}
	
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		status.Error = fmt.Errorf("directory does not exist: %s", dir)
		return status
	}
	
	// Check if .git directory exists
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		// Check for .git file (submodule or worktree)
		gitFile := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitFile); os.IsNotExist(err) {
			status.IsGitRepo = false
			return status
		}
	}
	
	status.IsGitRepo = true
	
	// Get current branch
	branch, err := getCurrentBranch(dir)
	if err != nil {
		status.Error = fmt.Errorf("failed to get current branch: %v", err)
		return status
	}
	status.Branch = branch
	
	// Check for uncommitted changes
	hasUncommitted, err := hasUncommittedChanges(dir)
	if err != nil {
		// Non-fatal error, continue with other checks
		status.Error = err
	} else {
		status.HasUncommitted = hasUncommitted
		status.IsClean = !hasUncommitted
	}
	
	// Get sync status with remote
	ahead, behind, remote, trackingBranch, err := getSyncStatus(dir, branch)
	if err != nil {
		// Non-fatal error, continue
		if status.Error == nil {
			status.Error = err
		}
	} else {
		status.Ahead = ahead
		status.Behind = behind
		status.Remote = remote
		status.TrackingBranch = trackingBranch
	}
	
	return status
}

// getCurrentBranch returns the current branch name or "HEAD" if detached
func getCurrentBranch(dir string) (string, error) {
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		// Might be detached HEAD
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Check if it's a detached HEAD
			cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
			cmd.Dir = dir
			if _, err := cmd.Output(); err == nil {
				return "HEAD", nil
			}
			return "", fmt.Errorf("failed to get branch: %v", exitErr)
		}
		return "", fmt.Errorf("failed to get branch: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// hasUncommittedChanges checks if there are uncommitted changes (staged or unstaged)
func hasUncommittedChanges(dir string) (bool, error) {
	// Use git status --porcelain=v1 which provides machine-readable output
	cmd := exec.Command("git", "status", "--porcelain=v1")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check uncommitted changes: %v", err)
	}
	return len(output) > 0, nil
}

// getSyncStatus returns the ahead/behind counts and remote tracking info
func getSyncStatus(dir, branch string) (ahead, behind int, remote, trackingBranch string, err error) {
	// If detached HEAD, we can't get sync status
	if branch == "HEAD" {
		return 0, 0, "", "", nil
	}
	
	// Get the tracking branch info
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", branch+"@{upstream}")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		// No upstream configured - this is normal
		return 0, 0, "", "", nil
	}
	
	upstream := strings.TrimSpace(string(output))
	if upstream == "" {
		return 0, 0, "", "", nil
	}
	
	// Parse upstream to get remote and branch
	// Format is typically "origin/main" or "remote-name/branch-name"
	parts := strings.Split(upstream, "/")
	if len(parts) >= 2 {
		remote = parts[0]
		trackingBranch = strings.Join(parts[1:], "/")
	} else {
		remote = upstream
		trackingBranch = branch
	}
	
	// Use git rev-list to count commits ahead and behind
	// Ahead: commits in branch not in upstream
	cmd = exec.Command("git", "rev-list", "--left-right", branch+"..."+upstream)
	cmd.Dir = dir
	output, err = cmd.Output()
	if err != nil {
		return 0, 0, remote, trackingBranch, fmt.Errorf("failed to count commits: %v", err)
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	ahead = 0
	behind = 0
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line[0] == '<' {
			behind++
		} else if line[0] == '>' {
			ahead++
		}
	}
	
	return ahead, behind, remote, trackingBranch, nil
}

// CheckRemoteConfigured checks if a remote is configured for the repository
func CheckRemoteConfigured(dir string) (bool, string, error) {
	cmd := exec.Command("git", "remote")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return false, "", fmt.Errorf("failed to check remotes: %v", err)
	}
	
	remotes := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(remotes) == 0 || (len(remotes) == 1 && remotes[0] == "") {
		return false, "", nil
	}
	
	// Return the first remote (usually "origin")
	return true, strings.TrimSpace(remotes[0]), nil
}

// Push pushes changes to the configured remote
func Push(dir string) error {
	// Check if remote is configured
	hasRemote, remote, err := CheckRemoteConfigured(dir)
	if err != nil {
		return fmt.Errorf("failed to check remote: %v", err)
	}
	if !hasRemote {
		return fmt.Errorf("no remote configured. Use 'git remote add' to configure a remote")
	}
	
	// Get current branch
	branch, err := getCurrentBranch(dir)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}
	
	// Push to remote
	cmd := exec.Command("git", "push", remote, branch)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git push failed: %v\n%s", err, string(output))
	}
	
	return nil
}

// Update fetches and merges changes from the remote
func Update(dir string) error {
	// Check if remote is configured
	hasRemote, remote, err := CheckRemoteConfigured(dir)
	if err != nil {
		return fmt.Errorf("failed to check remote: %v", err)
	}
	if !hasRemote {
		return fmt.Errorf("no remote configured. Use 'git remote add' to configure a remote")
	}
	
	// Get current branch
	branch, err := getCurrentBranch(dir)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}
	
	// Pull from remote (fetch + merge)
	cmd := exec.Command("git", "pull", remote, branch)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if it's a merge conflict
		outputStr := string(output)
		if strings.Contains(outputStr, "CONFLICT") || strings.Contains(outputStr, "merge conflict") {
			return fmt.Errorf("merge conflict detected. Resolve conflicts and commit, or use 'git merge --abort'")
		}
		return fmt.Errorf("git pull failed: %v\n%s", err, string(output))
	}
	
	return nil
}

// InitGitRepo initializes a git repository in the given directory
func InitGitRepo(dir string) error {
	// Check if already initialized
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
		// Non-fatal warning - we'll try to set reasonable defaults
		name := os.Getenv("PASS_GIT_NAME")
		if name == "" {
			name = os.Getenv("USERNAME")
			if name == "" {
				name = os.Getenv("USER")
				if name == "" {
					name = "Password Store User"
				}
			}
		}
		
		email := os.Getenv("PASS_GIT_EMAIL")
		if email == "" {
			email = os.Getenv("USERNAME") + "@localhost"
			if email == "@localhost" {
				email = os.Getenv("USER") + "@localhost"
			}
		}
		
		if err := runGitConfig(dir, "user.name", name); err != nil {
			// Still non-fatal
			return nil
		}
		if err := runGitConfig(dir, "user.email", email); err != nil {
			return nil
		}
	}
	
	return nil
}
