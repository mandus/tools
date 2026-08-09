# Windows `sleep` Command Specification

## Overview

This specification describes a small C utility, `sleep`, that provides a
GNU-`sleep`-compatible subset on Windows. Its primary consumer is the
[goal](https://codeberg.org/anaseto/goal) interpreter's `run` builtin, which
cannot delay execution because Windows ships no `sleep` executable.

## Status

- **Status**: Implemented ✅
- **Created**: 2026-08-09
- **Branch**: `feat/009-sleep-command`

## Background

In goal, shell-outs look like:

```k
run"curl""-V"
```

`run` execs the named program directly (it does not go through `cmd.exe`,
although Go's `LookPath` semantics mean `PATHEXT` entries such as `.cmd` are
resolved). On Windows there is no `sleep` on `%PATH%`, so:

```k
run"sleep""2"
```

fails with a "not available on %PATH%" style error.

## Problem Statement

The obvious Windows substitutes are all inadequate:

| Candidate | Problem |
|---|---|
| `timeout /t N` | Integer seconds only, and aborts with `ERROR: Input redirection is not supported` when stdin is not a console — precisely the case under `run`. |
| `ping -n N 127.0.0.1` | Integer seconds only; abuses the network stack; off-by-one semantics. |
| `powershell Start-Sleep` | ~200-400 ms of interpreter startup overhead per call. |
| `.cmd` shim | Resolvable via `PATHEXT`, but only wraps one of the above. |

## Goals

- Provide a real `sleep.exe` on `%PATH%` usable from goal's `run`.
- Support fractional seconds (`sleep 0.25`).
- Support the GNU suffixes `s`, `m`, `h`, `d`.
- Sum multiple arguments, as GNU `sleep` does (`sleep 1m 30`).
- Negligible startup overhead; no runtime dependencies.

## Non-Goals

- Full GNU coreutils compatibility (no `--version`, no locale handling).
- Cross-compilation. Unix platforms already ship `sleep`; the POSIX code path
  exists only so the source builds and can be tested anywhere.

## Design

Single translation unit, `sleep/sleep.c`:

- `parse_interval()` — `strtod` plus an optional single-character suffix.
  Rejects trailing garbage and negative values.
- `sleep_seconds()` — `Sleep()` on Windows, `nanosleep()` (with `EINTR`
  resume) elsewhere. Chunked at 1000 s so long waits cannot overflow the
  `DWORD` millisecond argument.
- Arguments are summed before sleeping.

### Interface

```
usage: sleep NUMBER[smhd]...
```

| Exit code | Meaning |
|---|---|
| 0 | Slept successfully, or `--help`/`-h` |
| 1 | Invalid time interval |
| 2 | No arguments given |

Errors are written to stderr as `sleep: invalid time interval '<arg>'`.

## Build & Install

Follows the existing per-tool convention (`gitprompt`, `pass`): a `build.sh`
and `Makefile` that emit into `../shell/`.

```sh
cd sleep && ./build.sh          # → ../shell/sleep[.exe]
./shell/make.sh                 # builds all tools, including sleep
```

Install by copying `shell/sleep.exe` to a directory on `%PATH%`.

## Acceptance Criteria

- [x] `goal -e 'say run"sleep""2"'` waits ~2 s and exits cleanly.
- [x] `goal -e 'run"sleep""0.25"'` honours fractional seconds.
- [x] `sleep 1m 30` sums to 90 s.
- [x] `sleep abc` prints an error and exits 1.
- [x] `sleep` with no arguments prints usage and exits 2.
- [x] `./shell/make.sh clean` removes the built binary.
