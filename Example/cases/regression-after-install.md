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
请使用 gitcode-regression skill，对当前环境的 GitCode CLI 做一次安装后冒烟验证。

测试仓库：infra-test/gctest1

请优先使用 `gitcode` 命令入口。Windows PowerShell 不要使用裸 `gc`；Linux/macOS 需要验证 `gitcode` 和 `gc` 两个入口。写操作只允许 dry-run。

请输出验证报告，包括 CLI 版本、OS/shell、entrypoint、auth、SSH、只读命令、dry-run 写路径、风险和未覆盖项。
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
