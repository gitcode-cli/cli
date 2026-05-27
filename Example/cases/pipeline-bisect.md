---
title: 流水线问题定位 — 对比正常/异常运行追查引入变更
description: 使用 pipeline-bisect skill 系统化对比正常/异常 CI 运行，从制品日志中获取真实错误，用 git log -S 精确锁定引入问题的 commit 和 PR
---

# 案例：流水线问题定位 — 对比正常/异常运行追查引入变更

## 场景

**流水线类问题**：某条 CI 流水线前几天运行正常，最新一次执行出现异常（如测试用例收集失败、构建报错等），需要快速定位是哪次代码变更引入的问题。

这类问题的核心特征：
- 有时间上的"正常期"和"异常期"分界
- 通常在 CI 制品中有完整日志
- 根因往往是某个 PR/commit 的合入

---

## 通用定位流程

```
                    ┌──────────────────────────────────────┐
                    │  1. 确定正常运行和异常运行的流水线任务  │
                    │     - 正常运行：最后一次已知正常的执行   │
                    │     - 异常运行：首次出现问题的执行       │
                    └──────────────────┬───────────────────┘
                                       │
                    ┌──────────────────▼───────────────────┐
                    │  2. 对比两个运行的差异点               │
                    │     - 输入参数、Docker 镜像、action 版本 │
                    │     - 关键步骤的执行路径和输出           │
                    │     - 制品中的详细日志                  │
                    └──────────────────┬───────────────────┘
                                       │
                    ┌──────────────────▼───────────────────┐
                    │  3. 从制品中获取第一手错误信息          │
                    │     - 收集 error logs / traceback     │
                    │     - 不要仅看流水线摘要（可能不完整）    │
                    └──────────────────┬───────────────────┘
                                       │
                    ┌──────────────────▼───────────────────┐
                    │  4. 根据错误类型定位引入变更            │
                    │     - ImportError → patch/依赖变更     │
                    │     - 构建失败 → Dockerfile/依赖变更    │
                    │     - 结果突变 → 逻辑/配置变更          │
                    └──────────────────┬───────────────────┘
                                       │
                    ┌──────────────────▼───────────────────┐
                    │  5. 用 git 工具精确锁定 commit         │
                    │     - git log -S "关键字" 搜索引入变更  │
                    │     - git log -- file 查看文件历史     │
                    │     - 输出：PR 号、作者、时间、关联 issue│
                    └──────────────────┬───────────────────┘
                                       │
                    ┌──────────────────▼───────────────────┐
                    │  6. 分析 commit 变更确认根因            │
                    │     - 变更了什么文件                    │
                    │     - 为什么会导致这个问题               │
                    │     - 确认修复方向                      │
                    └──────────────────────────────────────┘
```

---

## 实例：ONNX 测试用例收集失败

### 问题描述

`_torch-npu-upstream-collect.yml` 流水线中 `test/onnx/test_pytorch_onnx_onnxruntime.py` 收集到 0 个测试用例。前几天（May 25）正常，May 26 开始异常。

### Step 1：确定正常和异常的流水线任务

| | 正常运行 | 异常运行 |
|---|---|---|
| 链接 | [26388425776](https://github.com/Ascend/pytorch/actions/runs/26388425776) | [26449163884](https://github.com/Ascend/pytorch/actions/runs/26449163884) |
| 时间 | May 25 09:35 UTC | May 26 13:02 UTC |
| 结果 | ONNX 测试收集正常 | ONNX 测试 0 cases |

### Step 2：对比两个运行的差异

下载两个运行的 collect job 日志，对比关键步骤：

| 维度 | May 25（正常） | May 26（异常） |
|------|---------------|---------------|
| Docker 镜像 | `202605250326` | `202605250326`（相同） |
| `torch_env_patch.sh` 参数 | `--python=3.10 --verbose` | `--python=3.10`（无 verbose） |
| Patch 来源目录 | `ascend_pytorch/test_upstream/torch` | `pytorch-test-src/test_upstream/torch` |
| Patch 数量 | **25** | **27**（多了 2 个） |

**差异即线索**：异常运行使用了不同来源的 patch 目录（kerer-ai 仓库 vs Ascend 仓库），且多出 2 个 patch 文件。问题很可能出在这 2 个新增的 patch 上。

### Step 3：从制品中获取真实错误

下载异常运行的 `collect-cases-logs` 制品，解压 `collection_errors.tar.gz`：

```bash
gh run download 26449163884 --repo Ascend/pytorch \
  --name collect-cases-logs --dir /tmp/extract
tar -xzf /tmp/extract/__w/pytorch/pytorch/collection_errors.tar.gz
cat collection_errors/regular/onnx_test_pytorch_onnx_onnxruntime.log
```

**真实错误**（注意：与用户最初报告的错误不同，制品日志才是真相来源）：

```
ImportError: cannot import name 'skipIfOneDnnVersionLessThan'
  from 'pytorch_test_common'
```

### Step 4：根据错误类型定位

`ImportError` 意味着某个模块缺少函数定义。错误链路：

```
test/onnx/test_pytorch_onnx_onnxruntime.py (已被 patch 修改)
  → from pytorch_test_common import ... skipIfOneDnnVersionLessThan
    → pytorch_test_common.py (上游 PyTorch 原文件，未被修改)
      → 函数不存在 → ImportError
```

推断：某个 patch 修改了 `test_pytorch_onnx_onnxruntime.py` 添加了对 `skipIfOneDnnVersionLessThan` 的引用，但没有对应的 patch 在 `pytorch_test_common.py` 中添加函数定义。

### Step 5：git 精确定位引入 commit

```bash
# 搜索哪个 commit 引入了 skipIfOneDnnVersionLessThan 到 patch 文件中
git log --all --oneline -S "skipIfOneDnnVersionLessThan" \
  -- test_upstream/test/onnx/test_pytorch_onnx_onnxruntime.py.patch
```

输出：

```
581d77253 test(test_pytorch_onnx_onnxruntime.py) fix testcase test_quantized_conv3d/...
```

查看该 commit 详情：

```bash
git show 581d77253 --stat
# test/onnx/pytorch_test_common.py                   | 93 ++++++
# test_upstream/test/onnx/test_pytorch_onnx_onnxruntime.py.patch | 35 +++-
```

### Step 6：确认变更内容，定位根因

**Commit**: `581d77253`
- **PR**: Ascend/pytorch!36631
- **作者**: yuanqi1104
- **合入时间**: May 26 09:48 CST (May 26 01:48 UTC)
- **关联 issue**: https://gitcode.com/Ascend/pytorch/issues/2093

**变更分析**：

```
PR !36631 (commit 581d77253)
├── 新增: test/onnx/pytorch_test_common.py
│         └── 在 CI 仓库中直接添加了 skipIfOneDnnVersionLessThan 等函数定义
│         └── 但这个文件不会被 CI 流程复制到 pytest 运行的测试源码中
│
└── 修改: test_upstream/test/onnx/test_pytorch_onnx_onnxruntime.py.patch
          └── 在 import 中添加了 skipIfOneDnnVersionLessThan
          └── 在装饰器中使用了 @skipIfOneDnnVersionLessThan(3, 9, 0)
```

**根因**：CI 仓库中的直接文件修改（`test/onnx/pytorch_test_common.py`）不会自动传播到 PyTorch 上游测试源码。需要通过 `test_upstream/test/onnx/pytorch_test_common.py.patch` 以 patch 形式来修改上游文件。该 PR 仅创建了直接文件，缺少对应的 patch 文件，导致 `ImportError`。

---

![image.png](https://raw.gitcode.com/user-images/assets/9483585/5bdf5a1d-ff39-4af2-bcf7-bf63fdba2d55/image.png 'image.png')

## 经验总结

### 正确方向 vs 错误方向

| 错误方向（浪费时间） | 正确方向（快速定位） |
|---|---|
| 猜测 PyTorch 安装源有问题，对比华为/官方镜像 | 直接对比正常/异常运行的 CI 日志 |
| 猜测 CRLF 换行符导致 patch 异常 | 从制品中获取真实 traceback |
| 逐行审查 Docker 构建日志 | 用 `git log -S` 搜索引入变更的 commit |

### 流水线问题的通用排查口诀

```
一找二比三下载，四看五追六确认

① 找到正常和异常两个流水线任务链接
② 对比两者的参数、镜像、action 版本、关键步骤输出
③ 从异常运行的制品中下载详细错误日志（不要只看摘要）
④ 根据错误类型（ImportError/SyntaxError/构建失败）推理可能原因
⑤ 用 git log -S / git bisect 追溯引入问题的 commit
⑥ 分析 commit 变更内容，确认根因和修复方向
```

### 关键教训

- **制品中的错误日志是唯一的真相来源** — pipeline 摘要只说 "0 cases"，必须解压 `collection_errors.tar.gz` 才能看到真实 traceback
- **正常/异常对比是最高效的入口** — 不需要理解整个系统，只要找到差异点就能锁定方向
- **action 版本漂移会被对比暴露** — 正常和异常运行的 action 版本、patch 数量差异是本次定位的转折点
- **git log -S "关键字" 是定位利器** — 知道错误中出现的特定字符串后，一步就能找到引入它的 commit

---

## 推荐 skill

- **pipeline-bisect**: 流水线问题定位 skill，6 步定位流程 + 错误分类参考表 + 常见陷阱
- **Skill 仓库**: `git@gitcode.com:TTFHVV/skills.git`
- **触发方式**: 当描述中包含"流水线前几天正常现在失败"、"pipeline regression"、"定位引入问题的PR"、"对比正常和异常的流水线"等关键词时自动加载

## 关联 Skill

本案例已沉淀为一个可复用的 **pipeline-bisect** skill，详见 [TTFHVV/skills](https://gitcode.com/TTFHVV/skills) 仓库。
