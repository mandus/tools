# Pass Edit Implementation Tasks

## Overview

This document tracks the implementation tasks for the pass edit command and bug fixes as specified in `spec.md`.

## Task Breakdown

### Phase 1: Specification and Planning ✅
- [x] Create spec document (`specs/pass-edit/spec.md`)
- [x] Create tasks document (`specs/pass-edit/tasks.md`)
- [x] Review existing codebase and understand architecture

### Phase 2: Bug Fixes

#### 2.1 Fix Insert Overwrite Bug
- [ ] Modify `insertPassword()` in `cmd/insert.go`
  - [ ] Add check for existing file before creating new one
  - [ ] Return error: `pass: <path>: Already exists`
  - [ ] Update error handling
- [ ] Update tests in `cmd/insert_test.go`
  - [ ] Add test for overwrite prevention
  - [ ] Add test for error message format

#### 2.2 Fix List Command to Only Show Files
- [ ] Modify `listPasswords()` in `cmd/ls.go`
  - [ ] Skip directories in output (unless --dirs-only flag)
  - [ ] Only include `.gpg` files
  - [ ] Ensure `.gpg` extension is stripped from display
- [ ] Update tests in `cmd/ls_test.go`
  - [ ] Add test for file-only listing
  - [ ] Add test for directory filtering

### Phase 3: Edit Command Implementation

#### 3.1 Create Edit Command
- [ ] Create `cmd/edit.go`
  - [ ] Define `editCmd` cobra command
  - [ ] Implement `editPassword()` function
  - [ ] Add editor detection (`$EDITOR`, platform defaults)
  - [ ] Add temp file handling
  - [ ] Add decryption → edit → re-encryption flow
  - [ ] Add git integration
  - [ ] Add error handling
- [ ] Register command in `cmd/root.go`
  - [ ] Call `addEditCmd()` in `Execute()`

#### 3.2 Fuzzy Search Integration
- [ ] Update `FuzzySearchMode` enum in `cmd/fuzzy.go` or `cmd/tui/fuzzy.go`
  - [ ] Add `FuzzyModeEdit` constant
- [ ] Update TUI in `cmd/tui/models.go`
  - [ ] Handle edit mode in render
  - [ ] Update header and prompt text
- [ ] Update `RunInteractiveFuzzySearch` in `cmd/tui/fuzzy.go`
  - [ ] Add case for `FuzzyModeEdit`
  - [ ] Call appropriate edit function
- [ ] Update `cmd/fuzzy.go` (old implementation)
  - [ ] Add `FuzzyModeEdit` support
  - [ ] Update `RunInteractiveFuzzySearch`

#### 3.3 Helper Functions
- [ ] Create editor utility functions
  - [ ] `getEditor()` - returns editor command
  - [ ] `openInEditor(filePath string) error` - opens file and waits
  - [ ] `validateEditor(editor string) bool` - checks if editor exists

### Phase 4: Testing

#### 4.1 Unit Tests
- [ ] Create `cmd/edit_test.go`
  - [ ] Test path normalization
  - [ ] Test file existence check
  - [ ] Test editor detection
  - [ ] Test temp file operations
  - [ ] Test error scenarios
- [ ] Update `cmd/insert_test.go`
  - [ ] Test overwrite prevention
- [ ] Update `cmd/ls_test.go`
  - [ ] Test file-only listing

#### 4.2 Integration Tests
- [ ] Create `tests/edit_test.go`
  - [ ] Test insert → edit workflow
  - [ ] Test fuzzy search → edit workflow
  - [ ] Test git integration
  - [ ] Test various editor scenarios

#### 4.3 Manual Testing
- [ ] Test on Windows with Notepad
- [ ] Test on Unix with vi/vim
- [ ] Test with custom `$EDITOR`
- [ ] Test with nested directories
- [ ] Test with special characters in passwords

### Phase 5: Documentation

#### 5.1 Update README
- [ ] Update `pass/README.md`
  - [ ] Add edit command to usage examples
  - [ ] Add edit command to feature list
  - [ ] Document editor requirements
  - [ ] Document flags

#### 5.2 Update Command Documentation
- [ ] Update `pass/docs/` if applicable
  - [ ] Add edit command documentation
  - [ ] Update existing docs for insert and ls fixes

#### 5.3 Update Root README
- [ ] Update `README.md` in repository root
  - [ ] Mention pass edit command
  - [ ] Link to pass documentation

### Phase 6: Finalization

#### 6.1 Code Review
- [ ] Review all changes for consistency
- [ ] Check error handling
- [ ] Verify git integration
- [ ] Check cross-platform compatibility

#### 6.2 Cleanup
- [ ] Remove any debug code
- [ ] Fix any linting issues
- [ ] Ensure consistent code style

#### 6.3 Commit
- [ ] Stage all changes
- [ ] Write commit message following gitmoji conventions
- [ ] Commit to branch

## Priority Order

1. **High Priority (Must Have)**
   - Fix insert overwrite bug
   - Fix ls to only show files
   - Implement basic edit command
   - Add fuzzy search support for edit

2. **Medium Priority (Should Have)**
   - Add tests for all new functionality
   - Update documentation
   - Add editor validation

3. **Low Priority (Nice to Have)**
   - Custom editor flag
   - Empty password validation
   - Additional error scenarios

## Estimated Time

| Task | Estimate |
|------|----------|
| Spec and planning | 1 hour |
| Fix insert bug | 30 minutes |
| Fix ls bug | 30 minutes |
| Create edit command | 2 hours |
| Fuzzy search integration | 1 hour |
| Unit tests | 1 hour |
| Integration tests | 1 hour |
| Documentation | 1 hour |
| Finalization | 30 minutes |
| **Total** | **8.5 hours** |

## Dependencies

- Go 1.20+
- Existing pass tool packages
- GPG installed and configured
- Git installed and configured

## Blockers

None identified at this time.

## Notes

- All changes should follow the existing code style and conventions
- Error messages should match the format used in existing commands
- Git integration should be consistent with insert and rm commands
- Cross-platform compatibility is required (Windows, Linux, macOS)
