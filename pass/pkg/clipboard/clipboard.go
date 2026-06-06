// Package clipboard provides clipboard functionality for the pass tool.
// On Windows, it uses the built-in `clip` command.
// On other platforms, it may use different mechanisms.
package clipboard

import (
	"bytes"
	"os/exec"
)

// Copy copies text to the system clipboard.
// Uses platform-specific mechanisms.
func Copy(text string) error {
	// On Windows, use the clip command
	cmd := exec.Command("clip")
	
	// Write text to stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	
	if err := cmd.Start(); err != nil {
		return err
	}
	
	// Write the text
	if _, err := stdin.Write([]byte(text)); err != nil {
		stdin.Close()
		return err
	}
	
	stdin.Close()
	return cmd.Wait()
}

// CopyBytes copies raw bytes to the clipboard.
func CopyBytes(data []byte) error {
	return Copy(string(data))
}

// Clear clears the clipboard by copying an empty string.
func Clear() error {
	return Copy("")
}

// IsAvailable checks if clipboard functionality is available.
func IsAvailable() bool {
	// Try to find the clip command
	cmd := exec.Command("where", "clip")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// Read reads text from the clipboard.
// Note: This may not work on all platforms.
func Read() (string, error) {
	// Windows doesn't have a built-in way to read from clipboard via command line
	// This is a placeholder for future implementation
	// For now, we'll use a workaround if available
	
	// Try using powershell if available
	cmd := exec.Command("powershell", "-Command", "Get-Clipboard")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		return "", err
	}
	
	return out.String(), nil
}
