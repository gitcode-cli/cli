# Delivery Record: Issue #302

- **Title**: TestEditRun_AuthError / TestResolveRunMissingRepo fails on CI
- **Type**: bug (test infra)
- **Status**: merged
- **Loop**: /goal relay (stage 2 unblock)
- **PR**: #252
- **Date**: 2026-06-24

## State Transitions
| From | To | When | Evidence |
|------|----|------|----------|
| (discovered) | (fixed) | 2026-06-24 | CI failure in /goal stage 2 |
| fix committed | merged | 2026-06-24 | PR #252 merged |

## Key Artifacts
- Fix: pkg/cmd/pr/comment/resolve/resolve_test.go — t.Setenv("GC_TOKEN", "dummy-token")
- Unblocked #298 CI

## Notes
- Only TestResolveRunMissingRepo was fixed; issue title mentions TestEditRun_AuthError (different test)
- Pattern: CI lacks GC_TOKEN, AuthenticatedClient fails before reaching test target
