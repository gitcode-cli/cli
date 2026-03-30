# 共享 skill: pr-reviewer

## 目标

定义 gitcode-cli 项目中 PR 审查的共享场景和最低约束。

## 统一约束

- 审查流程以 `spec/workflows/review-workflow.md` 为准
- 质量门禁以 `spec/foundations/code-quality-gates.md` 为准
- 命令行为以 `gc pr review` 当前真实能力为准
- 平台不支持的动作必须明确说明，不能伪装成已支持

## 适配层说明

- Claude 适配：`.claude/skills/pr-reviewer/`
- Codex 适配：`.codex/skills/pr-reviewer/`
