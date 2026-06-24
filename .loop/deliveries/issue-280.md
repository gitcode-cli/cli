# Delivery Record: Issue #280

- **Title**: document release download --all and source archive filtering
- **Type**: docs
- **Status**: merged
- **Loop**: fullflow-main (4f73b1f3)
- **PR**: [#265](https://gitcode.com/gitcode-cli/cli/pulls/265)
- **Branch**: docs/issue-280
- **Started**: 2026-06-24 | **Merged**: 2026-06-24

## Gate Compliance

| # | Gate | Result | Evidence |
|---|------|--------|----------|
| 1 | 验证 | ✅ | docs 缺 --all 示例和过滤说明 |
| 2 | 开发 | ✅ | branch docs/issue-280, 1 file +9/-2 |
| 3 | 构建 | skipped | docs-only |
| 4 | UT | skipped | docs-only |
| 5 | Pre-commit | ✅ | 全部通过 |
| 6 | 实际命令 | skipped | docs-only |
| 7 | CI | skipped | docs-only |
| 8 | 风险分级 | ✅ | risk/low |
| + | 多角色评审 | skipped | docs-only, inline review |
| + | 合并 | ✅ | risk/low → 自动合并 |

## Key Artifacts
- Fix: docs/COMMANDS.md — --all example + "默认过滤 source archive" 说明
