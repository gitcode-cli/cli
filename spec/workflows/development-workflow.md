# 开发工作流程

本文档定义从 Issue 到主干合入的完整开发流程。

本文件的目标不是提供一份“建议清单”，而是为人工和 AI 定义一套不可随意跳步的状态化交付流程。

## 职责

- 定义从 issue 到 merge 的总流程
- 定义每个阶段的准入条件、允许动作和退出条件
- 约束 AI 只能沿受控流程推进，而不是自由发挥

## 适用场景

- 修复 issue
- 开发新功能
- 提交 PR
- 做本地自检、审查和合并前确认

## 必须

- 修复前先验证问题是否仍存在
- 不在 `main` 直接开发
- 补测试、做本地验证和实际命令验证
- 在提 PR 前完成安全审查
- 按状态推进 issue 和 PR，不得跳步
- 每次状态推进都留下结构化证据

## 禁止

- 看到 issue 就直接写代码
- 跳过验证、测试或安全检查
- 在主分支直接提交
- 把作者自检当作独立执行主体评审
- 在 PR 未合入主干前关闭“待修复” issue

## 同步要求

- 改动流程时同步 `spec/*`
- 改动命令行为时同步 `docs/COMMANDS.md`
- 改动 AI 协作或入口规则时同步 `AGENTS.md`、`CLAUDE.md` 和相关 skills

## 不负责什么

- 具体编码风格
- 测试覆盖细节
- 文档分层设计
- blocker 级问题的最终审查裁定

## 1. 流程总原则

本仓库对人工和 AI 一视同仁，统一遵守以下原则：

1. 先定义状态，再执行动作
2. 先提交证据，再推进状态
3. 作者实现、自检、独立执行主体评审必须分离
4. 本地自动检查和独立执行主体语义审查必须分层
5. 所有推进动作都必须可回放、可审计

## 2. 状态机

### 2.1 Issue 状态

Issue 使用以下状态标签推进：

| 状态标签 | 含义 | 谁可推进 |
|------|------|------|
| `status/triage` | 已受理，等待分类与补充信息 | 人工或 AI |
| `status/verified` | 已完成复现或确认需求有效 | 人工或 AI |
| `status/in-progress` | 已开始开发 | 人工或 AI |
| `status/blocked` | 当前无法继续推进 | 人工或 AI |
| `status/ready-for-review` | 代码与自检已完成，等待评审 | 人工或 AI |
| `status/merged` | 关联 PR 已合入主干 | 合并后更新 |
| `status/closed-no-fix` | 判定无效、重复、已修复或无需处理 | 人工或 AI |

### 2.2 PR 状态

PR 使用以下状态标签推进：

| 状态标签 | 含义 | 谁可推进 |
|------|------|------|
| `status/draft` | 已开工但未完成自检 | 作者 |
| `status/self-checked` | 作者已完成结构化自检 | 作者 |
| `status/ready-for-review` | 可以进入独立执行主体评审 | 作者 |
| `status/changes-requested` | 发现问题，需继续修改 | 评审者 |
| `status/approved` | 评审通过，可合并 | 评审者 |
| `status/merged` | 已合入主干 | 合并后更新 |

## 3. AI 执行边界

AI 在本仓库中必须遵守以下硬约束：

- 未进入 `status/verified` 的 issue，不得开始写代码
- 未创建开发分支前，不得修改实现文件
- 未完成自检证据前，不得把 PR 标记为 `status/ready-for-review`
- 不得把作者自己的 `gc pr review --comment` 当作独立执行主体评审通过
- 不得在 PR 未合入 `main` 前把待修复 issue 标记为完成
- 不得声称“已完成”但不给出测试、命令验证、文档同步证据

## 4. 标准流程

```
Issue 创建/受理
→ issue triage
→ issue verified
→ 创建分支
→ 开发与测试
→ 安全审查
→ issue in-progress
→ PR draft
→ 作者自检
→ 风险分级
→ 第一轮多角色评审
→ 发现问题则修复
→ 第二轮多角色评审（如有）
→ approved
→ merge
→ issue / PR merged
```

## 5. 分阶段执行要求

### 5.1 Triage

目标：确认 issue 类型、范围和下一步动作。

必须完成：

- 判断是 `bug`、`type/feature`、`type/docs` 还是其他类型
- 补基础标签和范围标签
- 判断是否需要补充复现信息

允许输出：

- issue comment：补充信息请求
- issue labels：类型、范围、状态

### 5.2 Verified

目标：确认问题真实存在，或确认需求确实要做。

必须完成：

- 用当前代码验证 issue 是否仍存在
- 检查时间线，避免修重复问题
- 给出“继续修复”或“关闭不修”的明确结论

结构化验证记录至少包含：

- 当前验证版本或分支
- 复现命令
- 实际结果
- 结论

### 5.3 In Progress

进入 `status/in-progress` 前必须满足：

- issue 已进入 `status/verified`
- 已创建非 `main` 开发分支

本阶段必须完成以下全部门禁，按执行顺序：

| 序号 | 门禁 | 要求 | 详细规范 | docs-only 是否跳过 |
|------|------|------|---------|:--:|
| 1 | 开发实现 | 修复或实现 issue 描述的功能 | — | — |
| 2 | 编写/补齐测试 | 覆盖修改路径的单元测试 | [测试指南](../foundations/testing-guide.md) | ✅ 可跳过 |
| 3 | 本地构建 | `go build -o ./gc ./cmd/gc` 成功 | — | ✅ 可跳过 |
| 4 | 单元测试 | `go test ./...` 全部通过 | [测试指南](../foundations/testing-guide.md) | ✅ 可跳过 |
| 5 | Pre-commit | 通过项目 `.pre-commit-config.yaml` 定义的全部 hooks | [代码质量门禁规范](../foundations/code-quality-gates.md) | — |
| 6 | 实际命令验证 | 至少一条真实命令，仅限 `infra-test/*` 仓库 | [测试流程](./test-workflow.md) | ✅ 可跳过 |
| 7 | 回归测试 | `./scripts/regression-core.sh` 通过（核心冒烟回归） | [测试流程](./test-workflow.md) | ✅ 可跳过 |
| 8 | 系统测试 | `go test -tags=system ./tests/system` 通过（只读路径，需本地认证） | [测试指南](../foundations/testing-guide.md) | ✅ 可跳过 |
| 9 | 远端 CI | 通过 GitHub Actions CI 全部 Job（lint/test/build/docker） | [CI 工作流规范](../delivery/ci-workflows.md) | ✅ 可跳过 |
| 10 | 风险分级 | 运行 `scripts/classify-change-risk.py --base origin/main` | [AI 本地开发流程](./ai-local-development-workflow.md) | — |

docs-only 改动可跳过门禁 2-4、6-9，但必须在自检中说明跳过理由。
CI 因环境原因（如 GitHub 镜像仓不可达）无法执行时，必须在自检中明确记录原因。

> **AI 协作者提示**：`/loop` 或 `/goal` 执行本阶段时，应逐项对照此表，每项完成后在对话中留下证据，未完成不得进入 5.5 自检。此表是 development-workflow.md 的权威门禁清单，其他文档的补充说明不得与此表冲突。

### 5.4 Security Review

进入作者自检前，至少检查：

- 无硬编码 token、password、secret
- 文档、测试和示例中未误写真实凭证
- 涉及认证、配置、权限、网络调用、删除或覆盖行为时，已对照 `spec/foundations/security.md` 检查
- 如存在安全影响，已在后续自检和评审中明确记录

### 5.5 Self Check

代码开发完成后，作者必须先做自检，不能直接等同于通过评审。

作者自检至少包含：

- 根因或实现理由
- 修改范围
- 测试结果
- 实际命令验证结果
- 安全审查结果
- 文档同步结果
- 风险点
- 未覆盖项

### 5.6 多角色评审

进入 `status/ready-for-review` 前必须完成多角色评审：

**第一轮评审（必须执行）**：
- 代码审查：检查代码逻辑、项目模式遵循
- 安全审查：检查凭证泄漏、Token 处理
- 测试审查：检查测试覆盖、边界条件
- 文档审查：检查文档同步、示例正确

**第二轮评审（发现问题时执行）**：
- 架构一致性：检查 Options 结构、函数命名
- API 契约：检查 API 参数、调用时序
- 边界条件：检查输入验证、错误处理
- 用户体验：检查帮助文本、错误消息

评审要求：
- 评审 Agent 必须与作者 Agent 是不同执行主体
- 每个角色输出结构化评审结论
- 发现 P0 问题必须修复后重新评审
- 非阻塞问题可创建 Issue 跟踪

### 5.7 Ready For Review

只有在以下条件全部满足后，才允许把 issue / PR 推到 `status/ready-for-review`：

- 自检证据完整
- 本地门禁已完成
- 真实命令验证已完成，或明确写出未执行原因
- 安全审查已完成，或明确写出本次改动不涉及安全敏感路径
- 文档同步已完成，或明确写出无需更新的依据
- 风险分级已完成，并已明确本次改动的评审策略
- **第一轮多角色评审已完成且全部通过**
- **第二轮多角色评审已完成**（如有问题）
- 评审汇总评论已添加到 PR

### 5.8 Merge

只有在以下条件全部满足后，才允许合并：

- 无 blocker 级问题
- 第一轮多角色评审全部通过
- 第二轮多角色评审全部通过（如有）
- PR 状态为 `status/approved`
- 相关 issue 和 PR 的状态、标签、评论记录完整
- 若为 `risk/high` 改动，已完成人工最终确认

合并后才能：

- 把 PR 标记为 `status/merged`
- 把对应 issue 标记为 `status/merged` 或正常关闭

## 6. 修复 Issue 前的验证要求

### 错误示例

```bash
# ❌ 看到 Issue 就立即创建分支、写代码
gc issue view 123 -R gitcode-cli/cli
git checkout -b bugfix/issue-123
# 直接开始修改...
```

### 正确示例

```bash
# ✅ 先验证问题是否存在
./gc issue view 123 -R gitcode-cli/cli

# 执行复现命令
./gc xxx command --params

# 补验证记录
gc issue comment 123 --body "## 验证记录

- 当前分支: main
- 复现命令: ./gc xxx command --params
- 实际结果: 复现成功
- 结论: 问题仍存在，进入修复" -R gitcode-cli/cli

# 更新状态后再开发
gc issue label 123 --add status/verified,status/in-progress -R gitcode-cli/cli
```

## 7. 最小标签集

Issue 和 PR 至少应使用以下标签维度：

- 类型：`type/bug`、`type/feature`、`type/docs`、`type/refactor`
- 状态：`status/triage`、`status/verified`、`status/in-progress`、`status/blocked`、`status/ready-for-review`、`status/merged`
- 风险：`risk/low`、`risk/medium`、`risk/high`
- 范围：开放命名空间，按命令/模块领域命名（示例：`scope/auth`、`scope/repo`、`scope/issue`、`scope/pr`、`scope/release`、`scope/docs`、`scope/api`、`scope/actions`、`scope/ci` 等，可按需扩展）

如果仓库当前尚未预建这些标签，先在规范和执行记录中按该命名约定使用。

## 8. 关闭规则

- issue 的“修复完成”不等于可以提前关闭
- 只有在 PR 已合入主干，或 issue 被明确判定为无需修复时，才能关闭 issue
- 如果 issue 已关闭但未合入主干，必须明确标记为“未完成主干合入”

## 9. 下一步去看哪里

- Issue 级状态与评论要求：看 [Issue 流程](./issue-workflow.md)
- PR 级状态、自检与合并要求：看 [PR 流程](./pr-workflow.md)
- 多角色评审角色与检查清单：看 [评审流程](./review-workflow.md)
- 本地与合并门禁：看 [代码质量门禁规范](../foundations/code-quality-gates.md)
- 远端 CI 与 AI 编排：看 [CI 工作流规范](../delivery/ci-workflows.md)
- AI 本地开发闭环编排：看 [AI 本地开发流程](./ai-local-development-workflow.md)

---

**最后更新**: 2026-07-13
