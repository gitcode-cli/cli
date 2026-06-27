# Issue #379 — Delivery Record

- **标题**: bug: .loop/deliveries/README.md Issue #371 重复行 — CI 值冲突
- **Issue**: [#379](https://gitcode.com/gitcode-cli/cli/issues/379)
- **PR**: [#309](https://gitcode.com/gitcode-cli/cli/pulls/309)
- **CI**: [✅](https://github.com/gitcode-cli/cli/actions/runs/28284858941) (8/8 jobs green)
- **风险**: low
- **类型**: bug (代码改动)
- **状态**: merged

## 问题

`.loop/deliveries/README.md` 中 Issue #371 出现两行:
- Line 9: CI=skipped 门禁=2/8 (旧数据)
- Line 63: CI=ok 门禁=5/8 (新数据)

根因: `process_tokens.py` 只更新第一个匹配行的 token 列，而 agent 在 finalize 阶段追加新行而非原地更新旧行。

## 修复

1. **`.loop/deliveries/README.md`**: 删除 #371 的 stale 行 (CI=skipped, gate=2/8)
2. **`.loop/scripts/process_tokens.py`**:
   - 收集所有匹配行 (不仅是第一个)
   - 保留最后一行 (最新)，删除之前的重复
   - 新增 `else` 分支: 无匹配行时在 `## 统计` 前插入新行
   - 增加表格行守卫 (必须首尾都是 `|`) 避免在引用链接上误匹配

## 门禁证据

| # | 门禁 | 状态 | 证据 |
|---|------|:----:|------|
| 1 | 实现 | ✅ | 删除重复行 + process_tokens.py 去重逻辑 |
| 2 | go test ./... | ✅ | 1259 passed in 96 packages |
| 3 | go build | ✅ | Build success |
| 4 | UT | ✅ | All passing |
| 5 | Pre-commit | ✅ | gc precommit check pass |
| 6 | 命令验证 | ✅ | regression-core.sh (auth+repo+issue flow pass) |
| 7 | CI | ✅ | [28284858941](https://github.com/gitcode-cli/cli/actions/runs/28284858941) — Lint + Test(3 OS) + Build(3 OS) + Docker all green |
| 8 | 风险分级 | ✅ | risk/low (classifier 误报 high — 关键词匹配，非安全变更) |

## 变更

| 文件 | 变更 | 说明 |
|------|:----:|------|
| `.loop/deliveries/README.md` | -1 | 删除 stale 重复行 |
| `.loop/scripts/process_tokens.py` | +9/-13 | 去重逻辑 + 行插入 fallback |

## Token 消耗

| 指标 | 值 |
|------|-----|
| 输入 tokens (cache miss) | 74,822 (75k) |
| 输出 tokens | 27,858 (28k) |
| 缓存命中 | 4,561,280 (4561k) |
| 缓存写入 | 0 |
| 总计 tokens | 102,680 (103k) |
| 成本 (DeepSeek) | ¥0.5056 (~$0.0702) |
| 耗时 | 659s |
| 轮次 | 92 |

> 计价: ¥3/M cache-miss + ¥0.025/M cache-hit + ¥6/M output
