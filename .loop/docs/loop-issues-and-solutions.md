# Loop 系统问题与解决方案

本文档记录 `/loop` 全流程交付系统在设计、实施和运行过程中遇到的所有问题及解决方案。

**最后更新**: 2026-06-27

---

## 一、架构问题

### 1. 上下文累积导致质量衰减

**现象**: 同一个 session 内连续处理多个 issue 后，模型开始丢失步骤、简化输出、跳过门禁。后期 issue 的 PR 评论变得极度简化。

**根因**: CronCreate 每次触发在当前 session 内注入消息，所有历史对话累积在上下文中。处理 38 个 issue 后上下文严重膨胀，模型"凭记忆脑补"而非严格按 spec 执行。

**方案**: `claude -p` 子进程独立上下文

```
Cron (父 session) → bash 脚本 → claude -p (子进程)
```

每次触发启动全新的 `claude -p` 进程，独立 session，互不干扰。父 session 只记录"启动 script + done"，上下文增量极小。

**相关文件**:
- `.loop/scripts/full-flow-run.sh`
- `.loop/prompts/full-flow-subprocess.md`

---

### 2. 同一时刻多 loop 冲突

**现象**: 全流程交付 loop (:00/:30) 和 PR 巡逻 loop (:05/:35) 可能同时操作同一文件或 git worktree。

**方案**:
- 错开 cron 调度（5 分钟偏移）
- worktree 命名加时间戳确保唯一：`.claude/worktrees/issue-N-timestamp`
- PID 锁防同类型 loop 重叠

---

### 3. 子进程被父进程退出连带 kill

**现象**: Bash 脚本启动 `claude -p` 后退出时，子进程被连带终止。导致多次交付中断。

**方案**: `nohup` 解耦

```bash
cat "$PROMPT_FILE" | nohup claude -p ... > "$JSONL_FILE" 2>&1
```

---

### 4. claude -p 初始化偶发卡死

**现象**: `claude -p` 加载 hook/plugin/skill 后停在 thinking_tokens 状态，不产出第一条 assistant 消息。JSONL 只有 init 和 hook 输出（49 行左右）。

**根因**: API 侧问题（模型响应超时、rate limit 等），脚本层面无法预防。

**方案**: `timeout 1200`（20 分钟）硬超时

```bash
cat "$PROMPT_FILE" | nohup timeout 1200 claude -p ...
```

正常交付 10-13 分钟，20 分钟给足余量。超时后 SIGKILL，下个 cron tick 自动重试。

---

### 5. Run in background 的双重后台陷阱

**现象**: 在 `run_in_background: true` 的 Bash 命令里又写了 `&`，导致 task 立即返回，脚本被截断。

```bash
# 错误
bash script.sh &    # run_in_background 已经异步，& 导致立即返回

# 正确
bash script.sh      # run_in_background 处理异步
```

**方案**: Cron prompt 和手动执行中去掉 `&`。

---

## 二、门禁与 CI 问题

### 6. CI 门禁被跳过

**现象**: 早期 `claude -p` 交付的 issue 没有跑 GitHub Actions CI，README 却乐观地填了 ✅。

**根因**: prompt 只写"按 spec 推进"，但子进程没实际读 spec，跳过 CI。

**方案**: prompt 显式列举 8 门禁表 + CI 执行命令

```markdown
| 7 | 远端 CI | docs-only 跳过；代码改动必须用 gh CLI 触发
              GitHub Actions CI 并等待全部 Job 通过，
              PR 评论中附 run URL |
CI 未跑就写 ✅ 算违规。
```

---

### 7. Docker CI 预存失败

**现象**: 每次 CI 的 Docker job 报 `Binary not found for your platform`，整体标记为 failure。

**根因**: CI workflow 的 "Build binary" 步骤构建 Go 二进制到 `./gc`，但 Python wrapper 期望 `gc_cli/bin/gc-linux-amd64`。

**方案**: 在 CI workflow Docker job 中添加 copy 步骤

```bash
mkdir -p gc_cli/bin
cp gc gc_cli/bin/gc-linux-amd64
```

---

### 8. macOS CI dyld 崩溃

**现象**: macOS test job 报 `dyld: missing LC_UUID load command → signal: abort trap`。

**根因**: Go linker 与 macOS 环境的兼容性问题，非代码 bug。预存问题。

**方案**: 识别为预存问题，不阻塞合并。其他平台 test 全部通过即可。

---

### 9. CI 修复后未立即验证

**现象**: Docker CI 修复推送到远端后，后续交付仍标记 CI 为 ⚠️。

**根因**: 修复已推送但历史 CI run 记录不变，需新 PR 触发后验证。

**方案**: 修复提交后第一个新 PR 会带新 CI run，验证通过后改为 ✅。

---

## 三、Token 追踪问题

### 9. 子进程 Token 无法获取

**现象**: 父 session 的 HUD 显示 Token 456M，但这是全部累积（含大量 cache），无法反映单个交付的成本。

**方案**: `claude -p --output-format stream-json --verbose`

`result` 对象包含完整 session 统计：
```json
{
  "usage": {
    "input_tokens": 70616,
    "cache_read_input_tokens": 5992704,
    "output_tokens": 26394
  },
  "modelUsage": {
    "deepseek-v4-pro[1m]": {
      "inputTokens": 136725,
      "outputTokens": 41882,
      "cacheReadInputTokens": 6609536
    }
  }
}
```

Post-processing 自动提取并注入交付文件。

---

### 10. Anthropic 计价不匹配

**现象**: 系统默认用 Anthropic 计价（$5.04/次），实际使用 DeepSeek 模型。

**方案**: 按 DeepSeek 费率重新计算

| 类型 | 单价 (¥/M tokens) |
|------|:--:|
| 输入 (cache miss) | ¥3 |
| 输入 (cache hit) | ¥0.025 |
| 输出 | ¥6 |

```python
cost = (model_in/1e6*3) + (model_cache/1e6*0.025) + (model_out/1e6*6)
```

**效果**: #327 成本从 Anthropic 的 $5.04 修正为 ¥0.83 (~$0.11)。

---

### 11. `grep -oP` 兼容性

**现象**: Bash 脚本用 `grep -oP 'ISSUE_NUM=\K\d+'` 提取 issue 号，部分环境 Perl 正则不可用。

**方案**: 替换为 Python 正则，并从 JSONL 直接搜索（跳过不可靠的 .log 中间文件）。

```python
m = re.search(r'ISSUE_NUM=(\d+)', content)
```

---

## 四、交付记录问题

### 12. Post-process 偶发失败

**现象**: Bash wrapper 退出码 1，token 注入和统计刷新未执行。JSONL 有完整的 token 数据但交付文件缺失。

**根因**: `set -e` 在 post-processing 阶段遇到非关键错误（如 Python 脚本异常）导致整个脚本退出。

**方案**: 整个 post-processing 段设置为 `set +e`，所有步骤变为 best-effort。

---

### 13. 交付文件随 worktree 丢失

**现象**: `claude -p` 在 worktree 里写 `issue-N.md`，删除 worktree 后文件丢失，未拷贝回主目录。

**方案**: prompt 要求在所有操作完成后，在 main 目录中重新创建丢失的文件。同时 post-process 从 JSONL 恢复 token 数据写入。

---

### 14. README 宽表格被 GitCode 截断

**现象**: GitCode blob 页面容器宽度不够，8+ 列表格最后一列（Token）被截断不显示。

**方案（待实施）**: 将 Token/成本列拆为独立表格，放在主汇总表下方。

---

### 15. 统计需手动更新

**现象**: README 底部的统计表（总 issue、已合并、平均门禁等）数据滞后。

**方案**: `scripts/count-deliveries.sh` 自动解析汇总表并刷新统计

每次交付完成后自动调用，统计项包括：
- 总 issue / 已合并 / 已关闭
- 含代码改动 / docs-only / risk/high
- 平均门禁 / Token 总消耗 / 总成本 / 今日交付

---

## 五、PR 管理问题

### 16. 孤儿 PR

**现象**: 子进程被 kill 或超时时，已创建但未合并的 PR 遗留为孤儿。已发生 3+ 次。

**方案**: prompt 增加"孤儿 PR 优先 rescue"步骤

1. 列出所有 open PR，过滤本人创建的
2. 完整读取 PR + Issue 所有评论
3. 对照 spec 逐项检查门禁完成度
4. 只补充真正缺失的（不重复已有工作）
5. 全部通过后合并

---

### 17. 子进程违规修改主目录

**现象**: `claude -p` 有时不创建 worktree 直接在主目录修改文件，违反"所有操作在 worktree 中"的规则。

**方案**: 在 prompt 中强调 worktree 是前置规则，违规变更在 push 前被识别并回退。

---

### 21. Token 消耗持续增长

**现象**: 单次交付 token 从 179k (#327) 涨到 369k (#351)，代码改动量不成比例（+2/-1 vs +19/-19）。

**根因**:
1. prompt 膨胀：指令反复重复
2. 每个 session 重复读 spec 文件
3. 孤儿 PR 检查无条件执行
4. docs-only 也走了完整评审流程

**方案**:
1. prompt 去重精简（1.5KB）
2. 孤儿 PR 检查延后到 triage 队列空时
3. docs-only 显式跳过评审
4. spec 关键规则内联到 prompt（待实施）

---

### 22. Token 注入 100% 失败——`"$TOKEN_FILE"` bash 引号截断

**现象**: Post-process 报告 "no issue num or token data, skipping delivery update"。Token `.json` 文件从未被创建。所有交付需手动注入 token 数据。

**初步误判**: 以为是从 `.log` 文本提取 ISSUE_NUM 的竞态条件。

**真正根因**: 脚本中 `python3 -c "..."` 内部使用了 Python 双引号 `"$TOKEN_FILE"`。在 bash 的 `"..."` 字符串中，内层 `"` 被 bash 解析为 `-c` 参数的结束符，导致整段 Python 代码被截断、静默失败。

```bash
# 错误 — bash 看到 " 就认为 -c 参数结束
echo "$RESULT" | python3 -c "
    with open("$TOKEN_FILE", "w") as f:  # ← 这个 " 终止了 -c
        json.dump({...}, f)
" >> log

# 正确 — Python 单引号，bash 不截断
echo "$RESULT" | python3 -c "
    with open('$TOKEN_FILE', 'w') as f:  # ← ' ' 对 bash 透明
        json.dump({...}, f)
" >> log
```

**方案**: 将 Python 代码中的所有路径引用改为单引号 `'$VAR'`。bash 在 `"..."` 内部仍然展开 `$VAR`，但 `'` 不会截断 bash 字符串。

**影响**: 此 bug 从 Loop 系统第一天就存在，导致 100% 的交付都需要手动注入 token 数据，直到 2026-06-27 13:50 修复。

---

### 23. exit code 首次归零

**现象**: 2026-06-27 12:49 的交付 (#351) 首次 exit code 0。

**根因**: 之前的 `set -e` 在 post-processing 遇到非关键错误即退出。

**方案**: 整个 post-processing 段设为 `set +e`。已验证：后续 2 次交付 exit code 0。

---

## 六、数据完整性问题

### 18. 代码变更统计缺失

**现象**: 大部分交付记录的变更列显示 `—`，无法追踪实际代码改动量。

**方案**: 从 git merge commit 的 `--stat` 提取

```bash
git diff <merge_sha>~1..<merge_sha> --stat
```

解析 `X files changed, Y insertions(+), Z deletions(-)`，写回 README 变更列。

---

### 19. 交付时间缺失

**现象**: 不知道每次交付的具体完成时间。

**方案**: 从 git merge commit 的 `%ci` 提取 commit date。

---

### 20. 成本列缺失

**现象**: 早期交付只记录 token 数，无成本信息。

**方案**:
- 新增"成本"列（DeepSeek ¥ 计价）
- 新增"完成时间"列
- `count-deliveries.sh` 统计总成本

---

## 七、汇总

### 问题总览

| # | 类别 | 问题 | 状态 | 方案摘要 |
|---|------|------|:--:|------|
| 1 | 架构 | 上下文累积→质量衰减 | ✅ | `claude -p` 独立 session |
| 2 | 架构 | 多 loop 并发冲突 | ✅ | 错开调度 + PID 锁 |
| 3 | 架构 | 子进程被父退出连带 kill | ✅ | `nohup` |
| 4 | 架构 | claude -p 初始化卡死 | ⚠️ | `timeout 1200` |
| 5 | 架构 | `&` 双重后台陷阱 | ✅ | 去掉 `&` |
| 6 | 门禁 | CI 门禁被跳过 | ✅ | prompt 显式 8 gate 表 |
| 7 | 门禁 | Docker CI 预存失败 | ✅ | `cp gc gc_cli/bin/` |
| 8 | 门禁 | macOS CI dyld 崩溃 | ⚠️ | 识别为预存，不阻塞 |
| 9 | 门禁 | CI 修复未验证 | ⚠️ | 等待新 PR 触发 |
| 10 | Token | 子进程 Token 无法获取 | ✅ | stream-json → post-process |
| 11 | Token | Anthropic 计价不匹配 | ✅ | DeepSeek 费率 |
| 12 | Token | `grep -oP` 兼容性 | ✅ | Python regex + JSONL 直搜 |
| 13 | Token | Token 消耗持续增长 | ✅ | prompt 精简 + docs-only 跳审 |
| 14 | Token | `"$TOKEN_FILE"` bash 截断 | ✅ | Python 单引号 `'$VAR'` |
| 15 | 交付 | Post-process 偶发失败 | ✅ | `set +e` best-effort |
| 16 | 交付 | 文件随 worktree 丢失 | ⚠️ | prompt 要求补写 |
| 17 | 交付 | README 表格被 GitCode 截断 | ⚠️ | 待拆为独立 Token 表 |
| 18 | 交付 | 统计需手动更新 | ✅ | `count-deliveries.sh` |
| 19 | 交付 | exit code 非零 | ✅ | `set +e` 已验证 |
| 20 | PR | 孤儿 PR | ✅ | prompt 优先 rescue |
| 21 | PR | 子进程违规改主目录 | ✅ | prompt 强调 worktree |
| 22 | 数据 | 代码变更统计缺失 | ✅ | git merge diff 回填 |
| 23 | 数据 | 交付时间缺失 | ✅ | git merge %ci |
| 24 | 数据 | 成本列缺失 | ✅ | DeepSeek ¥ 列 + 总成本统计 |
| 25 | Token | Token 注入 100% 失败 | ✅ | bash 引号 → Python 单引号 |

### 运行指标

| 指标 | 值 |
|------|-----|
| 总交付 | 55 issues |
| 已合并 | 44 |
| 含代码改动 | 41 |
| Token 总消耗 | 2.1M |
| 总成本 (DeepSeek) | ¥9.35 |
| 今日交付 (6/27) | 15 |
| 平均门禁 | 6.7/8 |
| 平均成本/交付 | ¥0.62 |
| 平均 Token/交付 | 144k |
| exit code 0 率 | 100% (3 consecutive) |
| 精简 prompt 后 Token | **70k** (from 369k) |

### 仍开放

| 风险 | 缓解 | 优先级 |
|------|------|:--:|
| claude -p 初始化卡死 | timeout 1200 | 低 |
| 文件随 worktree 丢失 | prompt 补写 | 低 |
| GitCode 表格截断 | 独立 Token 表 | 低 |
| macOS CI dyld | 标记预存 | 低 |

---

## 八、关键文件索引

| 文件 | 作用 |
|------|------|
| `.loop/scripts/full-flow-run.sh` | 编排器: PID 锁 + nohup + timeout + post-process |
| `.loop/prompts/full-flow-subprocess.md` | 子进程 prompt: 孤儿 PR rescue + 8 gate + CI + token |
| `scripts/count-deliveries.sh` | 自动刷新 README 统计 |
| `.loop/deliveries/README.md` | 交付汇总表 (11 列) |
| `.loop/deliveries/issue-N.md` | 单 issue 交付记录 |
| `.loop/registry/active.yaml` | 活跃 loop 注册 (runtime) |
| `.loop/history/*.jsonl` | stream-json 原始数据 (runtime) |
| `.github/workflows/ci.yml` | GitHub Actions CI 定义 |
