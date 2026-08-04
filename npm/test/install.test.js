// Unit tests for lib/install.js pure helpers (dir selection, completion
// targets, PATH membership). The bootstrap side effects (copy/symlink/spawn)
// are exercised end-to-end by the release workflow; here we lock the logic.
// Run: node --test (from npm/) or `npm test`.

"use strict";

const test = require("node:test");
const assert = require("node:assert");
const fs = require("fs");
const path = require("path");
const { chooseGlobalBinDir, completionTarget, dirOnPath } = require("../lib/install");

test("chooseGlobalBinDir returns a writable, existing dir on posix (regardless of /usr/local/bin)", () => {
  // Deterministic: on hosted runners /usr/local/bin may be writable, so we
  // only assert the invariant — the returned dir exists and is writable —
  // rather than which specific dir is chosen.
  const home = "/tmp/gc-install-home-" + process.pid;
  const dir = chooseGlobalBinDir(home, false);
  assert.ok(fs.existsSync(dir), `${dir} should exist`);
  // Writable probe (the function already guarantees writability by design).
  const probe = path.join(dir, ".gc-test-probe");
  fs.writeFileSync(probe, "");
  fs.unlinkSync(probe);
});

test("chooseGlobalBinDir returns the Windows per-user dir on win32 (no FS writability dependency)", () => {
  const home = "/u/home";
  const dir = chooseGlobalBinDir(home, true);
  assert.strictEqual(dir, path.join(home, "AppData", "Local", "gitcode-cli", "bin"));
});

test("completionTarget maps each shell to a standard path", () => {
  const home = "/u/home";
  assert.strictEqual(completionTarget("bash", home), path.join(home, ".local", "share", "bash-completion", "completions", "gc"));
  assert.strictEqual(completionTarget("zsh", home), path.join(home, ".zsh", "completions", "_gc"));
  assert.strictEqual(completionTarget("fish", home), path.join(home, ".config", "fish", "completions", "gc.fish"));
  assert.strictEqual(completionTarget("powershell", home), null);
});

test("dirOnPath reflects the current PATH delimiter", () => {
  // The real PATH contains some entries; a synthetic dir not in it is false.
  assert.strictEqual(dirOnPath("/this/dir/is/not/on/path"), false);
});
