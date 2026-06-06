# Edit Command

The `pass edit` command allows you to edit existing password files using your favorite text editor.

## Usage

```bash
# Edit a specific password
pass edit email/gmail.com/user

# Edit with fuzzy search
pass edit

# Edit without git commit
pass edit --no-commit email/gmail.com/user
pass edit -n email/gmail.com/user

# Edit with force flag (alias for --no-commit)
pass edit -f email/gmail.com/user
```

## How It Works

1. **Decrypt**: The password file is decrypted using GPG
2. **Open Editor**: The plaintext content is loaded into a temporary file and opened in your editor
3. **Wait for Save**: The program waits for you to save and close the editor
4. **Re-encrypt**: The modified content is read from the temporary file and re-encrypted
5. **Cleanup**: The temporary file is securely deleted (overwritten before removal)
6. **Git Commit**: The changes are committed to git (unless `--no-commit` flag is used)

## Editor Selection

The editor is determined in the following order:

1. **`$EDITOR` environment variable** - If set, this editor is used
2. **Platform default**:
   - Windows: `notepad`
   - Linux/macOS: `vi`

### Setting a Custom Editor

You can set your preferred editor by setting the `EDITOR` environment variable:

```bash
# Windows (Command Prompt)
set EDITOR=notepad++

# Windows (PowerShell)
$env:EDITOR = "notepad++"

# Linux/macOS (Bash/Zsh)
export EDITOR=nano
```

### Supported Editors

The edit command works with any editor that can be invoked as a command-line program:

- **GUI Editors**: Notepad, Notepad++, VS Code, Sublime Text, etc.
- **Terminal Editors**: vi, vim, emacs, nano, etc.
- **Custom**: Any executable in your PATH

## Temporary Files

- Created in the system temp directory with prefix `pass-edit-*.tmp`
- Set with restrictive permissions (0600 - readable/writable only by owner)
- Securely deleted after use (content is overwritten before file deletion)

## Git Integration

By default, changes are committed to git with a message like:
```
Edit email/gmail.com/user
```

To skip the git commit, use the `--no-commit` or `-n` flag:
```bash
pass edit --no-commit email/gmail.com/user
```

Or the `-f` (force) flag:
```bash
pass edit -f email/gmail.com/user
```

## Error Handling

| Error | Message | Exit Code |
|-------|---------|-----------|
| File not found | `pass: <path>: No such file or directory` | 1 |
| Is a directory | `pass: <path>: Is a directory` | 1 |
| Decryption failed | `pass: decryption failed: <error>` | 3 |
| Editor not found | `pass: failed to open editor: <error>` | 2 |
| Editor exit error | `pass: editor exited with status <code>` | 2 |
| Encryption failed | `pass: GPG encryption failed: <error>` | 4 |
| Git error | Warning to stderr (non-fatal) | 5 |

## Examples

### Edit a password
```bash
$ pass edit email/gmail.com/user
# Opens in editor, user makes changes and saves
Password updated successfully.
```

### Edit with fuzzy search
```bash
$ pass edit
# TUI appears, user types "gmail", selects "email/gmail.com/user", presses Enter
# Opens in editor, user makes changes and saves
Password updated successfully.
```

### Edit without git commit
```bash
$ pass edit --no-commit social/twitter.com/token
# Opens in editor, user makes changes and saves
Password updated successfully.
# No git commit is made
```

### Edit with custom editor
```bash
$ EDITOR=nano pass edit banking/chase.com/pin
# Opens in nano, user makes changes and saves
Password updated successfully.
```

## Security Considerations

- Temporary files are created with restrictive permissions (0600)
- Temporary files are securely deleted (overwritten) after use
- Content is encrypted at rest
- Git history is preserved (unless `--no-commit` is used)

## Related Commands

- `pass show` - View a password
- `pass insert` - Add a new password
- `pass rm` - Remove a password
- `pass ls` - List passwords
