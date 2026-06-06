# Pass Edit Command and Bug Fixes Specification

## Overview

This document specifies the implementation of the `pass edit` command and fixes for existing bugs in the `pass` tool.

### Issues Addressed

1. **Bug**: `pass insert <secret>` overwrites existing secrets - should fail with error
2. **Bug**: `pass ls` lists both directories and secret files - should only list secret files
3. **Feature**: `pass edit <secret>` to edit existing secrets
4. **Enhancement**: `pass edit` without arguments uses fuzzy finder to select secret to edit

## User Requirements

### Must Have

- [ ] `pass insert <secret>` must NOT overwrite existing secrets - return error instead
- [ ] `pass ls` must only list `.gpg` files (secret files), not directories
- [ ] `pass edit <path>` decrypts the secret, loads it in `$EDITOR`, re-encrypts on save
- [ ] `pass edit` without arguments enters fuzzy search mode to select secret to edit
- [ ] Edit command preserves git history with proper commit messages
- [ ] All operations respect existing flags (e.g., `--no-commit`)

### Should Have

- [ ] Edit command supports `--no-commit` flag to skip git commit
- [ ] Edit command validates that the file exists before attempting to edit
- [ ] Edit command handles editor errors gracefully
- [ ] Edit command preserves file permissions

### Nice to Have (Future)

- [ ] Support for editing multiple secrets at once
- [ ] Preview of secret content before editing
- [ ] Syntax highlighting for common secret types (JSON, YAML, etc.)

---

## Detailed Specifications

### 1. Fix: Insert Command Must Not Overwrite

**Current Behavior:**
```bash
pass insert email/gmail.com  # Creates email/gmail.com.gpg
pass insert email/gmail.com  # OVERWRITES existing file - BUG!
```

**Required Behavior:**
```bash
pass insert email/gmail.com  # Creates email/gmail.com.gpg
pass insert email/gmail.com  # Error: pass: email/gmail.com: Already exists
```

**Implementation:**
- Before creating/encrypting the file, check if the destination file already exists
- If it exists, return error: `pass: <path>: Already exists`
- Exit code: 1

**Error Message Format:**
```
pass: <path>: Already exists
```

### 2. Fix: List Command Must Only Show Files

**Current Behavior:**
```bash
$ pass ls
email
social
email/gmail.com
social/twitter.com
```
(Shows both directories and files)

**Required Behavior:**
```bash
$ pass ls
email/gmail.com
social/twitter.com
```
(Only shows `.gpg` files, with `.gpg` extension stripped)

**Implementation:**
- Modify `listPasswords()` function in `cmd/ls.go`
- Skip all directories in the output
- Only include `.gpg` files
- Strip `.gpg` extension from display

**Flags Behavior:**
- `--dirs-only` / `-d`: Show only directories (existing behavior)
- `--files-only` / `-f`: Show only files (default behavior for `pass ls`)
- `--recursive` / `-r`: Show full paths (existing behavior)

### 3. New Feature: Edit Command

**Synopsis:**
```
pass edit [OPTIONS] [<path>]
```

**Description:**
The edit command allows users to edit existing password files. It decrypts the file, loads the plaintext content into the user's default editor (from `$EDITOR` environment variable, or `vi` on Unix, `notepad` on Windows as fallback), and re-encrypts the file when the editor exits.

**Arguments:**
- `<path>`: Optional. The path of the password to edit. If omitted, enters fuzzy search mode.

**Options:**
- `-n, --no-commit`: Skip git commit after editing
- `-f, --force`: Alias for `--no-commit`

**Behavior:**

#### With explicit path:
1. Normalize the path (strip `.gpg` extension if present, add if missing)
2. Check if file exists
   - If not: `pass: <path>: No such file or directory` (exit code 1)
3. Decrypt the file to a temporary file
4. Open the temporary file in `$EDITOR`
5. Wait for editor to exit
6. Check editor exit code
   - If non-zero: `pass: editor exited with status <code>` (exit code 2)
7. Read the modified content from temporary file
8. Validate content is not empty (optional - allow empty passwords?)
9. Re-encrypt the content back to the original file
10. Securely delete the temporary file
11. If git repo exists and `--no-commit` not specified:
    - Run `git add <path>.gpg`
    - Run `git commit -m "Edit <path>"`
    - If git commands fail, display warning but continue
12. Display: "Password updated successfully."
13. Exit with code 0

#### Without explicit path (fuzzy search mode):
1. Enter fuzzy search mode (same as `pass` without args)
2. As user types, filter list of passwords
3. When user presses Enter on a selection:
   - Proceed with editing that password
4. If user presses Esc/Ctrl+C: exit without action

**Editor Selection:**
1. Check `$EDITOR` environment variable
2. If empty on Windows: use `notepad`
3. If empty on Unix: use `vi`
4. If editor command fails: return error

**Temporary File:**
- Create in system temp directory with prefix `pass-edit-`
- Use `.tmp` extension
- Set permissions to 0600 (readable/writable only by owner)
- Securely delete after use (overwrite before removal)

**Exit Codes:**
| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | File not found |
| 2 | Editor error |
| 3 | Decryption failed |
| 4 | Encryption failed |
| 5 | Git operation failed (non-fatal) |

**Examples:**
```bash
# Edit specific password
pass edit email/gmail.com/user

# Edit with fuzzy search
pass edit
# User types: gmail
# Selects: email/gmail.com/user
# Press Enter to edit

# Edit without git commit
pass edit --no-commit social/twitter.com/token

# Edit with force flag (alias for --no-commit)
pass edit -f banking/chase.com/pin
```

### 4. Fuzzy Search Mode for Edit

**Integration with existing TUI:**
- Add new mode: `FuzzyModeEdit` to the `FuzzySearchMode` enum
- Update TUI to handle edit mode
- When Enter is pressed in edit mode, call `editPassword(selected)`
- Update `RunInteractiveFuzzySearch` to support edit mode

**TUI Display for Edit Mode:**
- Header: "Select password to edit (Enter to edit, Esc to cancel):"
- Prompt: "Edit: "
- Help text should mention edit action

---

## Implementation Details

### 4.1 File Structure

```
pass/
├── cmd/
│   ├── edit.go          # NEW: Edit command
│   ├── edit_test.go     # NEW: Tests for edit command
│   ├── insert.go        # MODIFIED: Add overwrite check
│   ├── ls.go            # MODIFIED: Only list files
│   ├── root.go          # MODIFIED: Register edit command
│   ├── fuzzy.go         # MODIFIED: Add FuzzyModeEdit
│   └── tui/
│       └── fuzzy.go     # MODIFIED: Support edit mode
├── specs/
│   └── pass-edit/       # NEW: This spec
│       └── spec.md
├── tests/
│   ├── edit_test.go     # NEW: Integration tests
│   └── ...
└── README.md            # MODIFIED: Document edit command
```

### 4.2 Dependencies

**New dependencies:**
- None - uses existing packages

**Existing dependencies used:**
- `pkg/gpg` - for decryption and encryption
- `pkg/filesystem` - for file operations and temp files
- `pkg/git` - for git integration
- `cmd/tui` - for fuzzy search TUI

### 4.3 Editor Handling

```go
// getEditor returns the editor command to use
func getEditor() string {
    editor := os.Getenv("EDITOR")
    if editor != "" {
        return editor
    }
    // Platform-specific defaults
    if runtime.GOOS == "windows" {
        return "notepad"
    }
    return "vi"
}

// openInEditor opens a file in the editor and waits for it to close
func openInEditor(filePath string) error {
    editor := getEditor()
    cmd := exec.Command(editor, filePath)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    if err := cmd.Run(); err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            return fmt.Errorf("editor exited with status %d", exitErr.ExitCode())
        }
        return fmt.Errorf("failed to open editor: %v", err)
    }
    return nil
}
```

### 4.4 Edit Command Implementation

```go
// editCmd represents the edit command
var editCmd = &cobra.Command{
    Use:   "edit [OPTIONS] [<path>]",
    Short: "Edit a password",
    Long: `Edit an existing password. Decrypts the password, opens it in your editor,
and re-encrypts it when you save.

If a path is provided, edits that specific password.
If no path is provided, will enter interactive fuzzy search mode to select a password.`,
    Args: cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        var path string
        if len(args) > 0 {
            path = args[0]
        }
        
        noCommit, _ := cmd.Flags().GetBool("no-commit")
        force, _ := cmd.Flags().GetBool("force")
        
        if force {
            noCommit = true
        }
        
        if path != "" {
            return editPassword(path, noCommit)
        }
        
        // No path provided - enter fuzzy search mode
        return runEditFuzzySearch(noCommit)
    },
}
```

### 4.5 Insert Command Fix

```go
// insertPassword inserts a new password
func insertPassword(path string) error {
    // ... existing code ...
    
    // NEW: Check if file already exists
    if _, err := os.Stat(fullPath); err == nil {
        return fmt.Errorf("pass: %s: Already exists", path)
    }
    
    // ... rest of existing code ...
}
```

### 4.6 List Command Fix

```go
// listPasswords lists all passwords in the store
func listPasswords(subpath string) error {
    // ... existing setup code ...
    
    // Walk the directory tree
    var results []string
    err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        // Skip .git directory
        if info.IsDir() && info.Name() == ".git" {
            return filepath.SkipDir
        }
        
        // Skip the base path itself
        relPath, err := filepath.Rel(storeDir, path)
        if err != nil {
            return err
        }
        relPath = filesystem.NormalizePathForDisplay(relPath)
        if relPath == "." {
            return nil
        }
        
        // MODIFIED: Only include .gpg files (not directories)
        if !info.IsDir() && strings.HasSuffix(info.Name(), ".gpg") {
            passwordPath := strings.TrimSuffix(relPath, ".gpg")
            results = append(results, passwordPath)
        }
        
        return nil
    })
    
    // ... rest of existing code ...
}
```

---

## Error Handling

### Insert Errors

| Scenario | Error Message | Exit Code |
|----------|---------------|-----------|
| File already exists | `pass: <path>: Already exists` | 1 |
| Directory creation fails | `pass: failed to create directory: <error>` | 1 |
| Temp file creation fails | `pass: failed to create temp file: <error>` | 1 |
| GPG encryption fails | `pass: GPG encryption failed: <error>` | 1 |

### Edit Errors

| Scenario | Error Message | Exit Code |
|----------|---------------|-----------|
| File not found | `pass: <path>: No such file or directory` | 1 |
| Decryption fails | `pass: decryption failed: <error>` | 3 |
| Editor not found | `pass: editor not found: <editor>` | 2 |
| Editor exits with error | `pass: editor exited with status <code>` | 2 |
| Encryption fails | `pass: GPG encryption failed: <error>` | 4 |
| Git add fails | Warning to stderr | 5 (non-fatal) |
| Git commit fails | Warning to stderr | 5 (non-fatal) |

### List Errors

| Scenario | Error Message | Exit Code |
|----------|---------------|-----------|
| Store doesn't exist | (Create it automatically) | 0 |
| Permission denied | `pass: failed to walk directory: permission denied` | 1 |

---

## Testing Strategy

### 5.1 Unit Tests

**Insert command:**
- Test that insert fails when file already exists
- Test that insert succeeds when file doesn't exist
- Test error message format

**Edit command:**
- Test path normalization
- Test file existence check
- Test temp file creation and cleanup
- Test editor command construction
- Test re-encryption
- Test git integration

**List command:**
- Test that only files are listed (not directories)
- Test with nested directories
- Test with various flag combinations

### 5.2 Integration Tests

**Insert + Edit workflow:**
1. Insert a password
2. Edit that password
3. Verify content was updated
4. Verify git history

**Fuzzy search + Edit:**
1. Insert multiple passwords
2. Run `pass edit` without args
3. Select a password via fuzzy search
4. Edit it
5. Verify changes

### 5.3 End-to-End Tests

- Full workflow: insert → ls → edit → show
- Full workflow: insert → edit (fuzzy) → verify
- Git history verification

---

## Compatibility

### Unix pass Compatibility

| Feature | Unix pass | This implementation |
|---------|-----------|-------------------|
| `pass edit <path>` | ✓ | ✓ |
| Edit without args | ✗ | ✓ (with fuzzy search) |
| Git commit on edit | ✓ | ✓ (default) |
| `--no-commit` flag | Partial | ✓ |

### Editor Compatibility

- Supports any editor that can be invoked as a command
- Works with GUI editors (Notepad, VS Code, etc.)
- Works with terminal editors (vi, vim, emacs, nano, etc.)
- Handles editor exit codes properly

---

## Open Questions

### OQ-001: Should empty passwords be allowed?
**Status**: OPEN   
**Proposal**: Yes, allow empty passwords. Some users may want to store empty values.

### OQ-002: Should we validate the editor exists before opening?
**Status**: OPEN   
**Proposal**: Yes, check if editor command exists in PATH before attempting to open.

### OQ-003: Should we support a custom editor flag (e.g., `--editor`)?
**Status**: OPEN   
**Proposal**: Not in initial implementation. Users can set `$EDITOR` environment variable.

---

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `EDITOR` | Editor to use for editing passwords | `notepad` (Windows), `vi` (Unix) |
| `PASSWORD_STORE_DIR` | Password store location | `~/.password-store` |
| `PASS_GPG_ID` | GPG recipient key ID | Default key |
| `PASS_NO_COMMIT` | Skip git commits globally | false |

---

## References

- [Unix pass edit command](https://git.zx2c4.com/password-store/tree/man/pass.1.md)
- [Editor environment variable](https://en.wikipedia.org/wiki/EDITOR)
- [Temporary files in Go](https://pkg.go.dev/os#CreateTemp)

---

*Document Version: 1.0*   
*Last Updated: 2026-06-06*   
*Author: Mandu*   
*Status: Approved for Implementation*
