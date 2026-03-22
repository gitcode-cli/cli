# 里程碑 5: MR 功能

## 概述

实现 GitCode Merge Request 管理功能，包括创建、查看、列表、检出、合并和代码检视。**代码检视是本里程碑的重点功能**。

**预计工期**: 1.5 周

**依赖**: 里程碑 4 (Issue 功能)

**目标**: 用户能够通过命令行完整管理 MR，包括代码检视

---

## 任务清单

### MR-001: MR 创建

**优先级**: P0

**任务描述**:

实现 `gc mr create` 命令，创建新 MR。

**文件**:

```
pkg/cmd/mr/create/create.go
api/queries_mr.go
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
- [ ] 显示 MR URL

**示例**:

```bash
# 交互式创建
$ gc mr create
? Title: Feature: Add new feature
? Body: [Opens editor]
? Source branch: feature-branch
? Target branch: main
? Reviewers: user1, user2
✓ Created merge request !123
https://gitcode.com/owner/repo/merge_requests/123

# 命令行创建
$ gc mr create --title "Feature" --base main --head feature
✓ Created merge request !124

# 创建草稿
$ gc mr create --title "WIP Feature" --draft
✓ Created draft merge request !125

# 自动推送
$ gc mr create --push
✓ Pushed branch and created merge request !126
```

---

### MR-002: MR 列表

**优先级**: P0

**任务描述**:

实现 `gc mr list` 命令，列出 MRs。

**文件**:

```
pkg/cmd/mr/list/list.go
api/queries_mr.go
```

**验收标准**:

- [ ] 列出开放 MRs
- [ ] 按状态过滤
- [ ] 按作者/审核人过滤
- [ ] 按标签过滤
- [ ] 支持 JSON 输出

**示例**:

```bash
$ gc mr list
ID    TITLE                         STATUS    AUTHOR      UPDATED
!123  Feature: Add new feature      Open      user1       2 days ago
!122  Fix: Bug fix                  Draft     user2       1 week ago

$ gc mr list --state merged --author user1
$ gc mr list --reviewer user2 --limit 10
```

---

### MR-003: MR 查看

**优先级**: P0

**任务描述**:

实现 `gc mr view` 命令，查看 MR 详情。

**文件**:

```
pkg/cmd/mr/view/view.go
api/queries_mr.go
```

**验收标准**:

- [ ] 显示 MR 标题和状态
- [ ] 显示源分支和目标分支
- [ ] 显示 MR 内容
- [ ] 显示审核者、处理人
- [ ] 显示合并状态
- [ ] 显示评论

**示例**:

```bash
$ gc mr view 123
Open Feature: Add new feature
owner/repo!123 by username

feature-branch → main

Description of the MR...

Labels: enhancement
Reviewers: user1, user2
Assignees: user3

Merge status: Can be merged
Comments: 5

✓ View this MR: https://gitcode.com/owner/repo/merge_requests/123

$ gc mr view 123 --comments
# 显示所有评论
```

---

### MR-004: MR 检出

**优先级**: P1

**任务描述**:

实现 `gc mr checkout` 命令，检出 MR 分支。

**文件**:

```
pkg/cmd/mr/checkout/checkout.go
git/checkout.go
```

**验收标准**:

- [ ] 检出 MR 源分支
- [ ] 指定本地分支名
- [ ] 自动获取分支

**示例**:

```bash
$ gc mr checkout 123
From https://gitcode.com/owner/repo
* [new ref]         refs/heads/feature-branch -> feature-branch
✓ Checked out MR !123 to branch feature-branch

$ gc mr checkout 123 --branch my-feature
✓ Checked out MR !123 to branch my-feature
```

---

### MR-005: MR 合并

**优先级**: P1

**任务描述**:

实现 `gc mr merge` 命令，合并 MR。

**文件**:

```
pkg/cmd/mr/merge/merge.go
api/queries_mr.go
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
$ gc mr merge 123
✓ Merged merge request !123

$ gc mr merge 123 --squash
✓ Squash merged merge request !123

$ gc mr merge 123 --delete-branch
✓ Merged merge request !123 and deleted branch feature
```

---

### MR-006: MR 关闭/重开

**优先级**: P1

**任务描述**:

实现 `gc mr close` 和 `gc mr reopen` 命令。

**文件**:

```
pkg/cmd/mr/close/close.go
pkg/cmd/mr/reopen/reopen.go
api/queries_mr.go
```

**验收标准**:

- [ ] 关闭 MR
- [ ] 重开 MR
- [ ] 添加评论

**示例**:

```bash
$ gc mr close 123
✓ Closed merge request !123

$ gc mr close 123 --comment "No longer needed"
✓ Closed merge request !123 with comment

$ gc mr reopen 123
✓ Reopened merge request !123
```

---

### MR-007: 代码检视（重点功能）

**优先级**: P0

**任务描述**:

实现 `gc mr review` 命令，对 MR 进行代码检视。**这是本里程碑的重点功能**。

**文件**:

```
pkg/cmd/mr/review/review.go
api/queries_mr.go
api/queries_review.go
```

**审核类型**:

| 类型 | 参数 | 说明 |
|------|------|------|
| 批准 | --approve | 批准 MR |
| 请求修改 | --request-changes | 请求修改，阻止合并 |
| 仅评论 | --comment | 添加评论，不影响合并状态 |

**验收标准**:

- [ ] `--approve` 批准 MR
- [ ] `--request-changes` 请求修改
- [ ] `--comment` 仅评论
- [ ] 交互式审核流程
- [ ] 请求修改和评论必须有内容

**示例**:

```bash
# 批准 MR
$ gc mr review 123 --approve
✓ Approved merge request !123

# 批准并添加评论
$ gc mr review 123 --approve --body "LGTM! Good work."
✓ Approved merge request !123 with comment

# 请求修改
$ gc mr review 123 --request-changes --body "Please fix the following issues:
1. Add unit tests
2. Update documentation"
✓ Requested changes on merge request !123

# 仅评论
$ gc mr review 123 --comment --body "I have some questions about the implementation"
✓ Added comment to merge request !123

# 交互式审核
$ gc mr review 123
? Review type: [Approve, Request changes, Comment]
? Body: [Opens editor]
✓ Submitted review for merge request !123

# 从文件读取
$ gc mr review 123 --approve --body-file review.md
✓ Approved merge request !123
```

**交互式审核流程**:

```
1. 显示 MR 基本信息
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

### MR-008: MR 差异查看

**优先级**: P2

**任务描述**:

实现 `gc mr diff` 命令，查看 MR 代码变更。

**文件**:

```
pkg/cmd/mr/diff/diff.go
api/queries_mr.go
```

**验收标准**:

- [ ] 显示代码差异
- [ ] 支持分页浏览

**示例**:

```bash
$ gc mr diff 123
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

### MR-009: MR 就绪标记

**优先级**: P2

**任务描述**:

实现 `gc mr ready` 命令，标记 MR 为就绪或 WIP。

**文件**:

```
pkg/cmd/mr/ready/ready.go
api/queries_mr.go
```

**验收标准**:

- [ ] 标记为就绪
- [ ] 标记为 WIP

**示例**:

```bash
$ gc mr ready 123
✓ Marked merge request !123 as ready for review

$ gc mr ready 123 --wip
✓ Marked merge request !123 as work in progress
```

---

## API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/pulls` | GET | 列出 MRs |
| `/api/v5/repos/{owner}/{repo}/pulls` | POST | 创建 MR |
| `/api/v5/repos/{owner}/{repo}/pulls/{number}` | GET | 获取 MR |
| `/api/v5/repos/{owner}/{repo}/pulls/{number}` | PATCH | 更新 MR |
| `/api/v5/repos/{owner}/{repo}/pulls/{number}` | PUT | 合并 MR |
| `/api/v5/repos/{owner}/{repo}/pulls/{number}/reviews` | POST | 提交审核 |
| `/api/v5/repos/{owner}/{repo}/pulls/{number}/comments` | POST | 添加评论 |
| `/api/v5/repos/{owner}/{repo}/pulls/{number}/files` | GET | 获取变更文件 |

---

## 依赖关系

```
MR-001 (Create) ───┐
                   │
MR-002 (List) ─────┼─→ API Client
                   │
MR-003 (View) ─────┤
                   │
MR-004 (Checkout) ─┤
                   │
MR-005 (Merge) ────┤
                   │
MR-006 (Close) ────┤
                   │
MR-007 (Review) ───┤ ← 重点
                   │
MR-008 (Diff) ─────┤
                   │
MR-009 (Ready) ────┘
```

---

## 完成标准

里程碑 M5 完成需满足：

1. ✅ `gc mr create` 创建 MR 成功
2. ✅ `gc mr list` 列出 MRs
3. ✅ `gc mr view` 显示 MR 详情
4. ✅ `gc mr checkout` 检出 MR 分支
5. ✅ `gc mr merge` 合并 MR 成功
6. ✅ **`gc mr review --approve` 批准 MR**
7. ✅ **`gc mr review --request-changes` 请求修改**
8. ✅ 单元测试覆盖率 ≥ 80%

---

## 测试用例

### 单元测试

```bash
go test ./pkg/cmd/mr/... -v
go test ./pkg/cmd/mr/review/... -v
```

### 集成测试

```bash
go test -tags=integration ./pkg/cmd/mr/... -v
```

### 手动测试清单

- [ ] 创建 MR（交互式）
- [ ] 创建 MR（命令行）
- [ ] 创建草稿 MR
- [ ] 列出 MRs（各种过滤）
- [ ] 查看 MR 详情
- [ ] 检出 MR 分支
- [ ] 合并 MR
- [ ] Squash 合并
- [ ] 关闭 MR
- [ ] **批准 MR（核心）**
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
| 大 MR 加载慢 | 低 | 分页加载 |
| API 限流 | 中 | 实现重试机制 |

---

## 相关文档

- [issues-plan/06-module-mr.md](../06-module-mr.md)
- [gc-api-doc/doc/05-pull-requests.md](../../../gc-api-doc/doc/05-pull-requests.md)
- [gc-api-doc/test/test_pull_requests.py](../../../gc-api-doc/test/test_pull_requests.py)

---

**最后更新**: 2026-03-22