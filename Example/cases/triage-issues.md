---
title: 批量整理 Issue 队列
description: 使用 GitCode CLI 对 Issue 队列进行分类、优先级排序、补标签和重复项识别
---

# 批量整理 Issue 队列

## 场景

项目维护者需要整理长期积压的 Issue，识别重复问题、缺信息问题、高优先级问题，并补齐标签和评论。

## 推荐 skill

- `gitcode-issue-triage`

## 可直接执行的 Prompt

```text
请使用 gitcode-issue-triage skill，帮我整理 <owner/repo> 的 Issue 队列。

范围：
- 状态：open
- 数量上限：<limit>

请全程使用 `gitcode` 命令入口。先输出 triage 计划表、推荐标签变更和推荐评论，等我确认后再执行任何远端写操作。
```

## 预期产出

- 一份 issue 队列整理计划。
- 经确认后的标签和评论变更。
- 重复项和缺信息项列表。

## 价值

- 让维护者快速掌握 backlog 状态。
- 降低 issue 队列长期失控的风险。
- 为版本规划和人力分配提供依据。

## 复用方式

替换仓库和 `limit` 即可用于不同项目的定期 issue 运营。
