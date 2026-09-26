" evergarden.vim -- a forest-toned dark colorscheme for Vim
"
" Palette: https://evergarden.moe
" Spec:    specs/010-evergarden-colorscheme/spec.md
"
" Classic Vim only (8.0+). No Neovim groups, no Lua, no dependencies.
"
" Options:
"   let g:evergarden_italic      = 0       " disable italics (default 1)
"   let g:evergarden_transparent = 1       " keep terminal background (default 0)
"   let g:evergarden_keyword     = 'blue'  " keyword hue (default 'skye')
"                                          " any palette name: skye, aqua, blue,
"                                          " green, lime, orange, purple, pink...

hi clear
if exists('syntax_on')
  syntax reset
endif

set background=dark
let g:colors_name = 'evergarden'

let s:italic = get(g:, 'evergarden_italic', 1) ? 'italic' : 'NONE'
let s:transparent = get(g:, 'evergarden_transparent', 0)

" ---------------------------------------------------------------- palette ---
let s:p = {}
let s:p.red      = ['#f57f82', 210]
let s:p.orange   = ['#f7a182', 216]
let s:p.yellow   = ['#f5d098', 222]
let s:p.lime     = ['#dbe6af', 187]
let s:p.green    = ['#cbe3b3', 151]
let s:p.aqua     = ['#b3e3ca', 158]
let s:p.skye     = ['#b3e6db', 152]
let s:p.snow     = ['#afd9e6', 153]
let s:p.blue     = ['#b2caed', 111]
let s:p.purple   = ['#d2bdf3', 183]
let s:p.pink     = ['#f3c0e5', 218]
let s:p.cherry   = ['#fae6ef', 255]
let s:p.text     = ['#f8f9e8', 255]
let s:p.subtext1 = ['#adc9bc', 145]
let s:p.subtext0 = ['#96b4aa', 109]
let s:p.overlay2 = ['#839e9a', 246]
let s:p.overlay1 = ['#6f8788',  66]
let s:p.overlay0 = ['#58686d', 241]
let s:p.surface2 = ['#4a585c', 240]
let s:p.surface1 = ['#374145', 238]
let s:p.surface0 = ['#262f33', 236]
let s:p.base     = ['#1e2528', 235]
let s:p.mantle   = ['#191e21', 234]
let s:p.crust    = ['#171c1f', 233]
let s:p.none     = ['NONE', 'NONE']

" Keyword hue. Any palette name above; falls back to skye if unknown.
let s:kw = get(g:, 'evergarden_keyword', 'skye')
if !has_key(s:p, s:kw) || s:kw ==# 'none'
  let s:kw = 'skye'
endif

" ----------------------------------------------------------------- helper ---
" s:hi(group, fg, bg [, attr [, sp]])
function! s:hi(group, fg, bg, ...) abort
  let l:fg = s:p[a:fg]
  let l:bg = s:p[a:bg]
  let l:attr = a:0 > 0 && !empty(a:1) ? a:1 : 'NONE'
  let l:cmd = 'highlight ' . a:group
        \ . ' guifg=' . l:fg[0] . ' guibg=' . l:bg[0] . ' gui=' . l:attr
        \ . ' ctermfg=' . l:fg[1] . ' ctermbg=' . l:bg[1] . ' cterm=' . l:attr
  if a:0 > 1 && !empty(a:2)
    let l:cmd .= ' guisp=' . s:p[a:2][0]
  endif
  execute l:cmd
endfunction

let s:bg = s:transparent ? 'none' : 'base'

" -------------------------------------------------------------- ui chrome ---
call s:hi('Normal',          'text',     s:bg)
call s:hi('ColorColumn',     'none',     'surface0')
call s:hi('Conceal',         'overlay1', 'none')
call s:hi('Cursor',          'base',     'text')
call s:hi('lCursor',         'base',     'text')
call s:hi('CursorIM',        'base',     'text')
call s:hi('CursorColumn',    'none',     'surface0')
call s:hi('CursorLine',      'none',     'surface0')
call s:hi('Directory',       'blue',     'none')
call s:hi('EndOfBuffer',     'surface1', s:bg)
call s:hi('ErrorMsg',        'red',      'none', 'bold')
call s:hi('VertSplit',       'surface1', s:bg)
call s:hi('Folded',          'subtext0', 'surface0')
call s:hi('FoldColumn',      'overlay0', s:bg)
call s:hi('SignColumn',      'overlay0', s:bg)
call s:hi('IncSearch',       'base',     'orange')
call s:hi('CurSearch',       'base',     'orange')
call s:hi('Search',          'base',     'yellow')
call s:hi('LineNr',          'surface2', s:bg)
call s:hi('LineNrAbove',     'surface2', s:bg)
call s:hi('LineNrBelow',     'surface2', s:bg)
call s:hi('CursorLineNr',    'yellow',   'none', 'bold')
call s:hi('MatchParen',      'orange',   'surface2', 'bold')
call s:hi('ModeMsg',         'text',     'none', 'bold')
call s:hi('MoreMsg',         'green',    'none')
call s:hi('NonText',         'surface2', 'none')
call s:hi('Pmenu',           'subtext1', 'surface0')
call s:hi('PmenuSel',        'base',     'green', 'bold')
call s:hi('PmenuKind',       'purple',   'surface0')
call s:hi('PmenuKindSel',    'base',     'green')
call s:hi('PmenuExtra',      'overlay2', 'surface0')
call s:hi('PmenuExtraSel',   'base',     'green')
call s:hi('PmenuMatch',      'orange',   'surface0', 'bold')
call s:hi('PmenuMatchSel',   'base',     'green',    'bold')
call s:hi('PmenuSbar',       'none',     'surface1')
call s:hi('PmenuThumb',      'none',     'overlay0')
call s:hi('Question',        'green',    'none')
call s:hi('QuickFixLine',    'none',     'surface1', 'bold')
call s:hi('SpecialKey',      'overlay1', 'none')
call s:hi('StatusLine',      'text',     'surface0')
call s:hi('StatusLineNC',    'overlay1', 'mantle')
call s:hi('StatusLineTerm',  'base',     'green', 'bold')
call s:hi('StatusLineTermNC','overlay1', 'mantle')
call s:hi('TabLine',         'overlay2', 'mantle')
call s:hi('TabLineFill',     'none',     'crust')
call s:hi('TabLineSel',      'base',     'green', 'bold')
call s:hi('Title',           'green',    'none',  'bold')
call s:hi('Visual',          'none',     'surface1', 'bold')
call s:hi('VisualNOS',       'none',     'surface1')
call s:hi('WarningMsg',      'yellow',   'none')
call s:hi('Whitespace',      'surface1', 'none')
call s:hi('WildMenu',        'base',     'skye')
call s:hi('Terminal',        'text',     s:bg)

" popup windows / balloons / toolbar (GUI and +popupwin)
call s:hi('Menu',              'text',     'surface0')
call s:hi('Scrollbar',         'overlay0', 'surface0')
call s:hi('Tooltip',           'text',     'surface1')
call s:hi('ToolbarLine',       'none',     'surface0')
call s:hi('ToolbarButton',     'text',     'surface1', 'bold')
call s:hi('Popup',             'text',     'mantle')
call s:hi('PopupNotification', 'base',     'yellow')
call s:hi('PopupSelected',     'base',     'green',  'bold')
call s:hi('MessageWindow',     'text',     'surface0')

" ----------------------------------------------------------------- syntax ---
call s:hi('Comment',        'overlay1', 'none', s:italic)
call s:hi('Constant',       'orange',   'none')
call s:hi('String',         'green',    'none')
call s:hi('Character',      'aqua',     'none')
call s:hi('Number',         'orange',   'none')
call s:hi('Boolean',        'orange',   'none')
call s:hi('Float',          'orange',   'none')
call s:hi('Identifier',     'subtext1', 'none')
call s:hi('Function',       'blue',     'none')
call s:hi('Statement',      s:kw,       'none')
call s:hi('Conditional',    s:kw,       'none')
call s:hi('Repeat',         s:kw,       'none')
call s:hi('Label',          s:kw,       'none')
call s:hi('Keyword',        s:kw,       'none')
call s:hi('Operator',       'subtext0', 'none')
call s:hi('Exception',      'red',      'none')
call s:hi('PreProc',        'pink',     'none')
call s:hi('Include',        'pink',     'none')
call s:hi('Define',         'pink',     'none')
call s:hi('Macro',          'pink',     'none')
call s:hi('PreCondit',      'pink',     'none')
call s:hi('Type',           'yellow',   'none')
call s:hi('StorageClass',   'yellow',   'none')
call s:hi('Structure',      'yellow',   'none')
call s:hi('Typedef',        'yellow',   'none')
call s:hi('Special',        'purple',   'none')
call s:hi('SpecialChar',    'purple',   'none')
call s:hi('Tag',            'aqua',     'none')
call s:hi('Delimiter',      'overlay2', 'none')
call s:hi('SpecialComment', 'overlay2', 'none', s:italic)
call s:hi('Debug',          'red',      'none')
call s:hi('Underlined',     'snow',     'none', 'underline')
call s:hi('Ignore',         'overlay0', 'none')
call s:hi('Error',          'red',      'none', 'bold')
call s:hi('Todo',           'base',     'yellow', 'bold')
call s:hi('debugPC',        'none',     'surface1')
call s:hi('debugBreakpoint','base',     'red', 'bold')

" ------------------------------------------------------------------ diffs ---
call s:hi('DiffAdd',       'green',    'surface0')
call s:hi('DiffChange',    'yellow',   'surface0')
call s:hi('DiffDelete',    'red',      'surface0')
call s:hi('DiffText',      'base',     'yellow', 'bold')
call s:hi('diffAdded',     'green',    'none')
call s:hi('diffRemoved',   'red',      'none')
call s:hi('diffChanged',   'yellow',   'none')
call s:hi('diffFile',      'blue',     'none', 'bold')
call s:hi('diffLine',      'purple',   'none')
call s:hi('diffIndexLine', 'overlay2', 'none')

" ------------------------------------------------------------------ spell ---
call s:hi('SpellBad',   'none', 'none', 'undercurl', 'red')
call s:hi('SpellCap',   'none', 'none', 'undercurl', 'yellow')
call s:hi('SpellLocal', 'none', 'none', 'undercurl', 'skye')
call s:hi('SpellRare',  'none', 'none', 'undercurl', 'purple')

" ------------------------------------------------------------------ netrw ---
call s:hi('netrwDir',       'blue',     'none', 'bold')
call s:hi('netrwClassify',  'overlay1', 'none')
call s:hi('netrwExe',       'green',    'none')
call s:hi('netrwSymLink',   'skye',     'none', 'underline')
call s:hi('netrwTreeBar',   'surface2', 'none')
call s:hi('netrwMarkFile',  'base',     'orange')
call s:hi('netrwComment',   'overlay1', 'none', s:italic)
call s:hi('netrwHelpCmd',   s:kw,       'none')
call s:hi('netrwVersion',   'yellow',   'none')

" -------------------------------------------------------------- filetypes ---
call s:hi('htmlH1',           'green',  'none', 'bold')
call s:hi('htmlH2',           'aqua',   'none', 'bold')
call s:hi('markdownH1',       'green',  'none', 'bold')
call s:hi('markdownH2',       'aqua',   'none', 'bold')
call s:hi('markdownH3',       'skye',   'none', 'bold')
call s:hi('markdownCode',     'lime',   'none')
call s:hi('markdownUrl',      'snow',   'none', 'underline')
call s:hi('markdownLinkText', 'blue',   'none')
call s:hi('helpHyperTextJump','blue',   'none', 'underline')
call s:hi('helpExample',      'lime',   'none')

" ----------------------------------------------------- :terminal ansi ---
if has('terminal')
  let g:terminal_ansi_colors = [
        \ s:p.surface1[0], s:p.red[0],    s:p.green[0], s:p.yellow[0],
        \ s:p.blue[0],     s:p.purple[0], s:p.skye[0],  s:p.subtext1[0],
        \ s:p.surface2[0], s:p.orange[0], s:p.lime[0],  s:p.yellow[0],
        \ s:p.snow[0],     s:p.pink[0],   s:p.aqua[0],  s:p.text[0],
        \ ]
endif

delfunction s:hi
