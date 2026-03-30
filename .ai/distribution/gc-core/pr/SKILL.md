# gc-pr

使用 `gc` 完成 GitCode Pull Request 相关操作。

## 触发场景

- 创建 PR
- 列出或查看 PR
- 编辑 PR
- 检出 PR
- 合并 PR

## 常用命令

```bash
# 创建 PR
gc pr create -R owner/repo --title "New feature" --base main

# 指定 head 分支
gc pr create -R owner/repo --head feature-branch --title "Feature"

# 跨仓库 PR
gc pr create -R upstream/repo --fork myfork/repo --title "Feature"

# 列出 PR
gc pr list -R owner/repo --state open

# 查看 PR
gc pr view 123 -R owner/repo

# 编辑 PR
gc pr edit 123 -R owner/repo --title "New title"

# 检出 PR
gc pr checkout 123 -R owner/repo

# 合并 PR
gc pr merge 123 -R owner/repo
gc pr merge 123 -R owner/repo --squash
```

## 使用约束

- `pr create` 在很多场景下仍建议显式传 `-R`
- 未传 `--head` 时，CLI 会尝试使用当前分支
- 合并前应确认目标仓库策略和分支保护要求
