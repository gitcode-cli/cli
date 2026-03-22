# gitcode-cli

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

GitCode 官方命令行工具，为 GitCode 用户提供便捷的命令行操作体验。

## 安装

```bash
# Homebrew
brew install gitcode-com/tap/gc

# 或手动下载
curl -sSL https://gitcode.com/install.sh | sh
```

## 快速开始

```bash
# 登录 GitCode
gc auth login

# 克隆仓库
gc repo clone owner/repo

# 创建 Issue
gc issue create --title "Bug report"

# 创建 MR
gc mr create --title "New feature"

# 代码检视
gc mr review 123 --approve
```

## 功能特性

- **认证管理**: OAuth Device Flow、Token 认证、多账户支持
- **仓库操作**: 克隆、创建、Fork、查看仓库
- **Issue 管理**: 创建、查看、列表、关闭 Issue
- **MR 管理**: 创建、查看、检出、合并 MR
- **代码检视**: 批准、请求修改、添加评论

## 命令概览

| 命令 | 说明 |
|------|------|
| `gc auth` | 认证管理 |
| `gc repo` | 仓库操作 |
| `gc issue` | Issue 管理 |
| `gc mr` | MR 管理 |
| `gc config` | 配置管理 |
| `gc version` | 版本信息 |

## 文档

- [需求文档](./issues-plan/) - 完整需求规格和里程碑规划
- [CLAUDE.md](./CLAUDE.md) - AI 辅助开发指南
- [API 文档](https://gitcode.com/docs/api) - GitCode API 参考

## 技术栈

| 组件 | 技术 |
|------|------|
| 语言 | Go 1.21+ |
| 命令框架 | Cobra |
| 配置格式 | YAML |
| 目标平台 | Linux, macOS, Windows |

## 开发

```bash
# 克隆仓库
git clone https://gitcode.com/gitcode-com/gitcode-cli.git
cd gitcode-cli

# 构建
make build

# 测试
make test
```

## 贡献

欢迎贡献代码！请查看 [贡献指南](./CONTRIBUTING.md)。

## 许可证

[MIT License](./LICENSE)