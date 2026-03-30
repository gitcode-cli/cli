# gc-repo

使用 `gc` 完成 GitCode 仓库相关操作。

## 触发场景

- 查看仓库
- 创建仓库
- fork 仓库
- 删除仓库
- 列出仓库

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

# fork 仓库
gc repo fork owner/repo
gc repo fork owner/repo --clone

# 删除仓库
gc repo delete owner/repo
```

## 使用约束

- 删除仓库是危险操作，必须明确确认目标仓库
- `repo view` 在 Git 仓库内可尝试自动识别当前 remote
- fork 或 delete 前，先确认账号权限和目标 host
