# Pass Test Data

This directory contains test data and fixtures for the `pass` password manager tests.

## Ephemeral Test Keys

**Important**: This implementation uses **ephemeral** GPG keys that are generated on-the-fly during test execution. No pre-generated keys are committed to the repository.

## Directory Structure

```
pass/testdata/
└── gpg/                          # GPG test utilities (no committed keys)
    └── README.md                 # Documentation
```

## How It Works

1. When a test needs GPG keys, it calls `testhelper.SetupTestEnvWithGPGKeys()`
2. This generates temporary GPG keys in a temporary directory
3. The keys are used for the duration of the test
4. All keys and test data are automatically cleaned up after the test

## Benefits

✅ **No committed keys** - Keys are never stored in the repository  
✅ **Isolated** - Each test gets its own keys  
✅ **Secure** - No risk of key compromise  
✅ **Maintainable** - No need to update committed keys  

## Test Key Details

When generated, test keys include:

### Key Without Passphrase
- Name: Test User NoPass
- Email: test-nopass@example.com
- Passphrase: (none)
- Usage: Basic encryption/decryption tests

### Key With Passphrase
- Name: Test User WithPass
- Email: test-withpass@example.com
- Passphrase: `test-passphrase-123`
- Usage: Passphrase-protected encryption tests

## Performance Note

Key generation takes ~5-10 seconds. This is a one-time cost per test run and is acceptable for the security benefits.
