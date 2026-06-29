# 全流程交付演示优化方案

## 0. 现状问题

| # | 问题 | 严重度 |
|---|------|:--:|
| P1 | 每个 Phase 没有强制输出结构化总结框 — 10 个 Phase 只有 3 个输出了 | 🔴 |
| P2 | Issue/PR/comment 链接散落在工具返回值中，未集中呈现 | 🔴 |
| P3 | 跨 Issue 边界模糊，一次触发处理了 2 个 Issue | 🔴 |
| P4 | 无主控面板汇总所有交付指标 | 🟡 |
| P5 | 文档归档路径与 Issue comment 的对应关系不透明 | 🟡 |
| P6 | 完成后无"证据回放"能力 — 无法快速回顾全流程 | 🟡 |

---

## 1. 强制 Phase 输出模板

每个 Phase 结束时 AI **必须**输出以下结构化框，不得跳过任何 Phase：

```
┌─ Phase N/10: <阶段名> ───────────────────────────────────────┐
│ Issue:  #<N>  <标题前40字>                                    │
│ 状态:   ✅ 完成 / ⚠️ 阻塞 / ❌ 失败                           │
│ 产物:   <具体产出物，带可点击链接>                              │
│ 链接:   Issue: https://gitcode.com/gitcode-cli/cli/issues/<N> │
│         PR:    https://gitcode.com/gitcode-cli/cli/pulls/<N>  │
│         Doc:   .loop/deliveries/issue-<N>-<type>.md           │
│ 耗时:   <实际耗时>                                            │
│ 关键发现: <一句话>                                            │
│ 下一步:   Phase N+1: <名称>                                   │
└──────────────────────────────────────────────────────────────┘
```

### 各 Phase 产物清单

| Phase | 名称 | 必须输出的链接/索引 |
|:-----:|------|---------------------|
| 3 | 取 Issue | Issue URL, 标题, labels |
| 4 | Triage | 更新的 labels |
| 5 | Verified | 复现命令 + 结果 + Issue comment URL |
| 6 | 开发实现 | Gate 1-8 逐项结果表, commit SHA, PR URL |
| 7 | 安全审查 | 检查结论 + PR comment URL |
| 8 | 作者自检 | 9 项清单结论 + PR comment URL |
| 9 | 多角色评审 | 4 角色结论 + PR review URL |
| 10 | Merge | merge 确认 + 关闭状态 |

---

## 2. 主控面板 (Delivery Dashboard)

处理每个 Issue 时，维护 `.loop/deliveries/issue-<N>-dashboard.md`，实时更新：

```markdown
# Delivery Dashboard: Issue #<N>

**开始时间**: 2026-06-29 17:00  **结束时间**: 2026-06-29 18:03  **总耗时**: 63min

## 进度追踪

| Phase | 名称 | 状态 | 完成时间 | 产物链接 |
|:-----:|------|:----:|----------|----------|
| 0 | 需求分析 | ✅ | 17:05 | [analysis](.loop/deliveries/issue-N-analysis.md) · [comment](url) |
| 1 | 方案设计 | ✅ | 17:07 | [design](.loop/deliveries/issue-N-design.md) · [comment](url) |
| 2 | 开发计划 | ✅ | 17:08 | [plan](.loop/deliveries/issue-N-plan.md) · [comment](url) |
| 3 | 取 Issue | ✅ | 17:09 | [#N](https://gitcode.com/gitcode-cli/cli/issues/N) |
| 4 | Triage | ✅ | 17:10 | labels: type/bug, risk/medium |
| 5 | Verified | ✅ | 17:12 | [验证记录](comment-url) |
| 6 | 开发实现 | ✅ | 17:45 | [PR !N](https://gitcode.com/gitcode-cli/cli/pulls/N) |
| 7 | 安全审查 | ✅ | 17:48 | [审查结论](comment-url) |
| 8 | 作者自检 | ✅ | 17:52 | [自检报告](comment-url) |
| 9 | 多角色评审 | ✅ | 17:58 | [评审结论](comment-url) |
| 10 | Merge | ✅ | 18:03 | merged to main |

## Gate 结果

| # | Gate | 结果 | 证据 |
|---|------|:----:|------|
| 1 | 构建 | ✅ | `go build` passed |
| 2 | UT | ✅ | 1268 passed |
| 3 | Pre-commit | ✅ | 10/10 |
| 4 | 实际命令 | ✅ | `gc auth status` |
| 5 | CI | ⚠️ | macOS infra |
| 6 | 风险分级 | ✅ | risk/medium |

## 链接汇总

- **Issue**: https://gitcode.com/gitcode-cli/cli/issues/<N>
- **PR**: https://gitcode.com/gitcode-cli/cli/pulls/<N>
- **Commit**: <sha>
- **文档索引 comment**: <url>
- **自检 comment**: <url>
- **评审 comment**: <url>
- **CI Run**: <url>
```

---

## 3. Issue 边界强制规则

新增以下硬约束到 full-flow.md：

```
## 边界规则 (Anti-pattern 防护)

1. 每个 `/loop` 触发 **只处理 1 个 Issue**
2. Phase 10 (Merge) 完成后 **立即停止**，输出 ISSuE_NUM=<N>
3. 不得在未收到下次 cron 触发或用户显式指令的情况下开始新 Issue
4. 如果 triage 为空，输出 "当前无待处理 Issue，等待下次触发"
5. 孤儿 PR 模式 (triage 为空时找本人 open PR) 只处理 1 个
```

---

## 4. 链接集中输出规则

在每次 Phase 输出时，所有外部链接必须：

1. **可点击**: 使用完整 URL（`https://gitcode.com/gitcode-cli/cli/issues/408`）
2. **集中呈现**: 不散落在工具返回值中，由 AI 提取到总结框
3. **双路径引用**: 同时给出本地路径（`.loop/deliveries/issue-N-analysis.md`）和远端 URL（`Issue comment #177666274`）

示例：

```
### 本次交付链接

| 资源 | 链接 |
|------|------|
| Issue | [#408](https://gitcode.com/gitcode-cli/cli/issues/408) |
| PR | [!325](https://gitcode.com/gitcode-cli/cli/pulls/325) |
| 需求分析 | [本地](.loop/deliveries/issue-408-analysis.md) · [远端](https://gitcode.com/gitcode-cli/cli/issues/408#comment-177666274) |
| 方案设计 | [本地](.loop/deliveries/issue-408-design.md) · [远端](同上) |
| 开发计划 | [本地](.loop/deliveries/issue-408-plan.md) · [远端](同上) |
| 自检报告 | [comment](https://gitcode.com/gitcode-cli/cli/pulls/325#comment-aadda450) |
| CI Run | [#28356670397](https://github.com/gitcode-cli/cli/actions/runs/28356670397) |
```

---

## 5. 演示开场/收尾模板

### 开场（Phase 3: 取 Issue 时输出）

```
╔══════════════════════════════════════════════════════════════╗
║  全流程交付开始                                              ║
║  Issue:  #408 — refactor: multiple silent error ignores      ║
║  类型:   bug  |  风险: medium  |  范围: error-handling       ║
║  工作区: .claude/worktrees/issue-408-silent-errors           ║
║  时间:   2026-06-29 17:00 CST                                ║
╚══════════════════════════════════════════════════════════════╝
```

### 收尾（Phase 10: Merge 完成后输出）

```
╔══════════════════════════════════════════════════════════════╗
║  全流程交付完成                                              ║
║  Issue:  #408 → ✅ merged                                   ║
║  PR:     !325 → ✅ merged                                   ║
║  文件:   2 files, +25/-4                                     ║
║  门禁:   8/8 (CI 基础设施故障 1 项已记录)                    ║
║  耗时:   63 分钟                                             ║
║  ISSUE_NUM=408                                               ║
║                                                              ║
║  📊 交付面板: .loop/deliveries/issue-408-dashboard.md        ║
║  📋 全部交付: .loop/deliveries/issue-408.md                  ║
╚══════════════════════════════════════════════════════════════╝

等待下次 cron 触发 (2h 间隔)...
```

---

## 6. full-flow.md 修改方案

### 新增章节：`## 演示输出规范`

在现有 full-flow.md 末尾追加以下内容：

```markdown
## 演示输出规范

### 强制输出
AI 必须在每个 Phase 结束时输出 `┌─ Phase N/10 ─┐` 总结框，包含：
- 状态 (✅/⚠️/❌)
- 产物链接 (Issue URL, PR URL, comment URL, 本地文件路径)
- 关键发现
- 下一步

### 禁止跳过
以下 Phase 绝对不允许跳过总结框：
- Phase 0-2: 三文档完成时各输出一次
- Phase 5: Verified 完成时
- Phase 6: 开发实现完成时（必须附带 Gate 表）
- Phase 8: 自检完成时（必须附带 9 项清单）
- Phase 10: Merge 完成时（必须附带链接汇总表）

### 链接要求
- 所有远端 URL 使用完整 https:// 格式（确保终端可点击）
- 所有本地路径使用 .loop/deliveries/ 相对路径
- Phase 10 收尾时必须输出"本次交付链接"汇总表

### Issue 边界
- 每次 /loop 触发最多处理 1 个 Issue
- Merge 完成后立即结束，不开始新 Issue
- 孤儿 PR 模式同样只处理 1 个

### Dashboard 维护
处理开始时创建 `.loop/deliveries/issue-<N>-dashboard.md`
每完成一个 Phase 追加一行到进度表
处理完成时填充"链接汇总"章节
```

---

## 7. 演示效果对比

### Before（本次 session 实际表现）

```
Phase 0: ✅ (无总结框)
Phase 5: ✅ (1 个总结框)
Phase 6: ✅ (1 个总结框)
Phase 10: ✅ (1 个总结框)
→ 10 个 Phase 仅 3 个有结构化输出
→ Issue #408 完成后直接跳到 #405，无边界
→ 链接散落在 Bash 工具返回中
→ 无 dashboard 文件
```

### After（优化后期望表现）

```
╔═ 开场 ═══════════════════════════════════════════╗
┌─ Phase 0/10 ─┐ ┌─ Phase 1/10 ─┐ ... ┌─ Phase 10/10 ─┐
║ Issue #408                                    ║
║ PR !325  │  8/8 gates  │  63min               ║
║ 📊 dashboard  │  📋 delivery record           ║
╚═ 收尾 ═══════════════════════════════════════════╝
                        ↓
            等待 2h cron 触发
                        ↓
╔═ 开场 ═══════════════════════════════════════════╗
  (下一个 Issue)
```

---

## 8. 实施优先级

| 优先级 | 改动 | 工作量 |
|:--:|------|:--:|
| P0 | full-flow.md 追加"演示输出规范"章节 | 10 min |
| P0 | AI 在每个 Phase 结束后强制输出 `┌─ Phase N/10 ─┐` | prompt 约束 |
| P1 | 创建 dashboard 模板并在 Phase 3 初始化 | 已在 `.loop/deliveries/` 内 |
| P1 | Issue 边界规则（1 次 1 个） | 追加到 full-flow.md |
| P2 | 开场/收尾模板 | 追加到 full-flow.md |
| P2 | 链接汇总表（Phase 10 强制输出） | prompt 约束 |
| P0 | **CI 必须修复，无论是否本次引入** | 追加到 full-flow.md Gate 7 |

---

## 9. Gate 7 (CI) 结构化展示模板

CI 结果必须用结构化表格展示每个 Job 状态，不得用文字概述。

### 模板

```
┌─ Gate 7: CI ───────────────────────────────────────────────────┐
│ Run:  #<id>                                                     │
│ URL:  https://github.com/gitcode-cli/cli/actions/runs/<id>      │
│                                                                  │
│ | Job                  | 状态 | 说明                            │
│ |----------------------|:----:|---------------------------------|
│ | Test (ubuntu-latest) |  ✅  |                                 |
│ | Test (macos-latest)  |  ✅  |                                 |
│ | Test (windows-latest)|  ✅  |                                 |
│ | Lint                 |  ✅  |                                 |
│ | Docker               |  ✅  |                                 |
│ | Build (ubuntu-latest)|  ✅  |                                 |
│ | Build (macos-latest) |  ❌  | dyld: missing LC_UUID           |
│ | Build (windows-latest)| ⚠️  | cancelled (runner issue)        |
│                                                                  │
│ 通过: 6/8  失败: 1  取消: 1                                     │
│ 结论: ❌ macOS Build 失败，必须修复后再 merge                     │
│ 下一步: 分析根因 → 修复 → 重新触发 CI                             │
└──────────────────────────────────────────────────────────────────┘
```

### 规则

1. **CI 问题必须修复，无论是否本次修改引入**（来自 spec §5.3）
2. CI 有失败项时，Phase 6 不得标记为完成
3. 修复 CI 的 PR 需单独创建，关联原 Issue
4. CI 全绿后才能进入 Phase 7（安全审查）
5. 如果 CI 因环境不可达无法执行（如 GitHub 镜像仓不可达），必须在自检中明确记录原因

---

## 10. 违规审计清单

每次交付完成后，AI 必须对照此清单自检：

```
┌─ 交付合规审计 ─────────────────────────────────────────────────┐
│                                                                  │
│ □ 1. 每个 Phase 都有总结框 (10/10)                               │
│ □ 2. 每个 Phase 包含 Issue URL + PR URL                          │
│ □ 3. Gate 表逐项展示了每个 CI Job 状态（非文字概述）              │
│ □ 4. CI 失败项已修复（非跳过/标注即放行）                         │
│ □ 5. 一次触发只处理了 1 个 Issue                                 │
│ □ 6. 收尾有链接汇总表（Issue/PR/文档/CI 全部可点击）              │
│ □ 7. Dashboard 文件已创建并填充                                  │
│ □ 8. .loop/deliveries/issue-N.md 交付记录已写入                  │
│                                                                  │
│ 未通过项: <列出>                                                  │
│ 整改动作: <描述>                                                  │
└──────────────────────────────────────────────────────────────────┘
```

### 本次 session 审计结果

| # | 检查项 | 结果 | 说明 |
|---|--------|:--:|------|
| 1 | 每 Phase 总结框 | ❌ | 10 Phase 仅 3 个有框 |
| 2 | 链接跟随 Phase | ❌ | 散落在工具返回中 |
| 3 | CI Job 结构化表 | ❌ | 仅一段文字概述 |
| 4 | CI 失败即修复 | ❌ | macOS build 失败被标为"基础设施问题"放行 |
| 5 | 一次一个 Issue | ❌ | 触发了 2 个 (#408 → #405) |
| 6 | 收尾链接汇总 | ❌ | 无集中链接表 |
| 7 | Dashboard | ❌ | 未创建 |
| 8 | 交付记录 | ✅ | .loop/deliveries/issue-N.md 已写入 |
