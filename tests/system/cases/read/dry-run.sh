#!/usr/bin/env bash

log "dry-run: repo delete"
run_capture repo_delete_dry "$GC_BIN" repo delete "$SYSTEM_REPO" --dry-run
assert_contains "$repo_delete_dry" "Dry run"

log "dry-run: issue create"
run_capture issue_create_dry "$GC_BIN" issue create -R "$SYSTEM_REPO" --title "System test dry-run" --body "dry-run" --dry-run
assert_contains "$issue_create_dry" "Dry run"

log "dry-run: label delete"
run_capture label_delete_dry "$GC_BIN" label delete system-test-nonexistent -R "$SYSTEM_REPO" --dry-run
assert_contains "$label_delete_dry" "Dry run"
