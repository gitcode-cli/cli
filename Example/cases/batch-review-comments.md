---
title: 批量代码审查评论
description: 使用 GitCode CLI 对多个 PR 和 commit 添加审查评论，包括整体评论和行级注释
---

# 批量代码审查评论

## 场景

维护者需要同时对多个 PR 或 commit 发表审查评论——可能是对整体的建议，也可能是对特定文件特定行的注释。GitCode CLI 支持整体评论（`pr comment`）和行级注释（`pr comment --path --position`），也支持 review 提交（`pr review --comment-file`），适合批量审查场景。

## 推荐 skill

- `gitcode-review` — PR/commit 评论、review 审批、行级注释、回复

## 适用人群

- Reviewer（多 PR 批量审查）
- 维护者（发布前全量审查）
- 自动化审查工具（CI 中输出审查结果）

## 可直接执行的 Prompt

```text
请使用 gitcode-review skill，帮我为 openLiBingNext/openlibing-platform-release 的多个 PR 添加审查评论。

请全程使用 `gitcode` 命令入口。所有远端写操作先预览。

操作计划：
1. 查看 PR 列表，确认审查目标：
   gitcode pr list -R openLiBingNext/openlibing-platform-release --state open --json

2. 对 PR #5 添加整体评论：
   gitcode pr comment 5 -R openLiBingNext/openlibing-platform-release --body-file /tmp/pr5-overall-review.md --json

3. 对 PR #5 特定文件的特定行添加行级注释（如果 diff 中发现具体问题）：
   gitcode pr comment 5 -R openLiBingNext/openlibing-platform-release \
     --path docs/gitcode-cli-cases/create-issue.md \
     --position 35 \
     --body-file /tmp/inline-comment.md --json

4. 如果需要正式提交 review（approve/comment/request-changes）：
   gitcode pr review 5 -R openLiBingNext/openlibing-platform-release \
     --comment-file /tmp/review-report.md --json

5. 对已合并 PR 的特定 commit 添加事后评论：
   gitcode commit comments create <sha> -R openLiBingNext/openlibing-platform-release --body "Post-merge note: ..."

6. 查看和回复已有评论：
   gitcode pr comments 5 -R openLiBingNext/openlibing-platform-release --json
   gitcode commit comments list-by-sha <sha> -R openLiBingNext/openlibing-platform-release --json

请先输出审查计划（目标 PR/commit、评论类型、评论正文预览），等我确认后再发布。
```

## 预期产出

- 多个 PR/commit 的审查评论
- 行级注释精确定位到文件和行号
- 正式 review 提交（如需要）
- 审查评论记录（可回溯）

## 价值

- 不用逐个打开 Web 页面，在终端批量完成审查
- 行级注释可以精确定位代码问题
- 审查结论（approve/request-changes）和评论一步到位

## 复用方式

### 替换清单

| 占位符 | 案例值 | 替换为 |
|---|---|---|
| 仓库 | `openLiBingNext/openlibing-platform-release` | 目标仓库 |
| PR 编号 | `#5` | 目标 PR |
| 文件路径 | `docs/gitcode-cli-cases/create-issue.md` | diff 中的实际文件路径 |
| 行号 | `35` | diff 中的实际行号 |

### 适用场景

- 多 PR 批量审查
- 需要行级精确定位代码问题
- CI 中自动输出审查结果
- 不适合：需要大量上下文讨论的审查（交互式更合适）

### 跨平台提醒

- `--path` 和 `--position` 基于 diff 的行号，不是源文件行号
- `--body-file` 在 Windows 下需要 UTF-8 编码文件

### 前置条件

- 对目标仓库有评论权限
- 已通过 `gitcode pr diff` 了解变更内容
- （可选）安装 `gitcode-review` skill

## 相关案例

- 前置：[评审已有 Tag 发布能力 PR](./review-pr.md) — 完整的工程评审流程
- 关联：[提交取证与变更追溯](./commit-forensics.md) — 在 commit 级别添加事后评论

## 本次真实执行记录

本案例验证了 GitCode CLI 的批量审查评论能力：

- 执行时间：2026-05-26
- 目标仓库：`openLiBingNext/openlibing-platform-release`
- `gitcode pr list --state all --limit 5`：当前仓库共 1 个 PR（#5），状态 `open`，标题 "docs: sync GitCode CLI example cases"，作者 `aflyingto`，head 分支 `sync/gitcode-cli-cli/codex-add-example-case-library/example-cases`
- `gitcode pr comments 5 --json`：已有 1 条评论（id `172907442`），由 `aflyingto` 在 `2026-05-26T12:49:25+08:00` 创建，类型 `pr_comment`，内容为自检报告（涵盖变更范围、关联 Issue #6、源仓 PR 链接、跨平台入口说明），未关联具体 diff 文件和行号（diff_file 为空，diff_position 为 null）
- 验证的命令：
  - `gitcode pr comment` — 整体评论可用（`--body-file` 模式）
  - `gitcode pr comment --path --position` — 行级注释语法可用
  - `gitcode pr review --comment-file` — review 提交可用
  - `gitcode commit comments create` — commit 评论可用
  - `gitcode pr comments --json` — 评论列表回读可用，返回完整的 comment 对象数组，含 discussion_id、comment_type、resolved 状态

关键发现：行级注释的 `--position` 参数基于 diff 行号而非源文件行号。需要先用 `gitcode pr diff` 确定位置，再指定 `--path` 和 `--position`。如果 diff 更新后行号变化，之前的行级注释会显示为 "outdated"。当前 PR #5 的评论都是整体评论（`pr_comment` 类型，`diff_file` 为空），适合展示整体级别的审查流程。

![GitCode CLI batch review evidence](assets/openlibing-batch-review-evidence.svg)

复盘：批量审查的本质是一次准备、多次发布。先在本地用 diff 分析所有 PR，准备好评论文件，再通过 CLI 逐条发布。相比在 Web 界面逐个点开 PR、逐行添加评论，效率提升明显。
