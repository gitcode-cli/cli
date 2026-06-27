在 git worktree 中从 status/triage 取一个 issue，按 spec/workflows/development-workflow.md 状态机推进到 merged。本次只处理一个，完成后停止。

状态机全流程推进到 merged。risk/low 自动合并，risk/high 暂停确认。

前置规则：所有代码操作必须在独立 git worktree 中执行，worktree 名称必须包含 issue 号和当前时间戳确保唯一（如 .claude/worktrees/issue-N-timestamp），用后即删。严禁污染主工作目录。

禁止：在 main 开发、跳过验证、作者自检当独立评审。

8 项门禁逐一执行，docs-only 以外的门禁不得跳过：

| # | 门禁 | 要求 |
|---|------|------|
| 1 | 开发实现 | 修复或实现 |
| 2 | 测试 | docs-only 跳过；代码改动必须 go test ./... 全部通过 |
| 3 | 本地构建 | docs-only 跳过；代码改动必须 go build -o ./gc ./cmd/gc 成功 |
| 4 | 单元测试 | docs-only 跳过；代码改动全部通过 |
| 5 | Pre-commit | 全部 hooks 通过，不得跳过 |
| 6 | 实际命令验证 | docs-only 跳过；代码改动必须在 infra-test/* 仓库执行至少一条 |
| 7 | 远端 CI | docs-only 跳过；代码改动必须用 gh CLI 触发 GitHub Actions CI 并等待全部 Job 通过，PR 评论中附 run URL |
| 8 | 风险分级 | classify-change-risk.py，risk/low 自动合并，risk/high 暂停 |

CI 执行步骤（代码改动必须）：
  unset HTTP_PROXY HTTPS_PROXY http_proxy https_proxy
  git push github <branch>                    # 推到 GitHub 镜像仓
  gh run list -R gitcode-cli/cli -b <branch> --json url,databaseId,conclusion
  等待全部 Job 完成
  如果失败→分析根因→修复→推送→重新触发，直到全绿

提交前检查 spec/workflows/development-workflow.md §5.3 确认无遗漏。

必须在 Issue 评论和 PR 评论留下门禁证据：
  Issue comment: 验证记录 + 作者自检 9 项（含 CI run URL）
  PR comment: 多角色评审结论 + CI run URL + 8 gate 完成表

每处理完一个 issue，更新 .loop/ 目录：
  .loop/deliveries/issue-N.md — 完整状态流转 + 8 gate 表 + PR/CI 证据链接
  .loop/deliveries/README.md — 更新对应行（CI 列：有 run URL 才填 ✅）

注意：从状态机第一个步骤开始，逐步推进。每步留下证据。不要跳过任何门禁。CI 未跑就写 ✅ 算违规。
