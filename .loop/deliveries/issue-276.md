# Delivery Record: Issue #276

- **Title**: add --yes to state-changing JSON examples
- **Type**: docs
- **Status**: merged
- **Loop**: fullflow-main
- **PR**: [#254](https://gitcode.com/gitcode-cli/cli/pulls/254)
- **Branch**: docs/issue-276-2
- **Started**: 2026-06-24 | **Merged**: 2026-06-24

## Gate Compliance

| # | Gate | Result | Evidence |
|---|------|--------|----------|
| 1 | 验证 | ✅ | issue close/reopen/pr close 的 --json 示例缺少 --yes |
| 2 | 开发 | ✅ | branch docs/issue-276-2, 3 lines: `--json` → `--yes --json` |
| 3 | 构建 | skipped | docs-only |
| 4 | UT | skipped | docs-only |
| 5 | Pre-commit | ✅ | 全部通过 |
| 6 | 实际命令 | skipped | docs-only |
| 7 | CI | skipped | docs-only |
| 8 | 风险分级 | ✅ | risk/low |
| + | 多角色评审 | ⚠️ | inline only (should have been 2-role docs+security per spec) |
| + | 合并 | ✅ | PR #254 merged |

## Key Artifacts
- Fix: docs/COMMANDS.md — issue close/reopen/pr close JSON 示例加 --yes

## Notes
- **首次 /loop 执行，漏掉了多角色评审** — 这是后来补 spec 5.3 门禁表的直接原因
- Branch naming: docs/issue-276-2 (first attempt docs/issue-276 conflicted)
