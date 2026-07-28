#!/usr/bin/env bash
# Verify that Unix check mode is offline and side-effect free.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEMP_ROOT="$(mktemp -d "${TMPDIR:-/tmp}/gc-dev-setup-test.XXXXXX")"
trap 'rm -rf "$TEMP_ROOT"' EXIT

HOME_DIR="$TEMP_ROOT/home"
TOOLS_DIR="$TEMP_ROOT/tools"
mkdir -p "$HOME_DIR"
printf 'sentinel\n' >"$HOME_DIR/.profile"
printf 'sentinel\n' >"$HOME_DIR/.zprofile"

set +e
output="$(
    cd "$ROOT_DIR" &&
    HOME="$HOME_DIR" \
    GC_DEV_TOOLS_DIR="$TOOLS_DIR" \
    GC_TOKEN="must-not-appear" \
    GITCODE_TOKEN="must-not-appear" \
    GOPROXY="http://127.0.0.1:1" \
    bash scripts/dev-setup.sh --check 2>&1
)"
status=$?
set -e

# A clean CI host normally lacks the exact managed tools, so either result is
# valid; the contract under test is that check mode performs no installation.
[[ "$status" -eq 0 || "$status" -eq 1 ]]
[[ ! -e "$TOOLS_DIR" ]]
[[ "$(cat "$HOME_DIR/.profile")" == "sentinel" ]]
[[ "$(cat "$HOME_DIR/.zprofile")" == "sentinel" ]]
[[ "$output" != *"must-not-appear"* ]]

printf 'Unix dev setup check-mode test passed (exit=%s).\n' "$status"
