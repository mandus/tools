# evergarden.vim

A classic-Vim port of the [evergarden](https://evergarden.moe) palette — forest
tones, soft pastels, dark background. Single file, no dependencies, no plugin
manager. Vim 8.0+. **No Neovim support, by design.**

```
red #f57f82   orange #f7a182   yellow #f5d098   lime #dbe6af
green #cbe3b3  aqua #b3e3ca    skye #b3e6db     snow #afd9e6
blue #b2caed   purple #d2bdf3  pink #f3c0e5     cherry #fae6ef
```

## Install

```bash
cd evergarden && ./install.sh          # copy into ~/.vim/colors (and ~/vimfiles/colors)
./install.sh status                    # where is it?
./install.sh uninstall                 # remove
```

Or manually: drop `colors/evergarden.vim` into `~/.vim/colors/`.

## Use

In `~/.vimrc`:

```vim
set termguicolors     " strongly recommended
colorscheme evergarden
```

Without `termguicolors` the scheme falls back to xterm-256 approximations
(needs `set t_Co=256`).

## Options

Set before `:colorscheme`:

| Option | Default | Effect |
|---|---|---|
| `g:evergarden_italic` | `1` | Italic comments and special comments |
| `g:evergarden_transparent` | `0` | Keep the terminal background |
| `g:evergarden_keyword` | `'skye'` | Hue for `if`/`for`/`switch`/`function` etc. Any palette name |

```vim
let g:evergarden_transparent = 1
let g:evergarden_italic = 0
let g:evergarden_keyword = 'blue'
colorscheme evergarden
```

### Keyword hue

Control flow is the most-repeated token on screen, so it gets its own knob.
The default `skye` (`#b3e6db`) keeps it cool and out of the way; `Function`
stays blue, `Type` yellow, `String` green, `Constant` orange, and operators
and delimiters are muted so the keyword hue reads cleanly.

| Value | Feel |
|---|---|
| `skye` (default) | Cool cyan, low contrast against strings |
| `aqua` | Greener, softer; closer to `String` |
| `blue` | Cooler still — but shares the hue with `Function` |
| `orange` | Warm; shares the hue with `Constant`/`Number` |
| `purple` | The original mapping |

Audition them live in the preview with `:Kw skye`, `:Kw blue`, … (tab completes).

## Test it without installing

```bash
cd evergarden
./test/run.sh          # true colour
./test/run.sh 256      # 256-colour cterm fallback
```

Five tabs (`gt` / `gT` to move, `:qa` to quit):

| Tab | Shows |
|---|---|
| `syntax/hitest.vim` | Every defined highlight group, each in its own colours (Vim built-in) |
| `sample.c` | Comment / PreProc / Type / Statement / Constant / String / Special |
| `sample.vim` | Vimscript: functions, options, autocommands, registers |
| `sample.md` | Headings, code, links, `spell` on so `SpellBad` shows |
| diff | `DiffAdd` / `DiffChange` / `DiffDelete` / `DiffText` side by side |

The two Vim built-ins are worth knowing on their own:

```vim
:runtime syntax/hitest.vim      " all groups in their own colours
:runtime syntax/colortest.vim   " the cterm colour grid
```

## Covers

Every group in `:help highlight-groups`, plus diffs, spell, netrw, popup
windows, the toolbar, `termdebug` markers, markdown/help filetypes, and
`g:terminal_ansi_colors` for `:terminal`.

## Spec

`specs/010-evergarden-colorscheme/spec.md`
