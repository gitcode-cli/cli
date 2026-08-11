// Bootstrap ("install" subcommand) for @gitcode-cli/cli.
//
// Copies the bundled platform binary to a global bin directory and, on
// Linux/macOS, installs shell completions. The registry-isolated npx bootstrap works without a prior
// `npm i -g`. Dependency-free (Node built-ins only). Pure helpers are
// exported separately so they can be unit-tested without touching the FS.

"use strict";

const fs = require("fs");
const path = require("path");
const os = require("os");
const crypto = require("crypto");
const { spawnSync } = require("child_process");
const { resolveBinaryName, isSupported } = require("./platform");
const { normalizePath, writeInstallMetadata } = require("./install-metadata");
const pkg = require("../package.json");

const PLATFORMS_DIR = path.join(__dirname, "..", "bin", "platforms");
const BOOTSTRAP_HELPER = path.join(__dirname, "bootstrap-update-helper.js");

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

function pathExists(file) {
  try {
    fs.lstatSync(file);
    return true;
  } catch (error) {
    if (error.code === "ENOENT") return false;
    throw error;
  }
}

function isAllowedAliasSymlink(dst, allowedTarget) {
  if (!allowedTarget) return false;
  try {
    const targetStat = fs.lstatSync(allowedTarget, { bigint: true });
    if (!targetStat.isFile() || targetStat.isSymbolicLink()) return false;
    const linkTargetStat = fs.statSync(dst, { bigint: true });
    return linkTargetStat.isFile() &&
      linkTargetStat.dev === targetStat.dev &&
      linkTargetStat.ino === targetStat.ino;
  } catch (error) {
    if (["ENOENT", "EINVAL"].includes(error.code)) return false;
    throw error;
  }
}

function replacePath(src, dst, transactionID, options = {}) {
  const temp = `${dst}.tmp-${process.pid}-${crypto.randomBytes(8).toString("hex")}`;
  const backup = `${dst}.backup-${transactionID}`;
  const sourceStat = fs.lstatSync(src);
  if (!sourceStat.isFile() || sourceStat.isSymbolicLink()) {
    throw new Error(`refusing non-regular install source: ${src}`);
  }
  let hadOriginal = false;
  let moveOriginal = false;
  try {
    const stat = fs.lstatSync(dst);
    if (stat.isSymbolicLink() && isAllowedAliasSymlink(dst, options.allowedSymlinkTarget)) {
      moveOriginal = true;
    } else if (!stat.isFile() || stat.isSymbolicLink()) {
      throw new Error(`refusing non-regular install target: ${dst}`);
    }
    hadOriginal = true;
  } catch (error) {
    if (error.code !== "ENOENT") throw error;
  }
  try {
    fs.lstatSync(backup);
    throw new Error(`refusing existing transaction backup: ${backup}`);
  } catch (error) {
    if (error.code !== "ENOENT") throw error;
  }
  fs.copyFileSync(src, temp, fs.constants.COPYFILE_EXCL);
  ensureExec(temp);
  let backupReady = false;
  try {
    if (hadOriginal) {
      if (moveOriginal) {
        fs.renameSync(dst, backup);
        backupReady = true;
        if (!isAllowedAliasSymlink(backup, options.allowedSymlinkTarget)) {
          throw new Error(`refusing non-regular install target: ${dst}`);
        }
      } else {
        fs.copyFileSync(dst, backup, fs.constants.COPYFILE_EXCL);
        backupReady = true;
      }
    }
    renameReplace(temp, dst);
    return { dst, backup, hadOriginal };
  } catch (error) {
    let restoreError;
    try {
      if (backupReady) {
        if (!pathExists(backup)) {
          restoreError = new Error(`transaction backup missing: ${backup}`);
          restoreError.code = "EROLLBACK";
        } else {
          renameReplace(backup, dst);
        }
      }
    } catch (caught) {
      // Keep the transaction backup for manual recovery.
      restoreError = caught;
    }
    try {
      fs.unlinkSync(temp);
    } catch {
      // Preserve the original error.
    }
    if (restoreError) {
      throw new AggregateError(
        [error, restoreError],
        "path replacement failed and rollback was incomplete",
        { cause: error }
      );
    }
    throw error;
  }
}

function renameReplace(src, dst) {
  try {
    fs.renameSync(src, dst);
  } catch (error) {
    if (!["EEXIST", "EPERM"].includes(error.code)) throw error;
    fs.unlinkSync(dst);
    fs.renameSync(src, dst);
  }
}

function rollbackTransaction(records) {
  const failures = [];
  for (const record of [...records].reverse()) {
    if (record.hadOriginal) {
      try {
        if (!pathExists(record.backup)) {
          const error = new Error(`transaction backup missing: ${record.backup}`);
          error.code = "EROLLBACK";
          throw error;
        }
        renameReplace(record.backup, record.dst);
      } catch (error) {
        failures.push(error);
      }
      continue;
    }
    try {
      fs.unlinkSync(record.dst);
    } catch (error) {
      if (error.code !== "ENOENT") failures.push(error);
    }
  }
  if (failures.length) {
    throw new AggregateError(failures, `failed to restore ${failures.length} install path(s)`);
  }
}

function commitTransaction(records) {
  for (const record of records) {
    try {
      fs.unlinkSync(record.backup);
    } catch {
      // The install is already committed; an orphaned unique backup is safer
      // than rolling back a healthy installation because cleanup failed.
    }
  }
}

function sha256(file) {
  return crypto.createHash("sha256").update(fs.readFileSync(file)).digest("hex");
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
function dirOnPath(dir, env = process.env, isWin = process.platform === "win32") {
  const wanted = normalizePath(dir, isWin);
  const delimiter = isWin ? ";" : ":";
  return (env.PATH || env.Path || "")
    .split(delimiter)
    .filter(Boolean)
    .some((entry) => normalizePath(entry.replace(/^"|"$/g, ""), isWin) === wanted);
}

function parseInstallArgs(args) {
  const options = { targetDir: "" };
  for (let index = 0; index < args.length; index += 1) {
    if (args[index] === "--target-dir" && args[index + 1]) {
      options.targetDir = path.resolve(args[index + 1]);
      index += 1;
      continue;
    }
    throw new Error(`unknown install argument: ${args[index]}`);
  }
  return options;
}

function quotePowerShell(value) {
  return `'${String(value).replace(/'/g, "''")}'`;
}

function installHelp() {
  return [
    "Install bundled gc and gitcode binaries outside the npm package directory.",
    "",
    "Usage:",
    "  gitcode install [--target-dir <directory>]",
    "",
    "Flags:",
    "  --target-dir <directory>  Install into an explicit directory",
    "  -h, --help                Show this help",
    "",
  ].join("\n");
}

async function runInstall(args = []) {
  if (args.length === 1 && ["-h", "--help"].includes(args[0])) {
    process.stdout.write(installHelp());
    return;
  }
  const options = parseInstallArgs(args);
  const home = os.homedir();
  const isWin = process.platform === "win32";

  if (!isSupported(process.platform, process.arch)) {
    const guidance = process.platform === "win32" && process.arch === "arm64"
      ? "Windows arm64 is not shipped yet; build from source with Go."
      : "Use a release channel that explicitly lists this OS and architecture.";
    throw new Error(
      `no bundled binary for ${process.platform}/${process.arch}. ${guidance} ` +
        `https://gitcode.com/gitcode-cli/cli/releases`
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

  const dir = options.targetDir || chooseGlobalBinDir(home, isWin);
  fs.mkdirSync(dir, { recursive: true });

  const dst = path.join(dir, isWin ? "gc.exe" : "gc");
  const alias = path.join(dir, isWin ? "gitcode.exe" : "gitcode");
  const helper = path.join(dir, "gitcode-update-helper.js");
  const aliasOptions = isWin ? {} : { allowedSymlinkTarget: dst };
  const transactionID = `${process.pid}-${crypto.randomBytes(8).toString("hex")}`;
  const transaction = [];
  let versionLine;
  try {
    transaction.push(replacePath(src, dst, transactionID));
    transaction.push(replacePath(src, alias, transactionID, aliasOptions));
    transaction.push(replacePath(BOOTSTRAP_HELPER, helper, transactionID));
    if (sha256(src) !== sha256(dst) || sha256(src) !== sha256(alias)) {
      throw new Error("installed binary checksum verification failed");
    }

    const v = runGc(dst, ["version"]);
    if (v.status !== 0) {
      throw new Error(`installed binary health check failed: ${(v.stderr || "unknown error").trim()}`);
    }
    versionLine = (v.stdout || "").split("\n")[0] || "(gc version failed)";
    writeInstallMetadata(dir, {
      distribution: "npm-bootstrap",
      version: pkg.version,
      targetDir: dir,
      node: process.execPath,
      npm: process.env.npm_execpath || "",
      helper,
      sha256: sha256(src),
    });
    commitTransaction(transaction);
  } catch (error) {
    try {
      rollbackTransaction(transaction);
    } catch (rollbackError) {
      throw new AggregateError(
        [error, rollbackError],
        "installation failed and rollback was incomplete",
        { cause: error }
      );
    }
    throw error;
  }

  // Completions (posix only; Windows shell completion differs).
  const completions = isWin ? [] : installCompletions(dst, home);

  process.stdout.write(`Installed gc and gitcode to ${dir}\n`);
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
          `\nAdd ${dir} to your user PATH, then reopen terminals. PowerShell example:\n` +
          `  $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')\n` +
          `  [Environment]::SetEnvironmentVariable('Path', ${quotePowerShell(`${dir};`)} + $userPath, 'User')\n` +
          `  $env:Path = ${quotePowerShell(`${dir};`)} + $env:Path\n`
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
  process.stdout.write(`\nRun "${isWin ? "gitcode" : "gc"} --help" to get started.\n`);
}

module.exports = {
  runInstall, chooseGlobalBinDir, commitTransaction, completionTarget, dirOnPath,
  installHelp, parseInstallArgs, quotePowerShell, replacePath, rollbackTransaction,
};
