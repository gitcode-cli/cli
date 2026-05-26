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
请使用 GitCode CLI 帮我从 fork 仓库向上游仓库创建 Pull Request。

上下文：
- 上游仓库：<upstream-owner/upstream-repo>
- fork 仓库：<my-owner/my-fork-repo>
- 工作分支：<feature-branch>
- 目标分支：<base-branch，例如 main>

要求：
1. 全程使用 `gitcode` 命令，不使用 `gc`。
2. 优先使用 `gitcode-pr-create` skill；如果未安装该 skill，请按同等流程执行。
3. 所有代码下载、fetch、push 路径默认使用 SSH，先确认：
   - ssh -T git@gitcode.com
4. 创建 PR 前检查本地状态：
   - git status --short --branch
   - git remote -v
   - git log --oneline --decorate -5
   - git diff --stat <base-branch>...HEAD
5. 检查远端已有 PR，避免重复：
   - gitcode pr list -R <upstream-owner/upstream-repo> --head <feature-branch> --base <base-branch> --json
6. 生成 PR 描述文件，包含：
   - Summary
   - Verification
   - Risk
   - Related Issue
7. 使用 `--body-file` 创建 PR：
   - gitcode pr create -R <upstream-owner/upstream-repo> --fork <my-owner/my-fork-repo> --head <feature-branch> --base <base-branch> --title "<title>" --body-file <file> --json

我的变更说明：
<粘贴变更背景、测试结果、关联 issue、风险说明>

输出：
- 当前分支和远端检查结果
- PR 标题和正文预览
- 创建成功后的 PR 编号和链接
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
