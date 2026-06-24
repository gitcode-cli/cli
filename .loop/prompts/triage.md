# /goal: Issue Triage → Verified

## Prompt

```
/goal until issue #<ISSUE_NUMBER> is triaged AND verified:
  - 判断类型（bug/feature/docs/refactor），补标签
  - 复现问题或确认需求有效
  - 在 issue comment 中留下结构化验证记录
  - 标签更新为 status/verified
完成后更新 .loop/deliveries/issue-<ISSUE_NUMBER>.md 记录状态流转
```

## 替换参数

- `<ISSUE_NUMBER>`: 目标 issue 编号

## 评估器检查点

- Issue 标签是否包含 type/* 和 scope/*
- Issue comment 是否包含验证记录
- 标签是否包含 status/verified

## 验证记录模板

```
- 当前版本或分支: <branch/commit>
- 验证时间: <timestamp>
- 复现命令: <commands>
- 实际结果: <output>
- 预期结果: <expected>
- 结论: 问题仍存在/已修复/无需处理
```

## .loop/ 更新

完成后写 `.loop/deliveries/issue-<N>.md`:
```markdown
# Delivery Record: Issue #<N>
## State Transitions
| From | To | When | Evidence |
|------|----|------|----------|
| status/triage | status/verified | <timestamp> | <comment_url> |
```
