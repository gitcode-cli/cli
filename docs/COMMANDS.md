# GitCode CLI 命令使用指南

本文档提供 `gc` 命令行工具所有命令的实际使用示例。

## 前置准备

### 认证

```bash
# 方式一：设置环境变量（推荐）
export GC_TOKEN="your_gitcode_token"

# 永久生效，添加到 shell 配置
echo 'export GC_TOKEN="your_gitcode_token"' >> ~/.bashrc
source ~/.bashrc

# 方式二：交互式登录
gc auth login --token YOUR_TOKEN
```

### 测试仓库

本文档使用以下测试仓库：
- `infra-test/gctest1`

---

## 认证命令 (auth)

### auth login - 登录

```bash
# 交互式登录
gc auth login

# 使用 Token 登录
gc auth login --token YOUR_TOKEN
```

### auth status - 查看认证状态

```bash
gc auth status
```

输出示例：
```
gitcode.com
  ✓ Logged in as username (GC_TOKEN)
  ✓ Git operations protocol: https
```

### auth token - 显示 Token

```bash
gc auth token
```

### auth logout - 登出

```bash
gc auth logout
```

---

## 仓库命令 (repo)

### repo view - 查看仓库

```bash
# 查看仓库详情
gc repo view infra-test/gctest1

# 在浏览器中打开
gc repo view infra-test/gctest1 --web
```

### repo list - 列出仓库

```bash
# 列出自己的仓库
gc repo list

# 列出指定组织的仓库
gc repo list --owner infra-test

# 限制数量
gc repo list --limit 10

# 只列出公开仓库
gc repo list --visibility public
```

### repo create - 创建仓库

```bash
# 创建公开仓库
gc repo create my-repo --public

# 创建私有仓库
gc repo create my-repo --private

# 创建带描述的仓库
gc repo create my-repo --public --description "My project"
```

> **注意**: 在组织下创建仓库需要有组织的相应权限。

### repo fork - Fork 仓库

```bash
# Fork 仓库到自己的账户
gc repo fork owner/repo

# Fork 并克隆到本地
gc repo fork owner/repo --clone
```

---

## Issue 命令 (issue)

### issue create - 创建 Issue

```bash
# 创建 Issue
gc issue create -R infra-test/gctest1 --title "Bug: Something wrong" --body "Description here"

# 创建 Issue 并添加标签
gc issue create -R infra-test/gctest1 --title "Feature request" --body "Description" --label enhancement

# 指定受理人
gc issue create -R infra-test/gctest1 --title "Task" --body "Description" --assignee username
```

### issue list - 列出 Issues

```bash
# 列出所有开放的 Issues
gc issue list -R infra-test/gctest1

# 只列出已关闭的 Issues
gc issue list -R infra-test/gctest1 --state closed

# 按标签筛选
gc issue list -R infra-test/gctest1 --label bug

# 限制数量
gc issue list -R infra-test/gctest1 --limit 20
```

### issue view - 查看 Issue

```bash
# 查看 Issue 详情
gc issue view 1 -R infra-test/gctest1

# 查看评论
gc issue view 1 -R infra-test/gctest1 --comments

# 在浏览器中打开
gc issue view 1 -R infra-test/gctest1 --web
```

### issue close - 关闭 Issue

```bash
# 关闭 Issue
gc issue close 1 -R infra-test/gctest1
```

### issue reopen - 重开 Issue

```bash
# 重开 Issue
gc issue reopen 1 -R infra-test/gctest1
```

### issue comment - 添加评论

```bash
# 添加评论
gc issue comment 1 -R infra-test/gctest1 --body "This is a comment"
```

---

## Pull Request 命令 (pr)

### pr create - 创建 PR

```bash
# 创建 PR（需要先在分支上提交代码）
gc pr create -R infra-test/gctest1 --title "New feature" --body "Description"

# 指定基础分支
gc pr create -R infra-test/gctest1 --base main --title "Feature" --body "Description"

# 创建草稿 PR
gc pr create -R infra-test/gctest1 --title "WIP: Feature" --draft

# 从最后一次提交填充标题和内容
gc pr create -R infra-test/gctest1 --fill
```

### pr list - 列出 PRs

```bash
# 列出所有开放的 PRs
gc pr list -R infra-test/gctest1

# 只列出已关闭的 PRs
gc pr list -R infra-test/gctest1 --state closed

# 只列出已合并的 PRs
gc pr list -R infra-test/gctest1 --state merged

# 限制数量
gc pr list -R infra-test/gctest1 --limit 10
```

### pr view - 查看 PR

```bash
# 查看 PR 详情
gc pr view 1 -R infra-test/gctest1

# 查看评论
gc pr view 1 -R infra-test/gctest1 --comments

# 在浏览器中打开
gc pr view 1 -R infra-test/gctest1 --web
```

### pr diff - 查看 PR 差异

```bash
# 查看 PR 差异
gc pr diff 1 -R infra-test/gctest1
```

### pr checkout - 检出 PR 分支

```bash
# 检出 PR 到本地分支
gc pr checkout 1 -R infra-test/gctest1
```

### pr merge - 合并 PR

```bash
# 合并 PR（默认合并提交）
gc pr merge 1 -R infra-test/gctest1

# Squash 合并
gc pr merge 1 -R infra-test/gctest1 --squash

# Rebase 合并
gc pr merge 1 -R infra-test/gctest1 --rebase
```

### pr close - 关闭 PR

```bash
# 关闭 PR
gc pr close 1 -R infra-test/gctest1
```

### pr reopen - 重开 PR

```bash
# 重开 PR
gc pr reopen 1 -R infra-test/gctest1
```

### pr ready - 标记就绪状态

```bash
# 标记为就绪（取消草稿）
gc pr ready 1 -R infra-test/gctest1

# 标记为草稿
gc pr ready 1 -R infra-test/gctest1 --wip
```

---

## Release 命令 (release)

### release create - 创建 Release

```bash
# 创建 Release（建议包含 --notes 参数）
gc release create v1.0.0 -R infra-test/gctest1 --title "Version 1.0.0" --notes "Release notes"

# 创建预发布 Release
gc release create v1.0.0-beta -R infra-test/gctest1 --title "v1.0.0 Beta" --notes "Beta release" --prerelease

# 创建草稿 Release
gc release create v1.0.0 -R infra-test/gctest1 --title "v1.0.0" --notes "Draft" --draft

# 指定目标分支
gc release create v1.0.0 -R infra-test/gctest1 --title "v1.0.0" --notes "Release" --target main
```

> **注意**: `--notes` 参数是必需的，不带此参数可能返回 400 错误。

### release list - 列出 Releases

```bash
# 列出所有 Releases
gc release list -R infra-test/gctest1
```

### release view - 查看 Release

```bash
# 查看 Release 详情
gc release view v1.0.0 -R infra-test/gctest1

# 在浏览器中打开
gc release view v1.0.0 -R infra-test/gctest1 --web
```

### release upload - 上传资产

```bash
# 上传单个文件
gc release upload v1.0.0 app.zip -R infra-test/gctest1

# 上传多个文件
gc release upload v1.0.0 app.zip checksum.txt -R infra-test/gctest1
```

### release download - 下载资产

```bash
# 下载所有资产到当前目录
gc release download v1.0.0 -R infra-test/gctest1

# 下载到指定目录
gc release download v1.0.0 -R infra-test/gctest1 -o ./downloads/

# 下载指定文件
gc release download v1.0.0 app.zip -R infra-test/gctest1
```

---

## 标签命令 (label)

### label list - 列出标签

```bash
# 列出所有标签
gc label list -R infra-test/gctest1
```

### label create - 创建标签

```bash
# 创建标签
gc label create "bug" -R infra-test/gctest1 --color "#ff0000" --description "Bug report"
```

---

## 里程碑命令 (milestone)

### milestone list - 列出里程碑

```bash
# 列出所有里程碑
gc milestone list -R infra-test/gctest1
```

### milestone create - 创建里程碑

```bash
# 创建里程碑
gc milestone create "v1.0" -R infra-test/gctest1 --description "First release"
```

### milestone view - 查看里程碑

```bash
# 查看里程碑详情
gc milestone view 1 -R infra-test/gctest1
```

---

## 其他命令

### version - 显示版本

```bash
gc version
```

### help - 帮助

```bash
# 显示帮助
gc help

# 显示命令帮助
gc help issue
gc help issue create
```

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
| `NO_COLOR` | 禁用颜色输出 |

---

## 已知限制

以下功能受 GitCode API 限制，可能无法正常工作：

| 功能 | 限制说明 |
|------|----------|
| `repo fork` | 在某些情况下可能返回 400 错误 |
| `issue close/reopen` | 命令执行成功但状态可能不变化 |
| `pr review` | 可能返回 404 错误 |
| `milestone create/view` | 可能返回 400 错误 |
| `release edit/delete` | GitCode API 不返回 release ID |

---

**最后更新**: 2026-03-23