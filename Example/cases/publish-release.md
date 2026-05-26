---
title: 发布版本并上传 Release 资产
description: 使用 GitCode CLI 生成发布说明、创建 release、上传资产并完成发布验证
---

# 发布版本并上传 Release 资产

## 场景

项目维护者需要发布一个新版本，包含 release notes、tag 校验、release 创建、资产上传和发布后验证。

## 推荐 skill

- `gitcode-release-helper`
- 可辅助使用：`gitcode-release`

## 可直接执行的 Prompt

```text
请使用 gitcode-release-helper skill，帮我为 <owner/repo> 发布版本 <version>。

请全程使用 `gitcode` 命令入口；发布前先给我 release notes、资产清单和验证计划预览，等我确认后再创建 release。

输入：
- previous_tag: <previous_tag>
- version: <version>
- target_branch: <target-branch>
- asset_files: <asset list>
- 本版本重点变化：<粘贴变更摘要>
```

## 预期产出

- 一个 GitCode Release。
- 一份面向用户的发布说明。
- 上传完成并可下载验证的 release assets。

## 价值

- 降低发布流程漏步骤风险。
- 将 release notes、资产、验证记录统一留痕。
- 方便发布负责人和 AI 代理复用相同流程。

## 复用方式

替换版本号、前一版本 tag、目标分支和资产文件即可复用。
