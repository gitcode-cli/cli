# PR 流程

本文档定义 Pull Request 的完整生命周期管理流程。

## 流程概览

```
创建分支 → 开发代码 → 编写测试 → 提交代码 → 创建 PR → 合并 PR
```

## 1. 创建分支

### 分支命名规范

| 类型 | 命名格式 | 示例 |
|------|----------|------|
| BUG 修复 | `bugfix/issue-<number>` | `bugfix/issue-33` |
| 新特性 | `feature/issue-<number>` | `feature/issue-23` |
| 文档更新 | `docs/issue-<number>` | `docs/issue-5` |

### 创建步骤

```bash
# 1. 确保在 main 分支
git checkout main
git pull

# 2. 创建开发分支
git checkout -b feature/issue-23
```

## 2. 开发代码

### 开发规范
- 遵循 [编码规范](../foundations/coding-standards.md)
- 使用 [命令开发模板](../foundations/command-template.md)

### 本地测试
```bash
# 构建本地版本
go build -o ./gc ./cmd/gc

# 使用本地版本测试
./gc issue list -R owner/repo
```

## 3. 编写测试

参见 [测试流程](./test-workflow.md)

```bash
# 运行单元测试
go test ./pkg/cmd/xxx/...

# 进行实际命令测试
./gc xxx -R infra-test/gctest1
```

## 4. 提交代码

### 提交规范

使用 Conventional Commits：
- `feat:` 新功能
- `fix:` Bug 修复
- `docs:` 文档更新
- `test:` 测试相关
- `refactor:` 重构

### 提交步骤

```bash
# 1. 暂存更改
git add <files>

# 2. 提交（关联 Issue）
git commit -m "feat: add issue label command

Closes #23

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"

# 3. 推送分支
git push -u origin feature/issue-23
```

### 提交要求
- 单次提交不超过 800 行
- 完成一个功能点立即提交
- 提交后立即推送

## 5. 创建 PR

### 创建命令

```bash
gc pr create --title "feat: add issue label command" \
  --body "## 变更内容

- 新增 gc issue label 命令
- 支持添加、移除、列出标签

## 测试结果

- [x] 单元测试通过
- [x] 实际命令测试通过

Closes #23" \
  --base main \
  -R gitcode-cli/cli
```

### PR 描述要求
- 说明变更内容
- 关联 Issue（`Closes #xx` 或 `Fixes #xx`）
- 列出测试结果

## 6. 合并 PR

### 合并前检查
- [ ] 单元测试通过
- [ ] 实际命令测试通过
- [ ] Issue 已添加完成说明
- [ ] PR 已提交审查评论

### 合并命令

```bash
gc pr merge <number> -R gitcode-cli/cli
```

### 合并后操作

```bash
# 切换回 main 并拉取
git checkout main
git pull
```

## 完整流程示例

```bash
# 1. 创建分支
git checkout main && git pull
git checkout -b feature/issue-23

# 2. 开发代码
# ... 编写代码 ...

# 3. 编写测试
# ... 创建 xxx_test.go ...

# 4. 运行测试
go test ./pkg/cmd/issue/label/...

# 5. 实际命令测试
./gc issue label 1 --add bug -R infra-test/gctest1

# 6. 提交代码
git add .
git commit -m "feat(issue): add label command

Closes #23

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"

# 7. 推送分支
git push -u origin feature/issue-23

# 8. 创建 PR
gc pr create --title "feat: add issue label command" \
  --body "Closes #23" --base main -R gitcode-cli/cli

# 9. 添加 Issue 评论
gc issue comment 23 --body "已完成实现" -R gitcode-cli/cli

# 10. 合并 PR
gc pr merge <pr_number> -R gitcode-cli/cli

# 11. 拉取最新代码
git checkout main && git pull
```

## 检查清单

- [ ] 分支已创建
- [ ] 代码已开发
- [ ] 测试已编写
- [ ] 单元测试通过
- [ ] 实际命令测试通过
- [ ] 代码已提交（关联 Issue）
- [ ] 分支已推送
- [ ] PR 已创建
- [ ] Issue 已评论
- [ ] PR 已合并

---

**最后更新**: 2026-03-26
