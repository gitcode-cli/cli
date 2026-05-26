---
title: 从 Fork 分支创建 Pull Request
description: 使用 GitCode CLI 从 fork 仓库向上游仓库创建 Pull Request
---

# 从 Fork 分支创建 Pull Request

## 场景

外部贡献者没有上游仓库写权限，需要通过 fork 仓库提交分支，并向上游仓库创建 Pull Request。

## 推荐 skill

- `gitcode-pr-create`
- 可辅助使用：`gitcode-repo`

## 可直接执行的 Prompt

```text
请使用 gitcode-pr-create skill，帮我从 fork 仓库向上游仓库创建 Pull Request。

上下文：
- 上游仓库：<upstream-owner/upstream-repo>
- fork 仓库：<my-owner/my-fork-repo>
- 工作分支：<feature-branch>
- 目标分支：<base-branch，例如 main>

请全程使用 `gitcode` 命令入口；涉及代码传输默认使用 SSH。

我的变更说明：
<粘贴变更背景、测试结果、关联 issue、风险说明>

请先给出 PR 标题和正文预览，等我确认后再创建。
```

## 预期产出

- 从 fork 分支发起到上游仓库的 Pull Request。
- 带验证证据和风险说明的 PR 描述。

## 价值

- 降低外部贡献者参与门槛。
- 规避 Windows PowerShell 下 `gc` 命令别名冲突。
- 统一 fork PR 创建路径，减少维护者沟通成本。

## 复用方式

替换上游仓库、fork 仓库、工作分支和目标分支即可用于任意 GitCode 仓库。
