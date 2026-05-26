---
title: 定位引入问题的 PR — 以 Ascend/pytorch 测试用例突变为例
description: 使用 GitCode CLI 在 PR 海中快速定位引入回归问题的 PR，结合 gitcode-pr-impact-locator skill 实现自动化排查
---

# 定位引入问题的 PR — 以 Ascend/pytorch 测试用例突变为例

## 场景

在 Ascend/pytorch（torch-npu）v2.7.1 分支的日常开发中，发现测试用例泛化流程产出的结果突变：用例数从 67,693 降至 34,295，降幅 49.4%。缺失的用例集中在 `TestCommonCPU` 类，而 `TestCommonPRIVATEUSE1` 保留。

需要在最近合入的数十个 PR 中快速定位是哪个 PR 引入了问题。

## 推荐 skill

- `gitcode-pr-impact-locator`

## 可直接执行的 Prompt

```text
我需要定位 Ascend/pytorch 仓库 v2.7.1 分支上导致测试用例数量从 67k 降至 34k 的 PR。请使用 gitcode-pr-impact-locator skill。

问题详情：
- 仓库：Ascend/pytorch（https://gitcode.com/Ascend/pytorch）
- 分支：v2.7.1
- 时间范围：最近 10 天（2026-05-15 ~ 05-25）
- 问题文件：test/test_ops.py
- 症状：泛化出的测试用例约 67,693 降至 34,295，丢失的用例集中在 TestCommonCPU 类
- 已知线索：more.json（旧，67k 用例）和 less.json（新，34k 用例）的 diff 显示 TestCommonCPU 测试全部缺失

请按以下步骤排查：
1. 拉取 v2.7.1 最近 10 天合入的全部 PR
2. 按关键词（test_ops、common_device、common_utils、patch、skip）对 PR 评分排序
3. 深入查看高嫌疑 PR 的 diff 变更
4. 交叉验证关键文件（test_upstream/torch/testing/_internal/common_device_type.py.patch）的提交历史
5. 输出排序后的嫌疑 PR 列表及证据
```

## 预期产出

一份包含以下内容的定位报告：

| 定位步骤 | 使用的 gc CLI 命令 | 产出 |
|---------|-------------------|------|
| 拉取 PR 列表 | `gc pr list --state merged --base v2.7.1 --sort updated --direction desc` | 50+ 条合入 PR |
| 获取合入时间 | `curl` + GitCode REST API `/pulls` 端点（gc 暂未支持 merged_at） | 精确合入时间戳 |
| 多维度评分 | 关键词匹配 + 文件变更范围 + 时间窗口 | 嫌疑 PR 排序表 |
| 深入 diff | `gc pr diff` + `curl` `/pulls/{n}/files` 端点 | 核心代码变更证据 |
| 文件历史验证 | `curl` `/commits?path=...` 端点 | 关键文件修改时间线 |

**定位结论**：3 个高嫌疑 PR，其中 #36598（torch inductor patch add）为最可能根因——直接修改了控制设备类型测试类生成的 `common_device_type.py.patch`。

## 真实执行记录

本案例使用 gc 0.5.0 (commit: f60a2bb) 在 Ubuntu WSL2 环境下执行：

### 1. 拉取 PR 列表

```bash
$ gc pr list -R Ascend/pytorch --state merged --base v2.7.1 \
  --sort updated --direction desc --limit 30 --format table

NUMBER  STATE   AUTHOR      TITLE
#33558  merged  xiaoqi-zhou 将test_mish、test_silu等skip掉的用例重新补回来
#36518  merged  cuiduo      [fix]import_all_patch
#36542  merged  AACAES      [test]del test_bidirectional_lstm skip
#34013  merged  kkjocker    fix torch/test/inductor patch apply bug
#36598  merged  kkjocker    【community issue】torch inductor patch add
...
```

### 2. 对 PR 评分排序

使用问题域关键词 `test_ops`, `common_device`, `patch`, `skip` 筛选：

```
| PR      | 评分 | 关键词命中                                |
|---------|------|------------------------------------------|
| #36598  | ***  | torch inductor patch add                  |
| #34013  | ***  | fix torch/test/inductor patch apply bug   |
| #36518  | **   | import_all_patch                          |
| #35921  | **   | 删除跳过用例注解                          |
| #35620  | **   | fix for inductor patches                  |
```

### 3. 深入 #36598 变更

```bash
$ curl -s -H "Authorization: Bearer $GC_TOKEN" \
  "https://gitcode.com/api/v5/repos/Ascend/pytorch/pulls/36598/files"

test_upstream/test/nn/test_pooling.py.patch
test_upstream/torch/_inductor/utils.py.patch
test_upstream/torch/testing/_internal/common_device_type.py.patch  ← 关键文件
```

`common_device_type.py.patch` 中的 `filter_desired_device_types` 函数被大幅修改——该函数控制 PyTorch 测试框架中哪些设备类型的测试类被包含/排除。新增的 `_normalize_device_type` 归一化逻辑将 `npu:0 → npu → privateuse1` 的映射路径复杂化。

```bash
$ gc pr diff --repo Ascend/pytorch 36598
# 确认 diff 中 filter_desired_device_types 从 33 行扩展到 58 行
```

### 4. 交叉验证文件历史

```bash
$ curl -s -H "Authorization: Bearer $GC_TOKEN" \
  "https://gitcode.com/api/v5/repos/Ascend/pytorch/commits?path=test_upstream/torch/testing/_internal/common_device_type.py.patch&sha=v2.7.1"

2026-05-25  【community issue】torch inductor patch add   ← #36598 修改
2026-04-16  feat:Create a test_upstream folder            ← 初始创建
```

文件只在 05-25 被修改一次，与问题出现时间吻合。

### 5. 定位结论

| 排序 | PR | 等级 | 关键证据 |
|------|-----|------|---------|
| 1 | #36598 | 高 | 修改 `common_device_type.py.patch`，`filter_desired_device_types` 函数新增设备类型归一化逻辑 |
| 2 | #34013 | 中 | 修复 96 个 diff 文件的终止空行，影响 patch 拼接边界 |
| 3 | #36518 | 低 | 修改 inductor wrapper，主要影响编译行为而非测试发现 |

## 价值

- **效率提升**：将人工逐条查看 PR 的时间从数小时压缩到 10 分钟内
- **可复用**：适用于任何 GitCode 仓库的 PR 引入问题排查，只需替换仓库名、分支、关键词
- **可审计**：每步都有命令和输出，排查过程可追溯、可复现
- **暴露工具链缺口**：过程中识别出 gc CLI 缺少 `api` 子命令、`pr view` 统计为 0 等问题，已提交 [gitcode-cli/cli#246](https://gitcode.com/gitcode-cli/cli/issues/246)

## 复用方式

```text
我需要定位 <owner/repo> 仓库 <branch> 分支上导致 <症状描述> 的 PR。请使用 gitcode-pr-impact-locator skill。

问题详情：
- 仓库：<owner/repo>
- 分支：<branch>
- 时间范围：<时间窗口>
- 问题文件：<file_path>
- 症状：<具体现象>
- 已知线索：<如有>
```

替换尖括号内容后直接使用。
