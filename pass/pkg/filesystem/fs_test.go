package filesystem

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "forward slashes",
			input:    "email/gmail.com/user",
			expected: "email/gmail.com/user",
		},
		{
			name:     "backslashes",
			input:    "email\\gmail.com\\user",
			expected: "email/gmail.com/user",
		},
		{
			name:     "mixed slashes",
			input:    "email/gmail.com\\user",
			expected: "email/gmail.com/user",
		},
		{
			name:     "with dots",
			input:    "email/./gmail.com/../gmail.com/user",
			expected: "email/gmail.com/user",
		},
		{
			name:     "empty",
			input:    "",
			expected: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePath(tt.input)
			// Normalize expected to use OS separator
			expected := strings.ReplaceAll(tt.expected, "/", string(filepath.Separator))
			expected = strings.ReplaceAll(expected, "\\", string(filepath.Separator))
			// Clean the path to match filepath.Clean behavior
			expected = filepath.Clean(expected)
			if result != expected {
				t.Errorf("NormalizePath(%q) = %q, want %q", tt.input, result, expected)
			}
		})
	}
}

func TestNormalizePathForDisplay(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "backslashes to forward slashes",
			input:    "email\\gmail.com\\user",
			expected: "email/gmail.com/user",
		},
		{
			name:     "already forward slashes",
			input:    "email/gmail.com/user",
			expected: "email/gmail.com/user",
		},
		{
			name:     "mixed slashes",
			input:    "email\\gmail.com/user",
			expected: "email/gmail.com/user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePathForDisplay(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePathForDisplay(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSecureDelete(t *testing.T) {
	// Create a temp file
	tempFile, err := os.CreateTemp("", "pass-test-*.tmp")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		// Clean up if SecureDelete fails
		if _, err := os.Stat(tempFile.Name()); err == nil {
			os.Remove(tempFile.Name())
		}
	}()

	// Write some content
	content := []byte("test password content")
	if _, err := tempFile.Write(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Securely delete the file
	if err := SecureDelete(tempFile.Name()); err != nil {
		t.Fatalf("SecureDelete failed: %v", err)
	}

	// Verify file is gone
	if _, err := os.Stat(tempFile.Name()); !os.IsNotExist(err) {
		t.Errorf("File still exists after SecureDelete")
	}
}

func TestCopyToClipboard(t *testing.T) {
	// Skip on systems without clipboard
	if !IsClipboardAvailable() {
		t.Skip("Clipboard not available")
	}

	// Test copying simple text
	testText := "test password 123"
	if err := CopyToClipboard(testText); err != nil {
		t.Fatalf("CopyToClipboard failed: %v", err)
	}

	// Try to read back (may not work on all systems)
	// For now, just verify no error
	t.Log("Clipboard copy test passed (visual verification needed)")
}

func TestEnsurePasswordStore(t *testing.T) {
	// Create a temp directory for testing
	tempDir, err := os.MkdirTemp("", "pass-test-store")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test creating store
	storeDir := filepath.Join(tempDir, ".password-store")
	if err := EnsurePasswordStore(storeDir); err != nil {
		t.Fatalf("EnsurePasswordStore failed: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		t.Errorf("Password store directory was not created")
	}
}

// IsClipboardAvailable is a helper for tests
func IsClipboardAvailable() bool {
	// Simple check - try to run clip command
	cmd := exec.Command("clip")
	return cmd.Run() == nil
}
