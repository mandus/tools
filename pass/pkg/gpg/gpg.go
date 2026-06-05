// Package gpg provides GPG encryption and decryption functionality.
package gpg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// EncryptFile encrypts a file using GPG and saves it to the destination path.
// Uses the default GPG key or the one specified in PASS_GPG_ID environment variable.
func EncryptFile(srcPath, destPath string) error {
	// Get recipient from environment or use default
	recipient := os.Getenv("PASS_GPG_ID")
	
	args := []string{
		"--encrypt",
		"--armor", // ASCII-armored output for compatibility
		"--yes",   // Assume yes to prompts
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
	cmd := exec.Command("gpg", "--decrypt", filePath)
	
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		// Check if it's a decryption error
		if strings.Contains(stderr.String(), "decryption failed") ||
			strings.Contains(stderr.String(), "No secret key") {
			return "", fmt.Errorf("GPG decryption failed: %s", stderr.String())
		}
		return "", fmt.Errorf("gpg decryption failed: %v", err)
	}
	
	// Trim trailing newline if present
	output := strings.TrimSuffix(stdout.String(), "\n")
	return output, nil
}

// CheckGPG checks if GPG is installed and available.
func CheckGPG() error {
	cmd := exec.Command("gpg", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gpg: command not found. Please install GPG4Win or GnuPG")
	}
	return nil
}
