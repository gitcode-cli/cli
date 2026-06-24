# Delivery Record: Issue #279

- **Title**: fix release edit --help prerelease flag documentation
- **Type**: docs
- **Status**: merged
- **Loop**: fullflow-main (4f73b1f3)
- **PR**: #268
- **Date**: 2026-06-24

## State Transitions
| From | To | When | Evidence |
|------|----|------|----------|
| status/triage | status/verified | 2026-06-24 | help shows bare --prerelease, code uses StringVar(true/false) |
| status/verified | status/merged | 2026-06-24 | PR #268 merged (docs-only path) |

## Key Artifacts
- PR: #268 (merged)
- Fix: pkg/cmd/release/edit/edit.go — help text: bare flag→string value, remove draft/target from Long
- Risk: low → auto-merged
