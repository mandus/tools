#!/bin/sh
# install.sh -- install / uninstall the evergarden colorscheme
#
#   ./install.sh [install|uninstall|status]
#
# Copies colors/evergarden.vim into every Vim runtime directory found on this
# machine (~/vimfiles and ~/.vim). Copies (not symlinks) so it works on Windows.
# POSIX sh only: runs under dash, ash/busybox, ksh, bash.

set -eu

DIR=$(CDPATH= cd "$(dirname "$0")" && pwd)
SRC="$DIR/colors/evergarden.vim"
NAME="evergarden.vim"

die() { printf 'error: %s\n' "$1" >&2; exit 1; }

# Newline-separated list of target directories.
targets() {
  found=''
  if [ -d "$HOME/vimfiles" ]; then
    found="$found$HOME/vimfiles/colors
"
  fi
  if [ -d "$HOME/.vim" ]; then
    found="$found$HOME/.vim/colors
"
  fi
  # Nothing detected: fall back to the conventional Vim location.
  if [ -z "$found" ]; then
    found="$HOME/.vim/colors
"
  fi
  printf '%s' "$found"
}

cmd_install() {
  [ -f "$SRC" ] || die "missing $SRC"
  targets | while IFS= read -r dir; do
    [ -n "$dir" ] || continue
    mkdir -p "$dir"
    cp -f "$SRC" "$dir/$NAME"
    printf 'installed  %s\n' "$dir/$NAME"
  done
  printf '\nAdd to your vimrc:\n  colorscheme evergarden\n'
}

cmd_uninstall() {
  targets | while IFS= read -r dir; do
    [ -n "$dir" ] || continue
    if [ -f "$dir/$NAME" ]; then
      rm -f "$dir/$NAME"
      printf 'removed    %s\n' "$dir/$NAME"
    else
      printf 'absent     %s\n' "$dir/$NAME"
    fi
  done
}

cmd_status() {
  targets | while IFS= read -r dir; do
    [ -n "$dir" ] || continue
    if [ -f "$dir/$NAME" ]; then
      printf 'installed  %s\n' "$dir/$NAME"
    else
      printf 'missing    %s\n' "$dir/$NAME"
    fi
  done
}

case "${1:-install}" in
  install)   cmd_install ;;
  uninstall) cmd_uninstall ;;
  status)    cmd_status ;;
  -h|--help|help)
    sed -n '2,6p' "$0" | sed 's/^# \{0,1\}//'
    ;;
  *) die "unknown command: $1 (expected install|uninstall|status)" ;;
esac
