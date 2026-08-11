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

function createFileSymlinkOrSkip(t, target, link) {
  try {
    fs.symlinkSync(target, link, "file");
  } catch (error) {
    if (process.platform === "win32" && ["EPERM", "EACCES"].includes(error.code)) {
      t.skip("creating symlinks requires Windows Developer Mode or elevated privileges");
      return false;
    }
    throw error;
  }
  return true;
}

test("install transaction migrates a same-directory gitcode alias and restores it on rollback", (t) => {
  const root = fs.mkdtempSync(path.join(require("os").tmpdir(), "gc-install-alias-"));
  const source = path.join(root, "source");
  const target = path.join(root, "gc");
  const alias = path.join(root, "gitcode");
  fs.writeFileSync(source, "new");
  fs.writeFileSync(target, "old");
  if (!createFileSymlinkOrSkip(t, "gc", alias)) return;

  const transaction = [
    replacePath(source, target, "alias-rollback"),
    replacePath(source, alias, "alias-rollback", { allowedSymlinkTarget: target }),
  ];
  assert.strictEqual(fs.lstatSync(alias).isSymbolicLink(), false);
  assert.strictEqual(fs.readFileSync(alias, "utf8"), "new");

  rollbackTransaction(transaction);
  assert.strictEqual(fs.readFileSync(target, "utf8"), "old");
  assert.strictEqual(fs.lstatSync(alias).isSymbolicLink(), true);
  assert.strictEqual(fs.readlinkSync(alias), "gc");
  assert.strictEqual(fs.readFileSync(alias, "utf8"), "old");
  for (const record of transaction) assert.strictEqual(fs.existsSync(record.backup), false);
});

test("install commit removes the backup of a migrated same-directory alias", (t) => {
  const root = fs.mkdtempSync(path.join(require("os").tmpdir(), "gc-install-alias-commit-"));
  const source = path.join(root, "source");
  const target = path.join(root, "gc");
  const alias = path.join(root, "gitcode");
  fs.writeFileSync(source, "new");
  fs.writeFileSync(target, "old");
  if (!createFileSymlinkOrSkip(t, "gc", alias)) return;

  const record = replacePath(source, alias, "alias-commit", { allowedSymlinkTarget: target });
  commitTransaction([record]);
  assert.strictEqual(fs.lstatSync(alias).isSymbolicLink(), false);
  assert.strictEqual(fs.readFileSync(alias, "utf8"), "new");
  assert.strictEqual(fs.existsSync(record.backup), false);
});

test("install rejects aliases that point anywhere except the sibling gc binary", (t) => {
  const root = fs.mkdtempSync(path.join(require("os").tmpdir(), "gc-install-alias-reject-"));
  const source = path.join(root, "source");
  const target = path.join(root, "gc");
  const other = path.join(root, "other");
  const alias = path.join(root, "gitcode");
  fs.writeFileSync(source, "new");
  fs.writeFileSync(target, "old");
  fs.writeFileSync(other, "unrelated");
  if (!createFileSymlinkOrSkip(t, "other", alias)) return;

  assert.throws(
    () => replacePath(source, alias, "alias-reject", { allowedSymlinkTarget: target }),
    /refusing non-regular install target/
  );
  assert.strictEqual(fs.lstatSync(alias).isSymbolicLink(), true);
  assert.strictEqual(fs.readlinkSync(alias), "other");
  assert.strictEqual(fs.readFileSync(other, "utf8"), "unrelated");
});

test("install rejects a lexical sibling alias that resolves through an external directory symlink", (t) => {
  if (process.platform === "win32") {
    t.skip("nested symlink traversal is a POSIX migration boundary");
    return;
  }
  const root = fs.mkdtempSync(path.join(require("os").tmpdir(), "gc-install-alias-traversal-"));
  const external = fs.mkdtempSync(path.join(require("os").tmpdir(), "gc-install-alias-external-"));
  const source = path.join(root, "source");
  const target = path.join(root, "gc");
  const hop = path.join(root, "hop");
  const alias = path.join(root, "gitcode");
  const externalInner = path.join(external, "inner");
  const externalTarget = path.join(external, "gc");
  fs.writeFileSync(source, "new");
  fs.writeFileSync(target, "allowed");
  fs.mkdirSync(externalInner);
  fs.writeFileSync(externalTarget, "external");
  fs.symlinkSync(externalInner, hop, "dir");
  fs.symlinkSync("hop/../gc", alias, "file");
  assert.strictEqual(fs.readFileSync(alias, "utf8"), "external");

  assert.throws(
    () => replacePath(source, alias, "alias-traversal", { allowedSymlinkTarget: target }),
    /refusing non-regular install target/
  );
  assert.strictEqual(fs.readlinkSync(alias), "hop/../gc");
  assert.strictEqual(fs.readFileSync(externalTarget, "utf8"), "external");
  assert.strictEqual(fs.readFileSync(target, "utf8"), "allowed");
});

test("install still rejects a symlink at the primary gc target", (t) => {
  const root = fs.mkdtempSync(path.join(require("os").tmpdir(), "gc-install-primary-reject-"));
  const source = path.join(root, "source");
  const other = path.join(root, "other");
  const target = path.join(root, "gc");
  fs.writeFileSync(source, "new");
  fs.writeFileSync(other, "old");
  if (!createFileSymlinkOrSkip(t, "other", target)) return;

  assert.throws(() => replacePath(source, target, "primary-reject"), /refusing non-regular install target/);
  assert.strictEqual(fs.lstatSync(target).isSymbolicLink(), true);
  assert.strictEqual(fs.readFileSync(other, "utf8"), "old");
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

test("rollback reports a missing backup and preserves the current target", () => {
  const root = fs.mkdtempSync(path.join(require("os").tmpdir(), "gc-install-missing-backup-"));
  const source = path.join(root, "source");
  const target = path.join(root, "gc");
  fs.writeFileSync(source, "new");
  fs.writeFileSync(target, "old");
  const record = replacePath(source, target, "missing-backup");
  fs.unlinkSync(record.backup);

  assert.throws(
    () => rollbackTransaction([record]),
    (error) => error instanceof AggregateError &&
      error.errors.length === 1 &&
      error.errors[0].code === "EROLLBACK"
  );
  assert.strictEqual(fs.readFileSync(target, "utf8"), "new");
});

test("replacePath reports both the replacement and internal restore failures", () => {
  const root = fs.mkdtempSync(path.join(require("os").tmpdir(), "gc-install-internal-restore-"));
  const source = path.join(root, "source");
  const target = path.join(root, "gc");
  const backup = `${target}.backup-internal-restore`;
  fs.writeFileSync(source, "new");
  fs.writeFileSync(target, "old");
  const originalRenameSync = fs.renameSync;
  fs.renameSync = (src, dst) => {
    const error = new Error(src === backup ? "restore denied" : "replacement denied");
    error.code = "EIO";
    throw error;
  };
  try {
    assert.throws(
      () => replacePath(source, target, "internal-restore"),
      (error) => error instanceof AggregateError &&
        error.errors.length === 2 &&
        error.errors[0].message === "replacement denied" &&
        error.errors[1].message === "restore denied"
    );
  } finally {
    fs.renameSync = originalRenameSync;
  }
  assert.strictEqual(fs.readFileSync(target, "utf8"), "old");
  assert.strictEqual(fs.readFileSync(backup, "utf8"), "old");
});

test("replacePath reports a backup that disappears before internal restore", () => {
  const root = fs.mkdtempSync(path.join(require("os").tmpdir(), "gc-install-lost-backup-"));
  const source = path.join(root, "source");
  const target = path.join(root, "gc");
  const backup = `${target}.backup-lost-backup`;
  fs.writeFileSync(source, "new");
  fs.writeFileSync(target, "old");
  const originalRenameSync = fs.renameSync;
  let replacementAttempts = 0;
  fs.renameSync = (src, dst) => {
    if (src.startsWith(`${target}.tmp-`)) {
      replacementAttempts += 1;
      const error = new Error(replacementAttempts === 1 ? "replace collision" : "replacement failed");
      error.code = replacementAttempts === 1 ? "EPERM" : "EIO";
      if (replacementAttempts === 2) fs.unlinkSync(backup);
      throw error;
    }
    return originalRenameSync(src, dst);
  };
  try {
    assert.throws(
      () => replacePath(source, target, "lost-backup"),
      (error) => error instanceof AggregateError &&
        error.errors.length === 2 &&
        error.errors[0].message === "replacement failed" &&
        error.errors[1].code === "EROLLBACK"
    );
  } finally {
    fs.renameSync = originalRenameSync;
  }
  assert.strictEqual(fs.existsSync(target), false);
  assert.strictEqual(fs.existsSync(backup), false);
});
