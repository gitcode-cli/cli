# Delivery Record: Issue #276

- **Title**: add --yes to state-changing JSON examples
- **Type**: docs
- **Status**: merged
- **Loop**: fullflow-main
- **PR**: #254
- **Date**: 2026-06-24

## State Transitions
| From | To | When | Evidence |
|------|----|------|----------|
| status/triage | status/verified | 2026-06-24 | confirmed 3 lines missing --yes |
| status/verified | status/in-progress | 2026-06-24 | branch docs/issue-276-2 |
| status/in-progress | status/self-checked | 2026-06-24 | comment #176862928 |
| status/self-checked | status/merged | 2026-06-24 | PR #254 merged |

## Key Artifacts
- PR: #254 (merged)
- Fix: docs/COMMANDS.md — issue close/reopen/pr close --json 加 --yes

## Notes
- docs-only, skipped build/UT/CI per spec 5.3
- Review done inline (not formal 4-agent) — flagged as gap
