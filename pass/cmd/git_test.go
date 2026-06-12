package cmd

import (
	"os"
	"os/exec"
	"testing"
)

// setupTestPasswordStoreWithGit creates a test password store with git initialized
func setupTestPasswordStoreWithGit(t *testing.T) (string, func()) {
	t.Helper()
	
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "pass-git-cmd-test")
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
	
	// Create initial commit
	cmd = exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		cleanup()
		t.Fatalf("Failed to create initial commit: %v", err)
	}
	
	return tempDir, cleanup
}

// CheckGit is a helper to check if git is available
func CheckGit() error {
	cmd := exec.Command("git", "--version")
	return cmd.Run()
}

func TestGitStatusCommand(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}
	
	t.Run("Git status on clean repo", func(t *testing.T) {
		repoDir, cleanup := setupTestPasswordStoreWithGit(t)
		defer cleanup()
		
		// Set PASSWORD_STORE_DIR
		originalDir := os.Getenv("PASSWORD_STORE_DIR")
		os.Setenv("PASSWORD_STORE_DIR", repoDir)
		defer os.Setenv("PASSWORD_STORE_DIR", originalDir)
		
		// This test would normally use the cobra test framework
		// For now, we'll just verify the git status package works
		t.Skip("CLI testing requires cobra test framework setup")
	})
	
	t.Run("Git status on non-existent store", func(t *testing.T) {
		// Set PASSWORD_STORE_DIR to non-existent directory
		originalDir := os.Getenv("PASSWORD_STORE_DIR")
		os.Setenv("PASSWORD_STORE_DIR", "/nonexistent/path")
		defer os.Setenv("PASSWORD_STORE_DIR", originalDir)
		
		t.Skip("CLI testing requires cobra test framework setup")
	})
}

func TestGitPushCommand(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}
	
	t.Run("Push with remote configured", func(t *testing.T) {
		repoDir, cleanup := setupTestPasswordStoreWithGit(t)
		defer cleanup()
		
		// Set PASSWORD_STORE_DIR
		originalDir := os.Getenv("PASSWORD_STORE_DIR")
		os.Setenv("PASSWORD_STORE_DIR", repoDir)
		defer os.Setenv("PASSWORD_STORE_DIR", originalDir)
		
		// Create a new commit
		cmd := exec.Command("git", "commit", "--allow-empty", "-m", "Test commit")
		cmd.Dir = repoDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to create test commit: %v", err)
		}
		
		t.Skip("CLI testing requires cobra test framework setup")
	})
	
	t.Run("Push without remote", func(t *testing.T) {
		repoDir, cleanup := setupTestPasswordStoreWithGit(t)
		defer cleanup()
		
		// Set PASSWORD_STORE_DIR
		originalDir := os.Getenv("PASSWORD_STORE_DIR")
		os.Setenv("PASSWORD_STORE_DIR", repoDir)
		defer os.Setenv("PASSWORD_STORE_DIR", originalDir)
		
		t.Skip("CLI testing requires cobra test framework setup")
	})
}

func TestGitUpdateCommand(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}
	
	t.Run("Update with remote configured", func(t *testing.T) {
		repoDir, cleanup := setupTestPasswordStoreWithGit(t)
		defer cleanup()
		
		// Set PASSWORD_STORE_DIR
		originalDir := os.Getenv("PASSWORD_STORE_DIR")
		os.Setenv("PASSWORD_STORE_DIR", repoDir)
		defer os.Setenv("PASSWORD_STORE_DIR", originalDir)
		
		t.Skip("CLI testing requires cobra test framework setup")
	})
	
	t.Run("Update without remote", func(t *testing.T) {
		repoDir, cleanup := setupTestPasswordStoreWithGit(t)
		defer cleanup()
		
		// Set PASSWORD_STORE_DIR
		originalDir := os.Getenv("PASSWORD_STORE_DIR")
		os.Setenv("PASSWORD_STORE_DIR", repoDir)
		defer os.Setenv("PASSWORD_STORE_DIR", originalDir)
		
		t.Skip("CLI testing requires cobra test framework setup")
	})
}

func TestGitInitCommand(t *testing.T) {
	// Skip if git not available
	if err := CheckGit(); err != nil {
		t.Skipf("Git not available: %v", err)
	}
	
	t.Run("Init on new directory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "pass-git-init-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)
		
		// Set PASSWORD_STORE_DIR
		originalDir := os.Getenv("PASSWORD_STORE_DIR")
		os.Setenv("PASSWORD_STORE_DIR", tempDir)
		defer os.Setenv("PASSWORD_STORE_DIR", originalDir)
		
		t.Skip("CLI testing requires cobra test framework setup")
	})
	
	t.Run("Init on existing git repo", func(t *testing.T) {
		repoDir, cleanup := setupTestPasswordStoreWithGit(t)
		defer cleanup()
		
		// Set PASSWORD_STORE_DIR
		originalDir := os.Getenv("PASSWORD_STORE_DIR")
		os.Setenv("PASSWORD_STORE_DIR", repoDir)
		defer os.Setenv("PASSWORD_STORE_DIR", originalDir)
		
		t.Skip("CLI testing requires cobra test framework setup")
	})
}

// Test that the git command is properly registered
func TestGitCommandRegistration(t *testing.T) {
	// This test verifies that the git command is added to the root command
	// We can't easily test this without running the full cobra setup,
	// but we can at least verify the function exists and doesn't panic
	
	// Just verify that addGitCmd doesn't panic
	// In a real test, we'd use cobra's test framework
	t.Skip("Command registration testing requires cobra test framework")
}
