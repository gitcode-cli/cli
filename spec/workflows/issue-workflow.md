# Issue 流程

本文档定义 Issue 的完整生命周期管理流程。

## 流程概览

```
发现问题 → 创建 Issue → 打标签 → 验证问题 → 开发修复 → 关闭 Issue
```

## 1. 创建 Issue

### 触发条件
- 发现 Bug
- 需要新功能
- 文档需要更新

### 创建方式

```bash
# 方式一：命令行创建
gc issue create --title "Bug: 描述问题" --body "问题描述" -R owner/repo

# 方式二：在 Web 界面创建
# https://gitcode.com/owner/repo/issues/new
```

### Issue 模板

```markdown
## 问题描述

简要描述问题或需求。

## 复现步骤（Bug 适用）

1. 执行命令 `gc xxx`
2. 观察输出
3. 发现错误

## 期望结果

描述期望的正确行为。

## 实际结果

描述实际发生的情况。

## 环境

- 版本: v0.2.x
- 操作系统: Linux/macOS/Windows
```

## 2. 打标签

Issue 创建后立即打上合适的标签：

```bash
# 打标签
gc issue label <number> --add bug -R owner/repo
```

### 标签类型

| 标签 | 使用场景 |
|------|----------|
| `bug` | 错误修复 |
| `enhancement` | 功能增强/新特性 |
| `documentation` | 文档更新 |
| `help wanted` | 需要帮助 |
| `question` | 需要讨论 |

## 3. 验证问题

**在开始修复之前，必须先验证问题是否仍存在！**

### 验证步骤

1. **用当前版本验证**
   ```bash
   # 执行 Issue 中描述的复现步骤
   ./gc xxx --option value
   ```

2. **检查时间线**
   ```bash
   # 查看 Issue 创建时间
   gc issue view <number> -R owner/repo

   # 查看相关代码提交时间
   git log --oneline -- pkg/cmd/xxx/
   ```

3. **判断是否需要修复**
   - 如果问题已修复，在 Issue 中说明并关闭
   - 如果问题仍存在，继续修复流程

## 4. 开发修复

参见 [PR 流程](./pr-workflow.md)

## 5. 关闭 Issue

### 关闭时机
- 问题已修复并合并
- 需求已实现并合并
- 问题无效或重复

### 关闭方式

```bash
# 关闭 Issue
gc issue close <number> -R owner/repo

# 添加说明后关闭
gc issue comment <number> --body "问题已在 PR #xx 中修复" -R owner/repo
gc issue close <number> -R owner/repo
```

## 完整流程示例

```bash
# 1. 创建 Issue
gc issue create --title "fix: milestone create 命令报错" \
  --body "执行 gc milestone create 返回 400 错误" \
  -R gitcode-cli/cli

# 2. 打标签
gc issue label 33 --add bug -R gitcode-cli/cli

# 3. 验证问题
./gc milestone create "test" -R gitcode-cli/cli
# 确认问题存在

# 4. 开发修复（参见 PR 流程）
# ...

# 5. 添加完成说明
gc issue comment 33 --body "## 修复完成\n\n问题原因：...\n解决方案：..." \
  -R gitcode-cli/cli

# 6. 关闭 Issue
gc issue close 33 -R gitcode-cli/cli
```

## 检查清单

- [ ] Issue 已创建
- [ ] 标签已添加
- [ ] 问题已验证
- [ ] 修复已完成
- [ ] 完成说明已添加
- [ ] Issue 已关闭

---

**最后更新**: 2026-03-26