# Pass Replacement - Quick Reference

This is a concise reference for implementers and users of the pass replacement tool. For full details, see the specification documents.

---

## TL;DR

A Windows-compatible replacement for the Unix `pass` password manager that encrypts passwords with GPG and stores them in `~/.password-store/`.

---

## Installation

### Prerequisites
- Windows 11/ARM64 (or any Windows with Go 1.20+ support)
- [GPG4Win](https://www.gpg4win.org/) or GnuPG installed and in PATH
- [Git for Windows](https://gitforwindows.org/) installed and in PATH
- [Go (Golang)](https://golang.org/) 1.20+ installed (for compilation)
- A GPG key pair configured

### Quick Setup
1. Compile from source:
   ```bash
   go build -o pass.exe .
   ```
2. Place `pass.exe` in a directory in your PATH (e.g., `C:\Windows\System32\`)
3. Initialize the password store:
   ```bash
   pass ls
   ```

---

## Usage

### Core Commands

| Command | Description | Example |
|---------|-------------|---------|
| `pass insert <path>` | Insert new password | `pass insert email/gmail.com/user` |
| `pass <path>` | Show password | `pass email/gmail.com/user` |
| `pass -c <path>` | Copy to clipboard | `pass -c email/gmail.com/user` |
| `pass ls [<path>]` | List passwords | `pass ls` or `pass ls email/` |
| `pass find <string>` | Search passwords | `pass find gmail` |

### All Commands

```
# Insert
pass insert [OPTIONS] <path>
  Options:
    -e, --echo          Echo password while typing
    -m, --multiline    Allow multi-line password
    --no-commit        Skip git commit

# Show (default)
pass [OPTIONS] <path>
  Options:
    -c, --clip[board]  Copy to clipboard instead of stdout
    -o, --output <file> Write to file
    -q, --quiet        Suppress warnings
    --no-newline      Don't output trailing newline

# List
pass ls [OPTIONS] [<path>]
  Options:
    -r, --recursive    Show full paths (default)
    -d, --dirs-only    List only directories
    -f, --files-only   List only files

# Find
pass find [OPTIONS] <string>
  Options:
    -i, --ignore-case  Case-insensitive search

# Global Options
pass [COMMAND] [OPTIONS]
  Options:
    -v, --verbose      Verbose output
    --help, -h         Show help
    --version          Show version
```

---

## Examples

### Insert a Password
```bash
# Single-line password
$ pass insert email/gmail.com/myemail@gmail.com
Enter password for email/gmail.com/myemail@gmail.com: 
Retype password for email/gmail.com/myemail@gmail.com: 
Password inserted successfully.

# Multi-line password (e.g., private key)
$ pass insert --multiline certs/example.com/private-key
Enter password for certs/example.com/private-key (end with empty line):
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBg...
-----END PRIVATE KEY-----

Password inserted successfully.

# With echo (visible input)
$ pass insert --echo social/twitter.com/username
Enter password for social/twitter.com/username: mypassword
Retype password for social/twitter.com/username: mypassword
Password inserted successfully.

# Skip git commit
$ pass insert --no-commit temp/test
```

### Retrieve a Password
```bash
# Print to stdout
$ pass email/gmail.com/myemail@gmail.com
mypassword123

# Copy to clipboard
$ pass -c email/gmail.com/myemail@gmail.com
Copied email/gmail.com/myemail@gmail.com to clipboard.

# Write to file
$ pass -o password.txt email/gmail.com/myemail@gmail.com

# Quiet mode (no warnings)
$ pass -q email/gmail.com/myemail@gmail.com
```

### List Passwords
```bash
# List all
$ pass ls
email/gmail.com/myemail@gmail.com
banking/chase.com/account
social/twitter.com/username

# List under a path
$ pass ls email/
email/gmail.com/myemail@gmail.com

# List only directories
$ pass ls -d
email
banking
social

# List only files
$ pass ls -f
email/gmail.com/myemail@gmail.com
banking/chase.com/account
social/twitter.com/username
```

### Search Passwords
```bash
# Case-sensitive search
$ pass find gmail
email/gmail.com/myemail@gmail.com

# Case-insensitive search
$ pass find -i GMAIL
email/gmail.com/myemail@gmail.com
```

---

## File Structure

```
C:\Users\<username>\.password-store\  (or %USERPROFILE%\.password-store)
├── .git\                              (git repository)
│   ├── config
│   ├── HEAD
│   └── ...
├── email\                            (directory)
│   └── gmail.com\                    (directory)
│       └── myemail@gmail.com.gpg     (encrypted password file)
├── banking\                          (directory)
│   └── chase.com\                    (directory)
│       └── account.gpg               (encrypted password file)
└── social\                           (directory)
    └── twitter.com\                  (directory)
        └── username.gpg              (encrypted password file)
```

### File Contents

Each `.gpg` file contains a GPG-encrypted password. Example decrypted content:
```
mypassword123
```

Or for multi-line:
```
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7...
-----END PRIVATE KEY-----
```

---

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PASSWORD_STORE_DIR` | `%USERPROFILE%\.password-store` | Password store location |
| `PASS_GPG_ID` | Default GPG key | GPG recipient key ID |
| `PASS_CLIPBOARD_TIMEOUT` | 45 | Clipboard clear timeout (seconds) |
| `PASS_CLIPBOARD_CLEAR` | true | Enable clipboard auto-clear |
| `PASS_GIT_NAME` | Git config | Git user.name override |
| `PASS_GIT_EMAIL` | Git config | Git user.email override |

### Example Configuration
```bash
# Set in shell profile or system environment
# Windows (cmd):
set PASSWORD_STORE_DIR=D:\secure\passwords
set PASS_GPG_ID=ABCD1234
set PASS_CLIPBOARD_TIMEOUT=60

# Windows (PowerShell):
$env:PASSWORD_STORE_DIR = "D:\secure\passwords"
$env:PASS_GPG_ID = "ABCD1234"
$env:PASS_CLIPBOARD_TIMEOUT = 60

# Linux/macOS:
export PASSWORD_STORE_DIR="$HOME/secure/passwords"
export PASS_GPG_ID="ABCD1234"
export PASS_CLIPBOARD_TIMEOUT=60
```

---

## Git Integration

The password store is automatically initialized as a git repository on first use.

### Automatic Commits
- After `pass insert`, a git commit is automatically created
- Commit message: `Add <path>`
- Only local commit, **not pushed**
- Skip with `--no-commit` flag

### Manual Git Operations
```bash
# Navigate to password store
$ cd ~/.password-store

# Check status
$ git status

# Check history
$ git log --oneline

# Push to remote (manual)
$ git remote add origin <url>
$ git push -u origin main
```

---

## Security Features

### GPG Encryption
- All passwords encrypted with GPG
- Uses your GPG key pair
- Passphrase prompting handled by gpg-agent

### Secure Handling
- Plaintext passwords never stored on disk (except temporarily during insert)
- Temporary files securely deleted after use
- Clipboard auto-clears after 45 seconds (configurable)

### File Permissions
- Files created with restrictive permissions
- Only accessible to the owner

---

## Error Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error (file not found, invalid input) |
| 2 | GPG operation failed |
| 3 | File I/O error |
| 4 | Git operation failed |
| 5 | Clipboard operation failed |
| 6 | Password verification failed (mismatch or empty) |

---

## Troubleshooting

### Common Issues

**"pass: command not found"**
- Ensure `pass.exe` (or `pass` on Unix) is in your PATH
- Or use full path: `C:\path\to\pass.exe`

**"gpg: command not found"**
- Install GPG4Win or GnuPG
- Ensure `gpg.exe` is in your PATH
- Restart your shell after installation

**"git: command not found"**
- Install Git for Windows
- Ensure `git.exe` is in your PATH
- Git is optional for basic operations, but required for version control

**"No secret key" or "gpg: decryption failed"**
- Ensure you have a GPG key pair: `gpg --list-secret-keys`
- Ensure the password was encrypted to your key
- Check `PASS_GPG_ID` if you use multiple keys

**"No such file or directory"**
- Check the path is correct (case-sensitive)
- Use `pass ls` to see available passwords
- Ensure you're using forward slashes in paths

**"Permission denied"**
- Ensure you have permissions to the password store directory
- Try running as administrator (not recommended for security)

---

## Comparison with Unix pass

| Feature | Unix pass | Windows pass | Notes |
|---------|-----------|--------------|-------|
| Basic commands | ✅ | ✅ | insert, show, ls, find |
| GPG encryption | ✅ | ✅ | Same GPG backend |
| Git integration | ✅ | ✅ | Auto-commit on insert |
| Clipboard support | ✅ | ✅ | `-c` flag |
| Clipboard auto-clear | ❌ | ✅ | Windows-specific |
| Path separators | `/` | `/` or `\` | Normalized internally |
| Store location | `~/.password-store` | `%USERPROFILE%\.password-store` | Configurable |
| Multi-line passwords | ✅ | ✅ | `-m` flag |
| Tree view | ✅ | ❌ | Future enhancement |
| Edit command | ✅ | ❌ | Future enhancement |
| Generate command | ✅ | ❌ | Future enhancement |

---

## Implementation Notes

### For Developers

See the full specification documents:
- `docs/pass-replacement-spec.md` - Detailed technical specification
- `docs/pass-spec-kit.md` - Spec Kit aligned document
- `docs/pass-implementation-checklist.md` - Implementation tracking
- `docs/pass-decision-log.md` - Decision rationale

### File Structure (Implementation)
```
main.go               # Main entry point
go.mod                # Go module file
pass/
  cmd/
    insert.go          # Insert command
    show.go            # Show command
    ls.go              # List command
    find.go            # Find command
    root.go            # Root command
  pkg/
    gpg/
      gpg.go            # GPG operations
    git/
      git.go            # Git operations
    filesystem/
      fs.go             # Filesystem utilities
    clipboard/
      clipboard.go      # Clipboard operations
    config/
      config.go         # Configuration
tests/
  insert_test.go       # Tests for insert
  show_test.go         # Tests for show
  ls_test.go           # Tests for ls
  find_test.go         # Tests for find
docs/
  *.md                 # Documentation
```

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-06-05 | Initial quick reference |

---

*Last Updated: 2026-06-05*
