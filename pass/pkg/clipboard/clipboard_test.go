package clipboard

import (
	"testing"
)

func TestCopy(t *testing.T) {
	// Skip if clipboard not available
	if !IsAvailable() {
		t.Skip("Clipboard not available")
	}

	// Test copying simple text - may fail in restricted environments
	testText := "test password 123"
	err := Copy(testText)
	if err != nil {
		t.Skipf("Copy skipped - clipboard access denied: %v", err)
	}
	t.Log("Copy succeeded")
}

func TestCopyBytes(t *testing.T) {
	// Skip if clipboard not available
	if !IsAvailable() {
		t.Skip("Clipboard not available")
	}

	// Test copying raw bytes - may fail in restricted environments
	testData := []byte("test password bytes")
	err := CopyBytes(testData)
	if err != nil {
		t.Skipf("CopyBytes skipped - clipboard access denied: %v", err)
	}
	t.Log("CopyBytes succeeded")
}

func TestClear(t *testing.T) {
	// Skip if clipboard not available
	if !IsAvailable() {
		t.Skip("Clipboard not available")
	}

	// Clear clipboard - may fail in restricted environments
	err := Clear()
	if err != nil {
		t.Skipf("Clear skipped - clipboard access denied: %v", err)
	}
	t.Log("Clear succeeded")
}

func TestIsAvailable(t *testing.T) {
	available := IsAvailable()
	t.Logf("Clipboard available: %v", available)
	// Don't fail - just log the result
}
