// Safe npm-channel update checks and exact-version upgrades.

"use strict";

const fs = require("fs");
const path = require("path");
const { spawnSync } = require("child_process");
const pkg = require("../package.json");
const { expectedGlobalBin, readInstallMetadata, stateDir } = require("./install-metadata");

const PACKAGE = "@gitcode-cli/cli";
const TTL_MS = 24 * 60 * 60 * 1000;
const LOCK_STALE_MS = 15 * 60 * 1000;
const SENSITIVE_ENV_KEYS = ["GC_TOKEN", "GITCODE_TOKEN", "NPM_TOKEN", "NODE_AUTH_TOKEN"];

function updateStatePath(env = process.env) {
  return path.join(stateDir(env), "update-state.json");
}

function readJSON(file, fallback = {}) {
  try {
    return JSON.parse(fs.readFileSync(file, "utf8"));
  } catch {
    return fallback;
  }
}

function writeJSON(file, value) {
  fs.mkdirSync(path.dirname(file), { recursive: true, mode: 0o700 });
  const temp = `${file}.tmp-${process.pid}`;
  fs.writeFileSync(temp, `${JSON.stringify(value, null, 2)}\n`, { mode: 0o600 });
  try {
    fs.unlinkSync(file);
  } catch (error) {
    if (error.code !== "ENOENT") throw error;
  }
  fs.renameSync(temp, file);
}

function configDir(env = process.env) {
  if (env.GC_CONFIG_DIR) return env.GC_CONFIG_DIR;
  return path.join(require("os").homedir(), ".config", "gc");
}

function updaterEnvironment(env = process.env) {
  const clean = { ...env, GC_NO_UPDATE_CHECK: "1" };
  for (const key of SENSITIVE_ENV_KEYS) delete clean[key];
  return clean;
}

function appendLog(message, env = process.env) {
  const file = path.join(stateDir(env), "update.log");
  fs.mkdirSync(path.dirname(file), { recursive: true, mode: 0o700 });
  fs.appendFileSync(file, `${new Date().toISOString()} ${message}\n`, { mode: 0o600 });
}

function updateMode(env = process.env) {
  const fromEnv = (env.GC_UPDATE_MODE || "").toLowerCase();
  if (["auto", "notify", "off"].includes(fromEnv)) return fromEnv;
  const config = readJSON(path.join(configDir(env), "config.json"));
  const configured = (((config.hosts || {})["gitcode.com"] || {})["update.mode"] || "").toLowerCase();
  return ["auto", "notify", "off"].includes(configured) ? configured : "auto";
}

function disabledForInvocation(args = [], env = process.env) {
  const disabled = (env.GC_NO_UPDATE_CHECK || "").toLowerCase();
  return (
    disabled === "1" ||
    disabled === "true" ||
    String(env.CI || "").toLowerCase() === "true" ||
    args.includes("--no-update-check") ||
    args.includes("--no-interactive") ||
    updateMode(env) === "off"
  );
}

function stableVersion(value) {
  const match = String(value || "").trim().match(/^v?(\d+)\.(\d+)\.(\d+)$/);
  return match ? match.slice(1).map(Number) : null;
}

function compareVersions(left, right) {
  const a = stableVersion(left);
  const b = stableVersion(right);
  if (!a || !b) return null;
  for (let index = 0; index < 3; index += 1) {
    if (a[index] !== b[index]) return a[index] < b[index] ? -1 : 1;
  }
  return 0;
}

function npmCommand(metadata) {
  if (metadata && metadata.npm) {
    if (/\.m?js$/i.test(metadata.npm)) return { command: process.execPath, prefix: [metadata.npm] };
    return { command: metadata.npm, prefix: [] };
  }
  return { command: process.platform === "win32" ? "npm.cmd" : "npm", prefix: [] };
}

function runNpm(metadata, args, timeout = 15000) {
  const npm = npmCommand(metadata);
  return spawnSync(npm.command, [...npm.prefix, ...args], {
    encoding: "utf8",
    timeout,
    windowsHide: true,
    env: updaterEnvironment(),
  });
}

function checkLatest(metadata) {
  const result = runNpm(metadata, ["view", PACKAGE, "dist-tags.latest", "--json"], 10000);
  if (result.error) throw result.error;
  if (result.status !== 0) throw new Error((result.stderr || "npm registry check failed").trim());
  const latest = JSON.parse(result.stdout);
  if (!stableVersion(latest)) throw new Error(`registry latest is not a stable semantic version: ${latest}`);
  return latest;
}

function globalEntrypoint(metadata) {
  const bin = expectedGlobalBin(metadata.prefix, process.platform === "win32");
  return path.join(bin, process.platform === "win32" ? "gitcode.cmd" : "gitcode");
}

function healthCheck(metadata, expectedVersion) {
  const entrypoint = globalEntrypoint(metadata);
  const result = spawnSync(entrypoint, ["version", "--json"], {
    encoding: "utf8",
    timeout: 10000,
    windowsHide: true,
    env: updaterEnvironment(),
  });
  if (result.error) throw result.error;
  if (result.status !== 0) throw new Error((result.stderr || "updated CLI health check failed").trim());
  const version = JSON.parse(result.stdout).version;
  if (compareVersions(version, expectedVersion) !== 0) {
    throw new Error(`updated CLI reported ${version}, expected ${expectedVersion}`);
  }
}

function installExact(metadata, version) {
  const result = runNpm(
    metadata,
    ["install", "-g", `${PACKAGE}@${version}`, "--no-audit", "--no-fund"],
    120000
  );
  if (result.error) throw result.error;
  if (result.status !== 0) throw new Error((result.stderr || `npm install exited ${result.status}`).trim());
}

function acquireLock(file, now = Date.now()) {
  fs.mkdirSync(path.dirname(file), { recursive: true, mode: 0o700 });
  try {
    return fs.openSync(file, "wx", 0o600);
  } catch (error) {
    if (error.code !== "EEXIST") throw error;
    try {
      if (now - fs.statSync(file).mtimeMs > LOCK_STALE_MS) {
        fs.unlinkSync(file);
        return fs.openSync(file, "wx", 0o600);
      }
    } catch {
      // Another process owns or just released the lock.
    }
    return null;
  }
}

function releaseLock(file, descriptor) {
  if (descriptor == null) return;
  try {
    fs.closeSync(descriptor);
  } finally {
    try {
      fs.unlinkSync(file);
    } catch {
      // Best effort; stale locks are recovered on the next attempt.
    }
  }
}

function resultObject(status, latest, message) {
  return { status, distribution: "npm", current: pkg.version, latest: latest || "", message };
}

function performUpdate(options = {}) {
  const packageRoot = options.packageRoot || path.resolve(__dirname, "..");
  const metadata = options.metadata || readInstallMetadata(packageRoot);
  if (!metadata || !metadata.global || metadata.distribution !== "npm" || !metadata.prefix) {
    throw new Error("automatic update is available only for a global npm installation");
  }
  const latest = checkLatest(metadata);
  const comparison = compareVersions(pkg.version, latest);
  if (comparison == null) throw new Error(`cannot compare versions ${pkg.version} and ${latest}`);
  if (comparison >= 0) return resultObject("current", latest, `GitCode CLI ${pkg.version} is current.`);
  if (options.checkOnly || options.mode === "notify") {
    return resultObject("available", latest, `GitCode CLI ${latest} is available (current ${pkg.version}).`);
  }

  try {
    installExact(metadata, latest);
    healthCheck(metadata, latest);
    return resultObject("updated", latest, `Updated GitCode CLI ${pkg.version} -> ${latest}.`);
  } catch (error) {
    try {
      installExact(metadata, pkg.version);
      healthCheck(metadata, pkg.version);
      throw new Error(`${error.message}; restored ${pkg.version}`);
    } catch (rollbackError) {
      if (rollbackError.message.includes("restored")) throw rollbackError;
      throw new Error(`${error.message}; rollback failed: ${rollbackError.message}`);
    }
  }
}

function runUpdate(options = {}) {
  const stateFile = options.stateFile || updateStatePath();
  const lockFile = `${stateFile}.lock`;
  const descriptor = acquireLock(lockFile);
  if (descriptor == null) return resultObject("busy", "", "Another GitCode CLI update is already running.");
  try {
    const mode = options.mode || updateMode();
    const result = performUpdate({ ...options, mode });
    const now = Date.now();
    const state = readJSON(stateFile);
    state.lastChecked = new Date(now).toISOString();
    state.nextCheck = now + TTL_MS;
    state.summary = { message: result.message, shown: !options.background };
    writeJSON(stateFile, state);
    appendLog(`status=${result.status} current=${result.current} latest=${result.latest || "none"}`);
    return result;
  } catch (error) {
    const now = Date.now();
    const state = readJSON(stateFile);
    state.lastChecked = new Date(now).toISOString();
    state.nextCheck = now + TTL_MS;
    state.summary = { message: "Automatic update failed; run gitcode update for details.", shown: !options.background };
    writeJSON(stateFile, state);
    appendLog("status=error; run an explicit update for details");
    throw error;
  } finally {
    releaseLock(lockFile, descriptor);
  }
}

function shouldSchedule(args = [], env = process.env, now = Date.now()) {
  if (disabledForInvocation(args, env)) return false;
  const state = readJSON(updateStatePath(env));
  return !state.nextCheck || Number(state.nextCheck) <= now;
}

function showPendingSummary(stderr = process.stderr, env = process.env) {
  const file = updateStatePath(env);
  const state = readJSON(file);
  if (!state.summary || state.summary.shown) return;
  stderr.write(`${state.summary.message}\n`);
  state.summary.shown = true;
  writeJSON(file, state);
}

function showFirstRunNotice(stderr = process.stderr, env = process.env) {
  if (disabledForInvocation([], env)) return;
  const file = updateStatePath(env);
  const state = readJSON(file);
  if (state.noticeShown) return;
  stderr.write(
    "GitCode CLI installed by npm checks daily and automatically applies stable updates after commands.\n" +
      'Set "gitcode config set update.mode notify|off" or GC_NO_UPDATE_CHECK=1 to change this behavior.\n'
  );
  state.noticeShown = true;
  writeJSON(file, state);
}

module.exports = {
  PACKAGE,
  TTL_MS,
  acquireLock,
  appendLog,
  checkLatest,
  compareVersions,
  disabledForInvocation,
  performUpdate,
  readJSON,
  releaseLock,
  runUpdate,
  shouldSchedule,
  showFirstRunNotice,
  showPendingSummary,
  stableVersion,
  updateMode,
  updaterEnvironment,
  updateStatePath,
  writeJSON,
};
