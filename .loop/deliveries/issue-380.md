#380 docs: full-flow-subprocess.md 裸引用修复 | 2026-06-27
PR: #305 | risk: low | docs-only
Gates: G5(✅) G8(✅) | CI: skipped
Commit: cbeca2c

## Token 消耗

| 指标 | 值 |
|------|-----|
| 输入 tokens (cache miss) | 41,896 (42k) |
| 输出 tokens | 10,001 (10k) |
| 缓存命中 | 1,481,216 (1481k) |
| 缓存写入 | 0 |
| 总计 tokens | 51,897 (52k) |
| 成本 (DeepSeek) | ¥0.2227 (~$0.0309) |
| 耗时 | 177s |
| 轮次 | 47 |

> 计价: ¥3/M cache-miss + ¥0.025/M cache-hit + ¥6/M output
