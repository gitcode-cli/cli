---
name: gitcode-cli
description: |
  Use `gc` (GitCode CLI) for ALL GitCode repository operations. This is a custom CLI tool for GitCode platform, NOT GitHub's `gh` command.

  TRIGGER when: working with gitcode.com repositories, creating/viewing PRs, issues, releases, or any GitCode operations. Even if user doesn't explicitly mention "gc" or "gitcode", default to `gc` for repository operations in this project.

  IMPORTANT: Never use `gh` (GitHub CLI) for GitCode operations. The command is `gc`, not `gh`.
---

# GitCode CLI Skill

## 核心规则

**绝对禁止使用 `gh` 命令！** 这是 GitCode 项目，命令是 `gc`。

| 错误 | 正确 |
|------|------|
| `gh pr create` | `gc pr create` |
| `gh issue list` | `gc issue list` |
| `gh repo view` | `gc repo view` |

## 命令概览

```
gc <command> <subcommand> [flags]

Commands:
  auth        认证管理 (login, logout, status, token)
  repo        仓库操作 (clone, create, list, view, fork, delete)
  issue       Issue 管理 (create, list, view, close, reopen, comment)
  pr          PR 管理 (create, list, view, checkout, merge, close, reopen, review, diff, ready)
  label       标签管理 (create, list, delete)
  milestone   里程碑管理 (create, list, view, delete)
  release     Release 管理 (create, list, view, upload, download, delete)
  version     显示版本信息
```

## 认证

```bash
# 设置 Token（推荐）
export GC_TOKEN="your_token"

# 或登录
gc auth login --token YOUR_TOKEN

# 查看认证状态
gc auth status
```

## 常用命令示例

### 仓库操作
```bash
gc repo clone owner/repo
gc repo view owner/repo
gc repo create my-repo --public
gc repo fork owner/repo
```

### Issue 操作
```bash
gc issue create -R owner/repo --title "Bug" --body "Description"
gc issue list -R owner/repo
gc issue view 123 -R owner/repo
gc issue close 123 -R owner/repo
gc issue comment 123 --body "Comment" -R owner/repo
gc issue label 123 --add bug,enhancement -R owner/repo
```

### PR 操作
```bash
# 创建 PR（自动检测当前分支）
gc pr create -R owner/repo --title "Feature" --body "Description"

# 创建跨仓库 PR
gc pr create -R upstream/repo --fork myfork/repo --title "Feature"

# 查看/列出 PR
gc pr list -R owner/repo
gc pr view 456 -R owner/repo
gc pr view 456 --comments -R owner/repo

# PR 操作
gc pr checkout 456 -R owner/repo
gc pr merge 456 -R owner/repo
gc pr close 456 -R owner/repo
gc pr review 456 --approve -R owner/repo
gc pr review 456 --comment "Review comment" -R owner/repo
```

### Release 操作
```bash
gc release create v1.0.0 -R owner/repo --title "v1.0.0" --notes "Release notes"
gc release list -R owner/repo
gc release upload v1.0.0 file.zip -R owner/repo
gc release download v1.0.0 -R owner/repo
```

## 指定仓库

大部分命令支持 `-R, --repo` 参数指定仓库：
```bash
gc issue list -R owner/repo
gc pr create -R owner/repo --title "Title"
```

## 完整命令参考

详见项目文档：`docs/COMMANDS.md`

## 环境变量

| 变量 | 说明 |
|------|------|
| `GC_TOKEN` | 认证 Token（主要） |
| `GITCODE_TOKEN` | 备用 Token |