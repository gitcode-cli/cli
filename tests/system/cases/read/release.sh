#!/usr/bin/env bash

log "release: list json"
if probe_or_skip "release list" "$GC_BIN" release list -R "$SYSTEM_REPO" --json; then
  run_capture release_list_json "$GC_BIN" release list -R "$SYSTEM_REPO" --json
  printf '%s\n' "$release_list_json" | assert_json
fi

log "release: delete dry-run"
release_tag="${GC_SYSTEM_RELEASE_TAG:-v0.0.1-test}"
run_capture release_delete_dry "$GC_BIN" release delete "$release_tag" -R "$SYSTEM_REPO" --dry-run
assert_contains "$release_delete_dry" "Dry run: would delete release"
