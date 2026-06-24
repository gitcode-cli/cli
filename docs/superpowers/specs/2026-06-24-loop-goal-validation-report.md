# /goal 接力验证报告

本报告记录在 gitcode-cli 项目中对 Claude Code `/goal` 命令的四阶段接力验证结果。

- **验证日期**: 2026-06-24
- **验证目标**: Issue #298（fetchURLHost IPv6 解析 bug）+ 附带修复 #302
- **验证方法**: `/goal` 接力——每阶段单独设置目标，观察评估器行为
- **对应指南**: [docs/LOOP-GOAL-GUIDE.md](../../docs/LOOP-GOAL-GUIDE.md)

---

## 1. 验证设计

### 1.1 四个阶段

```
阶段 1: 开发+本地验证        阶段 2: 推送+CI 监控
/goal until test+build       /goal until CI success

阶段 3: 安全审查+自检        阶段 4: 标签+PR 收尾
/goal until self-check       /goal until label+PR
```

### 1.2 观察维度

- 评估器能否正确识别 PASS/FAIL 信号
- 评估器在涉及远端 API 的条件上是否有盲区
- 评估器在模板完整性检查上的表现
- 条件设计对 `/goal` 成功率的决定性影响

---

## 2. 各阶段详细记录

### 阶段 1: 开发 + 本地验证

**目标条件**:
```
go test ./pkg/cmd/pr/checkout/... passes
AND go test ./... passes
AND go build -o ./gc ./cmd/gc succeeds
```

**实际产出**:
- `fetchURLHost` IPv6 括号修复（+38 行）
- `scpHost` 辅助函数
- `TestFetchURLHost`（9 用例，含 4 个 IPv6）+ `TestScpHost`（2 用例）
- checkout 包: 16 passed | 全项目: 1187 passed | 构建: Success

**评估器判定**: ✅ 正确识别 PASS

**对话中可见的证据**:
```
Go test: 16 passed in 1 packages
Go test: 1187 passed in 95 packages
Go build: Success
BUILD: SUCCESS
```

**评价**: 评估器最强项。`go test` 输出中的 PASS/FAIL 和构建成功/失败是标准化的二进制信号，无歧义。

---

### 阶段 2: 推送 + CI 监控

**目标条件**:
```
gh run list --workflow=ci.yml --branch bugfix/issue-298 shows success
AND the CI run URL is recorded
```

**实际发生**:
1. 推送分支到 GitHub 镜像仓 → `workflow_dispatch` 触发 CI
2. 第一次 CI 运行 (#28071094618) 失败——Test 在所有 3 个平台失败
3. 根因: `TestResolveRunMissingRepo` 预存 bug (#302)，CI 无 `GC_TOKEN`
4. 评估器**连续 4 次拒绝**——严格对照条件 "shows success"，failure ≠ success
5. 被迫修复 #302（`t.Setenv("GC_TOKEN", "dummy-token")`）→ 重新触发 CI
6. 第二次 CI (#28071396966) 全绿 → 评估器通过

**评估器行为**: ✅ 严格且正确

**核心发现**:

| 发现 | 详细说明 |
|------|---------|
| 条件判定准确性 | CI 结论 `failure` 时评估器绝不放行——即使失败与本次改动无关 |
| 盲区 | 评估器无法区分"我的改动导致失败" vs "预存 bug 导致失败" |
| 卡死风险 | 全局条件 `shows success` 在有多预存失败的项目中可能永久无法达成 |
| 条件设计教训 | 应限定为 "涉及本次改动的测试通过"，而非 "全项目 succeeds" |

**关键教训**: `/goal` 用于 CI 监控时，条件必须缩小范围。推荐写法：

```bash
# ❌ 过于宽泛——预存失败会导致永远卡住
/goal until CI shows success

# ✅ 限定到本次改动相关包
/goal until CI test jobs pass for packages modified in this PR
AND the CI run URL is recorded
```

---

### 阶段 3: 安全审查 + 自检模板

**目标条件**:
```
issue #298 PR self-check is complete:
  - 根因、修改范围、测试结果、CI 证据
  - 安全审查结果、文档同步说明、风险分级
  - 未覆盖项已列出
AND add comment to issue #298 with the self-check record
```

**评估器判定**: ✅ 通过（无声息清除）

**自检模板 9 项完成情况**:

| # | 字段 | 内容 | 评估器可判定性 |
|---|------|------|--------------|
| 1 | 根因 | IPv6 括号解析 bug + #302 CI auth | 文本非空 → 可判定 |
| 2 | 修改范围 | 3 files, +88/-1 | 文本非空 → 可判定 |
| 3 | 测试结果 | 1187 passed | 含 PASS 标记 → 可判定 |
| 4 | CI 证据 | run ID + URL | 含 URL → 可判定 |
| 5 | 安全审查 | 无硬编码凭证 | 文本非空 → 可判定 |
| 6 | 文档同步 | 无需更新 | 文本非空 → 可判定 |
| 7 | 风险 | risk/low | 文本非空 → 可判定 |
| 8 | 未覆盖项 | 2 items | 文本非空 → 可判定 |
| 9 | Comment 发布 | ID 176846141 | "Added comment" → 可判定 |

**评价**: 模板完整性检查是评估器的可靠领域。每项只需要 "有内容 vs 空" 的二元判断，无歧义。

---

### 阶段 4: 标签 + PR 收尾

**目标条件**:
```
gc issue view 298 shows status/self-checked
AND a PR has been created for branch bugfix/issue-298
```

**实际产出**:
- Issue #298: `Labels: status/self-checked` ✅
- PR #252: `✓ Created PR #252 in gitcode-cli/cli` ✅
- 评估器判定: ✅ 通过（无声息清除）

**评价**: `gc` 输出中的 "Created PR" 和 "Labels: ... status/self-checked" 是明确的文本信号。

---

## 3. 跨阶段发现

### 3.1 多角色评审——必须手动

阶段 4 完成后，PR #252 进入多角色评审。按项目规范，评审使用 TeamCreate + 4 个独立 Agent：

| Agent | 结论 | 耗时 |
|-------|------|------|
| code-reviewer | ✅ approved | ~2min |
| test-reviewer | ✅ approved | ~2min |
| docs-reviewer | ✅ approved | ~2min |
| security-reviewer | ⚠️ 无输出 | N/A |

安全审查 Agent 出现会话异常——持续发送 idle 通知但未输出审查结论。**这是一个重要的可靠性观察**：Agent 在评审流程中可能静默失败，需要超时机制或手动补齐。

### 3.2 `/goal` vs `/loop` 选择指南

根据本次验证，更新选择矩阵：

| 场景 | 推荐 | 原因 |
|------|------|------|
| 本地可验证的单一阶段 | `/goal` | 评估器可靠，无盲区 |
| 涉及远端 API（CI、PR 状态） | `/loop` | 评估器无法直接验证 API 响应 |
| 模板/检查表完整性 | `/goal` | 二元判定，无歧义 |
| 需要独立语义判断（评审） | 手动 | 不能自动化 |
| 跨平台全流程 | `/loop` | 无评估器盲区问题 |

### 3.3 条件设计的黄金法则

```
条件 = 可观测 + 可衡量 + 限定范围
```

| 条件 | 可观测 | 可衡量 | 限定范围 | 评级 |
|------|--------|--------|---------|------|
| `go test ./... passes` | ✅ | ✅ | ❌ 太宽 | 中等 |
| `go test ./pkg/cmd/pr/checkout/... passes` | ✅ | ✅ | ✅ | 优秀 |
| `CI shows success` | ✅ | ✅ | ❌ 太宽 | 差 |
| `PR self-check template complete` | ✅ | ✅ | ✅ | 优秀 |
| `code quality is good` | ❌ | ❌ | ❌ | 不可用 |

---

## 4. 发现的问题与建议

### 4.1 评估器相关

| 问题 | 严重度 | 建议 |
|------|--------|------|
| 无法区分失败原因（我的改动 vs 预存 bug） | 中 | 条件限定到受影响包 |
| 远端 API 结果不可直接验证 | 高 | CI 类条件用 `/loop` 替代 `/goal` |
| Agent 静默失败（security-reviewer） | 中 | 评审流程增加超时 + 降级策略 |

### 4.2 流程相关

| 问题 | 严重度 | 建议 |
|------|--------|------|
| `workflow_dispatch` 需手动触发 CI | 低 | 在 `docs/LOOP-GOAL-GUIDE.md` 中补充触发说明 |
| `/goal` 卡死时用户需手动 `/goal clear` | 中 | 已知局限，在指南中明确说明 |
| 双平台（GitCode + GitHub）操作无统一入口 | 低 | 在技能或脚本中封装双平台 push+CI 触发 |

---

## 5. 结论

### 5.1 `/goal` 接力在本项目的可行性

**可行，但需遵守条件设计原则**。

- ✅ 阶段 1（本地验证）、阶段 3（自检模板）、阶段 4（标签收尾）——评估器表现可预测且正确
- ⚠️ 阶段 2（CI 监控）——条件过宽导致卡死，缩小范围后通过
- ⛔ 多角色评审——不可自动化，必须手动

### 5.2 推荐用法

```bash
# 精细模式（推荐给复杂改动）
/goal until <阶段 1 条件>  →  /goal until <阶段 3 条件>  →  /goal until <阶段 4 条件>
                          ↑
                      CI 用 /loop 监控（不用 /goal）

# 快速模式（推荐给小改动）
/loop   # 一条命令覆盖全部，配合 multi-role review 手动触发
```

### 5.3 未验证项

- `/loop` 自主模式的端到端验证（本次仅验证了 `/goal` 接力）
- 多 Agent 评审的可靠性（security-reviewer 的静默失败需要进一步调查）
- GPG 签名等安全措施的交互

---

## 6. 相关文档

| 文档 | 说明 |
|------|------|
| [docs/LOOP-GOAL-GUIDE.md](../../docs/LOOP-GOAL-GUIDE.md) | /loop 与 /goal 使用指南 |
| [spec/workflows/review-workflow.md](../../spec/workflows/review-workflow.md) | 多角色评审规范 |
| [spec/delivery/ci-workflows.md](../../spec/delivery/ci-workflows.md) | CI 工作流规范 |
| Issue #298 | fetchURLHost IPv6 解析 bug |
| Issue #302 | CI test auth 预存问题 |
| PR #252 | 本次验证关联的实现 PR |

---

**最后更新**: 2026-06-24
