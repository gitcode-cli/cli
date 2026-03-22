# MR 模块需求

本文档详细描述 gitcode-cli MR（Merge Request）模块的功能需求、验收标准和 API 映射。

## 模块概述

MR 模块提供 GitCode 仓库 Merge Request 的管理功能，包括创建、查看、列表、检出、合并、关闭和代码检视。**代码检视是本模块的重点功能。**

### 命令结构

```
gc mr <command>

Commands:
  create         Create a merge request
  list           List merge requests
  view           View a merge request
  checkout       Checkout a merge request
  merge          Merge a merge request
  close          Close a merge request
  reopen         Reopen a closed merge request
  review         Review a merge request (重点功能)
  diff           View changes in a merge request
  ready          Mark a merge request as ready
```

### MR 标识格式

| 格式 | 示例 | 描述 |
|------|------|------|
| 数字 | `123` | MR 编号 |
| URL | `https://gitcode.com/owner/repo/merge_requests/123` | 完整 URL |
| 分支名 | `feature-branch` | 源分支名 |
| 当前分支 | - | 自动检测 |

---

## MR-001: mr create - 创建 MR

### 功能描述

创建新的 Merge Request。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --title | -t | string | | MR 标题 |
| --body | -b | string | | MR 内容 |
| --body-file | -F | string | | 从文件读取内容 |
| --base | -B | string | | 目标分支 |
| --head | -H | string | | 源分支 |
| --draft | -d | bool | false | 创建为草稿 |
| --wip | -w | bool | false | 标记为 WIP |
| --label | -l | []string | | 添加标签 |
| --assignee | -a | []string | | 指派处理人 |
| --reviewer | -r | []string | | 指派审核人 |
| --milestone | -m | string | | 关联里程碑 |
| --push | | bool | false | 自动推送分支 |
| --fill | | bool | false | 自动填充标题和内容 |
| --template | -T | string | | 使用模板 |
| --web | | bool | false | 在浏览器中创建 |

### 使用示例

```bash
# 交互式创建
gc mr create

# 指定标题和内容
gc mr create --title "Feature: Add new feature" --body "Description"

# 创建为草稿
gc mr create --title "WIP Feature" --draft

# 指定分支
gc mr create --base main --head feature-branch

# 自动推送并创建
gc mr create --push

# 自动填充标题
gc mr create --fill

# 在浏览器中创建
gc mr create --web
```

### 验收标准

- [ ] 支持交互式输入标题和内容
- [ ] 支持命令行参数指定
- [ ] 支持草稿/WIP 标记
- [ ] 支持自动推送分支
- [ ] 支持自动填充标题（从提交）
- [ ] 显示创建成功的 MR URL

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/pulls` | POST | 创建 MR |

### 测试用例映射

- 参考 `gc-api-doc/test/test_pull_requests.py`

---

## MR-002: mr list - 列出 MRs

### 功能描述

列出仓库的 Merge Requests。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --state | -s | string | open | 状态过滤 (open/closed/merged/all) |
| --author | -A | string | | 作者过滤 |
| --assignee | -a | string | | 处理人过滤 |
| --reviewer | -r | string | | 审核人过滤 |
| --label | -l | []string | | 标签过滤 |
| --source-branch | | string | | 源分支过滤 |
| --target-branch | | string | | 目标分支过滤 |
| --search | -S | string | | 搜索关键字 |
| --limit | -L | int | 30 | 最大数量 |
| --draft | | bool | false | 只显示草稿 |
| --json | | []string | | JSON 输出 |
| --web | -w | bool | false | 在浏览器中查看 |

### 使用示例

```bash
# 列出开放的 MRs
gc mr list

# 列出已合并的 MRs
gc mr list --state merged

# 按作者过滤
gc mr list --author username

# 按标签过滤
gc mr list --label priority-high

# 只显示草稿
gc mr list --draft
```

### 输出示例

```
ID    TITLE                         STATUS    AUTHOR      UPDATED
!123  Feature: Add new feature      Open      user1       2 days ago
!122  Fix: Bug fix                  Draft     user2       1 week ago
```

### 验收标准

- [ ] 正确列出 MRs
- [ ] 支持状态过滤
- [ ] 支持作者/处理人/审核人过滤
- [ ] 支持标签过滤
- [ ] 支持分支过滤
- [ ] 支持草稿过滤
- [ ] 支持 JSON 输出

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/pulls` | GET | 列出 MRs |

### 测试用例映射

- 参考 `gc-api-doc/test/test_pull_requests.py`

---

## MR-003: mr view - 查看 MR

### 功能描述

查看 Merge Request 的详细信息。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --web | -w | bool | false | 在浏览器中查看 |
| --comments | -c | bool | false | 显示评论 |
| --json | | []string | | JSON 输出 |

### 使用示例

```bash
# 查看 MR
gc mr view 123

# 查看 MR 和评论
gc mr view 123 --comments

# 通过 URL 查看
gc mr view https://gitcode.com/owner/repo/merge_requests/123

# 查看当前分支关联的 MR
gc mr view

# 在浏览器中查看
gc mr view 123 --web
```

### 输出示例

```
Open Feature: Add new feature
owner/repo!123 by username

feature-branch → main

Description of the MR...

Labels: enhancement
Reviewers: user1, user2
Assignees: user3

Merge status: Can be merged
Comments: 5

View this MR on GitCode: https://gitcode.com/owner/repo/merge_requests/123
```

### 验收标准

- [ ] 正确显示 MR 标题和状态
- [ ] 显示源分支和目标分支
- [ ] 显示 MR 内容
- [ ] 显示审核者、处理人
- [ ] 显示合并状态
- [ ] 支持显示评论
- [ ] 支持 JSON 输出

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/pulls/{number}` | GET | 获取 MR |
| `/api/v5/repos/{owner}/{repo}/pulls/{number}/comments` | GET | 获取评论 |

### 测试用例映射

- 参考 `gc-api-doc/test/test_pull_requests.py`

---

## MR-004: mr checkout - 检出 MR

### 功能描述

检出 MR 的源分支到本地。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --branch | -b | string | | 本地分支名 |
| --detach | | bool | false | 以 detached HEAD 模式检出 |
| --fetch | | bool | true | 从远程获取分支 |

### 使用示例

```bash
# 检出 MR
gc mr checkout 123

# 检出到指定分支名
gc mr checkout 123 --branch my-feature

# 获取并检出
gc mr checkout 123 --fetch
```

### 验收标准

- [ ] 正确检出 MR 分支
- [ ] 支持指定本地分支名
- [ ] 支持自动获取分支
- [ ] 显示检出成功的确认信息

### API 端点

此命令主要使用 Git 操作，获取 MR 信息使用：
| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/pulls/{number}` | GET | 获取 MR 信息 |

### 测试用例映射

- 参考 `gc-api-doc/test/test_pull_requests.py`

---

## MR-005: mr merge - 合并 MR

### 功能描述

合并 Merge Request。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --message | -m | string | | 合并提交消息 |
| --squash | | bool | false | Squash 合并 |
| --squash-message | | string | | Squash 消息 |
| --delete-branch | -d | bool | false | 合并后删除源分支 |
| --auto | | bool | false | 只在检查通过时合并 |
| --rebase | | bool | false | Rebase 合并 |

### 使用示例

```bash
# 合并 MR
gc mr merge 123

# 合并并添加消息
gc mr merge 123 --message "Merge feature"

# Squash 合并
gc mr merge 123 --squash

# 合并后删除分支
gc mr merge 123 --delete-branch
```

### 验收标准

- [ ] 支持普通合并
- [ ] 支持 Squash 合并
- [ ] 支持 Rebase 合并
- [ ] 支持删除源分支
- [ ] 显示合并成功的确认信息

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/pulls/{number}/merge` | PUT | 合并 MR |

### 测试用例映射

- 参考 `gc-api-doc/test/test_pull_requests.py`

---

## MR-006: mr close - 关闭 MR

### 功能描述

关闭 Merge Request（不合并）。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --comment | -c | string | | 关闭时添加评论 |

### 使用示例

```bash
# 关闭 MR
gc mr close 123

# 关闭并添加评论
gc mr close 123 --comment "No longer needed"
```

### 验收标准

- [ ] 正确关闭 MR
- [ ] 支持添加关闭评论
- [ ] 显示关闭成功的确认信息

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/pulls/{number}` | PATCH | 更新 MR 状态 |

### 测试用例映射

- 参考 `gc-api-doc/test/test_pull_requests.py`

---

## MR-007: mr reopen - 重开 MR

### 功能描述

重新打开已关闭的 Merge Request。

### 使用示例

```bash
# 重开 MR
gc mr reopen 123
```

### 验收标准

- [ ] 正确重开 MR
- [ ] 显示重开成功的确认信息

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/pulls/{number}` | PATCH | 更新 MR 状态 |

---

## MR-008: mr review - 代码检视（重点功能）

### 功能描述

对 Merge Request 进行代码检视，支持批准、请求修改或仅评论。

### 审核类型

| 类型 | 参数 | 说明 |
|------|------|------|
| 批准 | --approve | 批准 MR |
| 请求修改 | --request-changes | 请求修改，阻止合并 |
| 仅评论 | --comment | 添加评论，不影响合并状态 |

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --approve | -a | bool | false | 批准 MR |
| --request-changes | -r | bool | false | 请求修改 |
| --comment | -c | bool | false | 仅评论 |
| --body | -b | string | | 审核内容 |
| --body-file | -F | string | | 从文件读取内容 |

### 使用示例

```bash
# 交互式审核
gc mr review 123

# 批准 MR
gc mr review 123 --approve

# 批准并添加评论
gc mr review 123 --approve --body "LGTM!"

# 请求修改
gc mr review 123 --request-changes --body "Please fix the tests"

# 仅评论
gc mr review 123 --comment --body "I have some questions"

# 从文件读取审核内容
gc mr review 123 --approve --body-file review.md
```

### 交互式审核流程

1. 显示 MR 基本信息（标题、分支）
2. 选择审核类型（批准/请求修改/仅评论）
3. 输入审核内容（Markdown 编辑器）
4. 确认提交

### 验收标准

- [ ] 支持 --approve 批准 MR
- [ ] 支持 --request-changes 请求修改
- [ ] 支持 --comment 仅评论
- [ ] 支持交互式审核流程
- [ ] 显示审核成功的确认信息
- [ ] 请求修改和评论必须有内容

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/pulls/{number}/reviews` | POST | 提交审核 |
| `/api/v5/repos/{owner}/{repo}/pulls/{number}/comments` | POST | 添加评论 |

### 测试用例映射

- 参考 `gc-api-doc/test/test_pull_requests.py`

### 审核最佳实践

```markdown
## 总体评价
[总结性评价]

## 具体建议

### 文件: src/main.go

**第 45 行:**
```go
// 建议修改
```

### 优点
- [列出优点]

### 需要改进
- [列出需要改进的地方]
```

---

## MR-009: mr diff - 查看 MR 差异

### 功能描述

查看 Merge Request 的代码变更。

### 使用示例

```bash
# 查看 MR 差异
gc mr diff 123
```

### 验收标准

- [ ] 正确显示代码差异
- [ ] 支持分页浏览

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/pulls/{number}/files` | GET | 获取变更文件 |

### 测试用例映射

- 参考 `gc-api-doc/test/test_pull_requests.py`

---

## MR-010: mr ready - 标记为就绪/WIP

### 功能描述

将 MR 标记为就绪审核或 WIP（Work In Progress）。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --wip | | bool | false | 标记为 WIP |

### 使用示例

```bash
# 标记为就绪
gc mr ready 123

# 标记为 WIP
gc mr ready 123 --wip
```

### 验收标准

- [ ] 正确更新 MR 状态
- [ ] 显示成功的确认信息

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/pulls/{number}` | PATCH | 更新 MR |

---

## 相关文档

- [gc-api-doc/doc/05-pull-requests.md](../../gc-api-doc/doc/05-pull-requests.md)
- [gc-api-doc/test/test_pull_requests.py](../../gc-api-doc/test/test_pull_requests.py)

---

**最后更新**: 2026-03-22