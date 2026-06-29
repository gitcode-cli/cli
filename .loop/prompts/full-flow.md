# /loop 全流程交付

## Prompt

~~~
/loop 2h 在 git worktree 中从 status/triage 取一个 issue，推进到 merged。每次只处理一个。

## 流程概览
Phase 0: 需求分析 → Phase 1: 方案设计 → Phase 2: 开发计划 →
Phase 3: 取 Issue → Phase 4: Triage → Phase 5: Verified →
Phase 6: 开发实现(8门禁) → Phase 7: 安全审查 → Phase 8: 作者自检(9项) →
Phase 9: 多角色评审 → Phase 10: Merge

## Phase 0-2: 开发前三文档
按 docs/superpowers/specs/ 下的模板输出，写入 Issue comment：
- analysis.md — 需求分析（问题定义/影响范围/根因/成功标准/约束）
- design.md — 方案设计（≥2 方案对比/推荐/关键决策/风险缓解）
- plan.md — 开发计划（文件清单/执行顺序/测试策略）

## Phase 3-10
状态机: triage→verified→in-progress→draft→self-checked→ready→approved→merged
risk/low 自动合，risk/high 暂停。
8 门禁按 spec/workflows/development-workflow.md §5.3 执行。

## 每步交付总结
每个 Phase 完成后打印:
┌─ Phase N/10: <阶段名> ────────────────────┐
│ 状态: ✅/❌/⚠️                               │
│ 产物: <comment/label/PR/commit>             │
│ 关键发现: <一句话>                           │
│ 下一步: Phase N+1: <名称>                   │
└────────────────────────────────────────────┘

## 证据
- Issue: 需求分析 + 方案设计 + 开发计划 + 验证记录 + 自检 9 项(含 CI URL)
- PR: 多角色评审结论 + CI URL + gate 表

## 交付
.loop/deliveries/issue-N.md（含 Phase 0-2 的 Issue comment ID 引用）
更新 README/INDEX/lessons。末尾输出 ISSUE_NUM=<N>。

## 孤儿 PR（triage 为空时）
gc pr list --state open --json，找本人非 draft PR→补缺失→合并。
~~~

## 设置方式

1. 替换间隔（推荐 30m）
2. 确认 `docs/superpowers/specs/` 下 analysis/design/plan 模板已是最新
3. 确认 `spec/workflows/development-workflow.md` §5.3 门禁表已是最新版本
4. 启动后写入 `.loop/registry/active.yaml`

## 门禁清单

AI 必须逐项执行，每项完成后留下证据：

| # | 门禁 | docs-only | 代码改动 | 证据 |
|---|------|:--:|:--:|------|
| 0a | 需求分析 | 必须 | 必须 | Issue comment（`docs/superpowers/specs/analysis.md` 模板） |
| 0b | 方案设计 | 必须 | 必须 | Issue comment（`docs/superpowers/specs/design.md` 模板） |
| 0c | 开发计划 | 必须 | 必须 | Issue comment（`docs/superpowers/specs/plan.md` 模板） |
| 1 | 验证 | 必须 | 必须 | Issue comment 中的复现记录 |
| 2 | 开发 | — | 必须 | 非 main 分支 + commits |
| 3 | 构建 | 跳过 | 必须 | `go build -o ./gc ./cmd/gc` 成功 |
| 4 | UT | 跳过 | 必须 | `go test ./...` 全部通过 |
| 5 | Pre-commit | 必须 | 必须 | 所有 hooks 通过 |
| 6 | 实际命令 | 跳过 | 必须 | `./gc <cmd> -R infra-test/gctest1` |
| 7 | CI | 跳过 | 必须 | `gh run list` 全绿 + run URL。CI 问题必须修复，无论是否本次修改引入 |
| 8 | 风险分级 | 必须 | 必须 | `scripts/classify-change-risk.py` |
| + | 多角色评审 | 2 角色 | 4 角色 | Agent 结论汇总到 PR |
| + | 合并 | 自动 | 自动* | risk/low 自动，risk/high 暂停 |

## 每步交付总结

```
┌─ Phase N/10: <阶段名> ────────────────────┐
│ 状态: ✅/❌/⚠️                               │
│ 产物: <comment/label/PR/commit>             │
│ 关键发现: <一句话>                           │
│ 下一步: Phase N+1: <名称>                   │
└────────────────────────────────────────────┘
```

## .loop/ 维护要求

每个 issue 处理完成后，必须更新：

### deliveries/issue-<N>.md

```markdown
# Delivery Record: Issue #<N>
- Title: <title>
- Type: <bug/feature/docs>
- Status: <merged/closed>

## Design Artifacts
- 需求分析: <Issue comment URL>
- 方案设计: <Issue comment URL>
- 开发计划: <Issue comment URL>

## State Transitions
| From | To | When | Evidence |

## Key Artifacts
- PR: !<N>, CI Run: <URL>, Commits: <sha>
```

### memory/INDEX.md
- 更新"当前活跃 Issue"列表
- 如有新发现，追加到"经验教训"
- 更新"待办"项

### 会话结束
- 写 `memory/YYYY-MM-DD-session.md` 记录本会话摘要
- 标记 `registry/active.yaml` 中的 loop 为 lost
