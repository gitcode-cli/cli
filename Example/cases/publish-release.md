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
请使用 GitCode CLI 帮我为 <owner/repo> 发布版本 <version>。

要求：
1. 全程使用 `gitcode` 命令，不使用 `gc`。
2. 优先使用 `gitcode-release-helper` skill；如果未安装该 skill，请按同等流程执行。
3. 先检查已有 release 和最近合并内容：
   - gitcode release list -R <owner/repo> --json
   - gitcode pr list -R <owner/repo> --state merged --limit 50 --json
   - git log <previous_tag>..HEAD --oneline --no-merges
4. 生成 RELEASE_NOTES.md，结构包括：
   - Added
   - Changed
   - Fixed
   - Security
   - Verification
   - Upgrade Notes
5. 确认 tag、target branch 和资产文件列表。
6. 创建 release：
   - gitcode release create <version> -R <owner/repo> --title "<version>" --notes-file RELEASE_NOTES.md --target <target-branch> --json
7. 上传资产：
   - gitcode release upload <version> <asset-files...> -R <owner/repo> --json
8. 发布后验证：
   - gitcode release view <version> -R <owner/repo> --json
   - gitcode release download <version> <asset-name> -R <owner/repo> -o ./release-verify/
9. 发布前做敏感信息检查，确认资产和 release notes 不包含 token、密码、私钥或内部敏感信息。

输入：
- previous_tag: <previous_tag>
- version: <version>
- target_branch: <target-branch>
- asset_files: <asset list>
- 本版本重点变化：<粘贴变更摘要>

输出：
- RELEASE_NOTES.md 内容
- release 创建结果
- asset 上传结果
- 发布后验证结果
- 未覆盖风险
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
