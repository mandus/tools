// Package config provides configuration management for the pass tool.
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration.
type Config struct {
	PasswordStoreDir string
	GPGID           string
	ClipboardTimeout int
	ClipboardClear   bool
	GitName         string
	GitEmail        string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() *Config {
	return &Config{
		PasswordStoreDir: getEnvWithDefault("PASSWORD_STORE_DIR", getDefaultStoreDir()),
		GPGID:           os.Getenv("PASS_GPG_ID"),
		ClipboardTimeout: getEnvAsIntWithDefault("PASS_CLIPBOARD_TIMEOUT", 45),
		ClipboardClear:   getEnvAsBoolWithDefault("PASS_CLIPBOARD_CLEAR", true),
		GitName:         getEnvWithDefault("PASS_GIT_NAME", ""),
		GitEmail:        getEnvWithDefault("PASS_GIT_EMAIL", ""),
	}
}

// getDefaultStoreDir returns the default password store directory.
func getDefaultStoreDir() string {
	// Use USERPROFILE on Windows, HOME on Unix
	if home := os.Getenv("USERPROFILE"); home != "" {
		return filepath.Join(home, ".password-store")
	}
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".password-store")
	}
	return ".password-store"
}

// getEnvWithDefault returns the value of an environment variable or a default.
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsIntWithDefault returns the value of an environment variable as int or a default.
func getEnvAsIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

// getEnvAsBoolWithDefault returns the value of an environment variable as bool or a default.
func getEnvAsBoolWithDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}
