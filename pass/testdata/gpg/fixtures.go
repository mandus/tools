// Package gpgtest provides test fixtures for GPG operations.
// This package contains test-only GPG keys and utilities for testing
// the pass password manager.
//
// ⚠️ IMPORTANT: These are test-only keys. NEVER use them for real password storage.
package gpgtest

import (
	"os"
	"path/filepath"
	"runtime"
)

// TestFixturesPath is the path to the GPG test fixtures directory.
// This is relative to the pass package directory.
const TestFixturesPath = "testdata/gpg"

// TestGNUPGHome is the path to the test GPG home directory.
// This contains test-only GPG keys.
const TestGNUPGHome = "testdata/gpg/test-gnupg-home"

// NoPassphraseKeyID is the GPG key ID for the test key without a passphrase.
// This key is used for basic encryption/decryption tests.
const NoPassphraseKeyID = "TEST_NO_PASSPHRASE_KEY"

// WithPassphraseKeyID is the GPG key ID for the test key with a passphrase.
// This key is used for testing passphrase-protected encryption.
const WithPassphraseKeyID = "TEST_WITH_PASSPHRASE_KEY"

// TestPassphrase is the known passphrase for the test key with passphrase.
// This is safe to include in tests because it's only used with test keys.
const TestPassphrase = "test-passphrase-123"

// TestPasswordStore is the path to the test password store directory.
const TestPasswordStore = "testdata/store"

// SetupTestEnvironment sets up the environment for testing.
// It configures GNUPGHOME to point to the test GPG home directory
// and PASSWORD_STORE_DIR to point to the test password store.
//
// Returns a function that should be deferred to restore the original environment.
func SetupTestEnvironment() func() {
	// Save original environment
	origGNUPGHOME := os.Getenv("GNUPGHOME")
	origPasswordStore := os.Getenv("PASSWORD_STORE_DIR")
	origPath := os.Getenv("PATH")

	// Get the directory of this file
	_, filename, _, _ := runtime.Caller(0)
	fixturesDir := filepath.Dir(filename)
	passDir := filepath.Dir(fixturesDir)

	// Set test environment
	os.Setenv("GNUPGHOME", filepath.Join(passDir, TestGNUPGHome))
	os.Setenv("PASSWORD_STORE_DIR", filepath.Join(passDir, TestPasswordStore))

	// Add the pass directory to PATH so we can run the pass binary if needed
	// This is useful for integration tests
	if origPath != "" {
		os.Setenv("PATH", filepath.Join(passDir)+string(os.PathListSeparator)+origPath)
	} else {
		os.Setenv("PATH", filepath.Join(passDir))
	}

	// Return cleanup function
	return func() {
		if origGNUPGHOME != "" {
			os.Setenv("GNUPGHOME", origGNUPGHOME)
		} else {
			os.Unsetenv("GNUPGHOME")
		}
		if origPasswordStore != "" {
			os.Setenv("PASSWORD_STORE_DIR", origPasswordStore)
		} else {
			os.Unsetenv("PASSWORD_STORE_DIR")
		}
		os.Setenv("PATH", origPath)
	}
}

// GetTestGNUPGHome returns the absolute path to the test GPG home directory.
func GetTestGNUPGHome() string {
	_, filename, _, _ := runtime.Caller(0)
	fixturesDir := filepath.Dir(filename)
	passDir := filepath.Dir(fixturesDir)
	return filepath.Join(passDir, TestGNUPGHome)
}

// GetTestPasswordStore returns the absolute path to the test password store.
func GetTestPasswordStore() string {
	_, filename, _, _ := runtime.Caller(0)
	fixturesDir := filepath.Dir(filename)
	passDir := filepath.Dir(fixturesDir)
	return filepath.Join(passDir, TestPasswordStore)
}

// CreateTestPasswordStore creates a clean test password store directory.
// Returns the path to the store and a cleanup function.
func CreateTestPasswordStore() (string, func()) {
	storePath := GetTestPasswordStore()
	
	// Remove existing test store
	os.RemoveAll(storePath)
	
	// Create the store directory
	if err := os.MkdirAll(storePath, 0700); err != nil {
		panic(err)
	}
	
	// Create .password-store subdirectory
	passwordStore := filepath.Join(storePath, ".password-store")
	if err := os.MkdirAll(passwordStore, 0700); err != nil {
		panic(err)
	}
	
	// Return cleanup function
	return storePath, func() {
		os.RemoveAll(storePath)
	}
}
