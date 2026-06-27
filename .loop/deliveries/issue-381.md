# Issue #381 — refactor: full-flow-run.sh 死代码（已由 dfb8b20 解决）

| 字段 | 值 |
|------|-----|
| Issue | [#381](https://gitcode.com/gitcode-cli/cli/issues/381) |
| 类型 | type/refactor |
| 范围 | scope/loop-script |
| 风险 | risk/low |
| 状态 | status/merged |
| 解决方式 | 已由 commit [dfb8b20](https://gitcode.com/gitcode-cli/cli/commit/dfb8b20) 解决 |
| 完成时间 | 2026-06-27 15:39 |

## 问题

`.loop/scripts/full-flow-run.sh:238-241` 存在死代码：`if line.endswith('|'):` 的两个分支执行完全相同的操作。

## 分析

Commit `dfb8b20` (refactor(loop): extract token processing to standalone Python script) 将所有内联 Python 代码从 `full-flow-run.sh` 提取到 `process_tokens.py`。死代码随旧代码一起移除。新实现使用 split-join 方式更新 README 行，无死代码。

## 验证

- 旧版本 `full-flow-run.sh:240` 存在 `if line.endswith('|'):` 双分支相同操作的死代码 ✓
- `dfb8b20` 移除了 197 行内联 Python 代码 ✓
- 新 `process_tokens.py:149-162` 无死代码 ✓
- 新代码为纯 Python 文件，无 shell 变量注入风险 ✓

## 门禁

| # | 门禁 | 结果 | 说明 |
|---|------|:--:|------|
| 1 | 实现 | N/A | 已由 dfb8b20 解决，无新代码 |
| 2 | 测试 | N/A | 非 Go 代码 |
| 3 | 构建 | N/A | 非 Go 代码 |
| 4 | UT | N/A | 非 Go 代码 |
| 5 | Pre-commit | N/A | 无新代码变更 |
| 6 | 命令验证 | N/A | 无 CLI 变更 |
| 7 | CI | N/A | 无 PR |
| 8 | 风险分级 | low | risk/low |

## 结论

Issue 在提交 dfb8b20 时已自动解决，无需额外代码变更。关闭并标记为 merged。

## Token 消耗

| 指标 | 值 |
|------|-----|
| 输入 tokens (cache miss) | 51,894 (52k) |
| 输出 tokens | 11,010 (11k) |
| 缓存命中 | 1,636,352 (1636k) |
| 缓存写入 | 0 |
| 总计 tokens | 62,904 (63k) |
| 成本 (DeepSeek) | ¥0.2627 (~$0.0365) |
| 耗时 | 188s |
| 轮次 | 42 |

> 计价: ¥3/M cache-miss + ¥0.025/M cache-hit + ¥6/M output
