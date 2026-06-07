# Pass Test Fixtures Specification

## Overview

This specification describes the implementation of **ephemeral** test fixtures for the `pass` password manager, ensuring that:

1. Tests do not rely on the user's personal GPG setup or password store
2. No personal secrets are committed to the repository
3. Tests with GPG keys (including passphrase-protected) work without blocking
4. All test data is automatically cleaned up after tests complete

## Key Design Decision: Ephemeral Keys

**All test GPG keys are generated on-the-fly during test execution.**

This means:
- ✅ No pre-generated keys committed to the repository
- ✅ Each test run gets fresh, isolated keys
- ✅ Maximum security - no static keys to compromise
- ✅ No maintenance burden for key rotation
- ⚠️ Tests are slightly slower due to key generation (~5-10 seconds)

## Architecture

### Test Helper Package

The `pass/internal/testhelper` package provides utilities for setting up isolated test environments:

```
pass/internal/testhelper/
└── testhelper.go          # Test environment utilities
```

### Core Functions

```go
// SetupTestEnv creates a clean test environment with temporary directories
func SetupTestEnv() (*TestEnv, error)

// SetupTestEnvWithGPGKeys creates environment AND generates ephemeral GPG keys
func SetupTestEnvWithGPGKeys() (*TestEnv, string, string, error)

// Cleanup restores original environment and removes temporary files
func (e *TestEnv) Cleanup()

// CreateTestPassword creates a test password in the test store
func (e *TestEnv) CreateTestPassword(path, content, recipient string) error

// ReadTestPassword reads and decrypts a test password
func (e *TestEnv) ReadTestPassword(path, passphrase string) (string, error)
```

### GPG Package Enhancements

The `pass/pkg/gpg` package has been enhanced to support batch mode operations:

```go
type GPGOptions struct {
    BatchMode      bool
    Passphrase     string
    Recipient      string
    UseAgent       bool
    PinentryMode   string
}

func EncryptFileWithOptions(srcPath, destPath string, opts GPGOptions) error
func DecryptFileWithOptions(filePath string, opts GPGOptions) (string, error)
func BatchGPGOptions(passphrase string) GPGOptions
```

## Test Key Generation

When `SetupTestEnvWithGPGKeys()` is called, it:

1. Creates a temporary GPG home directory
2. Generates a **key without passphrase** (for basic tests)
   - Name: Test User NoPass
   - Email: test-nopass@example.com
   - Passphrase: (none)
   
3. Generates a **key with passphrase** (for passphrase tests)
   - Name: Test User WithPass
   - Email: test-withpass@example.com
   - Passphrase: `test-passphrase-123`

4. Returns the key IDs for use in tests

## Batch Mode Operations

To prevent gpg-agent from blocking tests, we use:

```bash
# Encrypt in batch mode
gpg --batch --yes --encrypt --armor --recipient <key-id> --output <out> <in>

# Decrypt in batch mode without passphrase
gpg --batch --pinentry-mode loopback --decrypt <file>

# Decrypt in batch mode with passphrase
gpg --batch --passphrase <passphrase> --decrypt <file>
```

The `--pinentry-mode loopback` option allows GPG to bypass the pinentry prompt.

## Usage in Tests

### Basic Test (No GPG Required)

For tests that don't need actual GPG operations (e.g., path manipulation):

```go
func TestPathNormalization(t *testing.T) {
    env, err := testhelper.SetupTestEnv()
    if err != nil {
        t.Fatal(err)
    }
    defer env.Cleanup()
    
    // Test code here
}
```

### Test with GPG (Ephemeral Keys)

For tests that need actual GPG encryption/decryption:

```go
func TestEncryptDecrypt(t *testing.T) {
    // Set up environment with ephemeral GPG keys
    env, noPassKey, _, err := testhelper.SetupTestEnvWithGPGKeys()
    if err != nil {
        t.Skip("GPG not available for testing")
    }
    defer env.Cleanup()
    
    // Set the recipient
    os.Setenv("PASS_GPG_ID", noPassKey)
    os.Setenv("GNUPGHOME", env.GNUPGHome)
    
    // Create a test password
    if err := env.CreateTestPassword("test/password", "my-secret", noPassKey); err != nil {
        t.Fatal(err)
    }
    
    // Read it back
    content, err := env.ReadTestPassword("test/password", "")
    if err != nil {
        t.Fatal(err)
    }
    
    if content != "my-secret" {
        t.Errorf("Expected 'my-secret', got %q", content)
    }
}
```

### Test with Passphrase-Protected Keys

For tests that need passphrase-protected keys:

```go
func TestDecryptWithPassphrase(t *testing.T) {
    env, _, withPassKey, err := testhelper.SetupTestEnvWithGPGKeys()
    if err != nil {
        t.Skip("GPG not available for testing")
    }
    defer env.Cleanup()
    
    os.Setenv("PASS_GPG_ID", withPassKey)
    os.Setenv("GNUPGHOME", env.GNUPGHome)
    
    // Create password with passphrase-protected key
    if err := env.CreateTestPassword("test/password", "my-secret", withPassKey); err != nil {
        t.Fatal(err)
    }
    
    // Read with passphrase
    content, err := env.ReadTestPassword("test/password", "test-passphrase-123")
    if err != nil {
        t.Fatal(err)
    }
    
    if content != "my-secret" {
        t.Errorf("Expected 'my-secret', got %q", content)
    }
}
```

## Security Considerations

✅ **No Personal Data**: All keys are generated on-the-fly with test-only identities
✅ **Isolated Environment**: Each test gets its own GPG home and password store
✅ **Automatic Cleanup**: All test data is removed after tests complete
✅ **Safe Passphrase**: The test passphrase is hardcoded and only used with ephemeral keys
✅ **No Static Keys**: No keys are committed to the repository

## Requirements

### Must Have
- [x] Ephemeral GPG key generation during test execution
- [x] Isolated test environment (GNUPGHOME, PASSWORD_STORE_DIR)
- [x] Batch mode GPG operations
- [x] Passphrase support for tests
- [x] Automatic cleanup

### Should Have
- [x] Test helper package with easy-to-use functions
- [x] Support for both GPG and non-GPG tests
- [x] Documentation and examples

## Implementation Notes

### Performance

Key generation takes approximately 5-10 seconds. This is acceptable because:
1. It only happens once per test run (keys are reused within the same test)
2. Tests that don't need GPG are not affected
3. The security and maintainability benefits outweigh the performance cost

### GPG Availability

Tests that require GPG will skip if GPG is not available:
```go
if err := testhelper.SetupTestEnvWithGPGKeys(); err != nil {
    t.Skip("GPG not available for testing")
}
```

### Cross-Platform Support

The implementation uses:
- `os.MkdirTemp` for temporary directories
- `os.Setenv` for environment variables
- `exec.Command` for GPG operations
- All paths use `filepath.Join` for cross-platform compatibility

## Files Modified/Created

### New Files
- `pass/internal/testhelper/testhelper.go` - Test helper package
- `pass/pkg/gpg/gpg_fixture_test.go` - Tests using ephemeral keys

### Modified Files
- `pass/pkg/gpg/gpg.go` - Added batch mode support

### Removed Files
- Pre-generated key scripts and constants (using ephemeral approach instead)

## Verification

To verify the implementation:

```bash
# Run all tests
cd pass
go test ./...

# Run specific test packages
go test ./pkg/gpg/...
go test ./cmd/...

# Check that no keys are committed
git ls-files | grep -E "\.gpg|gnupg|secret|private" || echo "No keys committed - GOOD!"
```

## References

- [GitHub Spec Kit](https://github.com/github/spec-kit)
- [Repository AGENTS.md](../AGENTS.md)
