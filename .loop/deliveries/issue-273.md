# Delivery Record: Issue #273

- **Title**: issue edit --state should preserve title and normalize state
- **Type**: bug
- **Status**: merged
- **Loop**: fullflow-main (4f73b1f3)
- **PR**: [#258](https://gitcode.com/gitcode-cli/cli/pulls/258)
- **Branch**: bugfix/issue-273
- **Started**: 2026-06-24 | **Merged**: 2026-06-24

## Gate Compliance

| # | Gate | Result | Evidence |
|---|------|--------|----------|
| 1 | 验证 | ✅ | `open` 未规范化为 `reopen`，状态变更时 title 为空 |
| 2 | 开发 | ✅ | branch bugfix/issue-273, 1 file +13 |
| 3 | 构建 | ✅ | `go build -o ./gc ./cmd/gc` Success |
| 4 | UT | ✅ | `go test ./...` 1199 passed in 95 packages |
| 5 | Pre-commit | ✅ | 全部通过 |
| 6 | 实际命令 | ✅ | `./gc issue edit 1 --state open --yes -R infra-test/gctest1` ✅; `--state close` ✅; `--state reopen --yes --json` → title preserved ✅ |
| 7 | CI | ✅ | [run 28075508360](https://github.com/gitcode-cli/cli/actions/runs/28075508360) — All 8 ✅ |
| 8 | 风险分级 | ✅ | risk/medium |
| + | 多角色评审 | ✅ | inline 4/4 approved |
| + | 合并 | ✅ | PR #258 merged |

## Key Artifacts
- Fix: pkg/cmd/issue/edit/edit.go — open→reopen norm + GetIssue for title preservation

## Notes
- Python script needed for tab-indentation fix (Edit tool mismatch)
