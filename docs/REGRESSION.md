# 核心回归矩阵

本文档定义当前里程碑的最小真实命令回归集，对应脚本：

```bash
./scripts/regression-core.sh
```

## 目标

- 覆盖核心用户路径
- 覆盖关键错误路径
- 让实际命令验证可重复执行，而不是每个 issue 临时补一组命令

## 前置条件

```bash
export GC_TOKEN=your_token
go build -o ./gc ./cmd/gc
```

测试仓库限制：
- 只使用 `infra-test/gctest1` 或 `infra-test` 组织下仓库
- 不使用 `gitcode-cli/cli`
- 不使用个人仓库

## 默认回归集

`./scripts/regression-core.sh` 默认执行稳定的读路径和错误路径：

1. `auth login --token`
2. `auth status`
3. `auth token`
4. `auth logout`
5. `auth status` 登出后状态
6. `repo view`
7. `issue list --limit 1`
8. `issue view <list 返回的首个 issue>`
9. 非 Git 目录下的 `repo view` 错误路径

说明：
- 脚本会使用临时 `GC_CONFIG_DIR`，避免污染本地长期配置。
- 登录阶段会把环境变量 token 写入临时配置，然后取消环境变量覆盖，以验证 config-backed auth 流程。

## 可选写路径

`pr create` 属于有副作用的写路径，不默认执行。显式开启方式：

```bash
export GC_REGRESSION_WRITE=1
export GC_REGRESSION_PR_REPO=infra-test/some-test-repo
export GC_REGRESSION_PR_HEAD=feature-branch
export GC_REGRESSION_PR_BASE=main
./scripts/regression-core.sh
```

说明：
- 只有在明确提供目标仓库和 head 分支时才执行 `pr create`
- 这样可以避免默认污染测试仓库，同时仍把 `pr create` 纳入统一回归矩阵

## 推荐记录方式

在 PR 描述或 issue comment 中记录：

```markdown
## Regression

- [x] ./scripts/regression-core.sh
- [x] read-only checks passed
- [ ] write-path checks skipped
```
