#!/usr/bin/env bash

require_infra_repo "WRITE_REPO" "$WRITE_REPO"

if [[ -z "${GC_SYSTEM_PR_HEAD:-}" ]]; then
  skip "write pr create requires GC_SYSTEM_PR_HEAD"
  return 0
fi

log "write pr: create"
title="System Test PR $(date -u +%Y%m%d%H%M%S)"
label_suffix="$(date -u +%Y%m%d%H%M%S)-$$"
label_a="system-test-pr-a-$label_suffix"
label_b="system-test-pr-b-$label_suffix"
run_capture pr_create_out "$GC_BIN" pr create -R "$WRITE_REPO" --head "$GC_SYSTEM_PR_HEAD" --base "${GC_SYSTEM_PR_BASE:-main}" --title "$title" --body "Created by tests/system"
assert_contains "$pr_create_out" "Created PR #"
pr_number="$(printf '%s\n' "$pr_create_out" | sed -n 's/.*#\([0-9][0-9]*\).*/\1/p' | head -n 1)"
if [[ -z "$pr_number" ]]; then
  fail "failed to parse created PR number"
fi

cleanup_pr_and_labels() {
  "$GC_BIN" pr close "$pr_number" -R "$WRITE_REPO" --yes >/dev/null 2>&1 || true
  "$GC_BIN" label delete "$label_a" -R "$WRITE_REPO" --yes >/dev/null 2>&1 || true
  "$GC_BIN" label delete "$label_b" -R "$WRITE_REPO" --yes >/dev/null 2>&1 || true
}
trap cleanup_pr_and_labels RETURN

run_capture pr_view_json "$GC_BIN" pr view "$pr_number" -R "$WRITE_REPO" --json
printf '%s\n' "$pr_view_json" | assert_json
assert_contains "$pr_view_json" "\"number\""

log "write pr: label edit"
run_capture label_a_out "$GC_BIN" label create "$label_a" -R "$WRITE_REPO" --color '#3366cc'
assert_contains "$label_a_out" "$label_a"
run_capture label_b_out "$GC_BIN" label create "$label_b" -R "$WRITE_REPO" --color '#33aa66'
assert_contains "$label_b_out" "$label_b"

run_capture pr_add_labels_out "$GC_BIN" pr edit "$pr_number" -R "$WRITE_REPO" --add-label "$label_a"
assert_contains "$pr_add_labels_out" "Updated PR"
run_capture pr_legacy_labels_out "$GC_BIN" pr edit "$pr_number" -R "$WRITE_REPO" --labels "$label_b"
assert_contains "$pr_legacy_labels_out" "Updated PR"
run_capture pr_labels_json "$GC_BIN" pr view "$pr_number" -R "$WRITE_REPO" --json
assert_contains "$pr_labels_json" "$label_a"
assert_contains "$pr_labels_json" "$label_b"

run_capture pr_replace_labels_out "$GC_BIN" pr edit "$pr_number" -R "$WRITE_REPO" --replace-labels "$label_b" --yes
assert_contains "$pr_replace_labels_out" "Updated PR"
run_capture pr_labels_json "$GC_BIN" pr view "$pr_number" -R "$WRITE_REPO" --json
if [[ "$pr_labels_json" == *"$label_a"* ]]; then
  fail "replaced label set still contains $label_a"
fi
assert_contains "$pr_labels_json" "$label_b"

run_capture pr_clear_labels_out "$GC_BIN" pr edit "$pr_number" -R "$WRITE_REPO" --replace-labels "" --yes
assert_contains "$pr_clear_labels_out" "Updated PR"
run_capture pr_labels_json "$GC_BIN" pr view "$pr_number" -R "$WRITE_REPO" --json
if [[ "$pr_labels_json" == *"$label_b"* ]]; then
  fail "cleared label set still contains $label_b"
fi

run_capture pr_add_label_out "$GC_BIN" pr edit "$pr_number" -R "$WRITE_REPO" --add-label "$label_b"
assert_contains "$pr_add_label_out" "Updated PR"
run_capture pr_remove_label_out "$GC_BIN" pr edit "$pr_number" -R "$WRITE_REPO" --remove-label "$label_b"
assert_contains "$pr_remove_label_out" "Updated PR"
run_capture pr_labels_json "$GC_BIN" pr view "$pr_number" -R "$WRITE_REPO" --json
if [[ "$pr_labels_json" == *"$label_b"* ]]; then
  fail "removed label set still contains $label_b"
fi

run_capture pr_close_out "$GC_BIN" pr close "$pr_number" -R "$WRITE_REPO" --yes
assert_contains "$pr_close_out" "Closed"
