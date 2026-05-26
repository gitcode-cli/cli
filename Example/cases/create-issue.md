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
请使用 GitCode CLI 帮我向 <owner/repo> 提交一个高质量 Issue。

要求：
1. 全程使用 `gitcode` 命令，不使用 `gc`，确保 Windows PowerShell 和 Linux 体验一致。
2. 优先使用 `gitcode-issue-create` skill；如果未安装该 skill，请按同等流程执行。
3. 先检查认证状态：
   - gitcode version
   - gitcode auth status
4. 先搜索 open 和 closed issue，避免重复：
   - gitcode issue list -R <owner/repo> --state all --search "<关键词>" --json
5. 查询仓库已有标签：
   - gitcode label list -R <owner/repo> --json
6. 根据我提供的信息判断 issue 类型：bug / feature / enhancement / docs / question。
7. 如果关键信息缺失，请先向我追问，不要编造。
8. 生成 Markdown 格式 issue 正文，保存为临时文件后使用 `--body-file` 创建。
9. 创建前先使用 `--dry-run --json` 预览。
10. 我确认后再执行真实创建：
    - gitcode issue create -R <owner/repo> --title "<title>" --body-file <file> --label <labels> --json

我的原始描述：
<在这里粘贴问题、需求、复现步骤、截图说明、日志摘要等>

输出：
- 重复 issue 搜索结果摘要
- 建议标题
- 建议标签
- issue 正文预览
- 创建成功后的 issue 编号和链接
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

将 `<owner/repo>`、`<关键词>` 和原始描述替换为自己的仓库和问题即可复用。
