# Delivery Record: Issue #277

- **Title**: clarify that gc schema excludes help and completion commands
- **Type**: docs
- **Status**: merged
- **PR**: [#272](https://gitcode.com/gitcode-cli/cli/pulls/272)
- **Branch**: docs/issue-277
- **Date**: 2026-06-24

## Gate Compliance
| # | Gate | Result | Evidence |
|---|------|--------|----------|
| 1 | 验证 | ✅ | schema.go line 111-112 排除了 help/completion |
| 2 | 开发 | ✅ | docs/COMMANDS.md +1/-1 |
| 3-7 | — | skipped | docs-only |
| 8 | 风险 | ✅ | risk/low |
| + | 合并 | ✅ | auto-merged |
