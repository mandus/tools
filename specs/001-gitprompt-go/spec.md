# Spec: Fast Go-based Shell Prompt for Git

## Problem

The current `__git_branch` function in `shell/git_prompt.sh` forks 1–2 `git`
subprocesses on every interactive prompt. On Windows ARM64 (BusyBox ash) and
WSL, each `git` invocation costs ~50–150 ms of process-creation overhead,
making the prompt visibly sluggish. On native Linux the cost is lower
(~3 ms) but still adds up over a session.

## Goal

Reduce prompt-generation time to a level the user does not perceive
(<10 ms) on every supported platform, while preserving the existing
visual format and feature set.

## Supported Platforms

The user runs this prompt on three platforms:

| OS      | Arch  | Shell           |
|---------|-------|-----------------|
| Linux   | amd64 | bash / ash      |
| Linux   | arm64 | bash / ash      |
| Windows | arm64 | BusyBox ash     |

Notably, the Windows environment does **not** have GNU make available, so
the build path must not depend on it.

## Visual Behaviour (Unchanged)

The prompt segment must look identical to the current shell version:

```
\033[;32m(<branch>[(<tag>)][<sync>])\033[0m
```

- `<branch>` — current branch name, or `HEAD` if detached or pre-initial
- `(<tag>)` — nearest reachable annotated/lightweight tag, optional
- `<sync>` — one of `=`, `>`, `<`, `<>` (in-sync, ahead, behind, diverged),
  or empty when there is no configured upstream
- Outside a git repo: print nothing, exit 0

## Configuration (Unchanged)

- `GIT_PROMPT_DISABLE_TAGS=1` — skip tag lookup
- New: `GIT_PROMPT_BIN` — override path to the Go binary

## Non-Goals

- Replacing the existing `gco` fzf-based checkout helper
- Showing dirty-tree state (untracked / staged / modified files)
- Exact ahead/behind counts (only direction matters for the prompt)
- Supporting platforms beyond the three listed above
