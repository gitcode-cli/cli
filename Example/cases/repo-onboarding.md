---
title: 新成员上手 GitCode 仓库
description: 使用 GitCode CLI 和本地检查快速生成仓库上手指南
---

# 新成员上手 GitCode 仓库

## 场景

新成员加入项目，或者外部贡献者想快速了解一个 GitCode 仓库的结构、构建方式、测试方式和贡献路径。

## 推荐 skill

- `gitcode-repo-onboarding`

## 可直接执行的 Prompt

```text
请使用 gitcode-repo-onboarding skill，帮我为 <owner/repo> 生成一份新成员上手指南。

目标分支：<branch>

请全程使用 `gitcode` 命令入口；如需下载代码，使用 SSH。不要编造构建和测试命令，只基于仓库真实文件总结。

请输出 Markdown 格式的 onboarding-guide.md 内容。
```

## 预期产出

- 一份可发给新成员的仓库上手指南。
- 明确的构建、测试、贡献路径。

## 价值

- 缩短新成员熟悉项目的时间。
- 将仓库隐性知识转化为可复用文档。
- 适合售前、交付、开源协作和团队内部培训。

## 复用方式

替换 `<owner/repo>` 和 `<branch>` 即可对任意 GitCode 仓库生成上手指南。
