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

SSH default for code transfer:
- Code download and sync paths use SSH by default.
- `repo clone owner/repo` defaults to `git@gitcode.com:owner/repo.git` unless `--git-protocol https` or a saved config explicitly selects HTTPS.
- SSH-based code transfer requires a local SSH key with access to `git@gitcode.com`.


说明：
- 未显式传 `-R` 的命令，仍按各自命令说明决定是否支持从当前 Git 仓库自动推断。
- 传入 HTTPS 或 SSH 仓库地址时，CLI 会统一解析出目标仓库，不再要求手工改写成 `owner/repo`。

### Windows PowerShell 命令名和 stdin

Windows PowerShell 预置 `gc` 作为 `Get-Content` 别名。若在 PowerShell 中使用 CLI，推荐写完整命令名 `gitcode`；也可以显式调用 `gc.exe`。

当中文或其他非 ASCII 正文需要通过 `--body-file -` / `--comment-file -` 从 stdin 传入时，推荐使用 UTF-8 文件：

```powershell
Set-Content -Path body.md -Value "中文正文" -Encoding UTF8
gitcode issue create -R owner/repo --title "标题" --body-file body.md
```

若必须直接管道传入，请先设置 UTF-8 输出编码：

```powershell
$OutputEncoding = [System.Text.UTF8Encoding]::new($false)
"中文正文" | gitcode issue create -R owner/repo --title "标题" --body-file -
```

CLI 只会在显式 stdin 文本 flag（当前包括 `--body-file -` 和 `--comment-file -`）上拦截疑似已被 Windows PowerShell 损坏成 `???` 的输入，并在 stderr 提示正确用法；如果确实需要原样传入连续问号，可设置 `GITCODE_CLI_ALLOW_LOSSY_STDIN=1`。

当前自动推断边界：
- 仅显式接入 `cmdutil.ResolveRepo(...)` 的命令支持缺省 `-R` 时从当前 Git 仓库推断目标仓库，当前主要覆盖 `issue` 相关命令、`repo view/log/branch view`，以及 `pr list/view/issues`、`release list/view`、`commit view`、`label list`、`milestone list/view` 等"作用于当前仓库"的安全只读场景。
- 仍需显式传目标仓库参数的命令，通常是语义上操作“另一个仓库”的命令，例如 `repo sync --target-repo` 这类显式目标仓库场景。

### Agent-Friendly CLI 能力

当前版本已开始收口面向 AI 代理和脚本的 CLI 契约：

- 高频只读命令和高频写路径结果逐步支持 `--json`
- 删除类命令支持 `--dry-run`
- 非交互环境下删除命令未显式传 `--yes` 会直接失败，不再隐式等待输入
- 可通过 `gc schema` 查询命令树和单命令元数据

当前已支持 `--json` 的高频只读命令：

- `repo view`
- `repo list`
- `repo log`
- `repo branch view`
- `issue list`
- `issue view`
- `pr list`
- `pr view`
- `pr issues`
- `release list`
- `release view`
- `label list`
- `milestone list`

当前已支持 `--json` 的高频写路径命令：

- `issue create`
- `issue edit`
- `pr create`
- `pr edit`
- `pr merge`
- `repo create`
- `repo fork`
- `release create`
- `release edit`
- `release upload`
- `milestone edit`

其中 `issue list` 额外支持：

- `--format json|simple|table`
- `--time-format absolute|relative`
- `--template <go-template>`
- `--json` 与 `--format json` 等价，二者都应作为稳定机器可消费入口

`issue list` 的 `--format` 非法值应直接报用法错误，不应静默回退到默认格式。

`gc api` 可作为底层 API 调试入口，适合在 typed command 尚未覆盖时读取或调用 GitCode API；它输出远端原始响应，不额外包装 JSON。

### 退出码

`gc` 命令使用稳定的退出码语义，方便脚本和 AI 代理判断执行结果：

| 退出码 | 常量名 | 含义 | 典型场景 |
|--------|--------|------|----------|
| 0 | ExitSuccess | 命令执行成功 | 正常完成 |
| 1 | ExitError | 通用错误 | API 错误、网络错误 |
| 2 | ExitUsage | 参数用法错误 | 缺少必选参数、参数格式错误 |
| 3 | ExitNotFound | 资源不存在 | issue/pr/repo 不存在 |
| 4 | ExitAuth | 认证错误 | 未登录或 token 无效 |
| 5 | ExitConflict | 资源冲突 | PR merge 冲突 |

示例：
```bash
# 检查退出码
gc issue view abc  # 无效 issue 号
echo $?            # 输出 2 (ExitUsage)

gc issue view 999999  # 不存在的 issue
echo $?               # 输出 3 (ExitNotFound)

gc auth status  # 未认证时
echo $?         # 输出 4 (ExitAuth)
```

### 认证

```bash
# 方式一：设置环境变量（推荐）
export GC_TOKEN="your_gitcode_token"

# 永久生效，添加到 shell 配置
echo 'export GC_TOKEN="your_gitcode_token"' >> ~/.bashrc
source ~/.bashrc

# 方式二：非交互登录
echo "YOUR_TOKEN" | gc auth login --with-token
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

# 从 stdin 读取 Token 登录
echo "YOUR_TOKEN" | gc auth login --with-token

# 打开浏览器生成 Token 后继续登录
gc auth login --web
```

说明：
- `auth login --web` 会打开 GitCode Token 页面，然后继续在终端中读取你粘贴的 Token 完成登录。
- 登录成功后 token 会写入本地配置；若同时设置了 `GC_TOKEN` 或 `GITCODE_TOKEN`，环境变量优先。
- 未显式传 `--with-token` 时需要交互式 TTY；非交互环境会直接报错，避免命令挂起等待输入。

### auth status - 查看认证状态

```bash
gc auth status

# 查看指定主机的持久化认证状态
gc auth status --hostname gitcode.com

# 显示完整 token（谨慎使用）
gc auth status --show-token

# 输出 JSON
gc auth status --json
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

# 输出指定主机的已存储 token
gc auth token --hostname gitcode.com
```

说明：
- `auth token` 输出当前实际生效的 token，解析顺序与 `auth status` 一致。
- 显式传 `--hostname` 时，会读取该主机已存储的 token，不再被通用环境变量覆盖。

### auth logout - 登出

```bash
gc auth logout

# 非交互执行
gc auth logout --yes
```

说明：
- `auth logout` 会清理本地持久化认证信息。
- 若当前认证来自 `GC_TOKEN` 或 `GITCODE_TOKEN`，命令会提示你手动 `unset` 环境变量。
- 非交互场景中显式传 `--yes`。

---

## API 命令 (api)

### api - 调用 GitCode API

```bash
# 读取仓库 API 原始响应
gc api repos/infra-test/gctest1

# 读取 PR 文件列表
gc api repos/infra-test/gctest1/pulls/1/files

# 带查询参数的 API，包含 & 时请整体加引号
gc api 'repos/infra-test/gctest1/commits?path=README.md&sha=main'

# 指定 HTTP 方法和请求体文件
gc api repos/infra-test/gctest1/pulls/1 --method PATCH --input body.json

# 从 stdin 读取请求体
printf '{"title":"New title"}' | gc api repos/infra-test/gctest1/pulls/1 --method PATCH --input -

# 自定义请求头
gc api repos/infra-test/gctest1 --header 'Accept: application/json'
```

说明：
- endpoint 可写成 `repos/owner/repo` 或 `/api/v5/repos/owner/repo`；普通相对路径会自动补齐 `/api/v5/`。
- 认证沿用当前 `gc` 登录态或 `GC_TOKEN` / `GITCODE_TOKEN` 环境变量。
- 默认方法为 `GET`；传入 `--input` 但未指定 `--method` 时会自动使用 `POST`。
- 输出为远端原始响应 body，便于脚本继续交给 `python -m json.tool`、`jq` 或其他工具处理。

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

### repo branch view - 查看分支

```bash
# 查看分支详情
gc repo branch view main -R owner/repo

# 查看分支详情（当前仓库）
gc repo branch view main

# 输出 JSON
gc repo branch view main -R owner/repo --json
```

说明：
- `repo branch view` 显示指定分支的名称、保护状态和最新 commit 信息（ID、短 ID、标题、作者）。
- `--json` 输出分支对象，包含 `name`、`protected`、`commit.id`、`commit.short_id`、`commit.title`、`commit.message`、`commit.author.login`、`commit.committer.login`、`commit.created_at` 等字段。
- 分支不存在时返回明确错误。

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

# 结构化输出（与 --json 等价）
gc repo list --format json

# 简洁输出
gc repo list --format simple

# 表格输出
gc repo list --format table
```

### repo log - 查看仓库提交日志

```bash
# 查看最近提交
gc repo log -R infra-test/gctest1
gc repo log

# 查看指定分支上触碰某个文件的提交
gc repo log -R infra-test/gctest1 --file README.md --branch main

# 限制数量并输出 JSON
gc repo log -R infra-test/gctest1 --file README.md --branch main --limit 5 --json
```

说明：
- `repo log` 支持 `-R/--repo`，也支持在当前 Git 仓库中缺省 `-R` 自动推断目标仓库。
- `--file` 对应提交 API 的文件路径过滤，`--branch` 可传分支、tag 或 commit SHA。
- 文本输出会显示短 SHA、提交日期和提交信息首行；`--json` 输出远端提交对象数组。

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
  --yes \
  --json
```

说明：
- 该命令必须在本地 Git 仓库内执行
- `--source-dir` 是当前仓库内要同步的目录
- `--target-dir` 是目标仓库中的子目录，不能是仓库根目录
- 命令会自动创建同步分支、提交、推送并创建目标 PR
- `repo sync` clones and pushes the target repository over SSH; ensure an SSH key with access to `git@gitcode.com` is configured.
- 推送同步分支并创建目标 PR 前默认需要确认；非交互场景中显式传 `--yes`
- 如果目标目录内容与源目录一致，命令会直接返回“无变更”

### repo create - 创建仓库

```bash
# 创建公开仓库
gc repo create my-repo --public

# 创建私有仓库
gc repo create my-repo --private

# 创建带描述的仓库
gc repo create my-repo --public --description "My project"

# 创建后输出 JSON
gc repo create my-repo --private --json
```

> **注意**: 在组织下创建仓库需要有组织的相应权限。
> `--json` 只在成功创建后输出仓库对象；不会混入文本提示。

### repo fork - Fork 仓库

```bash
# Fork 仓库到自己的账户
gc repo fork owner/repo

# Fork 并克隆到本地
gc repo fork owner/repo --clone

# Fork 后输出 JSON
gc repo fork owner/repo --json
```

说明：
- `repo fork` 现在会按传入的 `owner/repo` 执行 fork，不再使用硬编码仓库路径。
- `--clone` 会在 fork 成功后将 fork 出来的仓库克隆到当前目录。
- `--json` 只在 fork 成功后输出 fork 仓库对象；不能与 `--clone` 同时使用。

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

# JSON 输出
gc repo stats --branch main -R infra-test/gctest1 --json
```

---

## Issue 命令 (issue)

### issue create - 创建 Issue

```bash
# 创建 Issue
gc issue create -R infra-test/gctest1 --title "Bug: Something wrong" --body "Description here"
gc issue create --title "Bug: Something wrong" --body "Description here"

# 从文件读取 body
gc issue create -R infra-test/gctest1 --title "Bug report" --body-file description.md

# 从 stdin 读取 body
echo "Description from stdin" | gc issue create -R infra-test/gctest1 --title "Bug report" --body-file -

# 创建 Issue 并添加标签
gc issue create -R infra-test/gctest1 --title "Feature request" --body "Description" --label enhancement

# 指定受理人
gc issue create -R infra-test/gctest1 --title "Task" --body "Description" --assignee username

# 使用模板路径创建
gc issue create -R infra-test/gctest1 --title "Feature request" --template-path .gitcode/ISSUE_TEMPLATE/feature.yaml

# 创建私有 Issue
gc issue create -R infra-test/gctest1 --title "Security report" --security-hole

# 传入高级字段（企业版）
gc issue create -R infra-test/gctest1 --title "Feature request" --issue-type "需求" --issue-severity "高"

# 通过 JSON 传入 custom_fields
gc issue create -R infra-test/gctest1 --title "Feature request" --custom-fields-json '[{"id":"field","value":"demo"}]'

# 从文件读取 custom_fields
gc issue create -R infra-test/gctest1 --title "Feature request" --custom-fields-file custom-fields.json

# 预演创建
gc issue create -R infra-test/gctest1 --title "Task" --body "Description" --dry-run

# 预演创建并输出 JSON
gc issue create -R infra-test/gctest1 --title "Task" --body "Description" --dry-run --json

# 创建后输出 JSON
gc issue create -R infra-test/gctest1 --title "Task" --body "Description" --json
```

说明：
- `issue create` 当前已支持 `--dry-run` 预演创建参数。
- `issue create --dry-run --json` 输出结构化预览，不执行真实创建。
- `--json` 只在成功创建并完成必要回读验证后输出 issue 对象；不会混入文本提示。
- 仅使用基础字段时，创建会继续走兼容的 repo 级 form 提交路径。
- 基础路径中的 `--assignee` 继续使用用户名输入，但客户端会先解析为 GitCode user ID，再提交到 issue API。
- 显式传入 `--template-path`、`--security-hole`、`--issue-type`、`--issue-severity`、`--custom-fields-json`、`--custom-fields-file` 时，会切换到 GitCode 文档化的 owner 级创建接口并透传高级字段。
- `--custom-fields-json` 与 `--custom-fields-file` 不能同时使用；两者都要求 JSON 顶层是 `object[]`。
- 模板路径支持仓库下的 `.gitcode`、`.github`、`.gitee` 目录；组织模板可能来自 `owner/.gitcode`，第一阶段仅支持显式传路径，不保证自动发现。
- 若 GitCode API 未实际应用 assignee，命令会成功完成创建并在 stderr 给出告警，避免自动化重试制造重复 issue。
- Windows PowerShell 中从 stdin 传中文正文时，建议使用 UTF-8 文件或先设置 `$OutputEncoding`，详见“Windows PowerShell 命令名和 stdin”。

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

# 查看所有评论（自动翻页获取全部）
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

# 非交互执行
gc issue close 1 -R infra-test/gctest1 --yes

# 输出 JSON
gc issue close 1 -R infra-test/gctest1 --json
```

说明：
- `-R` 在当前 Git 仓库目录执行时可省略，命令会自动推断目标仓库。
- `issue close` 属于写操作，默认需要确认；非交互场景中显式传 `--yes`。
- 命令会在关闭请求后验证 Issue 状态，避免服务端未实际关闭时误报成功。

### issue edit - 编辑 Issue

```bash
# 修改标题
gc issue edit 1 --title "New title" -R infra-test/gctest1
gc issue edit 1 --title "New title"

# 修改描述
gc issue edit 1 --body "New description" -R infra-test/gctest1

# 从文件读取新描述
gc issue edit 1 --body-file new-description.md -R infra-test/gctest1

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

# 更新后输出 JSON
gc issue edit 1 --title "Bug fix" -R infra-test/gctest1 --json
```

说明：
- `issue edit --assignee` 使用用户名输入，客户端会先解析为 GitCode user ID，再调用 issue 更新接口。
- 若 GitCode API 未实际应用 assignee，命令会返回失败并包含已更新 issue 的 URL，避免自动化流程静默误判。
- `--json` 只在成功更新并完成必要回读验证后输出 issue 对象；不会混入文本提示。

### issue reopen - 重开 Issue

```bash
# 重开 Issue
gc issue reopen 1 -R infra-test/gctest1
gc issue reopen 1

# 非交互执行
gc issue reopen 1 -R infra-test/gctest1 --yes

# 输出 JSON
gc issue reopen 1 -R infra-test/gctest1 --json
```

说明：
- `-R` 在当前 Git 仓库目录执行时可省略，命令会自动推断目标仓库。
- `issue reopen` 属于写操作，默认需要确认；非交互场景中显式传 `--yes`。
- 命令会在重开请求后验证 Issue 状态。

### issue comment - 添加评论

```bash
# 添加评论
gc issue comment 1 -R infra-test/gctest1 --body "This is a comment"
gc issue comment 1 --body "This is a comment"

# 从文件读取评论内容
gc issue comment 1 -R infra-test/gctest1 --body-file comment.txt

# 从 stdin 读取评论内容
echo "Comment from stdin" | gc issue comment 1 -R infra-test/gctest1 --body-file -

# 输出 JSON
gc issue comment 1 -R infra-test/gctest1 --body "This is a comment" --json
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

# 查看所有评论（--limit 0 或不指定 --limit，自动翻页获取全部）
gc issue comments 1 -R infra-test/gctest1
gc issue comments 1 -R infra-test/gctest1 --limit 0

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

# JSON 输出
gc issue prs 123 -R infra-test/gctest1 --json
```

说明：
- `issue create/list/view/close/reopen/comment/comments/edit/label/prs` 在当前 Git 仓库中可缺省 `-R`，CLI 会优先解析 `origin` remote；若没有 `origin`，则回退到第一个 remote。
- 若当前目录不是 Git 仓库，或仓库没有可用 remote，会返回明确错误并提示改用 `-R owner/repo`。
- `--mode 1` 显示 Mergeable 状态：绿色 `can merge` 表示可合并，红色 `cannot merge` 表示不可合并，黄色 `unknown` 表示未计算或未知。

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

# 从文件读取 PR 内容
gc pr create -R infra-test/gctest1 --title "New feature" --body-file description.md

# 从 stdin 读取 PR 内容
echo "Description from stdin" | gc pr create -R infra-test/gctest1 --title "New feature" --body-file -

# 指定 head 分支
gc pr create -R infra-test/gctest1 --head feature-branch --title "Feature" --body "Description"

# 指定基础分支
gc pr create -R infra-test/gctest1 --base main --title "Feature" --body "Description"

# 创建草稿 PR
gc pr create -R infra-test/gctest1 --title "WIP: Feature" --draft

# 创建跨仓库 PR（从 fork 到 upstream）；--fork 会自动把 head 规范化为 myfork:feature-branch
gc pr create -R upstream/repo --fork myfork/repo --head feature-branch --title "Feature"

# 等价写法：直接用 owner:branch 形式的 head，可省略 --fork
gc pr create -R upstream/repo --head myfork:feature-branch --title "Feature"

# 从最后一次提交填充标题和内容
gc pr create -R infra-test/gctest1 --fill

# 创建成功后在浏览器中打开 PR
gc pr create -R infra-test/gctest1 --title "New feature" --body "Description" --web

# 创建后输出 JSON
gc pr create -R infra-test/gctest1 --head feature-branch --title "Feature" --body "Description" --json
```

> **说明**: `--head` 参数可选，未指定时自动检测当前 Git 分支。
> `--body-file` 支持从文件读取 PR 内容；使用 `-` 可从 stdin 读取。`--body` 与 `--body-file` 不能同时使用。
> `--fill` 会使用最近一次 Git commit 的标题和正文补全未显式提供的 `--title` / `--body` / `--body-file`。
> `--web` 会在 PR 创建成功后打开新建 PR 页面。
> `--json` 只在成功创建后输出 PR 对象；不能与 `--web` 同时使用。
> 如果 GitCode 创建响应未返回 `body`，CLI 会在创建后尝试回读 PR；若回读仍未返回或无法确认远端 body，`--json` 会保持远端返回的空值并在 stderr 给出 warning，避免把本地提交内容伪装成远端事实。可用 `gitcode pr view <number> -R owner/repo --json` 再次核验。
> **跨仓库 PR**: 跨仓库（从 fork 到 upstream）的源仓库通过 `head="<fork_owner>:<branch>"` 表达，而不是已废弃的 `fork_path` 表单字段——后者在 GitCode v5 会错误解析源（upstream 同名分支 → 0 commits）甚至直接 403（见 #259）。使用 `--fork owner/repo` 时，CLI 会自动把 `--head` 规范化为 `<fork_owner>:<branch>`；若 `--head` 已是 `owner:branch` 形式则原样保留，可省略 `--fork`。
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

# 按里程碑筛选
gc pr list -R infra-test/gctest1 --milestone "v1.0"

# 限制数量
gc pr list -R infra-test/gctest1 --limit 10

# 排序与分页
gc pr list -R infra-test/gctest1 --sort updated --direction desc --page 2

# 自动翻页获取多页结果
gc pr list -R infra-test/gctest1 --paginate --per-page 100

# 按 PR 提交信息过滤
gc pr list -R infra-test/gctest1 --commit-message "fix login"

# 自动翻页后按提交信息过滤，仍可用 --limit 做结果上限
gc pr list -R infra-test/gctest1 --paginate --per-page 100 --limit 200 --commit-message "fix login" --json

# 输出 JSON
gc pr list -R infra-test/gctest1 --json

# 表格输出
gc pr list -R infra-test/gctest1 --format table
```

说明：
- `--paginate` 会从第一页开始连续读取多页结果，直到远端返回不足一页；不能与 `--page` 同时使用。
- `--per-page` 控制单页大小，未显式传 `--limit` 时默认每页 100；显式传 `--limit` 时会在本地截断到指定数量。
- `--commit-message` 会读取每个候选 PR 的提交列表并按提交信息子串匹配，适合从提交标题反查关联 PR。

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
- 文本输出包含里程碑信息（如果 PR 关联了里程碑）。
- `--time-format absolute|relative` 只影响文本详情和评论区中的时间展示，不改变 `--json` 结构。
- 如果 PR 详情 API 返回的 `additions`、`deletions`、`changed_files` 或 `commits` 为 0，CLI 会尝试通过 PR files/commits API 补齐统计；补齐失败时会给出 warning，但不阻断查看。
- `--json` 路径保持结构化输出，milestone、body、description、merged_at 等字段会自动包含在 JSON 中；其中 `body` 与 `description` 会基于远端返回互相补齐。

### pr issues - 查看 PR 关联的 Issues

```bash
# 查看 PR 关联的 Issues
gc pr issues 123 -R owner/repo

# 查看 PR 关联的 Issues（当前仓库）
gc pr issues 123

# 输出 JSON
gc pr issues 123 -R owner/repo --json
```

说明：
- `pr issues` 列出指定 PR 关联的 Issue 列表。
- `--json` 输出 Issue 对象数组，包含 `id`、`number`、`title`、`body`、`state`、`html_url`、`user`、`labels`、`milestone`、`created_at` 等字段；无关联 Issue 时输出 `[]`，不会混入文本提示。
- PR 不存在时返回明确错误；无关联 Issue 时文本输出提示"No linked issues found for PR #\<number\>"。

### pr comments - 查看 PR 评论

```bash
# 查看 PR 评论列表
gc pr comments 1 -R infra-test/gctest1

# 限制评论数量
gc pr comments 1 --limit 5 -R infra-test/gctest1

# JSON 输出
gc pr comments 1 -R infra-test/gctest1 --json
```

评论列表会显示 `Discussion ID`，可直接用于 `gc pr reply --discussion`。
`--json` 输出评论对象数组；无评论时输出 `[]`，不会混入文本提示。
当前 GitCode 公开 API 不支持通过 CLI 将 PR 评论标记为已解决或未解决；resolved 状态需要在 Web UI 中手动处理。
inline comment 会显示文件路径和 diff position 信息。

### pr comment - 添加 PR 评论

```bash
# 添加普通评论
gc pr comment 123 --body "This looks good" -R owner/repo

# 从文件读取评论内容
gc pr comment 123 --body-file comment.txt -R owner/repo

# 从 stdin 读取评论内容
echo "Comment from stdin" | gc pr comment 123 --body-file - -R owner/repo

# 添加行内评论 - 先获取文件路径
gc pr diff 123 -R owner/repo                        # 查看变更文件获取文件路径
gc pr comment 123 --body "代码逻辑正确" --path api/auth.go --position 1 -R owner/repo

# 输出 JSON
gc pr comment 123 --body "This looks good" --json
```

添加评论到 PR。支持普通评论和行内评论（inline comment）。

**行内评论注意事项**：
- 需要同时提供 `--path`（文件路径）和 `--position`（diff 行号）
- 文件路径必须是 diff 中显示的实际文件名（如 `test-cross-pr.txt`）
- diff 行号：新文件从 1 开始计数
- 可通过 `gc pr diff <number>` 查看变更文件列表来获取文件路径
- 如果文件名错误，会返回错误：`diff failed to be generated due to invalid params under position param`

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

# 合并后删除源分支
gc pr merge 1 -R infra-test/gctest1 --delete-branch --yes

# 合并后输出 JSON
gc pr merge 1 -R infra-test/gctest1 --yes --json
```

说明：
- `pr merge` 属于高风险写操作，默认需要确认。
- 非交互场景中显式传 `--yes`。
- `--delete-branch` 会在合并成功后调用远端分支删除接口；删除失败时命令返回失败。
- `--json` 只在合并和可选删除分支都完成后输出结构化结果，包含顶层 `number`、`merged`、`pull_request` 和可选 `deleted_branch`；不会混入文本提示。

### pr close - 关闭 PR

```bash
# 关闭 PR
gc pr close 1 -R infra-test/gctest1

# 非交互执行
gc pr close 1 -R infra-test/gctest1 --yes

# 输出 JSON
gc pr close 1 -R infra-test/gctest1 --json
```

说明：
- `pr close` 属于写操作，默认需要确认；非交互场景中显式传 `--yes`。

### pr reopen - 重开 PR

```bash
# 重开 PR
gc pr reopen 1 -R infra-test/gctest1

# 非交互执行
gc pr reopen 1 -R infra-test/gctest1 --yes

# 输出 JSON
gc pr reopen 1 -R infra-test/gctest1 --json
```

说明：
- `pr reopen` 属于写操作，默认需要确认；非交互场景中显式传 `--yes`。

### pr ready - 标记就绪状态

```bash
# 标记为就绪（取消草稿）
gc pr ready 1 -R infra-test/gctest1

# 标记为草稿
gc pr ready 1 -R infra-test/gctest1 --wip

# 非交互执行
gc pr ready 1 -R infra-test/gctest1 --ready --yes

# 输出 JSON
gc pr ready 1 -R infra-test/gctest1 --json
```

说明：
- `pr ready` 会修改 PR 草稿/就绪状态，默认需要确认；非交互场景中显式传 `--yes`。

### pr review - 评审 PR

```bash
# 评论 PR
gc pr review 1 --comment “评审意见” -R infra-test/gctest1

# 从文件读取评论
gc pr review 1 --comment-file review-notes.md -R infra-test/gctest1

# 从 stdin 读取评论
echo “评审意见” | gc pr review 1 --comment-file - -R infra-test/gctest1

# 批准 PR
gc pr review 1 --approve -R infra-test/gctest1

# 批准 PR 并附带评论
gc pr review 1 --approve --comment “LGTM” -R infra-test/gctest1

# 批准 PR 并从文件读取评论
gc pr review 1 --approve --comment-file self-check.md -R infra-test/gctest1

# 强制通过审批（管理员权限）
gc pr review 1 --approve --force -R infra-test/gctest1
```

说明：
- `--approve` 现在走 GitCode 实际可用的 `/pulls/:number/review` endpoint，不再命中错误的 `/reviews` 路径。
- `--approve --comment` 会先提交普通评论，再执行批准动作。
- `--comment-file` 支持从文件读取多行评论，使用 `-` 可从 stdin 读取。
- `--comment` 与 `--comment-file` 互斥，不能同时使用。
- Windows PowerShell 中从 stdin 传中文评论时，建议使用 UTF-8 文件或先设置 `$OutputEncoding`，详见“Windows PowerShell 命令名和 stdin”。
- GitCode 当前公开 API 不支持”request changes”动作，`--request` 会明确报错并提示改用 `--comment` 留下审查意见。

> **权限说明**: `--approve` 需要 GitCode 平台的”审批权限”，与 `gc pr merge` 的”合并权限”是两套独立权限体系。
> - 有合并权限的用户不一定有审批权限
> - PR 作者或不在审批人范围内的用户可能收到 403 Forbidden
> - 如遇权限错误，请使用 `--comment` 留下评审记录，或联系仓库管理员授予审批权限
> - `--force` 仅限管理员使用，用于强制通过审批门禁

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

# JSON 输出
gc pr edit 1 --title "新标题" -R infra-test/gctest1 --json
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
  --yes \
  --json
```

说明：
- `--source-pr` 支持两种格式：`owner/repo#number` 或完整 GitCode URL，例如 `https://gitcode.com/owner/repo/merge_requests/123`
- 命令会按原顺序逐个 cherry-pick 源 PR 的所有 commits 到目标仓库，保留提交边界
- `pr sync` clones, fetches, and pushes repositories over SSH; ensure an SSH key with access to `git@gitcode.com` is configured.
- 推送同步分支并创建目标 PR 前默认需要确认；非交互场景中显式传 `--yes`
- 新 PR 标题默认格式：`[sync] {源 PR 标题}`
- 新 PR 内容默认继承源 PR 内容并追加同步来源信息
- 如遇 cherry-pick 冲突，命令会报错并提示手动处理

---

## Release 命令 (release)

### release create - 创建 Release

```bash
# 创建 Release（建议包含 --notes 参数）
gc release create v1.0.0 -R infra-test/gctest1 --title "Version 1.0.0" --notes "Release notes"

# 从文件读取 Release Notes
gc release create v1.0.0 -R infra-test/gctest1 --title "Version 1.0.0" --notes-file CHANGELOG.md

# 创建预发布 Release
gc release create v1.0.0-beta -R infra-test/gctest1 --title "v1.0.0 Beta" --notes "Beta release" --prerelease

# 指定目标分支
gc release create v1.0.0 -R infra-test/gctest1 --title "v1.0.0" --notes "Release" --target main

# 创建后输出 JSON
gc release create v1.0.0 -R infra-test/gctest1 --title "v1.0.0" --notes "Release" --json
```

> **注意**: `--notes` 和 `--notes-file` 参数不能同时使用。
> `--draft` 当前不受 GitCode release create API 支持；CLI 会在发起远端创建前返回错误，避免创建出非草稿 Release。
> `--prerelease` 会使用 GitCode 的 `release_status=pre` 创建预发布，并在创建后回读确认状态。
> `--json` 只在成功创建后输出 release 对象；不会混入文本提示。

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

# 上传后输出 JSON
gc release upload v1.0.0 app.zip -R infra-test/gctest1 --json
```

说明：
- `--label` 参数当前不受 GitCode release upload API 支持；CLI 现在会直接报错，不再静默忽略。
- `--json` 只在所有文件上传完成后输出上传结果数组；每项包含 `name`、`path`、`size` 和 `content_type`，不会混入文本提示。

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

# 从文件读取说明
gc release edit v1.0.0 --notes-file RELEASE_NOTES.md -R infra-test/gctest1

# 标记为预发布
gc release edit v1.0.0 --prerelease true -R infra-test/gctest1

# 标记为正式发布
gc release edit v1.0.0 --prerelease false -R infra-test/gctest1

# JSON 输出
gc release edit v1.0.0 --title "New title" -R infra-test/gctest1 --json
```

说明：
- `--draft` 和 `--target` 参数当前不受 GitCode release edit API 支持，使用时会输出警告但继续执行其他修改。
- `--prerelease true` 将 release 标记为预发布状态（release_status=pre）。
- `--prerelease false` 将 release 标记为正式发布状态（release_status=latest）。
- 若只修改标题或只修改说明，API 会保留未修改字段的原始值。
- 支持包含斜杠的 tag 名称（如 `release/v1.0.0`）。

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

# JSON 输出
gc commit comments view 123 -R infra-test/gctest1 --json
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

# JSON 输出
gc commit comments list -R infra-test/gctest1 --json
```

### commit comments list-by-sha - 列出指定提交的评论

```bash
# 列出某提交的所有评论
gc commit comments list-by-sha abc123 -R infra-test/gctest1

# JSON 输出
gc commit comments list-by-sha abc123 -R infra-test/gctest1 --json
```

说明：
- `commit comments list` 和 `commit comments list-by-sha` 的 `--json` 输出评论对象数组；无评论时输出 `[]`。

---

## 标签命令 (label)

### label list - 列出标签

```bash
# 列出所有标签
gc label list -R infra-test/gctest1

# 结构化输出
gc label list -R infra-test/gctest1 --json
```

说明：
- `label list` 当前不提供 `--limit`，因为 GitCode labels list API 没有对应分页/limit 参数。

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

说明：
- `milestone list` 当前不提供 `--state` 或 `--limit`，因为 GitCode milestones list API 还没有对应筛选参数。

### milestone create - 创建里程碑

```bash
# 创建里程碑
gc milestone create "v1.0" -R infra-test/gctest1 --description "First release"
```

### milestone view - 查看里程碑

```bash
# 查看里程碑详情（包含关联 issues）
gc milestone view 1 -R infra-test/gctest1

# JSON 输出（包含 issues 数组和计数）
gc milestone view 1 -R infra-test/gctest1 --json

# 只查看里程碑元数据，不显示 issues
gc milestone view 1 -R infra-test/gctest1 --issues=false
```

说明：
- `milestone view` 默认显示里程碑关联的 issues，按状态分组（Closed/Open）。
- `--json` 输出包含 `issues` 数组、`total_issues`、`closed_issues`、`open_issues` 字段。
- `--issues=false` 只显示里程碑元数据，不获取和显示关联 issues。
- `--json` 不能与 `--web` 同时使用。

### milestone edit - 编辑里程碑

```bash
# 编辑里程碑标题
gc milestone edit 1 --title "New Title" -R infra-test/gctest1

# 编辑里程碑描述
gc milestone edit 1 --description "Updated description" -R infra-test/gctest1

# 从文件读取描述
gc milestone edit 1 --description-file milestone-desc.md -R infra-test/gctest1

# 关闭里程碑
gc milestone edit 1 --state closed -R infra-test/gctest1

# 重新打开里程碑
gc milestone edit 1 --state open -R infra-test/gctest1

# 编辑截止日期
gc milestone edit 1 --due-date "2024-12-31" -R infra-test/gctest1

# JSON 输出
gc milestone edit 1 --title "New Title" --json -R infra-test/gctest1

# 同时编辑多个字段
gc milestone edit 1 --title "v2.0" --description "Next release" --due-date "2025-01-31" -R infra-test/gctest1
```

说明：
- `milestone edit` 支持编辑标题、描述、状态和截止日期。
- `--state` 支持 `open` 和 `closed` 两个值。
- `--description-file` 从文件读取描述内容，支持多行文本。
- `--json` 输出更新后的里程碑对象。
- 至少需要提供一个编辑选项（`--title`, `--description`, `--description-file`, `--state`, `--due-date`）。

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

## Pre-commit 命令 (precommit)

`precommit` 命令组用于在提交代码前检查仓库的 pre-commit 配置与本地环境，确保提交时能正常拉起 pre-commit 检查。跨平台支持 Windows、Linux（x86/arm）、macOS。

### precommit check - 检查 pre-commit 配置与环境

检查流程：

1. 检测仓库根是否存在 `.pre-commit-config.yaml`（或 `.yml`）。无配置时视为"无需检查"，退出码 `0`。
2. 检测本地 `pre-commit` 工具是否安装。
3. 检测 git pre-commit hook 是否已初始化。
4. 可选：使用 `--run` 实际执行 `pre-commit run --all-files`。

环境缺失时，在交互式终端（stdin 为 TTY）下会自动安装并初始化；在非交互环境下需显式传 `--yes` 才会修改环境，否则仅诊断并报错。`--no-install` 表示只诊断、绝不修改环境。

```bash
# 检查环境是否就绪
gc precommit check

# 检查并实际拉起 pre-commit 检查
gc precommit check --run

# 仅诊断，不安装/初始化（不修改环境）
gc precommit check --no-install

# 非交互环境允许自动安装/初始化
gc precommit check --yes

# 机器可消费输出
gc precommit check --json
```

说明：

- 支持 `--json`：输出写入 stdout，字段为 `config_found`、`tool_installed`、`tool_version`、`hook_installed`、`actions_taken`、`run_result`、`run_output`、`ok`、`reason`、`install_failure_categories`（`run_output` 仅在 `run_result` 为 `failed` 时携带 `pre-commit run` 输出）。即使自动安装失败（退出码 `1`），`--json` 仍会输出结构化结果体（`reason=install_failed`），不会只剩退出码与 stderr 文本。
- `reason` 是稳定、机器可读的结果分类，便于脚本/agent 直接分支，取值：
  - `no_config`：仓库未配置 pre-commit（`ok=true`，属正常跳过）。
  - `tool_missing`：`pre-commit` 工具未安装（且未尝试 / 未授权安装）。
  - `hook_missing`：git pre-commit hook 未初始化。
  - `run_failed`：环境就绪但 `pre-commit run` 失败。
  - `install_failed`：已授权自动安装，但未能产出可用工具（安装尝试失败，或无可用安装器）；具体失败类型见 `install_failure_categories`。
  - `not_in_repo`：当前目录不在 git 仓库内。
  - 环境完全就绪（且 `--run` 通过或未请求）时 `reason` 省略（为空）。
- `install_failure_categories` 仅在 `reason=install_failed` 时出现，为机器可读的失败类型数组（按首次出现顺序去重）：`permission`（权限不足）/ `network`（网络失败）/ `toolchain`（缺少 Python/pip 工具链）。无法归类的失败不计入该数组（可能为空）。
- `--no-install` 与 `--yes` 互斥；hooks 本身运行失败时报"pre-commit checks failed"（区别于"环境未就绪"）。
- 退出码：`0` 就绪或无配置；`1` 环境未就绪 / 检查失败 / 非 Git 仓库 / 非交互且未授权修改环境；`2` 用法错误。
- 自动安装按工具可用性择优：`pipx` → `python3 -m pip install --user` → `python -m pip install --user`；都不可用时给出各平台手动安装指引。安装失败时按错误类型给出针对性指引（权限不足 / 网络失败 / 工具链缺失）。
- 不在 PATH 目录间复制二进制，始终在项目内调用。

---

## 其他命令

### version - 显示版本

```bash
gc version
```

Command-name note:
- When the CLI is launched as `gitcode` or `gitcode.exe`, `version`, `help`, `help --json`, `schema`, and shell completion output use `gitcode` as the command name.
- When the CLI is launched as `python -m gc_cli`, output uses `gitcode` as the command name.
- When the CLI is launched as `gc` or `gc.exe`, output continues to use `gc`.
- DEB/RPM packages install `/usr/bin/gitcode` as an alias of `/usr/bin/gc`; on Linux both commands are equivalent.

### help - 帮助

```bash
# 显示帮助
gc help

# 显示命令帮助
gc help issue
gc help issue create

# 搜索命令（按关键词搜索）
gc help --search pr
gc help --search issue

# 列出所有主题
gc help --topics

# 按主题过滤命令
gc help --topic pull-requests
gc help --topic issues
```

说明：
- `--search` 按关键词搜索命令名称、路径、描述和别名
- `--topics` 列出所有已定义的主题分类
- `--topic` 显示指定主题下的所有命令

### schema - 命令元数据

```bash
# 输出完整命令树
gc schema

# 输出单个命令的元数据
gc schema "issue view"
```

说明：
- 对带预定义取值的 flag，schema 会在 `enum` 字段中暴露合法值，例如 `format`、`time-format`、`method`、部分 `state`/`sort`/`direction` flag。

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
| `GC_HOST` | 默认 GitCode 主机（默认：gitcode.com）；必须是受信 hostname-only 值；已接入共享 host-aware 认证入口的业务命令会映射到对应 API 主机，且非默认 host 只使用该 host 的本地登录 token |
| `GC_TIMEOUT` | API 请求超时时间（默认：30s）；支持 Duration 格式如 `60s`、`2m`，或纯秒数如 `120` |
| `GC_DEBUG` | 启用 API 调试日志，输出重试、Rate Limit 等信息到 stderr |
| `GC_API_DEBUG` | 同 `GC_DEBUG`，启用 API 调试日志 |
| `NO_COLOR` | 禁用颜色输出 |

---

## 已知限制

以下功能受 GitCode API 限制，可能无法正常工作：

| 功能 | 限制说明 |
|------|----------|
| `repo fork` | 仓库路径已按用户输入解析，但 GitCode API 在部分仓库上仍可能返回 `400 Bad Request` |
| `milestone create/view` | 返回 400 错误，API 可能不支持 |
| `release delete` | GitCode API 不返回 release ID，且不支持按 tag 删除（#241） |

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

**最后更新**: 2026-05-04
