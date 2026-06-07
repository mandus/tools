package gpg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mandu/tools/pass/internal/testhelper"
)

// TestEncryptDecryptWithBatchMode tests encryption/decryption in batch mode
// using ephemeral GPG keys generated on-the-fly
func TestEncryptDecryptWithBatchMode(t *testing.T) {
	// Set up test environment with ephemeral GPG keys
	env, noPassKey, _, err := testhelper.SetupTestEnvWithGPGKeys()
	if err != nil {
		t.Skipf("GPG not available for testing: %v", err)
	}
	defer env.Cleanup()

	// Set the recipient to use the test key without passphrase
	os.Setenv("PASS_GPG_ID", noPassKey)
	defer os.Unsetenv("PASS_GPG_ID")

	// Set GNUPGHOME to use the test environment
	os.Setenv("GNUPGHOME", env.GNUPGHome)
	defer os.Unsetenv("GNUPGHOME")

	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "pass-gpg-batch-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "test password batch 456"
	if err := os.WriteFile(testFile, []byte(testContent), 0600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Encrypt the file with batch mode
	encryptedFile := filepath.Join(tempDir, "test.txt.gpg")
	opts := BatchGPGOptions("") // Batch mode without passphrase
	if err := EncryptFileWithOptions(testFile, encryptedFile, opts); err != nil {
		t.Fatalf("Failed to encrypt file with batch mode: %v", err)
	}
	defer os.Remove(encryptedFile)

	// Verify encrypted file exists
	if _, err := os.Stat(encryptedFile); os.IsNotExist(err) {
		t.Fatal("Encrypted file was not created")
	}

	// Decrypt the file with batch mode
	decrypted, err := DecryptFileWithOptions(encryptedFile, opts)
	if err != nil {
		t.Fatalf("Failed to decrypt file with batch mode: %v", err)
	}

	// Verify content matches
	if decrypted != testContent {
		t.Errorf("Decrypted content = %q, want %q", decrypted, testContent)
	}
}

// TestEncryptDecryptWithPassphrase tests encryption/decryption with passphrase-protected keys
func TestEncryptDecryptWithPassphrase(t *testing.T) {
	// Set up test environment with ephemeral GPG keys
	env, _, withPassKey, err := testhelper.SetupTestEnvWithGPGKeys()
	if err != nil {
		t.Skipf("GPG not available for testing: %v", err)
	}
	defer env.Cleanup()

	// Set the recipient to use the test key with passphrase
	os.Setenv("PASS_GPG_ID", withPassKey)
	defer os.Unsetenv("PASS_GPG_ID")

	// Set GNUPGHOME to use the test environment
	os.Setenv("GNUPGHOME", env.GNUPGHome)
	defer os.Unsetenv("GNUPGHOME")

	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "pass-gpg-passphrase-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "test password with passphrase 789"
	if err := os.WriteFile(testFile, []byte(testContent), 0600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Encrypt the file with batch mode
	encryptedFile := filepath.Join(tempDir, "test.txt.gpg")
	// Use the known test passphrase
	opts := BatchGPGOptions("test-passphrase-123")
	if err := EncryptFileWithOptions(testFile, encryptedFile, opts); err != nil {
		t.Fatalf("Failed to encrypt file with passphrase: %v", err)
	}
	defer os.Remove(encryptedFile)

	// Verify encrypted file exists
	if _, err := os.Stat(encryptedFile); os.IsNotExist(err) {
		t.Fatal("Encrypted file was not created")
	}

	// Decrypt the file with batch mode and passphrase
	decrypted, err := DecryptFileWithOptions(encryptedFile, opts)
	if err != nil {
		t.Fatalf("Failed to decrypt file with passphrase: %v", err)
	}

	// Verify content matches
	if decrypted != testContent {
		t.Errorf("Decrypted content = %q, want %q", decrypted, testContent)
	}
}

// TestCheckGPGBatch tests batch mode availability
func TestCheckGPGBatch(t *testing.T) {
	// Skip if GPG not available
	if err := CheckGPG(); err != nil {
		t.Skipf("GPG not available: %v", err)
	}

	// Test batch mode
	err := CheckGPGBatch()
	if err != nil {
		t.Logf("GPG batch mode check: %v", err)
	} else {
		t.Log("GPG batch mode is available")
	}
}
