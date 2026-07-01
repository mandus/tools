package gpg

import (
	"fmt"
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
		{
			name:     "operation cancelled",
			input:    "gpg: decryption failed: Operation cancelled",
			expected: "gpg: decryption failed: Operation cancelled",
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

func TestShouldRetryDecryption(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		stderr   string
		expected bool
	}{
		{
			name:     "operation cancelled",
			err:      fmt.Errorf("gpg: decryption failed: Operation cancelled"),
			stderr:   "",
			expected: true,
		},
		{
			name:     "gpg cancelled",
			err:      fmt.Errorf("gpg: cancelled"),
			stderr:   "",
			expected: true,
		},
		{
			name:     "no pinentry",
			err:      fmt.Errorf("gpg-agent: no pinentry"),
			stderr:   "",
			expected: true,
		},
		{
			name:     "no secret key",
			err:      fmt.Errorf("gpg: No secret key"),
			stderr:   "",
			expected: false,
		},
		{
			name:     "bad passphrase",
			err:      fmt.Errorf("gpg: Bad passphrase"),
			stderr:   "",
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			stderr:   "",
			expected: false,
		},
		{
			name:     "other error",
			err:      fmt.Errorf("some other error"),
			stderr:   "",
			expected: false,
		},
		{
			name:     "exit status 2",
			err:      fmt.Errorf("pass: GPG decryption failed: exit status 2"),
			stderr:   "",
			expected: true,
		},
		{
			name:     "tty error in stderr",
			err:      fmt.Errorf("pass: GPG decryption failed: exit status 2"),
			stderr:   "gpg: cannot open '/dev/tty': No such device or address",
			expected: true,
		},
		{
			name:     "operation cancelled in message",
			err:      fmt.Errorf("pass: decryption failed: gpg: decryption failed: Operation cancelled"),
			stderr:   "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldRetryDecryption(tt.err, tt.stderr)
			if result != tt.expected {
				t.Errorf("shouldRetryDecryption(%v, %q) = %v, want %v", tt.err, tt.stderr, result, tt.expected)
			}
		})
	}
}

func TestCheckGPGAgent(t *testing.T) {
	// Test with current environment
	isRunning := CheckGPGAgent()
	t.Logf("gpg-agent is running: %v", isRunning)
	
	// If running, we can test the positive case
	// If not running, we can only test the negative case
	if !isRunning {
		// This is expected in some test environments
		t.Log("gpg-agent is not running in test environment (expected)")
	}
}

func TestDefaultGPGOptions(t *testing.T) {
	opts := DefaultGPGOptions()
	
	if opts.BatchMode {
		t.Error("BatchMode should be false by default")
	}
	if opts.PinentryMode != "" {
		t.Errorf("PinentryMode should be empty by default, got %q", opts.PinentryMode)
	}
	if !opts.AllowPrompt {
		t.Error("AllowPrompt should be true by default")
	}
	if !opts.RetryOnCancel {
		t.Error("RetryOnCancel should be true by default")
	}
}

func TestBatchGPGOptions(t *testing.T) {
	opts := BatchGPGOptions("test-passphrase")
	
	if !opts.BatchMode {
		t.Error("BatchMode should be true for batch options")
	}
	if opts.Passphrase != "test-passphrase" {
		t.Errorf("Passphrase should be 'test-passphrase', got %q", opts.Passphrase)
	}
	if opts.PinentryMode != "loopback" {
		t.Errorf("PinentryMode should be 'loopback' for batch options, got %q", opts.PinentryMode)
	}
	if opts.AllowPrompt {
		t.Error("AllowPrompt should be false for batch options")
	}
	if opts.RetryOnCancel {
		t.Error("RetryOnCancel should be false for batch options")
	}
}
