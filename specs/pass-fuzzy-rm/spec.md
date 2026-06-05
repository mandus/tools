# Pass Fuzzy Search & Remove Feature Specification

## Overview

This document specifies two new features for the pass replacement tool:
1. **Interactive Fuzzy Search**: When invoked without arguments, `pass` enters an interactive fuzzy search mode that allows users to quickly find and select secrets by typing partial matches.
2. **Remove Command**: `pass rm <path>` removes a secret file and commits the removal to git.

Both features are designed to match the behavior of Unix fzf (fuzzy finder) and the standard `pass rm` command respectively.

---

## User Requirements

### Must Have
- [ ] Invoke `pass` without arguments to enter fuzzy search mode
- [ ] Real-time filtering as user types
- [ ] Fuzzy matching: characters must appear in order but not consecutively
- [ ] Best match auto-selected
- [ ] Matching characters visually highlighted
- [ ] Navigation: arrow keys, Enter, Esc, Ctrl+C
- [ ] Additional keybindings: Ctrl+A (move to start), Ctrl+K (clear from cursor to end), Ctrl+E (move to end)
- [ ] Press Enter on selection to show the password (or delete if invoked via `pass rm`)
- [ ] `pass rm <path>` removes the secret file
- [ ] `pass rm` without arguments enters fuzzy search mode, then deletes on Enter
- [ ] `pass rm` auto-commits removal to git with message "Remove <path>"
- [ ] No confirmation prompts for `rm` (git history is sufficient)

### Should Have
- [ ] Respect `.gitignore` - only show tracked/valid password files
- [ ] Case-insensitive matching by default
- [ ] Show match quality/ranking indicators
- [ ] Support `--no-commit` flag for `rm` to skip git
- [ ] Handle errors gracefully with clear messages

### Nice to Have (Future)
- [ ] Preview pane showing first line of matched secret
- [ ] Multiple selection for batch operations
- [ ] History of previous fuzzy searches
- [ ] Custom color schemes
- [ ] Mouse support in terminal

---

## Detailed Specifications

---

## 1. Fuzzy Search Feature

### 1.1 Invocation

**Command:**
```
pass
```

When `pass` is invoked without any arguments and without an explicit command (insert, ls, find, show, rm), it enters interactive fuzzy search mode.

**Alternative invocation:**
```
pass --interactive
pass -i
```

### 1.2 Fuzzy Matching Algorithm

**Subsequence Matching:**
- The query string must match as a subsequence of the target path
- Characters in the query must appear in the same order in the target
- Characters do NOT need to be consecutive in the target

**Example matches:**
- Query `twt` matches `social/twitter.com/admin` (t-w-t in order)
- Query `gm` matches `email/gmail.com/user` (g-m in order)
- Query `chaseb` matches `banking/chase.com/account` (c-h-a-s-e-b in order)

**Non-matches:**
- Query `twtt` does NOT match `twitter` (two t's required, only one available)
- Query `mtw` does NOT match `twitter` (m comes after t in query, but before t in target)

**Scoring Algorithm:**
```
score = 0
query_idx = 0

for each character in target:
    if query_idx < len(query) and character == query[query_idx]:
        # Bonus for consecutive matches
        if prev_char_was_match:
            score += 10  # Consecutive match bonus
        # Bonus for match at start of path component
        if character is at start of path component:
            score += 5
        score += 100 - (position_in_target * 2)  # Earlier matches better
        query_idx++
    else:
        score += 1  # Penalty for non-matching chars
        
# Penalty for long paths
score -= len(target) * 2

# Bonus for short paths
score += (max_path_length - len(target)) * 3

# Lower score is better
```

### 1.3 Display

**Layout:**
```
Passwords:
  email/gmail.com/user
> social/twitter.com/admin
  bank/chase.com/account
  
Search: tw
```

**Components:**
- **Header**: "Passwords:" label
- **List**: Scrollable list of matching passwords
- **Cursor**: `>` prefix indicates currently selected item
- **Prompt**: "Search: " at bottom with user's query
- **Matching characters**: Highlighted (using terminal colors if available, otherwise just display normally)

**Terminal Control:**
- Use ANSI escape codes for cursor movement and screen updates
- Clear and redraw entire display on each keystroke
- Support terminals without ANSI: fall back to simple redraw

### 1.4 Keybindings

| Key | Action |
|-----|--------|
| Any printable character | Add to query, re-filter list |
| Backspace | Remove last character from query |
| Delete | Remove character under cursor from query |
| Ctrl+A | Move cursor to start of query |
| Ctrl+E | Move cursor to end of query |
| Ctrl+K | Delete from cursor to end of query |
| Ctrl+L | Clear entire query |
| ↑ (Up Arrow) | Move selection up by 1 |
| ↓ (Down Arrow) | Move selection down by 1 |
| Page Up | Move selection up by page height |
| Page Down | Move selection down by page height |
| Home | Move to first item |
| End | Move to last item |
| Enter | Select current item (show or delete based on invocation) |
| Esc | Exit fuzzy search mode, return to shell |
| Ctrl+C | Exit fuzzy search mode, return to shell |
| Ctrl+D | Exit fuzzy search mode (EOF), return to shell |
| Tab | Toggle between search input and list navigation |

### 1.5 Selection Action

The action taken when Enter is pressed depends on how fuzzy search was invoked:

| Invocation | Action |
|------------|--------|
| `pass` (no args) | Show the selected password (to stdout) |
| `pass -c` (with clip flag) | Copy selected password to clipboard |
| `pass rm` | Delete the selected password file |
| `pass rm -c` | Delete and copy password to clipboard before deleting |

### 1.6 Filtering Behavior

- **Empty query**: Show all passwords
- **No matches**: Show "No matches found" message
- **Case sensitivity**: Case-insensitive by default
- **Git ignore**: Only show files that would be tracked by git (exclude `.git/` directory)
- **File types**: Only show `.gpg` files
- **Sorting**: Best matches first, then alphabetically

---

## 2. Remove Command

### 2.1 Synopsis

**Remove a single password:**
```
pass rm <path>
```

**Remove with fuzzy search:**
```
pass rm
```

**Remove without git commit:**
```
pass rm --no-commit <path>
pass rm -n <path>
```

### 2.2 Description

Removes a password file from the store and optionally commits the removal to git.

### 2.3 Arguments

- `<path>`: The path of the password to remove
  - Type: String, optional (if omitted, enters fuzzy search mode)
  - Can omit `.gpg` extension (automatically appended if not present)
  - Must be a valid path to an existing `.gpg` file

### 2.4 Options

- `-n, --no-commit`: Skip git commit after removal
- `-f, --force`: Skip git commit (alias for --no-commit)
- `-c, --clip`: Copy password to clipboard before deleting

### 2.5 Behavior

**With explicit path:**
1. Validate path is not empty
2. Construct file path (add `.gpg` if not present)
3. Check if file exists
   - If not, display error `pass: <path>: No such file or directory` and exit with code 1
4. If `-c` flag is set:
   - Decrypt and copy password to clipboard
   - Start clipboard clear timer
5. Remove the file
6. If git repo exists and `--no-commit` not specified:
   - Run `git rm <path>.gpg`
   - Run `git commit -m "Remove <path>"`
   - If git commands fail, display warning but continue
7. Display: "Password removed successfully."
8. Exit with code 0

**Without explicit path (fuzzy search mode):**
1. Enter fuzzy search mode
2. As user types, filter list of passwords
3. When user presses Enter on a selection:
   - If `-c` flag: copy to clipboard first
   - Remove the selected password file
   - Commit to git (unless --no-commit)
   - Display success message
   - Exit fuzzy search mode
4. If user presses Esc/Ctrl+C: exit without action

### 2.6 Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | File not found |
| 2 | Permission denied |
| 3 | Git operation failed (non-fatal, still removed) |

### 2.7 Examples

```bash
# Remove specific password
pass rm email/gmail.com/oldaccount

# Remove with fuzzy search
pass rm
# User types: gmail
# Selects: email/gmail.com/oldaccount
# Press Enter to delete

# Remove without git commit
pass rm --no-commit social/twitter.com/old

# Remove and copy to clipboard first
pass rm -c banking/chase.com/oldcard
```

---

## 3. Git Integration for Remove

### 3.1 Remove and Commit

```go
func RemoveAndCommit(filePath, message string) error {
    // Remove file
    if err := os.Remove(filePath); err != nil {
        return fmt.Errorf("pass: failed to remove file: %v", err)
    }
    
    // Git rm
    dir := filepath.Dir(filePath)
    base := filepath.Base(filePath)
    
    cmd := exec.Command("git", "rm", base)
    cmd.Dir = dir
    if err := cmd.Run(); err != nil {
        // Non-fatal: file is removed, just warn about git
        fmt.Fprintf(os.Stderr, "pass: warning: git rm failed: %v\n", err)
    }
    
    // Commit
    cmd = exec.Command("git", "commit", "-m", message)
    cmd.Dir = dir
    if err := cmd.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "pass: warning: git commit failed: %v\n", err)
    }
    
    return nil
}
```

### 3.2 No Confirmation

Per user requirement, no confirmation is needed for remove operations because:
- The password is still in git history
- Users can recover via `git checkout` if needed
- This matches the Unix pass behavior with `--force` flag

---

## 4. Implementation Details

### 4.1 File Structure

```
pass/
├── cmd/
│   ├── fuzzy.go       # Fuzzy search command and UI
│   ├── fuzzy_test.go  # Tests for fuzzy search
│   ├── rm.go          # Remove command
│   ├── rm_test.go     # Tests for remove command
│   └── root.go        # Modified to handle fuzzy mode
├── pkg/
│   ├── fuzzy/
│   │   ├── fuzzy.go   # Fuzzy matching algorithm
│   │   └── fuzzy_test.go
│   └── terminal/
│       ├── terminal.go # Terminal UI utilities
│       └── terminal_test.go
└── go.mod
```

### 4.2 Dependencies

**New dependencies:**
- `golang.org/x/term` - Already in use for password input
- No additional external dependencies needed

**Terminal handling:**
- Use `golang.org/x/term` for terminal size detection
- Use ANSI escape codes for cursor control
- Fall back to simple mode if terminal doesn't support ANSI

### 4.3 Fuzzy Matching Package

```go
// pkg/fuzzy/fuzzy.go

package fuzzy

// Match checks if query is a subsequence of target (case-insensitive)
func Match(query, target string) bool

// Score returns the match quality (lower is better)
func Score(query, target string) int

// Filter returns all items matching query, sorted by score
func Filter(query string, items []string) []MatchResult

type MatchResult struct {
    Path      string
    Score     int
    Matches   []int  // Indices of matching characters
}
```

### 4.4 Terminal UI Package

```go
// pkg/terminal/terminal.go

package terminal

// ClearScreen clears the terminal screen
func ClearScreen()

// MoveCursor moves cursor to (row, col)
func MoveCursor(row, col int)

// HideCursor hides the terminal cursor
func HideCursor()

// ShowCursor shows the terminal cursor
func ShowCursor()

// GetSize returns terminal width and height
func GetSize() (int, int, error)

// PrintAt prints text at specific position
func PrintAt(text string, row, col int)

// ReadKey reads a single key press (non-blocking if possible)
func ReadKey() (Key, error)

type Key struct {
    Rune rune
    IsArrow bool
    ArrowDir string // up, down, left, right
    IsCtrl bool
    CtrlChar rune
}
```

---

## 5. Error Handling

### 5.1 Fuzzy Search Errors

| Scenario | Behavior |
|----------|----------|
| Terminal too small | Show warning, proceed with minimal display |
| No terminal (piped) | Fall back to non-interactive mode, show error |
| No matches | Show "No matches found" |
| File read error | Show error, continue with available files |

### 5.2 Remove Errors

| Scenario | Behavior |
|----------|----------|
| File not found | Error: `pass: <path>: No such file or directory`, exit 1 |
| Permission denied | Error: `pass: failed to remove <path>: permission denied`, exit 2 |
| Not a .gpg file | Error: `pass: <path>: Not a password file`, exit 1 |
| Git rm fails | Warning to stderr, continue |
| Git commit fails | Warning to stderr, continue |

---

## 6. Testing Strategy

### 6.1 Unit Tests

**Fuzzy matching:**
- Test `Match()` with various query/target combinations
- Test `Score()` returns correct ordering
- Test `Filter()` returns properly sorted results
- Test case sensitivity
- Test empty query

**Terminal UI:**
- Test ANSI escape code generation
- Test key parsing
- Test cursor position calculations
- Test terminal size detection

**Remove command:**
- Test file path construction
- Test file existence check
- Test git integration
- Test error handling

### 6.2 Integration Tests

**Fuzzy search:**
- Test full interactive session (mock terminal input)
- Test navigation with arrow keys
- Test Enter to select
- Test Esc to exit
- Test with various terminal sizes

**Remove:**
- Test remove with explicit path
- Test remove with fuzzy search
- Test remove with --no-commit
- Test remove with -c flag
- Test error cases

### 6.3 End-to-End Tests

- Full workflow: insert → fuzzy search → show
- Full workflow: insert → fuzzy search → rm
- Full workflow: insert → rm <path>
- Git history verification after rm

---

## 7. Compatibility

### 7.1 Unix pass Compatibility

| Feature | Unix pass | This implementation |
|---------|-----------|-------------------|
| `pass rm <path>` | ✓ (with -f) | ✓ |
| No confirmation | ✓ (with -f) | ✓ (always) |
| Git commit on rm | ✓ (with -f) | ✓ (default) |
| Fuzzy search | ✗ (requires fzf) | ✓ (built-in) |

### 7.2 fzf Compatibility

| Feature | fzf | This implementation |
|---------|-----|-------------------|
| Subsequence matching | ✓ | ✓ |
| Real-time filtering | ✓ | ✓ |
| Ctrl+K to clear line | ✓ | ✓ |
| Ctrl+A/Ctrl+E navigation | ✓ | ✓ |
| ANSI color support | ✓ | ✓ (optional) |
| Mouse support | ✓ | ✗ (future) |
| Preview window | ✓ | ✗ (future) |

---

## 8. Open Questions (RESOLVED)

### OQ-001: Respect .gitignore?
**Status**: RESOLVED  
**Decision**: Yes - only show files that would be tracked by git. This means excluding the `.git/` directory and any files that match patterns in `.gitignore`.

### OQ-002: Confirmation before delete?
**Status**: RESOLVED  
**Decision**: No confirmation needed. The user explicitly stated that git history is sufficient, so we skip confirmation prompts.

### OQ-003: Keybindings?
**Status**: RESOLVED  
**Decision**: Use standard bindings (↑/↓ for navigation, Enter for select, Esc/Ctrl+C for exit) plus Ctrl+A (move to first char), Ctrl+K (remove from cursor position), Ctrl+E (move to end), and arrow left/right to move in the search input.

### OQ-004: Persist last query?
**Status**: RESOLVED  
**Decision**: No - fresh each invocation, as proposed.

---

## 9. Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PASS_FUZZY_TIMEOUT` | Timeout for fuzzy search selection (seconds) | 0 (no timeout) |
| `PASS_NO_COLOR` | Disable ANSI color codes | false |

---

## 10. References

- [fzf - Fuzzy Finder](https://github.com/junegunn/fzf)
- [Unix pass rm command](https://git.zx2c4.com/password-store/tree/man/pass.1.md)
- [ANSI escape codes](https://en.wikipedia.org/wiki/ANSI_escape_code)
- [golang.org/x/term](https://pkg.go.dev/golang.org/x/term)

---

*Document Version: 1.0*  
*Last Updated: 2026-06-05*  
*Author: Mandu*  
*Status: Approved for Implementation*
