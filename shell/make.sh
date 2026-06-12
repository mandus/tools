#!/bin/sh
# Wrapper script to build both gitprompt and pass.
#
# Usage:
#   ./make.sh            # build both tools
#   ./make.sh build      # same
#   ./make.sh build-all  # cross-compile both tools
#   ./make.sh clean      # clean both builds
#   ./make.sh update-deps # update dependencies for both

set -eu

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
REPO_ROOT=$(dirname "$SCRIPT_DIR")

cmd=${1:-build}

case "$cmd" in
    build|build-all|clean|update-deps)
        echo "Building gitprompt..."
        (cd "$REPO_ROOT/gitprompt" && ./build.sh "$cmd")
        
        echo "Building pass..."
        (cd "$REPO_ROOT/pass" && ./build.sh "$cmd")
        ;;
    *)
        echo "Usage: $0 [build|build-all|clean|update-deps]" >&2
        exit 1
        ;;
esac
