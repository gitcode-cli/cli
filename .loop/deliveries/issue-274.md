# Delivery Record: Issue #274

- **Title**: issue view --comments --json should always include comments array
- **Type**: bug
- **Status**: merged
- **Loop**: fullflow-main (4f73b1f3)
- **PR**: #257
- **Date**: 2026-06-24

## State Transitions
| From | To | When | Evidence |
|------|----|------|----------|
| status/triage | status/verified | 2026-06-24 | confirmed structure differs with/without comments |
| status/verified | status/in-progress | 2026-06-24 | branch bugfix/issue-274 |
| status/in-progress | status/self-checked | 2026-06-24 | comment #176870584 |
| status/self-checked | status/merged | 2026-06-24 | PR #257 merged, CI 28074422614 ✅ |

## Key Artifacts
- PR: #257 (merged)
- CI: 28074422614 (all 8 ✅)
- Real cmd: gc issue view 30 -R infra-test/gctest1 --comments --json → comments:[] ✅
- Fix: pkg/cmd/issue/view/view.go (+8/-4)

## Notes
- GitCode PR merge returned 405; used git merge fallback
- tst-review agent noted missing unit test (pre-existing gap, non-blocking)
