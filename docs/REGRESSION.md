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

`./scripts/regression-core.sh` 默认执行稳定的读路径、agent-friendly 契约检查和错误路径：

1. `auth login --token`
2. `auth status`
3. `auth token`
4. `auth logout`
5. `auth status` 登出后状态
6. `repo view`
7. `repo view --json`
8. `issue list --limit 1`
9. `issue list --json`
10. `issue view <list 返回的首个 issue>`
11. `issue view --json`
12. `pr list --json`
13. `release list --json`
14. `release delete <tag> --dry-run`
15. 非 Git 目录下的 `repo view` 错误路径

说明：
- 脚本会使用临时 `GC_CONFIG_DIR`，避免污染本地长期配置。
- 登录阶段会把环境变量 token 写入临时配置，然后取消环境变量覆盖，以验证 config-backed auth 流程。
- `release delete` 的 dry-run 默认使用 `GC_REGRESSION_RELEASE_TAG`，未显式设置时默认取 `v0.0.1-test`。

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

如果本次改动涉及 issue 标签写路径、issue 状态流转或 PR 状态流转，再补以下真实命令验证：

```bash
# issue label / close
issue_number=$(./gc issue create -R infra-test/gctest1 --title "Regression" --body "write-path" | sed -n 's/.*#\([0-9][0-9]*\).*/\1/p')
./gc issue label "$issue_number" -R infra-test/gctest1 --add bug
./gc issue label "$issue_number" -R infra-test/gctest1 --remove bug
./gc issue label "$issue_number" -R infra-test/gctest1 --list
./gc issue view "$issue_number" -R infra-test/gctest1 --json
./gc issue close "$issue_number" -R infra-test/gctest1

# pr close
./gc pr create -R infra-test/gctest1 --head <test-branch> --base main --title "Regression PR" --body "write-path"
./gc pr close <pr-number> -R infra-test/gctest1
./gc pr view <pr-number> -R infra-test/gctest1 --json
```

说明：
- 对写路径命令不能只看退出码，必须回读远端状态。
- `issue label --remove` 后，`--list` / `issue view --json` 不应再包含目标标签。
- `pr close` 后，`pr view --json` 的 `state` 应变为 `closed`。

## Agent-Friendly 契约补充回归

除核心脚本已覆盖的检查外，当前里程碑还推荐以下补充检查：

```bash
# 命令元数据
./gc schema
./gc schema "issue view"

# 结构化输出
./gc repo view infra-test/gctest1 --json
./gc issue list -R infra-test/gctest1 --limit 1 --json
./gc issue list -R infra-test/gctest1 --format json
./gc issue list -R infra-test/gctest1 --format table
./gc issue list -R infra-test/gctest1 --format simple
./gc issue list -R infra-test/gctest1 --time-format absolute
./gc issue list -R infra-test/gctest1 --time-format relative
./gc issue list -R infra-test/gctest1 --template '{{range .}}{{.Number}}{{"\n"}}{{end}}'
./gc issue view 1 -R infra-test/gctest1 --json
./gc issue view 1 -R infra-test/gctest1
./gc pr list -R infra-test/gctest1 --json
./gc pr view 1 -R infra-test/gctest1 --json
./gc pr view 1 -R infra-test/gctest1
./gc release list -R infra-test/gctest1 --json

# dry-run
./gc repo delete infra-test/gctest1 --dry-run
./gc label delete bug -R infra-test/gctest1 --dry-run
./gc milestone delete 1 -R infra-test/gctest1 --dry-run
./gc release delete v0.0.0 -R infra-test/gctest1 --dry-run
./gc issue create -R infra-test/gctest1 --title "Regression" --body "dry-run" --dry-run
```

说明：
- 上述 dry-run 命令不应执行真实写入。
- 非交互环境中未传 `--yes` 的删除命令应直接失败，不应挂起等待输入。
- `issue list --format yaml` 应返回用法错误，不应静默回退到默认输出。
- `issue list --json` 与 `issue list --format json` 应保持等价。
- `issue list --time-format absolute|relative` 仅影响文本展示，不应改变 JSON 输出。
- `issue list --template` 应输出模板结果，并与 `--json`、`--format` 的冲突保持稳定报错。
- `issue view` 和 `pr view` 的文本输出应保留稳定的详情布局，`--json` 继续输出结构化数据。

## 推荐记录方式

在 PR 描述或 issue comment 中记录：

```markdown
## Regression

- [x] ./scripts/regression-core.sh
- [x] read-only checks passed
- [ ] write-path checks skipped
```
