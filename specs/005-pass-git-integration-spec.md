# Pass Git Integration Specification

## Overview

This specification describes the implementation of git integration features for the `pass` password manager, providing users with visibility into the git status of their password store and the ability to perform git operations directly from the CLI and TUI.

## Status

- **Status**: Implemented ✅
- **Author**: @aasmundo
- **Created**: 2026-06-12
- **Last Updated**: 2026-06-12
- **Branch**: `feat/5-pass-git-integration`

## Background

The `pass` password manager stores passwords as GPG-encrypted files in a git repository (typically `~/.password-store/`). While `pass` already integrates with git for version control (auto-committing changes when inserting/removing passwords), there is currently no way to:

1. Check the git status of the password store
2. Push changes to a remote repository
3. Pull/fetch updates from a remote repository
4. View git status information in the TUI

This specification addresses these gaps by adding git status checking, push, and update commands, along with TUI integration.

## Goals

- Provide visibility into the git status of the password store
- Enable users to push and pull changes directly from `pass`
- Integrate git status information into the fuzzy search TUI
- Provide keyboard shortcuts for git operations in the TUI
- Maintain backward compatibility with existing functionality

## Non-Goals

- Implementing full git functionality (only status, push, fetch/pull)
- Supporting complex git workflows (rebasing, branching, etc.)
- Replacing the existing git integration for auto-committing
- Adding git history browsing or diff viewing

## User Stories

### As a pass user, I want to check the git status of my password store
So that I can see if I have uncommitted changes or if I'm behind the remote.

**Acceptance Criteria**:
- [x] `pass git` command shows current git status
- [x] Status indicates: clean/up-to-date, local changes, behind remote, ahead remote, dirty (uncommitted)
- [x] Color-coded output for easy identification (symbols: =, *, >, <, <>, !)
- [x] Works in both initialized and non-initialized git repos

### As a pass user, I want to push my changes to the remote
So that my password store is backed up.

**Acceptance Criteria**:
- [x] `pass git push` pushes changes to the configured remote
- [x] Handles push errors gracefully
- [x] Provides feedback on push success/failure
- [x] Respects existing git configuration (remote, branch)

### As a pass user, I want to update my local store from the remote
So that I have the latest passwords.

**Acceptance Criteria**:
- [x] `pass git update` fetches and merges changes from remote
- [x] Handles merge conflicts gracefully (non-fatal, with warning)
- [x] Provides feedback on update success/failure
- [x] Uses `git pull` (fetch + merge) by default

### As a pass TUI user, I want to see git status in the interface
So that I'm always aware of the sync state of my passwords.

**Acceptance Criteria**:
- [x] Git status displayed in TUI header or status bar
- [x] Status updates dynamically (on Ctrl+R, after push/update)
- [x] Shows same information as `pass git` command
- [x] Non-intrusive, doesn't interfere with password search

### As a pass TUI user, I want keyboard shortcuts for git operations
So that I can manage git without leaving the TUI.

**Acceptance Criteria**:
- [x] Keyboard shortcut to push changes (Ctrl+P)
- [x] Keyboard shortcut to update from remote (Ctrl+U)
- [x] Keyboard shortcut to refresh git status (Ctrl+R)
- [x] Visual feedback when operations complete (status updates)
- [x] Error messages displayed in TUI

## Technical Design

### Architecture

The implementation will consist of:

1. **Git Status Package** (`pkg/git/status.go`)
   - New functions for checking git status
   - Types representing git status states
   - Integration with existing `pkg/git` package

2. **Git Commands** (`cmd/git.go`)
   - New `pass git` command with subcommands
   - Integration with cobra command structure

3. **TUI Integration** (`cmd/tui/`)
   - Git status display in TUI
   - Keyboard shortcut handling
   - Status refresh mechanism

### Git Status States

```go
type GitStatus struct {
    IsGitRepo        bool
    IsClean          bool
    HasUncommitted   bool  // Dirty - unstaged or staged changes
    HasMergeConflict bool  // Merge conflicts present
    Ahead            int   // Commits ahead of remote
    Behind           int   // Commits behind remote
    Branch           string
    Remote           string
    TrackingBranch   string
    Error            error // Any error encountered
}

func GetGitStatus(storeDir string) GitStatus
```

**Status Interpretation**:
- `IsGitRepo == false`: Not a git repo
- `IsClean == true`: Up to date, no local changes
- `HasUncommitted == true`: Local changes not committed
- `Ahead > 0`: Local commits not pushed
- `Behind > 0`: Remote has commits not pulled

### Command Structure

```
pass git [command]

Commands:
  pass git           Show git status of password store
  pass git push      Push changes to remote
  pass git update    Pull changes from remote
  pass git status    Alias for 'pass git' (show status)

Flags:
  -v, --verbose   Show detailed git status information
```

### TUI Integration

#### Status Display

Git status will be displayed in the TUI header or as a status bar at the bottom:

```
┌─────────────────────────────────────────────────────────┐
│  Select password (Enter to show, Esc to cancel)            │
│  Git: main <> * (diverged, dirty)                          │
├─────────────────────────────────────────────────────────┤
│  > email/gmail.com/user                                    │
│    email/outlook.com/work                                  │
│    social/github.com                                       │
└─────────────────────────────────────────────────────────┘
Search: 
```

**Symbols**:
- `*` = uncommitted changes (dirty)
- `!` = merge conflicts
- `>` = ahead of remote
- `<` = behind remote
- `<>` = diverged (both ahead and behind)
- `=` = up to date

#### Keyboard Shortcuts

| Shortcut | Action | Description |
|----------|--------|-------------|
| `Ctrl+P` | Push | Push local changes to remote |
| `Ctrl+U` | Update | Pull changes from remote |
| `Ctrl+R` | Refresh | Refresh git status |
| `Ctrl+G` | Status | Show detailed git status |

### Error Handling

1. **No Git Repository**: If the password store is not a git repo, show a message suggesting initialization
2. **No Remote Configured**: For push/update, warn that no remote is configured
3. **Network Errors**: Display error message, don't block the TUI
4. **Merge Conflicts**: For update, warn about conflicts, provide instructions
5. **Authentication Errors**: For push, prompt for credentials if needed

### Dependencies

The implementation will use:
- Existing `pkg/git` package for git operations
- `github.com/go-git/go-git/v5` for more complex git status checking (optional, can use git CLI)
- Existing Bubble Tea framework for TUI

## Implementation Plan

### Phase 1: Core Git Status Functionality
1. [ ] Create `pkg/git/status.go` with `GetGitStatus()` function
2. [ ] Add tests for git status checking (with mock repos)
3. [ ] Create `cmd/git.go` with basic `pass git` command
4. [ ] Add tests for git command

### Phase 2: Git Operations
1. [ ] Implement `pass git push` command
2. [ ] Implement `pass git update` command
3. [ ] Add tests for push/update operations

### Phase 3: TUI Integration
1. [ ] Add git status display to TUI
2. [ ] Implement keyboard shortcuts for git operations
3. [ ] Add status refresh mechanism
4. [ ] Add tests for TUI git integration

### Phase 4: Polish and Documentation
1. [ ] Add color coding to status output
2. [ ] Add help text and documentation
3. [ ] Final testing and bug fixes

## Testing Strategy

### Unit Tests

1. **Git Status Package**
   - Test with mock git repositories (created in temp dirs)
   - Test all status states (clean, dirty, ahead, behind)
   - Test error handling (no git, no remote)
   - Use `go-git` for programmatic repo manipulation in tests

2. **Git Commands**
   - Test command parsing and argument handling
   - Test error messages
   - Mock git operations where possible

3. **TUI Integration**
   - Test status display formatting
   - Test keyboard shortcut handling
   - Test that git operations don't block TUI

### Integration Tests

1. **End-to-End Tests**
   - Create temp password store with git
   - Test full workflow: insert, check status, push, update
   - Test TUI with git status display

2. **Edge Cases**
   - Non-git password store
   - Git repo without remote
   - Git repo with merge conflicts
   - Network errors during push/pull

### Test Data

All tests will use ephemeral test data:
- Temporary directories for password stores
- Ephemeral git repositories (created and destroyed during tests)
- Mock git remotes (local bare repos for testing)
- No personal data or real credentials

Example test setup:
```go
func setupTestGitRepo(t *testing.T) (string, func()) {
    // Create temp directory
    tempDir, err := os.MkdirTemp("", "pass-git-test")
    require.NoError(t, err)
    
    // Initialize git repo
    cmd := exec.Command("git", "init")
    cmd.Dir = tempDir
    require.NoError(t, cmd.Run())
    
    // Configure user
    exec.Command("git", "config", "user.email", "test@example.com").Run()
    exec.Command("git", "config", "user.name", "Test User").Run()
    
    // Create a bare remote repo for testing
    remoteDir := filepath.Join(tempDir, "remote")
    exec.Command("git", "init", "--bare", remoteDir).Run()
    exec.Command("git", "remote", "add", "origin", remoteDir).Run()
    
    cleanup := func() {
        os.RemoveAll(tempDir)
    }
    
    return tempDir, cleanup
}
```

## Security Considerations

1. **No Credential Storage**: The implementation will not store git credentials
2. **Read-Only by Default**: Status checking only reads git information
3. **User Confirmation**: For destructive operations (push with force), require confirmation
4. **Error Isolation**: Git errors should not crash the TUI or expose sensitive information

## Backward Compatibility

1. **Existing Commands**: All existing `pass` commands continue to work
2. **Git Integration**: Existing auto-commit behavior unchanged
3. **Configuration**: No new configuration required for basic functionality
4. **Environment Variables**: Respect existing `PASSWORD_STORE_DIR` and git config

## Open Questions

1. **Should `pass git update` do a merge or rebase?**
   - Proposed: Use `git pull` (merge) by default, as it's safer for most users
   - Alternative: Add `--rebase` flag for users who prefer rebase

2. **Should we support multiple remotes?**
   - Proposed: Use the default remote (origin) and branch
   - Alternative: Add `--remote` and `--branch` flags

3. **Should push/update be available in non-interactive mode?**
   - Proposed: Yes, for scripting purposes
   - Consideration: Need to handle errors appropriately

4. **Should the TUI show a confirmation dialog for push/update?**
   - Proposed: Yes, for push (to prevent accidental pushes)
   - Proposed: No for update (pull is generally safe)
   - Alternative: Make it configurable

5. **Should we add `pass git init` to initialize a git repo?**
   - Proposed: Yes, for users who don't have git initialized
   - Alternative: Show helpful message when `pass git` is run on non-git repo

6. **Should we integrate with the existing `pkg/git` package or use `go-git`?**
   - Proposed: Extend existing `pkg/git` package with status functions
   - Consideration: `go-git` provides better programmatic access but adds dependency
   - Decision: Use git CLI for now (already a dependency), can add `go-git` later if needed

## Success Criteria

- [x] `pass git` command works and shows correct status
- [x] `pass git push` successfully pushes changes
- [x] `pass git update` successfully pulls changes
- [x] TUI shows git status information
- [x] TUI keyboard shortcuts work for git operations (defined, need runtime verification)
- [x] All tests pass
- [x] No personal data in tests or code
- [x] Documentation updated (README, specs)

## Appendix

### Git Status Output Examples

```
# Clean, up to date
$ pass git
Git status: clean (main)
Up to date with origin/main

# Local changes, not committed
$ pass git
Git status: dirty (main)
Uncommitted changes:
  modified:   email/gmail.com/user.gpg

# Ahead of remote
$ pass git
Git status: clean (main)
2 commits ahead of origin/main

# Behind remote
$ pass git
Git status: clean (main)
3 commits behind origin/main

# Diverged
$ pass git
Git status: clean (main)
2 commits ahead, 1 commit behind origin/main

# Not a git repo
$ pass git
Error: password store is not a git repository
Hint: Run 'pass git init' to initialize git
```

### TUI Status Display Examples

```
# Clean
Git: main = (up to date)

# Dirty
Git: main * (uncommitted changes)

# Ahead
Git: main > (ahead)

# Behind
Git: main < (behind)

# Diverged
Git: main <> (diverged)

# Multiple issues
Git: main <> * (diverged, uncommitted)
```
