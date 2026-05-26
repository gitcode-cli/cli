---
title: 安装或升级后做全量冒烟验证
description: 使用 GitCode CLI 在 Windows 和 Linux 上验证安装、入口命令、认证、SSH 和关键读命令
---

# 安装或升级后做全量冒烟验证

## 场景

用户安装或升级 GitCode CLI 后，需要确认 Windows 和 Linux 环境下命令入口、认证、SSH、JSON 输出和核心读命令都可用。

## 推荐 skill

- `gitcode-regression`

## 可直接执行的 Prompt

```text
请帮我对当前环境的 GitCode CLI 做一次安装后冒烟验证。

要求：
1. 全程优先使用 `gitcode` 命令。
2. Windows PowerShell 下不要使用裸 `gc`，如需验证别名冲突，可以使用 `gc.exe`。
3. Linux/macOS 下需要同时验证 `gitcode` 和 `gc` 两个入口。
4. 优先使用 `gitcode-regression` skill；如果未安装该 skill，请按同等流程执行。
5. 检查命令：
   - gitcode version
   - gitcode version --json
   - gitcode help --json
   - gitcode schema
   - gitcode schema "pr create"
   - gitcode auth status
6. 检查 SSH：
   - ssh -T git@gitcode.com
7. 使用安全测试仓库做只读验证：
   - gitcode repo view infra-test/gctest1 --json
   - gitcode issue list -R infra-test/gctest1 --state open --json
   - gitcode pr list -R infra-test/gctest1 --state open --json
   - gitcode release list -R infra-test/gctest1 --json
8. 写路径只允许 dry-run：
   - gitcode issue create -R infra-test/gctest1 --title "test: cli regression" --body "temporary test" --dry-run --json
9. 输出每条命令的通过/失败、关键错误和建议处理方式。

输出：
- CLI 版本
- OS 和 shell
- entrypoint 验证结果
- auth 验证结果
- SSH 验证结果
- 只读命令验证结果
- dry-run 写路径验证结果
- 风险和未覆盖项
```

## 预期产出

- 一份安装后冒烟验证报告。
- Windows/Linux 入口命令兼容性结论。
- 认证、SSH 和核心命令可用性结论。

## 价值

- 用户安装后能快速确认环境是否可用。
- 发布负责人可复用为发布后验证脚本的人工版。
- 能及时发现 Windows PowerShell `gc` 别名、Linux 包入口、SSH 权限等问题。

## 复用方式

将测试仓库替换为自己团队允许验证的 `infra-test/*` 仓库即可。
