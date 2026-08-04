// Unit tests for bin/gc.js behavior that can be exercised without the real
// gc binary: the install interception contract and the unsupported-platform
// error path. The exec path is covered end-to-end elsewhere.
// Run: node --test (from npm/) or `npm test`.

"use strict";

const test = require("node:test");
const assert = require("node:assert");
const { execFileSync } = require("child_process");
const path = require("path");

const WRAPPER = path.join(__dirname, "..", "bin", "gc.js");

function runWrapper(args, env) {
  try {
    const out = execFileSync(process.execPath, [WRAPPER, ...args], {
      encoding: "utf8",
      env: { ...process.env, ...env },
      timeout: 30000,
    });
    return { status: 0, stdout: out, stderr: "" };
  } catch (err) {
    return {
      status: err.status == null ? 1 : err.status,
      stdout: err.stdout || "",
      stderr: err.stderr || "",
    };
  }
}

test("wrapper exits 127 with a clear message when the binary is missing (ENOENT)", () => {
  // The platforms dir has no bundled binary in a source checkout (binaries are
  // built by the release workflow), so the wrapper must hit its ENOENT path.
  const r = runWrapper(["version"]);
  // 127 = command not found (the wrapper's ENOENT path)
  assert.strictEqual(r.status, 127);
  assert.match(r.stderr, /gc binary not found/);
});
