# Pass Replacement Implementation Checklist

This checklist tracks implementation progress for the pass replacement tool. Each item corresponds to a requirement in the main specification document.

## Phase 1: Core Infrastructure (MVP)

- [ ] **Project Setup**
  - [ ] Initialize Go module (`go mod init`)
  - [ ] Set up project structure
  - [ ] Configure build tags for Windows
  - [ ] Set up Go toolchain
  
- [ ] **Configuration**
  - [ ] Define `PASSWORD_STORE_DIR` with default `%USERPROFILE%\.password-store`
  - [ ] Support environment variable overrides
  - [ ] Create directory structure if it doesn't exist
  - [ ] Implement config package
  
- [ ] **GPG Integration**
  - [ ] Implement GPG encryption function (exec.Command)
  - [ ] Implement GPG decryption function (exec.Command)
  - [ ] Handle GPG passphrase prompting
  - [ ] Support `PASS_GPG_ID` environment variable
  - [ ] Error handling for GPG operations
  - [ ] Create gpg package
  
- [ ] **Filesystem Operations**
  - [ ] Implement path normalization (convert `/` to `\`)
  - [ ] Implement secure temporary file creation
  - [ ] Implement secure file deletion (overwrite before delete)
  - [ ] Handle invalid characters in paths
  - [ ] Create parent directories as needed
  - [ ] Create filesystem package
  
- [ ] **Main Entry Point**
  - [ ] Create `main.go` entry point
  - [ ] Parse command line arguments (using cobra or flag)
  - [ ] Route to appropriate command handlers
  - [ ] Global error handling
  - [ ] Verbose mode support

## Phase 2: Core Commands (MVP)

- [ ] **insert command**
  - [ ] Parse path argument
  - [ ] Validate path
  - [ ] Prompt for password (hidden input using golang.org/x/term)
  - [ ] Prompt for password confirmation
  - [ ] Verify passwords match and are not empty
  - [ ] Create directory structure
  - [ ] Write password to temporary file
  - [ ] Encrypt with GPG
  - [ ] Save as `<path>.gpg`
  - [ ] Securely delete temporary file
  - [ ] Display success message
  - [ ] Support `--no-commit` flag
  - [ ] Support `-e/--echo` flag
  - [ ] Support `-m/--multiline` flag

- [ ] **show command (default)**
  - [ ] Parse path argument
  - [ ] Construct file path (add `.gpg` if needed)
  - [ ] Check file exists
  - [ ] Decrypt with GPG
  - [ ] Output to stdout (default)
  - [ ] Support `-c/--clipboard` flag
  - [ ] Support `-o/--output` flag
  - [ ] Support `-q/--quiet` flag
  - [ ] Support `--no-newline` flag
  - [ ] Error handling for missing files
  - [ ] Error handling for decryption failures

- [ ] **ls command**
  - [ ] Parse optional subpath argument
  - [ ] Validate subpath exists
  - [ ] Recursively walk directory tree
  - [ ] Find all `.gpg` files
  - [ ] Strip `.gpg` extension from results
  - [ ] Strip password store prefix
  - [ ] Sort results alphabetically
  - [ ] Support `-r/--recursive` flag (default)
  - [ ] Support `-d/--dirs-only` flag
  - [ ] Support `-f/--files-only` flag
  - [ ] Display one result per line

- [ ] **find command**
  - [ ] Parse search string argument
  - [ ] Validate search string is not empty
  - [ ] Walk entire password store
  - [ ] Match search string against paths
  - [ ] Support `-i/--ignore-case` flag
  - [ ] Sort results alphabetically
  - [ ] Display matching paths

## Phase 3: Git Integration

- [ ] **Repository Initialization**
  - [ ] Check if `.password-store/.git` exists
  - [ ] Run `git init` if not exists
  - [ ] Configure git user.name and user.email if not set
  - [ ] Support `PASS_GIT_NAME` and `PASS_GIT_EMAIL` environment variables
  
- [ ] **Automatic Commits**
  - [ ] After successful insert, run `git add <path>.gpg`
  - [ ] Run `git commit -m "Add <path>"`
  - [ ] Handle git errors gracefully (non-fatal)
  - [ ] Support `--no-commit` flag to skip

## Phase 4: Clipboard Integration

- [ ] **Windows Clipboard**
  - [ ] Implement clipboard copy using `clip` command via exec.Command
  - [ ] Handle multi-line content
  - [ ] Error handling for clipboard operations
  - [ ] Create clipboard package

## Phase 5: Remove Command & Fuzzy Search

- [ ] **Fuzzy Matching Package (pkg/fuzzy/)**
  - [ ] Implement `Match(query, target string) bool` - subsequence check
  - [ ] Implement `Score(query, target string) int` - ranking algorithm
  - [ ] Implement `Filter(query string, items []string) []MatchResult`
  - [ ] Create fuzzy package with tests

- [ ] **Terminal UI Package (pkg/terminal/)**
  - [ ] Implement ANSI escape code utilities
  - [ ] Implement cursor control functions
  - [ ] Implement terminal size detection
  - [ ] Implement key reading with special key support
  - [ ] Create terminal package with tests

- [ ] **Fuzzy Search Command (cmd/fuzzy.go)**
  - [ ] Implement main fuzzy search loop
  - [ ] Implement display rendering
  - [ ] Implement query input handling
  - [ ] Implement list navigation
  - [ ] Implement match highlighting
  - [ ] Handle all keybindings (Ctrl+A, Ctrl+E, Ctrl+K, arrows, etc.)
  - [ ] Support different modes (show, clip, rm)
  - [ ] Create fuzzy command tests

- [ ] **Remove Command (cmd/rm.go)**
  - [ ] Implement rm command with cobra
  - [ ] Add flags: --no-commit/-n, --force/-f, --clip/-c
  - [ ] Implement removePassword() function
  - [ ] Handle explicit path removal
  - [ ] Handle fuzzy search mode for rm
  - [ ] Git integration (git rm + commit)
  - [ ] Create rm command tests

- [ ] **Git Integration Enhancements**
  - [ ] Add RemoveAndCommit() function to pkg/git/
  - [ ] Update git tests
  
- [ ] **Clipboard Auto-Clear**
  - [ ] Implement timer for auto-clear (using time.Timer)
  - [ ] Default timeout: 45 seconds
  - [ ] Support `PASS_CLIPBOARD_TIMEOUT` environment variable
  - [ ] Support `PASS_CLIPBOARD_CLEAR` to enable/disable
  - [ ] Overwrite clipboard with random data on clear

## Phase 5: Security Features

- [ ] **Secure Input**
  - [ ] Use `Read-Host -AsSecureString` for password input
  - [ ] Clear secure string from memory after use
  
- [ ] **Secure Temporary Files**
  - [ ] Generate cryptographically random temp file names
  - [ ] Use secure deletion (multiple overwrite passes)
  - [ ] Ensure temp files are on same volume (for secure delete)
  
- [ ] **File Permissions**
  - [ ] Set restrictive permissions on password files
  - [ ] Use `icacls` or PowerShell ACL cmdlets
  - [ ] Inherit from parent or set explicitly

## Phase 6: Error Handling & User Experience

- [ ] **Error Messages**
  - [ ] Format: `pass: <message>`
  - [ ] Clear, actionable messages
  - [ ] Context-specific errors
  
- [ ] **Exit Codes**
  - [ ] Implement all exit codes from spec
  - [ ] Return appropriate codes for all error conditions
  
- [ ] **Verbose Mode**
  - [ ] Implement `-v/--verbose` flag
  - [ ] Show debug information
  - [ ] Show commands being executed
  
- [ ] **Help System**
  - [ ] Implement `--help/-h` flag
  - [ ] Implement `--version` flag
  - [ ] Show usage for each command

## Phase 7: Testing

- [ ] **Test Infrastructure**
  - [ ] Set up test directory structure
  - [ ] Create test GPG keys
  - [ ] Mock external commands for unit tests
  - [ ] Test helper functions
  
- [ ] **Unit Tests**
  - [ ] Test GPG encryption/decryption
  - [ ] Test path handling
  - [ ] Test filesystem operations
  - [ ] Test argument parsing
  
- [ ] **Integration Tests**
  - [ ] Test full insert workflow
  - [ ] Test full show workflow
  - [ ] Test ls command
  - [ ] Test find command
  - [ ] Test git integration
  - [ ] Test clipboard functionality
  
- [ ] **End-to-End Tests**
  - [ ] Test complete user scenarios
  - [ ] Test error conditions
  - [ ] Test edge cases

## Phase 8: Documentation

- [ ] **User Documentation**
  - [ ] Create usage guide
  - [ ] Installation instructions
  - [ ] Configuration options
  - [ ] Command reference
  - [ ] Examples
  
- [ ] **Developer Documentation**
  - [ ] Code comments
  - [ ] Function documentation
  - [ ] Architecture overview
  - [ ] Contributing guide

## Phase 9: Packaging & Distribution

- [ ] **Installation Script**
  - [ ] Create installer or setup script
  - [ ] Verify dependencies (GPG, Git)
  - [ ] Create batch wrapper (`pass.bat`)
  
- [ ] **PowerShell Module**
  - [ ] Create module manifest (`pass.psd1`)
  - [ ] Package as installable module
  - [ ] Support PowerShell Gallery
  
- [ ] **Release Checklist**
  - [ ] Version bumping
  - [ ] Changelog
  - [ ] Test on clean system
  - [ ] Verify all features work

## Phase 10: Future Enhancements (Post-MVP)

- [ ] **pass edit** - Edit password in editor
- [ ] **pass generate** - Generate random password
- [x] **pass rm** - Remove password
- [ ] **pass mv/cp** - Move/copy password
- [ ] **pass git** - Pass-through git commands
- [ ] **pass tree** - Tree view
- [ ] **pass otp** - One-time password support
- [ ] **pass import/export** - Bulk operations
- [x] **pass fuzzy search** - Interactive fuzzy finder
- [x] **pass rm with fuzzy search** - Remove with fuzzy selection

---

## Implementation Order Recommendation

For rapid MVP delivery, implement in this order:

1. **Core Infrastructure** (Configuration, GPG, Filesystem)
2. **show command** (easiest to test)
3. **insert command** (core functionality)
4. **ls command** (useful for listing)
5. **find command** (search capability)
6. **Git Integration** (version control)
7. **Clipboard Integration** (quality of life)
8. **Security Features** (hardening)
9. **Error Handling & UX** (polish)
10. **Testing** (ongoing)
11. **Documentation** (ongoing)

---

## Estimated Complexity

| Component | Complexity | Effort |
|-----------|------------|--------|
| Core Infrastructure | Medium | 2-4 hours |
| insert command | Medium | 3-5 hours |
| show command | Medium | 2-3 hours |
| ls command | Low | 1-2 hours |
| find command | Low | 1-2 hours |
| Git Integration | Low | 1-2 hours |
| Clipboard Integration | Medium | 2-3 hours |
| Security Features | Medium | 2-3 hours |
| Error Handling | Medium | 2-3 hours |
| Testing | Medium | 3-5 hours |
| Documentation | Low | 2-3 hours |
| **Total** | | **20-35 hours** |

---

*Checklist Version: 1.0*
*Last Updated: 2026-06-05*
