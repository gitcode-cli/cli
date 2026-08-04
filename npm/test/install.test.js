// Unit tests for lib/install.js pure helpers (dir selection, completion
// targets, PATH membership). The bootstrap side effects (copy/symlink/spawn)
// are exercised end-to-end by the release workflow; here we lock the logic.
// Run: node --test (from npm/) or `npm test`.

"use strict";

const test = require("node:test");
const assert = require("node:assert");
const path = require("path");
const { chooseGlobalBinDir, completionTarget, dirOnPath } = require("../lib/install");

test("chooseGlobalBinDir returns the Windows per-user dir on win32", () => {
  const home = "/u/home";
  const dir = chooseGlobalBinDir(home, true);
  assert.strictEqual(dir, path.join(home, "AppData", "Local", "gitcode-cli", "bin"));
});

test("chooseGlobalBinDir falls back to ~/.local/bin when /usr/local/bin is not writable", () => {
  const home = "/tmp/definitely-not-writable-home-" + process.pid;
  const dir = chooseGlobalBinDir(home, false);
  assert.strictEqual(dir, path.join(home, ".local", "bin"));
  // the fallback dir is created
  const fs = require("fs");
  assert.ok(fs.existsSync(dir));
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
