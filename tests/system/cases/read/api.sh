#!/usr/bin/env bash

log "api: repo"
run_capture api_repo_json "$GC_BIN" api "repos/$SYSTEM_REPO"
printf '%s\n' "$api_repo_json" | assert_json
assert_contains "$api_repo_json" "\"full_name\""

log "api: repo commits"
if probe_or_skip "api commits" "$GC_BIN" api "repos/$SYSTEM_REPO/commits?path=README.md&sha=main"; then
  run_capture api_commits_json "$GC_BIN" api "repos/$SYSTEM_REPO/commits?path=README.md&sha=main"
  printf '%s\n' "$api_commits_json" | assert_json
fi
