#385 refactor: full-flow-run.sh 存在冗余 README_FILE 赋值 | 2026-06-27
PR: N/A (pre-existing fix) | risk: low | code-change
Gates: G5(✅) G8(✅) | CI: skipped (fix already on main)
Commit: dfb8b20 (refactor(loop): extract token processing to standalone Python script)

## 处置说明

Issue 描述的冗余 `README_FILE` 赋值（`full-flow-run.sh` 第 179 行和第 226 行完全相同）已在 commit `dfb8b20` 中修复。

该重构将约 200 行内联 Python 代码从 `full-flow-run.sh` 提取到独立的 `process_tokens.py`，重复的 `README_FILE` 赋值被一并消除：
- `full-flow-run.sh`: 60 行，不再包含 `README_FILE`
- `process_tokens.py`: 仅在 114 行定义一次 `readme_file`

## Token 消耗

| 指标 | 值 |
|------|-----|
| 输入 tokens (cache miss) | 43,835 (44k) |
| 输出 tokens | 10,122 (10k) |
| 缓存命中 | 1,541,888 (1542k) |
| 缓存写入 | 0 |
| 总计 tokens | 53,957 (54k) |
| 成本 (DeepSeek) | ¥0.2308 (~$0.0321) |
| 耗时 | 170s |
| 轮次 | 57 |

> 计价: ¥3/M cache-miss + ¥0.025/M cache-hit + ¥6/M output
