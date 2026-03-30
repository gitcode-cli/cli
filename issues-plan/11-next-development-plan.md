# 下一步开发计划

本文档用于承接 2026-03-30 之后的“收口与一致性”开发阶段。

当前项目已经完成一轮高覆盖度命令实现，但经过真实开发与回归后，仍有几类问题需要作为下一阶段主线处理：

1. 上下文解析还没有彻底收口，命令层仍有重复实现
2. GitCode API 兼容层仍有行为不稳定点，尤其是 PR review
3. 关键真实命令缺少成体系的回归矩阵
4. 文档和状态文件仍有历史失真
5. 主工作区存在构建产物和评估文件，工程边界不清

---

## 1. 本阶段目标

本阶段不再横向扩更多命令，而是优先完成以下工作：

- 统一 `owner/repo` 解析入口，消除命令层散落的 `parseRepo`
- 修复 `pr review` 在 GitCode 平台上的接口兼容问题
- 为核心命令建立最小但稳定的真实回归矩阵
- 刷新 `README.md`、`docs/COMMANDS.md`、`issues-plan/PROGRESS.md`
- 清理构建产物、评估输出和本地统计文件的仓库边界

---

## 2. 建议里程碑

建议新建里程碑：

- 名称：`v0.5.0`
- 主题：收口与一致性

目标定义：

- 核心命令行为一致
- 文档和真实行为一致
- 回归用例可重复执行
- PR / Issue / Repo 基础链路无明显平台兼容问题

---

## 3. 建议 Issue 拆分

### P0

1. `refactor(context): 统一仓库参数解析逻辑`
   范围：将命令层重复的 `parseRepo` 收口到统一入口，统一支持 `owner/repo`、HTTPS、SSH 三种格式。

2. `fix(pr): 修复 pr review 在 GitCode 上的 approve/request changes 兼容性`
   范围：核验真实接口，修复或降级 `pr review --approve` 与 `--request`。

3. `test(regression): 建立核心命令真实回归矩阵`
   范围：覆盖 `auth`、`repo view`、`issue list/view`、`pr create`、非 git 目录错误路径。

### P1

4. `docs(progress): 按当前真实状态重写交付进度说明`
   范围：重写 `issues-plan/PROGRESS.md`，将历史交付与当前待收口问题分开描述。

5. `docs(commands): 校准 README 与 COMMANDS 的上下文解析说明`
   范围：明确哪些命令支持自动推断，哪些必须显式传参。

6. `chore(repo): 清理构建产物、评估输出和本地统计文件`
   范围：治理 `.gitignore`、工作区噪音文件和评估目录边界。

---

## 4. 执行顺序

建议按以下顺序推进：

1. 统一仓库参数解析
2. 修复 PR review 兼容性
3. 建立核心真实回归矩阵
4. 刷新状态和命令文档
5. 清理工程边界与产物

---

## 5. 完成定义

满足以下条件后，可认为本轮收口完成：

- 命令层不再保留多套仓库解析实现
- `pr review --approve` / `--request` 有明确且可验证的行为
- 关键命令真实回归可重复通过
- `README.md`、`docs/COMMANDS.md`、`issues-plan/PROGRESS.md` 与代码现状一致
- 主工作区不再被构建产物和评估输出长期污染

---

**最后更新**: 2026-03-30
