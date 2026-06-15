# .loop

`.loop/` 保存 gitcode-cli 项目的 Loop Engineering 配置。

## 文件

- `project.yaml`：项目、仓库、事实源和外部标准包声明
- `policy.yaml`：本项目 loop 门禁策略
- `hooks.yaml`：本项目使用的 hook 阶段映射

## 边界

- `.loop/` 中的配置可以提交。
- `.loop/runtime/` 是本地临时缓存，不提交。
- 长期状态和证据写回 GitCode issue / PR。
