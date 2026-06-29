---
name: pipeline-bisect
description: 当 CI 流水线之前正常、最近突然失败时使用。提供系统化方法论：对比正常/异常运行、从制品中提取错误日志、追溯到具体 commit 或 PR。触发词：流水线回归、bisect、之前正常现在失败、对比流水线运行差异。
---

# Pipeline Bisect — 流水线问题快速定位

诊断流水线回归问题，遵循以下 6 步方法论。
核心原则：**不要猜测原因 — 直接对比正常和异常运行。** 两者之间的每一个差异都是线索。

## 前置准备：需要收集的信息

开始之前，先获取：

- **正常运行链接**：最后一次已知正常的流水线执行
- **异常运行链接**：首次（或当前）出现问题的流水线执行
- **症状描述**：哪个文件/步骤失败、观察到的错误信息

## 第 1 步：找 — 获取正常和异常运行的详细信息

确定两次流水线运行，提取任务结构：

```bash
# 列出最近的流水线运行
gh run list --repo OWNER/REPO --branch BRANCH --limit 20 \
  --json databaseId,createdAt,displayTitle,status,conclusion,workflowName

# 获取每次运行的 job 列表
gh run view RUN-ID --repo OWNER/REPO --json jobs \
  --jq '.jobs[] | "\(.name) | \(.conclusion) | id=\(.databaseId)"'
```

记录相关 job 的 ID（如 collect job、test job）。

## 第 2 步：比 — 对比正常运行和异常运行

下载两次运行的完整日志，系统化对比：

| 对比维度 | 检查方式 |
|---------|---------|
| Docker 镜像 | 日志中搜索 `quay.io\|docker image` |
| Action 版本 | 对比 workflow 文件中的 `uses: ...@REF` |
| 输入参数 | 对比 `with:` 代码块 |
| 脚本参数 | 对比 `--verbose` 等 flag 是否出现 |
| 源目录 | 对比脚本输出中的路径 |
| 数量 | 对比 patch/文件/用例数量 |

```bash
# 下载指定 job 的日志
gh run view RUN-ID --repo OWNER/REPO --job JOB-ID --log > /tmp/run.log
```

**关键洞察**：数量差异（如 25 vs 27 个 patch）、源目录不同，往往直接指向根因。

## 第 3 步：下载 — 从制品中提取真实错误

流水线摘要（如"收集到 0 个用例"）不包含完整错误信息。务必下载制品获取真实 traceback：

```bash
# 列出运行的所有制品
gh api repos/OWNER/REPO/actions/runs/RUN-ID/artifacts \
  --jq '.artifacts[] | "\(.name) | id=\(.id)"'

# 下载相关制品
gh run download RUN-ID --repo OWNER/REPO \
  --name ARTIFACT-NAME --dir /tmp/extract

# 解压并查找错误日志
find /tmp/extract -name "*.tar.gz" -exec tar -xzf {} -C /tmp/extract/ \;
find /tmp/extract -name "*PROBLEM-FILE-PATTERN*" -exec cat {} \;
```

**关键原则**：制品中的 traceback 是唯一真相来源。不要依赖流水线摘要，它可能被截断或产生误导。

## 第 4 步：看 — 分类错误并推理假设

根据真实错误信息，分类问题：

| 错误类型 | 常见原因 | 排查方向 |
|---------|---------|---------|
| `ImportError: cannot import name X` | Patch 添加了引用但未添加定义 | 检查 patch 是否修改了调用方，但被调用方缺少对应 patch |
| `IndentationError / SyntaxError` | Patch 应用错误，破坏了目标文件 | 检查 patch 上下文是否与目标文件版本匹配 |
| `ModuleNotFoundError` | 缺少依赖或路径错误 | 检查 pip install 步骤和 PYTHONPATH |
| 构建失败 | Dockerfile 或依赖变更 | 检查 `.ci/docker/` 和 `requirements*.txt` |
| 输出结果突变（数量变化） | 配置或过滤规则变更 | 检查 case_paths 配置和 skip 逻辑 |

## 第 5 步：追 — 用 Git 精确定位引入 commit

用错误信息中的关键词追溯变更：

```bash
# 搜索第一个引入特定字符串的 commit
git log --all --oneline -S "错误中的关键词" -- 相关文件路径

# 查看 commit 的完整变更
git show COMMIT --stat
git show COMMIT -- FILE-PATH

# 检查哪些分支包含该 commit
git branch -a --contains COMMIT
```

**产出**：commit hash、PR 编号、作者、合入时间、关联 issue 链接、变更文件列表。

## 第 6 步：确认 — 分析变更并确认根因

阅读 commit diff，回答以下问题：

1. **改了哪些文件** — 区分直接文件修改和 patch 文件修改
2. **CI 流程是否会采纳这个变更** — 仓库中的直接文件修改能否传递到 CI 执行环境？
3. **缺少了什么** — 是否需要但未包含对应的 patch、引用或配置更新？
4. **为什么正常期不受影响** — action 版本差异、条件触发、还是时间窗口？

**产出**：一句话根因 + 修复方向。

## 常见陷阱

- **不要从外部因素假设入手**（镜像源、换行符、操作系统差异）。先从直接对比两次运行开始。
- **不要只信流水线摘要** — 务必下载制品获取完整 traceback。
- **不要分析整个系统** — 聚焦于正常和异常运行之间的**差异**，而不是理解所有东西。
- **不要在第一个看似合理的解释处止步** — 通过本地复现或检查 git 历史来验证。

## 参考案例

详细实战案例（ONNX 测试收集失败，追溯到 PR !36631 中缺少 `pytorch_test_common.py.patch`），参见：`references/onnx_import_error_case.md`
