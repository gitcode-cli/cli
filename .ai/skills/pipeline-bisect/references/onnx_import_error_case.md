# 案例：ONNX 测试收集 ImportError

## 背景

`_torch-npu-upstream-collect.yml` 流水线中 `test/onnx/test_pytorch_onnx_onnxruntime.py` 收集到 0 个测试用例。May 25 正常，May 26 开始异常。

## 排查过程

### 第 1 步：找到两次运行

| | 正常运行 | 异常运行 |
|---|---|---|
| 链接 | [26388425776](https://github.com/Ascend/pytorch/actions/runs/26388425776) | [26449163884](https://github.com/Ascend/pytorch/actions/runs/26449163884) |
| 时间 | May 25 09:35 UTC | May 26 13:02 UTC |
| Collect job | id=77686187348 | id=77865923773 |

### 第 2 步：对比

| 维度 | May 25（正常） | May 26（异常） |
|------|--------------|---------------|
| Docker 镜像 | `202605250326` | `202605250326`（相同） |
| 脚本参数 | `--python=3.10 --verbose` | `--python=3.10`（无 verbose） |
| Patch 来源 | `ascend_pytorch/test_upstream/torch` | `pytorch-test-src/test_upstream/torch` |
| Patch 数量 | **25** | **27** |

不同的 patch 来源目录和数量差异是突破口。

### 第 3 步：下载制品

```bash
gh run download 26449163884 --repo Ascend/pytorch \
  --name collect-cases-logs --dir /tmp/extract
tar -xzf /tmp/extract/__w/pytorch/pytorch/collection_errors.tar.gz
cat collection_errors/regular/onnx_test_pytorch_onnx_onnxruntime.log
```

真实错误（与最初报告的不同）：

```
ImportError: cannot import name 'skipIfOneDnnVersionLessThan'
  from 'pytorch_test_common'
```

### 第 4 步：分类

ImportError — 上游 PyTorch 的 `pytorch_test_common.py` 中不存在该函数。某个 patch 添加了 import 但没有添加函数定义。

### 第 5 步：Git 追溯

```bash
git log --all --oneline -S "skipIfOneDnnVersionLessThan" \
  -- test_upstream/test/onnx/test_pytorch_onnx_onnxruntime.py.patch
```

输出：
```
581d77253 test(test_pytorch_onnx_onnxruntime.py) fix testcase test_quantized_conv3d/...
```

### 第 6 步：确认根因

**Commit**: `581d77253`
- **PR**: Ascend/pytorch!36631
- **作者**: yuanqi1104
- **合入时间**: May 26 09:48 CST
- **关联 Issue**: https://gitcode.com/Ascend/pytorch/issues/2093

该 commit：
1. 创建了 `test/onnx/pytorch_test_common.py`（CI 仓库中的直接文件，包含 `skipIfOneDnnVersionLessThan` 和 `useBackendOnednnOnArm` 函数定义）
2. 修改了 `test_upstream/test/onnx/test_pytorch_onnx_onnxruntime.py.patch`（添加了这些函数的 import 和使用）

**根因**：直接文件 `test/onnx/pytorch_test_common.py` 在 CI 准备阶段**不会**被复制到 PyTorch 测试源码中（只有 `test_upstream/` 目录下的文件会被复制）。patch 文件引用了这些函数，但没有对应的 patch 将它们添加到上游的 `pytorch_test_common.py` → `ImportError`。

**为什么 May 25 正常**：那次运行使用了不同的 action 版本，从 `ascend_pytorch/test_upstream/` 拉取 patch（25 个），不包含此 PR 的变更。May 26 回退到旧 action，使用 `pytorch-test-src/test_upstream/`（27 个 patch，包含此 PR）。

**修复方向**：创建 `test_upstream/test/onnx/pytorch_test_common.py.patch`，将函数定义添加到上游文件中。
