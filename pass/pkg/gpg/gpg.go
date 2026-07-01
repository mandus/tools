// Package gpg provides GPG encryption and decryption functionality.
package gpg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GPG error message patterns for detecting operation cancelled errors
const (
	ErrOperationCancelled = "Operation cancelled"
	ErrGPGCancelled       = "gpg: cancelled"
	ErrNoPinentry         = "gpg-agent: no pinentry"
	ErrCannotOpenTTY      = "cannot open '/dev/tty'"
	ErrNoTTY              = "No such device or address"
)

// GPGOptions contains options for GPG operations
// This allows tests to specify batch mode, passphrase, etc.
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

// EncryptFile encrypts a file using GPG and saves it to the destination path.
// Uses the default GPG key or the one specified in PASS_GPG_ID environment variable.
func EncryptFile(srcPath, destPath string) error {
	return EncryptFileWithOptions(srcPath, destPath, DefaultGPGOptions())
}

// EncryptFileWithOptions encrypts a file using GPG with custom options.
// This allows for batch mode, specific recipients, and passphrase handling.
func EncryptFileWithOptions(srcPath, destPath string, opts GPGOptions) error {
	// Get recipient from options, environment, or use default
	recipient := opts.Recipient
	if recipient == "" {
		recipient = os.Getenv("PASS_GPG_ID")
	}
	
	args := []string{
		"--encrypt",
		"--armor", // ASCII-armored output for compatibility
		"--yes",   // Assume yes to prompts
	}
	
	// Add batch mode if requested
	if opts.BatchMode {
		args = append(args, "--batch")
	}
	
	// Add pinentry mode if specified
	if opts.PinentryMode != "" {
		args = append(args, "--pinentry-mode", opts.PinentryMode)
	}
	
	if recipient != "" {
		args = append(args, "--recipient", recipient)
	} else {
		// Use default key (no explicit recipient)
		args = append(args, "--default-recipient-self")
	}
	
	// Add output file
	args = append(args, "--output", destPath)
	
	// Add source file
	args = append(args, srcPath)
	
	// Execute GPG
	cmd := exec.Command("gpg", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gpg encryption failed: %v", err)
	}
	
	return nil
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
	
	// Extract stderr from the error message if it contains it
	stderrStr := extractStderrFromError(err)
	
	// Check if we should retry - check both the error and stderr content
	if opts.RetryOnCancel && opts.AllowPrompt && shouldRetryDecryption(err, stderrStr) {
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

// decryptFileResult holds the result of a decryption attempt
type decryptFileResult struct {
	content    string
	err       error
	stderrStr  string
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

// shouldRetryDecryption checks if we should retry decryption with loopback pinentry-mode
// It examines both the error and stderr content to determine if a retry is appropriate
func shouldRetryDecryption(err error, stderrStr string) bool {
	if err == nil {
		return false
	}
	
	// Check the error message for known retryable conditions
	errStr := err.Error()
	if strings.Contains(errStr, ErrOperationCancelled) ||
		strings.Contains(errStr, ErrGPGCancelled) ||
		strings.Contains(errStr, ErrNoPinentry) ||
		strings.Contains(errStr, ErrCannotOpenTTY) ||
		strings.Contains(errStr, ErrNoTTY) {
		return true
	}
	
	// Also check the stderr content for retryable conditions
	if strings.Contains(stderrStr, ErrOperationCancelled) ||
		strings.Contains(stderrStr, ErrGPGCancelled) ||
		strings.Contains(stderrStr, ErrNoPinentry) ||
		strings.Contains(stderrStr, ErrCannotOpenTTY) ||
		strings.Contains(stderrStr, ErrNoTTY) ||
		strings.Contains(stderrStr, "cannot open") ||
		strings.Contains(stderrStr, "No such device or address") {
		return true
	}
	
	// Also check for exit status 2 which often indicates TTY/pinentry issues
	if strings.Contains(errStr, "exit status 2") {
		return true
	}
	
	return false
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

// extractStderrFromError extracts stderr content from an error message
func extractStderrFromError(err error) string {
	if err == nil {
		return ""
	}
	
	errStr := err.Error()
	// Look for "(stderr: " in the error message
	if idx := strings.Index(errStr, "(stderr: "); idx != -1 {
		// Extract everything after "(stderr: " until the end
		stderrStart := idx + len("(stderr: ")
		if stderrStart < len(errStr) {
			// Remove trailing parenthesis if present
			stderrContent := errStr[stderrStart:]
			if strings.HasSuffix(stderrContent, ")") {
				stderrContent = strings.TrimSuffix(stderrContent, ")")
			}
			return stderrContent
		}
	}
	return ""
}

// extractGPGError extracts a clean error message from GPG stderr output
func extractGPGError(stderr string) string {
	// Split by newlines and find the most relevant error line
	lines := strings.Split(stderr, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" && !strings.HasPrefix(line, "gpg:") {
			continue
		}
		if strings.Contains(line, "decryption failed") ||
			strings.Contains(line, "No secret key") ||
			strings.Contains(line, "bad passphrase") ||
			strings.Contains(line, "not a detached signature") ||
			strings.Contains(line, ErrOperationCancelled) ||
			strings.Contains(line, ErrCannotOpenTTY) ||
			strings.Contains(line, ErrNoTTY) {
			return line
		}
	}
	return stderr
}

// CheckGPG checks if GPG is installed and available.
func CheckGPG() error {
	cmd := exec.Command("gpg", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gpg: command not found. Please install GPG4Win or GnuPG")
	}
	return nil
}

// CheckGPGBatch checks if GPG is available and can run in batch mode.
// This is useful for tests to verify the environment is set up correctly.
func CheckGPGBatch() error {
	// Check basic GPG
	if err := CheckGPG(); err != nil {
		return err
	}
	
	// Check batch mode works
	cmd := exec.Command("gpg", "--batch", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gpg: batch mode not available: %v", err)
	}
	return nil
}

// HasSecretKey checks if there is at least one secret key available for decryption.
func HasSecretKey() bool {
	cmd := exec.Command("gpg", "--list-secret-keys")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = nil
	
	if err := cmd.Run(); err != nil {
		return false
	}
	
	// If there's any output, we have at least one secret key
	return stdout.Len() > 0
}

// GetDefaultRecipient returns the default GPG recipient key ID.
// Returns empty string if using default key.
func GetDefaultRecipient() string {
	// Check PASS_GPG_ID environment variable first
	if recipient := os.Getenv("PASS_GPG_ID"); recipient != "" {
		return recipient
	}
	
	// Try to get the first secret key
	cmd := exec.Command("gpg", "--list-secret-keys", "--with-colons")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	
	if err := cmd.Run(); err != nil {
		return ""
	}
	
	// Parse the output to find the first key
	lines := strings.Split(stdout.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "sec:") {
			// Format: sec:u:2048:1:4DB683CED8BB579C:1620000000:1700000000::::::::scESC:
			parts := strings.Split(line, ":")
			if len(parts) >= 5 {
				return parts[4] // Key ID
			}
		}
	}
	
	return ""
}
