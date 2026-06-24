# /loop 全流程交付

## Prompt

```
/loop 30m 从 status/triage 中取一个 issue，按 spec/workflows/development-workflow.md
状态机全流程推进到 merged。risk/low 自动合并，risk/high 暂停确认。
禁止：在 main 开发、跳过验证、作者自检当独立评审。
每处理完一个 issue，更新 .loop/ 目录：
  .loop/deliveries/issue-<N>.md — 记录状态流转 + PR/CI 证据链接
  .loop/memory/INDEX.md — 更新活跃 Issue、待办、经验教训
会话结束时写 .loop/memory/YYYY-MM-DD-session.md
```

## 设置方式

1. 替换间隔（推荐 30m 或 10m）
2. 确认 spec 5.3 门禁表已是最新版本
3. 启动后写入 `.loop/registry/active.yaml`

## 门禁清单

AI 必须逐项执行 spec 5.3 的 8 项门禁表，每项完成后留下证据：

| # | 门禁 | 证据 |
|---|------|------|
| 1 | 验证 | issue comment 中的复现记录 |
| 2 | 开发 | 非 main 分支 + commits |
| 3 | 构建 | `go build -o ./gc ./cmd/gc` 成功 |
| 4 | UT | `go test ./...` 全部通过 |
| 5 | Pre-commit | 所有 hooks 通过 |
| 6 | 实际命令 | `./gc <cmd> -R infra-test/gctest1` |
| 7 | CI | `gh run list` 全绿 + run URL |
| 8 | 风险分级 | `scripts/classify-change-risk.py` |
| + | 多角色评审 | 4 Agent 结论汇总到 PR |
| + | 合并 | risk/low 自动，risk/high 暂停 |

## .loop/ 维护要求

每个 issue 处理完成后，必须更新：

### deliveries/issue-<N>.md
```markdown
# Delivery Record: Issue #<N>
- Title: <title>
- Type: <bug/feature/docs>
- Status: <merged/closed>

## State Transitions
| From | To | When | Evidence |

## Key Artifacts
- PR: !<N>, CI Run: <URL>, Commits: <sha>
```

### memory/INDEX.md
- 更新"当前活跃 Issue"列表
- 如有新发现，追加到"经验教训"
- 更新"待办"项

### 会话结束
- 写 `memory/YYYY-MM-DD-session.md` 记录本会话摘要
- 标记 `registry/active.yaml` 中的 loop 为 lost
