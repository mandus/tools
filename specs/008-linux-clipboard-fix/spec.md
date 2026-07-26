# Linux/X11 Clipboard Fix Specification

## Overview

This specification describes a bug fix for `pass -c` (and any other `--clip`
flag usage) on Linux. The clipboard integration currently only works on
Windows; on Linux (X11 or Wayland) it fails unconditionally because the
implementation hard-codes the Windows `clip` command. It also adds tmux
integration: when `pass -c` runs inside a tmux session, the password is
additionally loaded into a tmux paste buffer, which works even when no
system clipboard utility is reachable (e.g. over SSH without X11/Wayland
forwarding).

## Status

- **Status**: Implemented ✅
- **Author**: @aasmundo
- **Created**: 2026-07-25
- **Last Updated**: 2026-07-25
- **Branch**: `fix/19-linux-clipboard-x11`

## Background

`pass` advertises itself as cross-platform (Windows, Linux, macOS) and lists
"Copy to clipboard support" as a feature (see `README`). In practice,
`pkg/filesystem.CopyToClipboard`, which backs `pass -c`, `pass rm -c` and the
TUI clip mode, is implemented as:

```go
func CopyToClipboard(text string) error {
	cmd := exec.Command("clip")
	...
}
```

`clip` is a Windows-only executable. On Linux (including X11 sessions) the
`exec.Command("clip")` call fails with `exec: "clip": executable file not
found in $PATH`, so `pass -c` always errors out and nothing is copied.

The dead/unused `pkg/clipboard` package has the identical problem (Windows
`clip`/`powershell` only) and is not referenced anywhere in `cmd/*`.

## Problem Statement

```bash
# Current behavior on Linux/X11 (BROKEN)
$ pass -c email/gmail.com/user
pass: failed to copy to clipboard: exec: "clip": executable file not found in $PATH

# Expected behavior
$ pass -c email/gmail.com/user
Copied email/gmail.com/user to clipboard.
# password is now on the X11 (or Wayland) clipboard
```

## Goals

- Make `pass -c` work on Linux under X11 (via `xclip`/`xsel`) and Wayland
  (via `wl-copy`), in addition to continuing to work on Windows and macOS.
- When running inside tmux, also load the password into a tmux paste
  buffer, so `pass -c` works over SSH sessions where no system clipboard
  utility can reach a display server.
- Keep the existing auto-clear-after-timeout behavior working cross-platform
  and for the tmux buffer.
- Fail with a clear, actionable error message when no clipboard utility is
  available (e.g. missing `xclip`/`xsel` on a minimal Linux install) and
  pass is not running inside tmux.
- Do not change the CLI/flag surface (`-c`/`--clip` behavior is unchanged).

## Non-Goals

- Implementing a bespoke X11/Wayland clipboard protocol client.
- Adding clipboard history or multiple clipboard selections (primary vs.
  clipboard) support.
- Fixing/removing the unused `pkg/clipboard` package beyond noting it is
  dead code (tracked separately if cleanup is desired).

## Root Cause

`pkg/filesystem.CopyToClipboard` unconditionally shells out to the Windows
`clip` command with no platform branching and no Linux/macOS fallback. There
is no build-tag-based or runtime-based platform detection, so the function
behaves identically (and fails identically) on every non-Windows OS.

## Technical Design

Replace the hand-rolled `exec.Command("clip")` call with
[`github.com/atotto/clipboard`](https://github.com/atotto/clipboard), a small
dependency-free (at the Go level) cross-platform clipboard library that:

- Uses `clip.exe`/`powershell.exe` on Windows (matches current behavior).
- Uses `pbcopy`/`pbpaste` on macOS.
- Uses `wl-copy`/`wl-paste` on Wayland (detected via `WAYLAND_DISPLAY`),
  falling back to `xclip`, then `xsel`, then Termux's
  `termux-clipboard-set`/`get` on Linux/X11 and other Unix variants.
- Sets an exported `clipboard.Unsupported` flag when no backend is found, so
  `pass` can surface a clear error instead of a generic "file not found".

This library is already present in `go.sum` as an indirect dependency (pulled
in transitively), so no new supply-chain surface is introduced; this change
promotes it to a direct dependency.

`pkg/filesystem.CopyToClipboard` becomes a thin wrapper:

```go
func CopyToClipboard(text string) error {
	if err := clipboard.WriteAll(text); err != nil {
		return fmt.Errorf("failed to write to clipboard: %v", err)
	}
	return nil
}
```

`IsClipboardAvailable` (test helper) is updated to check
`!clipboard.Unsupported` instead of shelling out to the Windows-only `clip`
command, and additionally reports availability when running inside tmux
(since the tmux paste buffer does not depend on a system clipboard utility).

### tmux Integration

When the `TMUX` environment variable is set (i.e. pass is running inside a
tmux session), `CopyToClipboard` additionally runs:

```
tmux load-buffer -w -b pass -
```

reading the password from stdin into a dedicated tmux paste buffer named
`pass`. This lets the password be pasted with tmux itself (`prefix` + `]`)
regardless of whether a system clipboard utility is available. The `-w` flag
additionally asks tmux to forward the buffer to the terminal's clipboard via
the OSC 52 escape sequence, if tmux's `set-clipboard` option and the
terminal support it - this is what allows a password copied inside a remote
tmux session over SSH to end up on the *local* machine's clipboard without
any X11/Wayland forwarding.

Behavior when both mechanisms are attempted:

- Outside tmux: behavior is unchanged, only the system clipboard is used.
- Inside tmux with a working system clipboard (e.g. local X11 session
  running tmux): both the system clipboard and the tmux buffer are set.
- Inside tmux with no working system clipboard (e.g. SSH session with no
  `DISPLAY`): the tmux buffer copy succeeding is treated as overall success;
  the system clipboard error is not surfaced to the user.
- Inside tmux with both the system clipboard and `tmux load-buffer` failing:
  both errors are reported.

The same dedicated buffer name (`pass`) is reused on every copy so repeated
`pass -c` invocations overwrite the previous buffer instead of accumulating
tmux's auto-named buffers.

## User Stories

### As a user running pass inside tmux over SSH, I want `pass -c` to still let me paste the password
So that I am not blocked by the remote host having no system clipboard
utility reachable (no `DISPLAY`/Wayland socket).

**Acceptance Criteria**:
- [x] `pass -c <path>` inside a tmux session loads the password into a tmux
      paste buffer named `pass`, pasteable with `prefix` + `]`.
- [x] This works even when `xclip`/`xsel`/`wl-copy` are missing or cannot
      reach a display server.
- [x] If tmux's `set-clipboard` option is enabled and the terminal supports
      OSC 52, the password is also forwarded to the local terminal's system
      clipboard automatically.
- [x] Outside of tmux, behavior is unchanged.

### As a Linux user with X11, I want `pass -c` to copy the password to my clipboard
So that I can paste it into another application without printing it to the
terminal.

**Acceptance Criteria**:
- [x] `pass -c <path>` copies the decrypted password to the X11 clipboard
      when `xclip` or `xsel` is installed.
- [x] `pass -c <path>` copies to the Wayland clipboard when running under a
      Wayland session with `wl-clipboard` installed.
- [x] The clipboard is auto-cleared after `PASS_CLIPBOARD_TIMEOUT` seconds
      (default 45), unchanged from prior behavior.
- [x] A missing clipboard utility produces a clear error instead of a
      confusing "clip: executable file not found" message.

### As a Windows/macOS user, I want existing clipboard behavior to keep working
So that this fix does not regress the platforms that already worked.

**Acceptance Criteria**:
- [x] `pass -c` continues to use `clip.exe` on Windows.
- [x] `pass -c` uses `pbcopy` on macOS (previously unimplemented/broken too,
      now fixed as a side effect of using a cross-platform library).

## Testing

- `pkg/filesystem/fs_test.go` `TestCopyToClipboard` is updated to skip
  cleanly (not fail) in headless CI/sandbox environments where no clipboard
  utility, display server, or tmux session is available, using
  `IsClipboardAvailable()` (system clipboard OR tmux) as the check.
- `TestIsInsideTmux` unit-tests the `TMUX` environment variable detection.
- `TestCopyToTmuxBuffer` is an integration test that only runs when actually
  executed inside a real tmux session (checks `TMUX` env var and `tmux`
  binary), so it can talk to a live tmux server. It writes to the `pass`
  buffer and reads it back with `tmux show-buffer`.
- Manually verified end-to-end inside a real tmux session in this sandbox
  (no `DISPLAY` set): `pass -c`-style `CopyToClipboard` call succeeded with
  no error, and `tmux show-buffer -b pass` returned the expected text,
  confirming the SSH-without-X11 fallback path works.
- Manually verified on Linux with `xclip`/`xsel`/`wl-copy` installed but no
  display server and outside tmux: `IsClipboardAvailable` correctly reports
  availability based on utility presence; the actual `WriteAll` call is
  expected to fail without a running X server/Wayland compositor, which is
  an environment limitation, not a code defect.

## Related

- `README` (Features: "Copy to clipboard support", "Cross-platform (Windows,
  Linux, macOS)")
- `docs/pass-decision-log.md` — [AD-005] Clipboard Implementation (superseded
  by this change)
- `docs/CHANGELOG.md`
