# Evidence 模型

## 目标

Evidence 用于解释为什么某个 loop 可以推进到下一状态。

## 证据类型

| 类型 | 示例 | 长期位置 |
| --- | --- | --- |
| verification | 复现命令、需求确认 | GitCode issue comment |
| plan | 修改范围、验证计划 | GitCode issue comment 或 PR body |
| test | `go test ./...`、真实命令验证 | GitCode issue / PR comment |
| ci | GitHub Actions run URL | GitCode PR comment |
| review | 独立评审结论 | GitCode PR review / comment |
| merge | merged PR、`origin/main` 包含关系 | GitCode issue / PR comment |
| archive | 归档判定 | GitCode issue / PR comment |

## 最小字段

- loop id
- phase
- actor
- timestamp
- repo
- issue 或 PR
- commit SHA（涉及代码或 CI 时必须）
- evidence URL 或命令结果摘要
- 下一步动作

## 存储边界

- GitCode issue / PR 保存长期证据
- GitHub Actions 保存 CI 原始事实
- `.loop/runtime` 只保存本地临时结构化缓存
- `loop-kits/templates` 只保存模板，不保存填好的证据
