# Delivery Record: Issue #279

- **Title**: fix release edit --help prerelease flag documentation
- **Type**: docs
- **Status**: merged
- **Loop**: fullflow-main (4f73b1f3)
- **PR**: [#268](https://gitcode.com/gitcode-cli/cli/pulls/268)
- **Branch**: bugfix/issue-279
- **Started**: 2026-06-24 | **Merged**: 2026-06-24

## Gate Compliance

| # | Gate | Result | Evidence |
|---|------|--------|----------|
| 1 | 验证 | ✅ | help `--prerelease` 裸写，代码 StringVar 需要 true/false |
| 2 | 开发 | ✅ | branch bugfix/issue-279, 1 file +7/-4 |
| 3 | 构建 | ✅ | `go build` Success, `./gc release edit --help` → 正确示例 |
| 4 | UT | ✅ | `go test ./...` 1199 passed |
| 5 | Pre-commit | ✅ | 全部通过 |
| 6 | 实际命令 | skipped | docs-only (help text 修改) |
| 7 | CI | skipped | docs-only |
| 8 | 风险分级 | ✅ | risk/low |
| + | 多角色评审 | skipped | docs-only, inline review |
| + | 合并 | ✅ | risk/low → 自动合并 |

## Key Artifacts
- Fix: pkg/cmd/release/edit/edit.go — Long(remove draft/target), Example(prerelease true/false)
