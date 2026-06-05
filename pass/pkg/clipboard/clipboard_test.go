package clipboard

import (
	"testing"
)

func TestCopy(t *testing.T) {
	// Skip if clipboard not available
	if !IsAvailable() {
		t.Skip("Clipboard not available")
	}

	// Test copying simple text
	testText := "test password 123"
	if err := Copy(testText); err != nil {
		t.Fatalf("Copy failed: %v", err)
	}
	t.Log("Copy succeeded")
}

func TestCopyBytes(t *testing.T) {
	// Skip if clipboard not available
	if !IsAvailable() {
		t.Skip("Clipboard not available")
	}

	// Test copying raw bytes
	testData := []byte("test password bytes")
	if err := CopyBytes(testData); err != nil {
		t.Fatalf("CopyBytes failed: %v", err)
	}
	t.Log("CopyBytes succeeded")
}

func TestClear(t *testing.T) {
	// Skip if clipboard not available
	if !IsAvailable() {
		t.Skip("Clipboard not available")
	}

	// Clear clipboard
	if err := Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}
	t.Log("Clear succeeded")
}

func TestIsAvailable(t *testing.T) {
	available := IsAvailable()
	t.Logf("Clipboard available: %v", available)
	// Don't fail - just log the result
}
