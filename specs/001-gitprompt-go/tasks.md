# Tasks: Fast Go-based Shell Prompt for Git

## Implementation

- [x] **T1** — Scaffold `gitprompt/` Go module with `go.mod` pinned to
      `github.com/go-git/go-git/v5`.
- [x] **T2** — Implement `main.go`:
      - [x] Open repo via `PlainOpenWithOptions{DetectDotGit: true}`
      - [x] Read HEAD; handle detached and pre-initial-commit cases
      - [x] Resolve upstream from `branch.<name>` config block
      - [x] Compute ahead/behind by BFS with 1000-commit cap
      - [x] Look up nearest tag (respecting `GIT_PROMPT_DISABLE_TAGS`)
      - [x] Emit `\033[;32m(branch[tag][sync])\033[0m`
      - [x] Exit silently outside a git repo
- [x] **T3** — Add `Makefile` with `build`, `build-all`, `clean` targets
      for linux/amd64, linux/arm64, windows/arm64.
- [x] **T4** — Add `build.sh` (POSIX `/bin/sh`) with the same targets, for
      Windows BusyBox ash where GNU make is unavailable. Use `go env GOEXE`
      to resolve the native binary suffix.
- [x] **T5** — Update `shell/git_prompt.sh`:
      - [x] Detect `GIT_PROMPT_BIN` once at source time
      - [x] Delegate to binary when set; keep shell fallback intact
      - [x] Document `GIT_PROMPT_BIN` and both build entry points
- [x] **T6** — Add `.gitignore` for compiled binaries
      (`shell/gitprompt*`).

## Verification

- [x] **V1** — Binary builds on Linux amd64 with `make build` and
      `./build.sh build`.
- [x] **V2** — Binary produces identical output to the shell fallback on
      the current repo (`(feat/gitprompt-go)`).
- [x] **V3** — Binary exits 0 with no output outside a git repo.
- [ ] **V4** — Cross-compile target `gitprompt-windows-arm64.exe` builds
      cleanly (deferred until user tests on their Windows ARM64 machine).
- [ ] **V5** — `build.sh build` on BusyBox ash produces a working
      `gitprompt.exe` (deferred to user verification).

## Documentation

- [x] **D1** — `spec.md`, `plan.md`, `tasks.md` under
      `specs/001-gitprompt-go/`.
- [x] **D2** — Inline comments in `git_prompt.sh` covering both build
      entry points and PATH requirement.
