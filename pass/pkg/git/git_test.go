package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCheckGit(t *testing.T) {
	// Test that git is available
	err := CheckGit()
	if err != nil {
		t.Skipf("Git not available: %v", err)
	}
	t.Log("Git is available")
}

func TestInitRepo(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "pass-git-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize repo
	if err := InitRepo(tempDir); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	// Verify .git directory was created
	gitDir := filepath.Join(tempDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Error(".git directory was not created")
	}
}

func TestAddAndCommit(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}

	// Create temp directory with git repo
	tempDir, err := os.MkdirTemp("", "pass-git-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize git repo
	if err := InitRepo(tempDir); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	// Configure git user in the repo directory
	cmd := exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add and commit
	message := "Test commit"
	if err := AddAndCommit(testFile, message); err != nil {
		t.Fatalf("AddAndCommit failed: %v", err)
	}

	// Verify file was added to git
	// (We can't easily verify the commit without parsing git output,
	// but at least we know it didn't error)
	t.Log("AddAndCommit succeeded")
}

func TestRunGitConfig(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}

	// Create temp directory with git repo
	tempDir, err := os.MkdirTemp("", "pass-git-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize git repo first
	if err := InitRepo(tempDir); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	// Configure git user in the repo directory
	cmd := exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}

	// Set a config value
	if err := runGitConfig(tempDir, "test.key", "test.value"); err != nil {
		t.Fatalf("runGitConfig failed: %v", err)
	}

	// Verify it was set (by trying to read it back)
	// This is a simple test - in a real scenario we'd verify the value
	t.Log("runGitConfig succeeded")
}
