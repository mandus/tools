#!/bin/sh
# Build helper for systems without GNU make (e.g. BusyBox ash on Windows arm64).
# Equivalent to the targets in ./Makefile.
#
# Usage:
#   ./build.sh            # native build → ../shell/gitprompt[.exe]
#   ./build.sh build      # same
#   ./build.sh build-all  # cross-compile linux amd64/arm64 + windows arm64
#   ./build.sh clean
#
# Requires Go 1.25+ (go-git v5.19 dependency).

set -eu

OUTDIR=../shell
LDFLAGS="-s -w"
# Resolves to ".exe" on Windows, "" elsewhere — so the native build is
# named correctly without sniffing uname.
EXT=$(go env GOEXE)

cmd=${1:-build}

case "$cmd" in
    build)
        go build -ldflags="$LDFLAGS" -o "$OUTDIR/gitprompt$EXT" .
        ;;
    build-all)
        GOOS=linux   GOARCH=amd64 go build -ldflags="$LDFLAGS" -o "$OUTDIR/gitprompt-linux-amd64" .
        GOOS=linux   GOARCH=arm64 go build -ldflags="$LDFLAGS" -o "$OUTDIR/gitprompt-linux-arm64" .
        GOOS=windows GOARCH=arm64 go build -ldflags="$LDFLAGS" -o "$OUTDIR/gitprompt-windows-arm64.exe" .
        ;;
    clean)
        rm -f "$OUTDIR/gitprompt" \
              "$OUTDIR/gitprompt.exe" \
              "$OUTDIR/gitprompt-linux-amd64" \
              "$OUTDIR/gitprompt-linux-arm64" \
              "$OUTDIR/gitprompt-windows-arm64.exe"
        ;;
    *)
        echo "Usage: $0 [build|build-all|clean]" >&2
        exit 1
        ;;
esac
