package gpg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mandu/tools/pass/internal/testhelper"
)

func TestCheckGPG(t *testing.T) {
	// Test that GPG is available
	err := CheckGPG()
	if err != nil {
		t.Skipf("GPG not available: %v", err)
	}
	// If we get here, GPG is available
	t.Log("GPG is available")
}

func TestHasSecretKey(t *testing.T) {
	// Set up test environment with GPG keys
	env, noPassKey, _, err := testhelper.SetupTestEnvWithGPGKeys()
	if err != nil {
		t.Skipf("Failed to set up test environment: %v", err)
	}
	defer env.Cleanup()

	// Set GNUPGHOME for this test
	os.Setenv("GNUPGHOME", env.GNUPGHome)
	os.Setenv("PASS_GPG_ID", noPassKey)

	hasKey := HasSecretKey()
	if !hasKey {
		t.Error("Expected to have secret key in test environment")
	}
}

func TestGetDefaultRecipient(t *testing.T) {
	// Set up test environment with GPG keys
	env, noPassKey, _, err := testhelper.SetupTestEnvWithGPGKeys()
	if err != nil {
		t.Skipf("Failed to set up test environment: %v", err)
	}
	defer env.Cleanup()

	// Set GNUPGHOME for this test
	os.Setenv("GNUPGHOME", env.GNUPGHome)
	os.Setenv("PASS_GPG_ID", noPassKey)

	recipient := GetDefaultRecipient()
	if recipient == "" {
		t.Error("Expected to get a default recipient in test environment")
	}
	t.Logf("Default GPG recipient: %s", recipient)
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	// Set up test environment with GPG keys
	env, noPassKey, _, err := testhelper.SetupTestEnvWithGPGKeys()
	if err != nil {
		t.Skipf("Failed to set up test environment: %v", err)
	}
	defer env.Cleanup()

	// Set GNUPGHOME for this test
	os.Setenv("GNUPGHOME", env.GNUPGHome)
	os.Setenv("PASS_GPG_ID", noPassKey)

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "pass-gpg-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "test password 123"
	if err := os.WriteFile(testFile, []byte(testContent), 0600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Encrypt the file with batch mode
	encryptedFile := filepath.Join(tempDir, "test.txt.gpg")
	opts := BatchGPGOptions("") // Batch mode without passphrase
	if err := EncryptFileWithOptions(testFile, encryptedFile, opts); err != nil {
		t.Fatalf("Failed to encrypt file: %v", err)
	}
	defer os.Remove(encryptedFile)

	// Verify encrypted file exists
	if _, err := os.Stat(encryptedFile); os.IsNotExist(err) {
		t.Fatal("Encrypted file was not created")
	}

	// Decrypt the file with batch mode
	decrypted, err := DecryptFileWithOptions(encryptedFile, opts)
	if err != nil {
		t.Fatalf("Failed to decrypt file: %v", err)
	}

	// Verify content matches
	if decrypted != testContent {
		t.Errorf("Decrypted content = %q, want %q", decrypted, testContent)
	}
}

func TestExtractGPGError(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "decryption failed",
			input:    "gpg: decryption failed: No secret key\ngpg: some other line",
			expected: "gpg: decryption failed: No secret key",
		},
		{
			name:     "no secret key",
			input:    "gpg: No secret key",
			expected: "gpg: No secret key",
		},
		{
			name:     "empty",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractGPGError(tt.input)
			if result != tt.expected {
				t.Errorf("extractGPGError(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
