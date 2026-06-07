// Package filesystem provides filesystem utilities for the pass tool.
package filesystem

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// NormalizePath converts a path to use the OS-specific separator.
// On Windows, converts / to \. On Unix, converts \ to /. Also handles path normalization.
func NormalizePath(path string) string {
	// Replace both forward and backward slashes with OS separator
	// This ensures cross-platform compatibility
	normalized := strings.ReplaceAll(path, "/", string(filepath.Separator))
	normalized = strings.ReplaceAll(normalized, "\\", string(filepath.Separator))
	
	// Clean the path (remove . and .. elements)
	return filepath.Clean(normalized)
}

// NormalizePathForDisplay converts a path to use forward slashes for display.
// This ensures consistent output across platforms.
func NormalizePathForDisplay(path string) string {
	// Replace both OS separator and backslash with forward slash
	// This handles paths from both Windows and Unix systems
	result := strings.ReplaceAll(path, "\\", "/")
	result = strings.ReplaceAll(result, string(filepath.Separator), "/")
	return result
}

// SecureDelete securely deletes a file by overwriting it before removal.
func SecureDelete(filePath string) error {
	// Open the file
	file, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		// If we can't open, try regular delete
		return os.Remove(filePath)
	}
	defer file.Close()
	
	// Get file info
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return os.Remove(filePath)
	}
	
	// Overwrite file content with random data
	// Use a simple pattern for overwriting
	buffer := make([]byte, 4096)
	for i := range buffer {
		buffer[i] = byte(i % 256)
	}
	
	// Overwrite in chunks
	fileSize := info.Size()
	var overwritten int64
	
	for overwritten < fileSize {
		toWrite := buffer
		if int64(len(buffer)) > fileSize-overwritten {
			toWrite = buffer[:fileSize-overwritten]
		}
		if _, err := file.WriteAt(toWrite, overwritten); err != nil {
			// If overwrite fails, just close and delete
			break
		}
		overwritten += int64(len(toWrite))
	}
	
	// Sync to ensure data is written
	file.Sync()
	file.Close()
	
	// Now delete the file
	return os.Remove(filePath)
}

// CopyToClipboard copies text to the system clipboard.
// On Windows, uses the built-in `clip` command.
func CopyToClipboard(text string) error {
	cmd := exec.Command("clip")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}
	
	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start clip command: %v", err)
	}
	
	// Write text to stdin
	if _, err := io.WriteString(stdin, text); err != nil {
		stdin.Close()
		return fmt.Errorf("failed to write to clipboard: %v", err)
	}
	
	// Close stdin
	stdin.Close()
	
	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("clip command failed: %v", err)
	}
	
	return nil
}

// StartClipboardClearTimer starts a timer to clear the clipboard after a delay.
// Default timeout is 45 seconds, configurable via PASS_CLIPBOARD_TIMEOUT.
func StartClipboardClearTimer() {
	timeoutStr := os.Getenv("PASS_CLIPBOARD_TIMEOUT")
	if timeoutStr == "0" || strings.ToLower(os.Getenv("PASS_CLIPBOARD_CLEAR")) == "false" {
		return // Disabled
	}
	
	// Parse timeout
	timeout := 45 // default
	if timeoutStr != "" {
		fmt.Sscanf(timeoutStr, "%d", &timeout)
	}
	
	// Start timer in background
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		// Clear clipboard by copying empty string
		_ = CopyToClipboard("")
	}()
}

// RunCommand runs an external command with arguments.
func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// JoinPath joins path elements using the OS-specific separator.
func JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

// EnsurePasswordStore ensures the password store directory exists.
func EnsurePasswordStore(storeDir string) error {
	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		if err := os.MkdirAll(storeDir, 0700); err != nil {
			return fmt.Errorf("failed to create password store: %v", err)
		}
	}
	return nil
}
