# Pass Replacement - Spec Kit Document

**Status**: Draft  
**Author**: @mandu  
**Created**: 2026-06-05  
**Last Updated**: 2026-06-05  

---

## 1. Purpose

Provide a Windows 11/ARM64-compatible replacement for the Unix `pass` password manager that offers equivalent basic functionality for managing GPG-encrypted password files.

### Why This Matters

- The native `pass` tool is not available on Windows 11/ARM64
- Users need a consistent password management experience across platforms
- GPG-based encryption provides strong security for stored credentials
- Git integration enables version control and backup of password store

---

## 2. Background

The Unix [password-store](https://www.passwordstore.org/) (`pass`) is a simple password manager that:
- Stores passwords as GPG-encrypted files in `~/.password-store/`
- Uses hierarchical directory structure for organization
- Provides CLI for insert, retrieve, list, and search operations
- Integrates with git for version control
- Supports clipboard integration

On Windows 11/ARM64, `pass` is unavailable, creating a gap for users who rely on this workflow.

### Current State

- No native Windows port of `pass`
- GPG is available via GPG4Win
- Git is available via Git for Windows
- Need a lightweight wrapper that provides `pass`-like interface

---

## 3. Goals

### In Scope

✅ **Core Functionality**
- `pass insert <path>` - Insert new password (prompt twice, encrypt, save)
- `pass <path>` - Retrieve and display password (decrypt, print to stdout)
- `pass -c <path>` - Retrieve and copy to clipboard
- `pass ls` - List all stored passwords
- `pass find <string>` - Search password paths for string

✅ **Data Management**
- Store passwords as GPG-encrypted files in `~/.password-store/`
- Use directory hierarchy matching password paths
- Automatic git commit on insert (local only, no push)

✅ **Security**
- GPG encryption/decryption
- Secure handling of plaintext passwords
- Clipboard auto-clear (configurable)

✅ **User Experience**
- Compatible command syntax with Unix `pass`
- Clear error messages
- GPG passphrase prompting when needed

### Out of Scope (v1.0)

❌ Advanced features for future versions:
- `pass edit` - Edit password in editor
- `pass generate` - Generate random password
- `pass rm` - Remove password
- `pass mv/cp` - Move/copy password
- `pass git` - Pass-through git commands
- `pass tree` - Tree view
- `pass otp` - One-time password support
- GUI interface
- Browser integration

---

## 4. User Stories

### As a user, I want to...

1. **Insert a new password** so I can store it securely
   ```
   pass insert email/gmail.com/myemail@gmail.com
   # Prompts: Enter password, Retype password
   # Creates: ~/.password-store/email/gmail.com/myemail@gmail.com.gpg
   # Commits to git
   ```

2. **Retrieve a password to stdout** so I can use it in scripts
   ```
   pass email/gmail.com/myemail@gmail.com
   # Prompts for GPG passphrase if needed
   # Prints: mysecretpassword123
   ```

3. **Copy a password to clipboard** so I can paste it easily
   ```
   pass -c email/gmail.com/myemail@gmail.com
   # Copies to clipboard, auto-clears after 45 seconds
   ```

4. **List all my passwords** so I can see what's stored
   ```
   pass ls
   # Output:
   # email/gmail.com/myemail@gmail.com
   # banking/chase.com/account
   # social/twitter.com/username
   ```

5. **Search for a password** so I can find it quickly
   ```
   pass find gmail
   # Output:
   # email/gmail.com/myemail@gmail.com
   # email/gmail.com/work@gmail.com
   ```

6. **Have version control** so my passwords are backed up
   ```
   # After insert, git commit is automatically created
   # cd ~/.password-store
   # git log
   # commit: Add email/gmail.com/myemail@gmail.com
   ```

---

## 5. Technical Requirements

### Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-001 | Implement `pass insert <path>` command | P0 |
| FR-002 | Prompt for password twice during insert | P0 |
| FR-003 | Encrypt password with GPG | P0 |
| FR-004 | Save as `<password-store>/<path>.gpg` | P0 |
| FR-005 | Implement `pass <path>` command | P0 |
| FR-006 | Decrypt and print password to stdout | P0 |
| FR-007 | Implement `pass -c <path>` clipboard option | P0 |
| FR-008 | Implement `pass ls` command | P0 |
| FR-009 | Implement `pass find <string>` command | P0 |
| FR-010 | Auto-commit to git on insert | P0 |
| FR-011 | Handle GPG passphrase prompting | P0 |
| FR-012 | Support `PASSWORD_STORE_DIR` environment variable | P1 |
| FR-013 | Support `PASS_GPG_ID` environment variable | P1 |
| FR-014 | Clipboard auto-clear after timeout | P1 |
| FR-015 | Support `--no-commit` flag for insert | P1 |

### Non-Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| NFR-001 | Must work on Windows 11/ARM64 | P0 |
| NFR-002 | Must require GPG installation | P0 |
| NFR-003 | Must require Git installation | P0 |
| NFR-004 | Plaintext passwords never stored on disk | P0 |
| NFR-005 | Temporary files securely deleted | P0 |
| NFR-006 | Clipboard cleared after 45 seconds | P0 |
| NFR-007 | Clear, actionable error messages | P0 |
| NFR-008 | Compatible with Unix pass command syntax | P1 |

### Assumptions

1. GPG (gpg.exe) is installed and in PATH
2. Git (git.exe) is installed and in PATH
3. User has a GPG key pair configured
4. gpg-agent is running and handles passphrase caching
5. User has appropriate permissions to create files and run git

### Dependencies

- **GPG4Win** or **GnuPG for Windows** (GPG 2.x)
- **Git for Windows**
- **Go (Golang) 1.20+** (for compilation)

---

## 6. Design

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        pass (Main)                             │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │ Argument     │  │ Command      │  │ Output       │       │
│  │ Parser       │─▶│ Router       │─▶│ Formatter    │       │
│  └──────────────┘  └──────────────┘  └──────────────┘       │
└─────────────────────────────────────────────────────────────┘
         │              │              │
         ▼              ▼              ▼
┌────────────────┐ ┌─────────────┐ ┌─────────────┐
│ Filesystem     │ │ GPG         │ │ Git         │
│ Module         │ │ Module      │ │ Module      │
│ - path utils   │ │ - encrypt   │ │ - init      │
│ - file ops     │ │ - decrypt   │ │ - add       │
│ - secure delete│ │ - key mgmt  │ │ - commit    │
└────────────────┘ └─────────────┘ └─────────────┘
         │              │
         └──────────────┴──────────────┘
                     │
              ┌──────▼──────┐
              │  Clipboard  │
              │   Module    │
              │ - copy      │
              │ - clear     │
              └─────────────┘
```

### Data Flow

#### Insert Flow
```
User Input: pass insert email/gmail.com/myemail
         │
         ▼
┌─────────────────────┐
│ Parse Arguments      │
│ - command: insert   │
│ - path: email/gmail… │
└─────────────┬───────┘
                │
                ▼
┌─────────────────────┐
│ Validate Path        │
│ - not empty          │
│ - valid characters   │
└─────────────┬───────┘
                │
                ▼
┌─────────────────────┐
│ Prompt for Password  │◄──────────────┐
│ - hidden input       │               │
│ - twice for verify   │               │
└─────────────┬───────┘               │
                │                     │
                ▼                     │
┌─────────────────────┐               │
│ Verify Passwords     │               │
│ - match?            │               │
│ - not empty?        │               │
└─────────────┬───────┘               │
                │ no                  │
                ▼                     │
        ┌─────────────────────┐        │
        │ Error: Mismatch      │        │
        │ Exit code 6          │        │
        └─────────────────────┘        │
                │ yes                 │
                ▼                     │
┌─────────────────────┐
│ Create Directory     │
│ Structure            │
└─────────────┬───────┘
                │
                ▼
┌─────────────────────┐
│ Write to Temp File   │
│ - secure location    │
│ - plaintext password │
└─────────────┬───────┘
                │
                ▼
┌─────────────────────┐
│ GPG Encrypt          │
│ gpg --encrypt ...    │
└─────────────┬───────┘
                │
                ▼
┌─────────────────────┐
│ Save as .gpg File    │
│ ~/.password-store/… │
└─────────────┬───────┘
                │
                ▼
┌─────────────────────┐
│ Secure Delete Temp   │
│ File                 │
└─────────────┬───────┘
                │
                ▼
┌─────────────────────┐
│ Git Add & Commit     │
│ git add ...          │
│ git commit -m "Add…" │
└─────────────┬───────┘
                │
                ▼
┌─────────────────────┐
│ Success Message      │
│ Exit code 0          │
└─────────────────────┘
```

#### Show Flow
```
User Input: pass email/gmail.com/myemail
         │
         ▼
┌─────────────────────┐
│ Parse Arguments      │
│ - command: show      │
│ - path: email/gmail… │
└─────────────┬───────┘
                │
                ▼
┌─────────────────────┐
│ Construct File Path  │
│ Add .gpg if needed   │
└─────────────┬───────┘
                │
                ▼
┌─────────────────────┐
│ Check File Exists    │
│ ~/.password-store/…  │
└─────────────┬───────┘
                │ no
                ▼
        ┌─────────────────────┐
        │ Error: Not found     │
        │ Exit code 1          │
        └─────────────────────┘
                │ yes
                ▼
┌─────────────────────┐
│ GPG Decrypt          │
│ gpg --decrypt …      │
│ (prompts if needed)  │
└─────────────┬───────┘
                │
                ▼
┌─────────────────────┐
│ Check for -c flag    │
│ Copy to clipboard?    │
└─────────────┬───────┘
                │ yes
                ▼
┌─────────────────────┐
│ Copy to Clipboard    │
│ Start clear timer    │
└─────────────────────┘
                │ no
                ▼
┌─────────────────────┐
│ Print to Stdout      │
└─────────────────────┘
                │
                ▼
┌─────────────────────┐
│ Exit code 0          │
└─────────────────────┘
```

### File Structure

```
~/.password-store/
├── .git/
│   ├── config
│   ├── HEAD
│   └── ...
├── email/
│   └── gmail.com/
│       └── myemail@gmail.com.gpg
├── banking/
│   └── chase.com/
│       └── account.gpg
└── social/
    └── twitter.com/
        └── username.gpg
```

Each `.gpg` file contains a single GPG-encrypted password (or multi-line content).

---

## 7. Command Reference

### Syntax
```
pass [OPTIONS] [COMMAND] [ARGS...]
pass [OPTIONS] <path>  # Default: show command
```

### Commands

| Command | Description | Example |
|---------|-------------|---------|
| `insert <path>` | Insert new password | `pass insert email/gmail.com/user` |
| `show <path>` | Show password | `pass email/gmail.com/user` |
| `ls [<path>]` | List passwords | `pass ls` or `pass ls email/` |
| `find <string>` | Search passwords | `pass find gmail` |

### Options

| Option | Description | Applies To |
|--------|-------------|------------|
| `-c, --clip[board]` | Copy to clipboard | show |
| `-e, --echo` | Echo password while typing | insert |
| `-m, --multiline` | Allow multi-line password | insert |
| `--no-commit` | Skip git commit | insert |
| `-o, --output <file>` | Write to file | show |
| `-q, --quiet` | Suppress warnings | show |
| `--no-newline` | No trailing newline | show |
| `-r, --recursive` | Full paths (default) | ls |
| `-d, --dirs-only` | List only directories | ls |
| `-f, --files-only` | List only files | ls |
| `-i, --ignore-case` | Case-insensitive search | find |
| `-v, --verbose` | Verbose output | all |
| `--help, -h` | Show help | all |
| `--version` | Show version | all |

---

## 8. Error Handling

### Error Messages

All errors follow the format: `pass: <message>`

| Error | Message | Exit Code |
|-------|---------|-----------|
| File not found | `pass: <path>: No such file or directory` | 1 |
| Password mismatch | `pass: password verification failed` | 6 |
| Empty password | `pass: password cannot be empty` | 6 |
| GPG encryption fail | `pass: GPG encryption failed: <error>` | 2 |
| GPG decryption fail | `pass: GPG decryption failed: <error>` | 2 |
| Git commit fail | `pass: git commit failed: <error>` | 4 |
| Clipboard fail | `pass: failed to copy to clipboard` | 5 |
| Invalid path | `pass: <path>: Invalid path` | 1 |

---

## 9. Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PASSWORD_STORE_DIR` | `%USERPROFILE%\.password-store` | Password store location |
| `PASS_GPG_ID` | Default GPG key | GPG recipient key ID |
| `PASS_GPG_OPTS` | (empty) | Additional GPG options |
| `PASS_GIT_NAME` | Git config | Git user.name override |
| `PASS_GIT_EMAIL` | Git config | Git user.email override |
| `PASS_CLIPBOARD_TIMEOUT` | 45 | Clipboard clear timeout (seconds) |
| `PASS_CLIPBOARD_CLEAR` | true | Enable clipboard auto-clear |

### Configuration File (Future)

Location: `%USERPROFILE%\.password-store\.passrc`

```ini
[core]
store_dir = C:\Users\user\.password-store
gpg_id = ABCD1234
clipboard_timeout = 45

[git]
name = User Name
email = user@example.com
```

---

## 10. Open Questions

1. **Should we support GPG armor (ASCII) mode?**
   - Unix `pass` uses `--armor` by default
   - Binary GPG files are more compact
   - Decision: Use `--armor` for compatibility

2. **How to handle GPG key selection?**
   - Option A: Use default key (simplest)
   - Option B: Require `PASS_GPG_ID` to be set
   - Option C: Prompt user to select key on first use
   - Decision: Use default key, allow override via `PASS_GPG_ID`

3. **Should clipboard auto-clear be on by default?**
   - Security vs. usability tradeoff
   - Decision: Yes, with configurable timeout

4. **How to handle line endings in passwords?**
   - Unix: LF
   - Windows: CRLF
   - Decision: Preserve input as-is, no conversion

5. **Should we support symlinks in password store?**
   - Unix `pass` supports symlinks
   - Windows symlink support varies
   - Decision: Support if possible, but not a priority

6. **What to do if GPG is not installed?**
   - Option A: Error with clear message
   - Option B: Offer to install GPG4Win
   - Decision: Error with message including download link

7. **What to do if Git is not installed?**
   - Similar to GPG handling
   - Decision: Error with message, but allow basic operations without git

---

## 11. Implementation Phases

### Phase 1: Core MVP (P0)
- [ ] Basic infrastructure (argument parsing, config)
- [ ] GPG encryption/decryption
- [ ] Filesystem operations
- [ ] `insert` command
- [ ] `show` command
- [ ] Basic error handling

### Phase 2: Essential Features (P0)
- [ ] `ls` command
- [ ] `find` command
- [ ] Git integration
- [ ] Clipboard support

### Phase 3: Polish (P1)
- [ ] Environment variable support
- [ ] Clipboard auto-clear
- [ ] Verbose mode
- [ ] Help system
- [ ] All flags and options

### Phase 4: Testing & Documentation
- [ ] Unit tests
- [ ] Integration tests
- [ ] User documentation
- [ ] Installation guide

---

## 12. Success Criteria

### MVP (v1.0) is Complete When:

✅ `pass insert <path>` works and encrypts password  
✅ `pass <path>` works and decrypts/displays password  
✅ `pass -c <path>` works and copies to clipboard  
✅ `pass ls` lists all passwords  
✅ `pass find <string>` searches passwords  
✅ Git commit happens automatically on insert  
✅ All P0 functional requirements are implemented  
✅ All P0 non-functional requirements are met  
✅ Basic error handling is in place  
✅ User can install and use the tool  

---

## 13. References

- [Password Store Official Site](https://www.passwordstore.org/)
- [pass GitHub Repository](https://github.com/zx2c4/password-store)
- [GnuPG Documentation](https://www.gnupg.org/documentation/)
- [Git Documentation](https://git-scm.com/doc)
- [Go Documentation](https://golang.org/doc/)
- [GitHub Spec Kit](https://github.com/github/spec-kit)

---

## 14. Appendix

### A. Example Session

```bash
# Initialize (first time)
$ pass ls
pass: warning: Password store not initialized. Creating...

# Insert a password
$ pass insert email/gmail.com/myemail@gmail.com
Enter password for email/gmail.com/myemail@gmail.com: 
Retype password for email/gmail.com/myemail@gmail.com: 
Password inserted successfully.

# Retrieve a password
$ pass email/gmail.com/myemail@gmail.com
mysecretpassword123

# Copy to clipboard
$ pass -c email/gmail.com/myemail@gmail.com
Copied email/gmail.com/myemail@gmail.com to clipboard.

# List all passwords
$ pass ls
email/gmail.com/myemail@gmail.com

# Insert another password
$ pass insert banking/chase.com/account
Enter password for banking/chase.com/account: 
Retype password for banking/chase.com/account: 
Password inserted successfully.

# List all passwords
$ pass ls
banking/chase.com/account
email/gmail.com/myemail@gmail.com

# Search
PS> pass find gmail
email/gmail.com/myemail@gmail.com

# Check git history
PS> cd ~/.password-store
PS> git log --oneline
abc1234 Add email/gmail.com/myemail@gmail.com
def5678 Add banking/chase.com/account
```

### B. File Contents

Before encryption (temp file):
```
mysecretpassword123
```

After encryption (`~/.password-store/email/gmail.com/myemail@gmail.com.gpg`):
```
-----BEGIN PGP MESSAGE-----

jA0ECwMC... (GPG armored content)
...
-----END PGP MESSAGE-----
```

### C. Git Repository

```
~/.password-store/.git/config:
[core]
	repositoryformatversion = 0
	filemode = false
	bare = false
	logallrefupdates = true

~/.password-store/.git/logs/HEAD:
... commit logs ...
```

---

*Spec Kit Version: 1.0*  
*Document Status: Draft*  
*Next Review: TBD*
