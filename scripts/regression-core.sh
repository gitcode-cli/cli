#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GC_BIN="${GC_BIN:-$ROOT_DIR/gc}"
READONLY_REPO="${GC_REGRESSION_REPO:-infra-test/gctest1}"
RELEASE_TAG="${GC_REGRESSION_RELEASE_TAG:-v0.0.1-test}"
RUN_WRITE_PATHS="${GC_REGRESSION_WRITE:-0}"

SOURCE_TOKEN="${GC_TOKEN:-${GITCODE_TOKEN:-}}"
if [[ -z "${SOURCE_TOKEN}" ]]; then
  echo "GC_TOKEN or GITCODE_TOKEN is required" >&2
  exit 1
fi

log() {
  printf '\n[%s]\n' "$1"
}

run_capture() {
  local __var_name="$1"
  shift

  local output
  if ! output="$("$@" 2>&1)"; then
    printf '%s\n' "$output" >&2
    return 1
  fi
  printf -v "$__var_name" '%s' "$output"
  printf '%s\n' "$output"
}

assert_contains() {
  local haystack="$1"
  local needle="$2"
  if [[ "$haystack" != *"$needle"* ]]; then
    echo "expected output to contain: $needle" >&2
    exit 1
  fi
}

run_expect_fail() {
  local __var_name="$1"
  shift

  local output
  set +e
  output="$("$@" 2>&1)"
  local status=$?
  set -e

  if [[ $status -eq 0 ]]; then
    echo "expected command to fail: $*" >&2
    exit 1
  fi
  printf -v "$__var_name" '%s' "$output"
  printf '%s\n' "$output"
}

log "Build"
(cd "$ROOT_DIR" && go build -o "$GC_BIN" ./cmd/gc)

TMP_CONFIG_DIR="$(mktemp -d)"
TMP_NON_GIT_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TMP_CONFIG_DIR" "$TMP_NON_GIT_DIR"
}
trap cleanup EXIT

export GC_CONFIG_DIR="$TMP_CONFIG_DIR"
unset GC_TOKEN GITCODE_TOKEN

log "Auth Login"
run_capture login_out "$GC_BIN" auth login --token "$SOURCE_TOKEN"
assert_contains "$login_out" "Logged in as"

log "Auth Status"
run_capture status_in_out "$GC_BIN" auth status
assert_contains "$status_in_out" "Logged in as"
assert_contains "$status_in_out" "(config)"

log "Auth Token"
token_out="$("$GC_BIN" auth token 2>&1)"
printf '%s\n' "[redacted]"
if [[ -z "${token_out//[$'\n\r\t ']}" ]]; then
  echo "auth token output is empty" >&2
  exit 1
fi

log "Auth Logout"
run_capture logout_out "$GC_BIN" auth logout
assert_contains "$logout_out" "Cleared stored authentication"

log "Auth Status After Logout"
run_capture status_out_out "$GC_BIN" auth status
assert_contains "$status_out_out" "Not logged in"

export GC_TOKEN="$SOURCE_TOKEN"

log "Repo View"
run_capture repo_view_out "$GC_BIN" repo view "$READONLY_REPO"
assert_contains "$repo_view_out" "$READONLY_REPO"

log "Repo View JSON"
run_capture repo_view_json "$GC_BIN" repo view "$READONLY_REPO" --json
assert_contains "$repo_view_json" "\"full_name\""

log "Issue List"
run_capture issue_list_out "$GC_BIN" issue list -R "$READONLY_REPO" --limit 1
assert_contains "$issue_list_out" "#"
ISSUE_NUMBER="$(printf '%s\n' "$issue_list_out" | awk '/^#/ {gsub(/^#/, "", $1); print $1; exit}')"
if [[ -z "$ISSUE_NUMBER" ]]; then
  echo "failed to determine issue number from issue list output" >&2
  exit 1
fi

log "Issue View"
run_capture issue_view_out "$GC_BIN" issue view "$ISSUE_NUMBER" -R "$READONLY_REPO"
assert_contains "$issue_view_out" "#$ISSUE_NUMBER"

log "Issue List JSON"
run_capture issue_list_json "$GC_BIN" issue list -R "$READONLY_REPO" --limit 1 --json
assert_contains "$issue_list_json" "\"number\""

log "Issue View JSON"
run_capture issue_view_json "$GC_BIN" issue view "$ISSUE_NUMBER" -R "$READONLY_REPO" --json
assert_contains "$issue_view_json" "\"number\""

log "PR List JSON"
run_capture pr_list_json "$GC_BIN" pr list -R "$READONLY_REPO" --limit 1 --json
assert_contains "$pr_list_json" "["

log "Release List JSON"
run_capture release_list_json "$GC_BIN" release list -R "$READONLY_REPO" --json
assert_contains "$release_list_json" "\"tag_name\""

log "Release Delete Dry Run"
run_capture release_delete_dry_run "$GC_BIN" release delete "$RELEASE_TAG" -R "$READONLY_REPO" --dry-run
assert_contains "$release_delete_dry_run" "Dry run: would delete release"

log "Non-Git Error Path"
run_expect_fail nongit_out bash -lc "cd '$TMP_NON_GIT_DIR' && '$GC_BIN' repo view"
assert_contains "$nongit_out" "not in a git repository"

if [[ "$RUN_WRITE_PATHS" == "1" ]]; then
  PR_REPO="${GC_REGRESSION_PR_REPO:-}"
  PR_HEAD="${GC_REGRESSION_PR_HEAD:-}"
  PR_BASE="${GC_REGRESSION_PR_BASE:-main}"
  PR_TITLE="${GC_REGRESSION_PR_TITLE:-Regression Probe}"
  PR_BODY="${GC_REGRESSION_PR_BODY:-Created by scripts/regression-core.sh}"

  if [[ -z "$PR_REPO" || -z "$PR_HEAD" ]]; then
    echo "GC_REGRESSION_PR_REPO and GC_REGRESSION_PR_HEAD are required when GC_REGRESSION_WRITE=1" >&2
    exit 1
  fi

  log "PR Create"
  run_capture pr_create_out "$GC_BIN" pr create -R "$PR_REPO" --head "$PR_HEAD" --base "$PR_BASE" --title "$PR_TITLE" --body "$PR_BODY"
  assert_contains "$pr_create_out" "Created PR #"
else
  log "PR Create"
  echo "skipped: set GC_REGRESSION_WRITE=1 plus GC_REGRESSION_PR_REPO and GC_REGRESSION_PR_HEAD to run write-path regression"
fi

log "Done"
echo "Core regression checks completed."
