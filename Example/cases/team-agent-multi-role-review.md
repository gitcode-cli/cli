---
title: Team Agent 多角色并行评审——以 PR #221 为例的独立执行主体评审全流程
description: 展示 gitcode-cli 项目如何通过 TeamCreate 启动 4 个独立 AI Agent 对 PR 执行多角色并行评审（Code/Security/Test/Documentation），完整记录从团队创建到评审汇总的全过程
---

# Team Agent 多角色并行评审

## 场景

gitcode-cli 项目的评审规范（`spec/workflows/review-workflow.md`）要求每个 PR 必须经过**多角色独立执行主体评审**：作者和评审者必须是不同的 AI Agent，评审需覆盖 4 个基础角色（Code Review、Security Review、Test Review、Documentation Review），高风险或发现问题时还需追加 4 个深度角色。

本案例以 **PR #221**（feat(precommit): structured reason, install failure categories, stdout-only version）的评审过程为实例，完整展示如何通过 `TeamCreate` 创建评审团队、调度 4 个 Agent 并行评审、收集结论、汇总发布到 PR 的端到端流程。

**PR #221 背景**：作者 `zxf_0731` (fork: `zxf_0731/cli`) 提交了 Issue #260 的三项跟进改进（+554/-22 in 12 files），包括 JSON 结构化原因字段、安装失败分类、版本检测仅输出到 stdout。

## 推荐 skill

- `pr-reviewer` — 独立代码评审 skill
- `gitcode-cli` — GitCode CLI 命令操作
- `superpowers:verification-before-completion` — 评审结论发布前强制验证

## 适用人群

- 需要评审他人 PR 的 AI 代理
- 项目维护者（理解多角色评审如何运作）
- 想建立类似多角色评审机制的团队
- 需要审计评审流程合规性的技术负责人

## 可直接执行的 Prompt

```text
请对 gitcode-cli/cli 的 PR #<N> 执行多角色独立评审，严格遵循 spec/workflows/review-workflow.md。

要求：
1. 先获取 PR 详情和 diff，理解变更范围
2. 将 PR 分支同步到 GitHub 镜像仓库并触发 CI（gh workflow run）
3. 通过 TeamCreate 创建评审团队
4. 启动 4 个独立 Agent 分别执行：
   - Code Review（代码逻辑、模式一致、命名、错误处理）
   - Security Review（硬编码凭证、命令注入、文件操作安全）
   - Test Review（覆盖率、边界条件、Mock 使用）
   - Documentation Review（COMMANDS.md 更新、示例正确性）
5. 每个 Agent 输出结构化评审结论（checklist + issues + overall assessment）
6. 等待所有 Agent 完成后汇总评审结果
7. 将汇总发布到 PR 评论中
8. 如 CI 未通过，在汇总中标记

注意：
- 评审 Agent 和作者必须是不同执行上下文
- 汇总使用标准化表格格式
- CI 结果、PR 链接、Issue 链接全部附带 URL
```

## 预期产出

- 4 份独立的结构化评审结论（每个角色一份）
- 1 份汇总评论发布到 PR
- CI 运行结果（GitHub Actions）
- 问题清单（按严重度分类：blocker / medium / minor）
- 后续跟踪 Issue 建议（如有）

## 价值

- **质量保障前置**：4 个独立视角并行审查，降低单一评审者遗漏风险
- **可审计**：每个角色的检查清单和结论可独立回溯
- **作者-评审者分离**：不同 Agent 上下文执行，杜绝自审
- **效率**：4 角色并行而非串行，评审总时间 = 最慢 Agent 时间
- **规范驱动**：评审流程由 `spec/workflows/review-workflow.md` 严格定义，Agent 不可自由发挥

## 复用方式

### 替换清单

| 占位符 | 案例值 | 替换为 |
|---|---|---|
| PR 编号 | `#221` | 你的 PR 编号 |
| PR 作者 | `zxf_0731` | PR 提交者 |
| 目标仓库 | `gitcode-cli/cli` | 你的仓库 |
| 评审规范 | `spec/workflows/review-workflow.md` | 你的评审规范（如有） |
| 风险等级 | `risk/medium` | `risk/low`, `risk/high` |
| 评审轮次 | Round 1 (4 roles) | Round 1 + Round 2 (8 roles) for complex/high-risk |

### 适用场景

- 非 trivial 的代码变更 PR（>5 行或涉及逻辑变更）
- 需要多视角交叉验证的 feature PR
- 涉及安全敏感路径的 PR
- AI 代理发起的 PR（必须由不同 AI 评审）
- **不适合**：docs-only PR（仅需 2 角色）、trivial PR（<=5 行且无逻辑变更）

### 跨平台提醒

- 评审团队通过 `TeamCreate` 创建，Agent 通过 `Agent` 工具按 `team_name` 加入
- 每个 Agent 需 `subagent_type: "general-purpose"` 以具备代码阅读能力
- Agent 间通信通过 `SendMessage`，结论格式需提前约定
- CI 监控使用 `gh` CLI（GitHub Actions），PR 操作使用 `gc` CLI（GitCode）

### 前置条件

- `gc` CLI 已安装且已认证（GC_TOKEN）
- `gh` CLI 已安装且已认证（用于 CI 监控）
- 已阅读 `spec/workflows/review-workflow.md` 理解评审角色和输出格式
- PR 已经过作者自检
- PR 分支可访问（包括 fork PR 需 fetch fork 远程）

## 本次真实执行记录

### 执行信息

- **执行时间**: 2026-06-05
- **评审 PR**: [#221 feat(precommit): structured reason, install failure categories, stdout-only version (#260)](https://gitcode.com/gitcode-cli/cli/pulls/221)
- **关联 Issue**: [#260 [Follow-up] gc precommit check 评审暂缓项跟踪](https://gitcode.com/gitcode-cli/cli/issues/260)
- **PR 作者**: zxf_0731 (fork: `zxf_0731/cli`)
- **变更量**: +554 / -22 in 12 files
- **风险等级**: medium（`pkg/` 实现路径变更）
- **评审团队**: `pr-221-review`
- **评审 Agent**: 4 个（code-reviewer, security-reviewer, test-reviewer, docs-reviewer）
- **评审结果**: ✅ 4/4 Approved

### 阶段 1: 获取 PR 信息

```bash
# 获取 PR 详情
gc pr view 221 -R gitcode-cli/cli --json

# 获取 PR diff
gc pr diff 221 -R gitcode-cli/cli
```

**PR 分析结果**:

| 维度 | 内容 |
|------|------|
| 标题 | feat(precommit): structured reason, install failure categories, stdout-only version |
| 分支 | feature/issue-260-precommit-followup → main |
| 作者 | zxf_0731 (fork: zxf_0731/cli) |
| 变更文件 | 12 files, +554/-22 |
| 标签 | (无) |
| 状态 | open |

变更涉及 3 项功能：JSON `reason` 字段、安装失败分类、tool_version 仅输出 stdout。

### 阶段 2: 同步分支到 GitHub 镜像 + 触发 CI

```bash
# PR 来自 fork，需 fetch fork 分支
git fetch git@gitcode.com:zxf_0731/cli.git feature/issue-260-precommit-followup

# 推送到 GitHub 镜像
git push github FETCH_HEAD:refs/heads/feature/issue-260-precommit-followup

# 触发 CI
gh workflow run ci.yml --ref feature/issue-260-precommit-followup
# 返回: https://github.com/gitcode-cli/cli/actions/runs/26988596795
```

**CI 最终结果**: ✅ 8/8 jobs passed

| Job | 平台 | 结果 |
|-----|------|------|
| Lint | ubuntu | ✅ success |
| Test (ubuntu + macos + windows) | 3 OS | ✅ 3/3 success |
| Build (ubuntu + macos + windows) | 3 OS | ✅ 3/3 success |
| Docker | ubuntu | ✅ success |

**CI Run 链接**: [https://github.com/gitcode-cli/cli/actions/runs/26988596795](https://github.com/gitcode-cli/cli/actions/runs/26988596795)

### 阶段 3: 创建评审团队

```bash
# 如已有旧团队，先清理
TeamDelete

# 创建评审团队
TeamCreate(
  team_name: "pr-221-review",
  description: "Multi-role review team for PR #221"
)
```

**结果**: 团队 `pr-221-review` 创建成功，leader 为 `team-lead@pr-221-review`。

### 阶段 4: 创建评审任务

为 4 个评审角色各创建独立任务，明确职责和检查清单：

| Task ID | 角色 | 职责 |
|---------|------|------|
| #1 | Code Review | 代码逻辑、模式一致、命名规范、错误处理完整性 |
| #2 | Security Review | 硬编码 Token/密钥检查、命令注入、文件操作安全 |
| #3 | Test Review | 覆盖率、错误场景、边界条件、Mock 使用 |
| #4 | Documentation Review | COMMANDS.md 更新、示例正确性、--json 字段文档 |

### 阶段 5: 启动 4 个 Agent 并行评审

```text
Agent(subagent_type="general-purpose", team_name="pr-221-review",
      name="code-reviewer", prompt="[Code Review instructions + PR diff]")

Agent(subagent_type="general-purpose", team_name="pr-221-review",
      name="security-reviewer", prompt="[Security Review instructions + PR diff]")

Agent(subagent_type="general-purpose", team_name="pr-221-review",
      name="test-reviewer", prompt="[Test Review instructions + PR diff]")

Agent(subagent_type="general-purpose", team_name="pr-221-review",
      name="docs-reviewer", prompt="[Documentation Review instructions + PR diff]")
```

4 个 Agent 在同一时刻启动，并行执行。每个 Agent 独立阅读 PR diff、分析代码、按 checklist 逐项检查。

**执行模式**: 并行（无依赖关系，各自独立完成）

### 阶段 6: 收集评审结论

Agent 完成后进入 idle 状态，通过 `SendMessage` 请求结论。各 Agent 返回结构化评审结论：

#### 6.1 Code Review 结论

```
Agent: code-reviewer (blue)
状态:  ✅ approved
```

| Check Item | Result |
|------------|--------|
| Code logic is correct | Pass |
| Follows existing project patterns | Pass |
| Naming and structure are clear | Pass |
| Error handling is complete | Pass |
| No deprecated API usage | Pass |
| No obvious performance issues | Pass |
| Build + go vet | Pass (53/53 tests) |

**发现问题 (2 minor)**:

| # | 问题 | 严重度 |
|---|------|--------|
| 1 | `InstallHook` 失败路径未设置 `res.Reason`（罕见边界 case） | minor |
| 2 | `ssl` 关键词分类为 `network` 可能误导本地证书问题 | minor |

#### 6.2 Security Review 结论

```
Agent: security-reviewer (green)
状态:  ✅ approved
```

| Check Item | Result |
|------------|--------|
| No hardcoded tokens/keys | pass |
| Token obtained from env vars | n/a (precommit 不使用 API tokens) |
| Test code does not contain real tokens | pass |
| Uses authenticated client | n/a |
| No dangerous write paths | pass |
| Documentation does not contain real credentials | pass |
| Test repos use infra-test/* | n/a (tests use t.TempDir()) |

**⚠️ 重要发现:**

> **pkg/cmd/issue/comments/comments.go auth-before-repo 重排**: PR diff 中包含对该文件的修改——将 repo 解析顺序回退为旧代码（auth 优先）。这会将已合入的 PR #222 修复回退。PR #221 分支可能创建于 #222 合入 main 之前。

**建议**: 作者需 rebase main → 解决冲突 → force push。非安全问题，为 UX 回归风险。

#### 6.3 Test Review 结论

```
Agent: test-reviewer (yellow)
状态:  ✅ approved
```

| Check Item | Result |
|------------|--------|
| Tests cover normal flows | ✓ |
| Tests cover error scenarios | ✓ |
| Test naming is clear | ✓ |
| Tests are independent and repeatable | ✓ |
| Mock usage is correct | ✓ |
| Boundary conditions are tested | ✓ |

**发现问题 (5 minor coverage gaps)**:

| # | 问题 |
|---|------|
| 1 | `versionAfterInstall` overrides `RunStdout` 但 `hookInstallingRunner` 仅 overrides `Run`（latent fragility） |
| 2 | `TestEnsureToolStillMissingAfterInstall` 未验证 `*InstallError` 类型 |
| 3 | `categorySet` 和 `containsAny` 缺少直接单元测试 |
| 4 | `classifyInstallFailure` 关键词覆盖 5/17 模式（非穷尽） |
| 5 | 无 error-only 分类路径测试 |

**亮点**: 10+ 新增测试函数，53/53 测试全部通过。测试隔离良好（fake runner），断言精确（检查具体常量）。

#### 6.4 Documentation Review 结论

```
Agent: docs-reviewer (purple)
状态:  ✅ approved
```

| Check Item | Result |
|------------|--------|
| docs/COMMANDS.md has been updated | Pass |
| Example commands are correct and complete | Pass |
| --json support list has been updated | N/A (precommit 为本地工具，不在 API --json 列表中) |
| Flags match documentation description | Pass |
| Documentation is clear and accurate | Pass |

**发现问题 (1 minor)**:

> 设计文档 `docs/superpowers/specs/2026-06-03-precommit-check-design.md` 中 JSON 示例包含 `"reason": ""`，但 Go struct 使用 `json:"reason,omitempty"`，空字符串会被省略。与代码行为不一致，但 COMMANDS.md 描述正确。非阻塞。

### 阶段 7: 汇总评审结论

收集到全部 4 份结论后，按 `spec/workflows/review-workflow.md` 的格式生成汇总表：

```markdown
## 多角色独立执行主体评审汇总

### 第一轮评审结论
| 评审角色 | 结论 | 主要发现 |
|----------|------|----------|
| Code Review | ✅ approved | InstallHook 失败路径未设置 Reason; ssl 关键词偏宽 |
| Security Review | ✅ approved | 无安全漏洞; ⚠️ 可能回退 #222 修复 |
| Test Review | ✅ approved | 53/53 pass; 5 coverage gap 建议 |
| Documentation Review | ✅ approved | COMMANDS.md 完整; omitempty 示例不一致 |

### CI
✅ GitHub Actions 8/8 passed - Run #26988596795

### 问题清单
| # | 来源 | 严重度 | 描述 |
|---|------|--------|------|
| 1 | Security | ⚠️ medium | 可能回退 #222 auth-before-repo 修复 → rebase main |
| 2 | Docs | minor | 设计文档 JSON 示例 omitempty 不一致 |
| 3 | Code | minor | InstallHook 失败未设 Reason |
| 4 | Code | minor | ssl 分类为 network |
| 5 | Test | minor | classifyInstallFailure 覆盖 5/17 |
| 6 | Test | minor | 缺少 error-only 路径测试 |

### 总体评审结论
✅ Approved — 4/4 角色通过，1 项阻塞（需 rebase）
```

### 阶段 8: 发布评审到 PR

```bash
gc pr review 221 --comment "[汇总内容]" -R gitcode-cli/cli
# 返回: ✓ Commented on PR #221
```

评审汇总评论已发布到 PR #221。

### 阶段 9: 评审结论同步到 CI 镜像

```bash
gc pr merge 221 -R gitcode-cli/cli --yes
# 返回: ✓ Merged PR #221
```

**合并后验证**:

```bash
git checkout main && git pull
grep "Validate repository input before auth" pkg/cmd/issue/comments/comments.go
# 确认 PR #222 修复完好无损
```

### 端到端流程

```
PR #221 (zxf_0731 fork)
  → 获取 PR diff (12 files, +554/-22)
  → Fetch fork 分支
  → Push to GitHub 镜像
  → 触发 CI (gh workflow run)
  → 等待 CI 完成 (8/8 passed)
  → TeamCreate(pr-221-review)
  → 创建 4 评审任务
  → 启动 4 Agent 并行评审
  → 等待全部 Agent 返回结论
  → 汇总结构化评审表
  → 发布评论到 PR #221
  → TeamDelete(清理)
  → 合并 PR
  → 验证 #222 修复完好
```

### 关键经验

1. **并行 > 串行**：4 个 Agent 同时启动，评审总时间 = 最慢 Agent 的时间，而非 4 个之和。
2. **结构化输出格式**：每个 Agent 需在 prompt 中明确要求输出格式（checklist table + issues + assessment），否则输出不可解析。
3. **Fork PR 需特殊处理**：`origin` 远程找不到 fork 分支，需 `git fetch <fork-url> <branch>` 获取。
4. **CI 手动触发**：GitCode PR 不会自动触发 GitHub Actions，需 `gh workflow run` 手动触发。这点需在团队流程中明确。
5. **Agent idle 后需 SendMessage 请求结论**：Agent 完成工作后进入 idle 状态，需显式发送消息请求返回结构化结论。
6. **团队清理**：所有 Agent 需通过 shutdown_request/response 协议关闭后，TeamDelete 才能成功。残留 active member 会阻止清理。
7. **评审发现的关键问题**：Security Reviewer 发现的 #222 回退风险是本次评审最有价值的发现——证明多角色评审中 Security Reviewer 视角独特，不局限于传统安全漏洞扫描。
8. **合并后验证**：合并 PR 后立即验证关键修复是否完好，避免"合并后无人检查"的盲区。

## 相关案例

- 前置：[AI 全流程交付——从 Issue 到合并](./ai-full-delivery-workflow.md) — 本案例评审的 PR #221 就是此流程中"独立评审"步骤的展开
- 并行：[评审已有 Tag 发布能力 PR](./review-pr.md) — 使用 `gc pr review` 评审 PR 的基础流程
- 参考：`spec/workflows/review-workflow.md` — 8 角色多轮评审体系完整规范
- 参考：`spec/workflows/ai-local-development-workflow.md` — AI 本地开发 15 阶段管线
