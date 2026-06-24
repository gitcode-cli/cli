# Delivery Record: Issue #272

- **Title**: repo stats decodes commit_statistics with wrong field names
- **Type**: bug
- **Status**: merged
- **Loop**: fullflow-main (4f73b1f3)
- **PR**: #264
- **Date**: 2026-06-24

## State Transitions
| From | To | When | Evidence |
|------|----|------|----------|
| status/triage | status/verified | 2026-06-24 | confirmed all-zero data despite API returning values |
| status/verified | status/in-progress | 2026-06-24 | branch bugfix/issue-272 |
| status/in-progress | status/self-checked | 2026-06-24 | comment #176877432 |
| status/self-checked | status/merged | 2026-06-24 | PR #264 merged (risk/high, human confirmed) |

## Key Artifacts
- PR: #264 (merged)
- CI: 28077718073 (all 8 ✅)
- Real cmd: gc repo stats -R gitcode-cli/cli --branch main → aflyingto 79351/13644 ✅
- Fix: api/queries_repo.go — 4 JSON tags (author_name→user_name, commits→commit_count, additions→add_lines, deletions→delete_lines)

## Notes
- Data went from all-zero to correct values after JSON tag fix
- risk/high false positive; human confirmed
