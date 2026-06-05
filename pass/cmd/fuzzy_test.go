package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCollectAllPasswordsFromStore(t *testing.T) {
	// Create a temp password store
	tempDir, err := os.MkdirTemp("", "pass-fuzzy-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set PASSWORD_STORE_DIR to temp dir
	orig := os.Getenv("PASSWORD_STORE_DIR")
	os.Setenv("PASSWORD_STORE_DIR", tempDir)
	defer os.Setenv("PASSWORD_STORE_DIR", orig)

	// Create some test password files
	passwords := []string{"email/gmail.com/user", "social/twitter.com/admin", "bank/chase.com/account"}
	for _, p := range passwords {
		fullPath := filepath.Join(tempDir, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath+".gpg", []byte("dummy"), 0600); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Collect passwords using the store directory
	got, err := collectAllPasswords(tempDir)
	if err != nil {
		t.Fatalf("collectAllPasswords failed: %v", err)
	}

	// Check that we got all passwords
	if len(got) != len(passwords) {
		t.Errorf("Expected %d passwords, got %d", len(passwords), len(got))
	}

	// Check that all expected passwords are present
	for _, want := range passwords {
		found := false
		for _, got := range got {
			if got == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Missing expected password: %s", want)
		}
	}
}

func TestUseFuzzyPath(t *testing.T) {
	// Create a temp password store
	tempDir, err := os.MkdirTemp("", "pass-fuzzy-path-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set PASSWORD_STORE_DIR to temp dir
	orig := os.Getenv("PASSWORD_STORE_DIR")
	os.Setenv("PASSWORD_STORE_DIR", tempDir)
	defer os.Setenv("PASSWORD_STORE_DIR", orig)

	// Create some test password files
	passwords := []string{"email/gmail.com/user", "social/twitter.com/admin"}
	for _, p := range passwords {
		fullPath := filepath.Join(tempDir, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath+".gpg", []byte("dummy"), 0600); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Test finding best match
	tests := []struct {
		query string
		want  string
	}{
		{"gm", "email/gmail.com/user"},
		{"tw", "social/twitter.com/admin"},
		{"email", "email/gmail.com/user"},
		{"nonexistent", ""},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			got, err := UseFuzzyPath(tt.query)
			if err != nil && tt.want != "" {
				t.Errorf("UseFuzzyPath(%q) error = %v, want %q", tt.query, err, tt.want)
				return
			}
			if got != tt.want {
				t.Errorf("UseFuzzyPath(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
}

func TestFuzzySearchModeEmptyStore(t *testing.T) {
	// Create a temp password store
	tempDir, err := os.MkdirTemp("", "pass-fuzzy-empty-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set PASSWORD_STORE_DIR to temp dir
	orig := os.Getenv("PASSWORD_STORE_DIR")
	os.Setenv("PASSWORD_STORE_DIR", tempDir)
	defer os.Setenv("PASSWORD_STORE_DIR", orig)

	// Don't create any password files - store is empty

	// fuzzySearchMode should return an error
	err = fuzzySearchMode()
	if err == nil {
		t.Error("Expected error for empty store")
	}
	if !strings.Contains(err.Error(), "no passwords found") {
		t.Errorf("Expected 'no passwords found' error, got: %v", err)
	}
}

func TestShowSelectedPassword(t *testing.T) {
	// This test verifies that showSelectedPassword correctly delegates to showPassword
	// We can't easily test the actual showPassword without a real GPG setup
	// So we just verify it doesn't panic with valid input
	
	// Set PASSWORD_STORE_DIR to a non-existent location
	// This will cause showPassword to fail, but we can verify the function works
	orig := os.Getenv("PASSWORD_STORE_DIR")
	os.Setenv("PASSWORD_STORE_DIR", "C:\\nonexistent\\path")
	defer os.Setenv("PASSWORD_STORE_DIR", orig)

	// This should return an error (file not found)
	err := showSelectedPassword("test/password")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}
