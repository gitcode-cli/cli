#!/usr/bin/env node
// Standalone updater copied next to npm-bootstrap binaries.

"use strict";

const fs = require("fs");
const os = require("os");
const path = require("path");
const { spawnSync } = require("child_process");

const PACKAGE = "@gitcode-cli/cli";
const TTL_MS = 24 * 60 * 60 * 1000;
const LOCK_STALE_MS = 15 * 60 * 1000;
const SENSITIVE_ENV_KEYS = ["GC_TOKEN", "GITCODE_TOKEN", "NPM_TOKEN", "NODE_AUTH_TOKEN"];

function parseArgs(args) {
  const options = { background: false, check: false, force: false, json: false, manifest: "", parentPid: 0 };
  for (let index = 0; index < args.length; index += 1) {
    const arg = args[index];
    if (arg === "--background") options.background = true;
    else if (arg === "--check") options.check = true;
    else if (arg === "--force") options.force = true;
    else if (arg === "--json") options.json = true;
    else if (arg === "--manifest" && args[index + 1]) options.manifest = path.resolve(args[++index]);
    else if (arg === "--parent-pid" && args[index + 1]) options.parentPid = Number(args[++index]);
    else throw new Error(`unknown updater argument: ${arg}`);
  }
  if (!options.manifest) options.manifest = path.join(__dirname, ".gitcode-install.json");
  return options;
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

function statePath(env = process.env) {
  let root = env.GC_STATE_DIR;
  if (!root && process.platform === "win32") {
    root = path.join(env.LOCALAPPDATA || path.join(os.homedir(), "AppData", "Local"), "gitcode-cli");
  }
  if (!root) root = path.join(env.XDG_STATE_HOME || path.join(os.homedir(), ".local", "state"), "gitcode-cli");
  return path.join(root, "update-state.json");
}

function updateMode(env = process.env) {
  const requested = String(env.GC_UPDATE_MODE || "").toLowerCase();
  if (["auto", "notify", "off"].includes(requested)) return requested;
  const configRoot = env.GC_CONFIG_DIR || path.join(os.homedir(), ".config", "gc");
  const config = readJSON(path.join(configRoot, "config.json"));
  const stored = (((config.hosts || {})["gitcode.com"] || {})["update.mode"] || "").toLowerCase();
  return ["auto", "notify", "off"].includes(stored) ? stored : "auto";
}

function updaterEnvironment(env = process.env) {
  const clean = { ...env, GC_NO_UPDATE_CHECK: "1" };
  for (const key of SENSITIVE_ENV_KEYS) delete clean[key];
  return clean;
}

function appendLog(message) {
  const file = path.join(path.dirname(statePath()), "update.log");
  fs.mkdirSync(path.dirname(file), { recursive: true, mode: 0o700 });
  fs.appendFileSync(file, `${new Date().toISOString()} ${message}\n`, { mode: 0o600 });
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

function npmCommand(manifest) {
  if (manifest.npm && /\.m?js$/i.test(manifest.npm)) return { command: manifest.node, prefix: [manifest.npm] };
  return { command: manifest.npm || (process.platform === "win32" ? "npm.cmd" : "npm"), prefix: [] };
}

function runNpm(manifest, args, timeout) {
  const npm = npmCommand(manifest);
  return spawnSync(npm.command, [...npm.prefix, ...args], {
    encoding: "utf8",
    timeout,
    windowsHide: true,
    env: updaterEnvironment(),
  });
}

function latestVersion(manifest) {
  const result = runNpm(manifest, ["view", PACKAGE, "dist-tags.latest", "--json"], 10000);
  if (result.error) throw result.error;
  if (result.status !== 0) throw new Error((result.stderr || "npm registry check failed").trim());
  const version = JSON.parse(result.stdout);
  if (!stableVersion(version)) throw new Error(`registry latest is not stable: ${version}`);
  return version;
}

function waitForParent(pid) {
  if (!pid) return;
  const deadline = Date.now() + 30000;
  while (Date.now() < deadline) {
    try {
      process.kill(pid, 0);
      Atomics.wait(new Int32Array(new SharedArrayBuffer(4)), 0, 0, 100);
    } catch {
      return;
    }
  }
  throw new Error(`parent process ${pid} did not exit before update timeout`);
}

function restorePrevious(manifest) {
  const names = process.platform === "win32" ? ["gc.exe", "gitcode.exe"] : ["gc", "gitcode"];
  for (const name of names) {
    const target = path.join(manifest.targetDir, name);
    const previous = `${target}.previous`;
    if (fs.existsSync(previous)) fs.copyFileSync(previous, target);
  }
}

function healthCheck(manifest, expected) {
  const entry = path.join(manifest.targetDir, process.platform === "win32" ? "gitcode.exe" : "gitcode");
  const result = spawnSync(entry, ["version", "--json"], {
    encoding: "utf8",
    timeout: 10000,
    windowsHide: true,
    env: updaterEnvironment(),
  });
  if (result.error) throw result.error;
  if (result.status !== 0) throw new Error((result.stderr || "bootstrap health check failed").trim());
  const actual = JSON.parse(result.stdout).version;
  if (compareVersions(actual, expected) !== 0) throw new Error(`updated CLI reported ${actual}, expected ${expected}`);
}

function installLatest(manifest, latest) {
  const args = [
    "exec",
    "--yes",
    `--package=${PACKAGE}@${latest}`,
    "--",
    "gitcode",
    "install",
    "--target-dir",
    manifest.targetDir,
  ];
  const result = runNpm(manifest, args, 120000);
  if (result.error) throw result.error;
  if (result.status !== 0) throw new Error((result.stderr || `npm exec exited ${result.status}`).trim());
  healthCheck(manifest, latest);
}

function acquireLock(file) {
  fs.mkdirSync(path.dirname(file), { recursive: true, mode: 0o700 });
  try {
    return fs.openSync(file, "wx", 0o600);
  } catch (error) {
    if (error.code === "EEXIST") {
      try {
        if (Date.now() - fs.statSync(file).mtimeMs > LOCK_STALE_MS) {
          fs.unlinkSync(file);
          return fs.openSync(file, "wx", 0o600);
        }
      } catch {
        // Another process owns or just released the lock.
      }
      return null;
    }
    throw error;
  }
}

function run(options) {
  const manifest = readJSON(options.manifest);
  if (manifest.distribution !== "npm-bootstrap" || !manifest.targetDir || !manifest.version) {
    throw new Error("invalid npm-bootstrap install manifest");
  }
  const stateFile = statePath();
  const state = readJSON(stateFile);
  if (!options.force && !options.check && state.nextCheck && Number(state.nextCheck) > Date.now()) {
    return { status: "cached", distribution: "npm-bootstrap", current: manifest.version, latest: "", message: "Update check is not due yet." };
  }
  const lock = `${stateFile}.bootstrap.lock`;
  const descriptor = acquireLock(lock);
  if (descriptor == null) return { status: "busy", distribution: "npm-bootstrap", current: manifest.version, latest: "", message: "Another update is running." };
  try {
    const latest = latestVersion(manifest);
    const comparison = compareVersions(manifest.version, latest);
    let result;
    if (comparison == null) throw new Error(`cannot compare ${manifest.version} and ${latest}`);
    if (comparison >= 0) {
      result = { status: "current", distribution: "npm-bootstrap", current: manifest.version, latest, message: `GitCode CLI ${manifest.version} is current.` };
    } else if (options.check || updateMode() === "notify") {
      result = { status: "available", distribution: "npm-bootstrap", current: manifest.version, latest, message: `GitCode CLI ${latest} is available.` };
    } else if (updateMode() === "off" && !options.force) {
      result = { status: "disabled", distribution: "npm-bootstrap", current: manifest.version, latest, message: "Automatic updates are disabled." };
    } else {
      waitForParent(options.parentPid);
      try {
        installLatest(manifest, latest);
      } catch (error) {
        restorePrevious(manifest);
        throw error;
      }
      result = { status: "updated", distribution: "npm-bootstrap", current: manifest.version, latest, message: `Updated GitCode CLI ${manifest.version} -> ${latest}.` };
    }
    const now = Date.now();
    state.lastChecked = new Date(now).toISOString();
    state.nextCheck = now + TTL_MS;
    state.summary = { message: result.message, shown: !options.background };
    writeJSON(stateFile, state);
    appendLog(`status=${result.status} current=${result.current} latest=${result.latest || "none"}`);
    return result;
  } finally {
    fs.closeSync(descriptor);
    try {
      fs.unlinkSync(lock);
    } catch {
      // Stale bootstrap locks are visible and can be removed by the user.
    }
  }
}

function main(args = process.argv.slice(2)) {
  let options;
  try {
    options = parseArgs(args);
    const result = run(options);
    if (!options.background) process.stdout.write(options.json ? `${JSON.stringify(result)}\n` : `${result.message}\n`);
    return 0;
  } catch (error) {
    if (options && options.background) {
      const file = statePath();
      const state = readJSON(file);
      state.nextCheck = Date.now() + TTL_MS;
      state.summary = { message: "Automatic update failed; run gitcode update for details.", shown: false };
      writeJSON(file, state);
      appendLog("status=error; run an explicit update for details");
    }
    if (!options || !options.background) {
      if (options && options.json) process.stdout.write(`${JSON.stringify({ status: "error", distribution: "npm-bootstrap", current: "", latest: "", message: error.message })}\n`);
      else process.stderr.write(`update failed: ${error.message}\n`);
    }
    return 1;
  }
}

if (require.main === module) process.exitCode = main();

module.exports = { compareVersions, main, parseArgs, stableVersion, updateMode };
