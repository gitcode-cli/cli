# Delivery: Issue #382

- **Issue**: [#382](https://gitcode.com/gitcode-cli/cli/issues/382) — docs: COMMANDS.md 全局 --json 表格与 9 个独立章节 --json 示例不同步
- **PR**: [#310](https://gitcode.com/gitcode-cli/cli/pulls/310) — docs: add --json examples to 9 command sections in COMMANDS.md
- **Type**: docs-only
- **Risk**: low
- **Branch**: worktree-issue-382-20260627
- **Date**: 2026-06-27

## Changes

Added `# Output as JSON` example lines to 9 command sections:
1. auth token
2. help
3. version
4. pr diff
5. issue label
6. label create
7. label delete
8. milestone create
9. repo delete

## Gates

| # | Gate | Status |
|---|------|--------|
| 1 | Implementation | ✅ |
| 2 | Tests | N/A (docs-only) |
| 3 | Build | N/A (docs-only) |
| 4 | UT | N/A (docs-only) |
| 5 | Pre-commit | ✅ all pass |
| 6 | Command verification | N/A (docs-only) |
| 7 | CI | N/A (no code changes) |
| 8 | Risk classification | risk/low |

## Verification

- Pre-commit: 10/10 pass
- Grep verification: all 9 sections contain `# Output as JSON`

## Token 消耗

| 指标 | 值 |
|------|-----|
| 输入 tokens (cache miss) | 69,168 (69k) |
| 输出 tokens | 11,748 (12k) |
| 缓存命中 | 2,059,904 (2060k) |
| 缓存写入 | 0 |
| 总计 tokens | 80,916 (81k) |
| 成本 (DeepSeek) | ¥0.3295 (~$0.0458) |
| 耗时 | 184s |
| 轮次 | 50 |

> 计价: ¥3/M cache-miss + ¥0.025/M cache-hit + ¥6/M output
