# 共享 skill: issue-reviewer

## 目标

定义 gitcode-cli 项目中 issue 分析、评论和流程推进的共享场景。

## 统一约束

- issue 生命周期以 `spec/workflows/issue-workflow.md` 为准
- 不同信息类型的事实边界以 `spec/governance/source-of-truth-matrix.md` 为准
- 问题修复前必须先验证问题仍然存在
- issue comment 应反映真实状态，不得提前宣布已完成
- issue 必须按 `status/triage -> status/verified -> status/in-progress -> status/ready-for-review -> status/merged` 推进
- issue 在 PR 未合入主干前不得关闭，除非明确判定为 `status/closed-no-fix`
- 不得把 `issues-plan/PROGRESS.md` 当成单个 issue / PR 实时状态真相源

## 最低记录要求

- 验证阶段必须补 `验证记录`
- 开发完成后必须补 `开发进度`
- 关闭时必须写明是“已合入主干关闭”还是“closed-no-fix 关闭”

## issue comment 模板

完整模板集合见：`docs/AI-TEMPLATES.md`

验证记录：

```markdown
## 验证记录

- 当前版本或分支:
- 复现命令:
- 实际结果:
- 结论:
```

开发进度：

```markdown
## 开发进度

- 根因:
- 主要修改:
- 测试:
- 实际命令验证:
- 安全影响:
- 风险或未覆盖项:
- 关联 PR:
```

## 适配层说明

- Claude 适配：`.claude/skills/issue-reviewer/`
- Codex 适配：`.codex/skills/issue-reviewer/`
