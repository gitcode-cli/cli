# Contributing to GitCode CLI

Thank you for your interest in contributing to GitCode CLI!

## Documentation

Before contributing, please read the following documentation:

| Document | Description |
|----------|-------------|
| [COMMANDS.md](./docs/COMMANDS.md) | Command usage guide - update when adding/modifying commands |
| [PACKAGING.md](./docs/PACKAGING.md) | Packaging and release guide - DEB/RPM build instructions |
| [CLAUDE.md](./CLAUDE.md) | AI-assisted development guide - coding conventions and workflow |
| [README.md](./README.md) | Project overview and installation guide |

> **Important**: When modifying command-related code, you MUST sync updates to `docs/COMMANDS.md`. See the documentation maintenance section in COMMANDS.md for details.

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

### Prerequisites

1. GitHub repository with write access
2. Docker Hub account (for Docker images)
3. GitHub tokens for Homebrew tap and Scoop bucket

### Create a Release

```bash
# 1. Update version in files
# 2. Commit changes
git commit -am "chore: prepare for release vX.Y.Z"

# 3. Create and push tag
git tag -a vX.Y.Z -m "Release vX.Y.Z"
git push origin main --tags

# 4. GitHub Actions will automatically:
#    - Run tests
#    - Build binaries for all platforms
#    - Create DEB/RPM packages
#    - Build and push Docker images
#    - Update Homebrew tap
#    - Update Scoop bucket
```

### Manual Release

```bash
# Install goreleaser
go install github.com/goreleaser/goreleaser/v2@latest

# Create release
make release

# Or create local snapshot
make release-local
```

## Required Secrets

For automated releases, configure these secrets in GitHub:

| Secret | Description |
|--------|-------------|
| `GITHUB_TOKEN` | Automatically provided by GitHub |
| `DOCKER_USERNAME` | Docker Hub username |
| `DOCKER_PASSWORD` | Docker Hub password/access token |
| `HOMEBREW_TAP_GITHUB_TOKEN` | Token for homebrew-tap repo |
| `SCOOP_BUCKET_GITHUB_TOKEN` | Token for scoop-bucket repo |
| `GPG_FINGERPRINT` | GPG key for signing (optional) |

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