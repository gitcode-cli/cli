---
title: GitCode CLI 应用案例库
description: GitCode CLI 在常见研发协作场景中的 AI 使用案例和可复用 prompt
---

# GitCode CLI 应用案例库

本目录用于沉淀 GitCode CLI 在各类业务场景中的应用案例，服务于推广宣传、用户使用指导和团队内部复用。

案例库采用总分结构：

- 本页是总入口，提供案例索引、适用场景和使用方式。
- `cases/` 目录下每个 Markdown 文件是一个独立案例，包含场景背景、推荐 skill、可直接复制执行的 prompt、预期产出、价值和复用建议。

## 使用前准备

建议先确认本地已安装最新版 GitCode CLI：

```bash
gitcode version
gitcode auth status
gitcode help --json
gitcode schema
```

跨平台命令约定：

- Windows PowerShell 推荐使用 `gitcode`，避免 `gc` 被 PowerShell 内置 `Get-Content` 别名抢占。
- Linux/macOS 下 `gitcode` 和 `gc` 均可用；案例中统一使用 `gitcode` 作为跨平台入口。
- 涉及代码下载、fork、sync、checkout 的场景默认使用 SSH，需先确认 `ssh -T git@gitcode.com` 可用。

## 可选 Skill 仓库

部分案例可配合 GitCode CLI 对外发布 skill 使用：

- 仓库地址：`git@gitcode.com:gitcode-cli/skills.git`
- Web 地址：<https://gitcode.com/gitcode-cli/skills>

案例中的 `推荐 skill` 字段用于提示 AI 客户端优先使用哪个 skill。即使未安装 skill，也可以直接复制 prompt 给 AI 执行，AI 应按 prompt 中的 GitCode CLI 命令约束完成任务。

## 案例索引

| 案例 | 适用对象 | 推荐 skill | 解决的问题 |
| --- | --- | --- | --- |
| [向指定仓库提交高质量 Issue](./cases/create-issue.md) | 产品、测试、开发、开源用户 | `gitcode-issue-create` | 把零散问题描述整理成可跟踪、可执行的 Issue。 |
| [从 Fork 分支创建 Pull Request](./cases/create-pr-from-fork.md) | 外部贡献者、协作开发者 | `gitcode-pr-create` | 在无主仓写权限时，通过 fork 分支提交贡献。 |
| [对 Pull Request 做工程评审](./cases/review-pr.md) | Reviewer、维护者、AI 评审代理 | `gitcode-pr-review` | 结构化发现风险、缺失测试和合并阻塞点。 |
| [发布版本并上传 Release 资产](./cases/publish-release.md) | 维护者、发布负责人 | `gitcode-release-helper` | 生成 release notes、创建 release、上传产物并验证。 |
| [新成员上手 GitCode 仓库](./cases/repo-onboarding.md) | 新成员、外部贡献者、售前/交付团队 | `gitcode-repo-onboarding` | 快速了解仓库、构建方式、测试方式和贡献路径。 |
| [做一次敏感信息与安全审查](./cases/security-review.md) | 安全、研发、发布负责人 | `gitcode-security-check` | 检查仓库或 PR 中的凭证泄漏、危险配置和常见代码风险。 |
| [批量整理 Issue 队列](./cases/triage-issues.md) | 项目经理、维护者、技术负责人 | `gitcode-issue-triage` | 对存量 Issue 分类、优先级排序、补标签和识别重复项。 |
| [同步目录到另一个仓库并创建 PR](./cases/sync-repo-directory.md) | 平台团队、文档团队、多仓维护者 | `gitcode-repo` | 把本仓库的一段目录内容同步到目标仓库并自动开 PR。 |
| [安装或升级后做全量冒烟验证](./cases/regression-after-install.md) | CLI 用户、发布负责人、测试人员 | `gitcode-regression` | 验证 `gitcode` 在 Windows/Linux 下的基础能力和关键命令。 |

## GitHub Pages 静态部署说明

本目录中的页面全部是 Markdown，并包含 GitHub Pages/Jekyll 可识别的 front matter。启用 GitHub Pages 后，可以直接将仓库根目录作为 Pages 来源，访问：

```text
https://<github-org>.github.io/<github-repo>/Example/
```

如 GitHub 镜像仓库使用 `main` 分支根目录发布，`Example/index.md` 会作为案例库入口页，`Example/cases/*.md` 会作为独立案例页面渲染。

## 维护原则

- 案例优先面向真实业务任务，而不是罗列命令参数。
- Prompt 必须可复制、可执行、可替换占位符。
- 案例中涉及代码下载、同步、PR checkout 的路径默认使用 SSH。
- 案例中不得包含真实 token、密码、私钥、内部敏感仓库名或不可公开数据。
- 如果 GitCode CLI 命令能力变化，应同步更新对应案例。
