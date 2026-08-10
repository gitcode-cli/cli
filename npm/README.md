# @gitcode-cli/cli

The open-source [GitCode CLI](https://gitcode.com/gitcode-cli/cli) distributed via npm with **bundled multi-platform binaries** (Linux x64/ARM64, macOS x64/ARM64, Windows x64). No separate download step; `gc` / `gitcode` works right after install.

GitCode CLI (`gc` / `gitcode`) is the community-developed, MIT-licensed command line tool for GitCode — bringing repositories, issues, pull requests, releases and Actions back to the terminal for developers, scripts, and AI agents. It is an independent open-source project, not published by the GitCode platform team.

## Install

### One-line bootstrap (no global npm install)

```bash
npx @gitcode-cli/cli@latest install
```

Copies the platform binary to a global bin dir (`/usr/local/bin` if writable, else `~/.local/bin`) and installs bash/zsh/fish completions.
On Windows it installs both `gc.exe` and `gitcode.exe`; use `gitcode` in PowerShell because `gc` is the built-in `Get-Content` alias.

### Global npm install

```bash
npm install -g @gitcode-cli/cli
gitcode version
```

## Installation diagnostics and updates

```bash
gitcode doctor install
gitcode doctor install --json
gitcode update --check
gitcode update
```

Global npm and npm-bootstrap installations default to `auto`: after a normal command, a detached helper checks the stable `latest` tag at most once every 24 hours and applies it for the next launch. Updates use a process lock, health check, and rollback. Registry, permission, or offline failures never change the completed command's exit code.

```bash
gitcode config set update.mode auto
gitcode config set update.mode notify
gitcode config set update.mode off
GC_NO_UPDATE_CHECK=1 gitcode version
```

`CI=true` and `--no-interactive` disable automatic application. The updater touches only `@gitcode-cli/cli`; it never invokes pip, Homebrew, apt, dnf, or rpm, and it never rewrites PATH. If another `gitcode` is earlier on PATH, npm's non-failing install check and `doctor install` report the exact candidates and remediation choices.

## Supported platforms

| OS | Arch |
| --- | --- |
| Linux | x64, arm64 |
| macOS | x64, arm64 |
| Windows | x64 |

Unsupported platforms fall back to PyPI/Homebrew/DEB/RPM/release binaries: https://gitcode.com/gitcode-cli/cli/releases

## Links

- Repository: https://gitcode.com/gitcode-cli/cli
- Issues: https://gitcode.com/gitcode-cli/cli/issues
- Command reference: https://gitcode.com/gitcode-cli/cli/blob/main/docs/COMMANDS.md
- Feishu/Lark notifications: https://gitcode.com/gitcode-cli/cli/blob/main/docs/LARK.md

## License

MIT
