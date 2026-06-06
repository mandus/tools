# pass - Windows-compatible password store

A Windows-compatible replacement for the Unix password-store tool. Manages GPG-encrypted password files in `~/.password-store/` with git integration.

## Features

- ✅ GPG-encrypted password storage
- ✅ Git integration for version control
- ✅ Interactive fuzzy search with TUI (Bubble Tea)
- ✅ Cross-platform (Windows, Linux, macOS)
- ✅ Copy to clipboard support

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

# List passwords
pass ls

# Find passwords
pass find gmail

# Insert new password
pass insert email/new.com
```

## TUI

When running `pass`, `pass -c`, or `pass rm` without arguments, an interactive TUI is displayed:
- Use ↑/↓ arrows or j/k to navigate
- Type to filter passwords
- Press Enter to select
- Press Esc, Ctrl+C, Ctrl+D, or Ctrl+Q to exit

## Installation

```bash
go build -o pass.exe
```

## Configuration

- `PASSWORD_STORE_DIR`: Custom password store location (default: `~/.password-store`)
- Requires GPG to be installed and configured

## Documentation

See `docs/` for detailed documentation and `specs/` for specifications.
