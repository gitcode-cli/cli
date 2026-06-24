# Delivery Record: Issue #251

- **Title**: TestViewRunUsesDetectedRepo fails - invalid time format
- **Type**: bug (test infra)
- **Status**: closed-no-fix
- **Loop**: fullflow-main (4f73b1f3)
- **Date**: 2026-06-24

## Gate Compliance

| # | Gate | Result | Evidence |
|---|------|--------|----------|
| 1 | 验证 | ✅ | `go test ./pkg/cmd/repo/view/ -run TestViewRunUsesDetectedRepo` → PASS |
| 2 | 开发 | N/A | 已有 `t.Setenv("GC_TOKEN", "test-token")` 修复 |
| 3-8 | — | skipped | 问题已修复 |

## Notes
- 与 #302 同类问题，已被先前 PR 修复
