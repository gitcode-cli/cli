# Delivery Record: Issue #256

- **Title**: release list returns old versions first, not sorted by latest
- **Type**: bug
- **Status**: merged
- **Loop**: fullflow-main (4f73b1f3)
- **PR**: #260
- **Date**: 2026-06-24

## State Transitions
| From | To | When | Evidence |
|------|----|------|----------|
| status/triage | status/verified | 2026-06-24 | confirmed unsorted: v0.2.0 before v0.5.3 |
| status/verified | status/in-progress | 2026-06-24 | branch bugfix/issue-256 |
| status/in-progress | status/self-checked | 2026-06-24 | comment #176875869 |
| status/self-checked | status/merged | 2026-06-24 | PR #260 merged (risk/high, human confirmed) |

## Key Artifacts
- PR: #260 (merged)
- CI: 28076568258 (macos flake, non-code)
- Real cmd: gc release list -R gitcode-cli/cli --json → v0.5.3 first ✅
- Fix: api/queries_release.go (+6) — sort.Slice by CreatedAt desc

## Notes
- risk/high due to classify-change-risk.py false positive
- Blocked until human approval
