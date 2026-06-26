---
name: gc-pr
description: GitCode CLI pull request operations — create, list, view, edit, merge, checkout, and comment on PRs.
---

# gc-pr

使用 `gc` 完成 GitCode Pull Request 相关操作。

## 触发场景

- 创建 PR
- 列出或查看 PR
- 编辑 PR
- 管理 PR 标签
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
gc pr list -R owner/repo --paginate --per-page 100 --json
gc pr list -R owner/repo --commit-message "fix login" --json

# 查看 PR
gc pr view 123 -R owner/repo --json

# 编辑 PR
gc pr edit 123 -R owner/repo --title "New title"

# 管理 PR 标签
gc pr label 123 --add bug,enhancement -R owner/repo
gc pr label 123 --remove bug -R owner/repo
gc pr label 123 --list -R owner/repo --json

# 检出 PR
gc pr checkout 123 -R owner/repo

# 合并 PR
gc pr merge 123 -R owner/repo
gc pr merge 123 -R owner/repo --squash
```

## 使用约束

- `pr create` 在很多场景下仍建议显式传 `-R`
- 未传 `--head` 时，CLI 会尝试使用当前分支
- `pr list --paginate` 用于跨页扫描；`--commit-message` 会读取候选 PR 的提交列表并按提交信息过滤
- `pr view --json` 会包含 `body`、`description`、`merged_at`，并在远端详情统计为 0 时尽量通过 files/commits API 补齐统计
- `pr label` 支持 `--add`、`--remove`、`--list` 和 `--json`；`--add` 接受逗号分隔的标签名列表
- 合并前应确认目标仓库策略和分支保护要求
- Windows PowerShell 中可将示例里的 `gc` 改为 `gitcode`；从 stdin 传中文/非 ASCII 正文到 `--body-file -` 时，优先使用 UTF-8 文件，直接管道前先设置 `$OutputEncoding = [System.Text.UTF8Encoding]::new($false)`
- `pr create --json` 若 warning 提示远端 body 未返回，不要把本地提交正文当作远端事实；使用 `gitcode pr view <number> -R owner/repo --json` 再核验
