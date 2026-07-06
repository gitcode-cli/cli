#!/usr/bin/env bash

require_infra_repo "WRITE_REPO" "$WRITE_REPO"

if [[ -z "${GC_SYSTEM_PR_HEAD:-}" ]]; then
  skip "write pr create requires GC_SYSTEM_PR_HEAD"
  return 0
fi

log "write pr: create"
title="System Test PR $(date -u +%Y%m%d%H%M%S)"
run_capture pr_create_out "$GC_BIN" pr create -R "$WRITE_REPO" --head "$GC_SYSTEM_PR_HEAD" --base "${GC_SYSTEM_PR_BASE:-main}" --title "$title" --body "Created by tests/system"
assert_contains "$pr_create_out" "Created PR #"
pr_number="$(printf '%s\n' "$pr_create_out" | sed -n 's/.*#\([0-9][0-9]*\).*/\1/p' | head -n 1)"
if [[ -z "$pr_number" ]]; then
  fail "failed to parse created PR number"
fi

cleanup_pr() {
  "$GC_BIN" pr close "$pr_number" -R "$WRITE_REPO" --yes >/dev/null 2>&1 || true
}
trap cleanup_pr RETURN

run_capture pr_view_json "$GC_BIN" pr view "$pr_number" -R "$WRITE_REPO" --json
printf '%s\n' "$pr_view_json" | assert_json
assert_contains "$pr_view_json" "\"number\""

run_capture pr_close_out "$GC_BIN" pr close "$pr_number" -R "$WRITE_REPO" --yes
assert_contains "$pr_close_out" "Closed"
