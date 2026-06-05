---
title: 创建 Issue 并关联提交 PR（端到端）
description: 使用 GitCode CLI 将已完成的分支变更整理为 Issue，再从 fork 仓库向目标仓库提交关联 PR
---

# 创建 Issue 并关联提交 PR（端到端）

## 场景

开发者在 fork 仓库完成代码变更后，需要将改动提交到上游仓库。完整流程是：先在上游仓库创建 Issue 描述变更背景和影响范围，再从 fork 分支提交 PR 并关联该 Issue，让 Reviewer 能快速理解改动动机和上下文。

本案例以 PyTorch NPU CI 工作流重构为例：在 fork 仓库 `kerer-sk/pytorch` 的 `v2.7.1_fast_test` 分支完成代码变更后，向上游 `Ascend/pytorch` 的 `v2.7.1` 分支提交 PR。

## 推荐 skill

- `gitcode-issue-create` — 撰写并创建高质量 Issue
- `gitcode-pr-create` — 从 fork 分支向目标仓库创建 PR

以上 skill 来自 [gitcode-cli/skills](https://gitcode.com/gitcode-cli/skills) 项目（`git@gitcode.com:gitcode-cli/skills.git`），可独立安装使用

## 适用人群

- 没有主仓写权限的外部贡献者
- 跨组织协作开发者
- 需要将变更动机和代码一起提交审阅的团队
- AI 代理辅助完成 Issue → PR 全流程

## 可直接执行的 Prompt

```text
请帮我完成以下操作：

1. 使用 gitcode-issue-create skill，在 <目标仓库> 创建 Issue，描述当前分支的变更内容。
2. 使用 gitcode-pr-create skill，从 <fork 仓库>:<源分支> 向 <目标仓库>:<目标分支> 创建 PR，关联上一步创建的 Issue。

上下文：
- 目标仓库：Ascend/pytorch
- fork 仓库：kerer-sk/pytorch
- 源分支：v2.7.1_fast_test
- 目标分支：v2.7.1

我的变更说明：
- 精简 CI 工作流中的 setup-npu-test-env 复合动作，去除冗余的 checkout、pip cache、PyTorch 安装步骤
- 统一工作流中的仓库引用为 Ascend/pytorch
- 涉及 12 个文件，+772/-559 行

请全程使用 `gitcode` 命令入口。先给出 Issue 预览，确认后再创建；Issue 创建后自动生成 PR 预览，确认后再提交。
```

## 预期产出

- 在目标仓库创建一个带有 `infra-tooling` 标签的 enhancement Issue
- 从 fork 仓库向目标仓库指定分支创建 PR，描述中通过 `Fixes #<issue>` 关联 Issue
- Issue 和 PR 均可直接回溯，形成完整的变更链路

## 价值

- Reviewer 在 PR 中可以直接跳转到 Issue 了解变更背景和验收标准
- Issue 作为变更的"为什么"，PR 作为"改了什么"，职责分离但互相引用
- 统一 `gitcode` CLI 路径，避免在 Web 表单中手动填写容易遗漏的字段

## 复用方式

### 替换清单

| 占位符 | 案例值 | 替换为 |
|---|---|---|
| 目标仓库 | `Ascend/pytorch` | 你的上游主仓 |
| fork 仓库 | `kerer-sk/pytorch` | 你的 fork 仓库 |
| 源分支 | `v2.7.1_fast_test` | 你的开发分支 |
| 目标分支 | `v2.7.1` | 目标仓库的合入分支（可能是 `master`/`main`） |
| Issue 类型 | `enhancement` | `bug` / `feature` / `enhancement` |
| 标签 | `infra-tooling` | 目标仓库实际标签名 |

### 适用场景

- 变更范围清晰，能用简短段落描述动机和影响
- 需要先在上游仓库创建 Issue 作为 PR 的关联上下文
- fork 仓库已推送到远端，目标分支明确
- 不适合：目标仓库默认不接受 fork PR、变更太小无需单独 Issue

### 前置条件

- `gitcode auth status` 确认已登录
- fork 仓库已在 GitCode 上创建，开发分支已推送
- 对目标仓库有 Issue 创建权限
- 了解目标仓库的标签体系和分支命名规范

## 本次真实执行记录

本案例在 `Ascend/pytorch` 上完成端到端 Issue → PR 链路：

- **Issue**: [#2136 enhance(ci): 替换工作流引用为 Ascend/pytorch 并精简 setup-npu-test-env 动作](https://gitcode.com/Ascend/pytorch/issues/2136)
- **PR**: [#36940 refactor(ci): 替换工作流引用为 Ascend/pytorch 并精简 setup-npu-test-env 动作](https://gitcode.com/Ascend/pytorch/merge_requests/36940)
- **源分支**: `kerer-sk/pytorch:v2.7.1_fast_test`
- **目标分支**: `Ascend/pytorch:v2.7.1`

![image.png](https://raw.gitcode.com/user-images/assets/9483585/8048a759-b3b0-456d-9d5f-e08afb504a76/image.png 'image.png')
![image.png](https://raw.gitcode.com/user-images/assets/9483585/e94f32ca-c349-48c3-aa85-345bd68480c8/image.png 'image.png')
![image.png](https://raw.gitcode.com/user-images/assets/9483585/d7af2559-d77c-47fa-b819-e11990073af1/image.png 'image.png')

### Issue 创建

```bash
gitcode issue create -R Ascend/pytorch \
  --title "enhance(ci): 替换工作流引用为 Ascend/pytorch 并精简 setup-npu-test-env 动作" \
  --label "infra-tooling" \
  --body "..."
```

结果：`Created issue #2136 in Ascend/pytorch`

### PR 创建

```bash
gitcode pr create -R Ascend/pytorch \
  --head kerer-sk/pytorch:v2.7.1_fast_test \
  --base v2.7.1 \
  --title "refactor(ci): 替换工作流引用为 Ascend/pytorch 并精简 setup-npu-test-env 动作" \
  --body "..."
```

结果：`Created PR #36940 in Ascend/pytorch`

### 关键经验

1. **Issue 先于 PR**：Issue 描述"为什么改"，PR 描述"改了什么"，两者通过 `Fixes #issue` 关联形成完整链路。
2. **变更预览确认**：Issue 和 PR 创建前都应先展示预览，确认标题、描述、标签、目标分支无误后再提交。
3. **标签要先查**：`gitcode label list -R <repo>` 查看目标仓库实际标签，避免使用不存在的标签名。
4. **跨仓库 PR 用 --head**：从 fork 提交 PR 时用 `--head <fork-owner>/<repo>:<branch>` 指定源分支。

## 相关案例

- 前置：[向发布平台提交高质量 Issue](./create-issue.md) — 单独创建 Issue 的详细流程
- 前置：[从 Fork 分支创建发布平台 PR](./create-pr-from-fork.md) — 单独创建 fork PR 的详细流程
- 后续：[评审 PR](./review-pr.md) — PR 创建后的代码审查流程
