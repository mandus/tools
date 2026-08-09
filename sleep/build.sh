#!/bin/sh
# Build helper for the sleep tool (C, no Go toolchain required).
# Equivalent to the targets in ./Makefile.
#
# Usage:
#   ./build.sh            # native build → ../shell/sleep[.exe]
#   ./build.sh build      # same
#   ./build.sh clean
#
# Requires a C compiler (cc/gcc/clang). Override with CC=...

set -eu

OUTDIR=../shell
CC=${CC:-gcc}
CFLAGS=${CFLAGS:--O2 -Wall -Wextra}
LDFLAGS=${LDFLAGS:--s}

# Resolves to ".exe" on Windows, "" elsewhere.
case "$(uname -s)" in
    MINGW*|MSYS*|CYGWIN*|*NT*) EXT=.exe ;;
    *)                         EXT=    ;;
esac

cmd=${1:-build}

case "$cmd" in
    build)
        mkdir -p "$OUTDIR"
        $CC $CFLAGS $LDFLAGS -o "$OUTDIR/sleep$EXT" sleep.c
        ;;
    clean)
        rm -f "$OUTDIR/sleep" "$OUTDIR/sleep.exe"
        ;;
    update-deps)
        # No dependencies; present so ./shell/make.sh can pass the target through.
        ;;
    build-all)
        # Cross-compilation is not wired up for the C tool; native only.
        "$0" build
        ;;
    *)
        echo "Usage: $0 [build|build-all|update-deps|clean]" >&2
        exit 1
        ;;
esac
