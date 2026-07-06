#!/usr/bin/env bash

log "auth: status"
run_capture auth_status_out "$GC_BIN" auth status
assert_contains "$auth_status_out" "Logged in as"

log "auth: token non-interactive guard"
if token_out="$(GC_TOKEN="fake-token-for-system-test" GITCODE_TOKEN= "$GC_BIN" auth token 2>&1)"; then
  fail "auth token unexpectedly succeeded in non-interactive mode"
fi
assert_contains "$token_out" "interactive confirmation"

log "auth: unauthenticated command exit code"
tmp_empty_config="$(mktemp -d)"
trap 'rm -rf "$tmp_empty_config"' RETURN
run_expect_status unauth_out 4 env GC_CONFIG_DIR="$tmp_empty_config" GC_TOKEN= GITCODE_TOKEN= "$GC_BIN" pr review 1 -R "$SYSTEM_REPO" --approve
assert_contains "$unauth_out" "not authenticated"
