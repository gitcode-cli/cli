# Delivery Record: Issue #400
- Title: security: auth.json file permissions rely on umask, missing explicit Chmod
- Type: security
- Risk: risk/medium
- Scope: scope/credential-storage
- Milestone: M1: Security Hardening
- Status: in-progress

## Design Artifacts
- 需求分析: .loop/deliveries/issue-400-analysis.md
- 方案设计: .loop/deliveries/issue-400-design.md
- 开发计划: .loop/deliveries/issue-400-plan.md

## State Transitions
| From | To | When | Evidence |
|------|----|------|----------|
| status/triage | status/verified | 2026-06-29 | 安全审查 Agent 2 确认 |
| status/verified | status/in-progress | 2026-07-06 | worktree bugfix/issue-400 开发，secureWriteFile 实现 + UT 完成 |

## Key Artifacts
- 分支: bugfix/issue-400
- 改动: 3 files（config.go, auth_config.go, auth_config_test.go）
- 同步主干: 2026-07-06 fast-forward 14 提交至 origin/main e65eb38，零冲突

## Gates Summary
| # | Gate | Result |
|---|------|--------|
| 0a-0c | 设计文档 | ✅ 本批补齐 |
| 1 | 验证 | ✅ |
| 2 | 开发 | ✅ 3 files |
| 3 | 构建 | ✅ go build -o ./gc ./cmd/gc，go build ./... 全包通过 |
| 4 | UT | ✅ 12 passed（含 3 新增 secureWriteFile 测试），-race 通过 |
| 5 | Lint | ✅ golangci-lint v2.12.2（CI 配置 --disable=errcheck --disable=staticcheck）0 issues；gofmt/vet 全绿 |
| 6 | 实际命令 | ⏩ 跳过（用户指示）— 待人工在 PR 评审时补：symlink 拒绝 + 权限硬化，清单见 issue-400-design.md |
| 7 | CI | ⏳ 待 PR 提交 |
| 8 | 风险分级 | ✅ classify-change-risk → **risk=high**（代码触及 auth/token/config 高风险关键词，属安全改动；按规范 high 风险需人工最终确认） |
| + | 合并 | ⏳ 待 PR（high 风险，独立 AI 评审后仍需人工最终确认） |

## 真实命令验证清单
见 .loop/deliveries/issue-400-design.md。AI 不参与 token 输入，由人工在 TTY 完成。
步骤 6 经用户指示跳过，待 PR 评审时由人工补做。

ISSUE_NUM=400
