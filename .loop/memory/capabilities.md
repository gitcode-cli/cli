# gitcode-cli AI 开发规范 — 关键能力总结

三层递进架构：**规范层**（`spec/`）定义门禁和流程作为唯一规则源，**运营层**（`.loop/`）通过 prompt 模板和交付记录将规则转为可执行循环，**适配层**（`.claude/`）以 hook 和 skills 将循环接入 Claude 运行时。

| 能力 | 实现方式 | 关键文件 |
|------|---------|---------|
| **状态机驱动流程** | Issue/PR 6状态标签强制流转，每步需证据 | `spec/workflows/development-workflow.md` §2 |
| **8 项门禁执行** | 门禁表（构建/UT/Pre-commit/实际命令/CI/风险/评审/合并），docs-only 降级规则 | `spec/workflows/development-workflow.md` §5.3 |
| **全流程自主交付** | `/loop` Cron 30min 从 triage 推进到 merged | `.loop/prompts/full-flow.md` |
| **阶段冲刺验证** | `/goal` 独立评估器逐项判定，local-only 条件最可靠 | `docs/LOOP-GOAL-GUIDE.md` |
| **多角色独立评审** | 4+4 Agent 分两轮（代码/安全/测试/文档 + 架构/API/边界/UX） | `spec/workflows/review-workflow.md` |
| **Issue+PR 双评论区证据** | 验证+自检落 Issue，评审+CI+gate 表落 PR | `.loop/prompts/full-flow.md` |
| **双平台 CI 协同** | GitCode 托管 + GitHub Actions CI + `gh` CLI 监控 | `spec/delivery/ci-workflows.md` |
| **GitHub 镜像自动同步** | `git push origin` → hook 自动 `git push github` | `.claude/settings.json` |
| **风险分级闸门** | `classify-change-risk.py`，low 自动合并 / high 暂停确认 | `spec/workflows/ai-local-development-workflow.md` §9 |
| **交付记录可审计** | 按 issue 的 8 gate 证据表 + 状态流转 + PR/CI 链接 | `.loop/deliveries/issue-<N>.md` |
| **经验跨会话积累** | lessons.md 团队共享 + INDEX.md 个人上下文 + session 摘要 | `.loop/memory/` |
| **Prompt 模板复用** | 11 个 `/goal` `/loop` 场景模板，含评估器检查点和 `.loop/` 维护指令 | `.loop/prompts/` |
| **Triage 自动化** | 批量 Issue 分类打标签（type/scope/status/risk） | `.loop/prompts/batch-triage.md` |
| **PR 巡逻审查** | 定时检查他人 PR 门禁证据，缺失反馈，齐全补 CI + 合并 | `.loop/prompts/pr-review-patrol.md` |
| **统计自动生成** | `count-deliveries.sh` 解析汇总表输出实时统计 | `scripts/count-deliveries.sh` |
| **Agent 友好 CLI 契约** | `--json` 输出 + `--dry-run` + `--yes` + 非交互拒绝 | `spec/foundations/agent-friendly-cli.md` |
| **安全加固** | symlink 拒绝 + 浏览器注入修复 + repo sync 符号链接防护 | `pkg/browser/` `pkg/cmd/release/` `pkg/cmd/repo/sync/` |

---

**最后更新**: 2026-06-25
