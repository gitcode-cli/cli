# GitCode CLI 命令使用指南

> 项目概述和功能介绍请参阅 [README.md](../README.md)，开发与规范入口请参阅 [spec/README.md](../spec/README.md)，打包发布请参阅 [PACKAGING.md](./PACKAGING.md)。

本文档提供 `gc` 命令行工具所有命令的实际使用示例。

## 前置准备

### 仓库参数格式

大多数接受仓库参数的命令现在统一支持以下格式：

```bash
owner/repo
https://gitcode.com/owner/repo
git@gitcode.com:owner/repo.git
```

说明：
- 未显式传 `-R` 的命令，仍按各自命令说明决定是否支持从当前 Git 仓库自动推断。
- 传入 HTTPS 或 SSH 仓库地址时，CLI 会统一解析出目标仓库，不再要求手工改写成 `owner/repo`。

当前自动推断边界：
- 仅显式接入 `cmdutil.ResolveRepo(...)` 的命令支持缺省 `-R` 时从当前 Git 仓库推断目标仓库，当前主要覆盖 `issue` 相关命令与 `repo view` 等“作用于当前仓库”的安全场景。
- 仍需显式传目标仓库参数的命令，通常是语义上操作“另一个仓库”的命令，例如 `repo sync --target-repo` 这类显式目标仓库场景。

### Agent-Friendly CLI 能力

当前版本已开始收口面向 AI 代理和脚本的 CLI 契约：

- 高频只读命令逐步支持 `--json`
- 删除类命令支持 `--dry-run`
- 非交互环境下删除命令未显式传 `--yes` 会直接失败，不再隐式等待输入
- 可通过 `gc schema` 查询命令树和单命令元数据

当前已支持 `--json` 的高频只读命令：

- `repo view`
- `repo list`
- `issue list`
- `issue view`
- `pr list`
- `pr view`
- `release list`
- `release view`
- `label list`
- `milestone list`

其中 `issue list` 额外支持：

- `--format json|simple|table`
- `--time-format absolute|relative`
- `--template <go-template>`
- `--json` 与 `--format json` 等价，二者都应作为稳定机器可消费入口

`issue list` 的 `--format` 非法值应直接报用法错误，不应静默回退到默认格式。

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

# 打开浏览器生成 Token 后继续登录
gc auth login --web
```

说明：
- `auth login --web` 会打开 GitCode Token 页面，然后继续在终端中读取你粘贴的 Token 完成登录。
- 登录成功后 token 会写入本地配置；若同时设置了 `GC_TOKEN` 或 `GITCODE_TOKEN`，环境变量优先。

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

说明：
- `auth token` 输出当前实际生效的 token，解析顺序与 `auth status` 一致。

### auth logout - 登出

```bash
gc auth logout
```

说明：
- `auth logout` 会清理本地持久化认证信息。
- 若当前认证来自 `GC_TOKEN` 或 `GITCODE_TOKEN`，命令会提示你手动 `unset` 环境变量。

---

## 仓库命令 (repo)

### repo view - 查看仓库

```bash
# 查看仓库详情
gc repo view infra-test/gctest1
gc repo view

# 在浏览器中打开
gc repo view infra-test/gctest1 --web

# 输出 JSON
gc repo view infra-test/gctest1 --json
```

说明：
- 在当前 Git 仓库中执行时，`gc repo view` 可缺省仓库参数；CLI 会优先解析 `origin` remote，若不存在则回退到第一个 remote。

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

# 输出 JSON
gc repo list --json

# 表格输出
gc repo list --format table
```

### repo sync - 同步目录到目标仓库并创建 PR

```bash
# 将当前仓库 docs/api 同步到目标仓库的 mirror/api 目录
gc repo sync \
  --target-repo infra-test/target-repo \
  --source-dir docs/api \
  --target-dir mirror/api

# 指定 base 分支和 PR 标题
gc repo sync \
  --target-repo infra-test/target-repo \
  --source-dir pkg/contracts \
  --target-dir mirror/contracts \
  --base main \
  --title "sync: update contracts"

# 结构化输出
gc repo sync \
  --target-repo infra-test/target-repo \
  --source-dir docs/api \
  --target-dir mirror/api \
  --json
```

说明：
- 该命令必须在本地 Git 仓库内执行
- `--source-dir` 是当前仓库内要同步的目录
- `--target-dir` 是目标仓库中的子目录，不能是仓库根目录
- 命令会自动创建同步分支、提交、推送并创建目标 PR
- 如果目标目录内容与源目录一致，命令会直接返回“无变更”

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

说明：
- `repo fork` 现在会按传入的 `owner/repo` 执行 fork，不再使用硬编码仓库路径。
- `--clone` 会在 fork 成功后将 fork 出来的仓库克隆到当前目录。

### repo delete - 删除仓库

```bash
# 删除仓库（危险操作，需确认）
gc repo delete owner/repo

# 预演删除
gc repo delete owner/repo --dry-run

# 非交互执行
gc repo delete owner/repo --yes
```

说明：
- 默认会要求输入仓库名确认。
- 在非交互环境中，未显式传 `--yes` 会直接失败。

### repo stats - 代码贡献统计

```bash
# 获取 main 分支代码贡献统计
gc repo stats --branch main -R infra-test/gctest1

# 按作者筛选
gc repo stats --branch main --author username -R infra-test/gctest1

# 仅显示个人统计
gc repo stats --branch main --only-self -R infra-test/gctest1

# 指定日期范围
gc repo stats --branch main --since 2024-01-01 --until 2024-12-31 -R infra-test/gctest1
```

---

## Issue 命令 (issue)

### issue create - 创建 Issue

```bash
# 创建 Issue
gc issue create -R infra-test/gctest1 --title "Bug: Something wrong" --body "Description here"
gc issue create --title "Bug: Something wrong" --body "Description here"

# 创建 Issue 并添加标签
gc issue create -R infra-test/gctest1 --title "Feature request" --body "Description" --label enhancement

# 指定受理人
gc issue create -R infra-test/gctest1 --title "Task" --body "Description" --assignee username

# 预演创建
gc issue create -R infra-test/gctest1 --title "Task" --body "Description" --dry-run
```

说明：
- `issue create` 当前已支持 `--dry-run` 预演创建参数。
- 创建时携带 `--label`、`--milestone`、`--assignee` 已走兼容的 form 提交路径。
- `--assignee` 继续使用用户名输入，但客户端会先解析为 GitCode user ID，再提交到 issue API。
- 若 GitCode API 未实际应用 assignee，命令会成功完成创建并在 stderr 给出告警，避免自动化重试制造重复 issue。

### issue list - 列出 Issues

```bash
# 列出所有开放的 Issues
gc issue list -R infra-test/gctest1
gc issue list

# 只列出已关闭的 Issues
gc issue list -R infra-test/gctest1 --state closed

# 列出所有状态的 Issues
gc issue list -R infra-test/gctest1 --state all

# 按标签筛选
gc issue list -R infra-test/gctest1 --label bug,enhancement

# 限制数量
gc issue list -R infra-test/gctest1 --limit 20

# 按里程碑筛选
gc issue list -R infra-test/gctest1 --milestone "v1.0"

# 按受理人筛选
gc issue list -R infra-test/gctest1 --assignee username

# 按创建者筛选
gc issue list -R infra-test/gctest1 --creator username

# 按更新时间排序
gc issue list -R infra-test/gctest1 --sort updated --direction desc

# 按创建时间筛选
gc issue list -R infra-test/gctest1 --created-after "2024-01-01"
gc issue list -R infra-test/gctest1 --created-before "2024-12-31"

# 按更新时间筛选
gc issue list -R infra-test/gctest1 --updated-after "2024-01-01"

# 关键字搜索
gc issue list -R infra-test/gctest1 --search "bug"

# 组合使用
gc issue list -R infra-test/gctest1 --state open --milestone "v1.0" --sort updated

# 输出 JSON
gc issue list -R infra-test/gctest1 --json

# 输出格式
gc issue list -R infra-test/gctest1 --format json
gc issue list -R infra-test/gctest1 --format simple
gc issue list -R infra-test/gctest1 --format table

# 时间格式
gc issue list -R infra-test/gctest1 --time-format absolute
gc issue list -R infra-test/gctest1 --time-format relative

# 自定义模板
gc issue list -R infra-test/gctest1 --template '{{range .}}#{{.Number}} {{.Title}}{{"\n"}}{{end}}'
```

说明：
- `--json` 继续作为兼容入口保留，适合脚本和代理调用。
- `--format json` 与 `--json` 输出一致。
- `--time-format` 只影响文本展示中的时间格式，不改变 JSON 结构。
- `--template` 使用 Go template 渲染 issue 列表，当前与 `--json`、`--format` 互斥。
- 非法 `--format` 值会返回错误，不会静默降级为默认输出。
- `--since`、`--created-after`、`--created-before`、`--updated-after`、`--updated-before` 支持 `YYYY-MM-DD` 和 ISO 8601 时间。
- CLI 会在请求前自动规范化为 GitCode API 可接受的 RFC3339 时间戳。

### issue view - 查看 Issue

```bash
# 查看 Issue 详情
gc issue view 1 -R infra-test/gctest1
gc issue view 1

# 查看评论
gc issue view 1 -R infra-test/gctest1 --comments

# 在浏览器中打开
gc issue view 1 -R infra-test/gctest1 --web

# 输出 JSON
gc issue view 1 -R infra-test/gctest1 --json

# 查看评论并输出 JSON
gc issue view 1 -R infra-test/gctest1 --comments --json

# 相对时间
gc issue view 1 -R infra-test/gctest1 --time-format relative
```

说明：
- `issue view` 的文本详情输出会使用更稳定的元信息排布，便于人工和代理阅读。
- `--time-format absolute|relative` 只影响文本详情和评论区中的时间展示，不改变 `--json` 结构。
- `--json` 路径保持结构化输出，不受文本排版变化影响。

### issue close - 关闭 Issue

```bash
# 关闭 Issue
gc issue close 1 -R infra-test/gctest1
gc issue close 1
```

### issue edit - 编辑 Issue

```bash
# 修改标题
gc issue edit 1 --title "New title" -R infra-test/gctest1
gc issue edit 1 --title "New title"

# 修改描述
gc issue edit 1 --body "New description" -R infra-test/gctest1

# 修改状态（close/reopen）
gc issue edit 1 --state close -R infra-test/gctest1
gc issue edit 1 --state reopen -R infra-test/gctest1

# 指派负责人
gc issue edit 1 --assignee username -R infra-test/gctest1
gc issue edit 1 --assignee user1 --assignee user2 -R infra-test/gctest1

# 设置标签
gc issue edit 1 --label bug,enhancement -R infra-test/gctest1

# 设置里程碑
gc issue edit 1 --milestone 5 -R infra-test/gctest1

# 设置为私有 Issue
gc issue edit 1 --security-hole -R infra-test/gctest1

# 组合使用
gc issue edit 1 --title "Bug fix" --assignee username --label bug --milestone 1 -R infra-test/gctest1
```

说明：
- `issue edit --assignee` 使用用户名输入，客户端会先解析为 GitCode user ID，再调用 issue 更新接口。
- 若 GitCode API 未实际应用 assignee，命令会成功完成更新并在 stderr 给出告警，避免自动化流程误把已成功更新当成失败重试。

### issue reopen - 重开 Issue

```bash
# 重开 Issue
gc issue reopen 1 -R infra-test/gctest1
gc issue reopen 1
```

### issue comment - 添加评论

```bash
# 添加评论
gc issue comment 1 -R infra-test/gctest1 --body "This is a comment"
gc issue comment 1 --body "This is a comment"

# 从文件读取评论内容
gc issue comment 1 -R infra-test/gctest1 --body-file comment.txt

# 从 stdin 读取评论内容
echo "Comment from stdin" | gc issue comment 1 -R infra-test/gctest1 --body-file -
```

### issue comment edit - 编辑 Issue 评论

```bash
# 按参数编辑评论
gc issue comment edit 166061383 -R infra-test/gctest1 --body "Updated comment"

# 按 --id 编辑评论
gc issue comment edit --id 166061383 -R infra-test/gctest1 --body "Updated comment"

# 从文件读取新内容
gc issue comment edit 166061383 -R infra-test/gctest1 --body-file comment.md
```

### issue comments - 列出 Issue 评论

```bash
# 列出评论
gc issue comments 1 -R infra-test/gctest1
gc issue comments 1

# 限制返回数量
gc issue comments 1 -R infra-test/gctest1 --limit 10

# 倒序排列
gc issue comments 1 -R infra-test/gctest1 --order desc

# 按更新时间筛选
gc issue comments 1 -R infra-test/gctest1 --since "2024-01-01T00:00:00+08:00"

# JSON 输出
gc issue comments 1 -R infra-test/gctest1 --json
```

### issue label - 管理 Issue 标签

```bash
# 添加标签
gc issue label 1 --add bug,enhancement -R infra-test/gctest1
gc issue label 1 --add bug,enhancement

# 移除标签
gc issue label 1 --remove bug -R infra-test/gctest1

# 列出标签
gc issue label 1 --list -R infra-test/gctest1
```

### issue prs - 查看 Issue 关联的 PRs

```bash
# 查看 Issue 关联的 Pull Requests
gc issue prs 123 -R infra-test/gctest1
gc issue prs 123

# 获取增强信息（包含可合并状态）
gc issue prs 123 --mode 1 -R infra-test/gctest1
```

说明：
- `issue create/list/view/close/reopen/comment/comments/edit/label/prs` 在当前 Git 仓库中可缺省 `-R`，CLI 会优先解析 `origin` remote；若没有 `origin`，则回退到第一个 remote。
- 若当前目录不是 Git 仓库，或仓库没有可用 remote，会返回明确错误并提示改用 `-R owner/repo`。

### issue relations - 查看仓库内 Issue / PR 关联表

```bash
# 查看仓库内所有 issue / PR 关联关系
gc issue relations -R infra-test/gctest1

# 输出 JSON 关系行
gc issue relations -R infra-test/gctest1 --json

# 只扫描开放 issue
gc issue relations -R infra-test/gctest1 --state open --limit 50
```

说明：
- 该命令会遍历仓库 issue，并获取每个 issue 关联的 PR。
- 文本输出按 PR 聚合，并同时显示关联 issue 的状态信息。
- `--json` 输出为关系行数组，每一行包含 `pr` 和 `issue` 两部分。

---

## Pull Request 命令 (pr)

### pr create - 创建 PR

```bash
# 创建 PR（自动检测当前分支作为 head）
gc pr create -R infra-test/gctest1 --title "New feature" --body "Description"

# 指定 head 分支
gc pr create -R infra-test/gctest1 --head feature-branch --title "Feature" --body "Description"

# 指定基础分支
gc pr create -R infra-test/gctest1 --base main --title "Feature" --body "Description"

# 创建草稿 PR
gc pr create -R infra-test/gctest1 --title "WIP: Feature" --draft

# 创建跨仓库 PR（从 fork 到 upstream）
gc pr create -R upstream/repo --fork myfork/repo --head feature-branch --title "Feature"

# 从最后一次提交填充标题和内容
gc pr create -R infra-test/gctest1 --fill

# 创建成功后在浏览器中打开 PR
gc pr create -R infra-test/gctest1 --title "New feature" --body "Description" --web
```

> **说明**: `--head` 参数可选，未指定时自动检测当前 Git 分支。
> `--fill` 会使用最近一次 Git commit 的标题和正文补全未显式提供的 `--title` / `--body`。
> `--web` 会在 PR 创建成功后打开新建 PR 页面。
> 当前分支解析已统一接入 `Factory.Branch`；若当前目录不是 Git 仓库或无法识别分支，会明确提示改用 `--head`。

### pr list - 列出 PRs

```bash
# 列出所有开放的 PRs
gc pr list -R infra-test/gctest1

# 只列出已关闭的 PRs
gc pr list -R infra-test/gctest1 --state closed

# 只列出已合并的 PRs
gc pr list -R infra-test/gctest1 --state merged

# 按 head / base 分支过滤
gc pr list -R infra-test/gctest1 --head feature/login --base main

# 限制数量
gc pr list -R infra-test/gctest1 --limit 10

# 排序与分页
gc pr list -R infra-test/gctest1 --sort updated --direction desc --page 2

# 输出 JSON
gc pr list -R infra-test/gctest1 --json

# 表格输出
gc pr list -R infra-test/gctest1 --format table
```

### pr view - 查看 PR

```bash
# 查看 PR 详情
gc pr view 1 -R infra-test/gctest1

# 查看评论
gc pr view 1 -R infra-test/gctest1 --comments

# 在浏览器中打开
gc pr view 1 -R infra-test/gctest1 --web

# 输出 JSON
gc pr view 1 -R infra-test/gctest1 --json

# 查看评论并输出 JSON
gc pr view 1 -R infra-test/gctest1 --comments --json

# 相对时间
gc pr view 1 -R infra-test/gctest1 --time-format relative
```

说明：
- `pr view` 的文本详情输出会使用更稳定的元信息排布，便于人工和代理阅读。
- `--time-format absolute|relative` 只影响文本详情和评论区中的时间展示，不改变 `--json` 结构。
- `--json` 路径保持结构化输出，不受文本排版变化影响。

### pr comments - 查看 PR 评论

```bash
# 查看 PR 评论列表
gc pr comments 1 -R infra-test/gctest1

# 限制评论数量
gc pr comments 1 --limit 5 -R infra-test/gctest1
```

评论列表会显示 `Discussion ID`，可直接用于 `gc pr reply --discussion`。
当前 GitCode 公开 API 不支持通过 CLI 将 PR 评论标记为已解决或未解决；resolved 状态需要在 Web UI 中手动处理。

### pr reply - 回复 PR 评论

```bash
# 回复评论讨论
gc pr reply 1 --discussion <discussion_id> --body "回复内容" -R infra-test/gctest1

# 使用简写
gc pr reply 1 -d <discussion_id> -b "回复内容" -R owner/repo
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

# 非交互执行
gc pr merge 1 -R infra-test/gctest1 --yes

# Squash 合并
gc pr merge 1 -R infra-test/gctest1 --method squash

# Rebase 合并
gc pr merge 1 -R infra-test/gctest1 --method rebase
```

说明：
- `pr merge` 属于高风险写操作，默认需要确认。
- 非交互场景中显式传 `--yes`。

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

### pr review - 评审 PR

```bash
# 评论 PR
gc pr review 1 --comment "评审意见" -R infra-test/gctest1

# 批准 PR
gc pr review 1 --approve -R infra-test/gctest1

# 批准 PR 并附带评论
gc pr review 1 --approve --comment "LGTM" -R infra-test/gctest1

# 强制通过审批（管理员权限）
gc pr review 1 --approve --force -R infra-test/gctest1
```

说明：
- `--approve` 现在走 GitCode 实际可用的 `/pulls/:number/review` endpoint，不再命中错误的 `/reviews` 路径。
- `--approve --comment` 会先提交普通评论，再执行批准动作。
- GitCode 当前公开 API 不支持“request changes”动作，`--request` 会明确报错并提示改用 `--comment` 留下审查意见。

### pr edit - 编辑 PR

```bash
# 修改标题
gc pr edit 1 --title "新标题" -R infra-test/gctest1

# 修改描述
gc pr edit 1 --body "新描述" -R infra-test/gctest1

# 设置草稿状态
gc pr edit 1 --draft true -R infra-test/gctest1

# 取消草稿状态
gc pr edit 1 --draft false -R infra-test/gctest1

# 添加标签
gc pr edit 1 --labels bug,enhancement -R infra-test/gctest1

# 设置里程碑
gc pr edit 1 --milestone 5 -R infra-test/gctest1
```

### pr test - 触发 PR 测试

```bash
# 触发测试
gc pr test 1 -R infra-test/gctest1

# 强制通过测试（管理员权限）
gc pr test 1 --force -R infra-test/gctest1
```

### pr sync - 同步 PR 到另一个仓库

```bash
# 同步 PR 到目标仓库
gc pr sync --source-pr owner/source-repo#123 --target-repo owner/target-repo

# 指定目标分支
gc pr sync --source-pr owner/source-repo#123 \
  --target-repo owner/target-repo \
  --base release/v1.0

# 自定义标题和内容
gc pr sync --source-pr owner/source-repo#123 \
  --target-repo owner/target-repo \
  --title "[sync] Fix login bug" \
  --body "从 owner/source-repo#123 同步"

# 创建草稿 PR
gc pr sync --source-pr owner/source-repo#123 \
  --target-repo owner/target-repo \
  --draft

# 结构化输出
gc pr sync --source-pr owner/source-repo#123 \
  --target-repo owner/target-repo \
  --json
```

说明：
- `--source-pr` 支持两种格式：`owner/repo#number` 或完整 GitCode URL，例如 `https://gitcode.com/owner/repo/merge_requests/123`
- 命令会按原顺序逐个 cherry-pick 源 PR 的所有 commits 到目标仓库，保留提交边界
- 新 PR 标题默认格式：`[sync] {源 PR 标题}`
- 新 PR 内容默认继承源 PR 内容并追加同步来源信息
- 如遇 cherry-pick 冲突，命令会报错并提示手动处理

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

# 输出 JSON
gc release list -R infra-test/gctest1 --json
```

说明：
- 文本输出中只有最新一个正式 release 会标记为 `(latest)`。
- 其他正式 release 会标记为 `(published)`；草稿和预发布仍分别显示 `(draft)`、`(pre-release)`。

### release view - 查看 Release

```bash
# 查看 Release 详情
gc release view v1.0.0 -R infra-test/gctest1

# 在浏览器中打开
gc release view v1.0.0 -R infra-test/gctest1 --web

# 输出 JSON
gc release view v1.0.0 -R infra-test/gctest1 --json
```

说明：
- 当 GitCode API 未返回资产大小时，文本输出会显示 `unknown size`，避免把未知值误写成 `0 bytes`。

### release upload - 上传资产

```bash
# 上传单个文件
gc release upload v1.0.0 app.zip -R infra-test/gctest1

# 上传多个文件
gc release upload v1.0.0 app.zip checksum.txt -R infra-test/gctest1
```

说明：
- `--label` 参数当前不受 GitCode release upload API 支持；CLI 现在会直接报错，不再静默忽略。

### release download - 下载资产

```bash
# 下载 latest release 的所有资产到当前目录
gc release download -R infra-test/gctest1

# 下载指定 release 的所有资产到当前目录
gc release download v1.0.0 -R infra-test/gctest1

# 下载到指定目录
gc release download v1.0.0 -R infra-test/gctest1 -o ./downloads/

# 下载指定文件
gc release download v1.0.0 app.zip -R infra-test/gctest1
```

### release edit - 编辑 Release

```bash
# 修改标题
gc release edit v1.0.0 --title "New title" -R infra-test/gctest1

# 修改说明
gc release edit v1.0.0 --notes "New release notes" -R infra-test/gctest1
```

说明：
- 若 GitCode 当前 release 查询响应未返回 release ID，`release edit` 会明确报错提示上游 API 限制，而不是给出模糊失败信息。

### release delete - 删除 Release

```bash
# 删除 Release
gc release delete v1.0.0 -R infra-test/gctest1

# 预演删除
gc release delete v1.0.0 -R infra-test/gctest1 --dry-run

# 非交互执行
gc release delete v1.0.0 -R infra-test/gctest1 --yes
```

说明：
- 若 GitCode 当前 release 查询响应未返回 release ID，`release delete` 会明确报错提示上游 API 限制。

---

## Commit 命令 (commit)

### commit view - 查看提交

```bash
# 查看提交详情
gc commit view abc123 -R infra-test/gctest1

# 显示变更文件
gc commit view abc123 -R infra-test/gctest1 --show-diff

# 输出 JSON 格式
gc commit view abc123 -R infra-test/gctest1 --json

# 在浏览器打开
gc commit view abc123 -R infra-test/gctest1 --web
```

### commit diff - 获取提交差异

```bash
# 获取提交 diff
gc commit diff abc123 -R infra-test/gctest1
```

### commit patch - 获取提交补丁

```bash
# 获取提交 patch
gc commit patch abc123 -R infra-test/gctest1
```

### commit comments create - 创建提交评论

```bash
# 创建评论
gc commit comments create abc123 --body "Nice work!" -R infra-test/gctest1
```

### commit comments view - 查看提交评论

```bash
# 查看指定评论
gc commit comments view 123 -R infra-test/gctest1
```

### commit comments edit - 编辑提交评论

```bash
# 编辑评论
gc commit comments edit 123 --body "Updated comment" -R infra-test/gctest1
```

### commit comments list - 列出仓库所有评论

```bash
# 列出所有评论
gc commit comments list -R infra-test/gctest1

# 分页
gc commit comments list -R infra-test/gctest1 --page 1 --per-page 50
```

### commit comments list-by-sha - 列出指定提交的评论

```bash
# 列出某提交的所有评论
gc commit comments list-by-sha abc123 -R infra-test/gctest1
```

---

## 标签命令 (label)

### label list - 列出标签

```bash
# 列出所有标签
gc label list -R infra-test/gctest1

# 结构化输出
gc label list -R infra-test/gctest1 --json
```

### label create - 创建标签

```bash
# 创建标签
gc label create "bug" -R infra-test/gctest1 --color "#ff0000" --description "Bug report"
```

### label delete - 删除标签

```bash
# 删除标签
gc label delete bug -R infra-test/gctest1

# 预演删除
gc label delete bug -R infra-test/gctest1 --dry-run

# 非交互执行
gc label delete bug -R infra-test/gctest1 --yes
```

---

## 里程碑命令 (milestone)

### milestone list - 列出里程碑

```bash
# 列出所有里程碑
gc milestone list -R infra-test/gctest1

# 结构化输出
gc milestone list -R infra-test/gctest1 --json
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

### milestone delete - 删除里程碑

```bash
# 删除里程碑
gc milestone delete 1 -R infra-test/gctest1

# 预演删除
gc milestone delete 1 -R infra-test/gctest1 --dry-run

# 非交互执行
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

### schema - 命令元数据

```bash
# 输出完整命令树
gc schema

# 输出单个命令的元数据
gc schema "issue view"
```

---

## 常用选项

| 选项 | 说明 |
|------|------|
| `-R, --repo owner/repo` | 指定仓库 |
| `--help` | 显示帮助 |
| `--limit N` | 限制结果数量 |
| `--web` | 在浏览器中打开 |
| `--json` | 输出结构化 JSON |
| `--dry-run` | 预演写操作而不执行 |

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
| `repo fork` | 仓库路径已按用户输入解析，但 GitCode API 在部分仓库上仍可能返回 `400 Bad Request` |
| `milestone create/view` | 返回 400 错误，API 可能不支持 |
| `release edit/delete` | GitCode API 不返回 release ID |

---

## 文档维护规范

**重要**：每次修改命令相关代码时，必须同步更新本文档！

### 同步更新要求

| 代码改动类型 | 需要更新的文档 |
|------------|--------------|
| 新增命令 | docs/COMMANDS.md、README.md |
| 新增子命令 | docs/COMMANDS.md |
| 修改命令参数/flags | docs/COMMANDS.md、README.md |
| 修改命令行为 | docs/COMMANDS.md |
| 删除命令 | docs/COMMANDS.md、README.md |

### 更新检查清单

开发完成后，确认以下检查项：

- [ ] 新命令已添加到 docs/COMMANDS.md
- [ ] README.md 命令概览已更新（如有新命令）
- [ ] 命令示例已验证可执行
- [ ] 参数说明与代码实现一致
- [ ] 已知限制表已更新（如有新的 API 限制）

### 常见问题

**Q: 如何确认文档与代码一致？**
```bash
# 查看所有命令
gc help

# 查看具体命令帮助
gc pr --help
gc issue --help
```

**Q: 文档更新顺序？**
1. 先更新 docs/COMMANDS.md（完整文档）
2. 再更新 README.md（概览文档）

---

**最后更新**: 2026-03-25
