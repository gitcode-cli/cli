#!/bin/bash
# Package script for gitcode-cli
# Usage: ./scripts/package.sh [version]
#
# This script builds binaries and creates DEB/RPM packages for all platforms.
# Run from project root directory.

set -e

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Change to project root
cd "${PROJECT_ROOT}"

# Get version from argument or nfpm config
VERSION=${1:-$(grep 'version:' nfpm-amd64.yaml | head -1 | sed 's/version: "//' | sed 's/"//')}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +%Y-%m-%d)

LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

echo "============================================"
echo "Packaging gc ${VERSION}"
echo "Commit: ${COMMIT}, Date: ${DATE}"
echo "============================================"

# Create dist directory
mkdir -p dist

# Build binaries
echo ""
echo "[1/4] Building Linux amd64 binary..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "${LDFLAGS}" -o dist/gc_linux_amd64 ./cmd/gc
echo "  ✓ dist/gc_linux_amd64"

echo ""
echo "[2/4] Building Linux arm64 binary..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "${LDFLAGS}" -o dist/gc_linux_arm64 ./cmd/gc
echo "  ✓ dist/gc_linux_arm64"

# Check nfpm availability
NFPMPATH=$(command -v nfpm 2>/dev/null || echo "$HOME/go/bin/nfpm")
if [ ! -x "$NFPMPATH" ]; then
    echo "Error: nfpm not found. Install it with:"
    echo "  go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest"
    exit 1
fi

# Create packages
echo ""
echo "[3/4] Creating amd64 packages..."
cd dist
"$NFPMPATH" pkg -f ../nfpm-amd64.yaml -p deb -t "gc_${VERSION}_amd64.deb" 2>&1 | sed 's/^/  /'
"$NFPMPATH" pkg -f ../nfpm-amd64.yaml -p rpm -t "gc-${VERSION}-1.x86_64.rpm" 2>&1 | sed 's/^/  /'
cd "${PROJECT_ROOT}"
echo "  ✓ gc_${VERSION}_amd64.deb"
echo "  ✓ gc-${VERSION}-1.x86_64.rpm"

echo ""
echo "[4/4] Creating arm64 packages..."
cd dist
"$NFPMPATH" pkg -f ../nfpm-arm64.yaml -p deb -t "gc_${VERSION}_arm64.deb" 2>&1 | sed 's/^/  /'
"$NFPMPATH" pkg -f ../nfpm-arm64.yaml -p rpm -t "gc-${VERSION}-1.aarch64.rpm" 2>&1 | sed 's/^/  /'
cd "${PROJECT_ROOT}"
echo "  ✓ gc_${VERSION}_arm64.deb"
echo "  ✓ gc-${VERSION}-1.aarch64.rpm"

# Summary
echo ""
echo "============================================"
echo "Package Summary"
echo "============================================"
ls -lh dist/gc*${VERSION}* 2>/dev/null | awk '{print "  " $9 ": " $5}'

echo ""
echo "Done! Packages are in dist/"
echo ""
echo "Next steps:"
echo "  1. Create release:"
echo "     ./gc release create ${VERSION} -R gitcode-cli/cli --title \"${VERSION}\" --notes \"Release notes\""
echo ""
echo "  2. Upload packages:"
echo "     ./gc release upload ${VERSION} dist/gc_${VERSION}_amd64.deb dist/gc_${VERSION}_arm64.deb dist/gc-${VERSION}-1.x86_64.rpm dist/gc-${VERSION}-1.aarch64.rpm -R gitcode-cli/cli"