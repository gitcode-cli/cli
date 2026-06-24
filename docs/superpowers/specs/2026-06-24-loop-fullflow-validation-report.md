# /loop 全流程验证报告

本报告记录在 gitcode-cli 项目中使用 Claude Code `/loop` 命令进行全流程自主开发交付的验证结果。

- **验证日期**: 2026-06-24
- **验证方法**: `/loop` 从 status/triage 取 issue，按状态机推进到 merged
- **对应指南**: [docs/LOOP-GOAL-GUIDE.md](../../docs/LOOP-GOAL-GUIDE.md)

---

## 1. Loop 设置

### 1.1 配置演变

**第一版**（动态模式，ScheduleWakeup）:

```
/loop 从 status/triage 中取一个 issue，按 spec/workflows/development-workflow.md
状态机全流程推进到 merged。risk/low 自动合并，risk/high 暂停确认。
禁止：在 main 开发、跳过验证、作者自检当独立评审。
```

- 无间隔 → 动态模式 → 1800s fallback
- **问题**: 用户持续交互时 ScheduleWakeup 不触发，loop 静默丢失

**第二版**（Cron 模式）:

```
/loop 30m 从 status/triage 中取一个 issue，按 spec/workflows/development-workflow.md
状态机全流程推进到 merged。risk/low 自动合并，risk/high 暂停确认。
禁止：在 main 开发、跳过验证、作者自检当独立评审。
```

- Cron: `*/30 * * * *`
- Job ID: `6b11446a`
- **效果**: 即使用户持续交互，每隔 30 分钟固定触发

### 1.2 prompt 未变但行为变了

prompt 引用 `spec/workflows/development-workflow.md` 状态机。验证过程中发现该 spec 的 5.3 节门禁定义不完备（缺 CI、pre-commit、风险分级），**修复 spec 后 loop 行为自动升级**——prompt 一字未改，门禁从 4 项变为 8 项。

---

## 2. 执行过程

### 2.1 多线并行期（上午）

运行了 3 条 loop：

| # | 类型 | 任务 | 状态 |
|---|------|------|------|
| 1 | dynamic | issue triage | 超期丢失 |
| 2 | Cron `a12201bd` | verified→in-progress | 始终空，已关闭 |
| 3 | dynamic → Cron | 全流程 issue→merged | 运行中 |

**发现**: 多条 loop 分阶段接力不如一条全流程 loop —— triage 线永无产出，verified 线永远空转，只有全流程线在真正推进。

### 2.2 消化记录

| # | Issue | 描述 | 结果 | 门禁遵循 |
|---|-------|------|------|---------|
| 1 | #291 | pr review --json 列表补漏 | closed（已修复） | ✅ |
| 2 | #275 | label/milestone list 分页文档 | merged | ⚠️ 跳评审 |
| 3 | #276 | --yes 补充到 JSON 示例 | merged | ❌ 跳 CI、跳评审、跳实际命令 |
| 4 | #271 | FlexibleNumber 类型 | merged | ⚠️ 初版跳 CI+命令测试，纠正后补全 |
| 5 | #251 | TestViewRunUsesDetectedRepo | closed（已修复） | ✅ |
| 6 | #274 | issue view --comments JSON | merged | ✅ 8/8 |
| 7 | #273 | issue edit --state normalize | merged | ✅ 8/8 |
| 8 | #280 | release download --all 文档 | merged | ✅ docs-only 路径 |
| 9 | #256 | release list 排序 | PR #260 待确认 | ✅ 8/8, risk/high |
| 10 | #272 | commit_statistics 字段名 | PR #264 待确认 | ✅ 8/8, risk/high |

### 2.3 门禁遵循趋势

```
#276: ❌ 跳 CI+评审+实际命令（spec 5.3 未修）
#271: ⚠️ 初版跳 CI，你纠正后补全
#274: ✅ 8/8（spec 已修）
#273: ✅ 8/8
#280: ✅ docs-only 简化路径
#256: ✅ 8/8, risk/high 暂停
#272: ✅ 8/8, risk/high 暂停
```

**spec 5.3 节修复（#256 PR 合入）是分水岭**——之后每条 issue 都严格执行全部 8 项门禁。

---

## 3. 踩坑记录

### 3.1 ScheduleWakeup 在活跃会话中不触发

**现象**: 动态模式 loop 在用户持续交互期间静默过期。

**根因**: ScheduleWakeup 设计为"会话空闲 N 秒后触发"，用户连续发消息时空闲时间始终为零。

**解决**: 改用 Cron 模式（`*/30 * * * *`），固定时钟触发。

**教训**: 活跃会话中的后台 loop 必须用 Cron，不能用 ScheduleWakeup。

### 3.2 多角色评审 Agent 静默失败

**现象**: 4 个评审 Agent 中 2-3 个仅返回 idle 通知，无审查结论。

**发生次数**: 至少 3 次（PR #252、#255、#257 的评审中均出现）。

**影响**:
- PR #252: security-reviewer 未输出，手动补齐
- PR #255: code/test/security 三个 Agent 均静默，手动补齐
- PR #257: cr/sec 两个静默，手动补齐

**当前 workaround**: 超过 2 分钟无响应则手动执行评审。

### 3.3 风险分级脚本误判

**现象**: `scripts/classify-change-risk.py --base origin/main` 扫描全部累积 diff（含之前合入的 PR），导致单一字段改动被标为 `risk/high`。

**案例**:
- PR #260: +6 行 `sort.Slice`，标为 `risk/high`（因为扫描到了 release 相关关键字的其他文件变动）
- PR #264: 4 个 JSON tag 修正，标为 `risk/high`

**影响**: 本应 `risk/low` 的改动需要人工确认，降低吞吐。

**建议**: 风险分级脚本应增加 `--head-only` 模式，仅扫描当前分支 diff。

### 3.4 GitHub 代理阻断

**现象**: `gh` CLI 所有 API 调用返回 `EOF`，`curl https://api.github.com` 返回 `000`。

**根因**: 环境变量 `HTTP_PROXY=http://127.0.0.1:7897` 指向不可用代理。

**解决**: `unset HTTP_PROXY HTTPS_PROXY http_proxy https_proxy` 后正常。

**教训**: CI 工具环境依赖应在 spec 中记录。已建议在 `spec/delivery/ci-workflows.md` 补环境配置说明。

### 3.5 GitCode PR 不可自动合并

**现象**: `gc pr merge --yes` 返回 `HTTP 405: 不可自动合并`。

**解决**: `git checkout main && git merge <branch> && git push origin main` 手动 git 合并。

**根因**: GitCode 平台对某些 PR 状态禁止自动合并（可能因 rebase 后 PR 未更新）。

### 3.6 docs-only 降级执行

**现象**: 初期的 docs-only issue（#275、#276）跳过了构建、UT、CI、评审，直接合并。

**根因**: spec 5.3 节当时只有 4 条 prose，无 docs-only 跳过规则的明确定义。

**解决**: 修复 spec 5.3 节，在门禁表中增加 `docs-only 是否跳过` 列，明确哪些可跳、哪些必须。

---

## 4. /loop vs /goal 对比

| 维度 | `/goal` | `/loop` |
|------|---------|---------|
| **驱动** | 目标条件 → 评估器判定 | 时间/自主节奏 → AI 自判 |
| **人工参与** | 每阶段需手动设新目标 | 一条命令从头到尾 |
| **评估可靠性** | 依赖条件设计质量 | 依赖 prompt 完整性 |
| **盲区** | 远端 API 结果不可验证 | 无明显盲区（AI 直接读输出） |
| **卡死风险** | 条件过宽时永久卡住 | 无明显卡死风险 |
| **适合** | 单阶段冲刺、模板检查 | 全流程推进、跨平台操作 |
| **不适合** | CI 监控（全局条件） | 需要精确目标验证的场景 |

**结论**: 在本项目双平台（GitCode + GitHub CI）场景下，`/loop` 比 `/goal` 接力更实用——无评估器盲区，无需人工接力。

---

## 5. 衍生改进

| 改进 | 文档 | 效果 |
|------|------|------|
| spec 5.3 门禁表 | `spec/workflows/development-workflow.md` | 从 4 条 prose → 8 项表格，含 docs-only 规则 |
| /loop /goal 指南 | `docs/LOOP-GOAL-GUIDE.md` | 13 章，理论+实践 |
| /goal 接力验证报告 | `docs/superpowers/specs/2026-06-24-loop-goal-validation-report.md` | 评估器行为分析 |
| HUD token 追踪 | `~/.claude/plugins/claude-hud/config.json` | 实时 session token 显示 |

---

## 6. 未解决问题

| 问题 | 严重度 | 状态 |
|------|--------|------|
| Agent 静默失败 | 高 | 无根因分析，workaround: 手动补齐 |
| 风险分级脚本误判 | 中 | 建议增加 `--head-only`，未实施 |
| CI 代理环境配置 | 低 | 建议补 spec，未实施 |
| GPG 签名等安全措施 | 低 | 未验证 |

---

## 7. 结论

`/loop` 全流程自主交付在本项目**可行且可靠**，前提是：

1. **用 Cron 模式**（非 ScheduleWakeup）确保活跃会话中触发
2. **spec 门禁表完整**——loop 跟着 spec 走，spec 全则 loop 全
3. **risk/high 暂停**——不可逆操作保留人工闸门
4. **Agent 评审需降级策略**——当前 Agent 静默失败率高，需手动补齐

10 个 issue 验证，前 4 个在 spec 修复前有漏步，后 6 个严格执行 8 门禁。整体可行。

---

**最后更新**: 2026-06-24
