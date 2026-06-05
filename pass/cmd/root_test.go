package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestGetPasswordStoreDir(t *testing.T) {
	// Save original env
	orig := os.Getenv("PASSWORD_STORE_DIR")
	defer os.Setenv("PASSWORD_STORE_DIR", orig)

	// Test default
	os.Unsetenv("PASSWORD_STORE_DIR")
	result := GetPasswordStoreDir()
	if result == "" {
		t.Error("GetPasswordStoreDir returned empty string")
	}
	t.Logf("Default store dir: %s", result)

	// Test with custom env
	os.Setenv("PASSWORD_STORE_DIR", "/custom/store")
	result = GetPasswordStoreDir()
	if result != "/custom/store" {
		t.Errorf("GetPasswordStoreDir = %q, want %q", result, "/custom/store")
	}
}

func TestIsClipboardFlagSet(t *testing.T) {
	// Save original state
	orig := clipFlag
	defer func() { clipFlag = orig }()

	// Test default (false)
	clipFlag = false
	if IsClipboardFlagSet() {
		t.Error("IsClipboardFlagSet should return false by default")
	}

	// Test when set
	clipFlag = true
	if !IsClipboardFlagSet() {
		t.Error("IsClipboardFlagSet should return true when flag is set")
	}
}

func TestRootCommandHelp(t *testing.T) {
	// Test that help doesn't error
	rootCmd.SetArgs([]string{"--help"})
	var buf bytes.Buffer
	rootCmd.SetOutput(&buf)
	
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Help command failed: %v", err)
	}
	
	if buf.Len() == 0 {
		t.Error("Help output is empty")
	}
	t.Log("Help command works")
}

func TestRootCommandVersion(t *testing.T) {
	// Test that version doesn't error
	rootCmd.SetArgs([]string{"--version"})
	var buf bytes.Buffer
	rootCmd.SetOutput(&buf)
	
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Version command failed: %v", err)
	}
	
	if buf.Len() == 0 {
		t.Error("Version output is empty")
	}
	t.Log("Version command works")
}
