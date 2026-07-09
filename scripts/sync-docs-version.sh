#!/usr/bin/env bash
# Sync release version strings in README.md and docs/PACKAGING.md.
#
# These two files are committed, user-facing docs whose download URLs and
# examples carry a pinned release version. Unlike the build-time version
# sync in scripts/package.sh and .github/workflows/release.yml (which cover
# nfpm-*.yaml / pyproject.toml / gc_cli/__init__.py), these docs are not
# synced automatically, so they lag behind every release (see #314).
#
# This script derives the target version from a single input, auto-detects
# the version currently referenced in the docs, and replaces it. Run it as
# a release step and commit the result.
#
# Usage:
#   ./scripts/sync-docs-version.sh <version> [--dry-run]
#
#   <version>   Target release version, e.g. v0.7.0 or 0.7.0
#   --dry-run   Print the planned changes without writing files
#
# Exit codes:
#   0  success (or already at target version, no-op)
#   1  invalid input / detection failure / residual old version after replace
#   2  usage error

set -euo pipefail

dry_run=0
target=""
for arg in "${1:-}" "${2:-}"; do
    case "${arg}" in
        --dry-run) dry_run=1 ;;
        "") ;;
        -*) echo "Unknown option: ${arg}" >&2; exit 2 ;;
        *) target="${arg}" ;;
    esac
done

if [[ -z "${target}" ]]; then
    echo "Usage: $0 <version> [--dry-run]" >&2
    echo "Example: $0 v0.7.0" >&2
    exit 2
fi

if [[ ! "${target}" =~ ^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]; then
    echo "Invalid version: ${target}" >&2
    echo "Expected vMAJOR.MINOR.PATCH or MAJOR.MINOR.PATCH (no prerelease)" >&2
    exit 1
fi

new_num="${target#v}"
new_tag="v${new_num}"

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/.." && pwd)"
readme="${repo_root}/README.md"
packaging="${repo_root}/docs/PACKAGING.md"

for f in "${readme}" "${packaging}"; do
    if [[ ! -f "${f}" ]]; then
        echo "Missing file: ${f}" >&2
        exit 1
    fi
done

# Detect the version currently referenced in README via the first
# releases/download/vX.Y.Z/ occurrence.
old_num="$(grep -oE 'releases/download/v[0-9]+\.[0-9]+\.[0-9]+/' "${readme}" | head -1 | sed -E 's#releases/download/v([0-9]+\.[0-9]+\.[0-9]+)/#\1#')"
if [[ -z "${old_num}" ]]; then
    echo "Could not detect current doc version in ${readme}" >&2
    echo "Expected a 'releases/download/vX.Y.Z/' URL" >&2
    exit 1
fi
old_tag="v${old_num}"

if [[ "${old_num}" == "${new_num}" ]]; then
    echo "Docs already reference ${new_tag}; nothing to do."
    exit 0
fi

echo "Syncing doc version: ${old_tag} -> ${new_tag}"
for f in "${readme}" "${packaging}"; do
    before=$(grep -c -E "v?${old_num//./\\.}" "${f}" || true)
    if [[ ${dry_run} -eq 1 ]]; then
        echo "  [dry-run] would touch ${before} line(s) in $(basename "${f}")"
        continue
    fi
    sed -i "s/v${old_num}/v${new_num}/g; s/${old_num}/${new_num}/g" "${f}"
    after_old=$(grep -c -E "v?${old_num//./\\.}" "${f}" || true)
    after_new=$(grep -c -E "v?${new_num//./\\.}" "${f}" || true)
    echo "  $(basename "${f}"): ${before} -> ${after_new} line(s) with ${new_tag}; residual ${old_tag}: ${after_old}"
    if [[ ${after_old} -gt 0 ]]; then
        echo "ERROR: residual ${old_tag} in $(basename "${f}") after replace" >&2
        exit 1
    fi
done

if [[ ${dry_run} -eq 1 ]]; then
    echo "[dry-run] no files written."
else
    echo "Done. Verify with: git diff --stat README.md docs/PACKAGING.md"
fi
