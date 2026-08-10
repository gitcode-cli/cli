"use strict";

const test = require("node:test");
const assert = require("node:assert");
const fs = require("fs");
const os = require("os");
const path = require("path");
const {
  acquireLock,
  compareVersions,
  disabledForInvocation,
  releaseLock,
  shouldSchedule,
  stableVersion,
  updateMode,
  updaterEnvironment,
  updateStatePath,
  writeJSON,
} = require("../lib/update");
const { parseArgs } = require("../lib/update-helper");

function stateEnv() {
  return { GC_STATE_DIR: fs.mkdtempSync(path.join(os.tmpdir(), "gc-update-state-")) };
}

test("accepts only stable semantic versions and compares without downgrading", () => {
  assert.deepStrictEqual(stableVersion("v1.2.3"), [1, 2, 3]);
  assert.strictEqual(stableVersion("1.2.3-rc.1"), null);
  assert.strictEqual(compareVersions("1.2.3", "1.3.0"), -1);
  assert.strictEqual(compareVersions("2.0.0", "1.9.9"), 1);
  assert.strictEqual(compareVersions("1.2.3", "1.2.3"), 0);
});

test("disables updates for explicit opt-out, CI, non-interactive, and off mode", () => {
  assert.strictEqual(disabledForInvocation([], { GC_NO_UPDATE_CHECK: "1" }), true);
  assert.strictEqual(disabledForInvocation([], { CI: "true" }), true);
  assert.strictEqual(disabledForInvocation(["--no-interactive"], {}), true);
  assert.strictEqual(disabledForInvocation(["--no-update-check"], {}), true);
  assert.strictEqual(disabledForInvocation([], { GC_UPDATE_MODE: "off" }), true);
  assert.strictEqual(disabledForInvocation([], { GC_UPDATE_MODE: "notify" }), false);
});

test("mode defaults to auto and honors the environment", () => {
  assert.strictEqual(updateMode({ GC_CONFIG_DIR: fs.mkdtempSync(path.join(os.tmpdir(), "gc-config-")) }), "auto");
  assert.strictEqual(updateMode({ GC_UPDATE_MODE: "notify" }), "notify");
});

test("TTL prevents a background check until it expires", () => {
  const env = stateEnv();
  assert.strictEqual(shouldSchedule([], env, 100), true);
  writeJSON(updateStatePath(env), { nextCheck: 200 });
  assert.strictEqual(shouldSchedule([], env, 100), false);
  assert.strictEqual(shouldSchedule([], env, 201), true);
});

test("cross-process lock permits only one owner", () => {
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), "gc-update-lock-"));
  const lock = path.join(dir, "update.lock");
  const first = acquireLock(lock);
  assert.notStrictEqual(first, null);
  assert.strictEqual(acquireLock(lock), null);
  releaseLock(lock, first);
  const second = acquireLock(lock);
  assert.notStrictEqual(second, null);
  releaseLock(lock, second);
});

test("update helper parser accepts check/json/background only", () => {
  assert.deepStrictEqual(parseArgs(["--check", "--json"]), {
    background: false,
    checkOnly: true,
    json: true,
  });
  assert.throws(() => parseArgs(["--channel", "next"]), /unknown update argument/);
});

test("updater child environment strips GitCode and npm credentials", () => {
  const env = updaterEnvironment({
    PATH: "/bin",
    GC_TOKEN: "gitcode-secret",
    GITCODE_TOKEN: "legacy-secret",
    NPM_TOKEN: "npm-secret",
    NODE_AUTH_TOKEN: "node-secret",
  });
  assert.strictEqual(env.PATH, "/bin");
  assert.strictEqual(env.GC_NO_UPDATE_CHECK, "1");
  for (const key of ["GC_TOKEN", "GITCODE_TOKEN", "NPM_TOKEN", "NODE_AUTH_TOKEN"]) {
    assert.strictEqual(env[key], undefined);
  }
});
