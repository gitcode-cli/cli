---
title: 同步目录到另一个仓库并创建 PR
description: 使用 GitCode CLI repo sync 将本地目录同步到目标仓库并自动创建 Pull Request
---

# 同步目录到另一个仓库并创建 PR

## 场景

平台团队维护一份公共规范、API 合约、文档或模板，需要定期同步到多个业务仓库，并在目标仓库自动创建 PR。

## 推荐 skill

- `gitcode-repo`

## 可直接执行的 Prompt

```text
请使用 gitcode-repo skill，帮我把当前仓库的 <source-dir> 同步到 <target-repo> 的 <target-dir>，并创建 Pull Request。

请全程使用 `gitcode` 命令入口；repo sync 涉及代码传输，默认使用 SSH。

输入：
- source_dir: <source-dir>
- target_repo: <target-repo>
- target_dir: <target-dir>
- base_branch: <base-branch>
- 同步目的：<说明为什么同步>

请先说明将执行的同步计划、PR 标题和正文，等我确认后再执行。
```

## 预期产出

- 目标仓库中的同步分支。
- 一个指向目标仓库的 Pull Request。

## 价值

- 适合多仓文档、合约、模板、配置同步。
- 减少人工复制文件和手动开 PR 的重复劳动。
- 使用 SSH 路径，跨 Windows/Linux 保持一致。

## 复用方式

替换源目录、目标仓库、目标目录和目标分支即可。
