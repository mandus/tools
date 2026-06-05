package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Save original env
	origEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, e := range origEnv {
			if idx := strings.Index(e, "="); idx >= 0 {
				os.Setenv(e[:idx], e[idx+1:])
			}
		}
	}()

	// Test default config
	cfg := LoadConfig()
	if cfg.PasswordStoreDir == "" {
		t.Error("PasswordStoreDir should not be empty")
	}
	if cfg.GPGID != "" {
		t.Errorf("GPGID should be empty by default, got %q", cfg.GPGID)
	}
	if cfg.ClipboardTimeout != 45 {
		t.Errorf("ClipboardTimeout should default to 45, got %d", cfg.ClipboardTimeout)
	}
	if !cfg.ClipboardClear {
		t.Error("ClipboardClear should default to true")
	}
}

func TestLoadConfigWithEnv(t *testing.T) {
	os.Setenv("PASSWORD_STORE_DIR", "/custom/store")
	os.Setenv("PASS_GPG_ID", "ABCD1234")
	os.Setenv("PASS_CLIPBOARD_TIMEOUT", "60")
	os.Setenv("PASS_CLIPBOARD_CLEAR", "false")
	os.Setenv("PASS_GIT_NAME", "Test User")
	os.Setenv("PASS_GIT_EMAIL", "test@example.com")

	cfg := LoadConfig()
	
	if cfg.PasswordStoreDir != "/custom/store" {
		t.Errorf("PasswordStoreDir = %q, want %q", cfg.PasswordStoreDir, "/custom/store")
	}
	if cfg.GPGID != "ABCD1234" {
		t.Errorf("GPGID = %q, want %q", cfg.GPGID, "ABCD1234")
	}
	if cfg.ClipboardTimeout != 60 {
		t.Errorf("ClipboardTimeout = %d, want 60", cfg.ClipboardTimeout)
	}
	if cfg.ClipboardClear {
		t.Error("ClipboardClear should be false")
	}
	if cfg.GitName != "Test User" {
		t.Errorf("GitName = %q, want %q", cfg.GitName, "Test User")
	}
	if cfg.GitEmail != "test@example.com" {
		t.Errorf("GitEmail = %q, want %q", cfg.GitEmail, "test@example.com")
	}
}

func TestGetDefaultStoreDir(t *testing.T) {
	// Save original env
	origUSERPROFILE := os.Getenv("USERPROFILE")
	origHOME := os.Getenv("HOME")
	defer func() {
		if origUSERPROFILE != "" {
			os.Setenv("USERPROFILE", origUSERPROFILE)
		} else {
			os.Unsetenv("USERPROFILE")
		}
		if origHOME != "" {
			os.Setenv("HOME", origHOME)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	// Test with USERPROFILE (Windows)
	os.Setenv("USERPROFILE", "C:\\Users\\test")
	os.Unsetenv("HOME")
	result := getDefaultStoreDir()
	expected := "C:\\Users\\test\\.password-store"
	if result != expected {
		t.Errorf("getDefaultStoreDir() = %q, want %q", result, expected)
	}

	// Test with HOME (Unix) - skip on Windows since HOME isn't used
	if os.PathSeparator == '/' {
		os.Unsetenv("USERPROFILE")
		os.Setenv("HOME", "/home/test")
		result = getDefaultStoreDir()
		expected := "/home/test/.password-store"
		if result != expected {
			t.Errorf("getDefaultStoreDir() = %q, want %q", result, expected)
		}
	}

	// Test fallback
	os.Unsetenv("USERPROFILE")
	os.Unsetenv("HOME")
	result = getDefaultStoreDir()
	expected = ".password-store"
	if result != expected {
		t.Errorf("getDefaultStoreDir() = %q, want %q", result, expected)
	}
}

func TestGetEnvWithDefault(t *testing.T) {
	os.Setenv("TEST_VAR", "value")
	defer os.Unsetenv("TEST_VAR")

	if result := getEnvWithDefault("TEST_VAR", "default"); result != "value" {
		t.Errorf("getEnvWithDefault = %q, want %q", result, "value")
	}
	if result := getEnvWithDefault("NONEXISTENT_VAR", "default"); result != "default" {
		t.Errorf("getEnvWithDefault = %q, want %q", result, "default")
	}
}

func TestGetEnvAsIntWithDefault(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	if result := getEnvAsIntWithDefault("TEST_INT", 10); result != 42 {
		t.Errorf("getEnvAsIntWithDefault = %d, want 42", result)
	}
	if result := getEnvAsIntWithDefault("NONEXISTENT_INT", 10); result != 10 {
		t.Errorf("getEnvAsIntWithDefault = %d, want 10", result)
	}
	if result := getEnvAsIntWithDefault("INVALID_INT", 10); result != 10 {
		t.Errorf("getEnvAsIntWithDefault = %d, want 10", result)
	}
}

func TestGetEnvAsBoolWithDefault(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"true", "true", true},
		{"1", "1", true},
		{"yes", "yes", true},
		{"false", "false", false},
		{"0", "0", false},
		{"no", "no", false},
		// Skip empty test - when env var doesn't exist, should return default
		// {"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv("TEST_BOOL", tt.value)
				defer os.Unsetenv("TEST_BOOL")
			}
			result := getEnvAsBoolWithDefault("TEST_BOOL", !tt.expected)
			if result != tt.expected {
				t.Errorf("getEnvAsBoolWithDefault = %v, want %v", result, tt.expected)
			}
		})
	}
}
