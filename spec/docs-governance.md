# 文档治理规范

本文件定义 gitcode-cli 仓库的文档体系、AI 协作文档分层、唯一真相源和分阶段实施方案。

## 职责

定义文档分层、唯一真相源、AI 协作入口边界和变更后的同步规则。

## 适用场景

- 修改命令行为说明
- 修改 AI 入口文档或 skills
- 调整项目规范、README、docs、issues-plan 之间的边界

## 必须

- 以 `spec/` 作为项目规则唯一来源
- 明确每类信息的唯一真相源
- 行为变化后检查相关入口文档和 skills

## 禁止

- 只改代码不检查文档同步
- 在 AI 文档中定义与 `spec/` 冲突的规则
- 把 `.claude/skills/` 当成跨 AI 的唯一 skill 来源

## 同步要求

- 命令、流程、发布、审查、AI 协作变化后都要检查对应入口文档和 skills
- 状态变化后同步 `issues-plan/PROGRESS.md`

## 不负责什么

- 命令实现细节
- 代码风格
- PR blocker 判定

## 1. 目标

本仓库的文档治理目标是：

- 为用户、开发者、维护者和 AI 协作者提供清晰入口
- 降低 README、docs、spec、skills 之间的重复
- 明确每类信息的唯一真相源
- 明确 Claude 与 Codex 并存时的协作边界
- 让文档同步成为正式交付要求，而不是个人习惯

## 2. 文档分层

### 2.1 README.md

`README.md` 是项目总入口，负责：

- 项目简介
- 安装与快速开始
- 文档导航
- 开发入口
- AI 协作入口

`README.md` 不应承载完整命令手册、完整开发流程或历史规划细节。

### 2.2 docs/

`docs/` 是用户文档层，负责：

- 命令行为说明
- 用户操作说明
- 安装、打包、发布的使用说明

其中 `docs/COMMANDS.md` 是命令行为唯一真相源。

### 2.3 spec/

`spec/` 是项目正式规范层，负责：

- 开发流程
- 编码规范
- 测试规范
- 安全规范
- 构建与打包规范
- 发布流程规范
- 代码质量门禁
- 文档治理规范

`spec/` 是项目规则唯一来源。其他文档不得定义与 `spec/` 冲突的规则。

### 2.4 issues-plan/

`issues-plan/` 是规划与状态层，负责：

- 阶段计划
- 历史规划
- 当前进度状态

其中 `issues-plan/PROGRESS.md` 是当前状态唯一真相源。

### 2.5 AI 协作入口层

AI 协作入口层由以下文件组成：

- `AGENTS.md`
- `CLAUDE.md`

这两个文件负责告诉不同 AI 客户端如何进入本仓库的规范体系，不负责重新定义项目规则。

### 2.6 AI skill 层

AI skill 采用三层结构：

- `.ai/skills/`：跨 AI 的共享 skill 真相源
- `.ai/distribution/`：可分发的通用 skill 包
- `.claude/skills/`：Claude 适配层
- `.codex/skills/`：Codex 适配层

当前仓库已补齐 `.ai/skills/`、`.claude/skills/` 和 `.codex/skills/` 的分层结构。
同时，`.ai/distribution/` 用于承载可脱离本仓库复用的 `gc` 通用 skill。

## 3. 唯一真相源

本仓库的唯一真相源定义如下：

- 项目总入口：`README.md`
- 命令行为：`docs/COMMANDS.md`
- 项目正式规范：`spec/`
- 当前状态：`issues-plan/PROGRESS.md`
- AI 共享 skill 真相源：`.ai/skills/`
- Codex 项目级入口：`AGENTS.md`
- Claude 项目级入口：`CLAUDE.md`

### 3.1 权威关系

权威关系按以下顺序生效：

1. `spec/` 定义项目正式规则
2. `docs/COMMANDS.md` 定义命令行为
3. `AGENTS.md` / `CLAUDE.md` 定义 AI 入口
4. `.ai/skills/` 定义共享场景技能
5. `.claude/skills/` / `.codex/skills/` 定义客户端适配

`AGENTS.md`、`CLAUDE.md`、`.claude/skills/`、`.codex/skills/` 不得定义与 `spec/` 冲突的项目规则。

## 4. 多 AI 协作规则

本仓库允许多人协作，并允许同时使用 Codex 和 Claude。

在此模式下，必须遵守以下规则：

- 任何开发流程、测试要求、安全要求都以 `spec/` 为准
- `AGENTS.md` 负责 Codex 侧入口说明
- `CLAUDE.md` 负责 Claude 侧入口说明
- `.claude/skills/` 不是跨 AI 的唯一 skill 来源
- `.codex/skills/` 必须与共享真相源保持一致

AI 协作文档属于正式项目文档，行为变更后必须纳入同步检查。

## 5. skill 体系设计

skill 分为两类：

- 项目专属 skill：仅适用于 gitcode-cli 仓库
- `gc` 通用 skill：可用于其他使用 `gc` 的项目

后续治理方向是：

- `.ai/skills/` 存放共享 skill 真相源
- `.claude/skills/` 存放 Claude 适配版本
- `.codex/skills/` 存放 Codex 适配版本

用户本地安装目录仅作为运行时副本，不作为项目规则真相源。

此外，需要明确区分两类 skill：

- 仓库协作 skill：服务 gitcode-cli 仓库自身开发，可依赖仓库内 `spec/`、`docs/` 和目录结构
- 通用 `gc` skill：服务外部项目使用 `gc`，不得依赖本仓库私有文档路径

当前仓库已提供最小同步工具：

- `scripts/sync-ai-skills.sh`

当前脚本边界：

- 会基于共享源更新 `.codex/skills/*` 的基础适配文件
- 只会为缺失的 `.claude/skills/*` 生成占位入口
- 不会覆盖现有 Claude skill 正文

## 6. 文档同步规则

### 6.1 命令行为变化

命令行为变化时，必须检查并同步：

- `README.md`
- `docs/COMMANDS.md`
- `AGENTS.md`
- `CLAUDE.md`
- 相关 AI skills

### 6.2 开发流程变化

开发流程变化时，必须检查并同步：

- `spec/*`
- `AGENTS.md`
- `CLAUDE.md`
- 相关 AI skills

### 6.3 审查流程变化

审查流程变化时，必须检查并同步：

- `spec/workflows/review-workflow.md`
- `AGENTS.md`
- `CLAUDE.md`
- 审查相关 skills

### 6.4 构建、打包、发布规则变化

构建、打包、发布规则变化时，必须检查并同步：

- `spec/build-and-package.md`
- `spec/release-process.md`
- `AGENTS.md`
- `CLAUDE.md`
- 相关 AI skills

### 6.5 状态变化

阶段计划或当前状态变化时，必须检查并同步：

- `issues-plan/PROGRESS.md`
- `issues-plan/README.md`

## 7. 规范补齐计划

当前 `spec/` 目录已具备开发、测试、安全、build、release、quality 和 workflow 基础规范，但仍需补齐以下文档：

- `spec/ci-workflows.md`

其中 `spec/ci-workflows.md` 放在最后阶段落地，因为当前 GitCode CI 条件尚未具备。

## 8. 分阶段实施方案

### 阶段 1：治理基线

目标：先建立规则和边界。

交付物：

- `spec/docs-governance.md`
- `spec/build-and-package.md`
- `spec/release-process.md`
- `spec/code-quality-gates.md`
- 更新 `spec/README.md`

验收标准：

- 明确各目录职责
- 明确唯一真相源
- 明确文档同步规则
- 明确后续规范补齐顺序

### 阶段 2：入口收口

目标：让用户、开发者、Codex、Claude 都知道先看哪里。

交付物：

- 调整 `README.md`
- 新增 `docs/README.md`
- 调整 `AGENTS.md`
- 调整 `CLAUDE.md`

验收标准：

- `README.md` 成为总入口
- `docs/` 和 `spec/` 的边界清晰
- AI 协作入口清晰

### 阶段 3：AI skill 体系重构

目标：建立共享真相源和客户端适配层。

交付物：

- 新增 `.ai/README.md`
- 新增 `.ai/skills/`
- 新增 `.codex/skills/`
- 梳理 `.claude/skills/`

验收标准：

- skill 真相源明确
- Claude 和 Codex 都有项目级适配层
- skill 体系支持跨项目复用

### 阶段 4：内容去重与迁移

目标：减少重复和漂移。

交付物：

- 压缩 `README.md`
- 收口 `docs/COMMANDS.md`
- 清理 `spec/` 与 AI 文档的重复
- 明确 `issues-plan/` 的历史与当前边界

验收标准：

- README 不再承载重复命令细节
- 命令行为只在 `docs/COMMANDS.md` 定义
- AI 文档不再复制通用规范

### 阶段 5：分发与同步工具

目标：让团队成员稳定复用同一套 skill 资产。

交付物：

- `scripts/sync-ai-skills.sh`
- 本地安装与同步说明

验收标准：

- 共享 skill 可同步到客户端适配层
- 多人环境下可复用、可版本化、可审查

### 阶段 6：CI 自动化

目标：在 GitCode CI 环境具备后再落地自动化约束。

交付物：

- `spec/ci-workflows.md`
- workflow 模板
- 质量门禁的 CI 映射

## 下一步去看哪里

- 如果你在修改命令行为，继续看 [docs/COMMANDS.md](../docs/COMMANDS.md)
- 如果你在准备提交，继续看 [代码质量门禁规范](./code-quality-gates.md)

验收标准：

- CI 规范基于真实环境，而不是纸面设计
- 本地门禁和 CI 门禁保持一致

## 9. 当前执行原则

在当前阶段，不进行大规模目录搬迁。

当前优先级固定为：

1. 建立治理基线
2. 补齐 build、release、quality 规范
3. 调整入口和 skill 分层
4. 最后处理 CI 自动化
