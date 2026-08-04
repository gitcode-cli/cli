// Bootstrap ("install" subcommand) for @gitcode-cli/cli.
//
// Copies the bundled platform binary to a global bin directory and installs
// shell completions, so `npx @gitcode-cli/cli install` works without a prior
// `npm i -g`. Dependency-free (Node built-ins only). Pure helpers are
// exported separately so they can be unit-tested without touching the FS.

"use strict";

const fs = require("fs");
const path = require("path");
const os = require("os");
const { spawnSync } = require("child_process");
const { resolveBinaryName, isSupported } = require("./platform");

const PLATFORMS_DIR = path.join(__dirname, "..", "bin", "platforms");

function bundledBinaryPath() {
  return path.join(PLATFORMS_DIR, resolveBinaryName(process.platform, process.arch));
}

/**
 * Choose a global bin dir. Prefer a no-sudo writable location; fall back to a
 * per-user dir under the home directory. Pure (no FS side effects beyond the
 * write probe on the candidate dir).
 */
function chooseGlobalBinDir(home, isWin) {
  if (isWin) {
    return path.join(home, "AppData", "Local", "gitcode-cli", "bin");
  }
  for (const dir of ["/usr/local/bin", path.join(home, ".local", "bin")]) {
    try {
      fs.mkdirSync(dir, { recursive: true });
      const probe = path.join(dir, ".gc-write-probe");
      fs.writeFileSync(probe, "");
      fs.unlinkSync(probe);
      return dir;
    } catch {
      // not writable; try next
    }
  }
  const dir = path.join(home, ".local", "bin");
  fs.mkdirSync(dir, { recursive: true });
  return dir;
}

// Completion target dirs per shell (user-writable, auto-loaded where possible).
// Pure: derives the path from the shell + home.
function completionTarget(shell, home) {
  switch (shell) {
    case "bash":
      return path.join(home, ".local", "share", "bash-completion", "completions", "gc");
    case "zsh":
      return path.join(home, ".zsh", "completions", "_gc");
    case "fish":
      return path.join(home, ".config", "fish", "completions", "gc.fish");
    default:
      return null;
  }
}

function ensureExec(file) {
  if (process.platform === "win32") return;
  try {
    fs.chmodSync(file, 0o755);
  } catch {
    /* best-effort */
  }
}

function copyFile(src, dst) {
  fs.copyFileSync(src, dst);
  ensureExec(dst);
}

function runGc(bin, args) {
  return spawnSync(bin, args, { encoding: "utf8" });
}

function installCompletions(bin, home) {
  const installed = [];
  for (const shell of ["bash", "zsh", "fish"]) {
    const res = runGc(bin, ["completion", shell]);
    if (res.status !== 0 || !res.stdout) continue;
    const target = completionTarget(shell, home);
    if (!target) continue;
    try {
      fs.mkdirSync(path.dirname(target), { recursive: true });
      fs.writeFileSync(target, res.stdout, { mode: 0o644 });
      installed.push(`${shell}: ${target}`);
    } catch {
      /* skip unwritable */
    }
  }
  return installed;
}

// Whether the per-user fallback dir is on PATH (pure).
function dirOnPath(dir) {
  return (process.env.PATH || "").split(path.delimiter).includes(dir);
}

async function runInstall() {
  const home = os.homedir();
  const isWin = process.platform === "win32";

  if (!isSupported(process.platform, process.arch)) {
    throw new Error(
      `no bundled binary for ${process.platform}/${process.arch}. ` +
        `Install via PyPI/Homebrew/DEB/RPM instead: https://gitcode.com/gitcode-cli/cli/releases`
    );
  }

  const src = bundledBinaryPath();
  if (!fs.existsSync(src)) {
    throw new Error(
      `bundled binary missing at ${src}. The npm package may be incomplete; ` +
        `reinstall @gitcode-cli/cli.`
    );
  }
  ensureExec(src);

  const dir = chooseGlobalBinDir(home, isWin);
  fs.mkdirSync(dir, { recursive: true });

  const dst = path.join(dir, isWin ? "gc.exe" : "gc");
  copyFile(src, dst);

  // gitcode alias (symlink on posix, copy on windows).
  if (!isWin) {
    const alias = path.join(dir, "gitcode");
    try {
      fs.unlinkSync(alias);
    } catch {
      /* fine */
    }
    try {
      fs.symlinkSync("gc", alias, "file");
    } catch {
      copyFile(src, alias);
    }
  }

  // Verify.
  const v = runGc(dst, ["version"]);
  const versionLine = (v.stdout || "").split("\n")[0] || "(gc version failed)";

  // Completions (posix only; Windows shell completion differs).
  const completions = isWin ? [] : installCompletions(dst, home);

  process.stdout.write(`Installed gc to ${dir}\n`);
  process.stdout.write(`  ${versionLine}\n`);
  if (completions.length) {
    process.stdout.write(`Shell completions installed:\n`);
    for (const c of completions) process.stdout.write(`  ${c}\n`);
  } else if (isWin) {
    process.stdout.write(`Shell completions: skipped on Windows. Run "gc completion bash|powershell" manually if needed.\n`);
  } else {
    process.stdout.write(`Shell completions: skipped (none writable). Run "gc completion bash|zsh|fish" manually.\n`);
  }

  // PATH hint.
  if (isWin) {
    if (!dirOnPath(dir)) {
      process.stdout.write(
        `\nAdd ${dir} to your PATH (run in PowerShell, then reopen terminals):\n` +
          `  setx PATH "${dir};%PATH%"\n`
      );
    }
  } else if (dir === path.join(home, ".local", "bin")) {
    if (!dirOnPath(dir)) {
      process.stdout.write(
        `\nAdd ${dir} to your PATH:\n` +
          `  echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc  # or ~/.zshrc\n`
      );
    }
  }
  process.stdout.write(`\nRun "gc --help" to get started.\n`);
}

module.exports = { runInstall, chooseGlobalBinDir, completionTarget, dirOnPath };
