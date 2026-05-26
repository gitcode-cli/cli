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
请使用 GitCode CLI 帮我为 <owner/repo> 生成一份新成员上手指南。

要求：
1. 全程使用 `gitcode` 命令，不使用 `gc`。
2. 优先使用 `gitcode-repo-onboarding` skill；如果未安装该 skill，请按同等流程执行。
3. 先获取仓库信息：
   - gitcode repo view <owner/repo> --json
   - gitcode repo stats --branch <branch> -R <owner/repo> --json
4. 如需下载代码，默认使用 SSH：
   - ssh -T git@gitcode.com
   - gitcode repo clone <owner/repo> --git-protocol ssh
5. 进入仓库后检查：
   - README.md
   - CONTRIBUTING.md
   - AGENTS.md / CLAUDE.md
   - Makefile
   - go.mod / package.json / pyproject.toml / Cargo.toml / pom.xml
   - docs/ 和 scripts/
6. 不要编造构建命令；只能基于仓库文件推断。
7. 输出一份 Markdown 上手指南，包含：
   - 项目简介
   - 目录结构
   - 本地环境准备
   - 构建命令
   - 测试命令
   - 常见开发流程
   - 如何提交 issue 和 PR
   - 建议第一个任务

输出文件：
- onboarding-guide.md
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
