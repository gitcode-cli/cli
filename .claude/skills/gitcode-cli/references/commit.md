# Commit 命令

> 本文档是 Claude 参考层，不是命令行为真相源。
> 命令行为以 `docs/COMMANDS.md` 和 `spec/` 为准。

## commit view - 查看提交

```bash
# 查看提交详情
gc commit view abc123 -R owner/repo

# 显示变更文件
gc commit view abc123 -R owner/repo --show-diff

# 输出 JSON 格式
gc commit view abc123 -R owner/repo --json

# 在浏览器打开
gc commit view abc123 -R owner/repo --web
```

## commit diff - 获取提交差异

```bash
# 获取提交 diff
gc commit diff abc123 -R owner/repo
```

## commit patch - 获取提交补丁

```bash
# 获取提交 patch
gc commit patch abc123 -R owner/repo
```

## commit comments - 提交评论管理

### 创建评论

```bash
gc commit comments create abc123 --body "Nice work!" -R owner/repo
```

### 查看评论

```bash
# 查看指定评论
gc commit comments view 123 -R owner/repo
```

### 编辑评论

```bash
gc commit comments edit 123 --body "Updated comment" -R owner/repo
```

### 列出评论

```bash
# 列出仓库所有评论
gc commit comments list -R owner/repo

# 列出指定提交的评论
gc commit comments list-by-sha abc123 -R owner/repo

# 分页
gc commit comments list -R owner/repo --page 1 --per-page 50
```

---

**最后更新**: 2026-03-26
