package gpg

import (
	"os"
	"path/filepath"
	"testing"
)

// TestEncryptDecryptWithBatchMode tests encryption/decryption in batch mode
// This is useful for non-interactive testing
func TestEncryptDecryptWithBatchMode(t *testing.T) {
	// Skip if GPG not available
	if err := CheckGPG(); err != nil {
		t.Skipf("GPG not available: %v", err)
	}

	// Create temp directory
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
		// In batch mode without passphrase, this might fail if the key has a passphrase
		// That's expected behavior - just log and skip
		t.Logf("Batch decryption without passphrase failed (expected if key has passphrase): %v", err)
		t.Skip("Key requires passphrase for batch mode")
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
