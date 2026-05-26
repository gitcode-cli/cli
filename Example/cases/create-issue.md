---
title: 向指定仓库提交高质量 Issue
description: 使用 GitCode CLI 和 gitcode-issue-create skill 将问题描述整理并提交到指定仓库
---

# 向指定仓库提交高质量 Issue

## 场景

用户发现某个仓库存在 bug、体验问题或功能缺口，希望把零散描述提交为一个结构清晰、可跟踪、可分派的 GitCode Issue。

## 推荐 skill

- `gitcode-issue-create`

## 适用人群

- 产品经理提交需求
- 测试人员提交缺陷
- 开源用户反馈问题
- AI 代理帮助用户整理 issue

## 可直接执行的 Prompt

```text
请使用 gitcode-issue-create skill，帮我向 <owner/repo> 提交一个高质量 Issue。

请全程使用 `gitcode` 命令入口；如果信息不足，先问我。

我的原始描述：
<在这里粘贴问题、需求、复现步骤、截图说明、日志摘要等>

请先给出 issue 预览，等我确认后再创建。
```

## 预期产出

- 一个标题清晰、正文完整、标签合理的 GitCode Issue。
- 可追溯的重复搜索结果。
- 可复用的 issue 模板。

## 价值

- 降低用户提交 issue 的门槛。
- 减少维护者反复追问背景、复现步骤和验收标准的成本。
- 提升 issue 后续进入开发、评审和交付流程的质量。

## 复用方式

将 `<owner/repo>` 和原始描述替换为自己的仓库和问题即可复用。
