#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
GC_BIN="${GC_BIN:-$ROOT_DIR/gc}"
SYSTEM_REPO="${GC_SYSTEM_REPO:-infra-test/gctest1}"
WRITE_REPO="${GC_SYSTEM_WRITE_REPO:-$SYSTEM_REPO}"
RUN_WRITE=0
BUILD=1

usage() {
  cat <<'EOF'
Usage: tests/system/run.sh [--read] [--write] [--all] [--repo infra-test/name] [--write-repo infra-test/name] [--no-build]

Runs real GitCode CLI system tests. All repository targets must be infra-test/*.

Modes:
  --read       Run read-only and dry-run cases (default)
  --write      Run explicit write-path cases
  --all        Run read and write cases

Environment:
  GC_BIN                  CLI binary path (default: ./gc)
  GC_SYSTEM_REPO          read-path test repository (default: infra-test/gctest1)
  GC_SYSTEM_WRITE_REPO    write-path test repository (default: GC_SYSTEM_REPO)
  GC_SYSTEM_PR_HEAD       existing test branch for write PR create case
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --read)
      RUN_WRITE=0
      shift
      ;;
    --write)
      RUN_WRITE=1
      shift
      ;;
    --all)
      RUN_WRITE=2
      shift
      ;;
    --repo)
      SYSTEM_REPO="${2:-}"
      shift 2
      ;;
    --write-repo)
      WRITE_REPO="${2:-}"
      shift 2
      ;;
    --no-build)
      BUILD=0
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

# shellcheck source=tests/system/lib.sh
source "$ROOT_DIR/tests/system/lib.sh"

require_infra_repo "GC_SYSTEM_REPO" "$SYSTEM_REPO"
require_infra_repo "GC_SYSTEM_WRITE_REPO" "$WRITE_REPO"

export ROOT_DIR GC_BIN SYSTEM_REPO WRITE_REPO

if [[ "$BUILD" == "1" ]]; then
  log "Build"
  (cd "$ROOT_DIR" && go build -o "$GC_BIN" ./cmd/gc)
fi

run_case() {
  local case_file="$1"
  # shellcheck source=/dev/null
  (source "$case_file")
}

run_read_cases() {
  log "Read Cases"
  for case_file in "$ROOT_DIR"/tests/system/cases/read/*.sh; do
    run_case "$case_file"
  done
}

run_write_cases() {
  log "Write Cases"
  for case_file in "$ROOT_DIR"/tests/system/cases/write/*.sh; do
    run_case "$case_file"
  done
}

case "$RUN_WRITE" in
  0)
    run_read_cases
    ;;
  1)
    run_write_cases
    ;;
  2)
    run_read_cases
    run_write_cases
    ;;
esac

log "Done"
echo "System tests completed."
