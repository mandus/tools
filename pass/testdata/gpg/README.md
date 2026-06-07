# GPG Test Fixtures

This directory contains GPG keys and configuration for testing purposes only.

## Important Security Notes

⚠️ **NEVER** commit real GPG keys or personal secrets to this repository.
⚠️ These are test-only keys generated specifically for automated testing.
⚠️ All test keys should be:
   - Generated with a known passphrase (for testing)
   - Used only in test environments
   - Never used for real password storage

## Test Key Generation

Test keys are generated using the `generate_test_keys.sh` script.
The script creates:
- A test GPG home directory (`test-gnupg-home`)
- Test keys with known passphrases
- Test keys without passphrases

## Usage in Tests

Tests should:
1. Set `GNUPGHOME` to point to the test GPG home directory
2. Use the test key IDs defined in this package
3. Clean up any test data after execution
4. Never rely on the user's personal GPG setup
