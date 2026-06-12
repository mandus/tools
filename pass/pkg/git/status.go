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
	// HasMergeConflict indicates whether there are merge conflicts
	HasMergeConflict bool
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
	
	// Merge conflict indicator
	if gs.HasMergeConflict {
		parts = append(parts, "!")
	}
	
	// Sync status - use simple symbols like git-prompt
	if gs.Ahead > 0 && gs.Behind > 0 {
		parts = append(parts, "<>")
	} else if gs.Ahead > 0 {
		parts = append(parts, ">")
	} else if gs.Behind > 0 {
		parts = append(parts, "<")
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
	
	if gs.HasMergeConflict {
		lines = append(lines, "Merge conflicts: yes")
	} else {
		lines = append(lines, "Merge conflicts: no")
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
	
	// Check for merge conflicts
	hasMergeConflict, err := hasMergeConflicts(dir)
	if err != nil {
		// Non-fatal
		if status.Error == nil {
			status.Error = err
		}
	} else {
		status.HasMergeConflict = hasMergeConflict
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

// hasMergeConflicts checks if there are merge conflicts
func hasMergeConflicts(dir string) (bool, error) {
	// Check for merge conflict markers or use git status
	cmd := exec.Command("git", "status", "--porcelain=v1")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check merge conflicts: %v", err)
	}
	
	// Check for conflict markers in git status output
	outputStr := string(output)
	if strings.Contains(outputStr, "UU") || strings.Contains(outputStr, "AA") ||
		strings.Contains(outputStr, "DD") || strings.Contains(outputStr, "both modified") ||
		strings.Contains(outputStr, "both added") || strings.Contains(outputStr, "both deleted") {
		return true, nil
	}
	
	return false, nil
}

// gitConfig holds git branch configuration
type gitConfig struct {
	Remote string
	Merge  string
}

// getGitConfig retrieves branch remote and merge configuration
func getGitConfig(dir, branch string) (gitConfig, error) {
	var cfg gitConfig
	
	// Get remote
	cmd := exec.Command("git", "config", "branch."+branch+".remote")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return cfg, err
	}
	cfg.Remote = strings.TrimSpace(string(output))
	
	// Get merge
	cmd = exec.Command("git", "config", "branch."+branch+".merge")
	cmd.Dir = dir
	output, err = cmd.Output()
	if err != nil {
		return cfg, err
	}
	cfg.Merge = strings.TrimPrefix(strings.TrimSpace(string(output)), "refs/heads/")
	
	return cfg, nil
}

// getUpstreamFromRef tries to get upstream from @{upstream} ref
func getUpstreamFromRef(dir, branch string) (gitConfig, error) {
	var cfg gitConfig
	
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", branch+"@{upstream}")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return cfg, err
	}
	
	upstream := strings.TrimSpace(string(output))
	if upstream == "" {
		return cfg, fmt.Errorf("no upstream")
	}
	
	// Strip refs/heads/ prefix if present
	upstream = strings.TrimPrefix(upstream, "refs/heads/")
	
	// Parse upstream to get remote and branch
	parts := strings.Split(upstream, "/")
	if len(parts) >= 2 {
		cfg.Remote = parts[0]
		cfg.Merge = strings.Join(parts[1:], "/")
	} else {
		cfg.Remote = upstream
		cfg.Merge = branch
	}
	
	return cfg, nil
}

// getSyncStatus returns the ahead/behind counts and remote tracking info
// Uses multiple strategies to determine sync status:
// 1. Try git ls-remote (queries actual remote state)
// 2. Fallback to git fetch --dry-run (simulates fetch, updates remote refs)
// 3. Fallback to local remote refs (cached, might be stale)
func getSyncStatus(dir, branch string) (ahead, behind int, remote, trackingBranch string, err error) {
	// If detached HEAD, we can't get sync status
	if branch == "HEAD" {
		return 0, 0, "", "", nil
	}
	
	var cfg gitConfig
	
	// First, try to get the upstream from git config (branch.<branch>.remote and branch.<branch>.merge)
	cfg, err = getGitConfig(dir, branch)
	if err != nil {
		// Non-fatal, try the @{upstream} approach
		cfg, err = getUpstreamFromRef(dir, branch)
		if err != nil {
			return 0, 0, "", "", nil
		}
	}
	
	if cfg.Remote == "" || cfg.Merge == "" {
		return 0, 0, "", "", nil
	}
	
	remote = cfg.Remote
	trackingBranch = cfg.Merge
	
	// Get local hash
	localHash, err := getCommitHash(dir, branch)
	if err != nil {
		return 0, 0, remote, trackingBranch, fmt.Errorf("failed to get local hash: %v", err)
	}
	
	// Try to get remote hash using multiple strategies
	var remoteHash string
	
	// Strategy 1: Use git ls-remote (most reliable, queries actual remote)
	remoteHash, err = getRemoteCommitHash(dir, remote, trackingBranch)
	if err == nil && remoteHash != "" {
		// Success with ls-remote
	} else {
		// Strategy 2: Try git fetch --dry-run to update remote refs
		if err := fetchRemote(dir, remote); err == nil {
			// After fetch, try local remote ref again
			remoteHash, err = getRemoteRefHash(dir, remote, trackingBranch)
		}
		
		// Strategy 3: Fallback to local remote ref without fetch
		if err != nil || remoteHash == "" {
			remoteHash, err = getRemoteRefHash(dir, remote, trackingBranch)
			if err != nil {
				// Can't determine sync status without remote hash
				return 0, 0, remote, trackingBranch, nil
			}
		}
	}
	
	// If hashes are the same, we're in sync
	if localHash == remoteHash {
		return 0, 0, remote, trackingBranch, nil
	}
	
	// Use git rev-list to count commits ahead and behind
	ahead, behind, err = countAheadBehind(dir, localHash, remoteHash)
	if err != nil {
		// Non-fatal error - just return 0,0 and the remote info we have
		// This can happen if git rev-list fails (e.g., repository issues)
		return 0, 0, remote, trackingBranch, nil
	}
	
	return ahead, behind, remote, trackingBranch, nil
}

// fetchRemote performs a dry-run fetch to update remote refs
func fetchRemote(dir, remote string) error {
	cmd := exec.Command("git", "fetch", "--dry-run", remote)
	cmd.Dir = dir
	return cmd.Run()
}

// getCommitHash gets the commit hash for a branch
func getCommitHash(dir, branch string) (string, error) {
	cmd := exec.Command("git", "rev-parse", branch)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getRemoteCommitHash uses git ls-remote to get the actual remote commit hash
func getRemoteCommitHash(dir, remote, branch string) (string, error) {
	// Strip refs/heads/ prefix if present (git config returns branch.merge as refs/heads/<name>)
	branchRef := strings.TrimPrefix(branch, "refs/heads/")
	cmd := exec.Command("git", "ls-remote", remote, "refs/heads/"+branchRef)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	// Output format: <hash>\trefs/heads/<branch>
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) >= 1 {
			return strings.TrimSpace(parts[0]), nil
		}
	}
	
	return "", fmt.Errorf("no hash found for %s/%s", remote, branch)
}

// getRemoteRefHash gets the remote hash from local remote reference
func getRemoteRefHash(dir, remote, branch string) (string, error) {
	remoteRefName := "refs/remotes/" + remote + "/" + branch
	cmd := exec.Command("git", "rev-parse", remoteRefName)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// countAheadBehind counts commits ahead and behind using git rev-list
// With --left-right flag:
// - '<' means commit is reachable only from first argument (local) = ahead
// - '>' means commit is reachable only from second argument (remote) = behind
func countAheadBehind(dir, localHash, remoteHash string) (int, int, error) {
	cmd := exec.Command("git", "rev-list", "--left-right", localHash+"..."+remoteHash)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	ahead := 0
	behind := 0
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// '<' = reachable only from local (first arg) = ahead
		// '>' = reachable only from remote (second arg) = behind
		if len(line) > 0 && line[0] == '<' {
			ahead++
		} else if len(line) > 0 && line[0] == '>' {
			behind++
		}
	}
	
	return ahead, behind, nil
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
		if strings.Contains(outputStr, "CONFLICT") || strings.Contains(outputStr, "merge conflict") ||
			strings.Contains(outputStr, "both modified") || strings.Contains(outputStr, "both added") ||
			strings.Contains(outputStr, "both deleted") {
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
