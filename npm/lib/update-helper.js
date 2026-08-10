#!/usr/bin/env node

"use strict";

const { runUpdate } = require("./update");

function parseArgs(args) {
  const options = { background: false, checkOnly: false, json: false };
  for (const arg of args) {
    if (arg === "--background") options.background = true;
    else if (arg === "--check") options.checkOnly = true;
    else if (arg === "--json") options.json = true;
    else throw new Error(`unknown update argument: ${arg}`);
  }
  return options;
}

function main(args = process.argv.slice(2)) {
  let options;
  try {
    options = parseArgs(args);
    const result = runUpdate(options);
    if (!options.background) {
      process.stdout.write(options.json ? `${JSON.stringify(result)}\n` : `${result.message}\n`);
    }
    return 0;
  } catch (error) {
    if (!options || !options.background) {
      if (options && options.json) {
        process.stdout.write(`${JSON.stringify({ status: "error", distribution: "npm", current: "", latest: "", message: error.message })}\n`);
      } else {
        process.stderr.write(`update failed: ${error.message}\n`);
      }
    }
    return 1;
  }
}

if (require.main === module) process.exitCode = main();

module.exports = { main, parseArgs };
