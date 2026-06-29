从 status/triage 取一个 issue，推进到 merged。每次只处理一个。

## 前置
- 独立 git worktree（`.claude/worktrees/issue-N-<ts>`），用后即删
- 禁止在 main 开发、跳过验证、作者自检当独立评审

## 流程概览

```
Phase 0: 需求分析 → Phase 1: 方案设计 → Phase 2: 开发计划 →
Phase 3: 取 Issue → Phase 4: Triage → Phase 5: Verified →
Phase 6: 开发实现(8门禁) → Phase 7: 安全审查 → Phase 8: 作者自检(9项) →
Phase 9: 多角色评审 → Phase 10: Merge
```

## Phase 0: 需求分析
读到 issue 后、写任何代码前，按 `docs/superpowers/specs/analysis.md` 模板输出需求分析：
- 写入 `.loop/deliveries/issue-N-analysis.md`
- Issue comment 只贴文件路径索引（如 `## 需求分析 → .loop/deliveries/issue-N-analysis.md`）
- 在当前窗口打印交付总结

## Phase 1: 方案设计
基于需求分析，按 `docs/superpowers/specs/design.md` 模板输出方案设计：
- 写入 `.loop/deliveries/issue-N-design.md`
- Issue comment 只贴文件路径索引
- 在当前窗口打印交付总结

## Phase 2: 开发计划
基于选定方案，按 `docs/superpowers/specs/plan.md` 模板输出开发计划：
- 写入 `.loop/deliveries/issue-N-plan.md`
- Issue comment 只贴文件路径索引
- 在当前窗口打印交付总结

## Phase 3: 取 Issue
`gc issue list -R gitcode-cli/cli --state opened --label status/triage --limit 5`，选最小 scope；若 triage 空→孤儿 PR 检查

## Phase 4-5: Triage → Verified
判定 docs-only 还是代码改动，走对应分支。补验证记录到 Issue。

## Phase 6: 开发实现（8 门禁）

| # | 门禁 | docs-only | 代码改动 |
|---|------|:--:|------|
| 1 | 实现 | — | 修复 |
| 2 | 测试 | 跳过 | go test ./... 全通过 |
| 3 | 构建 | 跳过 | go build 成功 |
| 4 | UT | 跳过 | 全通过 |
| 5 | Pre-commit | 必须 | 必须 |
| 6 | 命令验证 | 跳过 | infra-test/* 至少一条 |
| 7 | CI | 跳过 | gh CLI 触发 GitHub Actions，等待全绿，PR 附 run URL。CI 问题必须修复，无论是否本次修改引入 |
| 8 | 风险分级 | 必须 | classify-change-risk.py |

## Phase 7-10: 安全审查 → 自检 → 评审 → Merge
状态机: triage→verified→in-progress→draft→self-checked→ready→approved→merged

- risk/low 自动合，risk/high 暂停
- docs-only 被 classify-change-risk.py 误判为 high 时（文件名匹配，非内容分析），覆盖为实际风险级别继续推进，不阻塞
- docs-only 跳过评审；其余路径必须多角色独立评审

## 每步交付总结

每个 Phase 完成后，在当前窗口打印：

```
┌─ Phase N/10: <阶段名> ────────────────────┐
│ 状态: ✅/❌/⚠️                               │
│ 产物: <comment/label/PR/commit>             │
│ 关键发现: <一句话>                           │
│ 下一步: Phase N+1: <名称>                   │
└────────────────────────────────────────────┘
```

## 证据
- Issue: 验证记录 + 自检 9 项（含 CI URL）+ 三文档索引（analysis/design/plan 文件路径）
- PR: 评审结论 + CI URL + gate 表
- CI 未跑写 ✅ 算违规；CI 问题必须修复，无论是否本次修改引入

## 交付
创建 `.loop/deliveries/issue-N.md`（含 Phase 0-2 文档的 Issue comment ID 引用），更新 README。末尾输出 `ISSUE_NUM=<N>`。

## 孤儿 PR（仅 triage 为空时）
`gc pr list --state open --json`，找本人非 draft PR→完整读评论→对照 spec/workflows/development-workflow.md §5.3 补缺失→合并。
