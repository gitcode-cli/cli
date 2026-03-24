---
name: issue-reviewer
description: |
  Review and analyze open issues in GitCode CLI project. Analyze issue validity,
  add review comments, and apply appropriate labels.

  TRIGGER when: user asks to review issues, analyze issues, triage issues, or
  manage issue backlog.
---

# Issue Reviewer Skill

## 目的

自动化评审 GitCode CLI 项目的 Issue，帮助维护者进行 Issue 分类和优先级评估。

## 工作流程

### 步骤 1: 获取所有 open Issue

```bash
gc issue list -R gitcode-cli/cli --state open --limit 50
```

### 步骤 2: 逐个分析 Issue

对每个 Issue 执行：

```bash
gc issue view <number> -R gitcode-cli/cli
```

### 步骤 3: 分析维度

**检查以下方面：**

1. **需求清晰度**
   - 标题是否清晰描述问题
   - 是否有详细的问题描述
   - 是否有复现步骤（如果是 bug）
   - 是否有期望行为说明

2. **技术可行性**
   - 是否符合项目范围
   - 技术上是否可实现
   - 是否依赖未实现的功能

3. **重复性检查**
   - 是否已有类似 Issue
   - 是否已有相关 PR

4. **优先级评估**
   - 影响范围（高/中/低）
   - 紧急程度
   - 实现复杂度

### 步骤 4: 添加评论和标签

根据分析结果执行相应操作。

## 标签体系

### 状态标签

| 标签 | 使用场景 |
|------|----------|
| `accepted` | 已接纳，计划实现 |
| `needs-info` | 需要更多信息 |
| `duplicate` | 重复 Issue |
| `invalid` | 无效 Issue |
| `wontfix` | 不会修复 |
| `question` | 需要讨论 |

### 类型标签

| 标签 | 使用场景 |
|------|----------|
| `enhancement` | 功能增强 |
| `bug` | Bug 修复 |
| `documentation` | 文档相关 |
| `refactor` | 代码重构 |

### 优先级标签

| 标签 | 使用场景 |
|------|----------|
| `priority: high` | 高优先级 |
| `priority: medium` | 中优先级 |
| `priority: low` | 低优先级 |

### 社区标签

| 标签 | 使用场景 |
|------|----------|
| `good first issue` | 适合新手贡献者 |
| `help wanted` | 需要社区帮助 |

## 评论模板

### 接纳 Issue

```markdown
## Issue 分析

### 需求理解
[描述对需求的理解]

### 实现建议
[给出实现思路或建议]

### 优先级评估
**优先级**: [高/中/低]
**原因**: [说明原因]

### 后续步骤
- [ ] 创建开发分支
- [ ] 实现功能
- [ ] 编写测试

---
感谢提交 Issue！我们会尽快安排实现。
```

### 需要更多信息

```markdown
## 需要更多信息

为了更好地处理这个 Issue，请提供以下信息：

1. [具体需要的信息]
2. [其他补充信息]

建议使用以下模板完善 Issue 描述：

\`\`\`markdown
## 问题描述
[详细描述问题]

## 复现步骤
1. 步骤一
2. 步骤二
3. ...

## 期望行为
[描述期望的行为]

## 实际行为
[描述实际发生的情况]

## 环境
- 操作系统:
- 版本:
\`\`\`

---
请在 7 天内补充信息，否则 Issue 可能会被关闭。
```

### 重复 Issue

```markdown
## 重复 Issue

此 Issue 与以下 Issue 重复：
- #xxx: [Issue 标题]

请在该 Issue 下继续讨论，此 Issue 将被关闭。

---
感谢理解！
```

### 无效 Issue

```markdown
## Issue 分析

经过分析，此 Issue 被标记为无效，原因如下：

- [原因说明]

如有疑问，请提供更多信息后重新打开。

---
感谢理解！
```

## 执行命令

### 添加评论

```bash
gc issue comment <number> --body "评论内容" -R gitcode-cli/cli
```

### 添加标签

```bash
gc issue label <number> --add label1,label2 -R gitcode-cli/cli
```

### 关闭 Issue

```bash
gc issue close <number> -R gitcode-cli/cli
```

## 示例：评审单个 Issue

```
评审 Issue #15
```

执行流程：
1. `gc issue view 15 -R gitcode-cli/cli`
2. 分析 Issue 内容
3. 添加评论：`gc issue comment 15 --body "..." -R gitcode-cli/cli`
4. 添加标签：`gc issue label 15 --add accepted,enhancement -R gitcode-cli/cli`

## 示例：评审所有 Issue

```
评审所有 open Issue
```

执行流程：
1. `gc issue list -R gitcode-cli/cli --state open`
2. 遍历每个 Issue 执行分析
3. 为每个 Issue 添加评论和标签

## 注意事项

1. **谨慎关闭**: 不要轻易关闭 Issue，除非明确是重复或无效
2. **友好态度**: 评论保持友好和专业
3. **及时响应**: 尽快评审新提交的 Issue
4. **记录决策**: 在评论中说明决策原因