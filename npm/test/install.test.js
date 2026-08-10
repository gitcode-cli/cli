// Unit tests for lib/install.js pure helpers (dir selection, completion
// targets, PATH membership). The bootstrap side effects (copy/symlink/spawn)
// are exercised end-to-end by the release workflow; here we lock the logic.
// Run: node --test (from npm/) or `npm test`.

"use strict";

const test = require("node:test");
const assert = require("node:assert");
const fs = require("fs");
const path = require("path");
const {
  chooseGlobalBinDir, commitTransaction, completionTarget, dirOnPath, parseInstallArgs,
  installHelp, quotePowerShell, replacePath, rollbackTransaction,
} = require("../lib/install");

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

test("dirOnPath is case-insensitive and separator-insensitive on Windows", () => {
  const env = { PATH: `C:\\Tools\\GitCode\\;C:\\Windows` };
  assert.strictEqual(dirOnPath("c:\\tools\\gitcode", env, true), true);
});

test("parseInstallArgs accepts only an explicit target directory", () => {
  assert.deepStrictEqual(parseInstallArgs([]), { targetDir: "" });
  assert.strictEqual(parseInstallArgs(["--target-dir", "."]).targetDir, path.resolve("."));
  assert.throws(() => parseInstallArgs(["--force"]), /unknown install argument/);
});

test("install help documents the explicit target directory", () => {
  assert.match(installHelp(), /--target-dir <directory>/);
});

test("PowerShell guidance escapes single quotes in target directories", () => {
  assert.strictEqual(quotePowerShell("C:\\Users\\O'Brien\\bin;"), "'C:\\Users\\O''Brien\\bin;'");
});

test("install rollback restores only files touched by the current transaction", () => {
  const root = fs.mkdtempSync(path.join(require("os").tmpdir(), "gc-install-transaction-"));
  const source = path.join(root, "source");
  const first = path.join(root, "gc");
  const second = path.join(root, "gitcode");
  const untouched = path.join(root, "helper");
  fs.writeFileSync(source, "new");
  fs.writeFileSync(first, "first-current");
  fs.writeFileSync(second, "second-current");
  fs.writeFileSync(untouched, "untouched-current");
  fs.writeFileSync(`${untouched}.previous`, "stale-previous");

  const transaction = [replacePath(source, first, "test"), replacePath(source, second, "test")];
  rollbackTransaction(transaction);

  assert.strictEqual(fs.readFileSync(first, "utf8"), "first-current");
  assert.strictEqual(fs.readFileSync(second, "utf8"), "second-current");
  assert.strictEqual(fs.readFileSync(untouched, "utf8"), "untouched-current");
});

test("install commit removes transaction-scoped backups", () => {
  const root = fs.mkdtempSync(path.join(require("os").tmpdir(), "gc-install-commit-"));
  const source = path.join(root, "source");
  const target = path.join(root, "gc");
  fs.writeFileSync(source, "new");
  fs.writeFileSync(target, "old");
  const record = replacePath(source, target, "commit");
  commitTransaction([record]);
  assert.strictEqual(fs.readFileSync(target, "utf8"), "new");
  assert.strictEqual(fs.existsSync(record.backup), false);
});

test("rollback continues restoring remaining paths after one restore fails", () => {
  const root = fs.mkdtempSync(path.join(require("os").tmpdir(), "gc-install-rollback-failure-"));
  const source = path.join(root, "source");
  const target = path.join(root, "gc");
  const blocked = path.join(root, "blocked-directory");
  fs.writeFileSync(source, "new");
  fs.writeFileSync(target, "old");
  fs.mkdirSync(blocked);
  const good = replacePath(source, target, "failure-test");

  assert.throws(
    () => rollbackTransaction([good, { dst: blocked, backup: `${blocked}.backup`, hadOriginal: false }]),
    (error) => error instanceof AggregateError && error.errors.length === 1
  );
  assert.strictEqual(fs.readFileSync(target, "utf8"), "old");
});
