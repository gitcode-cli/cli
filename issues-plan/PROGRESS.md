# 项目进度与当前状态

本文档用于区分两件事：

- 历史交付：项目在 2026-03-22 到 2026-03-24 期间完成的首轮功能交付
- 当前状态：围绕真实可用性、一致性和回归能力进行的后续收口

这意味着“历史里程碑已完成”不等于“当前没有后续修复或收口工作”。

**最后更新**: 2026-03-30

---

## 当前状态

项目已经具备可用的认证、仓库、Issue、PR、Release 基础能力，但仍在进行一轮以真实行为校准为目标的收口。

当前收口里程碑：`v0.5.0`（Milestone `#307155`）

| 类型 | 状态 | 说明 |
|------|------|------|
| 基础功能交付 | 已完成 | M1-M7 的首轮功能已交付 |
| 正确性修复 | 持续收口中 | 修复硬编码、认证方式、上下文解析、PR review 兼容等问题 |
| 文档校准 | 持续收口中 | README、COMMANDS、spec、`.claude/skills` 需要跟随真实行为更新 |
| 真实回归 | 已建立最小矩阵 | 已新增 `scripts/regression-core.sh` |

---

## 当前收口清单

| Issue | 状态 | PR | 说明 |
|------|------|----|------|
| `#63` | ✅ 已完成 | `#50` | 统一仓库参数解析逻辑 |
| `#64` | ✅ 已完成 | `#51` | 修复 `pr review --approve`，并明确 `--request` 暂不支持 |
| `#65` | ✅ 已完成 | `#52` | 建立核心真实命令回归矩阵 |
| `#66` | 🚧 当前进行中 | - | 按真实状态重写进度说明 |
| `#67` | ✅ 已完成 | `#53` | 校准 README / COMMANDS / `.claude/skills` 上下文说明 |
| `#68` | 📋 待处理 | - | 清理构建产物、评估输出和本地统计文件 |

---

## 最近完成的关键修复

以下工作不是“新功能扩展”，而是对已有功能做真实行为校准：

| Issue | PR | 结果 |
|------|----|------|
| `#55` | `#43` | `repo fork` 不再硬编码仓库路径 |
| `#56` | `#44` | 清理 `access_token` query 认证，改用 Authorization Header |
| `#59` | `#45` | 修复 `auth login --web`、`pr create --fill`、`pr create --web` 的 silent no-op |
| `#60` / `#61` | `#46` | 统一 auth 配置与环境变量优先级 |
| `#58` | `#47` | 支持从当前 Git 仓库自动识别仓库 |
| `#57` | `#48` | 统一当前分支解析并接入 `pr create` |
| `#63` | `#50` | 仓库参数统一支持 `owner/repo`、HTTPS、SSH |
| `#64` | `#51` | `pr review --approve` 改走真实 endpoint，错误信息可诊断 |
| `#65` | `#52` | 增加可执行的核心真实回归脚本 |
| `#67` | `#53` | 同步 README / COMMANDS / `.claude/skills` 的真实说明 |

---

## 历史交付快照

下面的“已完成”表示对应模块已经完成首轮交付，不表示后续不再有修复。

| 里程碑 | 历史状态 | 首轮完成日期 | 说明 |
|--------|----------|--------------|------|
| M1 基础架构 | ✅ 已交付 | 2026-03-22 | root 命令、Factory、IOStreams、Git 封装等 |
| M2 认证功能 | ✅ 已交付 | 2026-03-22 | login/status/logout/token 等 |
| M3 仓库功能 | ✅ 已交付 | 2026-03-22 | clone/create/fork/view/list/delete |
| M4 Issue 功能 | ✅ 已交付 | 2026-03-22 | create/list/view/close/reopen/comment/label 等 |
| M5 PR 功能 | ✅ 已交付 | 2026-03-22 | create/list/view/checkout/merge/review 等 |
| M6 Release 功能 | ✅ 已交付 | 2026-03-23 | create/list/view/delete/edit/upload/download |
| M7 文档与基础设施 | ✅ 已交付 | 2026-03-23 | README、COMMANDS、SECURITY、打包等 |

---

## 当前已知边界

这些不是“文档遗漏”，而是当前真实行为边界：

- `gc pr review --approve` 当前可用，走 `/pulls/:number/review`
- `gc pr review --request` 当前会明确报错，因为 GitCode 公开 API 暂不支持 request changes
- 核心真实回归已收口到 `./scripts/regression-core.sh`
- `pr create` 已纳入回归矩阵，但默认不执行，需显式开启写路径参数

---

## 真实验证入口

优先使用以下入口验证当前状态：

```bash
go test ./...
go build -o ./gc ./cmd/gc
./scripts/regression-core.sh
```

补充说明：
- `./scripts/regression-core.sh` 默认覆盖 auth、repo view、issue list/view、非 Git 目录错误路径
- 写路径测试需要显式设置环境变量，详见 `docs/REGRESSION.md`

---

## 下一步

当前收口优先级：

1. 完成 `#66`，让状态文档不再误导
2. 完成 `#68`，清理构建产物和评估输出
3. 之后再决定是否进入下一轮 API 兼容或工程治理工作
