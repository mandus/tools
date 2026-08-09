#!/bin/sh
# Wrapper script to build all tools: gitprompt, pass and sleep.
#
# Usage:
#   ./make.sh            # build all tools
#   ./make.sh build      # same
#   ./make.sh build-all  # cross-compile (sleep is native-only)
#   ./make.sh clean      # clean all builds
#   ./make.sh update-deps # update dependencies where applicable

set -eu

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
REPO_ROOT=$(dirname "$SCRIPT_DIR")
PREFIX=${PREFIX:-$HOME/.local}
EXT=$(go env GOEXE)


cmd=${1:-build}

case "$cmd" in
    build|build-all|clean|update-deps)
        echo "Building gitprompt..."
        (cd "$REPO_ROOT/gitprompt" && ./build.sh "$cmd")
        
        echo "Building pass..."
        (cd "$REPO_ROOT/pass" && ./build.sh "$cmd")

        echo "Building sleep..."
        (cd "$REPO_ROOT/sleep" && ./build.sh "$cmd")
        ;;
	install)
		echo "Install to $PREFIX"
		cp $REPO_ROOT/shell/pass$EXT $PREFIX/bin
		cp $REPO_ROOT/shell/gitprompt$EXT $PREFIX/bin
		cp $REPO_ROOT/shell/sleep$EXT $PREFIX/bin
		;;
    *)
        echo "Usage: $0 [build|build-all|clean|update-deps]" >&2
        exit 1
        ;;
esac
