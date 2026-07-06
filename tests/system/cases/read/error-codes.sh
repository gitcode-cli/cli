#!/usr/bin/env bash

log "error: non-git repo view"
tmp_non_git="$(mktemp -d)"
trap 'rm -rf "$tmp_non_git"' RETURN
run_expect_status nongit_out 2 bash -lc "cd '$tmp_non_git' && '$GC_BIN' repo view"
assert_contains "$nongit_out" "not in a git repository"

log "error: missing commit"
run_expect_status commit_missing_out 3 "$GC_BIN" commit view definitely-missing-sha-for-system-test -R "$SYSTEM_REPO"
assert_contains "$commit_missing_out" "Not Found Commit"

log "error: invalid issue format"
run_expect_status issue_format_out 2 "$GC_BIN" issue list -R "$SYSTEM_REPO" --format yaml
assert_contains "$issue_format_out" "invalid"
