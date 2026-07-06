#!/usr/bin/env bash

log "repo: view"
run_capture repo_view_out "$GC_BIN" repo view "$SYSTEM_REPO"
assert_contains "$repo_view_out" "$SYSTEM_REPO"

log "repo: view json"
run_capture repo_view_json "$GC_BIN" repo view "$SYSTEM_REPO" --json
printf '%s\n' "$repo_view_json" | assert_json
assert_contains "$repo_view_json" "\"full_name\""

log "repo: log json"
if probe_or_skip "repo log" "$GC_BIN" repo log -R "$SYSTEM_REPO" --limit 1 --json; then
  run_capture repo_log_json "$GC_BIN" repo log -R "$SYSTEM_REPO" --limit 1 --json
  printf '%s\n' "$repo_log_json" | assert_json
fi

log "repo: stats json"
if probe_or_skip "repo stats" "$GC_BIN" repo stats -R "$SYSTEM_REPO" --json; then
  run_capture repo_stats_json "$GC_BIN" repo stats -R "$SYSTEM_REPO" --json
  printf '%s\n' "$repo_stats_json" | assert_json
fi
