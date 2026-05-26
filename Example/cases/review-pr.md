---
title: 对 Pull Request 做工程评审
description: 使用 GitCode CLI 查看 PR、diff、评论并输出结构化工程评审结论
---

# 对 Pull Request 做工程评审

## 场景

维护者或 Reviewer 需要对一个 GitCode Pull Request 做独立工程评审，重点关注缺陷、回归、安全风险、测试缺口和合并阻塞点。

## 推荐 skill

- `gitcode-pr-review`
- 可辅助使用：`gitcode-review`

## 可直接执行的 Prompt

```text
请使用 gitcode-pr-review skill，对 <owner/repo> 的 PR #<number> 做一次工程评审。

请全程使用 `gitcode` 命令入口。重点看行为回归、安全风险、缺失测试、文档同步和合并阻塞点。

请先输出评审报告，不要直接 approve；如果需要在 PR 上发表评论，先给我预览。
```

## 预期产出

- 一份结构化 PR 评审报告。
- 可选的 GitCode PR review 评论。

## 价值

- 让评审从“泛泛看过”变成可审计的工程判断。
- 帮助团队统一 blocker 和非 blocker 的分级。
- 减少漏测、漏文档、漏安全风险。

## 复用方式

替换 `<owner/repo>` 和 `<number>` 即可用于任意 GitCode PR。
