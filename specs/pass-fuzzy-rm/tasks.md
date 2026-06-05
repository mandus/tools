# Pass Fuzzy Search & Remove - Implementation Tasks

## Overview

This document breaks down the implementation of fuzzy search and rm features into actionable tasks.

---

## Task List

### Phase 1: Specification & Documentation
- [x] Create spec document (COMPLETED)
- [x] Create tasks document (COMPLETED)
- [x] Update main pass-replacement-spec.md
- [x] Update implementation checklist
- [x] Update decision log

### Phase 2: Core Packages

#### Fuzzy Matching Package (pkg/fuzzy/)
- [x] Create pkg/fuzzy/fuzzy.go
  - [x] Implement `Match(query, target string) bool` - subsequence check
  - [x] Implement `Score(query, target string) int` - ranking algorithm
  - [x] Implement `Filter(query string, items []string) []MatchResult` - filter and sort
  - [x] Implement `MatchResult` struct with Path, Score, MatchIndices
  - [x] Handle case-insensitive matching
  - [x] Handle empty query (return all)
- [x] Create pkg/fuzzy/fuzzy_test.go
  - [x] Test Match() with various inputs
  - [x] Test Score() returns correct ordering
  - [x] Test Filter() returns sorted results
  - [x] Test edge cases (empty query, empty items, no matches)

#### Terminal Package (pkg/terminal/)
- [x] Create pkg/terminal/terminal.go
  - [x] Implement ANSI escape code constants
  - [x] Implement `ClearScreen()`
  - [x] Implement `MoveCursor(row, col int)`
  - [x] Implement `HideCursor()` and `ShowCursor()`
  - [x] Implement `GetSize() (width, height, error)`
  - [x] Implement `PrintAt(text string, row, col int)`
  - [x] Implement `ReadKey() (Key, error)`
  - [x] Implement `Key` struct with Rune, IsArrow, ArrowDir, IsCtrl, CtrlChar
  - [x] Implement `SupportsANSI() bool` - detect terminal capability
  - [x] Implement `HighlightMatch()` for visual highlighting
  - [x] Implement utility functions (PadRight, PadLeft, Truncate, etc.)
- [x] Create pkg/terminal/terminal_test.go
  - [x] Test ANSI code generation
  - [x] Test cursor movement calculations
  - [x] Test key parsing
  - [x] Test utility functions

### Phase 3: Fuzzy Search Command

#### Fuzzy Search UI (cmd/fuzzy.go)
- [x] Implement fuzzy search mode entry point
- [x] Implement main loop for fuzzy search
- [x] Implement display rendering
  - [x] Header with mode-specific messages
  - [x] List of matching passwords with cursor (>) indicator
  - [x] Prompt "Search: " with query
  - [x] Handle terminal resizing
- [x] Implement query state management
  - [x] String buffer for query
  - [x] Cursor position in query
  - [x] Handle all keybindings (Ctrl+A, Ctrl+E, Ctrl+K, Ctrl+L, Ctrl+W, arrows, etc.)
- [x] Implement list navigation
  - [x] Selected index
  - [x] Scroll position
  - [x] Page up/down
- [x] Implement match highlighting
  - [x] Identify matching character positions
  - [x] Apply terminal color codes for highlighting
- [x] Implement fuzzy search with different invocations
  - [x] Default: show password on Enter
  - [x] With -c flag: copy to clipboard on Enter
  - [x] With rm command: delete on Enter

#### Fuzzy Search Tests (cmd/fuzzy_test.go)
- [x] Test fuzzy search initialization
- [x] Test query input handling
- [x] Test list filtering and sorting
- [x] Test navigation key handling
- [x] Test Enter key handling (different modes)
- [x] Test exit key handling (Esc, Ctrl+C, Ctrl+D)

### Phase 4: Remove Command

#### Remove Command (cmd/rm.go)
- [x] Implement rm command structure
  - [x] Command definition with cobra
  - [x] Flags: --no-commit/-n, --force/-f, --clip/-c
- [x] Implement `removePassword(path string) error`
  - [x] Validate path
  - [x] Construct full file path
  - [x] Check file exists
  - [x] If -c flag: decrypt and copy to clipboard
  - [x] Remove file with os.Remove()
  - [x] If git repo exists and not --no-commit:
    - [x] Run git rm
    - [x] Run git commit
- [x] Implement fuzzy search mode for rm without path
  - [x] Same as default fuzzy search but action is delete
  - [x] After selection: delete the file
  - [x] Handle -c flag with fuzzy search (copy then delete)

#### Remove Tests (cmd/rm_test.go)
- [x] Test remove with explicit path
- [x] Test remove with fuzzy search
- [x] Test remove with --no-commit
- [x] Test remove with -c flag
- [x] Test error: file not found
- [x] Test error: permission denied
- [x] Test: git integration

### Phase 5: Integration with Root Command

#### Update root.go
- [x] Modify default behavior when no args
  - [x] If no args and no command: enter fuzzy search mode (show)
  - [x] If no args with -c flag: enter fuzzy search mode (clip)
  - [x] Pass fuzzy search mode flag (normal, clip, rm)
- [x] Add rm command registration
- [x] Ensure fuzzy search respects global flags (-c)
- [x] Update Long description to mention fuzzy search

### Phase 6: Git Integration Enhancements

#### Update pkg/git/git.go
- [x] Add `RemoveAndCommit(filePath, message string) error`
  - [x] Remove file
  - [x] Git rm
  - [x] Git commit
  - [x] Handle errors gracefully (non-fatal)
- [x] Update git package with existing tests

### Phase 7: Testing & Validation

- [ ] Run all existing tests - ensure no regressions
- [ ] Run new tests - ensure features work
- [ ] Manual testing
  - [ ] Test fuzzy search with various inputs
  - [ ] Test rm with explicit path
  - [ ] Test rm with fuzzy search
  - [ ] Test all keybindings
  - [ ] Test edge cases
- [ ] Performance testing
  - [ ] Large password stores (1000+ entries)
  - [ ] Long paths

### Phase 8: Documentation Updates

- [x] Update pass-replacement-spec.md with fuzzy search and rm details
- [x] Update pass-decision-log.md with new decisions
- [x] Update pass-implementation-checklist.md with new tasks
- [x] Update specs/pass-fuzzy-rm/spec.md with full specification
- [ ] Verify all docs are consistent

---

## Priority Order

1. **P0 - Critical** (Blockers)
   - Fuzzy matching algorithm
   - Terminal input/output
   - Basic fuzzy search loop
   - rm command with explicit path

2. **P1 - High** (Core functionality)
   - Fuzzy search navigation
   - rm with fuzzy search
   - Git integration for rm
   - All keybindings

3. **P2 - Medium** (Polish)
   - Match highlighting
   - Error handling
   - Tests

4. **P3 - Low** (Nice to have)
   - Documentation updates
   - Performance optimizations

---

## Dependencies Between Tasks

```
Fuzzy Package <-- Terminal Package
    |                |
    v                v
Fuzzy Command    Remove Command
    |                |
    +----------------+
           |
           v
    Root Integration
           |
           v
    Git Integration
           |
           v
    Testing & Docs
```

---

## Estimated Timeline

| Phase | Tasks | Effort | Timeline |
|-------|-------|--------|----------|
| Phase 1 | Spec & Docs | 1-2 hours | Day 1 |
| Phase 2 | Core Packages | 3-4 hours | Day 1-2 |
| Phase 3 | Fuzzy Command | 4-5 hours | Day 2-3 |
| Phase 4 | Remove Command | 2-3 hours | Day 3 |
| Phase 5 | Integration | 2-3 hours | Day 3-4 |
| Phase 6 | Git Enhancements | 1-2 hours | Day 4 |
| Phase 7 | Testing | 3-4 hours | Day 4-5 |
| Phase 8 | Docs | 1-2 hours | Day 5 |
| **Total** | | **17-25 hours** | **1 week** |

---

## Acceptance Criteria

- [x] `pass` without args enters fuzzy search mode
- [x] Typing filters passwords in real-time
- [x] Arrow keys navigate selection
- [x] Enter shows selected password
- [x] Esc/Ctrl+C exits
- [x] Ctrl+A, Ctrl+E, Ctrl+K work in search input
- [x] Ctrl+W (delete word) works
- [x] Tab cycles through results
- [x] Page Up/Down for page navigation
- [x] Home/End keys work
- [x] `pass rm <path>` removes file and commits to git
- [x] `pass rm` enters fuzzy search, delete on Enter
- [x] `pass rm -c <path>` copies to clipboard before deleting
- [x] `pass rm --no-commit <path>` skips git
- [x] Matching characters are highlighted
- [ ] All existing tests still pass
- [ ] All new tests pass

---

*Tasks Version: 1.0*  
*Last Updated: 2026-06-05*
