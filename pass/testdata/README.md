# Pass Test Data

This directory contains test fixtures for the `pass` password manager tests.

## Structure

```
pass/testdata/
├── gpg/                          # GPG test fixtures
│   ├── README.md                 # Documentation for GPG test keys
│   ├── fixtures.go               # Go constants for test keys
│   ├── generate_test_keys.sh     # Script to generate test GPG keys
│   └── test-gnupg-home/          # Generated GPG home directory (not committed)
└── store/                        # Test password store (not committed)
    └── .password-store/          # Actual password store
```

## Important Security Notes

⚠️ **CRITICAL**: This directory MUST NOT contain any real GPG keys or personal secrets.

- All keys in `testdata/gpg/` are test-only keys generated specifically for testing
- Test keys use known passphrases that are safe to commit
- No personal GPG keys should ever be used or committed
- Test data is automatically cleaned up after tests complete

## Setup

### Generating Test Keys

To generate test GPG keys for the first time or to regenerate them:

```bash
cd pass/testdata/gpg
./generate_test_keys.sh
```

This will:
1. Create a `test-gnupg-home/` directory
2. Generate a GPG key without passphrase (for basic tests)
3. Generate a GPG key with passphrase `test-passphrase-123` (for passphrase tests)
4. Export key IDs to `key-ids.txt`

### Using Test Keys in Tests

Tests should use the test helper package:

```go
import "github.com/mandu/tools/pass/internal/testhelper"

func TestSomething(t *testing.T) {
    // Set up test environment
    env, err := testhelper.SetupTestEnv()
    if err != nil {
        t.Fatal(err)
    }
    defer env.Cleanup()
    
    // Or with GPG keys:
    env, noPassKey, withPassKey, err := testhelper.SetupTestEnvWithGPGKeys()
    if err != nil {
        t.Skip("GPG not available")
    }
    defer env.Cleanup()
    
    // Use the test keys
    os.Setenv("PASS_GPG_ID", noPassKey)
    
    // Create test passwords
    if err := env.CreateTestPassword("test/password", "my-secret", noPassKey); err != nil {
        t.Fatal(err)
    }
    
    // Read test passwords
    content, err := env.ReadTestPassword("test/password", "")
    if err != nil {
        t.Fatal(err)
    }
}
```

## Test Key Information

### Key Without Passphrase
- **Name**: Test User NoPass
- **Email**: test-nopass@example.com
- **Passphrase**: (none)
- **Usage**: Basic encryption/decryption tests

### Key With Passphrase
- **Name**: Test User WithPass
- **Email**: test-withpass@example.com
- **Passphrase**: `test-passphrase-123`
- **Usage**: Testing passphrase-protected encryption

## GPG Configuration

The test GPG home directory (`test-gnupg-home/`) contains:

- `gpg.conf`: Configuration for batch mode operations
  ```
  batch
  no-tty
  yes
  ```

- `gpg-agent.conf`: Configuration to prevent prompting
  ```
  pinentry-program /bin/false
  ```

## Batch Mode Operations

For non-interactive testing, use batch mode:

```bash
# Encrypt
gpg --batch --yes --encrypt --armor --recipient <key-id> --output <out> <in>

# Decrypt without passphrase
gpg --batch --pinentry-mode loopback --decrypt <file>

# Decrypt with passphrase
gpg --batch --passphrase <passphrase> --decrypt <file>
```

## Cleanup

All test data is automatically cleaned up when:
1. The test completes (using `defer env.Cleanup()`)
2. The test process exits

The cleanup removes:
- Temporary directories
- Generated GPG keys (in test environments)
- Test password store

## CI/CD Integration

In CI/CD environments:
1. The test keys should be pre-generated
2. Tests will automatically use the test environment
3. No user interaction is required
4. All tests should pass in a clean environment

## Verification

To verify the test setup is working:

```bash
cd pass

# Run all tests
go test ./...

# Run specific test packages
go test ./pkg/gpg/...
go test ./cmd/...

# Check for personal keys (should find none)
grep -r "personal\|secret\|private" testdata/ --exclude-dir=test-gnupg-home
```
