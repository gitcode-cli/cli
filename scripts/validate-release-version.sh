#!/usr/bin/env bash
# Validate a release workflow version input and print the package version.

set -euo pipefail

version="${1:-}"

if [[ -z "${version}" ]]; then
    echo "Invalid release version: value is required" >&2
    exit 1
fi

if [[ ! "${version}" =~ ^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-(alpha|beta|rc)\.(0|[1-9][0-9]*))?$ ]]; then
    echo "Invalid release version: ${version}" >&2
    echo "Expected format: vMAJOR.MINOR.PATCH, MAJOR.MINOR.PATCH, or vMAJOR.MINOR.PATCH-(alpha|beta|rc).N" >&2
    echo "Examples: v1.2.3, 1.2.3, v1.2.3-alpha.1, v1.2.3-beta.1, v1.2.3-rc.1; prerelease N must not have extra leading zeroes" >&2
    exit 1
fi

printf '%s\n' "${version#v}"
