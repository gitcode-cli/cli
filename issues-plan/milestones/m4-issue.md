# 里程碑 4: Issue 功能

## 概述

实现 GitCode Issue 管理功能，包括创建、查看、列表、关闭、重开和评论 Issue。

**预计工期**: 1 周

**依赖**: 里程碑 3 (仓库功能)

**目标**: 用户能够通过命令行管理 GitCode Issue

---

## 任务清单

### ISSUE-001: Issue 创建

**优先级**: P0

**任务描述**:

实现 `gc issue create` 命令，创建新 Issue。

**文件**:

```
pkg/cmd/issue/create/create.go
api/queries_issue.go
```

**功能**:

- 交互式创建
- 指定标题和内容
- 添加标签、指派者、里程碑
- Web 模式

**验收标准**:

- [ ] 交互式输入标题和内容
- [ ] 命令行参数指定
- [ ] 从文件读取内容
- [ ] 添加标签
- [ ] 指派处理人
- [ ] 显示 Issue URL

**示例**:

```bash
# 交互式创建
$ gc issue create
? Title: Bug: Something is wrong
? Body: [Opens editor]
? Labels: bug, priority-high
? Assignees: user1
✓ Created issue #123
https://gitcode.com/owner/repo/issues/123

# 命令行创建
$ gc issue create --title "Bug report" --body "Description" --label bug
✓ Created issue #124

# 从文件读取
$ gc issue create --title "Feature" --body-file feature.md
✓ Created issue #125
```

---

### ISSUE-002: Issue 列表

**优先级**: P0

**任务描述**:

实现 `gc issue list` 命令，列出 Issues。

**文件**:

```
pkg/cmd/issue/list/list.go
api/queries_issue.go
```

**功能**:

- 状态过滤
- 标签过滤
- 作者/处理人过滤
- 搜索
- 分页

**验收标准**:

- [ ] 列出开放 Issues
- [ ] 按状态过滤
- [ ] 按标签过滤
- [ ] 按作者过滤
- [ ] 支持 JSON 输出

**示例**:

```bash
$ gc issue list
ID    TITLE                        LABELS        ASSIGNEES   UPDATED
#123  Bug: Something is wrong      bug           user1       2 days ago
#122  Feature request              enhancement   user2       1 week ago

# 过滤
$ gc issue list --state closed --label bug
$ gc issue list --author user1 --limit 10
```

---

### ISSUE-003: Issue 查看

**优先级**: P0

**任务描述**:

实现 `gc issue view` 命令，查看 Issue 详情。

**文件**:

```
pkg/cmd/issue/view/view.go
api/queries_issue.go
```

**验收标准**:

- [ ] 显示 Issue 标题和状态
- [ ] 显示 Issue 内容
- [ ] 显示标签、处理人
- [ ] 显示评论
- [ ] 在浏览器中打开

**示例**:

```bash
$ gc issue view 123
Open Bug: Something is wrong
owner/repo#123 opened by username

Description of the issue...

Labels: bug, priority-high
Assignees: user1
Milestone: v1.0

✓ View this issue: https://gitcode.com/owner/repo/issues/123
```

---

### ISSUE-004: Issue 关闭/重开

**优先级**: P1

**任务描述**:

实现 `gc issue close` 和 `gc issue reopen` 命令。

**文件**:

```
pkg/cmd/issue/close/close.go
pkg/cmd/issue/reopen/reopen.go
api/queries_issue.go
```

**验收标准**:

- [ ] 关闭 Issue
- [ ] 重开 Issue
- [ ] 添加评论
- [ ] 显示操作确认

**示例**:

```bash
$ gc issue close 123
✓ Closed issue #123

$ gc issue close 123 --comment "Fixed in commit abc"
✓ Closed issue #123 with comment

$ gc issue reopen 123
✓ Reopened issue #123
```

---

### ISSUE-005: Issue 评论

**优先级**: P1

**任务描述**:

实现 `gc issue comment` 命令，添加评论。

**文件**:

```
pkg/cmd/issue/comment/comment.go
api/queries_issue.go
```

**验收标准**:

- [ ] 添加评论
- [ ] 从文件读取
- [ ] 编辑器编辑
- [ ] 显示评论 URL

**示例**:

```bash
$ gc issue comment 123 --body "This is a comment"
✓ Added comment to issue #123
https://gitcode.com/owner/repo/issues/123#note_456

$ gc issue comment 123 --editor
# 打开编辑器
```

---

### ISSUE-006: Issue 编辑

**优先级**: P2

**任务描述**:

实现 `gc issue edit` 命令，编辑 Issue。

**文件**:

```
pkg/cmd/issue/edit/edit.go
api/queries_issue.go
```

**验收标准**:

- [ ] 编辑标题
- [ ] 编辑内容
- [ ] 添加/移除标签
- [ ] 添加/移除指派者

**示例**:

```bash
$ gc issue edit 123 --title "New title"
✓ Updated issue #123

$ gc issue edit 123 --add-label confirmed --remove-label needs-triage
✓ Updated issue #123
```

---

### ISSUE-007: 标签管理

**优先级**: P2

**任务描述**:

支持 Issue 标签管理。

**文件**:

```
pkg/cmd/label/label.go
api/queries_label.go
```

**验收标准**:

- [ ] 列出标签
- [ ] 创建标签
- [ ] 编辑标签
- [ ] 删除标签

---

### ISSUE-008: 里程碑管理

**优先级**: P2

**任务描述**:

支持 Issue 里程碑管理。

**文件**:

```
pkg/cmd/milestone/milestone.go
api/queries_milestone.go
```

**验收标准**:

- [ ] 列出里程碑
- [ ] 创建里程碑
- [ ] 编辑里程碑
- [ ] 删除里程碑

---

## API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/issues` | GET | 列出 Issues |
| `/api/v5/repos/{owner}/{repo}/issues` | POST | 创建 Issue |
| `/api/v5/repos/{owner}/{repo}/issues/{number}` | GET | 获取 Issue |
| `/api/v5/repos/{owner}/{repo}/issues/{number}` | PATCH | 更新 Issue |
| `/api/v5/repos/{owner}/{repo}/issues/{number}/comments` | POST | 添加评论 |
| `/api/v5/repos/{owner}/{repo}/labels` | GET | 列出标签 |
| `/api/v5/repos/{owner}/{repo}/milestones` | GET | 列出里程碑 |

---

## 依赖关系

```
ISSUE-001 (Create) ─┐
                    │
ISSUE-002 (List) ───┼─→ API Client
                    │
ISSUE-003 (View) ───┤
                    │
ISSUE-004 (Close) ──┤
                    │
ISSUE-005 (Comment)─┤
                    │
ISSUE-006 (Edit) ───┤
                    │
ISSUE-007 (Labels) ─┤
                    │
ISSUE-008 (Milestones)┘
```

---

## 完成标准

里程碑 M4 完成需满足：

1. ✅ `gc issue create` 创建 Issue 成功
2. ✅ `gc issue list` 列出 Issues
3. ✅ `gc issue view` 显示 Issue 详情
4. ✅ `gc issue close/reopen` 操作成功
5. ✅ `gc issue comment` 添加评论成功
6. ✅ 单元测试覆盖率 ≥ 80%

---

## 测试用例

### 单元测试

```bash
go test ./pkg/cmd/issue/... -v
```

### 集成测试

```bash
go test -tags=integration ./pkg/cmd/issue/... -v
```

### 手动测试清单

- [ ] 创建 Issue（交互式）
- [ ] 创建 Issue（命令行）
- [ ] 列出 Issues（各种过滤）
- [ ] 查看 Issue 详情
- [ ] 查看 Issue 评论
- [ ] 关闭 Issue
- [ ] 重开 Issue
- [ ] 添加评论
- [ ] 编辑 Issue

---

## 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 长内容处理 | 低 | 支持文件和编辑器 |
| 标签不存在 | 中 | 提示创建或自动创建 |
| API 限流 | 中 | 实现分页和重试 |

---

## 相关文档

- [issues-plan/05-module-issue.md](../05-module-issue.md)
- [gc-api-doc/doc/04-issues.md](https://gitcode.com/afly-infra/gc-api-doc/blob/main/doc/04-issues.md)
- [gc-api-doc/test/test_issues.py](https://gitcode.com/afly-infra/gc-api-doc/blob/main/test/test_issues.py)

---

**最后更新**: 2026-03-22