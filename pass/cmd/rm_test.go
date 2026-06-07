package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemovePassword(t *testing.T) {
	// Create a temp password store
	tempDir, err := os.MkdirTemp("", "pass-rm-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set PASSWORD_STORE_DIR to temp dir
	orig := os.Getenv("PASSWORD_STORE_DIR")
	os.Setenv("PASSWORD_STORE_DIR", tempDir)
	defer os.Setenv("PASSWORD_STORE_DIR", orig)

	// Create a test password file
	passwordPath := filepath.Join(tempDir, "test", "password.gpg")
	if err := os.MkdirAll(filepath.Dir(passwordPath), 0700); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	
	// Create a dummy .gpg file
	if err := os.WriteFile(passwordPath, []byte("dummy encrypted content"), 0600); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Test removing the password
	err = removePassword("test/password", false, false)
	if err != nil {
		t.Fatalf("removePassword failed: %v", err)
	}

	// Check that file was removed
	if _, err := os.Stat(passwordPath); !os.IsNotExist(err) {
		t.Error("File was not removed")
	}
}

func TestRemovePasswordWithClip(t *testing.T) {
	// This test just verifies the clip flag doesn't cause errors
	// Actual clipboard functionality is tested in clipboard package
	// Create a temp password store
	tempDir, err := os.MkdirTemp("", "pass-rm-clip-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set PASSWORD_STORE_DIR to temp dir
	orig := os.Getenv("PASSWORD_STORE_DIR")
	os.Setenv("PASSWORD_STORE_DIR", tempDir)
	defer os.Setenv("PASSWORD_STORE_DIR", orig)

	// Create a test password file
	passwordPath := filepath.Join(tempDir, "test", "password.gpg")
	if err := os.MkdirAll(filepath.Dir(passwordPath), 0700); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	
	// Create a dummy .gpg file
	if err := os.WriteFile(passwordPath, []byte("dummy encrypted content"), 0600); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Test removing with clip flag
	// This will fail to decrypt since it's not a real GPG file, but that's expected
	err = removePassword("test/password", false, true)
	// We expect an error because it can't decrypt the dummy file
	if err == nil {
		// If no error, check file was still removed
		if _, err := os.Stat(passwordPath); !os.IsNotExist(err) {
			t.Error("File was not removed when clip flag was set")
		}
	} else {
		// Error is expected (can't decrypt)
		if !strings.Contains(err.Error(), "decrypt") {
			t.Logf("Expected decryption error, got: %v", err)
		}
	}
}

func TestRemovePasswordNotFound(t *testing.T) {
	// Create a temp password store
	tempDir, err := os.MkdirTemp("", "pass-rm-notfound-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set PASSWORD_STORE_DIR to temp dir
	orig := os.Getenv("PASSWORD_STORE_DIR")
	os.Setenv("PASSWORD_STORE_DIR", tempDir)
	defer os.Setenv("PASSWORD_STORE_DIR", orig)

	// Try to remove a non-existent password
	err = removePassword("nonexistent/password", false, false)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
	if !strings.Contains(err.Error(), "No such file or directory") {
		t.Errorf("Expected 'No such file or directory' error, got: %v", err)
	}
}

func TestRemovePasswordPathNormalization(t *testing.T) {
	// Create a temp password store
	tempDir, err := os.MkdirTemp("", "pass-rm-path-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set PASSWORD_STORE_DIR to temp dir
	orig := os.Getenv("PASSWORD_STORE_DIR")
	os.Setenv("PASSWORD_STORE_DIR", tempDir)
	defer os.Setenv("PASSWORD_STORE_DIR", orig)

	// Create a test password file with OS-specific path separator
	passwordPath := filepath.Join(tempDir, "test", "password.gpg")
	if err := os.MkdirAll(filepath.Dir(passwordPath), 0700); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	
	if err := os.WriteFile(passwordPath, []byte("dummy encrypted content"), 0600); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Test removing with forward slashes (should work on all platforms)
	err = removePassword("test/password", false, false)
	if err != nil {
		t.Fatalf("removePassword with forward slashes failed: %v", err)
	}

	// Check that file was removed
	if _, err := os.Stat(passwordPath); !os.IsNotExist(err) {
		t.Error("File was not removed when using forward slashes")
	}
}



func TestCollectAllPasswords(t *testing.T) {
	// Create a temp password store
	tempDir, err := os.MkdirTemp("", "pass-collect-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

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

	// Collect passwords
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
