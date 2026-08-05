# @gitcode-cli/cli

The open-source [GitCode CLI](https://gitcode.com/gitcode-cli/cli) distributed via npm with **bundled multi-platform binaries** (Linux x64/ARM64, macOS x64/ARM64, Windows x64). No separate download step; `gc` / `gitcode` works right after install.

GitCode CLI (`gc` / `gitcode`) is the community-developed, MIT-licensed command line tool for GitCode — bringing repositories, issues, pull requests, releases and Actions back to the terminal for developers, scripts, and AI agents. It is an independent open-source project, not published by the GitCode platform team.

## Install

### One-line bootstrap (no global npm install)

```bash
npx @gitcode-cli/cli@latest install
```

Copies the platform binary to a global bin dir (`/usr/local/bin` if writable, else `~/.local/bin`) and installs bash/zsh/fish completions.

### Global npm install

```bash
npm install -g @gitcode-cli/cli
gc version
```

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
