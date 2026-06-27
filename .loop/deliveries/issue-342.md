# Issue #342 — Delivery Record

## 概要

- **Issue**: [#342](https://gitcode.com/gitcode-cli/cli/issues/342)
- **PR**: [#?](https://gitcode.com/gitcode-cli/cli/pulls/300)
- **类型**: type/bug
- **风险**: risk/medium

## Token 消耗

| 指标 | 值 |
|------|-----|
| 输入 tokens (cache miss) | 291,451 (291k) |
| 输出 tokens | 71,548 (72k) |
| 缓存命中 | 8,319,744 (8.3M) |
| 总计 tokens | 362,999 (363k) |
| 成本 (DeepSeek) | ¥1.5116 (~$0.2099) |
| 耗时 | 761s |
| 轮次 | 112 |

> 计价: ¥3/M cache-miss + ¥0.025/M cache-hit + ¥6/M output
