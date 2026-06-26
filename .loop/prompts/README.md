# Prompt 模板索引

| # | 文件 | 命令 | 入口状态 | 出口状态 | 耗时 |
|---|------|------|---------|---------|------|
| 1 | `triage.md` | `/goal` | triage | verified | 2-5min |
| 2 | `develop-and-test.md` | `/goal` | verified | in-progress | 10-30min |
| 3 | `self-check.md` | `/goal` | in-progress | self-checked | 5-10min |
| 4 | `security-review.md` | `/goal` | in-progress | self-checked | 2-5min |
| 5 | `risk-classify.md` | `/goal` | self-checked | ready-for-review | 1min |
| 6 | `docs-sync.md` | `/goal` | any | any | 2-5min |
| 7 | `full-flow.md` | `/loop` | triage | merged | 30-60min |
| 8 | `ci-monitor.md` | `/loop` | in-progress | self-checked | 5-20min |
| 9 | `batch-triage.md` | `/loop` | triage | triage | 5-10min |
| 10 | `review-response.md` | `/loop` | ready-for-review | approved | 5-15min |
| 11 | `pr-review-patrol.md` | `/loop` | — | merged | 10-20min |
| 12 | `daily-audit.md` | `/loop` | — | — | 10-30min |

## 使用说明

1. 根据当前状态和任务类型选择模板
2. 复制 prompt 文本，替换 `<PLACEHOLDER>` 参数
3. 在 Claude Code 中运行 `/loop` 或 `/goal`
4. AI 会自动参考 `docs/LOOP-GOAL-GUIDE.md` 中的详细指南
