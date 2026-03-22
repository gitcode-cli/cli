# Issue 模块需求

本文档详细描述 gitcode-cli Issue 模块的功能需求、验收标准和 API 映射。

## 模块概述

Issue 模块提供 GitCode 仓库 Issue 的管理功能，包括创建、查看、列表、关闭、重开和评论。

### 命令结构

```
gc issue <command>

Commands:
  create    Create a new issue
  list      List issues in a repository
  view      View an issue
  close     Close an issue
  reopen    Reopen a closed issue
  comment   Add a comment to an issue
  edit      Edit an issue
```

### Issue 标识格式

| 格式 | 示例 | 描述 |
|------|------|------|
| 数字 | `123` | Issue 编号 |
| URL | `https://gitcode.com/owner/repo/issues/123` | 完整 URL |
| 当前分支 | - | 自动关联 |

---

## ISSUE-001: issue create - 创建 Issue

### 功能描述

创建新的 Issue。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --title | -t | string | | Issue 标题 |
| --body | -b | string | | Issue 内容 |
| --body-file | -F | string | | 从文件读取内容 |
| --label | -l | []string | | 添加标签 |
| --assignee | -a | []string | | 指派处理人 |
| --milestone | -m | string | | 关联里程碑 |
| --web | -w | bool | false | 在浏览器中创建 |

### 使用示例

```bash
# 交互式创建
gc issue create

# 指定标题和内容
gc issue create --title "Bug report" --body "Description here"

# 添加标签
gc issue create --title "Bug" --label bug,priority-high

# 指派处理人
gc issue create --title "Task" --assignee user1,user2

# 从文件读取内容
gc issue create --title "Feature" --body-file description.md

# 在浏览器中创建
gc issue create --web
```

### 验收标准

- [ ] 支持交互式输入标题和内容
- [ ] 支持命令行参数指定
- [ ] 支持从文件读取内容
- [ ] 支持添加标签
- [ ] 支持指派处理人
- [ ] 支持关联里程碑
- [ ] 显示创建成功的 Issue URL

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/issues` | POST | 创建 Issue |

### 测试用例映射

- 参考 `gc-api-doc/test/test_issues.py`

---

## ISSUE-002: issue list - 列出 Issues

### 功能描述

列出仓库的 Issues。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --state | -s | string | open | 状态过滤 (open/closed/all) |
| --label | -l | []string | | 标签过滤 |
| --assignee | -a | string | | 处理人过滤 |
| --author | -A | string | | 作者过滤 |
| --milestone | -m | string | | 里程碑过滤 |
| --search | -S | string | | 搜索关键字 |
| --limit | -L | int | 30 | 最大数量 |
| --page | -p | int | 1 | 页码 |
| --json | | []string | | JSON 输出 |
| --web | -w | bool | false | 在浏览器中查看 |

### 使用示例

```bash
# 列出开放的 Issues
gc issue list

# 列出已关闭的 Issues
gc issue list --state closed

# 按标签过滤
gc issue list --label bug,priority-high

# 按处理人过滤
gc issue list --assignee username

# 搜索 Issues
gc issue list --search "bug fix"

# 限制数量
gc issue list --limit 10
```

### 输出示例

```
ID    TITLE                        LABELS           ASSIGNEES    UPDATED
#123  Bug: Something is wrong      bug, high        user1        2 days ago
#122  Feature request              enhancement      user2        1 week ago
```

### 验收标准

- [ ] 正确列出 Issues
- [ ] 支持状态过滤
- [ ] 支持标签过滤
- [ ] 支持处理人过滤
- [ ] 支持作者过滤
- [ ] 支持搜索
- [ ] 支持分页
- [ ] 支持 JSON 输出

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/issues` | GET | 列出 Issues |

### 测试用例映射

- 参考 `gc-api-doc/test/test_issues.py`
- 参考 `gc-api-doc/test/test_labels.py`（标签）

---

## ISSUE-003: issue view - 查看 Issue

### 功能描述

查看 Issue 的详细信息。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --web | -w | bool | false | 在浏览器中查看 |
| --comments | -c | bool | false | 显示评论 |
| --json | | []string | | JSON 输出 |

### 使用示例

```bash
# 查看 Issue
gc issue view 123

# 查看 Issue 和评论
gc issue view 123 --comments

# 通过 URL 查看
gc issue view https://gitcode.com/owner/repo/issues/123

# 在浏览器中查看
gc issue view 123 --web
```

### 输出示例

```
Open Bug: Something is wrong
owner/repo#123 opened by username

Description of the issue...

Labels: bug, priority-high
Assignees: user1, user2
Milestone: v1.0
Comments: 5

View this issue on GitCode: https://gitcode.com/owner/repo/issues/123
```

### 验收标准

- [ ] 正确显示 Issue 标题和状态
- [ ] 显示 Issue 内容
- [ ] 显示标签、处理人、里程碑
- [ ] 支持显示评论
- [ ] 支持 JSON 输出

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/issues/{number}` | GET | 获取 Issue |
| `/api/v5/repos/{owner}/{repo}/issues/{number}/comments` | GET | 获取评论 |

### 测试用例映射

- 参考 `gc-api-doc/test/test_issues.py`

---

## ISSUE-004: issue close - 关闭 Issue

### 功能描述

关闭指定的 Issue。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --comment | -c | string | | 关闭时添加评论 |

### 使用示例

```bash
# 关闭 Issue
gc issue close 123

# 关闭并添加评论
gc issue close 123 --comment "Fixed in commit abc123"
```

### 验收标准

- [ ] 正确关闭 Issue
- [ ] 支持添加关闭评论
- [ ] 显示关闭成功的确认信息

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/issues/{number}` | PATCH | 更新 Issue 状态 |
| `/api/v5/repos/{owner}/{repo}/issues/{number}/comments` | POST | 添加评论 |

### 测试用例映射

- 参考 `gc-api-doc/test/test_issues.py`

---

## ISSUE-005: issue reopen - 重开 Issue

### 功能描述

重新打开已关闭的 Issue。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --comment | -c | string | | 重开时添加评论 |

### 使用示例

```bash
# 重开 Issue
gc issue reopen 123

# 重开并添加评论
gc issue reopen 123 --comment "Issue still exists"
```

### 验收标准

- [ ] 正确重开 Issue
- [ ] 支持添加评论
- [ ] 显示重开成功的确认信息

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/issues/{number}` | PATCH | 更新 Issue 状态 |

### 测试用例映射

- 参考 `gc-api-doc/test/test_issues.py`

---

## ISSUE-006: issue comment - 添加评论

### 功能描述

向 Issue 添加评论。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --body | -b | string | | 评论内容 |
| --body-file | -F | string | | 从文件读取内容 |
| --editor | -e | bool | false | 在编辑器中编辑 |
| --web | -w | bool | false | 在浏览器中评论 |

### 使用示例

```bash
# 添加评论
gc issue comment 123 --body "This is a comment"

# 从文件读取
gc issue comment 123 --body-file comment.md

# 在编辑器中编辑
gc issue comment 123 --editor
```

### 验收标准

- [ ] 支持命令行输入评论
- [ ] 支持从文件读取
- [ ] 支持在编辑器中编辑
- [ ] 显示添加成功的确认信息

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/issues/{number}/comments` | POST | 添加评论 |

### 测试用例映射

- 参考 `gc-api-doc/test/test_issues.py`

---

## ISSUE-007: issue edit - 编辑 Issue

### 功能描述

编辑 Issue 的标题、内容、标签等信息。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --title | -t | string | | 新标题 |
| --body | -b | string | | 新内容 |
| --body-file | -F | string | | 从文件读取内容 |
| --add-label | | []string | | 添加标签 |
| --remove-label | | []string | | 移除标签 |
| --add-assignee | | []string | | 添加处理人 |
| --remove-assignee | | []string | | 移除处理人 |

### 使用示例

```bash
# 编辑标题
gc issue edit 123 --title "New title"

# 编辑内容
gc issue edit 123 --body "New description"

# 添加标签
gc issue edit 123 --add-label bug,confirmed

# 移除标签
gc issue edit 123 --remove-label needs-triage
```

### 验收标准

- [ ] 支持编辑标题
- [ ] 支持编辑内容
- [ ] 支持添加/移除标签
- [ ] 支持添加/移除处理人
- [ ] 显示编辑成功的确认信息

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/issues/{number}` | PATCH | 更新 Issue |

### 测试用例映射

- 参考 `gc-api-doc/test/test_issues.py`

---

## 相关文档

- [gc-api-doc/doc/04-issues.md](../../gc-api-doc/doc/04-issues.md)
- [gc-api-doc/test/test_issues.py](../../gc-api-doc/test/test_issues.py)
- [gc-api-doc/test/test_labels.py](../../gc-api-doc/test/test_labels.py)
- [gc-api-doc/test/test_milestones.py](../../gc-api-doc/test/test_milestones.py)

---

**最后更新**: 2026-03-22