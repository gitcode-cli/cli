# Delivery Record: Issue #256

- **Title**: release list returns old versions first, not sorted by latest
- **Type**: bug
- **Status**: merged
- **Loop**: fullflow-main (4f73b1f3)
- **PR**: [#260](https://gitcode.com/gitcode-cli/cli/pulls/260)
- **Branch**: bugfix/issue-256
- **Started**: 2026-06-24 | **Merged**: 2026-06-24

## Gate Compliance

| # | Gate | Result | Evidence |
|---|------|--------|----------|
| 1 | 验证 | ✅ | `--limit 5` 返回 v0.2.0~v0.2.4，v0.5.3 不在列表 |
| 2 | 开发 | ✅ | branch bugfix/issue-256, 1 file +6 (+import sort) |
| 3 | 构建 | ✅ | `go build -o ./gc ./cmd/gc` Success |
| 4 | UT | ✅ | `go test ./...` 1199 passed in 95 packages |
| 5 | Pre-commit | ✅ | 全部通过 |
| 6 | 实际命令 | ✅ | `./gc release list -R gitcode-cli/cli --json` → v0.5.3 first ✅ |
| 7 | CI | ⚠️ | [run 28076568258](https://github.com/gitcode-cli/cli/actions/runs/28076568258) — macos auth flake (非代码) |
| 8 | 风险分级 | ⚠️ | risk/high (误判，classify-change-risk.py 扫描累积 diff) |
| + | 多角色评审 | ✅ | inline approved |
| + | 合并 | ✅ | 人工确认后合并 |

## Key Artifacts
- Fix: api/queries_release.go — `sort.Slice(releases, ...)` by CreatedAt desc

## Notes
- risk/high due to classify-change-risk.py false positive; human confirmed
- macos CI failure was auth-package test, unrelated to sort change
