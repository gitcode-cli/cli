# 经验教训（团队共享）

跨会话积累的操作经验，供所有 AI 协作者和人工维护者参考。

## Loop 运行

1. **后台 loop 用 Cron 模式** — ScheduleWakeup 在用户活跃时不触发，Cron 固定时钟可靠
2. **一条全流程 loop 优于多条接力** — triage 线和 verified 线没有独立存在价值
3. **spec 门禁表是 /loop 的基石** — spec 全则 loop 全，spec 漏则 loop 漏

## Agent 可靠性

4. **多角色评审 Agent ~50% 静默失败率** — 2 分钟无响应则手动补齐评审
5. **风险评估脚本误判** — `classify-change-risk.py` 扫描全部累积 diff，小改动也被标 risk/high
6. **建议增加 `--head-only` 模式** — 仅扫描当前分支 diff

## 环境

7. **GitHub 代理阻断** — `HTTP_PROXY` 环境变量导致 `gh` 返回 EOF，需 `unset` 后使用
8. **GitCode PR 405** — `gc pr merge --yes` 可能返回 405，fallback: `git merge` + `git push`
9. **CI 触发方式** — 项目 CI 由 `pull_request` 事件触发，push 到 GitHub 镜像仓后用 `gh workflow run` 手动 dispatch

## 门禁

10. **docs-only 不能全跳** — 仍需 pre-commit + 安全审查 + 2 角色评审（docs+security），见 spec 5.3
11. **实际命令测试必须用 infra-test/\*** — 不得用生产仓库测试

## git 操作

12. **`git commit --amend` 前置条件** — pre-commit hook 失败后不能直接 amend，需先 `git add` 修复文件再 commit

---

**最后更新**: 2026-06-24
