# GPG Test Utilities

This directory contains utilities for GPG testing in the `pass` password manager.

## Ephemeral Keys

**No GPG keys are committed to this directory.**

All test keys are generated on-the-fly during test execution using the `testhelper` package.

## Usage

Tests should use the `testhelper` package to set up GPG testing:

```go
import "github.com/mandu/tools/pass/internal/testhelper"

func TestWithGPG(t *testing.T) {
    env, noPassKey, withPassKey, err := testhelper.SetupTestEnvWithGPGKeys()
    if err != nil {
        t.Skip("GPG not available")
    }
    defer env.Cleanup()
    
    // Use the keys for testing
    os.Setenv("GNUPGHOME", env.GNUPGHome)
    os.Setenv("PASS_GPG_ID", noPassKey)
    
    // ... test code
}
```

## Generated Keys

When `SetupTestEnvWithGPGKeys()` is called, it generates:

1. **Key without passphrase**
   - Identity: Test User NoPass <test-nopass@example.com>
   - Key type: RSA 2048-bit
   - No expiration
   - No passphrase

2. **Key with passphrase**
   - Identity: Test User WithPass <test-withpass@example.com>
   - Key type: RSA 2048-bit
   - No expiration
   - Passphrase: `test-passphrase-123`

## Batch Mode

All GPG operations use batch mode to prevent interactive prompts:

```bash
# Encrypt
gpg --batch --yes --encrypt --armor --recipient <key-id> --output <out> <in>

# Decrypt without passphrase
gpg --batch --pinentry-mode loopback --decrypt <file>

# Decrypt with passphrase
gpg --batch --passphrase <passphrase> --decrypt <file>
```

## Security

✅ No personal keys are ever used or committed  
✅ All keys are ephemeral (generated and destroyed per test)  
✅ Test passphrase is safe to use in code  
✅ Isolated from user's personal GPG setup  
