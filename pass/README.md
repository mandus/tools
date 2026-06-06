# pass - Windows-compatible password store

A Windows-compatible replacement for the Unix password-store tool. Manages GPG-encrypted password files in `~/.password-store/` with git integration.

## Features

- ✅ GPG-encrypted password storage
- ✅ Git integration for version control
- ✅ Interactive fuzzy search with TUI (Bubble Tea)
- ✅ Cross-platform (Windows, Linux, macOS)
- ✅ Copy to clipboard support
- ✅ Edit existing passwords with your favorite editor

## Usage

```bash
# Show password
pass email/gmail.com

# Interactive fuzzy search
pass

# Copy to clipboard
pass -c email/gmail.com

# Remove password
pass rm email/old.com

# List passwords (files only, not directories)
pass ls

# Find passwords
pass find gmail

# Insert new password
pass insert email/new.com

# Edit existing password (opens in $EDITOR)
pass edit email/gmail.com

# Edit with fuzzy search
pass edit
```

## TUI

When running `pass`, `pass -c`, `pass rm`, or `pass edit` without arguments, an interactive TUI is displayed:
- Use ↑/↓ arrows to navigate the list
- Use ←/→ arrows to move cursor in search input
- Use Tab to cycle through results
- Type to filter passwords using **fuzzy matching** (subsequence matching)
- Press Enter to select (show, copy, remove, or edit based on command)
- Press Esc, Ctrl+C, Ctrl+D, or Ctrl+Q to exit

**Fuzzy Matching Examples:**
- Type `"g"` → shows all items containing "g"
- Type `"ga"` → shows items with "g" followed by "a" (in any position after)
- Type `"gmail"` → shows items like `email/gmail.com/user` where characters appear in order
- Matching characters are highlighted in the results

## Installation

```bash
go build -o pass.exe
```

## Configuration

- `PASSWORD_STORE_DIR`: Custom password store location (default: `~/.password-store`)
- `EDITOR`: Editor to use for editing passwords (default: `notepad` on Windows, `vi` on Unix)
- Requires GPG to be installed and configured

## Documentation

See `docs/` for detailed documentation and `specs/` for specifications.
