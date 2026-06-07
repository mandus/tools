# Pass Test Fixtures Specification

## Overview

This specification describes the implementation of test fixtures for the `pass` password manager, ensuring that:

1. Tests do not rely on the user's personal GPG setup or password store
2. No personal secrets are committed to the repository
3. Tests can run with GPG keys that have passphrases without blocking on gpg-agent
4. The test environment is isolated and self-contained

## Problem Statement

Currently, the `pass` test suite has the following issues:

1. **Personal GPG Setup Dependency**: Tests rely on the user's personal GPG keys being available
2. **gpg-agent Blocking**: Tests with passphrase-protected keys trigger gpg-agent and block forever
3. **Risk of Secret Leakage**: There's a risk of accidentally committing personal keys or passwords
4. **Non-Deterministic Tests**: Test results depend on the user's environment

## Requirements

### Must Have

- [ ] Test fixtures directory with generated GPG keys for testing
- [ ] Tests use isolated GPG home directory (`GNUPGHOME`) pointing to test fixtures
- [ ] Tests use isolated password store directory (`PASSWORD_STORE_DIR`) for test data
- [ ] Support for batch mode GPG operations (no gpg-agent prompts)
- [ ] Support for testing with passphrase-protected keys using `--passphrase` flag
- [ ] Cleanup of all test data after test execution
- [ ] No personal secrets committed to the repository

### Should Have

- [ ] Pre-generated test GPG keys committed to the repository (in `testdata/gpg/`)
- [ ] Script to regenerate test keys when needed
- [ ] Helper functions for creating test passwords and reading them back
- [ ] Documentation explaining how to use test fixtures

### Nice to Have

- [ ] CI/CD integration that automatically sets up test keys
- [ ] Multiple test key types (RSA, ECC, etc.)
- [ ] Test keys with different expiration dates
- [ ] Test keys with multiple subkeys

## Architecture

### Directory Structure

```
pass/
├── testdata/
│   ├── gpg/
│   │   ├── README.md                 # Documentation for test keys
│   │   ├── generate_test_keys.sh     # Script to generate test keys
│   │   ├── test-gnupg-home/          # Test GPG home directory
│   │   │   ├── pubring.kbx           # Public keyring
│   │   │   ├── secring.kbx           # Secret keyring (deprecated in newer GPG)
│   │   │   ├── private-keys-v1.d/     # Secret keys (newer GPG)
│   │   │   ├── gpg.conf              # GPG configuration
│   │   │   └── gpg-agent.conf        # gpg-agent configuration
│   │   ├── fixtures.go               # Go package with test constants
│   │   └── key-ids.txt               # Generated key IDs
│   └── store/                        # Test password store
│       └── .password-store/          # Actual store directory
└── internal/
    └── testhelper/
        └── testhelper.go             # Test helper functions
```

### Test Key Generation

Test keys are generated using the `generate_test_keys.sh` script:

1. Creates a temporary GPG home directory
2. Generates a test key **without** passphrase (for basic tests)
3. Generates a test key **with** passphrase (for passphrase tests)
4. Exports the keys and configuration
5. Outputs key IDs to `key-ids.txt`

The script uses:
- Key type: RSA 2048-bit
- No expiration date
- Known passphrase: `test-passphrase-123`
- Test email addresses that don't exist

### GPG Configuration for Tests

The test GPG configuration (`test-gnupg-home/gpg.conf`):

```
batch
no-tty
yes
```

The test gpg-agent configuration (`test-gnupg-home/gpg-agent.conf`):

```
pinentry-program /bin/false
```

This prevents gpg-agent from prompting for passphrases, which would block tests.

### Batch Mode Operations

For tests that need to run without user interaction, we use:

```bash
# Encryption
gpg --batch --yes --encrypt --armor --recipient <key-id> --output <out> <in>

# Decryption with passphrase
gpg --batch --passphrase <passphrase> --decrypt <file>

# Decryption without passphrase
gpg --batch --pinentry-mode loopback --decrypt <file>
```

The `--pinentry-mode loopback` option allows GPG to use the passphrase from the command line or environment without prompting.

## Implementation Details

### GPG Package Changes

The `pass/pkg/gpg` package will be extended with:

```go
// GPGOptions contains options for GPG operations
type GPGOptions struct {
    BatchMode    bool
    Passphrase   string
    Recipient    string
    UseAgent     bool
    PinentryMode string
}

// DefaultGPGOptions returns default options
func DefaultGPGOptions() GPGOptions

// BatchGPGOptions returns options for batch mode
func BatchGPGOptions(passphrase string) GPGOptions

// EncryptFileWithOptions encrypts with custom options
func EncryptFileWithOptions(srcPath, destPath string, opts GPGOptions) error

// DecryptFileWithOptions decrypts with custom options
func DecryptFileWithOptions(filePath string, opts GPGOptions) (string, error)

// CheckGPGBatch checks if batch mode is available
func CheckGPGBatch() error
```

### Test Helper Package

The `pass/internal/testhelper` package provides:

```go
// TestEnv holds test environment configuration
type TestEnv struct {
    TempDir       string
    GNUPGHome     string
    PasswordStore string
}

// SetupTestEnv creates a clean test environment
func SetupTestEnv() (*TestEnv, error)

// SetupTestEnvWithGPGKeys creates environment with GPG keys
func SetupTestEnvWithGPGKeys() (*TestEnv, string, string, error)

// Cleanup restores original environment
func (e *TestEnv) Cleanup()

// CreateTestPassword creates a test password
func (e *TestEnv) CreateTestPassword(path, content, recipient string) error

// ReadTestPassword reads a test password
func (e *TestEnv) ReadTestPassword(path, passphrase string) (string, error)
```

### Test Fixtures Package

The `pass/testdata/gpg` package provides constants:

```go
// Constants for test keys
const NoPassphraseKeyID = "..."
const WithPassphraseKeyID = "..."
const TestPassphrase = "test-passphrase-123"

// SetupTestEnvironment configures environment for tests
func SetupTestEnvironment() func()
```

## Testing Strategy

### Unit Tests

Unit tests that don't require actual GPG operations:
- Path manipulation
- String parsing
- Error handling

These tests don't need any special setup.

### Integration Tests

Integration tests that require GPG:
1. Set up test environment with `SetupTestEnv()` or `SetupTestEnvWithGPGKeys()`
2. Use the test GPG keys for encryption/decryption
3. Create test passwords in the test password store
4. Clean up with `Cleanup()`

Example:

```go
func TestEncryptDecrypt(t *testing.T) {
    env, noPassKey, _, err := testhelper.SetupTestEnvWithGPGKeys()
    if err != nil {
        t.Skip("GPG not available for testing")
    }
    defer env.Cleanup()
    
    // Set the recipient to use the test key
    os.Setenv("PASS_GPG_ID", noPassKey)
    
    // Create a test password
    if err := env.CreateTestPassword("test/password", "my-secret-123", noPassKey); err != nil {
        t.Fatal(err)
    }
    
    // Test show command
    content, err := env.ReadTestPassword("test/password", "")
    if err != nil {
        t.Fatal(err)
    }
    
    if content != "my-secret-123" {
        t.Errorf("Expected 'my-secret-123', got %q", content)
    }
}
```

### Passphrase Tests

Tests that require passphrase-protected keys:

```go
func TestDecryptWithPassphrase(t *testing.T) {
    env, _, withPassKey, err := testhelper.SetupTestEnvWithGPGKeys()
    if err != nil {
        t.Skip("GPG not available for testing")
    }
    defer env.Cleanup()
    
    // Create a password with the passphrase-protected key
    if err := env.CreateTestPassword("test/password", "my-secret-456", withPassKey); err != nil {
        t.Fatal(err)
    }
    
    // Read it back with the passphrase
    content, err := env.ReadTestPassword("test/password", "test-passphrase-123")
    if err != nil {
        t.Fatal(err)
    }
    
    if content != "my-secret-456" {
        t.Errorf("Expected 'my-secret-456', got %q", content)
    }
}
```

## Security Considerations

1. **Test Keys Only**: All keys in `testdata/gpg/` are test-only and should never be used for real password storage
2. **No Personal Data**: No personal keys, passwords, or secrets should ever be committed
3. **Cleanup**: All test data must be cleaned up after tests complete
4. **Isolation**: Tests must use isolated GPG home and password store directories
5. **Passphrase**: The test passphrase is hardcoded and safe to include in the repository

## CI/CD Integration

In CI/CD environments:

1. The test keys will be pre-generated and committed to the repository
2. Tests will use the committed test keys
3. No user interaction will be required
4. All tests should pass in a clean environment

## Branch and Commit Strategy

Following the repository's AGENTS.md guidelines:

- **Branch**: `feat/<number>-test-fixtures` (e.g., `feat/42-test-fixtures`)
- **Commits**: Use appropriate gitmojis:
  - ✨ feat: Add test fixture infrastructure
  - 🐛 fix: Prevent gpg-agent blocking in tests
  - 🧪 test: Add tests using test fixtures
  - 📝 docs: Add documentation for test fixtures

## Files to Modify

1. **New Files**:
   - `pass/testdata/gpg/README.md`
   - `pass/testdata/gpg/generate_test_keys.sh`
   - `pass/testdata/gpg/fixtures.go`
   - `pass/internal/testhelper/testhelper.go`
   - `specs/pass-test-fixtures-spec.md`

2. **Modified Files**:
   - `pass/pkg/gpg/gpg.go` (add options support)
   - `pass/pkg/gpg/gpg_test.go` (update to use fixtures)
   - `pass/cmd/*_test.go` (update to use test helpers)

3. **Generated Files** (not committed):
   - `pass/testdata/gpg/test-gnupg-home/*`
   - `pass/testdata/gpg/key-ids.txt`

## Verification

To verify the implementation:

1. Run all tests: `cd pass && go test ./...`
2. Verify no personal keys are used: `grep -r "personal\|secret\|private" pass/testdata/ --exclude-dir=test-gnupg-home`
3. Verify test keys are in the repository: `ls -la pass/testdata/gpg/test-gnupg-home/`
4. Verify tests pass without user interaction
