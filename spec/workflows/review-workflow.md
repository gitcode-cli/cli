# 评审流程

本文档定义代码评审的流程和规范。

## 流程概览

```
Issue 评论 → PR 审查评论 → 检查清单 → 通过评审
```

## 1. Issue 评论

### 评论时机
- 问题修复完成时
- 需求实现完成时
- 需要补充信息时

### 评论内容
- 如何解决的问题
- 做了哪些改动
- 测试结果

### 评论命令

```bash
gc issue comment <number> --body "## 修复完成

### 问题原因
描述问题的根本原因

### 解决方案
描述如何解决的

### 测试验证
- [x] 单元测试通过
- [x] 实际命令测试通过" \
  -R gitcode-cli/cli
```

## 2. PR 审查评论

### 评论时机
- 代码变更完成后
- 合并 PR 之前

### 评论内容
- 改动内容摘要
- 测试结果
- 需要注意的问题（如有）

### 评论命令

```bash
gc pr review <number> --comment "## 审查结果

### 改动内容
- 新增 xxx 命令
- 修复 xxx 问题

### 测试结果
- [x] 单元测试通过
- [x] 实际命令测试通过

### 注意事项
（如有）" \
  -R gitcode-cli/cli
```

## 3. 审查类型

### 批准（Approve）
```bash
gc pr review <number> --approve -R gitcode-cli/cli
```

### 请求修改（Request Changes）
```bash
gc pr review <number> --request -R gitcode-cli/cli
```

说明：GitCode 当前公开 API 暂不支持“request changes”动作；该命令会返回明确错误，实际审查流程请改用 `gc pr review <number> --comment "..." -R gitcode-cli/cli` 留下修改意见。

### 普通评论（Comment）
```bash
gc pr review <number> --comment "评论内容" -R gitcode-cli/cli
```

## 4. 审查检查清单

### 代码质量
- [ ] 代码符合编码规范
- [ ] 命名清晰易懂
- [ ] 无明显性能问题
- [ ] 无安全漏洞

### 测试覆盖
- [ ] 有单元测试
- [ ] 测试覆盖率达标
- [ ] 实际命令测试通过

### 文档同步
- [ ] docs/COMMANDS.md 已更新
- [ ] README.md 已更新（如有新命令）
- [ ] .claude/skills/ 文档已更新

### 提交规范
- [ ] commit 信息符合规范
- [ ] 关联了 Issue

## 5. 审查流程示例

```bash
# 1. 查看 PR 详情
gc pr view <number> -R gitcode-cli/cli

# 2. 查看 PR 代码变更
gc pr diff <number> -R gitcode-cli/cli

# 3. 提交审查评论
gc pr review <number> --comment "## 审查结果

### 改动内容
- 新增 issue label 命令

### 测试结果
- [x] 单元测试通过
- [x] 实际命令测试通过

代码质量良好，可以合并。" \
  -R gitcode-cli/cli

# 4. 如果满意，批准 PR
gc pr review <number> --approve -R gitcode-cli/cli
```

## 检查清单

- [ ] Issue 已添加完成说明
- [ ] PR 已提交审查评论
- [ ] 代码质量已检查
- [ ] 测试已验证
- [ ] 文档已同步

---

**最后更新**: 2026-03-26
