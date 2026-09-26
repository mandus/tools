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
cd evergarden && ./install.sh          # copy into ~/.vim/colors on every platform
./install.sh status                    # where is it?
./install.sh uninstall                 # remove
```

The installer uses only `~/.vim/colors/`, creating it if necessary, on Windows
and Unix alike. It never creates or modifies `~/vimfiles`. Uninstall removes
only `~/.vim/colors/evergarden.vim`, leaving directories and other files intact.

Or manually: drop `colors/evergarden.vim` into `~/.vim/colors/`.

### Not listed by `:colorscheme <Tab>`?

Check inside the Vim you use:

```vim
:echo expand('~')
:set runtimepath?
:echo globpath(&runtimepath, 'colors/evergarden.vim')
```

`runtimepath` must include `~/.vim`, the **parent** of `colors`, not the
`colors` directory itself. Native Windows Vim does not include `~/.vim` by
default.

Plugins under `~/.vim/bundle` can still work: managers such as Vundle add each
plugin's runtime root separately, without adding `~/.vim` itself. To keep
using `~/.vim/colors` on Windows, add this before the colorscheme in your vimrc:

```vim
set runtimepath^=~/.vim
colorscheme evergarden
```

If Vim uses a different home from the shell or a custom `runtimepath`, install
manually into a runtime root shown by `:set runtimepath?`.

## Use

In `~/.vimrc` (or `~/_vimrc` on Windows):

```vim
set runtimepath^=~/.vim   " needed when ~/.vim is not already on runtimepath
set termguicolors         " strongly recommended
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

## Installer regression tests

```sh
sh evergarden/test/install.sh
```

Tests use temporary homes to check install/status/uninstall into `~/.vim/colors`
and verify that `~/vimfiles` is untouched. If `vim` is available, they also
verify discovery, completion, loading, and removal with the documented
`set runtimepath^=~/.vim` configuration.

## Spec

`specs/010-evergarden-colorscheme/spec.md`
