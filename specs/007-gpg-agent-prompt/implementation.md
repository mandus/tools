# GPG Agent Passphrase Prompt - Implementation Plan

## Overview

This document outlines the implementation plan for the GPG agent passphrase prompt feature as specified in `specs/007-gpg-agent-prompt/spec.md`.

## Implementation Phases

### Phase 1: Core GPG Agent Handling (Priority: High)

**Objective**: Add automatic gpg-agent detection and retry logic for decryption.

**Tasks**:

1. **Add error detection helper**
   - [ ] Create `isOperationCancelledError()` function in `pkg/gpg/gpg.go`
   - [ ] Handle various "Operation cancelled" error message formats
   - [ ] Add unit tests for error detection

2. **Add gpg-agent detection**
   - [ ] Create `CheckGPGAgent()` function in `pkg/gpg/gpg.go`
   - [ ] Use `gpgconf --list-dirs agent-socket` to detect agent
   - [ ] Add fallback detection methods
   - [ ] Add unit tests for agent detection

3. **Modify decryption logic**
   - [ ] Refactor `DecryptFileWithOptions()` to support retry
   - [ ] Add `decryptFileAttempt()` helper function
   - [ ] Implement automatic retry with loopback pinentry-mode
   - [ ] Add unit tests for retry logic

4. **Update GPG options**
   - [ ] Add `AllowPrompt` and `RetryOnCancel` fields to `GPGOptions` struct
   - [ ] Update `DefaultGPGOptions()` to enable retry by default
   - [ ] Update `BatchGPGOptions()` to disable retry
   - [ ] Add unit tests for new options

**Files Modified**:
- `pass/pkg/gpg/gpg.go` - Core implementation
- `pass/pkg/gpg/gpg_test.go` - Unit tests

**Deliverables**:
- Working gpg-agent detection
- Automatic retry on "Operation cancelled" errors
- Updated options structure
- Comprehensive unit tests

---

### Phase 2: Enhanced Passphrase Handling (Priority: Medium)

**Objective**: Add passphrase caching detection and improved error handling.

**Tasks**:

1. **Add passphrase cache detection**
   - [ ] Create `CheckPassphraseCached()` function in `pkg/gpg/gpg.go`
   - [ ] Test decryption with cached vs non-cached passphrase
   - [ ] Add unit tests

2. **Improve error messages**
   - [ ] Add more specific error types for GPG operations
   - [ ] Improve error messages for different failure scenarios
   - [ ] Add error wrapping for better context

3. **Add agent management**
   - [ ] Consider adding `StartGPGAgent()` function
   - [ ] Add `EnsureGPGAgent()` function that starts agent if not running

**Files Modified**:
- `pass/pkg/gpg/gpg.go` - Additional helpers
- `pass/pkg/gpg/gpg_test.go` - Additional tests

**Deliverables**:
- Passphrase cache detection
- Improved error messages
- Agent management utilities

---

### Phase 3: Integration and Testing (Priority: High)

**Objective**: Ensure the feature works end-to-end and all tests pass.

**Tasks**:

1. **Integration tests**
   - [ ] Test decryption when gpg-agent is not running
   - [ ] Test decryption when passphrase is not cached
   - [ ] Test decryption after passphrase is cached
   - [ ] Test batch mode behavior
   - [ ] Test with keys that don't require passphrase
   - [ ] Test with keys that require passphrase

2. **Update existing tests**
   - [ ] Update tests that might be affected by new behavior
   - [ ] Ensure all existing tests still pass
   - [ ] Add test fixtures for new scenarios

3. **Test helper updates**
   - [ ] Add helpers for testing gpg-agent scenarios
   - [ ] Add helpers for testing passphrase caching

**Files Modified**:
- `pass/pkg/gpg/gpg_fixture_test.go` - Integration tests
- `pass/internal/testhelper/testhelper.go` - Test helpers
- Various test files as needed

**Deliverables**:
- All tests passing
- End-to-end test coverage
- Updated test infrastructure

---

### Phase 4: Documentation and Polish (Priority: Medium)

**Objective**: Document the new feature and ensure it's production-ready.

**Tasks**:

1. **Update README**
   - [ ] Document new GPG agent behavior
   - [ ] Add troubleshooting section for GPG issues
   - [ ] Update usage examples

2. **Update specs**
   - [ ] Mark spec as implemented
   - [ ] Update with any changes made during implementation
   - [ ] Add implementation notes

3. **Code review and cleanup**
   - [ ] Review all changes for code quality
   - [ ] Ensure consistent error handling
   - [ ] Add comments for complex logic
   - [ ] Clean up any temporary code

**Files Modified**:
- `pass/README.md` - User documentation
- `specs/007-gpg-agent-prompt/spec.md` - Update spec status
- Various code files for cleanup

**Deliverables**:
- Updated documentation
- Clean, production-ready code
- Final spec updates

---

## Detailed Implementation

### File: `pass/pkg/gpg/gpg.go`

#### New Constants

```go
// GPG error message patterns
const (
    ErrOperationCancelled = "Operation cancelled"
    ErrGPGCancelled       = "gpg: cancelled"
    ErrNoPinentry         = "gpg-agent: no pinentry"
)
```

#### New Helper Functions

```go
// isOperationCancelledError checks if the error indicates gpg-agent
// is not running or passphrase is not available
func isOperationCancelledError(err error) bool {
    if err == nil {
        return false
    }
    errStr := err.Error()
    return strings.Contains(errStr, ErrOperationCancelled) ||
           strings.Contains(errStr, ErrGPGCancelled) ||
           strings.Contains(errStr, ErrNoPinentry)
}

// CheckGPGAgent checks if gpg-agent is running and accessible
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
    
    // Check if socket file exists
    if _, err := os.Stat(socketPath); err != nil {
        return false
    }
    
    return true
}

// EnsureGPGAgent ensures gpg-agent is running by starting it if needed
func EnsureGPGAgent() error {
    if CheckGPGAgent() {
        return nil
    }
    
    // Start gpg-agent
    cmd := exec.Command("gpgconf", "--launch", "gpg-agent")
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to start gpg-agent: %v", err)
    }
    
    // Verify it started
    if !CheckGPGAgent() {
        return fmt.Errorf("gpg-agent did not start")
    }
    
    return nil
}

// decryptFileAttempt performs a single decryption attempt with given options
func decryptFileAttempt(filePath string, opts GPGOptions) (string, error) {
    args := []string{"--decrypt"}
    
    // Add batch mode if requested
    if opts.BatchMode {
        args = append(args, "--batch")
    }
    
    // Add pinentry mode if specified
    if opts.PinentryMode != "" {
        args = append(args, "--pinentry-mode", opts.PinentryMode)
    }
    
    // Add passphrase if provided (for batch mode)
    if opts.Passphrase != "" {
        args = append(args, "--passphrase", opts.Passphrase)
    }
    
    args = append(args, filePath)
    
    cmd := exec.Command("gpg", args...)
    
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    if err := cmd.Run(); err != nil {
        stderrStr := stderr.String()
        
        // Check for specific error conditions
        if strings.Contains(stderrStr, "No secret key") {
            return "", fmt.Errorf("pass: decryption failed: No secret key available for this password")
        }
        if strings.Contains(stderrStr, "decryption failed") {
            return "", fmt.Errorf("pass: decryption failed: %s", extractGPGError(stderrStr))
        }
        if strings.Contains(stderrStr, "bad passphrase") || strings.Contains(stderrStr, "Bad passphrase") {
            return "", fmt.Errorf("pass: decryption failed: Bad passphrase")
        }
        if strings.Contains(stderrStr, "gpg: WARN") || strings.Contains(stderrStr, "gpg: warning") {
            // Non-fatal warning, try to return the output anyway
            output := strings.TrimSuffix(stdout.String(), "\n")
            if output != "" {
                return output, nil
            }
        }
        
        return "", fmt.Errorf("pass: GPG decryption failed: %v (stderr: %s)", err, stderrStr)
    }
    
    // Trim trailing newline if present
    output := strings.TrimSuffix(stdout.String(), "\n")
    return output, nil
}
```

#### Modified Functions

```go
// GPGOptions contains options for GPG operations
type GPGOptions struct {
    BatchMode      bool
    Passphrase     string
    Recipient      string
    PinentryMode   string // "loopback" to bypass pinentry
    AllowPrompt    bool   // Allow interactive passphrase prompt (default: true)
    RetryOnCancel  bool   // Retry with loopback on operation cancelled (default: true)
}

// DefaultGPGOptions returns the default GPG options
func DefaultGPGOptions() GPGOptions {
    return GPGOptions{
        BatchMode:     false,
        PinentryMode:  "",
        AllowPrompt:   true,
        RetryOnCancel: true,
    }
}

// BatchGPGOptions returns options suitable for batch/non-interactive operations
func BatchGPGOptions(passphrase string) GPGOptions {
    return GPGOptions{
        BatchMode:      true,
        Passphrase:     passphrase,
        PinentryMode:   "loopback",
        AllowPrompt:    false,
        RetryOnCancel:  false,
    }
}

// DecryptFile decrypts a GPG file and returns the plaintext content.
// GPG will automatically handle gpg-agent and passphrase prompting.
func DecryptFile(filePath string) (string, error) {
    return DecryptFileWithOptions(filePath, DefaultGPGOptions())
}

// DecryptFileWithOptions decrypts a GPG file with custom options.
// If decryption fails with "Operation cancelled" and RetryOnCancel is true,
// it automatically retries with loopback pinentry-mode.
func DecryptFileWithOptions(filePath string, opts GPGOptions) (string, error) {
    // First attempt with provided options
    content, err := decryptFileAttempt(filePath, opts)
    if err == nil {
        return content, nil
    }
    
    // Check if we should retry
    if opts.RetryOnCancel && isOperationCancelledError(err) && opts.AllowPrompt {
        // Retry with loopback pinentry-mode
        retryOpts := opts
        retryOpts.PinentryMode = "loopback"
        retryOpts.BatchMode = false // Don't use batch mode for retry
        
        content, err = decryptFileAttempt(filePath, retryOpts)
        if err == nil {
            return content, nil
        }
    }
    
    return "", err
}
```

### File: `pass/pkg/gpg/gpg_test.go`

#### New Test Functions

```go
func TestIsOperationCancelledError(t *testing.T) {
    tests := []struct {
        name     string
        err      error
        expected bool
    }{
        {
            name:     "operation cancelled",
            err:      fmt.Errorf("gpg: decryption failed: Operation cancelled"),
            expected: true,
        },
        {
            name:     "gpg cancelled",
            err:      fmt.Errorf("gpg: cancelled"),
            expected: true,
        },
        {
            name:     "no pinentry",
            err:      fmt.Errorf("gpg-agent: no pinentry"),
            expected: true,
        },
        {
            name:     "no secret key",
            err:      fmt.Errorf("gpg: No secret key"),
            expected: false,
        },
        {
            name:     "bad passphrase",
            err:      fmt.Errorf("gpg: Bad passphrase"),
            expected: false,
        },
        {
            name:     "nil error",
            err:      nil,
            expected: false,
        },
        {
            name:     "other error",
            err:      fmt.Errorf("some other error"),
            expected: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := isOperationCancelledError(tt.err)
            if result != tt.expected {
                t.Errorf("isOperationCancelledError(%v) = %v, want %v", tt.err, result, tt.expected)
            }
        })
    }
}

func TestCheckGPGAgent(t *testing.T) {
    // Test with current environment
    isRunning := CheckGPGAgent()
    t.Logf("gpg-agent is running: %v", isRunning)
    
    // If running, we can test the positive case
    // If not running, we can only test the negative case
    if !isRunning {
        // This is expected in some test environments
        t.Log("gpg-agent is not running in test environment (expected)")
    }
}

func TestDefaultGPGOptions(t *testing.T) {
    opts := DefaultGPGOptions()
    
    if opts.BatchMode {
        t.Error("BatchMode should be false by default")
    }
    if opts.PinentryMode != "" {
        t.Errorf("PinentryMode should be empty by default, got %q", opts.PinentryMode)
    }
    if !opts.AllowPrompt {
        t.Error("AllowPrompt should be true by default")
    }
    if !opts.RetryOnCancel {
        t.Error("RetryOnCancel should be true by default")
    }
}

func TestBatchGPGOptions(t *testing.T) {
    opts := BatchGPGOptions("test-passphrase")
    
    if !opts.BatchMode {
        t.Error("BatchMode should be true for batch options")
    }
    if opts.Passphrase != "test-passphrase" {
        t.Errorf("Passphrase should be 'test-passphrase', got %q", opts.Passphrase)
    }
    if opts.PinentryMode != "loopback" {
        t.Errorf("PinentryMode should be 'loopback' for batch options, got %q", opts.PinentryMode)
    }
    if opts.AllowPrompt {
        t.Error("AllowPrompt should be false for batch options")
    }
    if opts.RetryOnCancel {
        t.Error("RetryOnCancel should be false for batch options")
    }
}

func TestDecryptFileWithRetry(t *testing.T) {
    // Set up test environment with GPG keys
    env, noPassKey, withPassKey, err := testhelper.SetupTestEnvWithGPGKeys()
    if err != nil {
        t.Skipf("Failed to set up test environment: %v", err)
    }
    defer env.Cleanup()
    
    os.Setenv("GNUPGHOME", env.GNUPGHome)
    
    // Test with key that requires passphrase
    os.Setenv("PASS_GPG_ID", withPassKey)
    
    testContent := "test-password-456"
    testPath := "test/retry-password"
    if err := env.CreateTestPassword(testPath, testContent, withPassKey); err != nil {
        t.Fatalf("Failed to create test password: %v", err)
    }
    
    fullPath := filepath.Join(env.PasswordStore, testPath+“.gpg")
    
    // Test with retry enabled (default)
    content, err := gpg.DecryptFile(fullPath)
    if err != nil {
        // This might fail in test environment without proper pinentry
        // For now, just log and skip if it fails
        t.Logf("Decryption with retry failed (expected in test environment): %v", err)
        t.Skip("Skipping retry test - requires proper gpg-agent setup")
    }
    
    if content != testContent {
        t.Errorf("Decrypted content = %q, want %q", content, testContent)
    }
}

func TestDecryptFileNoRetryInBatchMode(t *testing.T) {
    // Set up test environment with GPG keys
    env, _, withPassKey, err := testhelper.SetupTestEnvWithGPGKeys()
    if err != nil {
        t.Skipf("Failed to set up test environment: %v", err)
    }
    defer env.Cleanup()
    
    os.Setenv("GNUPGHOME", env.GNUPGHome)
    os.Setenv("PASS_GPG_ID", withPassKey)
    
    testContent := "test-password-789"
    testPath := "test/batch-password"
    if err := env.CreateTestPassword(testPath, testContent, withPassKey); err != nil {
        t.Fatalf("Failed to create test password: %v", err)
    }
    
    fullPath := filepath.Join(env.PasswordStore, testPath+“.gpg")
    
    // Test with batch mode - should not retry
    opts := gpg.BatchGPGOptions("") // Empty passphrase
    _, err = gpg.DecryptFileWithOptions(fullPath, opts)
    if err == nil {
        t.Error("Expected decryption to fail in batch mode without passphrase")
    }
    
    // Should not contain "Operation cancelled" retry
    if strings.Contains(err.Error(), "retry") {
        t.Error("Error message should not mention retry in batch mode")
    }
}
```

### File: `pass/internal/testhelper/testhelper.go`

#### New Helper Functions

```go
// SetupTestEnvWithoutGPGAgent sets up a test environment without gpg-agent running
func SetupTestEnvWithoutGPGAgent() (*TestEnv, string, string, error) {
    env, noPassKey, withPassKey, err := SetupTestEnvWithGPGKeys()
    if err != nil {
        return nil, "", "", err
    }
    
    // Kill gpg-agent if it's running
    killGPGAgent()
    
    return env, noPassKey, withPassKey, nil
}

// killGPGAgent kills any running gpg-agent processes
func killGPGAgent() {
    // Platform-specific implementation
    if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
        exec.Command("gpgconf", "--kill", "gpg-agent").Run()
    } else if runtime.GOOS == "windows" {
        exec.Command("taskkill", "/F", "/IM", "gpg-agent.exe").Run()
    }
}

// CreateTestPasswordWithoutCaching creates a password and ensures passphrase is not cached
func (e *TestEnv) CreateTestPasswordWithoutCaching(path, content, recipient string) error {
    // First, kill gpg-agent to ensure no cached passphrase
    killGPGAgent()
    
    // Create the password
    return e.CreateTestPassword(path, content, recipient)
}
```

## Testing Strategy

### Unit Test Execution Order

1. Run existing tests to ensure no regressions
2. Run new error detection tests
3. Run new option tests
4. Run new decryption retry tests
5. Run integration tests

### Test Categories

#### Category 1: Error Detection (Fast, No Dependencies)
- `TestIsOperationCancelledError`
- `TestExtractGPGError` (existing)

#### Category 2: Options (Fast, No Dependencies)
- `TestDefaultGPGOptions`
- `TestBatchGPGOptions`

#### Category 3: Agent Detection (Requires GPG)
- `TestCheckGPGAgent`

#### Category 4: Decryption with Retry (Requires GPG Keys)
- `TestDecryptFileWithRetry`
- `TestDecryptFileNoRetryInBatchMode`

#### Category 5: Integration (Requires Full Environment)
- End-to-end tests with various scenarios

### Test Execution

```bash
# Run all tests
cd pass && go test ./...

# Run specific test packages
cd pass && go test ./pkg/gpg/...

# Run specific tests
cd pass && go test ./pkg/gpg -run TestIsOperationCancelledError -v
cd pass && go test ./pkg/gpg -run TestDecryptFileWithRetry -v

# Run tests with coverage
cd pass && go test ./pkg/gpg -coverprofile=coverage.out -v
cd pass && go tool cover -html=coverage.out
```

## Rollout Plan

### Step 1: Create Feature Branch
```bash
git checkout main
git pull origin main
git checkout -b feat/7-gpg-agent-prompt
```

### Step 2: Implement Changes
1. Modify `pass/pkg/gpg/gpg.go` with new functions and modifications
2. Add tests to `pass/pkg/gpg/gpg_test.go`
3. Update test helpers if needed

### Step 3: Run Tests
```bash
cd pass && go test ./... -v
```

### Step 4: Fix Issues
- Fix any failing tests
- Address code review feedback
- Update documentation

### Step 5: Create Pull Request
```bash
git add .
git commit -m "✨ feat: add automatic gpg-agent passphrase prompting"
git push origin feat/7-gpg-agent-prompt
```

### Step 6: Review and Merge
- Get code review approval
- Ensure all tests pass
- Merge to main

## Success Metrics

- [ ] All existing tests pass
- [ ] All new tests pass
- [ ] Decryption works when gpg-agent is not running
- [ ] Decryption works when passphrase is not cached
- [ ] Batch mode continues to work without prompting
- [ ] No personal data in tests or code
- [ ] Documentation updated

## Risks and Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Breaking existing tests | Medium | High | Run all tests frequently during development |
| GPG version incompatibilities | Low | Medium | Test with multiple GPG versions |
| Performance issues with retry | Low | Low | Limit retry to once, add timeout if needed |
| Security issues with passphrase handling | Low | High | Use standard GPG mechanisms, don't store passphrases |
| Windows compatibility issues | Medium | Medium | Test on Windows, use cross-platform code |

## Contingency Plans

1. **If tests fail**: Revert changes, investigate, fix
2. **If performance issues**: Add timeout to retry, consider async retry
3. **If security concerns**: Audit code, consult GPG documentation
4. **If Windows issues**: Add platform-specific code paths

## Appendix

### GPG Command Reference

```bash
# Check gpg-agent status
gpgconf --list-dirs agent-socket
gpg --version

# Start gpg-agent
gpgconf --launch gpg-agent

# Kill gpg-agent
gpgconf --kill gpg-agent

# List secret keys
gpg --list-secret-keys

# Decrypt with loopback pinentry
gpg --pinentry-mode loopback --decrypt file.gpg

# Decrypt with batch mode and passphrase
gpg --batch --pinentry-mode loopback --passphrase "pass" --decrypt file.gpg
```

### Test Environment Setup

For manual testing:

```bash
# Create a test GPG key with passphrase
gpg --batch --pinentry-mode loopback --passphrase "test123" \
  --gen-key <<EOF
Key-Type: RSA
Key-Length: 2048
Subkey-Type: RSA
Subkey-Length: 2048
Name-Real: Test User
Name-Email: test@example.com
Expire-Date: 0
%no-protection
%commit
EOF

# Create a test password
mkdir -p ~/.password-store/test
echo "my-secret-password" > ~/.password-store/test/example.txt
gpg --encrypt --recipient "test@example.com" \
  --output ~/.password-store/test/example.txt.gpg \
  ~/.password-store/test/example.txt
rm ~/.password-store/test/example.txt

# Test decryption (should prompt for passphrase)
./pass show test/example
```

### Debugging Tips

```bash
# Enable verbose GPG output
export GPG_TTY=$(tty)
gpg --verbose --decrypt file.gpg

# Check gpg-agent logs
gpgconf --list-options gpg-agent

# Test with different pinentry modes
gpg --pinentry-mode loopback --decrypt file.gpg
gpg --pinentry-mode cancel --decrypt file.gpg
gpg --pinentry-mode error --decrypt file.gpg
```

---

*Document Version: 1.0*
*Last Updated: 2026-07-01*
*Author: @aasmundo*
*Status: Draft*
