# Issue #304: docs — 补齐 pr label 和 repo clone 命令文档

## 状态流转

| 时间 | 状态 | 说明 |
|------|------|------|
| 2026-06-26 | triage → verified | loop 接手 |
| 2026-06-26 | verified → in-progress | 开始实现 |
| 2026-06-26 | in-progress → merged | PR #288 merged |

## 交付

- **类型**: docs
- **PR**: [#288](https://gitcode.com/gitcode-cli/cli/merge_requests/288)
- **分支**: fix/issue-304-docs-pr-label-repo-clone
- **风险**: risk/low
- **来源**: daily-audit loop

## 门禁

| # | 门禁 | 状态 |
|---|------|------|
| 1 | 开发实现 | ✅ 3 files, +71/-1 |
| 2 | 测试 | ✅ docs-only 跳过 |
| 3 | 构建 | ✅ docs-only 跳过 |
| 4 | UT | ✅ docs-only 跳过 |
| 5 | Pre-commit | ✅ 10/10 |
| 6 | 实际命令 | ✅ docs-only 跳过 |
| 7 | CI | ✅ docs-only 跳过 |
| 8 | 风险分级 | risk/low |

## 修改文件

- `docs/COMMANDS.md` — +repo clone 小节，+pr label 小节，+json 列表
- `.ai/distribution/gc-core/pr/SKILL.md` — +pr label 场景/命令/约束
- `.ai/distribution/gc-core/repo/SKILL.md` — +clone 场景/命令/约束

## 评论证据

- Issue comment: [验证记录+自检9项](https://gitcode.com/gitcode-cli/cli/issues/304#177337215)
- PR comment: [门禁检查表](https://gitcode.com/gitcode-cli/cli/merge_requests/288#ca7c15eb807ba3a2e335badaeddc13913c10cece)
