#!/usr/bin/env bash

log "pr: fixture probe"
if ! probe_or_skip "pr list" "$GC_BIN" pr list -R "$SYSTEM_REPO" --limit 1 --json; then
  return 0
fi

log "pr: list json"
run_capture pr_list_json "$GC_BIN" pr list -R "$SYSTEM_REPO" --limit 1 --json
printf '%s\n' "$pr_list_json" | assert_json

log "pr: paginated list json"
run_capture pr_paginate_json "$GC_BIN" pr list -R "$SYSTEM_REPO" --paginate --per-page 1 --limit 1 --json
printf '%s\n' "$pr_paginate_json" | assert_json

PR_NUMBER="$(printf '%s\n' "$pr_list_json" | python3 -c 'import json,sys; data=json.load(sys.stdin); print(data[0].get("number","") if data else "")')"
if [[ -n "$PR_NUMBER" ]]; then
  log "pr: view json"
  run_capture pr_view_json "$GC_BIN" pr view "$PR_NUMBER" -R "$SYSTEM_REPO" --json
  printf '%s\n' "$pr_view_json" | assert_json
fi
