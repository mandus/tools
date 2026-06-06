---
name: pass-tui-update
version: 1.0.0
date: 2026-06-06
authors: [mandu]
status: draft
---

# Specification: Pass Program TUI Update with Bubble Tea

## Overview

This specification outlines the update to the `pass` password manager to replace the current buggy custom terminal handling with a proper TUI implementation using the [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) framework.

### Problem Statement

The current implementation has several critical issues:

1. **Buggy key handling**: The custom terminal handling in `pkg/terminal/terminal.go` has bugs with escape sequence detection that cause the program to freeze
2. **Inconsistent cross-platform behavior**: Keys are not caught correctly across different terminals and operating systems
3. **Poor user experience**: The fuzzy search UI lacks proper help information and has inconsistent behavior
4. **Maintenance burden**: Custom terminal handling code is complex and error-prone

### Solution

Replace the custom terminal handling with Bubble Tea, which provides:
- Cross-platform keyboard handling that works consistently on Windows, Linux, and macOS
- Built-in support for ANSI escape codes and terminal capabilities
- Well-tested, maintained library with a large ecosystem
- Better architecture for TUI applications

## Goals

### Primary Goals
- [ ] Fix all keyboard input bugs (arrow keys, Home, End, Delete, Page Up/Down)
- [ ] Ensure consistent behavior across Windows, Linux, and macOS
- [ ] Use full window size for listing as many completions as possible
- [ ] Display basic help information in the TUI
- [ ] Maintain all existing functionality

### Secondary Goals (Future)
- [ ] Add configurable number of lines to use in terminal (not implemented in this spec)
- [ ] Add more sophisticated styling and theming
- [ ] Add mouse support

## Scope

### In Scope
- Replace fuzzy search implementation with Bubble Tea
- Update `pass` (default command with fuzzy search)
- Update `pass -c` (copy with fuzzy search)
- Update `pass rm` (remove with fuzzy search)
- Update any other commands that use fuzzy search
- Maintain backward compatibility with existing CLI usage

### Out of Scope
- Replace non-interactive commands (e.g., `pass show <path>`, `pass ls`)
- Add new features beyond fixing the TUI
- Configurable terminal line count (noted for future implementation)

## Background

### Current Architecture

The current implementation has:
- `pkg/terminal/terminal.go`: Custom terminal handling with bugs
- `pkg/terminal/key_reader.go`: Custom key reader
- `cmd/fuzzy.go`: Fuzzy search UI logic
- `pkg/fuzzy/fuzzy.go`: Fuzzy matching algorithm (this can be reused)

The main issue is in `pkg/terminal/terminal.go` where the `ReadKey()` function incorrectly uses `UnreadByte()` which breaks escape sequence detection.

### Proposed Architecture

```
pass/
├── cmd/
│   ├── root.go          # CLI entry points (unchanged)
│   ├── show.go          # Non-interactive show (unchanged)
│   ├── ls.go            # Non-interactive ls (unchanged)
│   ├── find.go          # Non-interactive find (unchanged)
│   ├── insert.go        # Non-interactive insert (unchanged)
│   ├── rm.go            # Updated to use new TUI
│   └── tui/             # NEW: TUI components
│       ├── fuzzy.go     # Fuzzy search TUI using Bubble Tea
│       ├── models.go    # Bubble Tea models
│       └── styles.go    # Styling for TUI
├── pkg/
│   ├── fuzzy/           # Reuse existing fuzzy matching
│   ├── gpg/             # Reuse existing GPG handling
│   ├── filesystem/      # Reuse existing filesystem code
│   └── config/          # Reuse existing config
└── go.mod               # Add bubbletea dependencies
```

## Requirements

### Functional Requirements

#### FR-001: Fuzzy Search TUI
The TUI MUST provide fuzzy search functionality with the following characteristics:
- Display all passwords in the password store
- Filter passwords as user types
- Allow navigation using arrow keys (up/down)
- Allow selection using Enter key
- Exit using Esc, Ctrl+C, Ctrl+D, or Ctrl+Q
- Use full terminal width and height
- Display as many completions as possible in the available space

#### FR-002: Mode-Specific Behavior
The TUI MUST support different modes:
- **Show Mode**: Display selected password to stdout
- **Clip Mode**: Copy selected password to clipboard
- **Remove Mode**: Delete selected password (with confirmation)

#### FR-003: Keyboard Support
The TUI MUST support the following keys:
- Arrow keys: Navigate up/down through results
- Enter: Select highlighted item
- Esc: Exit/cancel
- Ctrl+C: Exit/cancel
- Ctrl+D: Exit/cancel
- Ctrl+Q: Exit/cancel
- Home: Jump to first item
- End: Jump to last item
- Page Up: Scroll up by page
- Page Down: Scroll down by page
- Tab: Cycle through results (optional)
- Backspace: Delete character from query
- Delete: Delete character after cursor
- Ctrl+A: Move cursor to start of query
- Ctrl+E: Move cursor to end of query
- Ctrl+K: Clear from cursor to end of query
- Ctrl+L: Clear entire query
- Ctrl+W: Delete word before cursor

#### FR-004: Display Requirements
The TUI MUST display:
- A header indicating the current mode (Search, Copy, Remove)
- A search prompt with the current query
- A list of matching passwords with the best matches first
- Match highlighting (characters that match the query)
- Current selection indicator (e.g., `>` prefix)
- Help text showing available keys
- Scroll position indicator if list is longer than screen

#### FR-005: Help Information
The TUI MUST display basic help information:
- At minimum: available navigation keys and action keys
- Should be visible without obscuring the main content
- Should be clear and concise

### Non-Functional Requirements

#### NFR-001: Cross-Platform Compatibility
The TUI MUST work correctly on:
- Windows (Windows Terminal, PowerShell, CMD)
- Linux (GNOME Terminal, Konsole, xterm, etc.)
- macOS (Terminal.app, iTerm2)

#### NFR-002: Performance
- Search and filtering MUST be responsive (< 100ms for 1000 passwords)
- Rendering MUST be smooth without flickering
- Must handle password stores with 1000+ entries efficiently

#### NFR-003: Accessibility
- Must work in terminals without ANSI color support (fallback to plain text)
- Must work in terminals with limited width (< 40 columns)
- Must work in terminals with limited height (< 10 rows)

#### NFR-004: Backward Compatibility
- Non-interactive commands MUST continue to work exactly as before
- Existing scripts and workflows MUST not be broken
- CLI flags and arguments MUST remain the same

### Technical Requirements

#### TR-001: Dependencies
Add the following dependencies to `go.mod`:
```
github.com/charmbracelet/bubbletea v0.25.0
github.com/charmbracelet/bubbles v0.16.1
github.com/charmbracelet/lipgloss v0.11.0
```

#### TR-002: Code Organization
- TUI-related code MUST be in a separate package (`cmd/tui/`)
- Existing fuzzy matching logic MUST be reused
- Terminal handling code in `pkg/terminal/` MAY be deprecated but not removed (for backward compatibility)

#### TR-003: Testing
- Unit tests for TUI models
- Integration tests for key handling
- Manual testing on all target platforms

## Design

### Component Design

#### TUI Model Structure

```go
type FuzzyModel struct {
    // Bubble Tea list component
    list list.Model
    
    // Search input
    input textinput.Model
    
    // Current mode
    mode FuzzySearchMode
    
    // All passwords (for filtering)
    allPasswords []string
    
    // Filtered passwords
    filteredPasswords []fuzzy.MatchResult
    
    // State
    loading   bool
    error     error
    quitting  bool
    selected  string
}
```

#### Mode Enum

```go
type FuzzySearchMode int

const (
    FuzzyModeShow FuzzySearchMode = iota
    FuzzyModeClip
    FuzzyModeRm
)
```

### Screen Layout

```
┌─────────────────────────────────────────────────────────────┐
│  Select password (Enter to show, Esc to cancel)                │
│                                                             │
│  Search: mypass                                              │
│                                                             │
│  > email/gmail.com/user                                     │
│    email/outlook.com/work                                   │
│    social/github.com                                       │
│    social/twitter.com                                       │
│    work/vpn/corporate                                       │
│                                                             │
│                                                             │
│  ↑/↓: Navigate  Enter: Select  Esc: Cancel  Ctrl+C: Quit      │
└─────────────────────────────────────────────────────────────┘
```

### Color Scheme

- Header: Bold white on dark background
- Search prompt: Cyan
- Selected item: Bright green with `>` prefix
- Match highlighting: Bright green
- Help text: Dimmed

### Keyboard Mapping

| Action | Keys |
|--------|------|
| Navigate up | ↑, k |
| Navigate down | ↓, j |
| Select | Enter |
| Exit | Esc, Ctrl+C, Ctrl+D, Ctrl+Q |
| Page up | PageUp, Ctrl+U |
| Page down | PageDown, Ctrl+D |
| Jump to first | Home, Ctrl+A |
| Jump to last | End, Ctrl+E |
| Clear query | Ctrl+L |
| Delete char before | Backspace |
| Delete char after | Delete |
| Delete word | Ctrl+W |
| Clear to end | Ctrl+K |

## Implementation Plan

### Phase 1: Setup (1 session)
1. Create `cmd/tui/` directory
2. Add Bubble Tea dependencies to `go.mod`
3. Create basic TUI skeleton
4. Ensure dependencies compile correctly

### Phase 2: Core TUI Implementation (2-3 sessions)
1. Implement `FuzzyModel` with list and input components
2. Implement key handling
3. Implement rendering with proper layout
4. Add match highlighting
5. Add help text display

### Phase 3: Mode Support (1-2 sessions)
1. Implement Show mode
2. Implement Clip mode
3. Implement Remove mode with confirmation
4. Test all modes

### Phase 4: Integration (1-2 sessions)
1. Update `cmd/root.go` to use new TUI for fuzzy search
2. Update `cmd/rm.go` to use new TUI
3. Ensure non-interactive commands still work
4. Comprehensive testing

### Phase 5: Polish and Testing (1-2 sessions)
1. Add styling and colors
2. Handle edge cases (empty store, no matches, etc.)
3. Performance optimization
4. Cross-platform testing

## File Changes

### New Files
- `cmd/tui/fuzzy.go` - Main TUI implementation
- `cmd/tui/models.go` - Bubble Tea models
- `cmd/tui/styles.go` - Styling definitions
- `cmd/tui/keys.go` - Key binding definitions

### Modified Files
- `go.mod` - Add Bubble Tea dependencies
- `cmd/root.go` - Update to use new TUI for fuzzy search
- `cmd/rm.go` - Update to use new TUI
- `cmd/fuzzy.go` - Deprecate or remove (after new TUI is working)

### Deprecated Files (Optional Removal)
- `pkg/terminal/terminal.go` - May be deprecated
- `pkg/terminal/key_reader.go` - May be deprecated

## Testing Strategy

### Unit Tests
- Test key handling in TUI
- Test filtering and sorting
- Test mode transitions
- Test error handling

### Integration Tests
- Test full fuzzy search workflow
- Test each mode (show, clip, rm)
- Test with various terminal sizes

### Manual Tests
- Test on Windows Terminal
- Test on Linux terminal
- Test on macOS Terminal
- Test with different terminal sizes
- Test with empty password store
- Test with large password store (1000+ entries)

### Test Cases

#### TC-001: Basic Navigation
1. Run `pass` with no arguments
2. Verify TUI starts
3. Press ↓ key
4. Verify selection moves down
5. Press ↑ key
6. Verify selection moves up

#### TC-002: Search Filtering
1. Run `pass` with no arguments
2. Type "gmail"
3. Verify list filters to show only matching passwords
4. Verify matches are highlighted

#### TC-003: Selection
1. Run `pass` with no arguments
2. Navigate to a password
3. Press Enter
4. Verify password is displayed

#### TC-004: Exit
1. Run `pass` with no arguments
2. Press Esc
3. Verify TUI exits
4. Repeat with Ctrl+C, Ctrl+D, Ctrl+Q

#### TC-005: Clip Mode
1. Run `pass -c` with no arguments
2. Navigate to a password
3. Press Enter
4. Verify password is copied to clipboard

#### TC-006: Remove Mode
1. Run `pass rm` with no arguments
2. Navigate to a password
3. Press Enter
4. Verify confirmation dialog appears
5. Confirm removal
6. Verify password is removed

#### TC-007: Edge Cases
1. Run `pass` with empty password store
2. Verify appropriate error message
3. Run `pass` with no matches for query
4. Verify "No matches" message

## Migration Strategy

### Step 1: Parallel Implementation
- Keep existing fuzzy search implementation
- Implement new TUI alongside
- Use feature flag to switch between implementations

### Step 2: Testing
- Test new TUI thoroughly
- Compare behavior with old implementation
- Fix any discrepancies

### Step 3: Cutover
- Replace old implementation with new TUI
- Monitor for issues
- Roll back if necessary

### Step 4: Cleanup
- Remove old implementation (optional)
- Remove feature flag
- Update documentation

## Rollback Plan

If the new TUI has issues:
1. Revert changes to `cmd/root.go` and `cmd/rm.go`
2. Keep new TUI code but don't use it
3. Fix issues in new TUI
4. Retry cutover

## Open Questions

1. **Should we keep the old implementation as a fallback?**
   - Pro: Safety net
   - Con: Maintenance burden
   - Recommendation: Remove after thorough testing

2. **Should we add a CLI flag to use the old implementation?**
   - Pro: Allows users to opt-out if they have issues
   - Con: Complexity
   - Recommendation: No, but document how to downgrade

3. **Should we add animation or other visual feedback?**
   - Pro: Better UX
   - Con: Complexity, potential performance impact
   - Recommendation: Not in initial implementation

4. **How should we handle very large password stores?**
   - Proposal: Implement pagination or lazy loading
   - Recommendation: Implement basic pagination first

## Future Enhancements

### Version 2.0
- Configurable number of lines to use in terminal
- Custom color schemes
- Mouse support
- Multi-select for batch operations
- Tree view for hierarchical password stores

### Version 3.0
- Password generator integration
- Password strength analysis
- Export/import functionality
- Sync status indicators

## Appendix

### References
- [Bubble Tea GitHub](https://github.com/charmbracelet/bubbletea)
- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/examples/tutorials)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Lip Gloss Styling](https://github.com/charmbracelet/lipgloss)
- [Current Pass Implementation](../cmd/fuzzy.go)
- [Current Terminal Handling](../pkg/terminal/terminal.go)

### Glossary
- **TUI**: Text-based User Interface
- **Bubble Tea**: A Go framework for building TUIs
- **ANSI**: American National Standards Institute escape codes for terminal control
- **Fuzzy Matching**: String matching that allows for typos and partial matches

### Revision History
- 2026-06-06: Initial draft created
