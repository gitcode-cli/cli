---
title: 飞书里程碑通报驱动的端到端交付——以 Issue #475 为例的 gitcode-delivery-notify 实战
description: 展示 gitcode-delivery-notify skill 如何用 gc lark send --to-self 在 AI 全流程交付（issue_verified→pr_created→ci→review→merged）的每个里程碑向本人飞书发送简短进展，并在修复阶段用 gc-api-doc 参考仓确认官方端点真相，以 bugfix #475→PR #448 为实例
---

# 飞书里程碑通报驱动的端到端交付

## 场景

AI 代理按 `spec/workflows/development-workflow.md` 跑一条完整的修复交付管线（Issue → 分支 → 修复 → 测试 → PR → CI → 评审 → 合并）。本案例在标准管线中织入两个横切资产，使用户**不必全程盯屏**、且修复涉及 API 时**不靠猜端点**：

1. **`gitcode-delivery-notify` skill**（横切编排器）：不替代 `gitcode-pr-create`/`gitcode-pr-review` 等交付 skill，只决定「何时、报什么」，委托 `gc lark send --to-self` 把简短摘要发到用户本人飞书（bot → user P2P，无需建群）。**在主流程的每个里程碑作为标准通报动作执行。**
2. **`gc-api-doc` 参考仓**（`git@gitcode.com:gitcode-cli/gc-api-doc.git`，官方 GitCode OpenAPI 全集）：修复涉及 API 端点/path/body 时，**作为权威真相源先查证再写代码**，不靠类比 GitHub CLI 或猜语义。

本案例以 **Issue #475**（`gc pr label` 同号跨资源误写 Issue 标签）→ **PR #448** → 合并为实例。该 bug 的修复恰恰是被 `gc-api-doc` 救活的：v1 方案用 `/merge_requests/{n}/labels` 实测 404，查 `gc-api-doc` 才确认官方端点是 `/pulls/{n}/labels`、body 是**裸 JSON 数组**（非 `{"labels":[]}` 对象）——这正是两个资产进主流程的典型场景。

**Issue #475 背景**：`gc pr label` 全程误调 Issue API（`GetIssue`/`AddIssueLabels`/`RemoveIssueLabel`），当 PR 与 Issue 同号时把标签写到 Issue。

## 推荐 skill 与参考资源

- `gitcode-delivery-notify` — 本案例主 skill（横切编排里程碑通报，独立仓库 [gitcode-cli/skills](https://gitcode.com/gitcode-cli/skills)）
- 参考资源：**`gc-api-doc` 仓**（`git@gitcode.com:gitcode-cli/gc-api-doc.git`）— 官方 OpenAPI 全集，修复 API 时的端点/path/body 权威真相源
- 配套交付 skill（主流程中实际使用，通报由 delivery-notify 编排）：
  - `gitcode-issue-review`、`gitcode-pr-create`、`gitcode-pr-review`、`gitcode-pr`、`gitcode-regression`

> 注意：`gitcode-delivery-notify` 是 standalone 编排器，不修改其他 skill；`gc-api-doc` 是只读参考，不进依赖；交付边界仍以 `spec/workflows/*` 为准。

## 适用人群

- 想让 AI 代理跑长交付任务、自己只看飞书进展的用户
- AI 代理（需要知道何时通报、报什么、修复 API 时先查权威端点）
- 流程设计者（建立「交付里程碑 → 个人通知」+「API 修复 → 先查证」机制）
- 异步协作团队

## 可直接执行的 Prompt

```text
请按 spec/workflows/development-workflow.md 全流程交付 gitcode-cli/cli 的 Issue #<N>。
主流程须织入两个横切资产：

A. 飞书里程碑通报（用 gitcode-delivery-notify skill）
1. 开始前确认 consent：读 ~/.config/gc/skills/lark-notify-consent.json；
   不存在且可交互时问我「通报哪些里程碑、发本人还是群、是否每条确认」，写入文件。
2. 每个里程碑到达即通报（简短：id+标题+URL+一行结论，禁止贴全文/日志/评审 dump）：
   - issue_verified：✅ Issue #N verified + 一行结论
   - pr_created：🆕 PR #N created + 标题 + html_url
   - ci：🔍 CI run #X {overall} + 每 job 一行 + run_url
   - review：👁️ Review PR #N: {approved|changes-requested} + 顶部发现 + html_url
   - merged：🎉 Merged PR #N ({method}) + issue #N closed
   - blocked：🛑 Blocked at {stage} + 阻塞点 + 需要我提供什么
3. 发送：gitcode lark send --to-self --markdown "<summary>" --json（--json 以便确认 ok:true）
4. 去重：同一 {milestone}:{id} 本会话只发一次；CI 抖动只发最终结论；review 状态变更才发。
5. 容错：发送失败只记一行 stderr 继续，不得重试刷屏，不得因通知失败阻断交付。

B. API 端点查证（修复涉及 API 调用时必做）
1. 先克隆/拉取 gc-api-doc 参考仓：git clone --depth 1 git@gitcode.com:gitcode-cli/gc-api-doc.git
2. 在 docs/api/<domain>/ 与 data/official/openapi.json 查目标端点的 path/method/body schema。
3. 端点/path/body 以 gc-api-doc 为准，不类比 GitHub CLI、不猜语义。
4. 现有 api 包函数若端点错误（如实测 404/400），按 gc-api-doc 修正路径与 body 形状。

安全：摘要只含 id/标题/URL/结论；不含 token、全文、评审 dump、日志；
gc lark send 提交前扫描 GC_TOKEN 值，命中拒绝，不要绕过。
```

## 预期产出

- 一条按里程碑推进的飞书消息流（本人 P2P），每条 ≤ 5 行
- 持久化 consent 文件（首次配置后复用）
- 修复涉及 API 时，端点/path/body 经 `gc-api-doc` 查证、与官方 OpenAPI 一致
- 交付本身完整完成（通报与查证是横切副产物，不阻断主线）
- 阻塞时一条明确「需要我提供什么」的消息

## 价值

- **异步可见**：长交付任务中用户无需盯屏，里程碑主动推送，阻塞才介入。
- **不靠猜端点**：`gc-api-doc` 提供官方 OpenAPI 全集作为权威真相源，避免类比 GitHub CLI 或猜语义导致的 404/400/跨资源误写（#475 正是猜 `/merge_requests/` 而非查 `/pulls/` 的代价）。
- **不刷屏**：模板强制简短 + 去重 + 状态变更才发，避免 CI 抖动产生噪音。
- **横切不侵入**：delivery-notify 只编排通报、gc-api-doc 只读参考，交付 skill 各司其职，可单独替换通报通道（本人/群）。
- **容错优先**：通知失败永不阻断交付；查证失败（仓不可达）则报告用户人工解决，不自行绕过。

## 复用方式

### 替换清单

| 占位符 | 案例值 | 替换为 |
|---|---|---|
| Issue | `#475` | 你的 Issue |
| PR | `#448` | 交付产生的 PR |
| 仓库 | `gitcode-cli/cli` | 目标仓库 |
| 通报目标 | `self`（本人 P2P） | `chat` + chat_id（需 bot 在群） |
| 里程碑集 | 全部 6 个 | 按需子集（如只 `merged`+`blocked`） |
| API 参考仓 | `gc-api-doc` | 项目对应的官方 OpenAPI 全集（如有） |

### 适用场景

- 任何按 spec 全流程交付且用户想要异步进展通知的任务
- AI 代理跑长管线（验证→修复→PR→CI→评审→合并）
- 修复涉及 GitCode REST API 调用（端点/path/body 不确定或现有函数实测报错）
- **不适合**：一次性问答、无需用户关注的短任务（通报反而打扰）；纯本地重构无 API 变更（可跳过 gc-api-doc 查证）

### 跨平台提醒

- opencode 的 user-level skill 装在 `~/.agents/skills/`（不是 `.claude/skills/`）；安装后需重启会话才进 `available_skills`，未进列表时手动按 SKILL.md 执行即可。
- `gc lark send --to-self` 自动解析当前 `lark-cli` 用户 open_id，无需手填；token 若 `needs_refresh`，发送时由 lark-cli 自动续期。
- `--to-self` 走 bot→user P2P，无需群成员资格；`--chat-id` 需 bot 在群且对外共享能力（外部群）。
- `gc-api-doc` 仓可一次性浅克隆到 `/tmp`，跨会话复用；OpenAPI 结构化定义在 `data/official/openapi.json`。

### 前置条件

- `gc` 已装且已认证（`GC_TOKEN`）
- `lark-cli` 已装且已登录（`lark-cli auth status` verified；`gc lark doctor` 通过）
- consent 已配置（文件或当次交互确认）
- 了解 `spec/workflows/development-workflow.md` 的里程碑定义
- （修复 API 时）`gc-api-doc` 仓可访问（SSH 或 HTTPS）

## 主流程（横切资产织入点）

标准交付管线 + 两个横切资产在每个阶段的标准动作：

| 阶段 | 交付动作 | 横切资产动作 |
|------|----------|--------------|
| 1. Verified | 复现 + 根因定位 + issue 评论 + label `status/verified,in-progress` | **delivery-notify**：发 `issue_verified` |
| 2. 分支 | 从 main 建非 main 分支 | — |
| 3. 修复 | 写代码 | **gc-api-doc**：涉及 API 调用时先查官方端点/path/body，现有 api 函数端点错误则按查证修正 |
| 4. 测试 | 单测 + 本地门禁（build/test/race/pre-commit）+ infra-test 实测 + 回归/系统测试 | — |
| 5. 风险/安全/docs | `classify-change-risk.py` + 安全审查 + 文档同步 | — |
| 6. PR | 自检 + `gc pr create`（body `Closes #N`）+ 标签 | **delivery-notify**：发 `pr_created` |
| 7. CI | GitCode 原生 CI（`.gitcode/workflows/ci.yml`） | **delivery-notify**：CI 完成发 `ci`（抖动只发最终） |
| 8. 评审 | 独立执行主体多角色评审；P0 修复、P1 should-fix、P2 建 follow-up Issue | **delivery-notify**：发 `review`（状态变更才发） |
| 9. 合并 | `gc pr merge` + 验证 issue 自动关闭 | **delivery-notify**：发 `merged` |
| (阻塞) | 某门禁无法继续（需人工确认/外部审批） | **delivery-notify**：发 `blocked`（含需要用户提供什么） |

> 关键：delivery-notify 在每个里程碑是**标准动作**而非可选装饰；gc-api-doc 在「修复」阶段涉及 API 时是**标准查证动作**而非事后补救。

## 本次真实执行记录

### 执行信息

- **执行时间**：2026-08-07
- **交付对象**：[Issue #475](https://gitcode.com/gitcode-cli/cli/issues/475) bug: pr label 误写同编号 Issue 标签
- **产出 PR**：[#448 fix(pr): pr label targets PR labels endpoint](https://gitcode.com/gitcode-cli/cli/merge_requests/448)
- **合并**：`925edb4 !448 merge`，Issue #475 自动关闭
- **通报里程碑**：5 个（issue_verified / pr_created / ci / review / merged；无 blocked）
- **通报通道**：本人 P2P（吴鹏飞，`ou_4968447f6c5aeec6638765d561e5c588`）
- **API 查证**：`gc-api-doc` 仓确认 `/pulls/{n}/labels` + 裸数组 body（救活 v1 的 404）

### 阶段 0：横切资产就位

```bash
# A. delivery-notify skill 装 opencode user-level 目录
cp -r /tmp/skills/gitcode-delivery-notify ~/.agents/skills/

# consent 文件（用户当次授权，落盘复用）
mkdir -p ~/.config/gc/skills
cat > ~/.config/gc/skills/lark-notify-consent.json <<'JSON'
{ "version":1, "target":"self", "chat_id":"",
  "milestones":["issue_verified","pr_created","ci","review","merged","blocked"],
  "confirm_each":false }
JSON
gc lark send --to-self --markdown "consent test" --json   # ok:true

# B. gc-api-doc 参考仓就位
git clone --depth 1 git@gitcode.com:gitcode-cli/gc-api-doc.git /tmp/gc-api-doc
```

### 里程碑通报（真实消息）

| 里程碑 | 触发信号 | 真实摘要 | message_id |
|--------|----------|----------|------------|
| `issue_verified` | Issue #475 加 `status/verified`+`in-progress` | `✅ Issue #475 verified` / 根因：pr/label/label.go 全程调 Issue API… | `om_x100b6860e759e888b17f1cfd36d720f` |
| `pr_created` | `gc pr create` 返回 number=448 | `🆕 PR #448 created` / 根因+修复一行结论 + URL | `om_x100b6860a34540b4b1247e757a72ed1` |
| `ci` | CI run #74 四 job 全 COMPLETED | `🔍 CI run #74 (PR #448) all passed` / Lint/Test/Build/Package ✅ + URL | ok:true |
| `review` | 独立评审 verdict=approved | `👁️ Review PR #448: approved (with nits)` / 无 P0，P1 已修，P2 建 Issue #498 + URL | ok:true |
| `merged` | `gc pr merge` merged:true + Issue closed | `🎉 Merged PR #448 (merge)` / issue #475 closed + 一行根治结论 | ok:true |

发送命令（统一形态）：

```bash
gc lark send --to-self --markdown "<short summary with URL>" --json
# 返回 { ok:true, data:{ message_id:"om_x..." } }
```

### 修复阶段的 gc-api-doc 查证（关键转折）

v1 方案用 `api.AddLabelsToPR`（走 `/merge_requests/{n}/labels`、body `{labels:[]}`），infra-test 实测 POST 返回 **404**。转折点：查 `gc-api-doc`：

```bash
# 1. 找 PR 标签端点文档
ls /tmp/gc-api-doc/docs/api/pull-requests/ | rg label
# → post/put/delete/get-api-v5-repos-owner-repo-pulls-number-labels*.md

# 2. 解析 openapi.json 的 path/body schema
python3 -c "
import json; d=json.load(open('/tmp/gc-api-doc/data/official/openapi.json'))
for p,m in d['paths'].items():
    if 'pulls/{number}/labels' in p:
        for meth,info in m.items():
            print(meth.upper(),p,'| body=',info.get('requestBody',{}).get('content',{}))
"
# → POST /api/v5/repos/{owner}/{repo}/pulls/{number}/labels  body: application/json, schema {"type":"array","items":{"type":"string"}}
# → PUT  同上（裸数组）
# → DELETE /api/v5/repos/{owner}/{repo}/pulls/{number}/labels/{name}（无 body）
```

结论：官方端点是 `/pulls/{n}/labels`（非 `/merge_requests/`），body 是**裸 JSON 数组**（非 `{labels:[]}` 对象）。据此修 `api.AddLabelsToPR`/`RemoveLabelFromPR`/`SetPRLabels`：路径改 `/pulls/`，`client.Post(path, labels, &result)` 直接传 `[]string`（`json.Marshal` → 裸数组）。重构建后 infra-test PR34/Issue34 同号实测 add/list/remove 全只动 PR。

### 其他要点

- **真实命令验证用 infra-test**：建 open PR #34 + 同号 Issue #34，新构建验证 add/list/remove 只动 PR、Issue 不变，验证后关闭 PR、删分支、清 marker。
- **系统测试沙箱坑**：测试沙箱不透传 `GC_TOKEN`（仅 env、无 auth.json），用 `gc auth login --with-token`（stdin 管道，不打印 token）持久化 auth.json 后通过；自检注明非代码问题。
- **评审发现前置建 Issue**：P2 预存在项（`--add` 无 trim/dedup、`SetPRLabels` 死代码、docs 自动推断名单）建 Issue #498 跟踪，合入前完成。

### 关键经验

1. **gc-api-doc 是 API 修复的权威真相源**：`/merge_requests/.../labels` 实测 404 后，靠 `gc-api-doc` 确认官方 `/pulls/.../labels` + 裸数组 body，否则会误判为「PR 不支持专用标签端点」而改走 `EditPR` 表单替换（语义不等价）。**修复涉及 API 调用时先查证再写代码，不靠类比 GitHub CLI 或猜语义。**
2. **通报是横切副产物，不阻断主线**：任何 `gc lark send` 失败只记一行 stderr 继续，绝不因通知失败卡住交付；本案 5 条全 ok，但即便失败也不影响 #475 合入。
3. **consent 必须前置且可持久化**：用户当次口头授权后落盘 `~/.config/gc/skills/lark-notify-consent.json`，后续会话免问；非交互 + 无 consent → 静默跳过，不报错。
4. **模板强制简短 + 去重**：摘要只含 id/标题/URL/结论；同一 `{milestone}:{id}` 本会话只发一次；CI 多次跳变只发最终结论。
5. **skill 装对位置**：opencode 是 `~/.agents/skills/`，不是 `.claude/skills/`；装错位置虽可能被部分 harness 加载，但会污染仓库工作树且不符合 opencode 约定。
6. **`--to-self` 免解析身份**：自动从 `lark-cli auth status` 取本人 open_id，无需 lark-contact 查找；token `needs_refresh` 由 lark-cli 自动续期。
7. **`blocked` 才需用户介入**：本案无阻塞故未触发；一旦某门禁无法继续，发 `🛑 Blocked at {stage}` + 需要用户提供什么，用户据此回归。
8. **安全双保险**：摘要不含 token/全文/日志；`gc lark send` 提交前扫描当前 `GC_TOKEN` 值，命中拒绝发送；飞书凭证留在 lark-cli OS keychain。
9. **通报与交付 skill 解耦**：delivery-notify 只编排「何时报什么」，交付动作仍由各自 skill/API 完成；通报通道可单独替换而不动交付流程。

## 相关案例

- 交付主线：[AI 全流程交付——从 Issue 到合并](./ai-full-delivery-workflow.md) — 本案例通报的交付管线本身（以 bugfix #250 为例的 15 阶段）
- 评审环节：[PR 评审反馈闭环](./review-feedback-loop.md) — 本案例 `review` 里程碑的评审环节展开
- 认证前置：[多环境 GitCode CLI 认证配置](./auth-setup.md) — `lark-cli` 登录与 `gc lark doctor` 前置
- 规范：`spec/workflows/development-workflow.md` — 里程碑信号定义
