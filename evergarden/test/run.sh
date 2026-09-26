#!/bin/sh
# run.sh -- preview evergarden without installing it
#
#   ./test/run.sh          # true colour (termguicolors)
#   ./test/run.sh 256      # 256-colour cterm fallback
#
# All Vim setup lives in test/vimrc (Vim allows at most 10 "-c" arguments).
# POSIX sh only: runs under dash, ash/busybox, ksh, bash.
#
# Tab 1: syntax/hitest.vim  -- every highlight group in its own colours
# Tab 2: sample.c           -- syntax in context
# Tab 3: sample.vim
# Tab 4: sample.md          -- spell on, to show SpellBad
# Tab 5: live diff          -- DiffAdd / DiffChange / DiffDelete / DiffText

set -eu

DIR=$(CDPATH= cd "$(dirname "$0")" && pwd)

if [ "${1:-}" = "256" ]; then
  EVERGARDEN_TEST_256=1
  export EVERGARDEN_TEST_256
fi

exec vim -N -u "$DIR/vimrc"
