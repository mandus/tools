# Pass Replacement Specification

## Overview

This document specifies a Windows-compatible replacement for the `pass` password manager (https://www.passwordstore.org/). The tool provides basic password management functionality by wrapping GnuPG (GPG) encryption/decryption operations on files stored in a hierarchical directory structure.

**Primary Goal**: Provide a drop-in replacement for basic `pass` functionality on Windows 11/ARM64 where the native `pass` tool is unavailable.

## Background

The Unix `pass` tool:
- Stores passwords as GPG-encrypted files in `~/.password-store/`
- Uses directory hierarchy to organize passwords (e.g., `email/gmail.com/username`)
- Integrates with `git` for version control
- Supports clipboard integration
- Provides search and listing capabilities

## User Requirements

### Must Have
- [ ] Insert new passwords with `pass insert <path>`
- [ ] Retrieve passwords with `pass <path>` (output to stdout)
- [ ] Retrieve passwords with `pass -c <path>` (copy to clipboard)
- [ ] List all passwords with `pass ls`
- [ ] Search passwords with `pass find <string>`
- [ ] Automatic git commit on insert (local commit only, no push)
- [ ] GPG encryption/decryption
- [ ] Store data in `~/.password-store/`

### Should Have
- [ ] Handle GPG passphrase prompting gracefully
- [ ] Support multi-line passwords
- [ ] Proper error handling for missing files/paths
- [ ] Windows path handling (convert `/` to `\` internally)

### Nice to Have (Future)
- [ ] `pass edit <path>` - Edit password in editor
- [ ] `pass generate <path> [length]` - Generate random password
- [x] `pass rm <path>` - Remove password
- [ ] `pass mv <old> <new>` - Move/rename password
- [ ] `pass git <args>` - Pass-through git commands
- [ ] Tree display with `pass tree`
- [ ] Password expiration tracking

## Architecture

### File Structure

```
~/.password-store/
├── .git/
│   └── ...
├── email/
│   └── gmail.com/
│       └── myemail@gmail.com.gpg
├── social/
│   └── twitter.com/
│       └── username.gpg
└── banking/
    └── chase.com/
        └── account.gpg
```

### Command Structure

```
pass [OPTIONS] COMMAND [ARGS...]
```

## Detailed Specifications

---

## 1. Configuration

### 1.1 Password Store Location

**Requirement**: The password store MUST be located at `~/.password-store/` on Windows.

**Implementation Notes**:
- On Windows, `~` resolves to `%USERPROFILE%`
- Default path: `%USERPROFILE%\.password-store`
- Support `PASSWORD_STORE_DIR` environment variable override

### 1.2 GPG Configuration

**Requirement**: The tool MUST use GnuPG for encryption/decryption.

**Assumptions**:
- GPG is installed and available in PATH as `gpg`
- User has a GPG key pair configured
- GPG agent handles passphrase caching

**Implementation Notes**:
- Use `gpg --encrypt --recipient <key-id>` for encryption
- Use `gpg --decrypt` for decryption
- Default recipient: use default GPG key (no explicit recipient needed if only one key)
- Support `PASS_GPG_ID` environment variable for specific recipient

### 1.3 Git Integration

**Requirement**: The password store MUST be a git repository, with automatic commits on insert.

**Implementation Notes**:
- Initialize git repo if `.password-store/.git` doesn't exist
- On insert: `git add <file>` and `git commit -m "Add <path>"`
- Commit author: use git config or default to system user
- No automatic push (user must push manually)

---

## 2. Command Line Interface

### 2.1 Global Options

```
-c, --clip[board]
    Copy the password to the clipboard instead of printing to stdout.
    Type: Boolean flag
    Applies to: show command, rm command
    
-v, --verbose
    Show verbose output for debugging.
    Type: Boolean flag
    Applies to: All commands
    
--version
    Display version information and exit.
    Type: Boolean flag
    
--help, -h
    Display help information and exit.
    Type: Boolean flag
    
-i, --interactive
    Force interactive fuzzy search mode (same as running `pass` without args).
    Type: Boolean flag
    Applies to: All commands
```

### 2.2 Command: insert

**Synopsis**:
```
pass insert [OPTIONS] <path>
```

**Description**: Insert a new password into the store. Prompts for password twice (for confirmation), encrypts it, and saves to file.

**Arguments**:
- `<path>`: The path where the password will be stored (e.g., `email/gmail.com/myemail`)
  - Type: String, required
  - Must be a valid path (no illegal characters for filesystem)
  - Path separators can be `/` or `\` (normalize to OS separator internally)

**Options**:
- `-e, --echo`: Echo the password while typing (default: hidden)
- `-m, --multiline`: Allow multi-line password input
- `--no-commit`: Skip git commit after insertion

**Behavior**:
1. Validate path is not empty
2. Check if password already exists at path (warn user)
3. Prompt: "Enter password for <path>:" (hidden input)
4. Prompt: "Retype password for <path>:" (hidden input)
5. If passwords don't match, display error and exit with code 1
6. If passwords match and are empty, display error and exit with code 1
7. Create directory structure if it doesn't exist
8. Write password to temporary file
9. Encrypt temporary file with GPG
10. Save encrypted file as `<password-store-dir>/<path>.gpg`
11. Delete temporary file securely (overwrite before delete)
12. If git repo exists, run `git add <path>.gpg` and `git commit -m "Add <path>"`
13. Display: "Password inserted successfully."

**Exit Codes**:
- 0: Success
- 1: Passwords don't match or are empty
- 2: GPG encryption failed
- 3: File write failed
- 4: Git commit failed (only if --no-commit not specified)

**Examples**:
```bash
pass insert email/gmail.com/myemail@gmail.com
pass insert banking/chase.com/account
pass insert --echo social/twitter.com/username
```

---

### 2.3 Command: show (default command)

**Synopsis**:
```
pass [OPTIONS] <path>
pass show [OPTIONS] <path>
```

**Description**: Retrieve and display a password from the store. If `-c` flag is set, copy to clipboard instead.

**Arguments**:
- `<path>`: The path of the password to retrieve
  - Type: String, required
  - Can omit `.gpg` extension (automatically appended if not present)

**Options**:
- `-c, --clip[board]`: Copy to clipboard instead of stdout
- `-o, --output <file>`: Write password to file instead of stdout
- `-q, --quiet`: Suppress warnings/errors (still exit with error code)
- `--no-newline`: Don't output trailing newline

**Behavior**:
1. Validate path is not empty
2. Construct file path: `<password-store-dir>/<path>.gpg` (add `.gpg` if not present)
3. Check if file exists
   - If not, display error "`pass: <path>: No such file or directory" and exit with code 1
4. Decrypt file with GPG
   - GPG will prompt for passphrase if needed (gpg-agent handles caching)
5. If decryption fails, display error and exit with code 2
6. If `-c` flag:
   - Copy decrypted content to clipboard
   - Display: "Copied <path> to clipboard."
   - Clear clipboard after 45 seconds (optional, see security considerations)
7. If `-o` flag:
   - Write decrypted content to specified file
8. Otherwise:
   - Print decrypted content to stdout
   - If `--no-newline`, don't add trailing newline
9. Exit with code 0

**Exit Codes**:
- 0: Success
- 1: File not found
- 2: GPG decryption failed
- 3: Clipboard copy failed
- 4: File write failed (for -o option)

**Examples**:
```bash
pass email/gmail.com/myemail@gmail.com
pass -c email/gmail.com/myemail@gmail.com
pass show banking/chase.com/account
pass email/gmail.com/myemail@gmail.com.gpg
```

---

### 2.4 Command: ls

**Synopsis**:
```
pass ls [OPTIONS] [<subpath>]
```

**Description**: List all passwords in the store, optionally filtered by a subpath.

**Arguments**:
- `<subpath>`: Optional subpath to list (e.g., `email/`)
  - Type: String, optional
  - If provided, list only passwords under this path

**Options**:
- `-r, --recursive`: Show full paths (default for pass compatibility)
- `-d, --dirs-only`: List only directories
- `-f, --files-only`: List only password files
- `-l, --long`: Long listing format with metadata

**Behavior**:
1. Construct base path: `<password-store-dir>/<subpath>` (if subpath provided)
2. If base path doesn't exist, display error and exit with code 1
3. Recursively walk directory tree starting from base path
4. For each `.gpg` file found:
   - Strip `.gpg` extension
   - Strip `<password-store-dir>/` prefix
   - If subpath was provided, strip that prefix too
5. Sort results alphabetically
6. Display each path on its own line

**Output Format**:
- Default: One path per line, relative to password store root or subpath
- With `-r` (default): Full paths from root
- With `-d`: Only directory names
- With `-f`: Only file names (without .gpg)

**Exit Codes**:
- 0: Success
- 1: Password store directory not found or subpath doesn't exist

**Examples**:
```bash
pass ls
pass ls email/
pass ls -d
pass ls -f
```

**Example Output**:
```
email/gmail.com/myemail@gmail.com
banking/chase.com/account
social/twitter.com/username
```

---

### 2.5 Command: find

**Synopsis**:
```
pass find [OPTIONS] <search-string>
```

**Description**: Search for passwords containing the search string anywhere in their path.

**Arguments**:
- `<search-string>`: String to search for
  - Type: String, required
  - Case-sensitive by default

**Options**:
- `-i, --ignore-case`: Case-insensitive search
- `-l, --long`: Show full paths (default)
- `-n, --name-only`: Show only matching part of path

**Behavior**:
1. Validate search string is not empty
2. Walk entire password store directory
3. For each `.gpg` file:
   - Get full path relative to password store root (without `.gpg`)
   - Check if path contains search string
     - If `-i` flag: case-insensitive comparison
4. Collect all matching paths
5. Sort results alphabetically
6. Display each matching path on its own line

**Exit Codes**:
- 0: Success (even if no matches found)
- 1: Search string is empty

**Examples**:
```bash
pass find gmail
pass find -i GMAIL
pass find chase
```

**Example Output**:
```
email/gmail.com/myemail@gmail.com
email/gmail.com/work@gmail.com
```

---

### 2.6 Command: rm

**Synopsis**:
```
pass rm [OPTIONS] [<path>]
```

**Description**: Remove a password from the store. If path is provided, removes that specific password. If no path is provided, enters interactive fuzzy search mode where the user can select a password to remove.

**Arguments**:
- `<path>`: The path of the password to remove
  - Type: String, optional
  - Can omit `.gpg` extension (automatically appended if not present)
  - Must be a valid path to an existing `.gpg` file

**Options**:
- `-n, --no-commit`: Skip git commit after removal
- `-f, --force`: Alias for --no-commit
- `-c, --clip[board]`: Copy password to clipboard before deleting

**Behavior**:

**With explicit path:**
1. Validate path is not empty
2. Construct file path (add `.gpg` if not present)
3. Check if file exists
   - If not, display error `pass: <path>: No such file or directory` and exit with code 1
4. If `-c` flag is set:
   - Decrypt the file and copy password to clipboard
   - Start clipboard clear timer (if enabled)
5. Remove the file with `os.Remove()`
6. If git repo exists and `--no-commit` not specified:
   - Run `git rm <path>.gpg`
   - Run `git commit -m "Remove <path>"`
   - If git commands fail, display warning but continue (non-fatal)
7. Display: "Password removed successfully."
8. Exit with code 0

**Without explicit path (fuzzy search mode):**
1. Enter interactive fuzzy search mode
2. As user types, filter list of passwords using fuzzy matching
3. Display matching passwords with best match selected
4. User navigates with arrow keys, types to filter
5. When user presses Enter on a selection:
   - If `-c` flag: decrypt and copy to clipboard first
   - Remove the selected password file
   - If git repo exists and not `--no-commit`: commit removal
   - Display success message
   - Exit fuzzy search mode
6. If user presses Esc/Ctrl+C: exit without action, exit with code 0

**Exit Codes**:
- 0: Success
- 1: File not found
- 2: Permission denied
- 3: Git operation failed (non-fatal, file still removed)

**Examples**:
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

### 2.7 Interactive Fuzzy Search Mode

**Synopsis**:
```
pass
pass --interactive
pass -i
```

**Description**: When invoked without any arguments or with the `--interactive` flag, `pass` enters an interactive fuzzy search mode. This allows users to quickly find and select secrets by typing partial matches, similar to the fzf tool.

**Fuzzy Matching**:
- Characters in the query must appear in the target path in the same order (subsequence match)
- Characters do NOT need to be consecutive
- Matching is case-insensitive
- Example: Query `twt` matches `social/twitter.com/admin` (t-w-t in order)

**Display**:
```
Passwords:
  email/gmail.com/user
> social/twitter.com/admin
  bank/chase.com/account
  
Search: tw
```

**Components**:
- **Header**: "Passwords:" label
- **List**: Scrollable list of matching passwords
- **Cursor**: `>` prefix indicates currently selected item
- **Prompt**: "Search: " at bottom with user's query
- **Matching characters**: Visually distinguished (highlighted if terminal supports it)

**Action on Selection**:
The action taken when Enter is pressed depends on how fuzzy search was invoked:

| Invocation | Action |
|------------|--------|
| `pass` or `pass -i` | Show the selected password (to stdout) |
| `pass -c` or `pass --clip` | Copy selected password to clipboard |
| `pass rm` | Delete the selected password file |
| `pass rm -c` | Copy to clipboard, then delete |

**Keybindings**:

| Key | Action |
|-----|--------|
| Any printable character | Add to query, re-filter list |
| Backspace | Remove last character from query |
| Delete | Remove character under cursor from query |
| Ctrl+A | Move cursor to start of query |
| Ctrl+E | Move cursor to end of query |
| Ctrl+K | Delete from cursor to end of query |
| Ctrl+L | Clear entire query |
| Ctrl+W | Delete word before cursor |
| ↑ (Up Arrow) | Move selection up by 1 |
| ↓ (Down Arrow) | Move selection down by 1 |
| Page Up | Move selection up by page height |
| Page Down | Move selection down by page height |
| Home | Move to first item |
| End | Move to last item |
| Enter | Select current item (perform action based on invocation) |
| Esc | Exit fuzzy search mode, return to shell |
| Ctrl+C | Exit fuzzy search mode, return to shell |
| Ctrl+D | Exit fuzzy search mode (EOF), return to shell |
| Tab | Toggle between search input and list navigation |
| ← (Left Arrow) | Move cursor left in query |
| → (Right Arrow) | Move cursor right in query |

**Filtering Behavior**:
- **Empty query**: Show all passwords
- **No matches**: Show "No matches found" message
- **Git ignore**: Only show `.gpg` files that would be tracked by git (exclude `.git/` directory)
- **File types**: Only show `.gpg` files
- **Sorting**: Best matches first (by score), then alphabetically

**Exit Codes**:
- 0: Success (selected and performed action, or exited normally)
- 1: Error occurred
- 2: No terminal available (when piped)

---

## 3. File Operations

### 3.1 File Naming

**Requirement**: Password files MUST use `.gpg` extension.

**Rules**:
- Input path `email/gmail.com/myemail` → File: `.password-store/email/gmail.com/myemail.gpg`
- Input path with `.gpg` extension: strip it before adding (idempotent)
- Invalid characters in path: Replace or reject based on OS

**Character Handling**:
- Windows reserved characters: `< > : " | ? *`
- Strategy: Replace with underscore `_` and warn user
- Forward slash `/` → Convert to OS path separator internally

### 3.2 File Encryption

**Process**:
```
1. Create temporary file with random name in temp directory
2. Write password content to temporary file (UTF-8 encoding)
3. Run: gpg --encrypt --recipient <recipient> --armor --output <dest>.gpg <temp-file>
   - --armor: ASCII-armored output (optional, but pass uses it)
   - --recipient: Use PASS_GPG_ID if set, otherwise use default key
4. Securely delete temporary file
5. Verify encrypted file exists and is readable
```

**GPG Options**:
- `--batch`: For non-interactive mode (but we want passphrase prompt)
- `--yes`: Assume yes to prompts
- `--no-tty`: Don't require TTY (for scripting)

### 3.3 File Decryption

**Process**:
```
1. Check file exists and is readable
2. Run: gpg --decrypt <file>.gpg
3. Capture stdout (decrypted content)
4. Handle stderr (GPG warnings, passphrase prompts)
5. Return decrypted content
```

**Note**: GPG handles passphrase prompting automatically via gpg-agent.

---

## 4. Git Integration

### 4.1 Repository Initialization

**Requirement**: The password store MUST be a git repository.

**Behavior**:
- On first insert (or explicitly via `pass git init`):
  - Check if `.password-store/.git` exists
  - If not, run `git init` in `.password-store`
  - Run `git config user.name` and `git config user.email` if not set
  - Create initial commit if no commits exist

### 4.2 Automatic Commits

**Requirement**: After successful insert, automatically commit the new file.

**Behavior**:
1. After encrypting and saving file:
   - Run `git add <path>.gpg`
   - Run `git commit -m "Add <path>"`
2. If git commands fail:
   - Warn user but don't fail the insert operation
   - Exit code for insert is still 0 (git failure is non-fatal)

**Skip Commit**:
- `--no-commit` flag on insert skips git operations

### 4.3 Git Configuration

**Environment Variables**:
- `PASS_GIT_NAME`: Git user.name override
- `PASS_GIT_EMAIL`: Git user.email override

---

## 5. Clipboard Integration

### 5.1 Windows Clipboard

**Requirement**: On Windows, use native clipboard API or `clip` command.

**Implementation Options**:
1. **Preferred**: Use PowerShell `Set-Clipboard` or `clip` command
   - `echo <password> | clip` (simple, works in cmd)
   - `Set-Clipboard -Value <password>` (PowerShell)
2. **Alternative**: Use .NET or Win32 API via PowerShell

**Behavior**:
- Copy decrypted content to clipboard
- Optionally clear clipboard after timeout (45 seconds default)
- Handle binary content (passwords are text, but be safe)

**Clear Clipboard**:
- After 45 seconds, overwrite clipboard with random data or empty string
- Implement as background process or timer
- Configurable via `PASS_CLIPBOARD_TIMEOUT` environment variable
- Set to 0 to disable auto-clear

---

## 6. Security Considerations

### 6.1 Passphrase Handling

**Requirement**: Never store passphrase in memory longer than necessary.

**Implementation**:
- Use secure input methods for passphrase (Windows: `Read-Host -AsSecureString`)
- Clear memory after use
- Rely on gpg-agent for passphrase caching

### 6.2 Temporary Files

**Requirement**: Securely delete temporary files containing plaintext passwords.

**Implementation**:
- Use cryptographically secure random names for temp files
- Overwrite file content before deletion
- Use `shred` equivalent or multiple overwrite passes
- On Windows: Use PowerShell `Clear-Content` + `Remove-Item` or Win32 API

### 6.3 Clipboard Security

**Requirement**: Minimize time sensitive data spends in clipboard.

**Implementation**:
- Auto-clear clipboard after 45 seconds (configurable)
- Warn user when copying to clipboard
- Consider requiring explicit opt-in for clipboard auto-clear

### 6.4 File Permissions

**Requirement**: Restrict access to password files.

**Implementation**:
- On Windows, set file permissions to owner-only
- Use `icacls` or PowerShell `Set-Acl`
- Inherit from parent directory or set explicitly

---

## 7. Error Handling

### 7.1 Error Messages

**Format**: `pass: <message>`

**Examples**:
- `pass: email/gmail.com/myemail: No such file or directory`
- `pass: password verification failed`
- `pass: GPG encryption failed: <gpg-error>`
- `pass: git commit failed: <git-error>`

### 7.2 Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error (file not found, invalid input) |
| 2 | GPG operation failed |
| 3 | File I/O error |
| 4 | Git operation failed (non-fatal for insert) |
| 5 | Clipboard operation failed |
| 6 | Password verification failed (mismatch or empty) |
| 7 | Remove operation failed |
| 8 | Fuzzy search error (no terminal, etc.) |

### 7.3 Verbose Mode

**Requirement**: `-v` or `--verbose` flag shows detailed error information.

**Behavior**:
- Show full error messages
- Show command being executed
- Show GPG/git command output

---

## 8. Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PASSWORD_STORE_DIR` | Password store directory | `%USERPROFILE%\.password-store` |
| `PASS_GPG_ID` | GPG recipient key ID | Default GPG key |
| `PASS_GPG_OPTS` | Additional GPG options | (empty) |
| `PASS_GIT_NAME` | Git user.name | System git config |
| `PASS_GIT_EMAIL` | Git user.email | System git config |
| `PASS_CLIPBOARD_TIMEOUT` | Clipboard clear timeout (seconds) | 45 |
| `PASS_CLIPBOARD_CLEAR` | Enable clipboard auto-clear | true |
| `PASS_NO_COLOR` | Disable ANSI color codes in fuzzy search | false |
| `PASS_FUZZY_TIMEOUT` | Timeout for fuzzy search (seconds, 0=none) | 0 |

---

## 9. Implementation Language

**Primary Choice**: Go (Golang)

**Rationale**:
- Cross-platform support (Windows, Linux, macOS)
- Excellent Windows 11/ARM64 support
- Single binary distribution (no runtime dependencies)
- Strong standard library for file operations, subprocess execution
- Built-in support for calling external commands (gpg, git)
- Good performance and memory management
- Well-suited for CLI tools
- Easy to compile and distribute
- Growing ecosystem for systems tools

**Alternatives Considered**:
- **PowerShell**: Native to Windows, but requires PowerShell runtime; less portable
- **Python**: Cross-platform, but requires Python installation; runtime dependency
- **Rust**: Similar benefits to Go, but more complex syntax and build system
- **C#/.NET**: Native to Windows, but requires .NET runtime; less portable
- **C/C++**: Native performance, but more complex development and memory management

**Decision**: Use Go (Golang) for implementation.

---

## 10. File Structure (Implementation)

```
main.go               # Main entry point
pass/                 # Main package
├── cmd/
│   ├── insert.go      # Insert command implementation
│   ├── show.go        # Show command implementation
│   ├── ls.go          # List command implementation
│   ├── find.go        # Find command implementation
│   ├── rm.go          # Remove command implementation
│   ├── fuzzy.go       # Fuzzy search command and UI
│   └── root.go        # Root command and global flags
├── pkg/
│   ├── gpg/
│   │   └── gpg.go      # GPG wrapper functions
│   ├── git/
│   │   └── git.go      # Git wrapper functions
│   ├── clipboard/
│   │   └── clipboard.go # Clipboard functions
│   ├── filesystem/
│   │   └── fs.go       # Filesystem utilities
│   ├── config/
│   │   └── config.go   # Configuration management
│   ├── fuzzy/
│   │   └── fuzzy.go    # Fuzzy matching algorithm
│   └── terminal/
│       └── terminal.go # Terminal UI utilities
├── internal/
│   └── utils/
│       └── helpers.go  # Internal utility functions
└── go.mod             # Go module file
tests/
├── insert_test.go     # Tests for insert command
├── show_test.go       # Tests for show command
├── ls_test.go         # Tests for ls command
└── find_test.go       # Tests for find command
docs/
  pass-replacement-spec.md  # This document
  pass-quick-reference.md  # Quick reference guide
  usage.md                 # User documentation
```

---

## 11. Dependencies

### 11.1 Required
- **GnuPG**: For encryption/decryption
  - Must be in PATH as `gpg`
  - Version: GPG4Win or GnuPG for Windows
  - Tested with: GnuPG 2.x

- **Git**: For version control
  - Must be in PATH as `git`
  - Version: Git for Windows

- **Go (Golang)**: For the pass binary itself
  - Version: Go 1.20+
  - Required for compilation

### 11.2 Optional
- **clip**: Windows clipboard utility (built-in)
- External dependencies are only needed at runtime for the compiled binary

---

## 12. Installation

### 12.1 Compilation
1. Ensure Go 1.20+ is installed
2. Clone the repository
3. Build the binary:
   ```bash
   go build -o pass.exe .
   ```
4. The compiled binary `pass.exe` is self-contained

### 12.2 Manual Installation
1. Place `pass.exe` in a directory in your PATH (e.g., `C:\Windows\System32\`)
2. Ensure GPG and Git are installed and in PATH
3. Initialize password store: `pass ls` (auto-creates directory)

### 12.3 Cross-Platform
For Linux/macOS:
```bash
go build -o pass .
chmod +x pass
sudo mv pass /usr/local/bin/
```

---

## 13. Testing Strategy

### 13.1 Unit Tests
- Test each function in isolation
- Mock external commands (gpg, git)
- Test path handling
- Test encryption/decryption

### 13.2 Integration Tests
- Test full command workflows
- Test with real GPG keys
- Test git integration
- Test clipboard functionality

### 13.3 End-to-End Tests
- Full user workflows
- Error scenarios
- Edge cases

### 13.4 Test Data
- Create test GPG keys for testing
- Use temporary directories for password store
- Clean up after tests

---

## 14. Compatibility

### 14.1 Unix pass Compatibility
- Command syntax should match where possible
- Exit codes should match
- Error messages should be similar

### 14.2 Differences from Unix pass
- Windows path handling
- Clipboard implementation (Unix uses xclip/xsel)
- GPG path/behavior differences
- Line ending handling

---

## 15. Future Enhancements

1. **pass edit**: Edit password in default editor
2. **pass generate**: Generate random password
3. **pass rm**: Remove password
4. **pass mv/cp**: Move/copy password
5. **pass git**: Pass-through git commands
6. **pass tree**: Tree view of password store
7. **pass otp**: One-time password support
8. **pass import/export**: Bulk operations
9. **GUI wrapper**: Optional GUI interface
10. **Browser integration**: Auto-fill support

---

## 16. Open Questions

1. **GPG Recipient**: Should we require explicit recipient or use default?
   - Decision: Use default key, allow override via PASS_GPG_ID

2. **Clipboard Auto-Clear**: Should this be on by default?
   - Decision: Yes, 45 seconds, configurable via PASS_CLIPBOARD_TIMEOUT

3. **Git Auto-Commit**: Should failures be fatal?
   - Decision: No, warn but don't fail insert

4. **Path Separators**: Should we support both `/` and `\`?
   - Decision: Yes, normalize internally

5. **Multi-line Passwords**: How to handle in clipboard?
   - Decision: Copy as-is, let clipboard handle it

6. **Empty Passwords**: Should we allow them?
   - Decision: No, reject empty passwords

---

## 17. References

- [Password Store (pass) Official Site](https://www.passwordstore.org/)
- [pass GitHub Repository](https://github.com/zx2c4/password-store)
- [GnuPG Documentation](https://www.gnupg.org/documentation/)
- [Git Documentation](https://git-scm.com/doc)
- [PowerShell Documentation](https://docs.microsoft.com/en-us/powershell/)

---

*Document Version: 1.0*
*Last Updated: 2026-06-05*
*Author: Spec-Driven Development*
