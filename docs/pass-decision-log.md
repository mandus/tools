# Pass Replacement - Decision Log

This document tracks design decisions made during the specification and implementation of the pass replacement tool. Each decision includes the context, options considered, the chosen approach, and the rationale.

---

## Format

```
## [Decision-ID] Decision Title

**Status**: Accepted | Rejected | Superseded  
**Date**: YYYY-MM-DD  
**Context**: Description of the problem or question  
**Decision**: The chosen solution  
**Rationale**: Why this decision was made  
**Alternatives Considered**: Other options that were rejected  
**Consequences**: Implications of this decision  
**Related**: Links to specs, issues, or other decisions
```

---

## Architecture Decisions

### [AD-001] Implementation Language: Go (Golang)

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: Need to choose a language for implementing the pass replacement tool on Windows.

**Decision**: Use Go (Golang) as the primary implementation language.

**Rationale**:
- Cross-platform support (Windows, Linux, macOS)
- Excellent Windows 11/ARM64 support
- Single binary distribution (no runtime dependencies)
- Strong standard library for file operations and subprocess execution
- Built-in support for calling external commands (gpg, git, clip)
- Good performance and memory management
- Well-suited for CLI tools
- Easy to compile and distribute
- Growing ecosystem for systems tools
- User preference for Go over PowerShell

**Alternatives Considered**:
1. **PowerShell**: Native to Windows, but requires PowerShell runtime; less portable; runtime dependency
2. **Python**: Cross-platform, but requires Python installation; runtime dependency
3. **Rust**: Similar benefits to Go, but more complex syntax and build system
4. **C#/.NET**: Native to Windows, but requires .NET runtime; less portable
5. **C/C++**: Native performance, but more complex development and memory management

**Consequences**:
- Tool will be a compiled Go binary (`pass.exe` on Windows, `pass` on Unix)
- No runtime dependencies (self-contained binary)
- Can be distributed as a single file
- Cross-platform compatible
- Requires Go toolchain for compilation (but not for running)
- Better alignment with user's preferences

**Related**: pass-replacement-spec.md (Section 9)

---

### [AD-002] Password Store Location

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: Where should the password store be located on Windows?

**Decision**: Use `%USERPROFILE%\.password-store` as the default location, with `PASSWORD_STORE_DIR` environment variable override.

**Rationale**:
- Matches Unix convention (`~/.password-store`)
- `%USERPROFILE%` is the Windows equivalent of `~`
- Environment variable override provides flexibility
- Consistent with how Unix pass works

**Alternatives Considered**:
1. **`%APPDATA%\password-store`**: More "Windows-like" but differs from Unix pass
2. **`%LOCALAPPDATA%\password-store`**: Similar to above
3. **Configurable only via config file**: Less flexible than environment variable
4. **Prompt user on first run**: Adds friction to initial setup

**Consequences**:
- Default path: `C:\Users\<username>\.password-store`
- Users can override via `PASSWORD_STORE_DIR`
- Must handle path creation if it doesn't exist
- Must handle both forward and backward slashes in paths

**Related**: pass-replacement-spec.md (Section 1.1)

---

### [AD-003] GPG Recipient Selection

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: How should the GPG recipient (encryption key) be selected?

**Decision**: Use the default GPG key, with `PASS_GPG_ID` environment variable override.

**Rationale**:
- Simplest approach for most users (single key)
- Matches Unix pass behavior (uses default key if not specified)
- Environment variable allows power users to specify a different key
- Avoids prompting user for key selection on every operation

**Alternatives Considered**:
1. **Require `PASS_GPG_ID` to be set**: More secure but less user-friendly
2. **Prompt user to select key on first use**: Adds setup friction
3. **Encrypt to all available keys**: Overkill, may not be desired
4. **Use symmetric encryption**: Less secure, defeats purpose of GPG key pairs

**Consequences**:
- Most users won't need to configure anything
- Power users can set `PASS_GPG_ID` for specific keys
- If multiple keys exist, GPG will use the default
- Must handle case where no default key exists (error with helpful message)

**Related**: pass-replacement-spec.md (Section 1.2)

---

### [AD-004] File Encryption Format

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: Should we use ASCII-armored (--armor) or binary GPG encryption?

**Decision**: Use ASCII-armored GPG encryption (`--armor` flag).

**Rationale**:
- Unix pass uses `--armor` by default
- ASCII-armored files are text-based and easier to work with
- Can be viewed/edited with text tools (though encrypted)
- More portable across systems
- Binary files may cause issues with some text processing tools

**Alternatives Considered**:
1. **Binary GPG files**: More compact, but less compatible with text tools
2. **User choice via config**: Adds complexity for minimal benefit

**Consequences**:
- Encrypted files will have `.gpg` extension and contain ASCII text
- Files will be slightly larger than binary GPG files
- Compatible with Unix pass password stores (can share stores across platforms)

**Related**: pass-replacement-spec.md (Section 3.2), pass-spec-kit.md (Open Question 1)

---

### [AD-005] Clipboard Implementation

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: How should clipboard operations be implemented on Windows?

**Decision**: Use the built-in `clip` command as the primary method, executed via Go's `exec.Command`.

**Rationale**:
- `clip` is available on all modern Windows versions
- Simple and reliable
- Works via standard subprocess execution in Go
- No additional dependencies required

**Alternatives Considered**:
1. **Win32 API via Go syscall**: More complex, platform-specific
2. **Third-party tools (xclip for Windows)**: Adds external dependencies
3. **COM object via WScript.Shell**: Works but more verbose and Windows-specific
4. **Go clipboard library**: External dependency, but could be an option

**Consequences**:
- Implementation: Use `exec.Command("clip")` with stdin pipe
- Must handle multi-line content correctly
- Must handle special characters in passwords
- Will work in any Windows shell context
- For cross-platform, may need platform-specific implementations

**Related**: pass-replacement-spec.md (Section 5)

---

### [AD-006] Clipboard Auto-Clear

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: Should the clipboard be automatically cleared after copying a password, and if so, when?

**Decision**: Yes, auto-clear the clipboard after 45 seconds by default, configurable via `PASS_CLIPBOARD_TIMEOUT` environment variable. Set `PASS_CLIPBOARD_TIMEOUT=0` to disable.

**Rationale**:
- Security best practice: minimize time sensitive data spends in clipboard
- 45 seconds is a reasonable balance between security and usability
- Configurable to allow users to disable if they prefer
- Matches behavior of other password managers

**Alternatives Considered**:
1. **No auto-clear**: Less secure, but simpler
2. **Clear immediately after paste**: Not possible to detect paste events reliably
3. **Clear on next pass command**: Unpredictable timing
4. **Different default timeout**: 30s, 60s, etc. - 45s is a common choice

**Consequences**:
- After copying, start a timer
- After timeout, overwrite clipboard with random data or empty string
- Requires background process or PowerShell job
- Must handle cases where script exits before timeout
- User can disable via environment variable

**Related**: pass-replacement-spec.md (Section 5.1), pass-spec-kit.md (Open Question 3)

---

### [AD-002] Re-evaluation: Switch from PowerShell to Go

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: User requested that the implementation should NOT be a PowerShell script, and suggested Go as a good alternative.

**Decision**: Switch implementation language from PowerShell to Go (Golang).

**Rationale**:
- User explicitly requested non-PowerShell implementation
- Go provides better cross-platform support
- Single binary distribution is more portable
- No runtime dependencies (unlike PowerShell which requires PowerShell runtime)
- Better for distribution and deployment
- Strong standard library for systems programming
- User preference for Go

**Alternatives Considered**:
1. **Stick with PowerShell**: Rejected due to user requirement
2. **Python**: Cross-platform but requires Python runtime
3. **Rust**: Good alternative but more complex
4. **C#/.NET**: Windows-centric, requires runtime

**Consequences**:
- All spec documents updated to reflect Go implementation
- Project structure changed from PowerShell modules to Go packages
- Build process changed from script to compilation
- Installation changed from script placement to binary placement
- Better cross-platform compatibility

**Supersedes**: AD-001 (original PowerShell decision)

**Related**: All spec documents

---

## Design Decisions

### [DD-001] Command Syntax Compatibility

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: Should the command syntax exactly match Unix pass, or can we make Windows-specific adjustments?

**Decision**: Match Unix pass syntax as closely as possible, with Windows-specific adjustments only where necessary.

**Rationale**:
- Users familiar with Unix pass can use this tool without learning new syntax
- Consistency across platforms
- Only differences should be unavoidable (e.g., path separators)

**Alternatives Considered**:
1. **Windows-specific syntax**: Easier to implement but breaks compatibility
2. **Hybrid approach**: Support both syntaxes, adds complexity

**Consequences**:
- Paths use forward slashes (`/`) as in Unix, converted to backslashes internally
- Command names and flags match Unix pass
- Error messages should be similar to Unix pass
- Exit codes should match Unix pass where applicable

**Related**: pass-replacement-spec.md (Section 14.1)

---

### [DD-002] Default Command

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: What should happen when the user runs `pass <path>` without an explicit command?

**Decision**: Treat it as `pass show <path>` (the default command is show).

**Rationale**:
- Matches Unix pass behavior
- Most common operation (retrieving a password)
- Intuitive: `pass <path>` shows the password at that path
- Consistent with other CLI tools (e.g., `cat <file>`)

**Alternatives Considered**:
1. **Error if no command specified**: Less user-friendly
2. **Show help**: Not useful for the common case
3. **List all passwords**: Not what user expects when specifying a path

**Consequences**:
- `pass <path>` is equivalent to `pass show <path>`
- Must handle this in argument parsing
- Help text should reflect this

**Related**: pass-replacement-spec.md (Section 2.3)

---

### [DD-003] Git Commit Messages

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: What format should git commit messages use?

**Decision**: Use `"Add <path>"` for insert operations.

**Rationale**:
- Simple and descriptive
- Matches common git commit message conventions
- Consistent with Unix pass behavior
- Easy to parse if needed

**Alternatives Considered**:
1. **`"pass insert <path>"`**: More explicit but redundant
2. **`"New password: <path>"`**: Similar to chosen option
3. **Customizable via config**: Adds complexity for minimal benefit
4. **Include timestamp**: Not necessary, git handles this

**Consequences**:
- Commit message: `Add email/gmail.com/myemail@gmail.com`
- Consistent across all insert operations
- Easy to filter in git log

**Related**: pass-replacement-spec.md (Section 4.2)

---

### [DD-004] Error Message Format

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: What format should error messages use?

**Decision**: Use `pass: <message>` format, matching Unix pass.

**Rationale**:
- Consistency with Unix pass
- Easy to parse in scripts
- Clear that the error came from pass
- Standard error message format

**Alternatives Considered**:
1. **`Error: <message>`**: Different from Unix pass
2. **`[pass] <message>`**: More verbose
3. **Just `<message>`**: Less clear source of error

**Consequences**:
- All errors start with `pass: `
- Examples: `pass: file not found`, `pass: password verification failed`
- Scripts can grep for `pass: ` to detect errors

**Related**: pass-replacement-spec.md (Section 7.1)

---

### [DD-005] Path Handling

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: How should path separators be handled?

**Decision**: Accept both forward slashes (`/`) and backslashes (`\`) in input, normalize to OS separator internally, but store paths with forward slashes in the password store.

**Rationale**:
- Forward slashes match Unix pass behavior
- Backslashes are native to Windows, so should be accepted
- Internal normalization ensures consistent behavior
- Storing with forward slashes allows potential cross-platform compatibility

**Alternatives Considered**:
1. **Only accept forward slashes**: Less Windows-friendly
2. **Only accept backslashes**: Less Unix-like
3. **Convert all to backslashes**: Breaks cross-platform compatibility

**Consequences**:
- Input: `pass insert email\gmail.com\user` or `pass insert email/gmail.com/user` both work
- Internal: Normalized to OS path separator for file operations
- Storage: Files use forward slashes in their logical paths
- Display: Use forward slashes in output (ls, find, etc.)

**Related**: pass-replacement-spec.md (Section 3.1)

---

### [DD-006] Invalid Character Handling

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: How should invalid filesystem characters in paths be handled?

**Decision**: Replace Windows reserved characters (`< > : " | ? *`) with underscores (`_`) and warn the user.

**Rationale**:
- Prevents file creation errors
- Preserves most of the intended path
- Warning informs user of the change
- Simple to implement

**Alternatives Considered**:
1. **Reject paths with invalid characters**: More strict, but may break user expectations
2. **URL-encode invalid characters**: Less readable
3. **Remove invalid characters**: Could cause ambiguity
4. **Prompt user for replacement**: Adds friction

**Consequences**:
- `pass insert aux:<con>` becomes `aux_con` with a warning
- User must be aware of the transformation
- Could cause confusion if user expects exact path
- Document this behavior

**Related**: pass-replacement-spec.md (Section 3.1)

---

## Implementation Decisions

### [ID-001] Git Commit on Insert

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: Should git commit failures cause the insert operation to fail?

**Decision**: No, git commit failures should be non-fatal. Warn the user but still consider the insert successful.

**Rationale**:
- The primary goal of insert is to store the password securely
- Git is a secondary feature (version control)
- User may not have git configured properly
- User may not want git integration
- Git errors shouldn't prevent password storage

**Alternatives Considered**:
1. **Fail insert on git error**: More strict, but could prevent password storage
2. **Skip git if not configured**: Similar to chosen, but less explicit
3. **Retry git commit**: Adds complexity, may not succeed

**Consequences**:
- Insert succeeds even if git commit fails
- Warning message: `pass: warning: git commit failed: <error>`
- User can use `--no-commit` to explicitly skip git
- User can manually commit later

**Related**: pass-replacement-spec.md (Section 4.2)

---

### [ID-002] GPG Passphrase Prompting

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: How should GPG passphrase prompting be handled?

**Decision**: Rely on gpg-agent for passphrase caching and prompting. Let GPG handle the passphrase prompt directly.

**Rationale**:
- gpg-agent is designed for this purpose
- Provides consistent behavior with other GPG operations
- Handles passphrase caching (configurable timeout)
- Secure passphrase entry
- No need to reimplement passphrase handling

**Alternatives Considered**:
1. **Prompt in PowerShell, pass to GPG**: More control but less secure
2. **Require passphrase in environment variable**: Insecure
3. **Disable passphrase caching**: Less user-friendly

**Consequences**:
- GPG will prompt for passphrase when gpg-agent cache is empty/expired
- User can configure gpg-agent cache timeout
- No passphrase handling code in our tool
- Consistent with Unix pass behavior

**Related**: pass-replacement-spec.md (Section 1.2)

---

### [ID-003] Temporary File Security

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: How should temporary files containing plaintext passwords be handled?

**Decision**: Create temporary files in a secure location (temp directory), use cryptographically random names, and securely delete them (overwrite before deletion).

**Rationale**:
- Minimizes risk of plaintext password exposure
- Random names prevent predictable file locations
- Secure deletion prevents recovery from disk
- Temp directory is typically on same volume (required for secure delete)

**Alternatives Considered**:
1. **Keep in memory only**: Not possible for GPG encryption (requires file)
2. **Use RAM disk**: Not available on all systems
3. **Simple delete**: Less secure, password could be recovered

**Consequences**:
- Temp files created in `%TEMP%` or `$env:TEMP`
- File names: random GUID or similar
- Secure deletion: overwrite with random data multiple times before delete
- Must handle errors during secure deletion (still delete if overwrite fails)

**Related**: pass-replacement-spec.md (Section 3.2), (Section 6.2)

---

### [ID-004] Multi-line Password Support

**Status**: Accepted  
**Date**: 2026-06-05  
**Context**: How should multi-line passwords be handled?

**Decision**: Support multi-line passwords via `-m/--multiline` flag on insert. Read until empty line or EOF.

**Rationale**:
- Some passwords/credentials are multi-line (e.g., private keys, certificates)
- Unix pass supports this
- Simple to implement in PowerShell
- Clear flag to indicate multi-line mode

**Alternatives Considered**:
1. **Always support multi-line**: Ambiguous when user finishes input
2. **Use special delimiter**: Less intuitive
3. **Only single-line**: Limits functionality

**Consequences**:
- Without `-m`: Single-line input (Enter submits)
- With `-m`: Multi-line input (empty line or Ctrl+D/Ctrl+Z submits)
- Multi-line passwords stored and retrieved as-is
- Clipboard handles multi-line content correctly

**Related**: pass-replacement-spec.md (Section 2.2)

---

## Open Decisions

These decisions have not yet been finalized and need further discussion:

### [OD-001] Symlink Support

**Status**: Pending  
**Date**: 2026-06-05  
**Context**: Should we support symlinks in the password store?

**Options**:
1. Support symlinks (Unix pass does)
2. Ignore symlinks (treat as regular files)
3. Warn about symlinks but allow them
4. Reject symlinks

**Considerations**:
- Unix pass supports symlinks
- Windows symlink support varies (requires admin for some operations)
- Symlinks could be used for organization
- Security implications of symlinks

---

### [OD-002] Cross-Platform Store Compatibility

**Status**: Pending  
**Date**: 2026-06-05  
**Context**: Should a password store created on Windows be usable on Unix, and vice versa?

**Options**:
1. Full compatibility (goal)
2. Best-effort compatibility
3. Windows-only, no compatibility guarantee

**Considerations**:
- File paths: Unix uses `/`, Windows uses `\`
- Line endings: Unix LF, Windows CRLF
- GPG behavior differences between platforms
- Git behavior differences

**Current Decision**: Aim for compatibility, but it's a nice-to-have, not a requirement.

---

## Superseded Decisions

None yet.

---

## Rejected Decisions

None yet.

---

*Decision Log Version: 1.0*  
*Last Updated: 2026-06-05*
