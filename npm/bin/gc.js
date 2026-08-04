// Wrapper entry for @gitcode-cli/cli.
//
// - Default: resolve the bundled platform binary and exec it with the
//   remaining args (so `gc version`, `gc issue list`, ... all work).
// - Special subcommand `install`: run the Node-side bootstrap (copy binary to
//   a global bin dir, install shell completions) so `npx @gitcode-cli/cli
//   install` works without `npm i -g`.
//
// This file must stay dependency-free (no require of external packages) so
// the npm package needs no install step beyond shipping the binaries.

"use strict";

const fs = require("fs");
const path = require("path");
const { spawn } = require("child_process");
const { resolveBinaryName } = require("../lib/platform");

// PLATFORMS_DIR can be overridden (e.g. for deterministic unit tests that
// point at an empty dir to exercise the ENOENT path). Defaults to the bundled
// platforms directory shipped with this package.
const PLATFORMS_DIR = process.env.GC_PLATFORMS_DIR || path.join(__dirname, "platforms");

function binPath() {
  const name = resolveBinaryName(process.platform, process.arch);
  return path.join(PLATFORMS_DIR, name);
}

function runBinary(args) {
  const p = binPath();
  // Ensure the exec bit survives tarball extraction on posix.
  if (process.platform !== "win32") {
    try {
      fs.chmodSync(p, 0o755);
    } catch {
      /* best-effort; ENOENT handled below */
    }
  }
  const child = spawn(p, args, { stdio: "inherit" });
  child.on("error", (err) => {
    if (err.code === "ENOENT") {
      process.stderr.write(
        `gc binary not found at ${p}. ` +
          `Run "npm install -g @gitcode-cli/cli" or "npx @gitcode-cli/cli install" first.\n`
      );
      process.exit(127);
    }
    process.stderr.write(`failed to run gc: ${err.message}\n`);
    process.exit(1);
  });
  child.on("exit", (code, signal) => {
    if (signal) {
      // Forward the signal so shells report it correctly (do not force exit 1).
      try {
        process.kill(process.pid, signal);
      } catch {
        process.exit(1);
      }
      return;
    }
    process.exit(code == null ? 1 : code);
  });
}

function main() {
  const args = process.argv.slice(2);

  // Intercept the install subcommand (lark-cli-style bootstrap). If the Go
  // CLI grows a real "gc install" later, prefer forwarding "--help"/unknown
  // flags through and reserve only the bare "install" first token.
  if (args[0] === "install") {
    const { runInstall } = require("../lib/install");
    runInstall().catch((err) => {
      process.stderr.write(`install failed: ${err && err.message ? err.message : err}\n`);
      process.exit(1);
    });
    return;
  }

  runBinary(args);
}

main();
