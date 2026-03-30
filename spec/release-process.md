# 发布流程规范

本文件定义 gitcode-cli 的正式发布流程、版本约束、发布检查和发布后验证要求。

## 1. 目标

发布流程需要满足以下目标：

- 发布动作可复现
- 版本号和产物命名一致
- 发布说明与真实产物一致
- 发布前后都可验证
- 文档与安装方式同步更新

## 2. 适用范围

本规范适用于：

- 版本发布准备
- release notes 编写
- tag 创建
- release 产物上传
- 发布后验证

本规范不定义 CI 自动发布细节。当前 GitCode CI 条件未具备，CI 相关内容留到最后阶段单独定义。

## 3. 权威边界

- 正式发布流程以本文件为准
- 发布用命令和产物说明以本文件和 `docs/PACKAGING.md` 共同约束
- 历史性的 GitHub Actions 说明保留在根目录 `RELEASE.md`，但不能替代本规范

## 4. 版本规则

项目版本遵循语义化版本：

```text
vMAJOR.MINOR.PATCH[-PRERELEASE]
```

示例：

- `v1.0.0`
- `v1.0.1`
- `v1.1.0`
- `v2.0.0`
- `v1.0.0-beta.1`

发布时必须保证：

- git tag 与 release tag 一致
- 产物文件名中的版本与 tag 一致
- release notes 中的下载示例与实际版本一致

## 5. 发布前置条件

正式发布前必须满足：

- 目标改动已合并到 `main`
- 当前主线处于可发布状态
- 本地构建与相关测试已完成
- 真实命令验证已完成
- 文档已同步
- 无未解决的 blocker 级问题

## 6. 标准发布流程

当前仓库的标准发布流程如下：

1. 切换到最新 `main`
2. 运行测试与本地构建验证
3. 准备 release notes
4. 使用标准脚本构建发布产物
5. 创建 tag
6. 创建 release
7. 上传 release 产物
8. 执行发布后验证

### 6.1 获取最新主线

```bash
git checkout main
git pull origin main
```

### 6.2 最低发布前验证

```bash
go test ./...
go build -o ./gc ./cmd/gc
./gc version
```

如发布涉及命令行为变更，还必须执行真实命令验证。

### 6.3 构建发布产物

优先使用：

```bash
./scripts/package.sh <version> release
```

### 6.4 创建 release

当前仓库的 CLI 路径为：

```bash
gc release create <tag> -R gitcode-cli/cli --title "<title>" --notes "<notes>"
gc release upload <tag> <files...> -R gitcode-cli/cli
```

### 6.5 创建 tag

```bash
git tag -a vX.Y.Z -m "Release vX.Y.Z"
git push origin vX.Y.Z
```

如果当前发布依赖平台或流程限制，允许先完成 release 说明和产物验证，再执行推送动作，但不得跳过验证。

## 7. Release Notes 规则

release notes 必须满足：

- 描述本次更新内容
- 明确修复的 issue 或功能范围
- 提供完整安装方式
- 下载链接必须是完整路径
- 版本号必须与实际产物一致

不允许：

- 只写文件名，不写完整下载地址
- 使用与实际版本不一致的安装命令
- 在未验证产物存在前写入下载示例

## 8. 发布后验证

正式发布后至少完成以下验证：

- release 页面存在且信息正确
- 发布产物名称正确
- 下载链接可访问
- 至少抽样验证一个安装路径
- `gc version` 可正常输出版本信息

若发布包含 DEB / RPM / wheel，建议至少各抽样验证一种常用安装路径。

## 9. 文档同步要求

发布流程、版本策略或安装方式变化时，必须同步检查：

- `README.md`
- `docs/PACKAGING.md`
- `docs/COMMANDS.md`
- `spec/release-process.md`
- `spec/build-and-package.md`
- `AGENTS.md`
- `CLAUDE.md`
- 相关 AI skills

## 10. 禁止事项

以下行为不允许出现：

- 未验证产物即创建正式 release
- release notes 中使用错误版本号
- 发布说明仍引用旧文件名或旧下载路径
- 用个人临时脚本替代仓库标准流程而不更新文档
- 把平台不支持的自动发布能力写成既成事实

## 11. 当前执行基线

在当前无 GitCode CI 环境的前提下，发布执行基线为：

1. 以本地验证为主
2. 以脚本化构建为主
3. 以人工确认 release notes 和产物一致性为主
4. CI 自动化规则放到后续阶段单独补齐
