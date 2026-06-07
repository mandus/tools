// Package testhelper provides utilities for testing the pass package.
// This package contains helpers for setting up test environments,
// creating test fixtures, and managing test GPG keys.
//
// ⚠️ IMPORTANT: This package is for testing only. Never use it in production code.
// All test keys and data should be isolated from the user's personal setup.
package testhelper

import (
  "bytes"
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "runtime"
  "strings"
)

// TestEnv holds the test environment configuration
// It manages temporary directories and GPG settings for testing
type TestEnv struct {
  // Original environment values
  OrigGNUPGHOME     string
  OrigPasswordStore string
  OrigPath          string
  
  // Test paths
  TempDir       string
  GNUPGHome     string
  PasswordStore string
  
  // Cleanup functions
  cleanupFuncs []func()
}

// SetupTestEnv sets up a clean test environment for pass testing.
// It creates temporary directories and configures GPG to use test keys.
//
// Returns a TestEnv that should be cleaned up with Cleanup() when done.
//
// Example usage:
//
//  func TestSomething(t *testing.T) {
//      env, err := testhelper.SetupTestEnv()
//      if err != nil {
//          t.Fatal(err)
//      }
//      defer env.Cleanup()
//      
//      // Run tests here
//  }
func SetupTestEnv() (*TestEnv, error) {
  // Save original environment
  env := &TestEnv{
    OrigGNUPGHOME:     os.Getenv("GNUPGHOME"),
    OrigPasswordStore: os.Getenv("PASSWORD_STORE_DIR"),
    OrigPath:          os.Getenv("PATH"),
  }
  
  // Create temporary directory for test data
  tempDir, err := os.MkdirTemp("", "pass-test-env")
  if err != nil {
    return nil, err
  }
  env.TempDir = tempDir
  
  // Set up GPG home directory
  env.GNUPGHome = filepath.Join(tempDir, "gnupg-home")
  if err := os.MkdirAll(env.GNUPGHome, 0700); err != nil {
    os.RemoveAll(tempDir)
    return nil, err
  }
  // Normalize path for GPG compatibility (use forward slashes on Windows)
  env.GNUPGHome = filepath.ToSlash(env.GNUPGHome)
  // On Windows, convert drive letter to Unix-style path (e.g., C:/ -> /c/)
  if runtime.GOOS == "windows" && len(env.GNUPGHome) > 1 && env.GNUPGHome[1] == ':' {
    env.GNUPGHome = "/" + strings.ToLower(string(env.GNUPGHome[0])) + env.GNUPGHome[2:]
  }
  
  // Set up password store directory
  env.PasswordStore = filepath.Join(tempDir, "password-store")
  if err := os.MkdirAll(env.PasswordStore, 0700); err != nil {
    os.RemoveAll(tempDir)
    return nil, err
  }
  // Normalize path for consistency
  env.PasswordStore = filepath.ToSlash(env.PasswordStore)
  
  // Configure environment
  os.Setenv("GNUPGHOME", env.GNUPGHome)
  os.Setenv("PASSWORD_STORE_DIR", env.PasswordStore)
  
  // Get the pass directory to add to PATH
  _, filename, _, ok := runtime.Caller(1)
  if !ok {
    // Fallback: try to find pass directory
    wd, _ := os.Getwd()
    for {
      if strings.Contains(wd, "pass") || strings.HasSuffix(wd, "pass") {
        break
      }
      parent := filepath.Dir(wd)
      if parent == wd {
        break
      }
      wd = parent
    }
    filename = filepath.Join(wd, "testhelper.go")
  }
  passDir := filepath.Dir(filename)
  
  // Add pass directory to PATH
  newPath := passDir
  if env.OrigPath != "" {
    newPath = passDir + string(os.PathListSeparator) + env.OrigPath
  }
  os.Setenv("PATH", newPath)
  
  // Register cleanup
  env.cleanupFuncs = append(env.cleanupFuncs, func() {
    os.RemoveAll(tempDir)
  })
  
  return env, nil
}

// SetupTestEnvWithGPGKeys sets up a test environment with pre-generated GPG keys.
// This is useful for tests that need actual GPG encryption/decryption.
//
// The keys are generated on-the-fly if they don't exist in the test GPG home.
// Returns a TestEnv and the key IDs, or an error.
func SetupTestEnvWithGPGKeys() (*TestEnv, string, string, error) {
  env, err := SetupTestEnv()
  if err != nil {
    return nil, "", "", err
  }
  
  // Generate test GPG keys
  noPassphraseKeyID, withPassphraseKeyID, err := GenerateTestGPGKeys(env.GNUPGHome)
  if err != nil {
    env.Cleanup()
    return nil, "", "", err
  }
  
  return env, noPassphraseKeyID, withPassphraseKeyID, nil
}

// Cleanup restores the original environment and cleans up temporary files.
func (e *TestEnv) Cleanup() {
  // Restore original environment
  if e.OrigGNUPGHOME != "" {
    os.Setenv("GNUPGHOME", e.OrigGNUPGHOME)
  } else {
    os.Unsetenv("GNUPGHOME")
  }
  
  if e.OrigPasswordStore != "" {
    os.Setenv("PASSWORD_STORE_DIR", e.OrigPasswordStore)
  } else {
    os.Unsetenv("PASSWORD_STORE_DIR")
  }
  
  if e.OrigPath != "" {
    os.Setenv("PATH", e.OrigPath)
  } else {
    os.Unsetenv("PATH")
  }
  
  // Run all cleanup functions in reverse order
  for i := len(e.cleanupFuncs) - 1; i >= 0; i-- {
    e.cleanupFuncs[i]()
  }
}

// GenerateTestGPGKeys generates test GPG keys in the specified GPG home directory.
// Returns the key IDs for the keys without and with passphrase.
func GenerateTestGPGKeys(gpgHome string) (string, string, error) {
  // Check if GPG is available
  cmd := exec.Command("gpg", "--version")
  if err := cmd.Run(); err != nil {
    return "", "", err
  }
  
  // Set GNUPGHOME for the gpg commands
  cmd.Env = append(os.Environ(), "GNUPGHOME="+gpgHome)
  
  // Generate key without passphrase
  noPassphraseKeyID, err := generateTestKey(gpgHome, "Test User NoPass", "test-nopass@example.com", "", true)
  if err != nil {
    return "", "", err
  }
  
  // Generate key with passphrase
  withPassphraseKeyID, err := generateTestKey(gpgHome, "Test User WithPass", "test-withpass@example.com", "test-passphrase-123", false)
  if err != nil {
    return "", "", err
  }
  
  return noPassphraseKeyID, withPassphraseKeyID, nil
}

// generateTestKey generates a single test GPG key
func generateTestKey(gpgHome, name, email, passphrase string, noPassphrase bool) (string, error) {
  args := []string{
    "--batch",
    "--pinentry-mode", "loopback",
    "--gen-key",
  }
  
  if !noPassphrase && passphrase != "" {
    args = append(args, "--passphrase", passphrase)
  }
  
  // Create the key specification
  keySpec := []string{
    "Key-Type: RSA",
    "Key-Length: 2048",
    "Subkey-Type: RSA",
    "Subkey-Length: 2048",
    "Name-Real: " + name,
    "Name-Email: " + email,
    "Expire-Date: 0",
    "%commit",
  }
  
  if noPassphrase {
    keySpec = append([]string{"%no-protection"}, keySpec...)
  }
  
  cmd := exec.Command("gpg", args...)
  cmd.Env = append(os.Environ(), "GNUPGHOME="+gpgHome)
  cmd.Stdin = strings.NewReader(strings.Join(keySpec, "\n"))
  
  var stderr bytes.Buffer
  cmd.Stderr = &stderr
  
  if err := cmd.Run(); err != nil {
    return "", fmt.Errorf("failed to generate key: %v (stderr: %s)", err, stderr.String())
  }
  
  // Get the key ID
  cmd = exec.Command("gpg", "--list-keys", "--with-colons", email)
  cmd.Env = append(os.Environ(), "GNUPGHOME="+gpgHome)
  var stdout bytes.Buffer
  cmd.Stdout = &stdout
  
  if err := cmd.Run(); err != nil {
    return "", err
  }
  
  // Parse the output to find the key
  lines := strings.Split(stdout.String(), "\n")
  for _, line := range lines {
    if strings.HasPrefix(line, "pub:") {
      parts := strings.Split(line, ":")
      if len(parts) >= 5 {
        return parts[4], nil // Key ID
      }
    }
  }
  
  return "", fmt.Errorf("could not find generated key for %s", email)
}

// CreateTestPassword creates a test password file in the password store.
// The password is encrypted using the specified key ID.
func (e *TestEnv) CreateTestPassword(path, content, recipient string) error {
  // Create the directory structure
  fullPath := filepath.Join(e.PasswordStore, path)
  if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
    return err
  }
  
  // Create a temporary plaintext file
  tempFile := filepath.Join(e.TempDir, "temp-password.txt")
  if err := os.WriteFile(tempFile, []byte(content), 0600); err != nil {
    return err
  }
  defer os.Remove(tempFile)
  
  // Encrypt it
  encryptedPath := fullPath + ".gpg"
  cmd := exec.Command("gpg", 
    "--batch",
    "--yes",
    "--encrypt",
    "--armor",
    "--recipient", recipient,
    "--output", encryptedPath,
    tempFile)
  
  cmd.Env = append(os.Environ(), "GNUPGHOME="+e.GNUPGHome)
  
  if err := cmd.Run(); err != nil {
    return err
  }
  
  return nil
}

// ReadTestPassword reads and decrypts a test password file.
// Returns the decrypted content.
func (e *TestEnv) ReadTestPassword(path, passphrase string) (string, error) {
  fullPath := filepath.Join(e.PasswordStore, path+".gpg")
  
  args := []string{"--batch", "--decrypt"}
  
  if passphrase != "" {
    args = append(args, "--passphrase", passphrase)
  }
  
  args = append(args, fullPath)
  
  cmd := exec.Command("gpg", args...)
  cmd.Env = append(os.Environ(), "GNUPGHOME="+e.GNUPGHome)
  
  var stdout, stderr bytes.Buffer
  cmd.Stdout = &stdout
  cmd.Stderr = &stderr
  
  if err := cmd.Run(); err != nil {
    return "", fmt.Errorf("decryption failed: %v (stderr: %s)", err, stderr.String())
  }
  
  return strings.TrimSuffix(stdout.String(), "\n"), nil
}
