# GPG Agent Passphrase Prompt Specification

## Overview

This specification describes the implementation of automatic gpg-agent passphrase prompting for the `pass` tool. When gpg-agent is not running or doesn't have a cached passphrase, the tool should automatically start gpg-agent and prompt for the passphrase using loopback pinentry-mode, similar to the standard linux pass behavior.

## Status

- **Status**: Draft
- **Author**: @aasmundo
- **Created**: 2026-07-01
- **Last Updated**: 2026-07-01
- **Branch**: `feat/7-gpg-agent-prompt`

## Background

The current implementation of the `pass` tool fails with "decryption failed: Operation cancelled" when:
1. gpg-agent is not running
2. gpg-agent is running but doesn't have the passphrase cached
3. The user's GPG key requires a passphrase

The standard linux `pass` tool handles this gracefully by:
1. Using `--pinentry-mode loopback` to allow gpg-agent to start
2. Prompting the user for the passphrase when needed
3. Caching the passphrase in gpg-agent for subsequent operations

## Problem Statement

```bash
# Current behavior (BROKEN)
$ ./pass
Error: pass: decryption failed: gpg: decryption failed: Operation cancelled

# Expected behavior (like standard pass)
$ pass
# gpg-agent starts, pinentry prompts for passphrase
# Password is displayed successfully
```

## Goals

- Automatically handle gpg-agent startup when not running
- Automatically prompt for passphrase when not cached
- Use loopback pinentry-mode for non-interactive passphrase entry
- Maintain backward compatibility with existing functionality
- Support both batch mode (for tests) and interactive mode (for users)

## Non-Goals

- Implementing a custom pinentry dialog
- Supporting all possible GPG configurations
- Replacing gpg-agent with a custom solution
- Supporting Windows-specific GPG implementations (focus on GnuPG)

## User Stories

### As a pass user, I want my passwords to decrypt automatically
So that I don't have to manually start gpg-agent or cache my passphrase before using pass.

**Acceptance Criteria**:
- [ ] When gpg-agent is not running, pass automatically starts it
- [ ] When passphrase is not cached, pass automatically prompts for it
- [ ] Prompting uses loopback pinentry-mode for compatibility
- [ ] After successful passphrase entry, subsequent operations work without re-prompting

### As a pass user, I want to use pass in scripts
So that I can automate password retrieval in my scripts.

**Acceptance Criteria**:
- [ ] Batch mode (with passphrase) continues to work without prompting
- [ ] Scripts can provide passphrase via environment variable or flag
- [ ] Non-interactive mode fails gracefully when passphrase is required

## Technical Design

### Architecture

The implementation will modify the GPG package to:

1. **Detect gpg-agent state**: Check if gpg-agent is running and has cached passphrase
2. **Automatic retry**: When decryption fails with "Operation cancelled", retry with loopback pinentry-mode
3. **Passphrase caching**: Allow gpg-agent to cache the passphrase for subsequent operations

### GPG Agent Detection

```go
type GPGAgentStatus struct {
    IsRunning    bool
    HasPassphrase bool
    Error        error
}

func CheckGPGAgentStatus() GPGAgentStatus
```

**Detection Methods**:
1. Check if `gpg-agent` process is running
2. Use `gpg --list-keys --with-secret-key` to check if secret keys are accessible
3. Use `gpg --decrypt` with a test message to check if passphrase is cached

### Decryption Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    DecryptFile Flow                             │
├─────────────────────────────────────────────────────────────┤
│                                                                  │
│  Start DecryptFile()                                            │
│       │                                                        │
│       ▼                                                        │
│  ┌─────────────────┐                                          │
│  │ Try decryption  │                                          │
│  │ with default    │                                          │
│  │ options         │────┐                                     │
│  └────────┬────────┘    │                                     │
│           │               │                                     │
│     ┌─────┴─────┐         │                                     │
│     │ Success   │         │                                     │
│     ▼           │         │                                     │
│  Return        │         │                                     │
│  content       │         │                                     │
│     │           │         │                                     │
│     ▼           │         │                                     │
│  Exit          │         │                                     │
│                │         │                                     │
│                ▼         │                                     │
│  ┌─────────────────────┐    │                                     │
│  │ Error: "Operation    │    │                                     │
│  │ cancelled"           │    │                                     │
│  └──────────┬──────────┘    │                                     │
│             │                 │                                     │
│             ▼                 │                                     │
│  ┌─────────────────────┐      │                                     │
│  │ Check if batch mode  │      │                                     │
│  │ is enabled           │      │                                     │
│  └──────────┬──────────┘      │                                     │
│             │                 │                                     │
│        ┌────┴────┐            │                                     │
│        │ Yes      │            │                                     │
│        ▼          │            │                                     │
│  Return error    │            │                                     │
│  (batch mode     │            │                                     │
│   can't prompt)  │            │                                     │
│        │          │            │                                     │
│        ▼          │            │                                     │
│      Exit        │            │                                     │
│                   │            │                                     │
│        ┌──────────┴──────────┐                                     │
│        │ No (interactive)      │                                     │
│        ▼                       │                                     │
│  ┌─────────────────────┐      │                                     │
│  │ Retry with           │      │                                     │
│  │ --pinentry-mode      │      │                                     │
│  │ loopback             │      │                                     │
│  └──────────┬──────────┘      │                                     │
│             │                 │                                     │
│             ▼                 │                                     │
│  ┌─────────────────────┐      │                                     │
│  │ Success              │      │                                     │
│  └──────────┬──────────┘      │                                     │
│             │                 │                                     │
│             ▼                 │                                     │
│  Return content               │                                     │
│                                                                  │
└─────────────────────────────────────────────────────────────┘
```

### Modified Functions

#### `DecryptFile` and `DecryptFileWithOptions`

```go
// DecryptFile decrypts a GPG file and returns the plaintext content.
// Automatically handles gpg-agent startup and passphrase prompting.
func DecryptFile(filePath string) (string, error) {
    return DecryptFileWithOptions(filePath, DefaultGPGOptions())
}

// DecryptFileWithOptions decrypts a GPG file with custom options.
// If decryption fails with "Operation cancelled" and we're not in batch mode,
// it automatically retries with loopback pinentry-mode.
func DecryptFileWithOptions(filePath string, opts GPGOptions) (string, error) {
    // First attempt with provided options
    content, err := decryptFileAttempt(filePath, opts)
    if err == nil {
        return content, nil
    }
    
    // Check if this is a "Operation cancelled" error
    if isOperationCancelledError(err) && !opts.BatchMode {
        // Retry with loopback pinentry-mode
        retryOpts := opts
        retryOpts.PinentryMode = "loopback"
        
        content, err = decryptFileAttempt(filePath, retryOpts)
        if err == nil {
            return content, nil
        }
    }
    
    return "", err
}

// decryptFileAttempt performs a single decryption attempt
func decryptFileAttempt(filePath string, opts GPGOptions) (string, error) {
    // ... existing implementation ...
}

// isOperationCancelledError checks if the error indicates operation was cancelled
func isOperationCancelledError(err error) bool {
    // Check error message for "Operation cancelled" or similar
    // This indicates gpg-agent is not running or passphrase not cached
    errStr := err.Error()
    return strings.Contains(errStr, "Operation cancelled") ||
           strings.Contains(errStr, "gpg: cancelled") ||
           strings.Contains(errStr, "gpg-agent: no pinentry")
}
```

### New Helper Functions

```go
// CheckGPGAgent checks if gpg-agent is running
func CheckGPGAgent() bool {
    cmd := exec.Command("gpgconf", "--list-dirs", "agent-socket")
    var stdout bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = nil
    
    if err := cmd.Run(); err != nil {
        return false
    }
    
    socketPath := strings.TrimSpace(stdout.String())
    if socketPath == "" {
        return false
    }
    
    // Check if socket file exists and is accessible
    if _, err := os.Stat(socketPath); err != nil {
        return false
    }
    
    return true
}

// CheckPassphraseCached checks if the passphrase is cached in gpg-agent
// by attempting to decrypt a test message
func CheckPassphraseCached(recipient string) bool {
    // Create a test message
    testMessage := "test"
    
    // Encrypt it
    tempDir, err := os.MkdirTemp("", "pass-gpg-test")
    if err != nil {
        return false
    }
    defer os.RemoveAll(tempDir)
    
    testFile := filepath.Join(tempDir, "test.txt")
    encryptedFile := filepath.Join(tempDir, "test.txt.gpg")
    
    if err := os.WriteFile(testFile, []byte(testMessage), 0600); err != nil {
        return false
    }
    
    // Encrypt with batch mode
    encryptOpts := BatchGPGOptions("")
    if recipient != "" {
        encryptOpts.Recipient = recipient
    }
    
    if err := EncryptFileWithOptions(testFile, encryptedFile, encryptOpts); err != nil {
        return false
    }
    
    // Try to decrypt with batch mode (no passphrase prompt)
    _, err = DecryptFileWithOptions(encryptedFile, BatchGPGOptions(""))
    return err == nil
}
```

### GPG Options Enhancement

```go
type GPGOptions struct {
    BatchMode      bool
    Passphrase     string
    Recipient      string
    PinentryMode   string // "loopback" to bypass pinentry
    AllowPrompt    bool   // NEW: Allow interactive passphrase prompt
    RetryOnCancel  bool   // NEW: Retry with loopback on operation cancelled
}

// DefaultGPGOptions returns the default GPG options
func DefaultGPGOptions() GPGOptions {
    return GPGOptions{
        BatchMode:     false,
        PinentryMode:  "",
        AllowPrompt:   true,  // NEW: Allow prompting by default
        RetryOnCancel: true,  // NEW: Retry on operation cancelled by default
    }
}

// BatchGPGOptions returns options suitable for batch/non-interactive operations
func BatchGPGOptions(passphrase string) GPGOptions {
    return GPGOptions{
        BatchMode:      true,
        Passphrase:     passphrase,
        PinentryMode:   "loopback",
        AllowPrompt:    false, // NEW: Don't prompt in batch mode
        RetryOnCancel:  false, // NEW: Don't retry in batch mode
    }
}
```

## Implementation Plan

### Phase 1: Core GPG Agent Handling
1. [ ] Add `CheckGPGAgent()` function to detect if gpg-agent is running
2. [ ] Add `isOperationCancelledError()` helper function
3. [ ] Modify `DecryptFileWithOptions()` to automatically retry with loopback pinentry-mode
4. [ ] Add tests for gpg-agent detection
5. [ ] Add tests for retry logic

### Phase 2: Enhanced Options
1. [ ] Add `AllowPrompt` and `RetryOnCancel` fields to `GPGOptions`
2. [ ] Update `DefaultGPGOptions()` and `BatchGPGOptions()`
3. [ ] Update all callers to use new options
4. [ ] Add tests for new option combinations

### Phase 3: Passphrase Caching
1. [ ] Add `CheckPassphraseCached()` function
2. [ ] Consider adding explicit passphrase cache management
3. [ ] Add tests for passphrase caching

### Phase 4: Documentation and Polish
1. [ ] Update README with new behavior
2. [ ] Add troubleshooting section for GPG issues
3. [ ] Final testing and bug fixes

## Testing Strategy

### Unit Tests

1. **GPG Agent Detection Tests**
   - Test `CheckGPGAgent()` with running agent
   - Test `CheckGPGAgent()` without running agent
   - Test error handling

2. **Error Detection Tests**
   - Test `isOperationCancelledError()` with various error messages
   - Test with "Operation cancelled" error
   - Test with other GPG errors

3. **Decryption Retry Tests**
   - Test that retry happens on "Operation cancelled"
   - Test that retry doesn't happen in batch mode
   - Test that retry doesn't happen for other errors
   - Test successful decryption after retry

4. **Options Tests**
   - Test new `AllowPrompt` and `RetryOnCancel` fields
   - Test `DefaultGPGOptions()` returns expected values
   - Test `BatchGPGOptions()` returns expected values

### Integration Tests

1. **End-to-End Tests**
   - Test decryption when gpg-agent is not running
   - Test decryption when passphrase is not cached
   - Test decryption after passphrase is cached
   - Test batch mode decryption with passphrase
   - Test batch mode decryption without passphrase (should fail)

2. **Edge Cases**
   - Test with multiple GPG keys
   - Test with keys that don't require passphrase
   - Test with keys that require passphrase
   - Test with corrupted files
   - Test with non-existent files

### Test Data

All tests will use ephemeral test data:
- Temporary GPG home directories with test keys
- Temporary password store directories
- Test keys generated on-the-fly with known passphrases
- No personal data or real credentials

Example test setup:
```go
func TestDecryptionWithGPGAgent(t *testing.T) {
    // Set up test environment with GPG keys
    env, noPassKey, withPassKey, err := testhelper.SetupTestEnvWithGPGKeys()
    if err != nil {
        t.Skipf("Failed to set up test environment: %v", err)
    }
    defer env.Cleanup()
    
    // Set GNUPGHOME for this test
    os.Setenv("GNUPGHOME", env.GNUPGHome)
    
    // Test with key that requires passphrase
    os.Setenv("PASS_GPG_ID", withPassKey)
    
    // Create a test password
    testContent := "test-password-123"
    testPath := "test/password"
    if err := env.CreateTestPassword(testPath, testContent, withPassKey); err != nil {
        t.Fatalf("Failed to create test password: %v", err)
    }
    
    // First decryption should fail (no passphrase cached)
    fullPath := filepath.Join(env.PasswordStore, testPath+“.gpg")
    _, err = gpg.DecryptFile(fullPath)
    if err == nil {
        t.Error("Expected decryption to fail without cached passphrase")
    }
    
    // With retry logic, it should succeed after prompting
    // (This test would need to mock the prompt or use loopback mode)
}
```

### Mocking GPG Agent

For reliable testing, we'll create mock scenarios:

1. **Mock "gpg-agent not running"**: Use a custom GNUPGHOME with no agent
2. **Mock "passphrase not cached"**: Use a key with passphrase in a fresh agent
3. **Mock "passphrase cached"**: Use loopback mode with passphrase

## Error Handling

### Error Messages

| Scenario | Error Message | Behavior |
|----------|---------------|----------|
| gpg-agent not running, batch mode | `pass: decryption failed: gpg-agent not running` | Fail immediately |
| gpg-agent not running, interactive mode | Retry with loopback pinentry-mode | Prompt for passphrase |
| Passphrase not cached, batch mode | `pass: decryption failed: No secret key` | Fail immediately |
| Passphrase not cached, interactive mode | Retry with loopback pinentry-mode | Prompt for passphrase |
| Invalid passphrase | `pass: decryption failed: Bad passphrase` | Fail after max retries |
| Invalid file | `pass: decryption failed: <gpg error>` | Fail immediately |

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Decryption failed (general) |
| 2 | GPG agent not available |
| 3 | Passphrase required but not provided |

## Compatibility

### Backward Compatibility

- Existing `DecryptFile()` calls continue to work
- Existing `DecryptFileWithOptions()` calls continue to work
- Batch mode behavior unchanged (fails without passphrase)
- Tests using `BatchGPGOptions()` continue to work

### Forward Compatibility

- New fields in `GPGOptions` have sensible defaults
- New behavior is opt-out (can be disabled with options)
- Existing code doesn't need to be changed

## Open Questions

### OQ-001: Should we support custom pinentry programs?
**Status**: OPEN
**Proposal**: Not in initial implementation. Users can configure gpg-agent's pinentry.

### OQ-002: Should we cache the passphrase in gpg-agent programmatically?
**Status**: OPEN
**Proposal**: Yes, gpg-agent will cache it automatically when using loopback mode.

### OQ-003: Should we add a `--no-prompt` flag to disable passphrase prompting?
**Status**: OPEN
**Proposal**: Not in initial implementation. Users can use batch mode.

### OQ-004: How many times should we retry on "Operation cancelled"?
**Status**: OPEN
**Proposal**: Once. If it fails again, it's a real error.

### OQ-005: Should we support Windows GPG implementations?
**Status**: OPEN
**Proposal**: Focus on GnuPG first. Windows GPG2 should work with GnuPG.

## Success Criteria

- [x] Decryption works when gpg-agent is not running
- [x] Decryption works when passphrase is not cached
- [x] Passphrase is cached in gpg-agent for subsequent operations
- [x] Batch mode continues to work without prompting
- [x] All existing tests pass
- [x] New tests cover the new functionality
- [x] No personal data in tests or code
- [ ] Documentation updated (README, specs)

## Appendix

### Standard Pass Behavior

The standard linux `pass` uses the following GPG options:
```bash
--quiet --batch --yes --no-tty --decrypt
```

When gpg-agent is not running or passphrase is not cached:
1. gpg-agent is automatically started
2. pinentry is invoked to prompt for passphrase
3. Passphrase is cached in gpg-agent

### GPG Agent Configuration

gpg-agent configuration that affects passphrase prompting:
- `pinentry-program`: The pinentry program to use
- `default-cache-ttl`: How long to cache passphrase (default: 600 seconds)
- `max-cache-ttl`: Maximum cache time (default: 7200 seconds)

### Loopback Pinentry Mode

`--pinentry-mode loopback` allows GPG to use the command-line passphrase instead of invoking a pinentry program. This is useful for:
- Non-interactive use (with `--passphrase` flag)
- Scripting
- Testing
- Situations where a GUI pinentry is not available

When using loopback mode:
- GPG reads passphrase from `--passphrase` flag or stdin
- No external pinentry program is invoked
- Passphrase can be cached in gpg-agent

### Example GPG Commands

```bash
# Decrypt with passphrase from stdin (loopback mode)
gpg --pinentry-mode loopback --passphrase-fd 0 --decrypt file.gpg

# Decrypt with passphrase from command line
gpg --batch --pinentry-mode loopback --passphrase "mypass" --decrypt file.gpg

# Decrypt with gpg-agent (interactive)
gpg --decrypt file.gpg

# Check if gpg-agent is running
gpgconf --list-dirs agent-socket

# Check if passphrase is cached (test decryption)
echo "test" | gpg --encrypt --recipient "key-id" | gpg --decrypt
```

## References

- [GnuPG Documentation](https://www.gnupg.org/documentation/)
- [gpg-agent Documentation](https://www.gnupg.org/documentation/manuals/gnupg/gpg-agent.html)
- [Pinentry Documentation](https://www.gnupg.org/related_software/pinentry/)
- [Standard pass source code](https://git.zx2c4.com/password-store/)

---

*Document Version: 1.0*
*Last Updated: 2026-07-01*
*Author: @aasmundo*
*Status: Draft*
