---
title: 发布平台敏感信息与安全审查
description: 使用 GitCode CLI 对 openLiBing 发布平台执行敏感信息、凭证和常见安全风险检查
---

# 发布平台敏感信息与安全审查

## 场景

发布平台连接 Jenkins、OBS、MongoDB、PostgreSQL、Redis、华为云 SDK 和 Git 平台，配置文件和脚本天然容易出现凭证、内网地址、调试配置和命令执行风险。这个案例用于发布前、合并前或开源前做一次集中安全检查。

## 推荐 skill

- `gitcode-security-check`

## 可直接执行的 Prompt

```text
请使用 gitcode-security-check skill，对 openLiBingNext/openlibing-platform-release 做一次敏感信息与安全审查。

审查范围：
- 分支或 PR：master
- 路径：src/main/resources, Dockerfile, docker-compose.yml, .devcontainer, scripts, start*.sh, docs, src/main/java/com/openlibing/platformrelease/common/utils, src/main/java/com/openlibing/platformrelease/business/service/impl

请全程使用 `gitcode` 命令入口；如需下载代码，使用 SSH。不要输出完整 secret。如果发现真实凭证，请明确建议撤销和轮换。

请重点检查：
- `application*.yaml`、`.env` 风格配置、Docker 配置中是否有真实密码、AK/SK、token；
- Jenkins URL、OBS endpoint、Redis/Mongo/PostgreSQL 配置是否含生产信息；
- `ExecCmdUtil`、`LinuxCommandInjectCheck`、`OpenEulerJenkinsOperation`、`FileFromRepoUtil` 中是否存在命令注入或 SSRF 风险；
- 文档和测试中是否误写真实内部地址或凭证；
- release 资产候选文件是否可能包含本地配置。

请给出按 Critical / High / Medium / Low 分类的审查报告，并说明是否建议阻塞合并或发布。
```

## 预期产出

- 一份围绕发布平台配置、脚本、Jenkins/OBS/Git 相关代码的安全审查报告。
- 对真实凭证、疑似占位符、测试数据、内部地址分别给出判断。
- 明确是否建议阻塞当前发布或合并。

## 价值

- 发布平台属于连接外部系统的中枢服务，提前检查能显著降低凭证泄漏和命令注入风险。
- 让安全审查从“扫一遍关键字”变成“围绕业务连接点审查”。
- 形成可复制的发布前安全门禁。

## 复用方式

复用时替换仓库、分支/PR、配置目录和外部系统连接点即可。对其他 Java 服务，可保留配置、脚本、命令执行、网络调用四类检查。
