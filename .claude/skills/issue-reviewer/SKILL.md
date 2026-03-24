---
name: issue-reviewer
description: |
  Review and analyze GitCode CLI project issues. Assess validity, add comments, and apply labels.

  TRIGGER when: user asks to review issues, analyze issues, triage issues, manage issue backlog,
  evaluate issue acceptance, or check issue quality. Also trigger for phrases like "评审Issue",
  "分析Issue", "Issue评审", "帮我看看这个Issue".
---

# Issue Reviewer

自动化评审 GitCode CLI 项目的 Issue，给出分析意见和标签建议。

## 评审流程

### 1. 获取 Issue 信息

```bash
gc issue view <number> -R gitcode-cli/cli
```

### 2. 分析维度

| 维度 | 检查项 |
|------|--------|
| 需求清晰度 | 标题是否明确、描述是否完整、是否有复现步骤 |
| 技术可行性 | 是否符合项目范围、是否可实现 |
| 重复性 | 是否已有类似 Issue 或 PR |
| 优先级 | 影响范围、紧急程度、实现复杂度 |

### 3. 给出结论

根据分析结果：
- **接纳**: 给出实现建议和优先级
- **需要信息**: 说明缺少什么信息
- **重复/无效**: 说明原因并关闭

## 标签体系

| 标签 | 使用场景 |
|------|----------|
| `accepted` | 已接纳，计划实现 |
| `needs-info` | 需要更多信息 |
| `duplicate` | 重复 Issue |
| `invalid` | 无效 Issue |
| `wontfix` | 不会修复 |
| `enhancement` | 功能增强 |
| `bug` | Bug 修复 |
| `priority: high/medium/low` | 优先级 |
| `good first issue` | 适合新手 |

## 执行命令

```bash
# 添加评论
gc issue comment <number> --body "评论内容" -R gitcode-cli/cli

# 添加标签
gc issue label <number> --add label1,label2 -R gitcode-cli/cli

# 关闭 Issue
gc issue close <number> -R gitcode-cli/cli
```

## 注意事项

1. 不要轻易关闭 Issue，除非明确是重复或无效
2. 评论保持友好和专业
3. 在评论中说明决策原因
4. 检查现有代码实现，避免重复工作