package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListPasswords(t *testing.T) {
	// Create a temp password store
	tempDir, err := os.MkdirTemp("", "pass-ls-test")
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
	}

	for _, p := range passwords {
		// Create directory structure
		fullPath := filepath.Join(tempDir, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		
		// Create .gpg file
		if err := os.WriteFile(fullPath+".gpg", []byte("encrypted"), 0600); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Test listing all - call the function directly
	err = listPasswords("")
	if err != nil {
		t.Fatalf("listPasswords failed: %v", err)
	}
	t.Log("List all passwords works")
}

func TestListPasswordsSubpath(t *testing.T) {
	// Create a temp password store
	tempDir, err := os.MkdirTemp("", "pass-ls-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set PASSWORD_STORE_DIR to temp dir
	orig := os.Getenv("PASSWORD_STORE_DIR")
	os.Setenv("PASSWORD_STORE_DIR", tempDir)
	defer os.Setenv("PASSWORD_STORE_DIR", orig)

	// Create test password files
	passwords := []string{
		"email/gmail.com/user1",
		"email/gmail.com/user2",
		"bank/chase/account",
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

	// Test listing subpath - call the function directly
	err = listPasswords("email/")
	if err != nil {
		t.Fatalf("listPasswords failed: %v", err)
	}
	t.Log("List subpath works")
}

func TestListPasswordsNonexistentSubpath(t *testing.T) {
	// Create a temp password store
	tempDir, err := os.MkdirTemp("", "pass-ls-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set PASSWORD_STORE_DIR to temp dir
	orig := os.Getenv("PASSWORD_STORE_DIR")
	os.Setenv("PASSWORD_STORE_DIR", tempDir)
	defer os.Setenv("PASSWORD_STORE_DIR", orig)

	// Test listing nonexistent subpath
	err = listPasswords("nonexistent/")
	if err == nil {
		t.Error("Expected error for nonexistent subpath")
	}
	if !strings.Contains(err.Error(), "No such file or directory") {
		t.Errorf("Expected 'No such file or directory' error, got: %v", err)
	}
	t.Log("Nonexistent subpath error handling works")
}
