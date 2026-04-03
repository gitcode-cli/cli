#!/usr/bin/env python3
"""Classify change risk for local diffs.

Usage examples:
  ./scripts/classify-change-risk.py --base origin/main
  ./scripts/classify-change-risk.py spec/workflows/review-workflow.md
"""

from __future__ import annotations

import argparse
import subprocess
import sys
from pathlib import Path


LOW_PREFIXES = ("docs/", "spec/", ".ai/", ".codex/", ".claude/")
LOW_FILES = {"README.md", "AGENTS.md", "CLAUDE.md", "CONTRIBUTING.md", "RELEASE.md"}
HIGH_KEYWORDS = (
    "auth",
    "config",
    "delete",
    "confirm",
    "release",
    "security",
    "token",
    "permission",
    "credential",
)
MEDIUM_PREFIXES = ("cmd/", "pkg/", "internal/", "api/", "git/")
MEDIUM_SUFFIXES = (".go",)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("paths", nargs="*", help="Optional explicit file paths")
    parser.add_argument(
        "--base",
        help="Git base ref used to diff changed files when no explicit paths are passed",
    )
    return parser.parse_args()


def run_lines(cmd: list[str]) -> list[str]:
    result = subprocess.run(cmd, check=False, capture_output=True, text=True)
    if result.returncode != 0:
        raise SystemExit(result.stderr.strip() or f"failed to run: {' '.join(cmd)}")
    return [line.strip() for line in result.stdout.splitlines() if line.strip()]


def changed_files(base: str) -> list[str]:
    paths: list[str] = []
    seen: set[str] = set()
    commands = [
        ["git", "diff", "--name-only", f"{base}...HEAD"],
        ["git", "diff", "--name-only"],
        ["git", "diff", "--name-only", "--cached"],
        ["git", "ls-files", "--others", "--exclude-standard"],
    ]

    for cmd in commands:
        for path in run_lines(cmd):
            if path in seen:
                continue
            seen.add(path)
            paths.append(path)
    return paths


def classify_path(path: str) -> tuple[str, str]:
    normalized = path.replace("\\", "/")
    lower = normalized.lower()

    if any(keyword in lower for keyword in HIGH_KEYWORDS):
        return ("high", f"{normalized}: matched high-risk keyword")

    if normalized in LOW_FILES or normalized.startswith(LOW_PREFIXES):
        return ("low", f"{normalized}: documentation or workflow asset")

    if normalized.startswith(MEDIUM_PREFIXES) or normalized.endswith(MEDIUM_SUFFIXES):
        return ("medium", f"{normalized}: runtime or command implementation path")

    return ("medium", f"{normalized}: default medium risk")


def aggregate(levels: list[str]) -> str:
    if "high" in levels:
        return "high"
    if "medium" in levels:
        return "medium"
    return "low"


def main() -> int:
    args = parse_args()

    paths = [str(Path(path)) for path in args.paths]
    if not paths:
        if not args.base:
            print("error: pass explicit paths or --base", file=sys.stderr)
            return 2
        paths = changed_files(args.base)

    if not paths:
        print("risk=low")
        print("- no changed files detected")
        return 0

    findings = [classify_path(path) for path in paths]
    overall = aggregate([level for level, _reason in findings])

    print(f"risk={overall}")
    for level, reason in findings:
        print(f"- {level}: {reason}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
