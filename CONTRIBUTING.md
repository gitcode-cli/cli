# Contributing to GitCode CLI

Thank you for your interest in contributing to GitCode CLI!

## Documentation

Before contributing, please read the following documentation:

### Core Documentation

| Document | Description |
|----------|-------------|
| [README.md](./README.md) | Project overview and installation guide |
| [CLAUDE.md](./CLAUDE.md) | AI-assisted development guide |

### Development Specifications (spec/)

| Document | Description |
|----------|-------------|
| [编码规范](./spec/coding-standards.md) | Naming, file structure, error handling |
| [测试指南](./spec/testing-guide.md) | Unit tests, actual command testing |
| [命令开发模板](./spec/command-template.md) | Command development template |
| [安全规范](./spec/security.md) | Token management, sensitive info protection |

### Development Workflows (spec/workflows/)

| Document | Description |
|----------|-------------|
| [Issue 流程](./spec/workflows/issue-workflow.md) | Issue lifecycle management |
| [PR 流程](./spec/workflows/pr-workflow.md) | PR creation and merge |
| [评审流程](./spec/workflows/review-workflow.md) | Code review process |
| [测试流程](./spec/workflows/test-workflow.md) | Testing workflow |
| [构建打包流程](./spec/workflows/build-workflow.md) | Build and package |
| [Release 流程](./spec/workflows/release-workflow.md) | Release process |

### Other Documentation

| Document | Description |
|----------|-------------|
| [docs/COMMANDS.md](./docs/COMMANDS.md) | Command usage guide |
| [docs/PACKAGING.md](./docs/PACKAGING.md) | Packaging guide |

> **Important**: When modifying command-related code, you MUST sync updates to `docs/COMMANDS.md`.

## Development Setup

### Prerequisites

- Go 1.22+
- Make
- Docker (optional)
- GoReleaser (optional, for releases)

### Building

```bash
# Clone the repository
git clone https://gitcode.com/gitcode-cli/cli.git
cd gitcode-cli

# Install dependencies
make deps

# Build
make build

# Run tests
make test

# Run
make run
```

### Development Commands

```bash
make build          # Build binary
make run            # Run application
make test           # Run tests
make test-coverage  # Run tests with coverage
make fmt            # Format code
make lint           # Run linter
make completions    # Generate shell completions
```

## Packaging

### Binary

```bash
# Build for current platform
make build

# Build for all platforms
make build-all
```

### Docker

```bash
# Build Docker image
make docker

# Build with specific tag
make docker DOCKER_TAG=v1.0.0

# Run Docker container
make docker-run

# Multi-platform build
make docker-all
```

### Linux Packages (DEB/RPM)

```bash
# Install goreleaser
go install github.com/goreleaser/goreleaser/v2@latest

# Create snapshot release (local testing)
make release-local

# Check generated packages in dist/
```

### Homebrew

```bash
# Local formula is in homebrew/gc.rb
# After release, users can install with:
brew install gitcode-com/tap/gc
```

### Scoop (Windows)

```bash
# Manifest is in scoop/gc.json
# After release, users can install with:
scoop bucket add gitcode-com https://github.com/gitcode-com/scoop-bucket
scoop install gc
```

## Release Process

> 详细发布流程请参阅 [RELEASE.md](./RELEASE.md)。

### Quick Start

```bash
# 创建并推送标签触发自动发布
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions 自动执行：
# - 构建 RPM/DEB 包（x86_64, arm64）
# - 发布到 GitHub Release
# - 发布 PyPI 包 (pip install gitcode-cli)
```

### Required Secrets

| Secret | Description |
|--------|-------------|
| `GITHUB_TOKEN` | Automatically provided by GitHub |

> PyPI 使用 Trusted Publishing（OIDC），无需配置 API Token。

## Directory Structure

```
.
├── cmd/gc/              # Main application entry
├── pkg/                 # Public packages
│   ├── cmd/             # Command implementations
│   ├── cmdutil/         # Utilities for commands
│   ├── iostreams/       # I/O stream management
│   └── testutil/        # Testing utilities
├── internal/            # Private packages
├── api/                 # API client
├── git/                 # Git operations
├── completions/         # Shell completions
├── scripts/             # Package scripts
├── homebrew/            # Homebrew formula
├── scoop/               # Scoop manifest
├── .github/workflows/   # GitHub Actions
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose
├── .goreleaser.yaml     # GoReleaser config
└── Makefile             # Build automation
```

## Code Style

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `make fmt` before committing
- Run `make lint` to check for issues
- Write tests for new functionality

## Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Run linter: `make lint`
6. Commit with conventional commits
7. Push and create PR

## Commit Convention

We use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `test:` Test changes
- `refactor:` Code refactoring
- `chore:` Maintenance tasks

## License

By contributing, you agree that your contributions will be licensed under the MIT License.