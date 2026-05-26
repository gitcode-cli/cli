---
title: 整理发布平台 Issue 队列
description: 使用 GitCode CLI 对 openLiBing 发布平台 Issue 队列进行分类、优先级排序、补标签和重复项识别
---

# 整理发布平台 Issue 队列

## 场景

`openLiBingNext/openlibing-platform-release` 当前有 5 个 open issue，覆盖 JavaDoc 清理、文件传输优化、附件测试覆盖、已有 Tag 发布、发布决策可靠性。维护者需要判断哪些已经有 PR、哪些需要补标签、哪些适合新手、哪些应优先进入版本计划。

## 推荐 skill

- `gitcode-issue-triage`

## 可直接执行的 Prompt

```text
请使用 gitcode-issue-triage skill，帮我整理 openLiBingNext/openlibing-platform-release 的 Issue 队列。

范围：
- 状态：open
- 数量上限：20

请全程使用 `gitcode` 命令入口。请结合当前已知上下文重点处理：
- Issue #1：HwCloudClient JavaDoc/TODO 清理，适合 `type/chore`、`scope/common`，已有 PR #1；
- Issue #2：文件传输子系统优化，已有 PR #2，但 PR 当前可能存在 mergeable=false，需要标记风险；
- Issue #3：附件管理模块单元测试覆盖，适合 `type/test`、`scope/attachment`，已有 PR #3；
- Issue #4：已有 Tag 发布能力，已有 PR #4，属于 `enhancement`、`scope/release`；
- Issue #5：发布决策异步可靠性和制品级进度追踪，属于 `enhancement`、`scope/release`，尚需实现。

先输出 triage 计划表、推荐标签变更、推荐评论和建议关闭/保留策略，等我确认后再执行任何远端写操作。
```

## 预期产出

- 一份针对 5 个 open issue 的 triage 表。
- 标记“已有 PR 但未合并”“适合新手”“高优先级发布可靠性”的任务。
- 给出标签体系补齐建议，例如 `scope/release`、`scope/attachment`、`scope/file-transfer`、`type/test`。

## 价值

- 帮维护者把当前 open issue 与 open PR 对齐，避免 issue 已有 PR 但状态没有推进。
- 快速识别 Issue #5 这类对发布可靠性影响更高的需求。
- 为版本规划、贡献者分工和标签治理提供一份可执行清单。

## 复用方式

复用时替换仓库、当前 issue 列表和标签体系即可。建议每次版本规划前运行一次。

## 本次真实执行记录

本案例已对目标仓执行 triage 类远端操作：

- 新增标签：`type/docs`、`scope/docs`
- 更新 Issue：[#6 docs: 补充 GitCode CLI 应用案例文档](https://gitcode.com/openLiBingNext/openlibing-platform-release/issues/6)
- Issue #6 当前标签：`enhancement`、`scope/docs`、`type/docs`
- 补充评论：记录 `gitcode repo sync` 已创建 PR #5，并记录 Windows PowerShell stdin 中文编码问题的绕过方式

![GitCode CLI issue evidence](assets/openlibing-issue-evidence.svg)

复盘：triage 不是简单“加标签”。它需要把需求对象、同步 PR、发现的问题和后续验证路径串起来。对于已有标签体系较少的仓库，应克制新增标签，只补最少的类型和范围标签，避免把案例演示变成标签污染。
