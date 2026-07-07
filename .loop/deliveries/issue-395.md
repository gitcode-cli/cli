# Delivery Record: Issue #395
- Title: security: git.RemoteURL lacks -- separator, risk of git option injection
- Type: security
- Risk: risk/medium
- Scope: scope/git
- Milestone: M1: Security Hardening
- Status: self-checked

## Design Artifacts
- 需求分析: .loop/deliveries/issue-395-analysis.md
- 方案设计: .loop/deliveries/issue-395-design.md
- 开发计划: .loop/deliveries/issue-395-plan.md

## State Transitions
| From | To | When | Evidence |
|------|----|------|----------|
| status/triage | status/verified | 2026-07-07 | RemoteURL 确认缺 -- + ValidateRef |
| status/verified | status/in-progress | 2026-07-07 | bugfix/issue-395 开发 |
| status/in-progress | status/self-checked | 2026-07-07 | 作者自检完成，本地验证全绿 |

## Key Artifacts
- 分支: bugfix/issue-395
- 改动: 2 files（git.go, git_test.go）

## Gates Summary
| # | Gate | Result |
|---|------|--------|
| 1 | 验证 | ✅ |
| 2 | 开发 | ✅ 2 files |
| 3 | 构建 | ✅ go build ./... |
| 4 | UT | ✅ TestRemoteURLRejectsOptionInjection 全 PASS，-race 通过 |
| 5 | Lint | ✅ 0 issues |
| 6 | 实际命令 | ⏩ 豁免（内部 git 封装，UT 充分覆盖） |
| 7 | CI | ⏳ 待 PR 提交 |
| 8 | 风险分级 | ✅ risk=medium |
| + | 合并 | ⏳ 待 PR |

ISSUE_NUM=395
