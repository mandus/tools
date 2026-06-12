#!/bin/sh
# Build helper for systems without GNU make (e.g. BusyBox ash on Windows arm64).
# Equivalent to the targets in ./Makefile.
#
# Usage:
#   ./build.sh            # native build → ../shell/pass[.exe]
#   ./build.sh build      # same
#   ./build.sh build-all  # cross-compile linux amd64/arm64 + windows arm64
#   ./build.sh clean
#
# Requires Go 1.20+.

set -eu

OUTDIR=../shell
LDFLAGS="-s -w"
# Resolves to ".exe" on Windows, "" elsewhere — so the native build is
# named correctly without sniffing uname.
EXT=$(go env GOEXE)

cmd=${1:-build}

case "$cmd" in
    build)
        go build -ldflags="$LDFLAGS" -o "$OUTDIR/pass$EXT" .
        ;;
    build-all)
        GOOS=linux   GOARCH=amd64 go build -ldflags="$LDFLAGS" -o "$OUTDIR/pass-linux-amd64" .
        GOOS=linux   GOARCH=arm64 go build -ldflags="$LDFLAGS" -o "$OUTDIR/pass-linux-arm64" .
        GOOS=windows GOARCH=arm64 go build -ldflags="$LDFLAGS" -o "$OUTDIR/pass-windows-arm64.exe" .
        ;;
    clean)
        rm -f "$OUTDIR/pass" \
              "$OUTDIR/pass.exe" \
              "$OUTDIR/pass-linux-amd64" \
              "$OUTDIR/pass-linux-arm64" \
              "$OUTDIR/pass-windows-arm64.exe"
        ;;
    update-deps)
        go get -u ./...
        go mod tidy
        ;;
    *)
        echo "Usage: $0 [build|build-all|update-deps|clean]" >&2
        exit 1
        ;;
esac
