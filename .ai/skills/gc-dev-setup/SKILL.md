# 共享 skill: gc-dev-setup

## 目标

定义 gitcode-cli 项目本地开发环境检查和初始化的共享场景。

## 统一约束

- 本地构建规则以 `spec/delivery/build-and-package.md` 为准
- 本地测试规则以 `spec/foundations/testing-guide.md` 为准
- 真实命令验证只能使用 `infra-test/*`
- 仓库内 AI 本地开发闭环以 `spec/workflows/ai-local-development-workflow.md` 为准

## 最小检查顺序

1. 确认 Go 环境和依赖可用
2. 确认 `GC_TOKEN` 或等价认证路径可用
3. 运行 `go test ./...`
4. 运行 `go build -o ./gc ./cmd/gc`
5. 运行 `./gc version`
6. 运行 `./scripts/regression-core.sh`
7. 如涉及命令行为改动，补 `infra-test/*` 上的真实命令验证

## 适配层说明

- Claude 适配：`.claude/skills/gc-dev-setup/`
- Codex 适配：`.codex/skills/gc-dev-setup/`
