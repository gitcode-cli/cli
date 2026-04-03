# 共享 skill: pr-reviewer

## 目标

定义 gitcode-cli 项目中 PR 审查的共享场景和最低约束。

## 统一约束

- 审查流程以 `spec/workflows/review-workflow.md` 为准
- 质量门禁以 `spec/foundations/code-quality-gates.md` 为准
- 不同信息类型的事实边界以 `spec/governance/source-of-truth-matrix.md` 为准
- 命令行为以 `gc pr review` 当前真实能力为准
- 平台不支持的动作必须明确说明，不能伪装成已支持
- 作者自检与独立评审必须分离
- PR 必须按 `status/draft -> status/self-checked -> status/ready-for-review -> status/approved -> status/merged` 推进

## 最低记录要求

- 作者必须先在 PR 中补 `作者自检`
- 独立评审必须单独给出 `评审结论`
- 只有独立评审才能把 PR 视为 `approved`

## PR 评论模板

完整模板集合见：`docs/AI-TEMPLATES.md`

作者自检：

```markdown
## 作者自检

- 根因或实现理由:
- 主要修改:
- 单元测试:
- 构建:
- 实际命令验证:
- 安全审查:
- 文档同步:
- 风险:
- 未覆盖项:
```

评审结论：

```markdown
## 评审结论

- 发现:
- blocker:
- 安全检查:
- 结论:
```

## 适配层说明

- Claude 适配：`.claude/skills/pr-reviewer/`
- Codex 适配：`.codex/skills/pr-reviewer/`
