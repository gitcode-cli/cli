#!/usr/bin/env bash

require_infra_repo "WRITE_REPO" "$WRITE_REPO"

log "write issue: create, view, close"
title="System Test $(date -u +%Y%m%d%H%M%S)"
body="Created by tests/system; safe to close."

run_capture issue_create_out "$GC_BIN" issue create -R "$WRITE_REPO" --title "$title" --body "$body"
assert_contains "$issue_create_out" "#"
issue_number="$(printf '%s\n' "$issue_create_out" | sed -n 's/.*#\([0-9][0-9]*\).*/\1/p' | head -n 1)"
if [[ -z "$issue_number" ]]; then
  fail "failed to parse created issue number"
fi

cleanup_issue() {
  "$GC_BIN" issue close "$issue_number" -R "$WRITE_REPO" --yes >/dev/null 2>&1 || true
}
trap cleanup_issue RETURN

run_capture issue_view_json "$GC_BIN" issue view "$issue_number" -R "$WRITE_REPO" --json
printf '%s\n' "$issue_view_json" | assert_json
assert_contains "$issue_view_json" "\"number\""

run_capture issue_close_out "$GC_BIN" issue close "$issue_number" -R "$WRITE_REPO" --yes
assert_contains "$issue_close_out" "Closed"
