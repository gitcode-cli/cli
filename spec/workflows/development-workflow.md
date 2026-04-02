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
- 把作者自检当作独立评审
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
3. 作者实现、自检、独立评审必须分离
4. 本地自动检查和人工语义审查必须分层
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
| `status/ready-for-review` | 可以进入独立评审 | 作者 |
| `status/changes-requested` | 发现问题，需继续修改 | 评审者 |
| `status/approved` | 评审通过，可合并 | 评审者 |
| `status/merged` | 已合入主干 | 合并后更新 |

## 3. AI 执行边界

AI 在本仓库中必须遵守以下硬约束：

- 未进入 `status/verified` 的 issue，不得开始写代码
- 未创建开发分支前，不得修改实现文件
- 未完成自检证据前，不得把 PR 标记为 `status/ready-for-review`
- 不得把作者自己的 `gc pr review --comment` 当作独立评审通过
- 不得在 PR 未合入 `main` 前把待修复 issue 标记为完成
- 不得声称“已完成”但不给出测试、命令验证、文档同步证据

## 4. 标准流程

```
Issue 创建/受理
→ issue triage
→ issue verified
→ 创建分支
→ 开发与测试
→ issue in-progress
→ PR draft
→ 作者自检
→ issue / PR ready-for-review
→ 独立评审
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

本阶段必须完成：

- 开发实现
- 编写或补齐测试
- 执行本地构建与相关测试
- 执行至少一个真实命令验证

### 5.4 Self Check

代码开发完成后，作者必须先做自检，不能直接等同于通过评审。

作者自检至少包含：

- 根因或实现理由
- 修改范围
- 测试结果
- 实际命令验证结果
- 文档同步结果
- 风险点
- 未覆盖项

### 5.5 Ready For Review

只有在以下条件全部满足后，才允许把 issue / PR 推到 `status/ready-for-review`：

- 自检证据完整
- 本地门禁已完成
- 真实命令验证已完成，或明确写出未执行原因
- 文档同步已完成，或明确写出无需更新的依据

### 5.6 Merge

只有在以下条件全部满足后，才允许合并：

- 无 blocker 级问题
- 存在独立评审结论
- PR 状态为 `status/approved`
- 相关 issue 和 PR 的状态、标签、评论记录完整

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
- 范围：`scope/auth`、`scope/repo`、`scope/issue`、`scope/pr`、`scope/release`、`scope/docs`

如果仓库当前尚未预建这些标签，先在规范和执行记录中按该命名约定使用。

## 8. 关闭规则

- issue 的“修复完成”不等于可以提前关闭
- 只有在 PR 已合入主干，或 issue 被明确判定为无需修复时，才能关闭 issue
- 如果 issue 已关闭但未合入主干，必须明确标记为“未完成主干合入”

## 9. 下一步去看哪里

- Issue 级状态与评论要求：看 [Issue 流程](./issue-workflow.md)
- PR 级状态、自检与合并要求：看 [PR 流程](./pr-workflow.md)
- 作者自检与独立评审边界：看 [评审流程](./review-workflow.md)
- 本地与合并门禁：看 [代码质量门禁规范](../foundations/code-quality-gates.md)

---

**最后更新**: 2026-04-02
