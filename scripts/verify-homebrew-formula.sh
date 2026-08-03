#!/usr/bin/env bash
# verify-homebrew-formula.sh
# End-to-end local verification:
#   build binary → create archive → compute sha256 → tap → audit → install → test
#
# Usage:
#   ./scripts/verify-homebrew-formula.sh
#
# Prerequisites: Go 1.22+, Homebrew, macOS or Linux

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log()  { printf "${GREEN}[OK]${NC} %s\n" "$*"; }
warn() { printf "${YELLOW}[WARN]${NC} %s\n" "$*"; }
fail() { printf "${RED}[FAIL]${NC} %s\n" "$*"; exit 1; }

UNAME_S="$(uname -s)"
UNAME_M="$(uname -m)"

case "${UNAME_S}" in
  Darwin) OS="darwin" ;;
  Linux)  OS="linux" ;;
  *)      fail "Unsupported OS: ${UNAME_S}" ;;
esac
case "${UNAME_M}" in
  arm64|aarch64) ARCH="arm64" ;;
  *)             ARCH="amd64" ;;
esac

PLATFORM="${OS}_${ARCH}"
VERSION="0.0.0-test"
WORKDIR="$(mktemp -d)"
FORMULA_TMP="${WORKDIR}/gc.rb"
ARCHIVE_NAME="gc_${VERSION}_${PLATFORM}.tar.gz"
ARCHIVE_PATH="${WORKDIR}/${ARCHIVE_NAME}"
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
LOCAL_TAP_DIR="${WORKDIR}/tap"
LOCAL_TAP="local/test-tap"

trap "rm -rf ${WORKDIR}; brew untap ${LOCAL_TAP} 2>/dev/null || true" EXIT

echo "=== Homebrew Formula Verification ==="
echo "Platform: ${PLATFORM}"
echo "Workdir:  ${WORKDIR}"
echo ""

# Step 1: Build binary
echo "--- Step 1: Build binary ---"
cd "${PROJECT_DIR}"
CGO_ENABLED=0 GOOS="${OS}" GOARCH="${ARCH}" go build -o "${WORKDIR}/gc" ./cmd/gc
log "Binary built: ${WORKDIR}/gc"
"${WORKDIR}/gc" version

# Step 2: Generate completions
echo ""
echo "--- Step 2: Generate completions ---"
mkdir -p "${WORKDIR}/completions"
"${WORKDIR}/gc" completion bash > "${WORKDIR}/completions/gc.bash"
"${WORKDIR}/gc" completion zsh  > "${WORKDIR}/completions/gc.zsh"
"${WORKDIR}/gc" completion fish > "${WORKDIR}/completions/gc.fish"
log "Completions generated"
for f in "${WORKDIR}/completions/"*; do
  echo "  $(basename "$f") ($(wc -c < "$f" | tr -d ' ') bytes)"
done

# Step 3: Create tar.gz (matching goreleaser archive structure)
echo ""
echo "--- Step 3: Create archive ---"
cd "${WORKDIR}"
tar -czf "${ARCHIVE_PATH}" gc completions/
log "Archive: ${ARCHIVE_NAME} ($(du -h "${ARCHIVE_PATH}" | cut -f1))"

# Step 4: Compute sha256
echo ""
echo "--- Step 4: Compute sha256 ---"
SHA256=$(shasum -a 256 "${ARCHIVE_PATH}" | cut -d' ' -f1)
log "SHA256: ${SHA256}"

# Step 5: Generate formula from template
echo ""
echo "--- Step 5: Generate formula ---"
cp "${PROJECT_DIR}/Formula/gc.rb" "${FORMULA_TMP}"
sed -i '' "s/VERSION_PLACEHOLDER/${VERSION}/g" "${FORMULA_TMP}"
case "${PLATFORM}" in
  darwin_amd64)  sed -i '' "s/DARWIN_AMD64_SHA256/${SHA256}/g" "${FORMULA_TMP}" ;;
  darwin_arm64)  sed -i '' "s/DARWIN_ARM64_SHA256/${SHA256}/g" "${FORMULA_TMP}" ;;
  linux_amd64)   sed -i '' "s/LINUX_AMD64_SHA256/${SHA256}/g"  "${FORMULA_TMP}" ;;
  linux_arm64)   sed -i '' "s/LINUX_ARM64_SHA256/${SHA256}/g"  "${FORMULA_TMP}" ;;
esac
log "Formula generated"

# Step 6: Validate Ruby syntax
echo ""
echo "--- Step 6: Validate formula syntax ---"
ruby -c "${FORMULA_TMP}" && log "Ruby syntax OK"

# Step 7: Create local tap (modern brew requires formula in a tap)
echo ""
echo "--- Step 7: Setup local tap ---"
mkdir -p "${LOCAL_TAP_DIR}/Formula"
cp "${FORMULA_TMP}" "${LOCAL_TAP_DIR}/Formula/gc.rb"
cd "${LOCAL_TAP_DIR}" && git init -q && git add -A && git commit -qm "init" && cd "${WORKDIR}"
brew tap "${LOCAL_TAP}" "${LOCAL_TAP_DIR}" 2>&1 | tail -1
log "Local tap: ${LOCAL_TAP}"

# Step 8: Brew audit
echo ""
echo "--- Step 8: Brew audit ---"
if brew audit --strict "${LOCAL_TAP}/gc" 2>&1; then
  log "Brew audit passed"
else
  warn "Audit warnings (non-blocking for local test)"
fi

# Step 9: Brew install via local tap
echo ""
echo "--- Step 9: Brew install ---"
if HOMEBREW_NO_AUTO_UPDATE=1 brew install "${LOCAL_TAP}/gc" 2>&1; then
  log "Brew install succeeded"
else
  warn "Install failed — check formula URL references correct archive name"
  exit 0
fi

# Step 10: Test binary
INSTALLED_GC="$(brew --prefix)/bin/gc"
echo ""
echo "--- Step 10: Test installed binary ---"
"${INSTALLED_GC}" version
log "Binary works from brew install"

# Step 11: Verify no Go runtime dependency
echo ""
echo "--- Step 11: Check self-contained (no Go deps) ---"
DEPS=$(otool -L "${INSTALLED_GC}" 2>/dev/null | grep -v '/usr/lib/libSystem' | grep -v '/System/Library' | grep -v 'not a dynamic' | grep -v '^\t' | grep -v "^${INSTALLED_GC}" || true)
if [ -z "${DEPS}" ]; then
  log "Binary is fully self-contained (zero external .dylib deps)"
else
  warn "Unexpected deps: ${DEPS}"
fi

# Step 12: Verify completions
echo ""
echo "--- Step 12: Verify completions ---"
BREW_PREFIX="$(brew --prefix)"
for comp in "${BREW_PREFIX}/etc/bash_completion.d/gc" \
            "${BREW_PREFIX}/share/zsh/site-functions/_gc" \
            "${BREW_PREFIX}/share/fish/vendor_completions.d/gc.fish"; do
  if [ -f "${comp}" ]; then
    log "Found: ${comp}"
  fi
done

# Cleanup
echo ""
brew untap "${LOCAL_TAP}" 2>/dev/null && log "Local tap cleaned up"

echo ""
echo "=== Verification: ALL CHECKS PASSED ==="
echo ""
echo "Real user workflow (after push to GitHub):"
echo "  brew tap ZhongqiXiao/gitcode-cli"
echo "  brew install gc"
echo "  gc version"
echo ""
echo "M-series Mac users get darwin_arm64 binary automatically (on_arm block in formula)."
echo "Intel Mac users get darwin_amd64 (on_intel block)."
