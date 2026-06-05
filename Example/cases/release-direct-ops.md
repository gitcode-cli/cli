---
title: Release 直接操作
description: 不通过 release-helper，直接用 gitcode release 命令管理 release 的查看、编辑、资产上传下载和删除
---

# Release 直接操作

## 场景

`gitcode-release-helper` 适合规划型发布流程，但有时维护者只需要快速操作：查看已有 release、下载某次发布的资产、编辑 release 说明、删除错误的 release。这个案例展示 release 命令族的直接操作。

## 推荐 skill

- `gitcode-release` — 来自 [gitcode-cli/skills](https://gitcode.com/gitcode-cli/skills) 项目（`git@gitcode.com:gitcode-cli/skills.git`），可独立安装使用

## 适用人群

- 维护者（日常 release 管理）
- 测试人员（下载历史版本进行测试）
- DevOps（自动化 release 资产管理）

## 可直接执行的 Prompt

```text
请使用 gitcode-release skill，帮我直接管理 openLiBingNext/openlibing-platform-release 的 release。

请全程使用 `gitcode` 命令入口。删除操作必须先 dry-run。

操作计划：
1. 列出所有 release：
   gitcode release list -R openLiBingNext/openlibing-platform-release --json

2. 查看特定 release 详情，确认资产和状态：
   gitcode release view v0.0.0-gitcode-cli-case-demo -R openLiBingNext/openlibing-platform-release --json

3. 下载 release 资产（如文档或 jar）：
   gitcode release download v0.0.0-gitcode-cli-case-demo -R openLiBingNext/openlibing-platform-release --dir ./downloads

4. 如需更新 release 说明：
   gitcode release edit v0.0.0-gitcode-cli-case-demo -R openLiBingNext/openlibing-platform-release --body-file /tmp/release-notes.md --json

5. 如需上传额外资产：
   gitcode release upload v0.0.0-gitcode-cli-case-demo -R openLiBingNext/openlibing-platform-release ./new-asset.jar

6. 如要删除错误 release（危险操作）：
   gitcode release delete v0.0.0-gitcode-cli-case-demo -R openLiBingNext/openlibing-platform-release --dry-run
   # 确认后再执行：
   gitcode release delete v0.0.0-gitcode-cli-case-demo -R openLiBingNext/openlibing-platform-release --yes --json

请先展示当前 release 列表，等我确认后再执行写操作。
```

## 预期产出

- 当前所有 release 的清单
- 指定 release 的详情（tag、状态、资产列表）
- 下载的资产文件（如执行下载）
- 编辑/上传/删除操作的确认结果

## 价值

- 快速查看和下载历史版本，不需要通过 Web 界面
- 对错误 release 可以快速清理
- 适合在 CI 中自动化 release 资产管理

## 复用方式

### 替换清单

| 占位符 | 案例值 | 替换为 |
|---|---|---|
| 仓库 | `openLiBingNext/openlibing-platform-release` | 目标仓库 |
| Tag | `v0.0.0-gitcode-cli-case-demo` | 目标 release tag |
| 下载目录 | `./downloads` | 你的本地目录 |

### 适用场景

- 下载历史版本的构建产物进行测试
- 更新 release notes 和资产
- 清理错误或过期的 release
- 不适合：需要完整发布规划和 release notes 生成的场景（使用 release-helper）

### 跨平台提醒

- `release download` 资产路径使用目标平台的分隔符
- `release delete` 不可逆，务必先 `--dry-run`

### 前置条件

- 对目标仓库有 release 管理权限
- （可选）安装 `gitcode-release` skill

## 相关案例

- 关联：[发布 openLiBing 发布平台版本](./publish-release.md) — release-helper 的规划式发布
- 关联：[对发布平台仓库做 CLI 冒烟验证](./regression-after-install.md) — 验证 release 相关命令可用

## 本次真实执行记录

本案例验证了 `openLiBingNext/openlibing-platform-release` 的 release 直接操作命令：

- 执行时间：2026-05-26
- 当前 release 数量：1
- 演示 release：`v0.0.0-gitcode-cli-case-demo`
- Release 详情：tag_name `v0.0.0-gitcode-cli-case-demo`，name "GitCode CLI case demo release"，target_commitish `3886ca261fcd7bf2fd62d45b9f6a9caa558a0106`，draft `false`，prerelease `false`，created_at `2026-05-26T12:56:31+08:00`，author `aflyingto`。Release body 包含关联 Issue #6 和 PR !5 的说明。资产共 5 个：`v0.0.0-gitcode-cli-case-demo.zip`、`.tar.gz`、`.tar.bz2`、`.tar`（自动生成的源码归档）以及 `index.md`（手动上传的文档资产）。
- 验证的命令：
  - `gitcode release list --json` — 可用，返回数组，当前 1 个 release
  - `gitcode release view <tag> --json` — 可用，返回 tag、状态、资产列表（含 browser_download_url）
  - `gitcode release download <tag>` — 可下载单个资产
  - `gitcode release upload <tag> <file>` — 可上传追加资产
  - `gitcode release edit <tag> --body-file` — 可更新 release notes
  - `gitcode release delete <tag> --dry-run` — dry-run 返回成功

关键提醒：`release delete` 不可逆！务必先 `--dry-run` 确认。另外 release 创建后 `draft`/`prerelease` 字段可能与预期不一致，用 `release view --json` 回读确认。本例中 release body 末尾有一个 `\r\n`（CRLF），来自 `--body-file` 的源文件编码，在 Windows 环境下需注意。

![GitCode CLI release direct ops evidence](assets/openlibing-release-direct-ops-evidence.svg)

复盘：`release-helper` 覆盖的是规划型发布，而直接操作覆盖的是维护型操作。两种场景互补：规划用 helper，维护用直接命令。
