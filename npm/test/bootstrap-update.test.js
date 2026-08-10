"use strict";

const test = require("node:test");
const assert = require("node:assert");
const fs = require("fs");
const os = require("os");
const path = require("path");
const { compareVersions, parseArgs, stableVersion, updateMode } = require("../lib/bootstrap-update-helper");

test("bootstrap updater accepts only stable versions", () => {
  assert.deepStrictEqual(stableVersion("1.2.3"), [1, 2, 3]);
  assert.strictEqual(stableVersion("1.2.3-beta.1"), null);
  assert.strictEqual(compareVersions("1.2.3", "1.2.4"), -1);
});

test("bootstrap updater parses detached update arguments", () => {
  const manifest = path.join(fs.mkdtempSync(path.join(os.tmpdir(), "gc-manifest-")), "install.json");
  assert.deepStrictEqual(
    parseArgs(["--background", "--force", "--parent-pid", "42", "--manifest", manifest]),
    { background: true, check: false, force: true, json: false, manifest, parentPid: 42 }
  );
});

test("bootstrap updater honors enterprise update mode", () => {
  assert.strictEqual(updateMode({ GC_UPDATE_MODE: "notify" }), "notify");
  assert.strictEqual(updateMode({ GC_UPDATE_MODE: "off" }), "off");
});
