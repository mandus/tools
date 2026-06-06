package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewModel tests that a new TUI model can be created
func TestNewModel(t *testing.T) {
	passwords := []string{"email/test", "social/github"}
	model := NewModel(passwords, FuzzyModeShow)
	
	if model == nil {
		t.Fatal("NewModel returned nil")
	}
	
	if model.mode != FuzzyModeShow {
		t.Errorf("Expected mode FuzzyModeShow, got %v", model.mode)
	}
	
	if len(model.allPasswords) != len(passwords) {
		t.Errorf("Expected %d passwords, got %d", len(passwords), len(model.allPasswords))
	}
	
	if model.loading {
		t.Error("Model should not be loading initially")
	}
	
	if model.quitting {
		t.Error("Model should not be quitting initially")
	}
}

// TestModelModes tests all modes
func TestModelModes(t *testing.T) {
	passwords := []string{"test/password"}
	
	modes := []FuzzySearchMode{FuzzyModeShow, FuzzyModeClip, FuzzyModeRm}
	for _, mode := range modes {
		model := NewModel(passwords, mode)
		if model == nil {
			t.Errorf("NewModel returned nil for mode %v", mode)
		}
		if model.mode != mode {
			t.Errorf("Expected mode %v, got %v", mode, model.mode)
		}
	}
}

// TestGetTitle tests the title generation
func TestGetTitle(t *testing.T) {
	tests := []struct {
		mode     FuzzySearchMode
		expected string
	}{
		{FuzzyModeShow, "Select password (Enter to show, Esc to cancel)"},
		{FuzzyModeClip, "Select password to copy (Enter to copy, Esc to cancel)"},
		{FuzzyModeRm, "Select password to remove (Enter to delete, Esc to cancel)"},
	}
	
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := getTitle(tt.mode)
			if result != tt.expected {
				t.Errorf("getTitle(%v) = %q, want %q", tt.mode, result, tt.expected)
			}
		})
	}
}

// TestGetPrompt tests the prompt generation
func TestGetPrompt(t *testing.T) {
	tests := []struct {
		mode     FuzzySearchMode
		expected string
	}{
		{FuzzyModeShow, "Search: "},
		{FuzzyModeClip, "Copy: "},
		{FuzzyModeRm, "Remove: "},
	}
	
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := getPrompt(tt.mode)
			if result != tt.expected {
				t.Errorf("getPrompt(%v) = %q, want %q", tt.mode, result, tt.expected)
			}
		})
	}
}

// TestItem tests the password list item
func TestItem(t *testing.T) {
	item := item{path: "test/password", matchScore: 100, matchIndices: []int{0, 5}}
	
	if item.Title() != "test/password" {
		t.Errorf("Title() = %q, want %q", item.Title(), "test/password")
	}
	
	if item.Description() != "" {
		t.Errorf("Description() = %q, want empty string", item.Description())
	}
	
	if item.FilterValue() != "test/password" {
		t.Errorf("FilterValue() = %q, want %q", item.FilterValue(), "test/password")
	}
}

// TestFilterList tests the filtering functionality
func TestFilterList(t *testing.T) {
	passwords := []string{"email/gmail.com", "social/github.com", "work/vpn"}
	model := NewModel(passwords, FuzzyModeShow)
	
	// Initially should show all passwords
	if model.list.Items() == nil || len(model.list.Items()) != len(passwords) {
		t.Errorf("Expected %d items initially, got %d", len(passwords), len(model.list.Items()))
	}
	
	// Set a query that filters the list
	model.input.SetValue("gmail")
	model.filterList()
	
	// Should only show gmail password
	filtered := model.list.Items()
	if len(filtered) != 1 {
		t.Errorf("Expected 1 filtered item, got %d", len(filtered))
	}
	
	// Check it's the right password
	if len(filtered) > 0 {
		item := filtered[0].(item)
		if !strings.Contains(item.path, "gmail") {
			t.Errorf("Filtered item %q doesn't contain 'gmail'", item.path)
		}
	}
}

// TestHelpView tests the help text
func TestHelpView(t *testing.T) {
	help := helpView()
	if help == "" {
		t.Error("Help view should not be empty")
	}
	
	// Should contain navigation keys
	if !strings.Contains(help, "↑") || !strings.Contains(help, "↓") {
		t.Error("Help should contain arrow key navigation")
	}
	
	// Should contain action keys
	if !strings.Contains(help, "Enter") || !strings.Contains(help, "Esc") {
		t.Error("Help should contain action keys")
	}
}

// TestTruncate tests the truncate function
func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello", 5, "hello"},
		{"hello world", 5, "he..."},
		{"hello", 3, "hel"},
		{"hello", 2, "he"},
		{"hello", 1, "h"},
		{"hello", 0, ""},
		{"", 5, ""},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := truncate(tt.input, tt.length)
			if result != tt.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.length, result, tt.expected)
			}
		})
	}
}

// TestPasswordDelegate tests the custom delegate
func TestPasswordDelegate(t *testing.T) {
	delegate := NewPasswordDelegate()
	
	if delegate.Height() != 1 {
		t.Errorf("Height() = %d, want 1", delegate.Height())
	}
	
	if delegate.Spacing() != 0 {
		t.Errorf("Spacing() = %d, want 0", delegate.Spacing())
	}
}

// TestCollectAllPasswords tests password collection
func TestCollectAllPasswords(t *testing.T) {
	// Create a temp password store
	tempDir, err := os.MkdirTemp("", "pass-collect-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create some test password files
	passwords := []string{"email/test", "social/github"}
	for _, p := range passwords {
		fullPath := filepath.Join(tempDir, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath+".gpg", []byte("dummy"), 0600); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}
	
	// Collect passwords
	got, err := CollectAllPasswords(tempDir)
	if err != nil {
		t.Fatalf("CollectAllPasswords failed: %v", err)
	}
	
	// Check that we got all passwords
	if len(got) != len(passwords) {
		t.Errorf("Expected %d passwords, got %d", len(passwords), len(got))
	}
	
	// Check that all expected passwords are present
	for _, want := range passwords {
		found := false
		for _, got := range got {
			if got == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Missing expected password: %s", want)
		}
	}
}
