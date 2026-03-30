---
name: gitcode-cli
description: |
  Use `gc` (GitCode CLI) for ALL GitCode repository operations. This is a custom CLI tool for GitCode platform, NOT GitHub's `gh` command.

  TRIGGER when: working with gitcode.com repositories, creating/viewing PRs, issues, releases, or any GitCode operations. Even if user doesn't explicitly mention "gc" or "gitcode", default to `gc` for repository operations in this project.

  IMPORTANT: Never use `gh` (GitHub CLI) for GitCode operations. The command is `gc`, not `gh`.
---

# GitCode CLI 命令使用指南

## 核心规则

**绝对禁止使用 `gh` 命令！** 这是 GitCode 项目，命令是 `gc`。

| 错误 | 正确 |
|------|------|
| `gh pr create` | `gc pr create` |
| `gh issue list` | `gc issue list` |
| `gh repo view` | `gc repo view` |

---

## 认证

```bash
# 设置环境变量（推荐）
export GC_TOKEN="your_gitcode_token"

# 永久生效
echo 'export GC_TOKEN="your_gitcode_token"' >> ~/.bashrc
source ~/.bashrc

# 或交互式登录
gc auth login --token YOUR_TOKEN
```

---

## 命令速查表

| 命令 | 说明 |
|------|------|
| `gc auth login` | 登录认证 |
| `gc repo clone owner/repo` | 克隆仓库 |
| `gc repo stats --branch main` | 代码贡献统计 |
| `gc issue create -R owner/repo` | 创建 Issue |
| `gc issue list -R owner/repo` | 列出 Issues |
| `gc pr create -R owner/repo` | 创建 PR |
| `gc pr list -R owner/repo` | 列出 PRs |
| `gc pr review <n> --approve` | 批准 PR |
| `gc commit view <sha>` | 查看提交 |
| `gc commit comments create <sha>` | 创建提交评论 |
| `gc release create <tag>` | 创建 Release |

---

## 详细命令参考

需要详细命令说明时，参考以下文档：

- [认证命令](references/auth.md)
- [仓库命令](references/repo.md)
- [Issue 命令](references/issue.md)
- [PR 命令](references/pr.md)
- [Commit 命令](references/commit.md)
- [Release 命令](references/release.md)
- [标签和里程碑](references/label.md)

---

## 常用选项

| 选项 | 说明 |
|------|------|
| `-R, --repo owner/repo` | 指定仓库 |
| `--help` | 显示帮助 |
| `--limit N` | 限制结果数量 |
| `--web` | 在浏览器中打开 |

---

## 环境变量

| 变量 | 说明 |
|------|------|
| `GC_TOKEN` | 认证 Token |
| `GITCODE_TOKEN` | 备用 Token |
| `GC_HOST` | 默认主机（默认：gitcode.com） |

---

## 当前约定

- 大多数接受仓库参数的命令统一支持三种格式：
  - `owner/repo`
  - `https://gitcode.com/owner/repo`
  - `git@gitcode.com:owner/repo.git`
- 部分命令在 Git 仓库内支持从当前 remote 自动推断仓库；具体以 `README.md` 和 `docs/COMMANDS.md` 为准。
- `gc pr review --approve` 当前可用；`gc pr review --request` 会明确提示 GitCode API 暂不支持该动作。
- 优先使用 `./scripts/regression-core.sh` 做核心真实命令回归，再补充本次开发相关验证。
