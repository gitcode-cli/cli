---
title: 做一次敏感信息与安全审查
description: 使用 GitCode CLI 下载仓库并执行敏感信息、凭证和常见安全风险检查
---

# 做一次敏感信息与安全审查

## 场景

发布前、合并前或开源前，需要检查仓库或 PR 是否包含真实凭证、私钥、敏感配置、危险调用或常见安全风险。

## 推荐 skill

- `gitcode-security-check`

## 可直接执行的 Prompt

```text
请对 <owner/repo> 做一次敏感信息与安全审查。

要求：
1. 全程使用 `gitcode` 命令，不使用 `gc`。
2. 优先使用 `gitcode-security-check` skill；如果未安装该 skill，请按同等流程执行。
3. 如需下载代码，使用 SSH：
   - ssh -T git@gitcode.com
   - gitcode repo clone <owner/repo> --git-protocol ssh
4. 审查范围：
   - 分支或 PR：<branch-or-pr>
   - 路径：<paths>
5. 检查维度：
   - 硬编码 token、password、secret、api key
   - 私钥、证书、`.env`、credentials 文件
   - 日志输出中的敏感信息
   - 不安全 TLS/HTTP 配置
   - shell/SQL/模板注入风险
   - 权限绕过、认证跳过、debug 后门
   - release 资产或文档中的敏感信息
6. 不要在输出中打印完整 secret；只允许展示文件位置、变量名、前后 4 位或哈希摘要。
7. 如果发现真实凭证，明确建议撤销和轮换。

输出：
- 审查范围
- Findings，按 Critical / High / Medium / Low 分类
- 确认的敏感信息
- 疑似误报
- 修复建议
- 是否建议阻塞合并或发布
```

## 预期产出

- 一份安全审查报告。
- 明确的敏感信息处理建议。
- 合并或发布风险判断。

## 价值

- 减少凭证泄漏和不安全配置上线风险。
- 为开源、发布、合并前检查提供标准化证据。

## 复用方式

替换仓库、分支/PR 和路径范围即可应用到不同项目。
