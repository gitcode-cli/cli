# 全流程交付

## 任务

在 git worktree 中从 status/triage 取一个 issue，推进到 merged。每次只处理一个。

## 流程概览
Phase 0: 需求分析 → Phase 1: 方案设计 → Phase 2: 开发计划 →
Phase 3: 取 Issue → Phase 4: Triage → Phase 5: Verified →
Phase 6: 开发实现(8门禁) → Phase 7: 安全审查 → Phase 8: 作者自检(9项) →
Phase 9: 多角色评审 → Phase 10: Merge

## Phase 0-2: 开发前三文档
按 docs/superpowers/specs/ 下的模板输出，写入 .loop/deliveries/ 目录：
- analysis.md → .loop/deliveries/issue-N-analysis.md
- design.md → .loop/deliveries/issue-N-design.md
- plan.md → .loop/deliveries/issue-N-plan.md
Issue comment 只贴文件路径索引，不贴全文。

## Phase 3-10
状态机: triage→verified→in-progress→draft→self-checked→ready→approved→merged
risk/low 自动合，risk/high 暂停。
8 门禁按 spec/workflows/development-workflow.md §5.3 执行。

## 证据
- Issue: 验证记录 + 自检 9 项(含 CI URL) + 三文档索引
- PR: 多角色评审结论 + CI URL + gate 表

## 交付
.loop/deliveries/issue-N.md（含 Phase 0-2 的 Issue comment ID 引用）
更新 README/INDEX/lessons。末尾输出 ISSUE_NUM=<N>。

## 孤儿 PR（triage 为空时）
gc pr list --state open --json，找本人非 draft PR→补缺失→合并。

## 门禁清单

AI 必须逐项执行 `spec/workflows/development-workflow.md` §5.3，每项完成后留下证据：

| # | 门禁 | docs-only | 代码改动 | 证据 |
|---|------|:--:|:--:|------|
| 0a | 需求分析 | 必须 | 必须 | `.loop/deliveries/issue-N-analysis.md` |
| 0b | 方案设计 | 必须 | 必须 | `.loop/deliveries/issue-N-design.md` |
| 0c | 开发计划 | 必须 | 必须 | `.loop/deliveries/issue-N-plan.md` |
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

## 演示输出规范

### 强制 Phase 总结框

AI 必须在 **每个 Phase 结束** 时输出结构化总结框，不得跳过任何 Phase：

```
┌─ Phase N/10: <阶段名> ───────────────────────────────────────┐
│ Issue:  #<N>  <标题>                                         │
│ 状态:   ✅/⚠️/❌                                              │
│ 产物:   <具体产出物>                                          │
│ 链接:   Issue: https://gitcode.com/gitcode-cli/cli/issues/<N> │
│         PR:    https://gitcode.com/gitcode-cli/cli/pulls/<N>  │
│         Doc:   .loop/deliveries/issue-<N>-<type>.md           │
│ 关键发现: <一句话>                                            │
│ 下一步: Phase N+1: <名称>                                     │
└──────────────────────────────────────────────────────────────┘
```

### Gate 7 (CI) 强制要求

**CI 问题必须修复，无论是否本次修改引入。**

CI 结果必须用结构化表格展示每个 Job：

```
┌─ Gate 7: CI ───────────────────────────────────────────────────┐
│ Run:  #<id>  URL: https://github.com/gitcode-cli/cli/actions/runs/<id> │
│ | Job                  | 状态 | 说明                           │
│ | Build (macOS)        |  ❌  | dyld: missing LC_UUID          │
│ | Build (ubuntu)       |  ✅  |                                │
│ 通过: N/8  失败: M                                              │
│ 结论: ❌ 有失败项，必须修复后再 merge                            │
└──────────────────────────────────────────────────────────────────┘
```

- CI 有失败项时，Phase 6 不得标记为完成
- 修复 CI 的 PR 需单独创建并关联原 Issue
- CI 因环境不可达无法执行时（如 GitHub 镜像仓不可达），须在自检中明确记录

### Issue 边界

- 每次 /loop 触发 **只处理 1 个 Issue**
- Phase 10 Merge 完成后 **立即停止**
- 不得在未收到下次触发或用户显式指令的情况下开始新 Issue
- 孤儿 PR 模式同样只处理 1 个

### 开场模板

```
╔══════════════════════════════════════════════════════════════╗
║  全流程交付开始                                              ║
║  Issue:  #<N> — <标题>                                      ║
║  类型:   <type>  |  风险: <risk>  |  范围: <scope>          ║
║  工作区: .claude/worktrees/issue-<N>-<slug>                 ║
║  时间:   <timestamp>                                        ║
╚══════════════════════════════════════════════════════════════╝
```

### 收尾模板

Phase 10 完成后输出：

```
╔══════════════════════════════════════════════════════════════╗
║  全流程交付完成                                              ║
║  Issue:  #<N> → ✅ merged                                   ║
║  PR:     !<N> → ✅ merged                                   ║
║  文件:   N files, +N/-M                                      ║
║  门禁:   8/8 (例外项: <列出并说明>)                           ║
║                                                              ║
║  ## 本次交付链接                                              ║
║  | 资源     | 链接                                          | ║
║  | Issue    | https://gitcode.com/.../issues/<N>            | ║
║  | PR       | https://gitcode.com/.../pulls/<N>             | ║
║  | 需求分析 | .loop/deliveries/issue-<N>-analysis.md        | ║
║  | CI Run   | https://github.com/.../actions/runs/<id>      | ║
║                                                              ║
║  ISSUE_NUM=<N>                                               ║
║  📊 面板: .loop/deliveries/issue-<N>-dashboard.md            ║
╚══════════════════════════════════════════════════════════════╝
```

### 交付合规审计

每次交付完成后自检：

```
□ 1. 每个 Phase 都有总结框 (10/10)
□ 2. 每个 Phase 包含 Issue URL + PR URL
□ 3. Gate 7 CI 用结构化表格展示了每个 Job
□ 4. CI 失败项已修复，非跳过
□ 5. 一次触发只处理了 1 个 Issue
□ 6. 收尾有链接汇总表
□ 7. Dashboard 文件已创建
□ 8. .loop/deliveries/issue-N.md 已写入

未通过项: <列出>
```
