# Pass Fuzzy Search & Remove - Implementation Tasks

## Overview

This document breaks down the implementation of fuzzy search and rm features into actionable tasks.

---

## Task List

### Phase 1: Specification & Documentation
- [ ] Create spec document (COMPLETED)
- [ ] Create tasks document (IN PROGRESS)
- [ ] Update main pass-replacement-spec.md
- [ ] Update implementation checklist
- [ ] Update decision log

### Phase 2: Core Packages

#### Fuzzy Matching Package (pkg/fuzzy/)
- [ ] Create pkg/fuzzy/fuzzy.go
  - [ ] Implement `Match(query, target string) bool` - subsequence check
  - [ ] Implement `Score(query, target string) int` - ranking algorithm
  - [ ] Implement `Filter(query string, items []string) []MatchResult` - filter and sort
  - [ ] Implement `MatchResult` struct with Path, Score, MatchIndices
  - [ ] Handle case-insensitive matching
  - [ ] Handle empty query (return all)
- [ ] Create pkg/fuzzy/fuzzy_test.go
  - [ ] Test Match() with various inputs
  - [ ] Test Score() returns correct ordering
  - [ ] Test Filter() returns sorted results
  - [ ] Test edge cases (empty query, empty items, no matches)

#### Terminal Package (pkg/terminal/)
- [ ] Create pkg/terminal/terminal.go
  - [ ] Implement ANSI escape code constants
  - [ ] Implement `ClearScreen()`
  - [ ] Implement `MoveCursor(row, col int)`
  - [ ] Implement `HideCursor()` and `ShowCursor()`
  - [ ] Implement `GetSize() (width, height, error)`
  - [ ] Implement `PrintAt(text string, row, col int)`
  - [ ] Implement `ReadKey() (Key, error)`
  - [ ] Implement `Key` struct with Rune, IsArrow, ArrowDir, IsCtrl, CtrlChar
  - [ ] Implement `SupportsANSI() bool` - detect terminal capability
- [ ] Create pkg/terminal/terminal_test.go
  - [ ] Test ANSI code generation
  - [ ] Test cursor movement calculations
  - [ ] Test key parsing

### Phase 3: Fuzzy Search Command

#### Fuzzy Search UI (cmd/fuzzy.go)
- [ ] Implement fuzzy search mode entry point
- [ ] Implement main loop for fuzzy search
- [ ] Implement display rendering
  - [ ] Header "Passwords:"
  - [ ] List of matching passwords with cursor
  - [ ] Prompt "Search: " with query
  - [ ] Handle terminal resizing
- [ ] Implement query state management
  - [ ] String buffer for query
  - [ ] Cursor position in query
  - [ ] Handle all keybindings
- [ ] Implement list navigation
  - [ ] Selected index
  - [ ] Scroll position
  - [ ] Page up/down
- [ ] Implement match highlighting
  - [ ] Identify matching character positions
  - [ ] Apply terminal color codes for highlighting
- [ ] Implement fuzzy search with different invocations
  - [ ] Default: show password on Enter
  - [ ] With -c flag: copy to clipboard on Enter
  - [ ] With rm command: delete on Enter

#### Fuzzy Search Tests (cmd/fuzzy_test.go)
- [ ] Test fuzzy search initialization
- [ ] Test query input handling
- [ ] Test list filtering and sorting
- [ ] Test navigation key handling
- [ ] Test Enter key handling (different modes)
- [ ] Test exit key handling (Esc, Ctrl+C, Ctrl+D)

### Phase 4: Remove Command

#### Remove Command (cmd/rm.go)
- [ ] Implement rm command structure
  - [ ] Command definition with cobra
  - [ ] Flags: --no-commit/-n, --force/-f, --clip/-c
- [ ] Implement `removePassword(path string) error`
  - [ ] Validate path
  - [ ] Construct full file path
  - [ ] Check file exists
  - [ ] If -c flag: decrypt and copy to clipboard
  - [ ] Remove file with os.Remove()
  - [ ] If git repo exists and not --no-commit:
    - [ ] Run git rm
    - [ ] Run git commit
- [ ] Implement fuzzy search mode for rm without path
  - [ ] Same as default fuzzy search but action is delete
  - [ ] After selection: delete the file

#### Remove Tests (cmd/rm_test.go)
- [ ] Test remove with explicit path
- [ ] Test remove with fuzzy search
- [ ] Test remove with --no-commit
- [ ] Test remove with -c flag
- [ ] Test error: file not found
- [ ] Test error: permission denied
- [ ] Test: git integration

### Phase 5: Integration with Root Command

#### Update root.go
- [ ] Modify default behavior when no args
  - [ ] Detect if invoked as `pass rm` (check os.Args)
  - [ ] If no args and no command: enter fuzzy search mode
  - [ ] Pass fuzzy search mode flag (normal, clip, rm)
- [ ] Add rm command registration
- [ ] Ensure fuzzy search respects global flags (-c)

### Phase 6: Git Integration Enhancements

#### Update pkg/git/git.go
- [ ] Add `RemoveAndCommit(filePath, message string) error`
  - [ ] Remove file
  - [ ] Git rm
  - [ ] Git commit
  - [ ] Handle errors gracefully (non-fatal)
- [ ] Update tests

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

- [ ] Update README if exists
- [ ] Update pass-quick-reference.md
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

- [ ] `pass` without args enters fuzzy search mode
- [ ] Typing filters passwords in real-time
- [ ] Arrow keys navigate selection
- [ ] Enter shows selected password
- [ ] Esc/Ctrl+C exits
- [ ] Ctrl+A, Ctrl+E, Ctrl+K work in search input
- [ ] `pass rm <path>` removes file and commits to git
- [ ] `pass rm` enters fuzzy search, delete on Enter
- [ ] `pass rm -c <path>` copies to clipboard before deleting
- [ ] `pass rm --no-commit <path>` skips git
- [ ] All existing tests still pass
- [ ] All new tests pass

---

*Tasks Version: 1.0*  
*Last Updated: 2026-06-05*
