# evergarden Vim Colorscheme Specification

## Overview

A Vim colorscheme that ports the [evergarden](https://evergarden.moe) palette —
a forest-toned dark scheme inspired by *The Legend of Zelda: The Minish Cap* —
into a single self-contained Vimscript file for **classic Vim only**.

## Status

- **Status**: Implemented ✅
- **Created**: 2026-02-14
- **Branch**: `feat/010-evergarden-colorscheme`

## Background

`curl -L evergarden.moe` serves the palette as ANSI-coloured plain text. The
upstream project ships a Lua-only Neovim plugin, which is useless to a classic
Vim user. This repository needs a dependency-free Vimscript variant that works
in Vim 8+, both GUI and 256-colour terminals.

## Goals

- Reproduce the upstream palette exactly (24 named colours).
- Single file, no plugin manager, no Lua, no runtime dependencies.
- Work in Vim 8.0+, GUI and terminal.
- Degrade gracefully to xterm-256 approximations when `termguicolors` is off.
- Cover every group in `:help highlight-groups`, plus diffs, spell, netrw,
  popup/toolbar, and `termdebug`.
- Ship an installer that is symmetric: install and uninstall.

## Non-Goals

- **Neovim.** No Lua module, no `@*` Tree-sitter captures, no `Diagnostic*` or
  `Lsp*` groups, no `NormalNC`/`Winbar`/`MsgArea`/`WinSeparator`. Tree-sitter
  capture names are not even legal group names in Vim (`W18: Invalid character
  in group name`).
- Light variant (upstream publishes the dark palette only).
- Third-party plugin groups.
- 16-colour terminal support.

## Palette

Verbatim from `evergarden.moe`, with the xterm-256 fallback chosen by nearest
sRGB distance and then hand-adjusted so adjacent hues stay distinguishable.

| Name | Hex | cterm |
|---|---|---|
| red | `#f57f82` | 210 |
| orange | `#f7a182` | 216 |
| yellow | `#f5d098` | 222 |
| lime | `#dbe6af` | 187 |
| green | `#cbe3b3` | 151 |
| aqua | `#b3e3ca` | 158 |
| skye | `#b3e6db` | 152 |
| snow | `#afd9e6` | 153 |
| blue | `#b2caed` | 111 |
| purple | `#d2bdf3` | 183 |
| pink | `#f3c0e5` | 218 |
| cherry | `#fae6ef` | 255 |
| text | `#f8f9e8` | 255 |
| subtext1 | `#adc9bc` | 145 |
| subtext0 | `#96b4aa` | 109 |
| overlay2 | `#839e9a` | 246 |
| overlay1 | `#6f8788` | 66 |
| overlay0 | `#58686d` | 241 |
| surface2 | `#4a585c` | 240 |
| surface1 | `#374145` | 238 |
| surface0 | `#262f33` | 236 |
| base | `#1e2528` | 235 |
| mantle | `#191e21` | 234 |
| crust | `#171c1f` | 233 |

## Semantic Mapping

| Role | Colour |
|---|---|
| Normal text | text on base |
| Comment | overlay1, italic |
| Keyword / Statement / Conditional / Repeat / Label | `g:evergarden_keyword`, default skye |
| Function | blue |
| String | green |
| Constant / Number / Boolean | orange |
| Type / StorageClass | yellow |
| PreProc / Macro / Include | pink |
| Operator | subtext0 |
| Delimiter | overlay2 |
| Special / SpecialChar | purple |
| Tag | aqua |
| Identifier / Variable | subtext1 |
| Error / Exception | red |

### Keyword hue

Control-flow keywords are the most repeated token on screen, so the hue is a
single option rather than a hard-coded choice. The first draft used purple,
which read as too loud against the green/orange of strings and literals. The
default is now skye (`#b3e6db`): cool, recessive, and distinct from blue
`Function` and green `String`.

`g:evergarden_keyword` accepts any palette name and falls back to skye if the
name is unknown, so a typo cannot produce an unreadable buffer. Operators moved
to subtext0 and delimiters to overlay2 to keep punctuation from competing;
Special/SpecialChar inherited purple so the hue is still present, but rarely.

`test/vimrc` exposes `:Kw <name>` (with completion) and `$EVERGARDEN_KEYWORD`
for auditioning without editing the scheme.

## Design

### Structure

```
evergarden/
  colors/evergarden.vim   # the colorscheme
  install.sh              # install / uninstall / status
  test/run.sh             # preview harness (no install required)
  test/sample.{c,vim,md}  # in-context syntax samples
  test/diff_{a,b}.txt     # diff-mode fixture
  README.md
```
`colors/evergarden.vim` defines the palette as a script-local dictionary and
emits highlights through one `s:hi()` helper that writes `guifg`/`guibg`/`gui`
and `ctermfg`/`ctermbg`/`cterm` in a single `highlight` command. `NONE` is used
for absent attributes so the file is safe to re-source.

### Italics

Comments and special comments use italics, gated behind `g:evergarden_italic`
(default `1`). Terminals without italic support render them upright; no
fallback logic is needed.

### Transparency

`g:evergarden_transparent` (default `0`) clears the background of `Normal`,
`SignColumn`, `EndOfBuffer`, `LineNr`, `FoldColumn`, `VertSplit`, and
`Terminal` so the terminal's own background shows through.

### Vim only

The group list is taken from `:help highlight-groups`. Neovim-only names are
omitted entirely: `@*` captures raise `W18` in Vim, and `Diagnostic*`, `Lsp*`,
`NormalNC`, `NormalFloat`, `FloatBorder`, `MsgArea`, `Winbar`, `WinSeparator`,
`Substitute` are dead weight. Vim-specific groups are included instead:
`LineNrAbove`/`LineNrBelow`, `PmenuMatch`/`PmenuMatchSel`, `StatusLineTerm`,
`Terminal`, `Popup`/`PopupSelected`/`PopupNotification`, `MessageWindow`,
`ToolbarLine`/`ToolbarButton`, `Menu`/`Scrollbar`/`Tooltip`, `debugPC`,
`debugBreakpoint`.

### Terminal Colours

`g:terminal_ansi_colors` is set from the palette (guarded by `has('terminal')`)
so `:terminal` buffers match.

### Installation

`install.sh` copies (not symlinks — Windows-friendly) the colorscheme into
every detected Vim runtime directory:

- `$HOME/vimfiles/colors/` (Windows Vim)
- `$HOME/.vim/colors/` (Unix Vim)

Subcommands: `install` (default), `uninstall`, `status`.

## Testing

Visual verification uses Vim's own runtime files rather than a bespoke
preview: `syntax/hitest.vim` renders every defined group in its own colours,
and `syntax/colortest.vim` renders the cterm grid. `test/run.sh` wires those
together with in-context samples, loading the scheme via `set rtp+=.` so
nothing has to be installed.

| Case | Expectation |
|---|---|
| `install.sh install` | File present in each detected colors dir |
| `install.sh status` | Reports installed/missing per target |
| `install.sh uninstall` | Files removed; empty `colors/` dirs left in place |
| Re-run `install.sh install` | Idempotent overwrite, exit 0 |
| `test/run.sh` | Five tabs, no errors, all groups legible |
| `test/run.sh 256` | Same, via cterm fallback |
| `:Kw <name>` in the preview | Keyword hue changes, no errors |
| `let g:evergarden_keyword = 'bogus'` | Silently falls back to skye |
| `vim -Nu NONE -c 'set tgc' -c 'colo evergarden' -c q` | Exits 0, no errors |
| `:colorscheme evergarden` twice | No `E` errors (re-sourcing is safe) |
| Plain Vim 8, `:colorscheme evergarden` | No `W18` warnings |
| Vim without `+terminal` | No error from the `g:terminal_ansi_colors` block |
| `&t_Co == 256`, no termguicolors | cterm attributes applied |

## References

- <https://evergarden.moe>
- `:help highlight-groups`, `:help group-name`
