// Package gpg provides GPG encryption and decryption functionality.
package gpg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GPGOptions contains options for GPG operations
// This allows tests to specify batch mode, passphrase, etc.
type GPGOptions struct {
	BatchMode      bool
	Passphrase     string
	Recipient      string
	PinentryMode   string // "loopback" to bypass pinentry
}

// DefaultGPGOptions returns the default GPG options
func DefaultGPGOptions() GPGOptions {
	return GPGOptions{
		BatchMode:    false,
		PinentryMode: "",
	}
}

// BatchGPGOptions returns options suitable for batch/non-interactive operations
func BatchGPGOptions(passphrase string) GPGOptions {
	return GPGOptions{
		BatchMode:      true,
		Passphrase:     passphrase,
		PinentryMode:   "loopback",
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
// GPG will prompt for passphrase if needed (handled by gpg-agent).
func DecryptFile(filePath string) (string, error) {
	return DecryptFileWithOptions(filePath, DefaultGPGOptions())
}

// DecryptFileWithOptions decrypts a GPG file with custom options.
// This allows for batch mode, passphrase specification, etc.
func DecryptFileWithOptions(filePath string, opts GPGOptions) (string, error) {
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
			strings.Contains(line, "not a detached signature") {
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
