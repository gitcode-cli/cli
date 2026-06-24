# Delivery Record: Issue #273

- **Title**: issue edit --state should preserve title and normalize state
- **Type**: bug
- **Status**: merged
- **Loop**: fullflow-main (4f73b1f3)
- **PR**: #258
- **Date**: 2026-06-24

## State Transitions
| From | To | When | Evidence |
|------|----|------|----------|
| status/triage | status/verified | 2026-06-24 | confirmed open→reopen missing, title not preserved |
| status/verified | status/in-progress | 2026-06-24 | branch bugfix/issue-273 |
| status/in-progress | status/self-checked | 2026-06-24 | comment #176873474 |
| status/self-checked | status/merged | 2026-06-24 | PR #258 merged, CI 28075508360 ✅ |

## Key Artifacts
- PR: #258 (merged)
- CI: 28075508360 (all 8 ✅)
- Real cmd: gc issue edit 1 --state open/reopen/close -R infra-test/gctest1 ✅
- Fix: pkg/cmd/issue/edit/edit.go (+13)

## Notes
- Had edit tool tab-indentation issue; used python fallback
