#!/bin/bash
# Build script for gitcode-cli
# Usage: ./scripts/build.sh [version]

set -e

# Get version from argument or git tag
VERSION=${1:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +%Y-%m-%d)

# macOS: drop -s to preserve LC_UUID (dyld requires it)
if [[ "$(uname -s)" == "Darwin" ]]; then
  STRIP_FLAGS="-w"
else
  STRIP_FLAGS="-s -w"
fi
LDFLAGS="${STRIP_FLAGS} -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

echo "Building gc ${VERSION} (commit: ${COMMIT}, date: ${DATE})"

# Build for current platform
mkdir -p bin
go build -ldflags "${LDFLAGS}" -o bin/gc ./cmd/gc

echo "Binary built at bin/gc"
