---
name: gc-dev-setup
description: Set up local development environment for gitcode-cli project. Use this skill when user says "初始化本地开发环境", "搭建本地开发环境", "init dev environment", "setup local dev", or when user has development needs and you should first check if local dev environment is ready. The check standard is: local code compiles without errors.
---

# GC Development Environment Setup

Set up the local development environment for the gitcode-cli (gc) project.

## When to Use This Skill

Trigger this skill proactively when:
- User says "初始化本地开发环境" or "搭建本地开发环境"
- User says "init dev environment", "setup local dev", "初始化开发环境"
- User has development needs (creating features, fixing bugs) - first verify the environment is ready
- User just pulled the code and needs to start working

## Environment Check Standard

The local development environment is considered "ready" when:
- Go is installed and accessible
- Project builds successfully (`./gc` binary exists and runs)
- No compilation errors

## Workflow

### Step 1: Pull Latest Code

```bash
git pull origin main
# or current branch
git pull
```

If not in a git repository or need to clone first, guide the user to clone the project.

### Step 2: Check Go Environment

Verify Go is installed:

```bash
go version
```

If Go is not installed:
1. On Ubuntu/Debian: `sudo apt update && sudo apt install -y golang-go`
2. Or install from https://go.dev/dl/
3. For China users, may need to set GOPROXY: `export GOPROXY=https://goproxy.cn,direct`

### Step 3: Build the Project

Build the gc binary:

```bash
# Set GOPROXY for China users if needed
export GOPROXY=https://goproxy.cn,direct

# Build
go build -o ./gc ./cmd/gc
```

### Step 4: Verify Build

Run the version command to verify:

```bash
./gc version
```

If the command runs without errors, the environment is ready. No need to check authentication status.

Expected output:
```
gc version dev
  commit: none
  built:  unknown
https://gitcode.com/gitcode-cli/cli
```

## Quick Environment Check

Use this to quickly verify the environment is ready:

```bash
# Do a fresh build and verify
GOPROXY=https://goproxy.cn,direct go build -o ./gc ./cmd/gc && ./gc version
```

If no errors, the environment is ready.

## Common Issues

### Go not found
- Install Go: `sudo apt install -y golang-go` (Ubuntu)
- Or download from https://go.dev/dl/

### Build fails with network errors
- Set GOPROXY: `export GOPROXY=https://goproxy.cn,direct`

### Permission denied
- Make sure not using `sudo` for go build
- Check file permissions

### Old version in PATH
- Do NOT copy `./gc` to `~/bin/` or other PATH directories
- Always use `./gc` directly from the project directory

## Output Format

After setup, provide a simple summary:

```
✓ 本地开发环境已就绪
- Go 版本: go1.22.x
- gc 二进制: ./gc
- 状态: 构建成功，命令正常
```

No need to check authentication status. The environment is considered ready once `./gc version` runs without errors.