# Changelog

All notable changes to the pass tool are documented in this file.

## [Unreleased]

### Added
- **Edit command**: `pass edit <path>` allows editing existing passwords in your favorite editor
- **Fuzzy search for edit**: `pass edit` without arguments opens fuzzy search to select a password to edit
- **Edit mode in TUI**: Added support for edit mode in the Bubble Tea TUI

### Fixed
- **Insert overwrite prevention**: `pass insert <path>` now fails with error `pass: <path>: Already exists` if the file already exists
- **List command**: `pass ls` now only lists secret files (`.gpg` files), not directories

### Changed
- **List command behavior**: Directories are no longer shown in `pass ls` output unless `--dirs-only` flag is used
- **TUI modes**: Added `FuzzyModeEdit` to support editing via fuzzy search

### Documentation
- Updated README with edit command usage
- Added `edit-command.md` with detailed documentation
- Updated `tui.md` to include edit mode

## Implementation Details

### Edit Command Workflow
1. Decrypt the password file
2. Create a secure temporary file (0600 permissions)
3. Open the temporary file in the user's editor (from `$EDITOR` or platform default)
4. Wait for the editor to exit
5. Read the modified content
6. Re-encrypt the content back to the original file
7. Securely delete the temporary file
8. Commit changes to git (unless `--no-commit` flag is used)

### Editor Selection
- Priority: `$EDITOR` environment variable → Platform default (`notepad` on Windows, `vi` on Unix)

### Security
- Temporary files are created with restrictive permissions
- Temporary files are securely overwritten before deletion
- All operations respect existing file permissions

### Files Modified
- `cmd/edit.go` - New file with edit command implementation
- `cmd/edit_test.go` - New file with edit command tests
- `cmd/insert.go` - Added file existence check before insert
- `cmd/insert_test.go` - Added tests for overwrite prevention
- `cmd/ls.go` - Modified to only list files, not directories
- `cmd/ls_test.go` - Updated tests for file-only listing
- `cmd/root.go` - Registered edit command
- `cmd/fuzzy.go` - Added `FuzzyModeEdit` constant and support
- `cmd/tui/fuzzy.go` - Added `FuzzyModeEdit` constant
- `cmd/tui/models.go` - Updated getTitle and getPrompt for edit mode
- `README.md` - Updated with edit command documentation
- `docs/tui.md` - Updated with edit mode
- `docs/edit-command.md` - New file with detailed edit command docs

### Files Added
- `cmd/edit.go`
- `cmd/edit_test.go`
- `cmd/insert_test.go`
- `cmd/ls_test.go`
- `docs/edit-command.md`
- `specs/pass-edit/spec.md`
- `specs/pass-edit/tasks.md`

## Testing

All existing tests continue to pass. New tests added for:
- Edit command functionality
- Insert overwrite prevention
- List command file-only output
- Editor detection
- Error handling

## Backward Compatibility

All changes are backward compatible. Existing functionality remains unchanged:
- `pass insert` now prevents overwrites (this is a bug fix, not a breaking change)
- `pass ls` now only shows files (this is a bug fix, not a breaking change)
- `pass edit` is a new command
