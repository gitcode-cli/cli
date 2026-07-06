#!/usr/bin/env bash

log "issue: fixture probe"
if ! probe_or_skip "issue list" "$GC_BIN" issue list -R "$SYSTEM_REPO" --limit 1 --json; then
  return 0
fi

log "issue: list"
run_capture issue_list_out "$GC_BIN" issue list -R "$SYSTEM_REPO" --limit 1
assert_contains "$issue_list_out" "#"
ISSUE_NUMBER="$(printf '%s\n' "$issue_list_out" | awk '/^#/ {gsub(/^#/, "", $1); print $1; exit}')"
if [[ -z "$ISSUE_NUMBER" ]]; then
  fail "failed to determine issue number from issue list output"
fi

log "issue: view"
run_capture issue_view_out "$GC_BIN" issue view "$ISSUE_NUMBER" -R "$SYSTEM_REPO"
assert_contains "$issue_view_out" "#$ISSUE_NUMBER"

log "issue: list json"
run_capture issue_list_json "$GC_BIN" issue list -R "$SYSTEM_REPO" --limit 1 --json
printf '%s\n' "$issue_list_json" | assert_json
assert_contains "$issue_list_json" "\"number\""

log "issue: format json"
run_capture issue_format_json "$GC_BIN" issue list -R "$SYSTEM_REPO" --limit 1 --format json
printf '%s\n' "$issue_format_json" | assert_json

log "issue: view json"
run_capture issue_view_json "$GC_BIN" issue view "$ISSUE_NUMBER" -R "$SYSTEM_REPO" --json
printf '%s\n' "$issue_view_json" | assert_json
assert_contains "$issue_view_json" "\"number\""
