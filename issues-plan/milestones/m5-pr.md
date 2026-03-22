# 里程碑 5: PR 功能

## 概述

实现 GitCode Pull Request 管理功能，包括创建、查看、列表、检出、合并和代码检视。**代码检视是本里程碑的重点功能**。

**预计工期**: 1.5 周

**依赖**: 里程碑 4 (Issue 功能)

**目标**: 用户能够通过命令行完整管理 PR，包括代码检视

---

## 任务清单

### PR-001: PR 创建

**优先级**: P0

**任务描述**:

实现 `gc pr create` 命令，创建新 PR。

**文件**:

```
pkg/cmd/pr/create/create.go
api/queries_pr.go
```

**功能**:

- 交互式创建
- 指定源分支和目标分支
- 设置标题和内容
- 草稿/WIP 标记
- 自动推送分支

**验收标准**:

- [ ] 交互式输入标题和内容
- [ ] 指定源分支和目标分支
- [ ] 支持草稿/WIP 标记
- [ ] 自动推送分支
- [ ] 显示 PR URL

**示例**:

```bash
# 交互式创建
$ gc pr create
? Title: Feature: Add new feature
? Body: [Opens editor]
? Source branch: feature-branch
? Target branch: main
? Reviewers: user1, user2
✓ Created pull request !123
https://gitcode.com/owner/repo/pull/123

# 命令行创建
$ gc pr create --title "Feature" --base main --head feature
✓ Created pull request !124

# 创建草稿
$ gc pr create --title "WIP Feature" --draft
✓ Created draft pull request !125

# 自动推送
$ gc pr create --push
✓ Pushed branch and created pull request !126
```

---

### PR-002: PR 列表

**优先级**: P0

**任务描述**:

实现 `gc pr list` 命令，列出 PRs。

**文件**:

```
pkg/cmd/pr/list/list.go
api/queries_pr.go
```

**验收标准**:

- [ ] 列出开放 PRs
- [ ] 按状态过滤
- [ ] 按作者/审核人过滤
- [ ] 按标签过滤
- [ ] 支持 JSON 输出

**示例**:

```bash
$ gc pr list
ID    TITLE                         STATUS    AUTHOR      UPDATED
!123  Feature: Add new feature      Open      user1       2 days ago
!122  Fix: Bug fix                  Draft     user2       1 week ago

$ gc pr list --state merged --author user1
$ gc pr list --reviewer user2 --limit 10
```

---

### PR-003: PR 查看

**优先级**: P0

**任务描述**:

实现 `gc pr view` 命令，查看 PR 详情。

**文件**:

```
pkg/cmd/pr/view/view.go
api/queries_pr.go
```

**验收标准**:

- [ ] 显示 PR 标题和状态
- [ ] 显示源分支和目标分支
- [ ] 显示 PR 内容
- [ ] 显示审核者、处理人
- [ ] 显示合并状态
- [ ] 显示评论

**示例**:

```bash
$ gc pr view 123
Open Feature: Add new feature
owner/repo!123 by username

feature-branch → main

Description of the PR...

Labels: enhancement
Reviewers: user1, user2
Assignees: user3

Merge status: Can be merged
Comments: 5

✓ View this PR: https://gitcode.com/owner/repo/pull/123

$ gc pr view 123 --comments
# 显示所有评论
```

---

### PR-004: PR 检出

**优先级**: P1

**任务描述**:

实现 `gc pr checkout` 命令，检出 PR 分支。

**文件**:

```
pkg/cmd/pr/checkout/checkout.go
git/checkout.go
```

**验收标准**:

- [ ] 检出 PR 源分支
- [ ] 指定本地分支名
- [ ] 自动获取分支

**示例**:

```bash
$ gc pr checkout 123
From https://gitcode.com/owner/repo
* [new ref]         refs/heads/feature-branch -> feature-branch
✓ Checked out PR !123 to branch feature-branch

$ gc pr checkout 123 --branch my-feature
✓ Checked out PR !123 to branch my-feature
```

---

### PR-005: PR 合并

**优先级**: P1

**任务描述**:

实现 `gc pr merge` 命令，合并 PR。

**文件**:

```
pkg/cmd/pr/merge/merge.go
api/queries_pr.go
```

**功能**:

- 普通合并
- Squash 合并
- Rebase 合并
- 删除源分支

**验收标准**:

- [ ] 普通合并
- [ ] Squash 合并
- [ ] Rebase 合并
- [ ] 合并后删除分支

**示例**:

```bash
$ gc pr merge 123
✓ Merged pull request !123

$ gc pr merge 123 --squash
✓ Squash merged pull request !123

$ gc pr merge 123 --delete-branch
✓ Merged pull request !123 and deleted branch feature
```

---

### PR-006: PR 关闭/重开

**优先级**: P1

**任务描述**:

实现 `gc pr close` 和 `gc pr reopen` 命令。

**文件**:

```
pkg/cmd/pr/close/close.go
pkg/cmd/pr/reopen/reopen.go
api/queries_pr.go
```

**验收标准**:

- [ ] 关闭 PR
- [ ] 重开 PR
- [ ] 添加评论

**示例**:

```bash
$ gc pr close 123
✓ Closed pull request !123

$ gc pr close 123 --comment "No longer needed"
✓ Closed pull request !123 with comment

$ gc pr reopen 123
✓ Reopened pull request !123
```

---

### PR-007: 代码检视（重点功能）

**优先级**: P0

**任务描述**:

实现 `gc pr review` 命令，对 PR 进行代码检视。**这是本里程碑的重点功能**。

**文件**:

```
pkg/cmd/pr/review/review.go
api/queries_pr.go
api/queries_review.go
```

**审核类型**:

| 类型 | 参数 | 说明 |
|------|------|------|
| 批准 | --approve | 批准 PR |
| 请求修改 | --request-changes | 请求修改，阻止合并 |
| 仅评论 | --comment | 添加评论，不影响合并状态 |

**验收标准**:

- [ ] `--approve` 批准 PR
- [ ] `--request-changes` 请求修改
- [ ] `--comment` 仅评论
- [ ] 交互式审核流程
- [ ] 请求修改和评论必须有内容

**示例**:

```bash
# 批准 PR
$ gc pr review 123 --approve
✓ Approved pull request !123

# 批准并添加评论
$ gc pr review 123 --approve --body "LGTM! Good work."
✓ Approved pull request !123 with comment

# 请求修改
$ gc pr review 123 --request-changes --body "Please fix the following issues:
1. Add unit tests
2. Update documentation"
✓ Requested changes on pull request !123

# 仅评论
$ gc pr review 123 --comment --body "I have some questions about the implementation"
✓ Added comment to pull request !123

# 交互式审核
$ gc pr review 123
? Review type: [Approve, Request changes, Comment]
? Body: [Opens editor]
✓ Submitted review for pull request !123

# 从文件读取
$ gc pr review 123 --approve --body-file review.md
✓ Approved pull request !123
```

**交互式审核流程**:

```
1. 显示 PR 基本信息
   - 标题
   - 源分支 → 目标分支
   - 作者

2. 选择审核类型
   - Approve (批准)
   - Request changes (请求修改)
   - Comment (仅评论)

3. 输入审核内容
   - 打开编辑器
   - 或从文件读取

4. 确认提交
   - 显示审核摘要
   - 确认提交
```

**最佳实践提示**:

```markdown
## 审核模板

### 总体评价
[总结性评价]

### 具体建议

#### 文件: src/main.go

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

### PR-008: PR 差异查看

**优先级**: P2

**任务描述**:

实现 `gc pr diff` 命令，查看 PR 代码变更。

**文件**:

```
pkg/cmd/pr/diff/diff.go
api/queries_pr.go
```

**验收标准**:

- [ ] 显示代码差异
- [ ] 支持分页浏览

**示例**:

```bash
$ gc pr diff 123
diff --git a/src/main.go b/src/main.go
--- a/src/main.go
+++ b/src/main.go
@@ -1,5 +1,10 @@
 package main

+import "fmt"
+
 func main() {
-    println("Hello")
+    fmt.Println("Hello, World!")
 }
```

---

### PR-009: PR 就绪标记

**优先级**: P2

**任务描述**:

实现 `gc pr ready` 命令，标记 PR 为就绪或 WIP。

**文件**:

```
pkg/cmd/pr/ready/ready.go
api/queries_pr.go
```

**验收标准**:

- [ ] 标记为就绪
- [ ] 标记为 WIP

**示例**:

```bash
$ gc pr ready 123
✓ Marked pull request !123 as ready for review

$ gc pr ready 123 --wip
✓ Marked pull request !123 as work in progress
```

---

## API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/pulls` | GET | 列出 PRs |
| `/api/v5/repos/{owner}/{repo}/pulls` | POST | 创建 PR |
| `/api/v5/repos/{owner}/{repo}/pulls/{number}` | GET | 获取 PR |
| `/api/v5/repos/{owner}/{repo}/pulls/{number}` | PATCH | 更新 PR |
| `/api/v5/repos/{owner}/{repo}/pulls/{number}` | PUT | 合并 PR |
| `/api/v5/repos/{owner}/{repo}/pulls/{number}/reviews` | POST | 提交审核 |
| `/api/v5/repos/{owner}/{repo}/pulls/{number}/comments` | POST | 添加评论 |
| `/api/v5/repos/{owner}/{repo}/pulls/{number}/files` | GET | 获取变更文件 |

---

## 依赖关系

```
PR-001 (Create) ───┐
                   │
PR-002 (List) ─────┼─→ API Client
                   │
PR-003 (View) ─────┤
                   │
PR-004 (Checkout) ─┤
                   │
PR-005 (Merge) ────┤
                   │
PR-006 (Close) ────┤
                   │
PR-007 (Review) ───┤ ← 重点
                   │
PR-008 (Diff) ─────┤
                   │
PR-009 (Ready) ────┘
```

---

## 完成标准

里程碑 M5 完成需满足：

1. ✅ `gc pr create` 创建 PR 成功
2. ✅ `gc pr list` 列出 PRs
3. ✅ `gc pr view` 显示 PR 详情
4. ✅ `gc pr checkout` 检出 PR 分支
5. ✅ `gc pr merge` 合并 PR 成功
6. ✅ **`gc pr review --approve` 批准 PR**
7. ✅ **`gc pr review --request-changes` 请求修改**
8. ✅ 单元测试覆盖率 ≥ 80%

---

## 测试用例

### 单元测试

```bash
go test ./pkg/cmd/pr/... -v
go test ./pkg/cmd/pr/review/... -v
```

### 集成测试

```bash
go test -tags=integration ./pkg/cmd/pr/... -v
```

### 手动测试清单

- [ ] 创建 PR（交互式）
- [ ] 创建 PR（命令行）
- [ ] 创建草稿 PR
- [ ] 列出 PRs（各种过滤）
- [ ] 查看 PR 详情
- [ ] 检出 PR 分支
- [ ] 合并 PR
- [ ] Squash 合并
- [ ] 关闭 PR
- [ ] **批准 PR（核心）**
- [ ] **请求修改（核心）**
- [ ] **添加评论**
- [ ] 查看差异
- [ ] 标记就绪/WIP

---

## 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 审核状态同步 | 中 | 实现状态轮询 |
| 冲突处理 | 中 | 提供解决提示 |
| 大 PR 加载慢 | 低 | 分页加载 |
| API 限流 | 中 | 实现重试机制 |

---

## 相关文档

- [issues-plan/06-module-pr.md](../06-module-pr.md)
- [gc-api-doc/doc/05-pull-requests.md](../../../gc-api-doc/doc/05-pull-requests.md)
- [gc-api-doc/test/test_pull_requests.py](../../../gc-api-doc/test/test_pull_requests.py)

---

**最后更新**: 2026-03-22