---
name: gc-repo
description: GitCode CLI repository operations — create, list, view, fork, clone, sync, and manage branches.
---

# gc-repo

使用 `gc` 完成 GitCode 仓库相关操作。

## 触发场景

- 查看仓库
- 创建仓库
- fork 仓库
- 删除仓库
- 列出仓库
- 查看仓库或文件提交历史

## 仓库参数格式

常见格式：

```text
owner/repo
https://gitcode.com/owner/repo
git@gitcode.com:owner/repo.git
```

## 常用命令

```bash
# 查看仓库
gc repo view owner/repo

# 当前 Git 仓库里直接查看
gc repo view

# 创建仓库
gc repo create my-repo --public

# 列出仓库
gc repo list
gc repo list --owner my-org

# 查看提交历史
gc repo log -R owner/repo
gc repo log -R owner/repo --file README.md --branch main --json

# fork 仓库
gc repo fork owner/repo
gc repo fork owner/repo --clone

# 删除仓库
gc repo delete owner/repo
```

## 使用约束

- 删除仓库是危险操作，必须明确确认目标仓库
- `repo view` 和 `repo log` 在 Git 仓库内可尝试自动识别当前 remote
- `repo log --file ... --branch ... --json` 适合脚本或 AI 追踪具体文件在指定分支上的提交历史
- fork 或 delete 前，先确认账号权限和目标 host
- Windows PowerShell 中可将示例里的 `gc` 改为 `gitcode`，避免 `gc` 被内置 `Get-Content` 别名覆盖
