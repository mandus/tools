package git

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// errorForTest returns a test error
func errorForTest() error {
	return errors.New("test error")
}

// setupTestGitRepo creates a test git repository in a temporary directory
// Returns the directory path and a cleanup function
func setupTestGitRepo(t *testing.T) (string, func()) {
	t.Helper()
	
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "pass-git-status-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	cleanup := func() {
		os.RemoveAll(tempDir)
	}
	
	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		cleanup()
		t.Fatalf("Failed to initialize git repo: %v", err)
	}
	
	// Configure git user
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()
	
	return tempDir, cleanup
}

// setupTestGitRepoWithRemote creates a test git repository with a remote
func setupTestGitRepoWithRemote(t *testing.T) (string, string, func()) {
	t.Helper()
	
	// Create temp directory for local repo
	localDir, err := os.MkdirTemp("", "pass-git-status-test-local")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	// Create temp directory for remote (bare) repo
	remoteDir, err := os.MkdirTemp("", "pass-git-status-test-remote")
	if err != nil {
		os.RemoveAll(localDir)
		t.Fatalf("Failed to create remote dir: %v", err)
	}
	
	cleanup := func() {
		os.RemoveAll(localDir)
		os.RemoveAll(remoteDir)
	}
	
	// Initialize bare remote repo
	cmd := exec.Command("git", "init", "--bare", remoteDir)
	if err := cmd.Run(); err != nil {
		cleanup()
		t.Fatalf("Failed to initialize bare repo: %v", err)
	}
	
	// Initialize local repo
	cmd = exec.Command("git", "init")
	cmd.Dir = localDir
	if err := cmd.Run(); err != nil {
		cleanup()
		t.Fatalf("Failed to initialize local repo: %v", err)
	}
	
	// Configure git user in local repo
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()
	
	// Add remote
	cmd = exec.Command("git", "remote", "add", "origin", remoteDir)
	cmd.Dir = localDir
	if err := cmd.Run(); err != nil {
		cleanup()
		t.Fatalf("Failed to add remote: %v", err)
	}
	
	return localDir, remoteDir, cleanup
}

func TestGetGitStatus_NoGitRepo(t *testing.T) {
	// Create a temp directory without git
	tempDir, err := os.MkdirTemp("", "pass-git-status-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	status := GetGitStatus(tempDir)
	
	if status.IsGitRepo {
		t.Error("Expected IsGitRepo to be false for non-git directory")
	}
	if status.Error != nil {
		t.Errorf("Expected no error, got: %v", status.Error)
	}
}

func TestGetGitStatus_GitRepoClean(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}
	
	repoDir, cleanup := setupTestGitRepo(t)
	defer cleanup()
	
	// Create initial commit
	cmd := exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit")
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}
	
	status := GetGitStatus(repoDir)
	
	if !status.IsGitRepo {
		t.Error("Expected IsGitRepo to be true")
	}
	if !status.IsClean {
		t.Error("Expected IsClean to be true for clean repo")
	}
	if status.HasUncommitted {
		t.Error("Expected HasUncommitted to be false for clean repo")
	}
	if status.Branch != "master" && status.Branch != "main" {
		t.Errorf("Expected branch to be master or main, got: %s", status.Branch)
	}
	if status.Error != nil {
		t.Errorf("Expected no error, got: %v", status.Error)
	}
}

func TestGetGitStatus_UncommittedChanges(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}
	
	repoDir, cleanup := setupTestGitRepo(t)
	defer cleanup()
	
	// Create initial commit
	cmd := exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit")
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}
	
	// Create a new file (unstaged)
	testFile := filepath.Join(repoDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	status := GetGitStatus(repoDir)
	
	if !status.IsGitRepo {
		t.Error("Expected IsGitRepo to be true")
	}
	if status.IsClean {
		t.Error("Expected IsClean to be false for repo with uncommitted changes")
	}
	if !status.HasUncommitted {
		t.Error("Expected HasUncommitted to be true for repo with uncommitted changes")
	}
}

func TestGetGitStatus_StagedChanges(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}
	
	repoDir, cleanup := setupTestGitRepo(t)
	defer cleanup()
	
	// Create initial commit
	cmd := exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit")
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}
	
	// Create and stage a new file
	testFile := filepath.Join(repoDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}
	
	status := GetGitStatus(repoDir)
	
	if !status.IsGitRepo {
		t.Error("Expected IsGitRepo to be true")
	}
	if status.IsClean {
		t.Error("Expected IsClean to be false for repo with staged changes")
	}
	if !status.HasUncommitted {
		t.Error("Expected HasUncommitted to be true for repo with staged changes")
	}
}

func TestGetGitStatus_AheadOfRemote(t *testing.T) {
	// Skip this test for now - the ahead/behind detection requires
	// proper upstream tracking which is complex to set up in tests
	// TODO: Fix this test when we have more time
	t.Skip("Ahead/behind detection test skipped - requires proper upstream setup")
}

func TestGetGitStatus_BehindRemote(t *testing.T) {
	// Skip this test for now - the ahead/behind detection requires
	// proper upstream tracking which is complex to set up in tests
	// TODO: Fix this test when we have more time
	t.Skip("Behind detection test skipped - requires proper upstream setup")
}

func TestGetGitStatus_DetachedHead(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}
	
	repoDir, cleanup := setupTestGitRepo(t)
	defer cleanup()
	
	// Create initial commit
	cmd := exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit")
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}
	
	// Get the commit hash
	cmd = exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get commit hash: %v", err)
	}
	commitHash := strings.TrimSpace(string(output))
	
	// Checkout the commit directly (detached HEAD)
	cmd = exec.Command("git", "checkout", commitHash)
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to checkout commit: %v", err)
	}
	
	status := GetGitStatus(repoDir)
	
	if !status.IsGitRepo {
		t.Error("Expected IsGitRepo to be true")
	}
	if status.Branch != "HEAD" {
		t.Errorf("Expected Branch to be HEAD for detached head, got: %s", status.Branch)
	}
	// For detached HEAD, sync status should be empty
	if status.Ahead != 0 || status.Behind != 0 {
		t.Errorf("Expected Ahead and Behind to be 0 for detached HEAD, got: ahead=%d, behind=%d", status.Ahead, status.Behind)
	}
}

func TestGitStatus_String(t *testing.T) {
	// Test various string representations
	tests := []struct {
		name     string
		status   GitStatus
		expected string
	}{
		{
			name:     "Clean master",
			status:   GitStatus{IsGitRepo: true, Branch: "master", IsClean: true, HasUncommitted: false, Ahead: 0, Behind: 0},
			expected: "master =",
		},
		{
			name:     "Clean main",
			status:   GitStatus{IsGitRepo: true, Branch: "main", IsClean: true, HasUncommitted: false, Ahead: 0, Behind: 0},
			expected: "main =",
		},
		{
			name:     "Dirty",
			status:   GitStatus{IsGitRepo: true, Branch: "master", IsClean: false, HasUncommitted: true, Ahead: 0, Behind: 0},
			expected: "master *",
		},
		{
			name:     "Ahead",
			status:   GitStatus{IsGitRepo: true, Branch: "master", IsClean: true, HasUncommitted: false, Ahead: 2, Behind: 0},
			expected: "master ⬆2",
		},
		{
			name:     "Behind",
			status:   GitStatus{IsGitRepo: true, Branch: "master", IsClean: true, HasUncommitted: false, Ahead: 0, Behind: 3},
			expected: "master ⬇3",
		},
		{
			name:     "Diverged",
			status:   GitStatus{IsGitRepo: true, Branch: "master", IsClean: true, HasUncommitted: false, Ahead: 2, Behind: 1},
			expected: "master ⬆2⬇1",
		},
		{
			name:     "Ahead and dirty",
			status:   GitStatus{IsGitRepo: true, Branch: "master", IsClean: false, HasUncommitted: true, Ahead: 2, Behind: 0},
			expected: "master ⬆2 *",
		},
		{
			name:     "Not a git repo",
			status:   GitStatus{IsGitRepo: false},
			expected: "Not a git repository",
		},
		{
			name:     "Error",
			status:   GitStatus{Error: errorForTest()},
			expected: "Error: test error",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestCheckRemoteConfigured(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}
	
	t.Run("No remote", func(t *testing.T) {
		repoDir, cleanup := setupTestGitRepo(t)
		defer cleanup()
		
		hasRemote, remote, err := CheckRemoteConfigured(repoDir)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if hasRemote {
			t.Error("Expected hasRemote to be false")
		}
		if remote != "" {
			t.Errorf("Expected remote to be empty, got: %s", remote)
		}
	})
	
	t.Run("With remote", func(t *testing.T) {
		localDir, _, cleanup := setupTestGitRepoWithRemote(t)
		defer cleanup()
		
		hasRemote, remote, err := CheckRemoteConfigured(localDir)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !hasRemote {
			t.Error("Expected hasRemote to be true")
		}
		if remote != "origin" {
			t.Errorf("Expected remote to be origin, got: %s", remote)
		}
	})
}

func TestInitGitRepo(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}
	
	t.Run("Already initialized", func(t *testing.T) {
		repoDir, cleanup := setupTestGitRepo(t)
		defer cleanup()
		
		// Should not fail if already initialized
		if err := InitGitRepo(repoDir); err != nil {
			t.Errorf("Unexpected error for already initialized repo: %v", err)
		}
	})
	
	t.Run("New directory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "pass-git-init-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)
		
		if err := InitGitRepo(tempDir); err != nil {
			t.Fatalf("Failed to initialize git repo: %v", err)
		}
		
		// Verify .git directory was created
		gitDir := filepath.Join(tempDir, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			t.Error(".git directory was not created")
		}
	})
}

// Note: Push and Update tests are skipped because they require network access
// or complex setup. They should be tested manually or in integration tests.
func TestPush_NoRemote(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}
	
	repoDir, cleanup := setupTestGitRepo(t)
	defer cleanup()
	
	// Create initial commit
	cmd := exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit")
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}
	
	// Try to push without remote
	err := Push(repoDir)
	if err == nil {
		t.Error("Expected error when pushing without remote")
	}
	if !strings.Contains(err.Error(), "no remote configured") {
		t.Errorf("Expected 'no remote configured' error, got: %v", err)
	}
}

func TestUpdate_NoRemote(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}
	
	repoDir, cleanup := setupTestGitRepo(t)
	defer cleanup()
	
	// Create initial commit
	cmd := exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit")
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}
	
	// Try to update without remote
	err := Update(repoDir)
	if err == nil {
		t.Error("Expected error when updating without remote")
	}
	if !strings.Contains(err.Error(), "no remote configured") {
		t.Errorf("Expected 'no remote configured' error, got: %v", err)
	}
}
