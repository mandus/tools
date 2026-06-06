// Package terminal provides terminal UI utilities for the pass tool.
// It handles ANSI escape codes, cursor control, key reading, and terminal size detection.
package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"golang.org/x/term"
)

// Key represents a key press event with information about the key.
type Key struct {
	Rune      rune  // The rune value of the key
	IsArrow   bool  // Whether this is an arrow key
	ArrowDir  string // Direction for arrow keys: "up", "down", "left", "right"
	IsCtrl    bool  // Whether Ctrl was pressed
	CtrlChar  rune  // The control character (e.g., 'A' for Ctrl+A)
	IsEscape  bool  // Whether this is the Escape key
	IsEnter   bool  // Whether this is Enter/Return
	IsBackspace bool // Whether this is Backspace
	IsDelete  bool  // Whether this is Delete
	IsTab     bool  // Whether this is Tab
	IsHome    bool  // Whether this is Home
	IsEnd     bool  // Whether this is End
	IsPageUp  bool  // Whether this is Page Up
	IsPageDown bool // Whether this is Page Down
}

// String returns a string representation of the key for debugging.
func (k Key) String() string {
	if k.IsArrow {
		return fmt.Sprintf("Arrow(%s)", k.ArrowDir)
	}
	if k.IsCtrl {
		return fmt.Sprintf("Ctrl+%c", k.CtrlChar)
	}
	if k.IsEscape {
		return "Escape"
	}
	if k.IsEnter {
		return "Enter"
	}
	if k.IsBackspace {
		return "Backspace"
	}
	if k.IsDelete {
		return "Delete"
	}
	if k.IsTab {
		return "Tab"
	}
	if k.IsHome {
		return "Home"
	}
	if k.IsEnd {
		return "End"
	}
	if k.IsPageUp {
		return "PageUp"
	}
	if k.IsPageDown {
		return "PageDown"
	}
	if k.Rune > 0 {
		return string(k.Rune)
	}
	return "Unknown"
}

// SupportsANSI checks if the terminal supports ANSI escape codes.
func SupportsANSI() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// ANSI escape code constants
const (
	// Clear screen and move cursor to home position
	ClearScreen = "\033[2J\033[H"

	// Hide cursor
	HideCursor = "\033[?25l"

	// Show cursor
	ShowCursor = "\033[?25h"

	// Move cursor up by n lines
	CursorUp = "\033[%dA"

	// Move cursor down by n lines
	CursorDown = "\033[%dB"

	// Move cursor right by n columns
	CursorRight = "\033[%dC"

	// Move cursor left by n columns
	CursorLeft = "\033[%dD"

	// Move cursor to specific position (1-indexed)
	CursorTo = "\033[%d;%dH"

	// Save cursor position
	SaveCursor = "\033[s"

	// Restore cursor position
	RestoreCursor = "\033[u"

	// Clear from cursor to end of line
	ClearToEOL = "\033[K"

	// Clear from cursor to end of screen
	ClearToEOS = "\033[J"

	// Reset all attributes
	Reset = "\033[0m"

	// Text colors (foreground)
	ColorBlack   = "\033[30m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"

	// Bright text colors
	ColorBrightBlack   = "\033[90m"
	ColorBrightRed     = "\033[91m"
	ColorBrightGreen   = "\033[92m"
	ColorBrightYellow  = "\033[93m"
	ColorBrightBlue    = "\033[94m"
	ColorBrightMagenta = "\033[95m"
	ColorBrightCyan    = "\033[96m"
	ColorBrightWhite   = "\033[97m"

	// Background colors
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"

	// Text attributes
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
	Reverse   = "\033[7m"
	Hidden    = "\033[8m"
)

// ClearScreenFunc clears the terminal screen.
func ClearScreenFunc() {
	if SupportsANSI() {
		fmt.Print(ClearScreen)
	}
}

// MoveCursorFunc moves the cursor to the specified row and column (1-indexed).
func MoveCursorFunc(row, col int) {
	if SupportsANSI() {
		fmt.Printf(CursorTo, row, col)
	}
}

// HideCursorFunc hides the terminal cursor.
func HideCursorFunc() {
	if SupportsANSI() {
		fmt.Print(HideCursor)
	}
}

// ShowCursorFunc shows the terminal cursor.
func ShowCursorFunc() {
	if SupportsANSI() {
		fmt.Print(ShowCursor)
	}
}

// GetSize returns the terminal width and height.
func GetSize() (int, int, error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

// PrintAt prints text at the specified position.
func PrintAt(text string, row, col int) {
	if SupportsANSI() {
		fmt.Printf("\033[%d;%dH%s", row, col, text)
	} else {
		fmt.Print(text)
	}
}

// KeyReader provides buffered key reading with support for special keys.
type KeyReader struct {
	reader   *bufio.Reader
	rawMode  bool
	oldState *term.State
	mu       sync.Mutex
}

// NewKeyReader creates a new KeyReader.
func NewKeyReader() (*KeyReader, error) {
	kr := &KeyReader{
		reader: bufio.NewReader(os.Stdin),
	}
	return kr, nil
}

// ReadKey reads a single key press, handling special keys like arrows.
func (kr *KeyReader) ReadKey() (Key, error) {
	kr.mu.Lock()
	defer kr.mu.Unlock()

	// Read a single byte
	b, err := kr.reader.ReadByte()
	if err != nil {
		return Key{}, err
	}

	// Handle Ctrl+C (interrupt)
	if b == 3 {
		return Key{IsCtrl: true, CtrlChar: 'C', Rune: 3}, nil
	}

	// Handle Ctrl+D (EOF)
	if b == 4 {
		return Key{IsCtrl: true, CtrlChar: 'D', Rune: 4}, nil
	}

	// Handle Escape sequences (ANSI)
	if b == 27 {
		// Check if this is a standalone Escape or the start of a sequence
		// Use Buffered() to check if there are bytes available without blocking
		if kr.reader.Buffered() == 0 {
			return Key{IsEscape: true}, nil
		}
		next1, err := kr.reader.ReadByte()
		if err != nil {
			return Key{IsEscape: true}, nil
		}

		// If next byte is '[', it's an ANSI escape sequence
		if next1 == '[' {
			// Read more bytes to determine the sequence
			next2, err := kr.reader.ReadByte()
			if err != nil {
				return Key{IsEscape: true}, nil
			}

			// Check for arrow keys and other special keys
			switch {
			case next2 >= 'A' && next2 <= 'D':
				// Arrow keys
				var dir string
				switch next2 {
				case 'A':
					dir = "up"
				case 'B':
					dir = "down"
				case 'C':
					dir = "right"
				case 'D':
					dir = "left"
				}
				// Read and ignore any trailing bytes
				for {
					_, err := kr.reader.ReadByte()
					if err != nil {
						break
					}
				}
				return Key{IsArrow: true, ArrowDir: dir}, nil

			case next2 == '1':
				// Could be Home (ESC[1~)
				next3, err := kr.reader.ReadByte()
				if err != nil {
					return Key{IsEscape: true}, nil
				}
				if next3 == '~' {
					return Key{IsHome: true}, nil
				}
				return Key{IsEscape: true}, nil

			case next2 == '2':
				// Could be Insert (ESC[2~)
				next3, err := kr.reader.ReadByte()
				if err != nil {
					return Key{IsEscape: true}, nil
				}
				if next3 == '~' {
					return Key{}, nil // Ignore Insert
				}
				return Key{IsEscape: true}, nil

			case next2 == '3':
				// Could be Delete (ESC[3~)
				next3, err := kr.reader.ReadByte()
				if err != nil {
					return Key{IsEscape: true}, nil
				}
				if next3 == '~' {
					return Key{IsDelete: true}, nil
				}
				return Key{IsEscape: true}, nil

			case next2 == '4':
				// Could be End (ESC[4~)
				next3, err := kr.reader.ReadByte()
				if err != nil {
					return Key{IsEscape: true}, nil
				}
				if next3 == '~' {
					return Key{IsEnd: true}, nil
				}
				return Key{IsEscape: true}, nil

			case next2 == '5':
				// Could be Page Up (ESC[5~)
				next3, err := kr.reader.ReadByte()
				if err != nil {
					return Key{IsEscape: true}, nil
				}
				if next3 == '~' {
					return Key{IsPageUp: true}, nil
				}
				return Key{IsEscape: true}, nil

			case next2 == '6':
				// Could be Page Down (ESC[6~)
				next3, err := kr.reader.ReadByte()
				if err != nil {
					return Key{IsEscape: true}, nil
				}
				if next3 == '~' {
					return Key{IsPageDown: true}, nil
				}
				return Key{IsEscape: true}, nil

			case next2 == 'H':
				// Home key
				return Key{IsHome: true}, nil
			case next2 == 'F':
				// End key
				return Key{IsEnd: true}, nil

			default:
				// Unknown escape sequence - could be standalone Escape
				return Key{IsEscape: true}, nil
			}
		}

		// Standalone Escape
		return Key{IsEscape: true}, nil
	}

	// Handle Ctrl+key combinations (Ctrl+A = 1, Ctrl+B = 2, etc.)
	// Note: Line feed (10) and carriage return (13) are excluded as they should be treated as Enter
	if b >= 1 && b <= 26 && b != 10 && b != 13 {
		// Ctrl+A is 1, Ctrl+B is 2, ..., Ctrl+Z is 26
		return Key{IsCtrl: true, CtrlChar: rune('A' + int(b) - 1), Rune: rune(b)}, nil
	}

	// Handle special Ctrl keys (27-31)
	if b >= 27 && b <= 31 {
		var ctrlChar rune
		switch b {
		case 27:
			ctrlChar = '['
		case 28:
			ctrlChar = '\\'
		case 29:
			ctrlChar = ']'
		case 30:
			ctrlChar = '^'
		case 31:
			ctrlChar = '_'
		}
		return Key{IsCtrl: true, CtrlChar: ctrlChar, Rune: rune(b)}, nil
	}

	// Handle special keys
	switch b {
	case 13, 10:
		return Key{IsEnter: true, Rune: '\n'}, nil
	case 8, 127:
		return Key{IsBackspace: true, Rune: 8}, nil
	case 9:
		return Key{IsTab: true, Rune: '\t'}, nil
	case 27:
		// Already handled above, but just in case
		return Key{IsEscape: true}, nil
	}

	// Regular character
	return Key{Rune: rune(b)}, nil
}

// Close cleans up the KeyReader.
func (kr *KeyReader) Close() error {
	kr.mu.Lock()
	defer kr.mu.Unlock()

	if kr.oldState != nil {
		if err := term.Restore(int(os.Stdin.Fd()), kr.oldState); err != nil {
			return err
		}
		kr.oldState = nil
	}
	return nil
}

// EnableRawMode enables raw terminal mode for better key reading.
func (kr *KeyReader) EnableRawMode() error {
	kr.mu.Lock()
	defer kr.mu.Unlock()

	if kr.rawMode {
		return nil
	}

	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	kr.oldState = state
	kr.rawMode = true
	return nil
}

// DisableRawMode disables raw terminal mode.
func (kr *KeyReader) DisableRawMode() error {
	kr.mu.Lock()
	defer kr.mu.Unlock()

	if !kr.rawMode {
		return nil
	}

	if kr.oldState != nil {
		if err := term.Restore(int(os.Stdin.Fd()), kr.oldState); err != nil {
			return err
		}
		kr.oldState = nil
	}
	kr.rawMode = false
	return nil
}

// Colorize wraps text with ANSI color codes if terminal supports it.
func Colorize(text string, colorCode string) string {
	if !SupportsANSI() {
		return text
	}
	return colorCode + text + Reset
}

// HighlightMatch takes a path and match indices and returns the path with
// matching characters highlighted using ANSI color codes.
func HighlightMatch(path string, matchIndices []int) string {
	if !SupportsANSI() || len(matchIndices) == 0 {
		return path
	}

	// Sort indices
	sortedIndices := make([]int, len(matchIndices))
	copy(sortedIndices, matchIndices)
	for i := 0; i < len(sortedIndices)-1; i++ {
		for j := i + 1; j < len(sortedIndices); j++ {
			if sortedIndices[i] > sortedIndices[j] {
				sortedIndices[i], sortedIndices[j] = sortedIndices[j], sortedIndices[i]
			}
		}
	}

	var result strings.Builder
	prevIdx := 0

	for _, idx := range sortedIndices {
		// Add the non-matching part
		if idx > prevIdx {
			result.WriteString(path[prevIdx:idx])
		}
		// Add the matching character with highlighting
		if idx < len(path) {
			result.WriteString(Colorize(string(path[idx]), ColorBrightGreen))
			prevIdx = idx + 1
		}
	}

	// Add the remaining part
	if prevIdx < len(path) {
		result.WriteString(path[prevIdx:])
	}

	return result.String()
}

// GetMatchHighlight returns a function that can be used to highlight matches.
func GetMatchHighlight(matchIndices []int) func(string) string {
	return func(s string) string {
		return HighlightMatch(s, matchIndices)
	}
}

// RepeatString repeats a string n times.
func RepeatString(s string, n int) string {
	var result strings.Builder
	for i := 0; i < n; i++ {
		result.WriteString(s)
	}
	return result.String()
}

// PadRight pads a string to the specified width with spaces on the right.
func PadRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + RepeatString(" ", width-len(s))
}

// PadLeft pads a string to the specified width with spaces on the left.
func PadLeft(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return RepeatString(" ", width-len(s)) + s
}

// Truncate truncates a string to the specified length, adding ellipsis if truncated.
func Truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	if length <= 0 {
		return ""
	}
	if length == 1 {
		return "."
	}
	if length == 2 {
		return ".."
	}
	if length == 3 {
		return "..."
	}
	return s[:length-3] + "..."
}


