# 设计:`gc precommit check`(提交前 pre-commit 配置与环境校验)

- 关联 Issue: gitcode-cli/cli #258(后续跟踪 #260)
- 日期: 2026-06-03(后续跟踪项更新于 2026-06-04)
- 状态: 已实现(PR #220);#260 跟踪 JSON `reason` 字段 / 安装失败分类 / 版本探测改用 stdout-only

## 1. 背景与目标

使用 gitcode-cli 提交代码时,工具不会感知项目是否配置了 pre-commit。
当项目存在 pre-commit 配置但本地环境未安装或未初始化时,提交会绕过本应触发的
代码检查,导致未经检查的代码进入仓库。

目标:提供一个命令,在提交代码前检测项目是否配置了 pre-commit、本地工具是否安装、
git hook 是否已初始化,必要时自动安装并初始化,并可选地实际拉起 pre-commit 检查。
需兼容 Windows、Linux(x86_64 / arm)、macOS。

### 关键前提(实现可行性)

`gitcode-cli` **没有本地提交命令**:`gc commit` 仅用于查看/对比**远端** GitCode 提交
(见 `pkg/cmd/commit/commit.go`)。本地 git 仅被 `pr create`、`repo clone/checkout/sync`
间接调用,薄封装在 `git/git.go`。因此本特性以**独立校验命令**落地,而非挂接到某个提交流程。

## 2. 命令与 flag

```
gc precommit check [flags]
```

| Flag | 默认 | 作用 |
|---|---|---|
| `--run` | false | 校验通过后实际执行 `pre-commit run --all-files` |
| `--no-install` | false | 仅诊断,绝不修改环境(默认即"允许在交互式终端自动安装/初始化",故不单设 `--install`) |
| `--yes, -y` | false | 非交互环境下确认执行会修改环境的自动安装/初始化(与 `--no-install` 互斥) |
| `--json` | false | 结构化输出到 stdout |

### 安全约束(对齐 `spec/foundations/agent-friendly-cli.md` §3/§4)

自动安装 pre-commit、执行 `pre-commit install` 都会**修改用户环境**,属于 mutating 操作:

- TTY 环境:默认提示并执行自动安装/初始化。
- 非 TTY 且未给 `--yes`:**不静默修改环境**,立即报错并打印可操作的修复命令。
- `--no-install`:仅诊断,绝不修改环境。

这样既满足"尽量自动装+初始化"的诉求,又不违反"非交互不可隐式改环境"。

## 3. 检查流水线(check 状态机)

1. **定位仓库**:用 `git.RootDir()`。非 git 仓库 → 退出码 1,提示。
2. **检测配置**:仓库根是否存在 `.pre-commit-config.yaml` 或 `.pre-commit-config.yml`。
   - 无配置 → 输出"未配置 pre-commit,跳过",**退出码 0**(无配置不是错误)。
3. **检测工具**:`pre-commit --version` 是否可执行。
   - 缺失且允许安装 → 进入 §4 安装;否则提示 + 退出码 1。
4. **检测初始化**:仓库 `.git/hooks/pre-commit` 是否存在且由 pre-commit 安装。
   - 未初始化且允许 → 执行 `pre-commit install`;否则提示 + 退出码 1。
5. **可选拉起**:`--run` 时执行 `pre-commit run --all-files`。失败时**不透传 pre-commit 原始退出码**,而是按统一错误模型以退出码 `1` 结束,并把检查输出经 `RunOutput` 反馈到 stderr。

## 4. 跨平台自动安装策略

pre-commit 是 Python 包,按工具**可用性择优**,不写死 OS 路径,天然适配 arm/x86 与三大系统:

1. 已在 PATH(`pre-commit --version` 成功)→ 跳过安装。
2. 按顺序尝试可用的安装器,**任一失败则继续尝试下一个**(不因 pipx 失败而中止):
   `pipx install pre-commit` → `python3 -m pip install --user pre-commit` →
   `python -m pip install --user pre-commit`;Windows 额外追加 `py -m pip install --user pre-commit`(Windows Python launcher)。
3. 每次安装后重新探测 `pre-commit --version` 确认成功;成功即返回,否则记录该次失败原因并继续。
4. 全部不可用/全部失败 → **不自动安装 Python**(过于侵入),报错并给出各平台安装指引(brew / apt / 官网)及各次失败详情。
5. 安装失败时按错误输出分类(权限不足 / 网络失败 / 工具链缺失),在聚合错误中追加针对性修复指引;多次同类失败的指引去重,顺序稳定。分类基于 installer(pip/pipx)输出关键字,经 `CommandRunner` 注入,保持可测试。

## 5. 代码结构

### 核心逻辑 `pkg/precommit/`(纯逻辑,可单测)

- `runner.go` — 定义注入式 `CommandRunner` 接口(包一层 `os/exec`),测试用 fake 替换,
  使单测**不依赖真实 pre-commit / python**。接口提供 `Run`(合并 stdout+stderr,用于安装/run 的完整诊断)与 `RunStdout`(仅 stdout,用于版本探测,避免 stderr warning 干扰解析)两种执行方式。
- `detect.go` — 查找配置文件、git 根、`pre-commit --version`、hook 初始化状态。
- `install.go` — 探测 pipx/pip、执行安装与 `pre-commit install`。
- `check.go` — 编排上述流水线,返回结构化结果 `Result`。

### 命令层

- `pkg/cmd/precommit/precommit.go` — 命令组。
- `pkg/cmd/precommit/check/check.go` — 子命令,遵循 Factory + IOStreams + `runF` 注入模式
  (对齐 `pkg/cmd/auth/status/status.go`)。
- 在 `pkg/cmd/root/root.go` 注册 `precommit`。

## 6. 输出与退出码

- 文本:逐步骤 ✓/✗ 与修复建议;警告/提示写 stderr,主数据写 stdout。
- `--json`(仅 stdout):

```json
{
  "config_found": true,
  "tool_installed": true,
  "tool_version": "3.7.0",
  "hook_installed": true,
  "actions_taken": ["installed pre-commit via pipx", "ran pre-commit install"],
  "run_result": "passed",
  "run_output": "",
  "ok": true,
  "reason": ""
}
```

> `run_output` 仅在 `run_result == "failed"` 时携带 `pre-commit run` 的输出,便于定位失败原因;成功或未 `--run` 时省略。

> `reason` 是稳定、机器可读的结果分类,供脚本/agent 直接分支,不必解析文本。取值:
> `no_config`(无配置跳过,`ok=true`)、`tool_missing`、`hook_missing`、`run_failed`、`install_failed`、`not_in_repo`(命令层在非 git 仓库时设置)。环境完全就绪时省略(为空)。

> `install_failed`:已授权自动安装但未能产出可用工具(安装尝试失败,或无可用安装器;退出码 `1`)。即使在此硬失败路径,`--json` 仍输出结构化结果体而非仅退出码 + stderr 文本——`Check` 在返回 error 的同时回填 `reason`/`install_failure_categories`,命令层在 `--json` 下先写出该结果体再返回 CLIError(人类可读的失败详情仍走 stderr)。`install_failure_categories` 为机器可读的失败类型数组(`permission` / `network` / `toolchain`,按首见顺序去重,无法归类的失败不计入,可能为空),与人类可读的安装失败指引(权限/网络/工具链)同源于 `classifyInstallFailure`。

- 退出码(对齐 agent-friendly-cli §6):
  - `0` 就绪(或无配置跳过)
  - `1` 未就绪 / 检查失败 / 非 git 仓库 / 非交互拒绝修改环境
  - `2` 参数或用法错误

## 7. 测试策略(TDD)

- `pkg/precommit/*_test.go`:用 fake `CommandRunner` + 临时目录覆盖分支:
  无配置 / 有配置无工具 / 工具在但未初始化 / 全就绪 / `--run` 通过 / `--run` 失败 /
  非交互拒绝安装 / 安装失败。
- `pkg/cmd/precommit/check/check_test.go`:用 `runF` + test IOStreams 断言文本、JSON、退出码。
- 单测**不真跑** pre-commit;真实命令验证留给 `infra-test/*` 手动回归。

## 8. 文档同步

- `docs/COMMANDS.md`:新增 `precommit check` 说明。
- `.ai/skills/gitcode-cli/` 与 `.claude/skills/gitcode-cli/` 命令速查表。
- 如影响 AI 使用路径:`docs/AI-GUIDE.md`、`AGENTS.md`。
- 如纳入回归矩阵:`docs/REGRESSION.md`。

## 9. 已确认的设计决策

- 落地形式:独立命令 `gc precommit check`。
- 修复策略:尽量自动安装 + 初始化,但非 TTY 需 `--yes`,`--no-install` 可纯诊断。
- 拉起:默认仅校验,`--run` 时执行 `pre-commit run --all-files`。
- 检测范围:仅 pre-commit 框架(`.pre-commit-config.yaml/.yml`),不含 husky/lefthook。

## 10. 非目标(YAGNI)

- 不支持 husky、lefthook 等其他 hook 框架。
- 不自动安装 Python 运行时。
- 不实现本地 `git commit` 封装命令。
