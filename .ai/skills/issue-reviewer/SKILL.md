# 共享 skill: issue-reviewer

## 目标

定义 gitcode-cli 项目中 issue 分析、评论和流程推进的共享场景。

## 统一约束

- issue 生命周期以 `spec/workflows/issue-workflow.md` 为准
- 问题修复前必须先验证问题仍然存在
- issue comment 应反映真实状态，不得提前宣布已完成

## 适配层说明

- Claude 适配：`.claude/skills/issue-reviewer/`
- Codex 适配：`.codex/skills/issue-reviewer/`
