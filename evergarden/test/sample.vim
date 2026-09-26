" sample.vim -- vimscript highlighting sample
" TODO: keywords, options, registers, autocommands, functions.

set nocompatible
set expandtab shiftwidth=4 softtabstop=4
setlocal listchars=tab:>-,trail:~,nbsp:+

let s:name    = 'evergarden'
let s:count   = 42
let s:ratio   = 0.618
let s:list    = ['red', 'green', 'blue']
let s:dict    = {'fg': '#f8f9e8', 'bg': '#1e2528'}
let g:pattern = '\v^\s*(let|call)\s+'

function! s:Greet(who, ...) abort
  if a:0 > 0 && !empty(a:1)
    echomsg printf('hello %s (%d extra)', a:who, a:0)
  elseif a:who ==# 'world'
    echohl WarningMsg | echo 'generic greeting' | echohl NONE
  else
    for l:item in s:list
      echo l:item . ' => ' . get(s:dict, 'fg', 'NONE')
    endfor
  endif
  return v:true
endfunction

command! -nargs=? -bang Greet call s:Greet(<q-args>)

augroup evergarden_sample
  autocmd!
  autocmd BufWritePre *.vim  silent! %s/\s\+$//e
  autocmd ColorScheme *      echomsg 'colorscheme: ' . g:colors_name
augroup END

nnoremap <silent> <leader>g :Greet<CR>
inoremap <expr> <Tab> pumvisible() ? "\<C-n>" : "\<Tab>"

try
  call s:Greet('world', 1)
catch /^Vim\%((\a\+)\)\=:E/
  echoerr 'failed: ' . v:exception
finally
  unlet! s:count
endtry
