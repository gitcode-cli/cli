---
title: Issue 实现前评审
description: 开发前使用 GitCode CLI 对 Issue 做完整性检查，确认需求清晰、验收标准可测、信息完备
---

# Issue 实现前评审

## 场景

开发者在认领 Issue 后直接开始写代码，经常发现做到一半需求不清晰、验收标准无法量化、或者 Issue 描述的 API 端点实际不存在。这个案例展示如何在动手之前用 GitCode CLI 系统化审查 Issue 的就绪状态。

## 推荐 skill

- `gitcode-issue-review` — 来自 [gitcode-cli/skills](https://gitcode.com/gitcode-cli/skills) 项目（`git@gitcode.com:gitcode-cli/skills.git`），可独立安装使用

## 适用人群

- 开发者（认领 Issue 后、动手前）
- 技术负责人（评估 Issue 是否可分配）
- AI 代理（辅助分析 Issue 就绪度）

## 可直接执行的 Prompt

```text
请使用 gitcode-issue-review skill，帮我在实现前评审 openLiBingNext/openlibing-platform-release 的 Issue #5。

请全程使用 `gitcode` 命令入口。

Issue #5 主题：发布决策异步可靠性和制品级 release_result 进度追踪。

请重点检查：
1. 需求是否完整：背景、目标、期望行为是否清晰？
2. 验收标准是否可测：能否写出具体的测试用例？
3. 关联信息是否充分：PR 链接、标签、里程碑是否已设置？
4. 技术可行性：涉及 release_result 表、异步任务、Jenkins 回调，是否需要 API 变更？
5. 缺失信息：是否有需要向提交者追问的？

请输出 Issue 评审报告，包含就绪度判断（Ready / Needs Clarification / Blocked）和建议的下一步。
```

## 预期产出

- 一份 Issue 就绪度评审报告
- 列出缺失信息清单（需要向提交者追问的）
- 建议的就绪度判断和下一步行动

## 价值

- 减少"开发到一半才发现需求不清楚"的返工
- 把隐性假设和模糊约束显式化
- 帮助维护者判断 Issue 是否已具备分配条件

## 复用方式

### 替换清单

| 占位符 | 案例值 | 替换为 |
|---|---|---|
| 仓库 | `openLiBingNext/openlibing-platform-release` | 目标仓库 |
| Issue 编号 | `#5` | 待分析的 Issue |
| 技术上下文 | release_result 表、Jenkins 回调 | 目标 Issue 涉及的技术模块 |

### 适用场景

- 复杂 Issue（涉及数据库变更、API 变更、多模块联动）
- Bug Issue 需要确认可复现性和根因
- 不适合：文档修正、配置变更等简单 Issue

### 跨平台提醒

- 无特殊跨平台差异

### 前置条件

- 对目标仓库有读权限
- 了解目标项目的技术栈和模块结构
- （可选）安装 `gitcode-issue-review` skill

## 相关案例

- 前置：[向发布平台提交高质量 Issue](./create-issue.md) — 确保 Issue 已按要求创建
- 后续：[从 Fork 分支创建发布平台 PR](./create-pr-from-fork.md) — Issue 确认就绪后开始实现
- 关联：[整理发布平台 Issue 队列](./triage-issues.md) — 批量 triage 后对高优先级 Issue 做深度分析

## 本次真实执行记录

本案例对 `openLiBingNext/openlibing-platform-release` 的 Issue #5 执行了实现前评审：

- 评审时间：2026-05-26
- Issue #5：发布决策异步可靠性和制品级 release_result 进度追踪
- `gitcode issue view 5 --comments --json`：Issue 状态 open，标签 enhancement，无里程碑，无指派人，comments=0，创建于 2026-05-20，无更新记录
- `gitcode issue prs 5 --json`：无关联 PR（返回空数组）
- 就绪度判断：**Needs Clarification** — Issue 描述了目标（发布结果追踪可视化）和方案 A 的实现方向（~50 行改动，不改表结构），验收标准清晰（4 条 checklist），但缺少里程碑分配、无指派人、无关联 PR，说明尚未进入开发阶段
- 缺失信息：里程碑未分配（无法确定版本归属）、无指派人（未明确由谁实现）、方案 A 被描述为"最小补丁"但未评估方案 B 是否存在

关键建议：在开始实现前，需要补充里程碑分配（建议关联到 v1.1 发布结果追踪增强），并将验收标准从"status=1 预插记录"细化为具体的测试场景（并发发布时的记录隔离、超时后的状态兜底）。

![GitCode CLI issue review evidence](assets/openlibing-issue-review-evidence.svg)

复盘：实现前评审的价值在于把隐性假设（"不改表结构就够"）变成显式约束（确认方案 A 是最终选择、确认无 API 变更需求）。这一步花 15 分钟能省下数小时的方向偏差。
