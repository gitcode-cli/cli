# Loop 状态机

## 标准状态

```text
discovered
-> triaged
-> verified
-> planned
-> executing
-> self_checked
-> review_requested
-> ci_waiting
-> ci_passed / ci_failed
-> approved
-> merged
-> archived
```

## 状态定义

| 状态 | 含义 | 最小证据 |
| --- | --- | --- |
| `discovered` | 发现任务或候选问题 | issue、扫描记录或用户请求 |
| `triaged` | 完成基本分类和范围判断 | 类型、范围、风险初判 |
| `verified` | 问题或需求已确认有效 | 验证记录或需求确认 |
| `planned` | 已形成执行计划 | 修改范围和验证计划 |
| `executing` | 正在实施 | 非 main 分支和关联 issue |
| `self_checked` | 作者完成自检 | 测试、命令验证、文档同步、安全检查 |
| `review_requested` | 请求独立评审 | PR、自检记录、风险说明 |
| `ci_waiting` | 等待镜像 CI | PR head SHA 和 mirror repo |
| `ci_passed` | CI 通过 | GitHub Actions run URL 和 commit SHA |
| `ci_failed` | CI 失败 | 失败 job、日志 URL、下一步 |
| `approved` | 独立评审通过 | review 结论 |
| `merged` | 合入 GitCode 主干 | merged PR 和 `origin/main` 验证 |
| `archived` | 完成知识归档判定 | archive decision |

## 状态记录

- 长期状态写回 GitCode issue / PR
- 本地 `.loop/runtime` 只作临时缓存
- 状态推进事件应符合 `loop-kits/schemas/loop-event.schema.json`

## 禁止跳转

- 未 `verified` 不得进入代码实施
- 未 `self_checked` 不得请求独立评审
- 未绑定 commit SHA 不得判定 CI 通过
- 未 merged PR + `origin/main` 验证不得判定完成
