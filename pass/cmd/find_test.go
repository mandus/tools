package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindPasswords(t *testing.T) {
	// Create a temp password store
	tempDir, err := os.MkdirTemp("", "pass-find-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set PASSWORD_STORE_DIR to temp dir
	orig := os.Getenv("PASSWORD_STORE_DIR")
	os.Setenv("PASSWORD_STORE_DIR", tempDir)
	defer os.Setenv("PASSWORD_STORE_DIR", orig)

	// Create some test password files
	passwords := []string{
		"email/gmail.com/user1",
		"email/gmail.com/user2",
		"bank/chase/account",
		"social/twitter.com/user",
	}

	for _, p := range passwords {
		fullPath := filepath.Join(tempDir, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath+".gpg", []byte("encrypted"), 0600); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Test finding "gmail" - call function directly
	// Reset ignoreCaseFlag
	ignoreCaseFlag = false
	
	err = findPasswords("gmail")
	if err != nil {
		t.Fatalf("findPasswords failed: %v", err)
	}
	t.Log("Find passwords works")
}

func TestFindPasswordsCaseInsensitive(t *testing.T) {
	// Create a temp password store
	tempDir, err := os.MkdirTemp("", "pass-find-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set PASSWORD_STORE_DIR to temp dir
	orig := os.Getenv("PASSWORD_STORE_DIR")
	os.Setenv("PASSWORD_STORE_DIR", tempDir)
	defer os.Setenv("PASSWORD_STORE_DIR", orig)

	// Create some test password files
	passwords := []string{
		"email/GMAIL.com/user1",
		"email/gmail.com/user2",
	}

	for _, p := range passwords {
		fullPath := filepath.Join(tempDir, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath+".gpg", []byte("encrypted"), 0600); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Test finding "gmail" with case-insensitive flag
	ignoreCaseFlag = true
	
	err = findPasswords("gmail")
	if err != nil {
		t.Fatalf("findPasswords failed: %v", err)
	}
	t.Log("Find passwords case-insensitive works")
}

func TestFindPasswordsEmptyString(t *testing.T) {
	// Test that empty search string returns error
	ignoreCaseFlag = false
	
	err := findPasswords("")
	if err == nil {
		t.Error("Expected error for empty search string")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("Expected 'empty' in error message, got: %v", err)
	}
	t.Log("Empty search string error handling works")
}
