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
- `aflyingto/gitcode-cli`

---

## 认证命令 (auth)

### auth login - 登录

```bash
# 交互式登录
gc auth login

# 使用 Token 登录
gc auth login --token YOUR_TOKEN

# 从 stdin 读取 Token
echo "YOUR_TOKEN" | gc auth login --with-token
```

### auth status - 查看认证状态

```bash
gc auth status
```

### auth logout - 登出

```bash
gc auth logout
```

### auth token - 显示 Token

```bash
gc auth token
```

---

## 仓库命令 (repo)

### repo clone - 克隆仓库

```bash
# 克隆仓库
gc repo clone infra-test/gctest1

# 克隆到指定目录
gc repo clone infra-test/gctest1 my-project

# 使用 SSH 协议
gc repo clone infra-test/gctest1 --git-protocol ssh
```

### repo create - 创建仓库

```bash
# 创建公开仓库
gc repo create my-new-repo --public

# 创建私有仓库
gc repo create my-private-repo --private

# 创建并添加远程
gc repo create my-repo --public --add-remote

# 在组织下创建仓库
gc repo create infra-test/new-repo --public
```

### repo list - 列出仓库

```bash
# 列出自己的仓库
gc repo list

# 列出指定用户的仓库
gc repo list --owner aflyingto

# 列出组织的仓库
gc repo list --org infra-test

# 限制数量
gc repo list --limit 10
```

### repo view - 查看仓库

```bash
# 查看仓库详情
gc repo view infra-test/gctest1

# 在浏览器中打开
gc repo view infra-test/gctest1 --web
```

### repo fork - Fork 仓库

```bash
# Fork 仓库
gc repo fork infra-test/gctest1

# Fork 到指定组织
gc repo fork infra-test/gctest1 --org my-org

# Fork 并克隆
gc repo fork infra-test/gctest1 --clone
```

### repo delete - 删除仓库

```bash
# 删除仓库（需要确认）
gc repo delete owner/repo

# 跳过确认
gc repo delete owner/repo --yes
```

---

## Issue 命令 (issue)

### issue create - 创建 Issue

```bash
# 创建 Issue
gc issue create -R infra-test/gctest1 --title "Bug: Something wrong" --body "Description here"

# 从文件读取内容
gc issue create -R infra-test/gctest1 --title "Feature request" --body-file ./issue.md

# 添加标签
gc issue create -R infra-test/gctest1 --title "Bug report" --label bug,high-priority

# 指定受理人
gc issue create -R infra-test/gctest1 --title "Task" --assignee username
```

### issue list - 列出 Issues

```bash
# 列出所有 Issues
gc issue list -R infra-test/gctest1

# 只列出开放的 Issues
gc issue list -R infra-test/gctest1 --state open

# 只列出已关闭的 Issues
gc issue list -R infra-test/gctest1 --state closed

# 按标签筛选
gc issue list -R infra-test/gctest1 --label bug

# 按受理人筛选
gc issue list -R infra-test/gctest1 --assignee username

# 限制数量
gc issue list -R infra-test/gctest1 --limit 20
```

### issue view - 查看 Issue

```bash
# 查看 Issue 详情
gc issue view 123 -R infra-test/gctest1

# 在浏览器中打开
gc issue view 123 -R infra-test/gctest1 --web

# 查看评论
gc issue view 123 -R infra-test/gctest1 --comments
```

### issue close - 关闭 Issue

```bash
# 关闭 Issue
gc issue close 123 -R infra-test/gctest1

# 关闭并添加评论
gc issue close 123 -R infra-test/gctest1 --comment "Fixed in PR #456"
```

### issue reopen - 重开 Issue

```bash
# 重开 Issue
gc issue reopen 123 -R infra-test/gctest1
```

### issue comment - 添加评论

```bash
# 添加评论
gc issue comment 123 -R infra-test/gctest1 --body "This is a comment"

# 从文件读取评论
gc issue comment 123 -R infra-test/gctest1 --body-file ./comment.md
```

### issue edit - 编辑 Issue

```bash
# 编辑标题
gc issue edit 123 -R infra-test/gctest1 --title "New title"

# 编辑内容
gc issue edit 123 -R infra-test/gctest1 --body "New description"

# 添加标签
gc issue edit 123 -R infra-test/gctest1 --add-label bug,help-wanted

# 移除标签
gc issue edit 123 -R infra-test/gctest1 --remove-label wontfix
```

---

## Pull Request 命令 (pr)

### pr create - 创建 PR

```bash
# 创建 PR（自动检测当前分支）
gc pr create -R infra-test/gctest1 --title "New feature" --body "Description"

# 指定基础分支和特性分支
gc pr create -R infra-test/gctest1 --base main --head feature-branch --title "Feature" --body "Desc"

# 创建草稿 PR
gc pr create -R infra-test/gctest1 --title "WIP: Feature" --draft

# 从文件读取描述
gc pr create -R infra-test/gctest1 --title "Feature" --body-file ./pr-desc.md

# 指定受理人
gc pr create -R infra-test/gctest1 --title "Feature" --assignee username

# 添加标签
gc pr create -R infra-test/gctest1 --title "Feature" --label enhancement
```

### pr list - 列出 PRs

```bash
# 列出所有 PRs
gc pr list -R infra-test/gctest1

# 只列出开放的 PRs
gc pr list -R infra-test/gctest1 --state open

# 只列出已关闭的 PRs
gc pr list -R infra-test/gctest1 --state closed

# 只列出已合并的 PRs
gc pr list -R infra-test/gctest1 --state merged

# 按作者筛选
gc pr list -R infra-test/gctest1 --author username

# 按标签筛选
gc pr list -R infra-test/gctest1 --label bug

# 限制数量
gc pr list -R infra-test/gctest1 --limit 10
```

### pr view - 查看 PR

```bash
# 查看 PR 详情
gc pr view 456 -R infra-test/gctest1

# 在浏览器中打开
gc pr view 456 -R infra-test/gctest1 --web

# 查看评论
gc pr view 456 -R infra-test/gctest1 --comments
```

### pr checkout - 检出 PR 分支

```bash
# 检出 PR 到本地分支
gc pr checkout 456 -R infra-test/gctest1

# 检出到指定分支名
gc pr checkout 456 -R infra-test/gctest1 --branch my-feature
```

### pr merge - 合并 PR

```bash
# 合并 PR（默认合并提交）
gc pr merge 456 -R infra-test/gctest1

# Squash 合并
gc pr merge 456 -R infra-test/gctest1 --squash

# Rebase 合并
gc pr merge 456 -R infra-test/gctest1 --rebase

# 删除分支
gc pr merge 456 -R infra-test/gctest1 --delete-branch
```

### pr close - 关闭 PR

```bash
# 关闭 PR
gc pr close 456 -R infra-test/gctest1

# 关闭并添加评论
gc pr close 456 -R infra-test/gctest1 --comment "Closing this PR"
```

### pr reopen - 重开 PR

```bash
# 重开 PR
gc pr reopen 456 -R infra-test/gctest1
```

### pr review - 代码检视

```bash
# 批准 PR
gc pr review 456 -R infra-test/gctest1 --approve

# 批准并添加评论
gc pr review 456 -R infra-test/gctest1 --approve --body "Looks good!"

# 请求修改
gc pr review 456 -R infra-test/gctest1 --request-changes --body "Please fix the bug"

# 添加评论（不批准/不拒绝）
gc pr review 456 -R infra-test/gctest1 --comment --body "Just a comment"
```

### pr diff - 查看 PR 差异

```bash
# 查看 PR 差异
gc pr diff 456 -R infra-test/gctest1

# 使用特定 diff 工具
gc pr diff 456 -R infra-test/gctest1 --color always
```

### pr ready - 标记就绪状态

```bash
# 标记为就绪（取消草稿）
gc pr ready 456 -R infra-test/gctest1

# 标记为草稿
gc pr ready 456 -R infra-test/gctest1 --wip
```

---

## Release 命令 (release)

### release create - 创建 Release

```bash
# 创建 Release
gc release create v1.0.0 -R infra-test/gctest1 --title "Version 1.0.0"

# 创建带说明的 Release
gc release create v1.0.0 -R infra-test/gctest1 --title "v1.0.0" --notes "First stable release"

# 从文件读取说明
gc release create v1.0.0 -R infra-test/gctest1 --title "v1.0.0" --notes-file RELEASE_NOTES.md

# 创建草稿 Release
gc release create v1.0.0 -R infra-test/gctest1 --title "v1.0.0" --draft

# 创建预发布 Release
gc release create v1.0.0 -R infra-test/gctest1 --title "v1.0.0-beta" --prerelease

# 指定目标分支
gc release create v1.0.0 -R infra-test/gctest1 --title "v1.0.0" --target main
```

### release list - 列出 Releases

```bash
# 列出所有 Releases
gc release list -R infra-test/gctest1

# 限制数量
gc release list -R infra-test/gctest1 --limit 10
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

# 上传特定类型文件
gc release upload v1.0.0 myapp.tar.gz myapp.rpm -R infra-test/gctest1
```

### release download - 下载资产

```bash
# 下载所有资产
gc release download v1.0.0 -R infra-test/gctest1

# 下载到指定目录
gc release download v1.0.0 -R infra-test/gctest1 -o ./downloads/

# 下载指定文件
gc release download v1.0.0 app.zip -R infra-test/gctest1

# 下载包括源码包
gc release download v1.0.0 -R infra-test/gctest1 --all
```

### release edit - 编辑 Release

```bash
# 修改标题
gc release edit v1.0.0 -R infra-test/gctest1 --title "New Title"

# 修改说明
gc release edit v1.0.0 -R infra-test/gctest1 --notes "Updated notes"

# 从文件读取说明
gc release edit v1.0.0 -R infra-test/gctest1 --notes-file NEW_NOTES.md

# 修改草稿状态
gc release edit v1.0.0 -R infra-test/gctest1 --draft true

# 修改预发布状态
gc release edit v1.0.0 -R infra-test/gctest1 --prerelease false
```

### release delete - 删除 Release

```bash
# 删除 Release（需要确认）
gc release delete v1.0.0 -R infra-test/gctest1

# 跳过确认
gc release delete v1.0.0 -R infra-test/gctest1 --yes
```

---

## 标签命令 (label)

### label create - 创建标签

```bash
# 创建标签
gc label create bug -R infra-test/gctest1 --description "Bug report" --color ff0000

# 从预设创建
gc label create enhancement -R infra-test/gctest1
```

### label list - 列出标签

```bash
# 列出所有标签
gc label list -R infra-test/gctest1

# 搜索标签
gc label list -R infra-test/gctest1 --search bug
```

### label delete - 删除标签

```bash
# 删除标签
gc label delete bug -R infra-test/gctest1 --yes
```

---

## 里程碑命令 (milestone)

### milestone create - 创建里程碑

```bash
# 创建里程碑
gc milestone create "v1.0" -R infra-test/gctest1 --description "First release"

# 指定截止日期
gc milestone create "v2.0" -R infra-test/gctest1 --due-date "2024-12-31"
```

### milestone list - 列出里程碑

```bash
# 列出所有里程碑
gc milestone list -R infra-test/gctest1

# 按状态筛选
gc milestone list -R infra-test/gctest1 --state open
gc milestone list -R infra-test/gctest1 --state closed
```

### milestone view - 查看里程碑

```bash
# 查看里程碑详情
gc milestone view 1 -R infra-test/gctest1
```

### milestone delete - 删除里程碑

```bash
# 删除里程碑
gc milestone delete 1 -R infra-test/gctest1 --yes
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

### completion - 生成补全脚本

```bash
# Bash
gc completion bash > /etc/bash_completion.d/gc

# Zsh
gc completion zsh > "${fpath[1]}/_gc"

# Fish
gc completion fish > ~/.config/fish/completions/gc.fish
```

---

## 常用选项

| 选项 | 说明 |
|------|------|
| `-R, --repo owner/repo` | 指定仓库 |
| `--help` | 显示帮助 |
| `--yes` | 跳过确认 |
| `--limit N` | 限制结果数量 |
| `--json` | JSON 格式输出 |
| `--web` | 在浏览器中打开 |

---

## 环境变量

| 变量 | 说明 |
|------|------|
| `GC_TOKEN` | 认证 Token |
| `GITCODE_TOKEN` | 备用 Token |
| `GC_HOST` | 默认主机（默认：gitcode.com） |
| `NO_COLOR` | 禁用颜色输出 |
| `GC_CONFIG_DIR` | 配置目录（默认：~/.config/gc） |