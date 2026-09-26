" Name:         evergarden
" Description:  A forest-toned dark colorscheme inspired by The Minish Cap. Evergarden.moe used as basis
" Author:       AI / mandus
" Maintainer:   mandus 
" URL:          https://evergarden.moe
" License:      Vim License (see `:help license`)
" Last Change:  2026.09.26
"
" Spec: specs/010-evergarden-colorscheme/spec.md
"
" Options:
"   let g:evergarden_italic      = 0       " disable italics (default 1)
"   let g:evergarden_transparent = 1       " keep terminal background (default 0)
"   let g:evergarden_keyword     = 'blue'  " keyword hue (default 'skye'); any
"                                          " palette name: skye, aqua, blue,
"                                          " green, lime, orange, purple, pink...

set background=dark

hi clear
let g:colors_name = 'evergarden'

let s:t_Co = has('gui_running') ? 16777216 : str2nr(&t_Co)
let s:tgc = has('termguicolors') && &termguicolors

let g:terminal_ansi_colors = ['#374145', '#f57f82', '#cbe3b3', '#f5d098', '#b2caed', '#d2bdf3', '#b3e6db', '#adc9bc', '#4a585c', '#f7a182', '#dbe6af', '#f5d098', '#afd9e6', '#f3c0e5', '#b3e3ca', '#f8f9e8']

hi! link Boolean Constant
hi! link Conditional Statement
hi! link CursorIM Cursor
hi! link CursorLineFold FoldColumn
hi! link CursorLineSign SignColumn
hi! link Define PreProc
hi! link Float Constant
hi! link Include PreProc
hi! link Keyword Statement
hi! link Label Statement
hi! link LineNrAbove LineNr
hi! link LineNrBelow LineNr
hi! link Macro PreProc
hi! link Number Constant
hi! link PopupSelected PmenuSel
hi! link PreCondit PreProc
hi! link PreInsert NonText
hi! link Repeat Statement
hi! link SpecialChar Special
hi! link StatusLineTermNC StatusLineNC
hi! link StorageClass Type
hi! link Structure Type
hi! link TabPanel TabLine
hi! link TabPanelFill TabLineFill
hi! link Terminal Normal
hi! link Typedef Type
hi! link VertSplitNC VertSplit
hi! link diffAdded Added
hi! link diffChanged Changed
hi! link diffRemoved Removed
hi! link lCursor Cursor
hi! link markdownH1 htmlH1
hi! link markdownH2 htmlH2
hi! link markdownUrl Underlined
hi! link netrwComment Comment
hi! link netrwHelpCmd Statement

hi Normal guifg=#f8f9e8 guibg=#1e2528 guisp=NONE gui=NONE ctermfg=255 ctermbg=235 cterm=NONE term=NONE
hi Added guifg=#cbe3b3 guibg=NONE guisp=NONE gui=NONE ctermfg=151 ctermbg=NONE cterm=NONE term=NONE
hi Changed guifg=#f5d098 guibg=NONE guisp=NONE gui=NONE ctermfg=222 ctermbg=NONE cterm=NONE term=NONE
hi Character guifg=#b3e3ca guibg=NONE guisp=NONE gui=NONE ctermfg=158 ctermbg=NONE cterm=NONE term=NONE
hi ColorColumn guifg=NONE guibg=#262f33 guisp=NONE gui=NONE ctermfg=NONE ctermbg=236 cterm=NONE term=reverse
hi Comment guifg=#6f8788 guibg=NONE guisp=NONE gui=italic ctermfg=66 ctermbg=NONE cterm=italic term=bold
hi Conceal guifg=#6f8788 guibg=NONE guisp=NONE gui=NONE ctermfg=66 ctermbg=NONE cterm=NONE term=NONE
hi Constant guifg=#f7a182 guibg=NONE guisp=NONE gui=NONE ctermfg=216 ctermbg=NONE cterm=NONE term=NONE
hi CurSearch guifg=#1e2528 guibg=#f7a182 guisp=NONE gui=NONE ctermfg=235 ctermbg=216 cterm=NONE term=reverse
hi Cursor guifg=#1e2528 guibg=#f8f9e8 guisp=NONE gui=NONE ctermfg=235 ctermbg=255 cterm=NONE term=reverse
hi CursorColumn guifg=NONE guibg=#262f33 guisp=NONE gui=NONE ctermfg=NONE ctermbg=236 cterm=NONE term=NONE
hi CursorLine guifg=NONE guibg=#262f33 guisp=NONE gui=NONE ctermfg=NONE ctermbg=236 cterm=NONE term=underline
hi CursorLineNr guifg=#f5d098 guibg=NONE guisp=NONE gui=bold ctermfg=222 ctermbg=NONE cterm=bold term=bold
hi Debug guifg=#f57f82 guibg=NONE guisp=NONE gui=NONE ctermfg=210 ctermbg=NONE cterm=NONE term=NONE
hi Delimiter guifg=#839e9a guibg=NONE guisp=NONE gui=NONE ctermfg=246 ctermbg=NONE cterm=NONE term=NONE
hi DiffAdd guifg=#cbe3b3 guibg=#262f33 guisp=NONE gui=NONE ctermfg=151 ctermbg=236 cterm=NONE term=reverse
hi DiffChange guifg=#f5d098 guibg=#262f33 guisp=NONE gui=NONE ctermfg=222 ctermbg=236 cterm=NONE term=NONE
hi DiffDelete guifg=#f57f82 guibg=#262f33 guisp=NONE gui=NONE ctermfg=210 ctermbg=236 cterm=NONE term=reverse
hi DiffText guifg=#1e2528 guibg=#f5d098 guisp=NONE gui=bold ctermfg=235 ctermbg=222 cterm=bold term=reverse
hi Directory guifg=#b2caed guibg=NONE guisp=NONE gui=NONE ctermfg=111 ctermbg=NONE cterm=NONE term=NONE
hi EndOfBuffer guifg=#374145 guibg=#1e2528 guisp=NONE gui=NONE ctermfg=238 ctermbg=235 cterm=NONE term=NONE
hi Error guifg=#f57f82 guibg=NONE guisp=NONE gui=bold ctermfg=210 ctermbg=NONE cterm=bold term=bold,reverse
hi ErrorMsg guifg=#f57f82 guibg=NONE guisp=NONE gui=bold ctermfg=210 ctermbg=NONE cterm=bold term=bold,reverse
hi Exception guifg=#f57f82 guibg=NONE guisp=NONE gui=NONE ctermfg=210 ctermbg=NONE cterm=NONE term=NONE
hi FoldColumn guifg=#58686d guibg=#1e2528 guisp=NONE gui=NONE ctermfg=241 ctermbg=235 cterm=NONE term=NONE
hi Folded guifg=#96b4aa guibg=#262f33 guisp=NONE gui=NONE ctermfg=109 ctermbg=236 cterm=NONE term=NONE
hi Function guifg=#b2caed guibg=NONE guisp=NONE gui=NONE ctermfg=111 ctermbg=NONE cterm=NONE term=NONE
hi Identifier guifg=#adc9bc guibg=NONE guisp=NONE gui=NONE ctermfg=145 ctermbg=NONE cterm=NONE term=NONE
hi Ignore guifg=#58686d guibg=NONE guisp=NONE gui=NONE ctermfg=241 ctermbg=NONE cterm=NONE term=NONE
hi IncSearch guifg=#1e2528 guibg=#f7a182 guisp=NONE gui=NONE ctermfg=235 ctermbg=216 cterm=NONE term=bold,reverse,underline
hi LineNr guifg=#4a585c guibg=#1e2528 guisp=NONE gui=NONE ctermfg=240 ctermbg=235 cterm=NONE term=NONE
hi MatchParen guifg=#f7a182 guibg=#4a585c guisp=NONE gui=bold ctermfg=216 ctermbg=240 cterm=bold term=bold,underline
hi Menu guifg=#f8f9e8 guibg=#262f33 guisp=NONE gui=NONE ctermfg=255 ctermbg=236 cterm=NONE term=NONE
hi MessageWindow guifg=#f8f9e8 guibg=#262f33 guisp=NONE gui=NONE ctermfg=255 ctermbg=236 cterm=NONE term=NONE
hi ModeMsg guifg=#f8f9e8 guibg=NONE guisp=NONE gui=bold ctermfg=255 ctermbg=NONE cterm=bold term=bold
hi MoreMsg guifg=#cbe3b3 guibg=NONE guisp=NONE gui=NONE ctermfg=151 ctermbg=NONE cterm=NONE term=NONE
hi NonText guifg=#4a585c guibg=NONE guisp=NONE gui=NONE ctermfg=240 ctermbg=NONE cterm=NONE term=NONE
hi Operator guifg=#96b4aa guibg=NONE guisp=NONE gui=NONE ctermfg=109 ctermbg=NONE cterm=NONE term=NONE
hi Pmenu guifg=#adc9bc guibg=#262f33 guisp=NONE gui=NONE ctermfg=145 ctermbg=236 cterm=NONE term=reverse
hi PmenuBorder guifg=#6f8788 guibg=#262f33 guisp=NONE gui=NONE ctermfg=66 ctermbg=236 cterm=NONE term=NONE
hi PmenuExtra guifg=#839e9a guibg=#262f33 guisp=NONE gui=NONE ctermfg=246 ctermbg=236 cterm=NONE term=reverse
hi PmenuExtraSel guifg=#1e2528 guibg=#cbe3b3 guisp=NONE gui=NONE ctermfg=235 ctermbg=151 cterm=NONE term=NONE
hi PmenuKind guifg=#d2bdf3 guibg=#262f33 guisp=NONE gui=NONE ctermfg=183 ctermbg=236 cterm=NONE term=reverse
hi PmenuKindSel guifg=#1e2528 guibg=#cbe3b3 guisp=NONE gui=NONE ctermfg=235 ctermbg=151 cterm=NONE term=NONE
hi PmenuMatch guifg=#f7a182 guibg=#262f33 guisp=NONE gui=bold ctermfg=216 ctermbg=236 cterm=bold term=bold
hi PmenuMatchSel guifg=#1e2528 guibg=#cbe3b3 guisp=NONE gui=bold ctermfg=235 ctermbg=151 cterm=bold term=bold
hi PmenuSbar guifg=NONE guibg=#374145 guisp=NONE gui=NONE ctermfg=NONE ctermbg=238 cterm=NONE term=reverse
hi PmenuSel guifg=#1e2528 guibg=#cbe3b3 guisp=NONE gui=bold ctermfg=235 ctermbg=151 cterm=bold term=bold
hi PmenuShadow guifg=#6f8788 guibg=#171c1f guisp=NONE gui=NONE ctermfg=66 ctermbg=233 cterm=NONE term=NONE
hi PmenuThumb guifg=NONE guibg=#58686d guisp=NONE gui=NONE ctermfg=NONE ctermbg=241 cterm=NONE term=NONE
hi Popup guifg=#f8f9e8 guibg=#191e21 guisp=NONE gui=NONE ctermfg=255 ctermbg=234 cterm=NONE term=NONE
hi PopupBorder guifg=#6f8788 guibg=#191e21 guisp=NONE gui=NONE ctermfg=66 ctermbg=234 cterm=NONE term=NONE
hi PopupNotification guifg=#1e2528 guibg=#f5d098 guisp=NONE gui=NONE ctermfg=235 ctermbg=222 cterm=NONE term=NONE
hi PopupTitle guifg=#f8f9e8 guibg=#191e21 guisp=NONE gui=bold ctermfg=255 ctermbg=234 cterm=bold term=bold
hi PreProc guifg=#f3c0e5 guibg=NONE guisp=NONE gui=NONE ctermfg=218 ctermbg=NONE cterm=NONE term=NONE
hi Question guifg=#cbe3b3 guibg=NONE guisp=NONE gui=NONE ctermfg=151 ctermbg=NONE cterm=NONE term=standout
hi QuickFixLine guifg=NONE guibg=#374145 guisp=NONE gui=bold ctermfg=NONE ctermbg=238 cterm=bold term=bold
hi Removed guifg=#f57f82 guibg=NONE guisp=NONE gui=NONE ctermfg=210 ctermbg=NONE cterm=NONE term=NONE
hi Scrollbar guifg=#58686d guibg=#262f33 guisp=NONE gui=NONE ctermfg=241 ctermbg=236 cterm=NONE term=NONE
hi Search guifg=#1e2528 guibg=#f5d098 guisp=NONE gui=NONE ctermfg=235 ctermbg=222 cterm=NONE term=reverse
hi SignColumn guifg=#58686d guibg=#1e2528 guisp=NONE gui=NONE ctermfg=241 ctermbg=235 cterm=NONE term=reverse
hi Special guifg=#d2bdf3 guibg=NONE guisp=NONE gui=NONE ctermfg=183 ctermbg=NONE cterm=NONE term=NONE
hi SpecialComment guifg=#839e9a guibg=NONE guisp=NONE gui=italic ctermfg=246 ctermbg=NONE cterm=italic term=bold
hi SpecialKey guifg=#6f8788 guibg=NONE guisp=NONE gui=NONE ctermfg=66 ctermbg=NONE cterm=NONE term=bold
hi SpellBad guifg=NONE guibg=NONE guisp=#f57f82 gui=undercurl ctermfg=210 ctermbg=NONE cterm=underline term=underline
hi SpellCap guifg=NONE guibg=NONE guisp=#f5d098 gui=undercurl ctermfg=222 ctermbg=NONE cterm=underline term=underline
hi SpellLocal guifg=NONE guibg=NONE guisp=#b3e6db gui=undercurl ctermfg=152 ctermbg=NONE cterm=underline term=underline
hi SpellRare guifg=NONE guibg=NONE guisp=#d2bdf3 gui=undercurl ctermfg=183 ctermbg=NONE cterm=underline term=underline
hi Statement guifg=#b3e6db guibg=NONE guisp=NONE gui=NONE ctermfg=152 ctermbg=NONE cterm=NONE term=NONE
hi StatusLine guifg=#f8f9e8 guibg=#262f33 guisp=NONE gui=NONE ctermfg=255 ctermbg=236 cterm=NONE term=bold,reverse
hi StatusLineNC guifg=#6f8788 guibg=#191e21 guisp=NONE gui=NONE ctermfg=66 ctermbg=234 cterm=NONE term=bold,underline
hi StatusLineTerm guifg=#1e2528 guibg=#cbe3b3 guisp=NONE gui=bold ctermfg=235 ctermbg=151 cterm=bold term=bold,reverse
hi String guifg=#cbe3b3 guibg=NONE guisp=NONE gui=NONE ctermfg=151 ctermbg=NONE cterm=NONE term=NONE
hi TabLine guifg=#839e9a guibg=#191e21 guisp=NONE gui=NONE ctermfg=246 ctermbg=234 cterm=NONE term=bold,underline
hi TabLineFill guifg=NONE guibg=#171c1f guisp=NONE gui=NONE ctermfg=NONE ctermbg=233 cterm=NONE term=NONE
hi TabLineSel guifg=#1e2528 guibg=#cbe3b3 guisp=NONE gui=bold ctermfg=235 ctermbg=151 cterm=bold term=bold,reverse
hi Tag guifg=#b3e3ca guibg=NONE guisp=NONE gui=NONE ctermfg=158 ctermbg=NONE cterm=NONE term=NONE
hi Title guifg=#cbe3b3 guibg=NONE guisp=NONE gui=bold ctermfg=151 ctermbg=NONE cterm=bold term=bold
hi TitleBar guifg=#f8f9e8 guibg=#374145 guisp=NONE gui=NONE ctermfg=255 ctermbg=238 cterm=NONE term=NONE
hi TitleBarNC guifg=#6f8788 guibg=#262f33 guisp=NONE gui=NONE ctermfg=66 ctermbg=236 cterm=NONE term=NONE
hi Todo guifg=#1e2528 guibg=#f5d098 guisp=NONE gui=bold ctermfg=235 ctermbg=222 cterm=bold term=bold,reverse
hi ToolbarButton guifg=#f8f9e8 guibg=#374145 guisp=NONE gui=bold ctermfg=255 ctermbg=238 cterm=bold term=bold,reverse
hi ToolbarLine guifg=NONE guibg=#262f33 guisp=NONE gui=NONE ctermfg=NONE ctermbg=236 cterm=NONE term=reverse
hi Tooltip guifg=#f8f9e8 guibg=#374145 guisp=NONE gui=NONE ctermfg=255 ctermbg=238 cterm=NONE term=NONE
hi Type guifg=#f5d098 guibg=NONE guisp=NONE gui=NONE ctermfg=222 ctermbg=NONE cterm=NONE term=NONE
hi Underlined guifg=#afd9e6 guibg=NONE guisp=NONE gui=underline ctermfg=153 ctermbg=NONE cterm=underline term=underline
hi VertSplit guifg=#374145 guibg=#1e2528 guisp=NONE gui=NONE ctermfg=238 ctermbg=235 cterm=NONE term=NONE
hi Visual guifg=NONE guibg=#374145 guisp=NONE gui=bold ctermfg=NONE ctermbg=238 cterm=bold term=reverse
hi VisualNOS guifg=NONE guibg=#374145 guisp=NONE gui=NONE ctermfg=NONE ctermbg=238 cterm=NONE term=reverse
hi WarningMsg guifg=#f5d098 guibg=NONE guisp=NONE gui=NONE ctermfg=222 ctermbg=NONE cterm=NONE term=standout
hi Whitespace guifg=#374145 guibg=NONE guisp=NONE gui=NONE ctermfg=238 ctermbg=NONE cterm=NONE term=NONE
hi WildMenu guifg=#1e2528 guibg=#b3e6db guisp=NONE gui=NONE ctermfg=235 ctermbg=152 cterm=NONE term=bold
hi debugBreakpoint guifg=#1e2528 guibg=#f57f82 guisp=NONE gui=bold ctermfg=235 ctermbg=210 cterm=bold term=bold
hi debugPC guifg=NONE guibg=#374145 guisp=NONE gui=NONE ctermfg=NONE ctermbg=238 cterm=NONE term=reverse
hi diffFile guifg=#b2caed guibg=NONE guisp=NONE gui=bold ctermfg=111 ctermbg=NONE cterm=bold term=bold
hi diffIndexLine guifg=#839e9a guibg=NONE guisp=NONE gui=NONE ctermfg=246 ctermbg=NONE cterm=NONE term=NONE
hi diffLine guifg=#d2bdf3 guibg=NONE guisp=NONE gui=NONE ctermfg=183 ctermbg=NONE cterm=NONE term=NONE
hi helpExample guifg=#dbe6af guibg=NONE guisp=NONE gui=NONE ctermfg=187 ctermbg=NONE cterm=NONE term=NONE
hi helpHyperTextJump guifg=#b2caed guibg=NONE guisp=NONE gui=underline ctermfg=111 ctermbg=NONE cterm=underline term=underline
hi htmlH1 guifg=#cbe3b3 guibg=NONE guisp=NONE gui=bold ctermfg=151 ctermbg=NONE cterm=bold term=bold
hi htmlH2 guifg=#b3e3ca guibg=NONE guisp=NONE gui=bold ctermfg=158 ctermbg=NONE cterm=bold term=bold
hi markdownCode guifg=#dbe6af guibg=NONE guisp=NONE gui=NONE ctermfg=187 ctermbg=NONE cterm=NONE term=NONE
hi markdownH3 guifg=#b3e6db guibg=NONE guisp=NONE gui=bold ctermfg=152 ctermbg=NONE cterm=bold term=bold
hi markdownLinkText guifg=#b2caed guibg=NONE guisp=NONE gui=NONE ctermfg=111 ctermbg=NONE cterm=NONE term=NONE
hi netrwClassify guifg=#6f8788 guibg=NONE guisp=NONE gui=NONE ctermfg=66 ctermbg=NONE cterm=NONE term=NONE
hi netrwDir guifg=#b2caed guibg=NONE guisp=NONE gui=bold ctermfg=111 ctermbg=NONE cterm=bold term=bold
hi netrwExe guifg=#cbe3b3 guibg=NONE guisp=NONE gui=NONE ctermfg=151 ctermbg=NONE cterm=NONE term=NONE
hi netrwMarkFile guifg=#1e2528 guibg=#f7a182 guisp=NONE gui=NONE ctermfg=235 ctermbg=216 cterm=NONE term=reverse
hi netrwSymLink guifg=#b3e6db guibg=NONE guisp=NONE gui=underline ctermfg=152 ctermbg=NONE cterm=underline term=underline
hi netrwTreeBar guifg=#4a585c guibg=NONE guisp=NONE gui=NONE ctermfg=240 ctermbg=NONE cterm=NONE term=NONE
hi netrwVersion guifg=#f5d098 guibg=NONE guisp=NONE gui=NONE ctermfg=222 ctermbg=NONE cterm=NONE term=NONE

" Options --------------------------------------------------------------- {{{
if !get(g:, 'evergarden_italic', 1)
  hi Comment gui=NONE cterm=NONE
  hi SpecialComment gui=NONE cterm=NONE
endif

if get(g:, 'evergarden_transparent', 0)
  hi Normal guibg=NONE ctermbg=NONE
  hi EndOfBuffer guibg=NONE ctermbg=NONE
  hi FoldColumn guibg=NONE ctermbg=NONE
  hi LineNr guibg=NONE ctermbg=NONE
  hi SignColumn guibg=NONE ctermbg=NONE
  hi VertSplit guibg=NONE ctermbg=NONE
endif

if exists('g:evergarden_keyword')
  let s:hues = {'red': ['#f57f82', 210], 'orange': ['#f7a182', 216],
        \ 'yellow': ['#f5d098', 222], 'lime': ['#dbe6af', 187],
        \ 'green': ['#cbe3b3', 151], 'aqua': ['#b3e3ca', 158],
        \ 'skye': ['#b3e6db', 152], 'snow': ['#afd9e6', 153],
        \ 'blue': ['#b2caed', 111], 'purple': ['#d2bdf3', 183],
        \ 'pink': ['#f3c0e5', 218], 'cherry': ['#fae6ef', 255],
        \ 'text': ['#f8f9e8', 255], 'subtext1': ['#adc9bc', 145],
        \ 'subtext0': ['#96b4aa', 109]}
  if has_key(s:hues, g:evergarden_keyword)
    let s:kw = s:hues[g:evergarden_keyword]
    execute 'hi Statement guifg=' . s:kw[0] . ' ctermfg=' . s:kw[1]
    unlet s:kw
  endif
  unlet s:hues
endif
" }}}

if s:tgc || s:t_Co >= 256
  finish
endif

" 16- and 8-colour terminals are out of scope (see the spec):
" the 256-colour attributes above are left in place.
if s:t_Co >= 8
  finish
endif

" Monochrome terminals: groups that are defined only by a link
" still need their own term attributes.
if s:t_Co >= 0
  hi CursorLineFold term=underline
  hi CursorLineSign term=underline
  hi StatusLineTermNC term=bold,underline
  hi Terminal term=NONE
  finish
endif

" vim: et ts=8 sw=2 sts=2
