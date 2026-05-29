# Plan: Fast Go-based Shell Prompt for Git

## Approach

Replace the per-prompt `git` subprocess calls with a single statically-linked
Go binary (`gitprompt`) that reads `.git/` internals directly. The shell
script delegates to the binary when available and falls back to the existing
shell implementation otherwise — so users without a built binary keep
working behaviour.

## Architecture

```
shell/git_prompt.sh    — sources at shell start; detects binary once
shell/gitprompt[.exe]  — compiled Go binary (gitignored)
gitprompt/main.go      — Go source
gitprompt/Makefile     — build for Linux dev machines
gitprompt/build.sh     — POSIX build script (Windows ash, no make required)
```

## Binary Discovery

At source time, `git_prompt.sh` resolves `GIT_PROMPT_BIN` once via
`command -v gitprompt` and caches the path. Per-prompt cost is one `[ -n
"$GIT_PROMPT_BIN" ]` test plus the binary invocation — no `command -v`
fork on the hot path.

Users put `tools/shell/` on their `PATH` (which they likely already do
since `git_prompt.sh` lives there). The variable can be set explicitly
to override.

## Git Reading Strategy

Use `github.com/go-git/go-git/v5` for correctness:

- Handles loose objects, pack files, delta chains, packed-refs, worktrees
- Resolves annotated tag indirection
- Reads upstream config from `.git/config`

Trading binary size (~6 MB) for implementation correctness. On Linux the
binary is ~6 ms warm — within budget.

## Ahead / Behind

BFS the commit graph from both local and remote HEADs into two reachability
sets, then compute set difference. Cap traversal at **1000 commits** per
side to bound worst-case runtime; for a prompt we only care about
direction, not exact counts, so capping is acceptable.

## Tag Lookup

Build a `commit-hash → tag-name` map (dereferencing annotated tags), then
BFS from HEAD and return the first commit that appears in the map. Skipped
entirely when `GIT_PROMPT_DISABLE_TAGS=1`.

## Build System

Two equivalent entry points so neither platform group is left out:

- **Makefile** — preferred on Linux developer machines (`make build`)
- **build.sh** — POSIX `/bin/sh` script (Windows BusyBox ash, no make)

Both produce identical output. `build.sh` uses `go env GOEXE` so the
native build target is named `gitprompt.exe` on Windows automatically.

## Fallback Path

If the binary is missing or fails to start, the shell function falls
through to the original `git status --porcelain=v2` + `git describe`
implementation. No regression for users who haven't built the binary yet.

## Risks

- **go-git binary size** — 6 MB is acceptable for tooling; not user-facing
- **go-git Go version requirement** — pinned to Go 1.25; documented in
  build scripts
- **Pack format changes** — go-git tracks upstream git; low risk
- **Ahead/behind cap of 1000** — only affects repos with massive
  divergence; prompt still shows correct direction
