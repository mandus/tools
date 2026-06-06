package terminal

import (
	"strings"
	"testing"
)

func TestColorizeNoANSI(t *testing.T) {
	// Test that Colorize returns the original text when ANSI is not supported
	// We can't easily mock SupportsANSI, so we test the logic indirectly
	result := Colorize("hello", ColorRed)
	// In non-ANSI mode, should return the text unchanged
	// In ANSI mode, will have color codes
	if !SupportsANSI() {
		if result != "hello" {
			t.Errorf("Colorize in non-ANSI mode: got %q, want %q", result, "hello")
		}
	}
}

func TestHighlightMatch(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		matchIndices []int
		want         string
	}{
		{
			name:         "no matches",
			path:         "email/gmail.com/user",
			matchIndices: []int{},
			want:         "email/gmail.com/user",
		},
		{
			name:         "empty path",
			path:         "",
			matchIndices: []int{0, 1},
			want:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HighlightMatch(tt.path, tt.matchIndices)
			if got != tt.want {
				t.Errorf("HighlightMatch(%q, %v) = %q, want %q", tt.path, tt.matchIndices, got, tt.want)
			}
		})
	}
}

func TestRepeatString(t *testing.T) {
	tests := []struct {
		s      string
		n      int
		want   string
	}{
		{"a", 3, "aaa"},
		{"ab", 2, "abab"},
		{"x", 0, ""},
		{"y", 1, "y"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := RepeatString(tt.s, tt.n)
			if got != tt.want {
				t.Errorf("RepeatString(%q, %d) = %q, want %q", tt.s, tt.n, got, tt.want)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		s     string
		width int
		want  string
	}{
		{"hello", 10, "hello     "},
		{"hello", 5, "hello"},
		{"hello", 3, "hello"},
		{"", 5, "     "},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := PadRight(tt.s, tt.width)
			if got != tt.want {
				t.Errorf("PadRight(%q, %d) = %q, want %q", tt.s, tt.width, got, tt.want)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	tests := []struct {
		s     string
		width int
		want  string
	}{
		{"hello", 10, "     hello"},
		{"hello", 5, "hello"},
		{"hello", 3, "hello"},
		{"", 5, "     "},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := PadLeft(tt.s, tt.width)
			if got != tt.want {
				t.Errorf("PadLeft(%q, %d) = %q, want %q", tt.s, tt.width, got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		s      string
		length int
		want   string
	}{
		{"hello", 10, "hello"},
		{"hello", 5, "hello"},
		{"hello", 4, "h..."},    // length 4: 1 char + "..." = 4
		{"hello", 3, "..."},   // length 3: just "..."
		{"hello", 2, ".."},    // length 2: ".."
		{"hello", 1, "."},     // length 1: "."
		{"hello", 0, ""},
		{"", 5, ""},
		{"abc", 3, "abc"},
		{"abc", 2, ".."},
		{"abc", 1, "."},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := Truncate(tt.s, tt.length)
			if got != tt.want {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tt.s, tt.length, got, tt.want)
			}
		})
	}
}

func TestKeyString(t *testing.T) {
	tests := []struct {
		key  Key
		want string
	}{
		{Key{Rune: 'a'}, "a"},
		{Key{IsArrow: true, ArrowDir: "up"}, "Arrow(up)"},
		{Key{IsCtrl: true, CtrlChar: 'A'}, "Ctrl+A"},
		{Key{IsEscape: true}, "Escape"},
		{Key{IsEnter: true}, "Enter"},
		{Key{IsBackspace: true}, "Backspace"},
		{Key{IsDelete: true}, "Delete"},
		{Key{IsTab: true}, "Tab"},
		{Key{IsHome: true}, "Home"},
		{Key{IsEnd: true}, "End"},
		{Key{IsPageUp: true}, "PageUp"},
		{Key{IsPageDown: true}, "PageDown"},
		{Key{}, "Unknown"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := tt.key.String()
			if got != tt.want {
				t.Errorf("Key.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetMatchHighlight(t *testing.T) {
	matchIndices := []int{0, 2, 4}
	highlight := GetMatchHighlight(matchIndices)

	result := highlight("abcdef")
	// Just verify it returns a non-empty string
	if result == "" {
		t.Error("GetMatchHighlight should return a non-empty string")
	}
	// In non-ANSI mode, should return original if no matches
	// But we have match indices, so it should still process
	if !SupportsANSI() && !strings.Contains(result, "abcdef") {
		t.Errorf("Expected result to contain original text, got %q", result)
	}
}

func TestSupportsANSI(t *testing.T) {
	// Just test that it doesn't panic
	supports := SupportsANSI()
	t.Logf("Terminal supports ANSI: %v", supports)
	// This is environment-dependent, so we just check it returns a bool
	_ = supports
}
