---
name: gitcode-label-stats
description: |
  Statistics of PR label distribution across repositories in a GitCode organization.
  Supports multiple labels, monthly breakdown, and ratio analysis.
  TRIGGER when: user asks to count PRs by label, measure label adoption rate,
  compare labeled vs total PRs, generate label statistics report, or phrases like
  "标签统计", "ai-assisted占比", "PR标签分布", "标签覆盖率", "label stats",
  "label adoption", "组织PR统计".
---

# gitcode-label-stats

统计 GitCode 组织下各仓库的 PR 标签分布，支持多标签、月度维度和占比分析。

## 适用范围

- 按组织维度统计所有仓库的 PR 标签分布
- 支持同时追踪多个标签（如 `ai-assisted`、`bug`、`enhancement`）
- 支持按月度维度统计标签 PR 数量及占比
- 输出 CSV 报表，便于进一步分析

## 前置条件

- GitCode CLI 已安装且已认证：`gitcode auth status`
- 拥有目标组织的读取权限
- 跨平台：Windows 使用 `gitcode`，Linux/macOS 可使用 `gc`

## 工作流程

### Step 1: 获取组织仓库列表

```bash
gitcode repo list -o <org> --limit 100 --json
```

- 若组织仓库超过 100 个，需分页获取（增大 `--limit` 或多次调用）
- 从 JSON 输出中提取每个仓库的 `name` 字段
- 注意：列表中可能包含已删除的幽灵仓库（列表有缓存），后续步骤会自然跳过

### Step 2: 逐仓库获取 PR 列表

对每个仓库，使用 `gitcode pr list` 获取所有 PR：

```bash
gitcode pr list -R <org>/<repo> --state all --limit 100 --json
```

- `--state all` 包含已合并和已关闭的 PR
- 若 PR 数量超过 100 个，需分页处理：多次调用并递增 `--page` 参数（如 `--page 2`）
- 若仓库不存在或已删除，CLI 会返回错误，直接跳过该仓库

### Step 3: 解析 PR 标签与时间

对每个 PR 的 JSON 数据，提取以下字段：

| 字段 | 用途 |
|------|------|
| `number` | PR 编号 |
| `created_at` | 创建时间（用于月度统计） |
| `labels[].name` | 标签名称列表 |

判断 PR 是否携带目标标签：

```text
遍历 pr.labels 数组，检查 label.name 是否在目标标签集合中
```

按 `created_at` 的年月（格式 `YYYY-MM`）归类到对应月份。

### Step 4: 汇总统计

对每个仓库计算：

| 统计项 | 计算方式 |
|--------|----------|
| 总 PR 数 | 该仓库所有 PR 计数 |
| 各目标标签 PR 数 | 携带对应标签的 PR 计数 |
| 标签占比 | 标签 PR 数 / 总 PR 数 × 100% |
| 月度总 PR 数 | 按年月筛选后的 PR 计数 |
| 月度标签 PR 数 | 按年月筛选后携带标签的 PR 计数 |
| 月度标签占比 | 月度标签 PR 数 / 月度总 PR 数 × 100% |

### Step 5: 输出报表

生成 **两份** CSV 文件：

#### 文件 1：全量数据（`<label>-pr-stats.csv`）

包含所有仓库，含标签 PR 数为 0 的仓库：

```text
Repository, Total PRs, <label> PRs, <label> Ratio,
<YYYY-MM> Total PRs, <YYYY-MM> <label> PRs, <YYYY-MM> <label> Ratio,
..., Overall <label> Ratio
```

#### 文件 2：过滤数据（`<label>-pr-stats-filtered.csv`）

过滤掉标签 PR 数为 0 的仓库，TOTAL 行基于过滤后的仓库重新计算（分子分母均不含被过滤仓库）：

- 过滤条件：`<label> PRs == 0` 的仓库行移除
- TOTAL 行重算：对过滤后的仓库行，累加各数值列，重新计算占比
- 适用于分析"实际使用该标签的仓库群体"的采纳情况

CSV 写入需使用 UTF-8 编码：

```powershell
# Windows PowerShell
[System.IO.File]::WriteAllLines($csvPath, $csvLines, [System.Text.Encoding]::UTF8)

# Linux/macOS
$csvLines | Out-File -FilePath $csvPath -Encoding utf8
```

## 常见问题

| 问题 | 原因 | 解决方案 |
|------|------|----------|
| 仓库列表中有仓库但 CLI 报错 | 仓库已删除，列表有缓存 | 跳过该仓库，记录为不可访问 |
| CLI JSON 输出解析失败 | 终端编码或输出折行 | 将输出写入临时文件后用 UTF-8 读取 |
| PR 数量超过 100 | 单页上限 | 使用 `--page` 参数逐页获取 |
| Windows 中文写入 CSV 乱码 | PowerShell 默认编码非 UTF-8 | 使用 `[System.IO.File]::WriteAllLines` 并指定 UTF-8 编码 |
| 标签 PR 占比显示 0% | 月度内无标签 PR 或除零 | 分母为 0 时占比显示 `0%` 或 `N/A` |

## 输出示例

### 文件 1：全量数据（`ai-assisted-pr-stats.csv`）

```text
Repository,Total PRs,ai-assisted PRs,ai-assisted Ratio,2026-04 Total,2026-04 ai-assisted,2026-04 Ratio,2026-05 Total,2026-05 ai-assisted,2026-05 Ratio
repo-a,340,26,7.6%,95,0,0%,99,6,6.1%
repo-b,369,12,3.3%,33,0,0%,29,2,6.9%
repo-c,22,11,50%,8,1,12.5%,4,4,100%
repo-d,50,0,0%,10,0,0%,12,0,0%
repo-e,30,0,0%,5,0,0%,8,0,0%
TOTAL,811,49,6%,151,1,0.7%,152,12,7.9%
```

### 文件 2：过滤数据（`ai-assisted-pr-stats-filtered.csv`）

移除 ai-assisted PRs 为 0 的仓库（repo-d、repo-e），TOTAL 行重新计算：

```text
Repository,Total PRs,ai-assisted PRs,ai-assisted Ratio,2026-04 Total,2026-04 ai-assisted,2026-04 Ratio,2026-05 Total,2026-05 ai-assisted,2026-05 Ratio
repo-a,340,26,7.6%,95,0,0%,99,6,6.1%
repo-b,369,12,3.3%,33,0,0%,29,2,6.9%
repo-c,22,11,50%,8,1,12.5%,4,4,100%
TOTAL,731,49,6.7%,136,1,0.7%,132,12,9.1%
```

## 使用示例

以下为用户向 AI agent 发出的提示词示例，agent 收到后按本 skill 工作流程执行。

### 示例 1：单标签 + 月度维度（完整示例）

> 统计 openLiBing 组织下所有仓库的 ai-assisted 标签 PR 数量，加上每个仓库的总 PR 数、4 月和 5 月的月度数据及占比，输出 CSV

agent 行为：Step 1 获取仓库列表 → Step 2 逐仓库获取 PR → Step 3 匹配 `ai-assisted` 标签，按 `2026-04` 和 `2026-05` 归类 → Step 4 汇总 → Step 5 输出两份 CSV（全量 + 过滤）

输出列：Repository, Total PRs, ai-assisted PRs, ai-assisted Ratio, 2026-04 Total, 2026-04 ai-assisted, 2026-04 Ratio, 2026-05 Total, 2026-05 ai-assisted, 2026-05 Ratio

### 示例 2：多标签 + 月度维度

> 统计 openLiBing 组织各仓库的 PR 标签分布，标签包括 ai-assisted 和 bug，加上 4 月和 5 月的月度数据，输出 CSV

agent 行为：同上流程，但标签集合为 `{ai-assisted, bug}`，月度筛选 `2026-04` 和 `2026-05`，CSV 列包含两组标签列和两组月度列。

### 示例 3：指定月份范围

> 统计 my-org 组织的 enhancement 标签 PR，按季度维度（2026-Q1 和 Q2）输出

agent 行为：月度筛选改为按季度聚合（Q1 = 01/02/03 月，Q2 = 04/05/06 月），其余流程不变。

### 示例 4：仅统计特定仓库子集

> 统计 openLiBing/openlibing-cicd 和 openLiBing/openlibing-web 两个仓库的 ai-assisted PR 占比

agent 行为：跳过 Step 1，直接对指定仓库执行 Step 2-5。

### 示例 5：英文提示词

> Count PRs with label "ai-assisted" across all repos in the openLiBing org, include monthly breakdown for April and May

agent 行为：同示例 2，标签为 `{ai-assisted}`，月度为 `2026-04` 和 `2026-05`。

## 规则

- 使用 `gitcode` 命令确保跨平台兼容性（Windows 下 `gc` 是 `Get-Content` 别名）
- 所有含中文的远端写操作必须用 `--body-file` + UTF-8 编码写盘
- CSV 输出必须使用 UTF-8 编码
- 遇到不可访问的仓库时跳过并记录，不中断整体统计
- 分母为 0 时占比显示 `0%`，不执行除法
- 大量仓库时注意 API 速率限制，适当添加延迟
- **必须生成两份 CSV**：全量版和过滤版（过滤掉标签 PR 数为 0 的仓库）
- 过滤版的 TOTAL 行必须基于过滤后的仓库重新累加计算，不能复用全量版的 TOTAL
- 若当前月份不是完整月份（如月中统计），列名需标注截止日期，如 `Jun Total PRs (as of 2026-06-09)`
