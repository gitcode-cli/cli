---
name: gitcode-release
description: Use GitCode CLI release commands to list, view, create, edit, upload, download, and delete releases. Trigger for GitCode release publishing, asset upload/download, release notes, or version delivery.
---

# gitcode-release

Manage GitCode releases with GitCode CLI.

## Command Entry

Use `gitcode` for cross-platform instructions. Linux/macOS may use `gc`; Windows PowerShell should use `gitcode`.

## Inspect Existing Releases

```bash
gitcode release list -R owner/repo --json
gitcode release view v1.0.0 -R owner/repo --json
gitcode release view v1.0.0 -R owner/repo --web
```

## Create or Edit

Use files for release notes:

```bash
gitcode release create v1.0.0 -R owner/repo --title "v1.0.0" --notes-file RELEASE_NOTES.md --target main --json
gitcode release create v1.0.0-rc1 -R owner/repo --title "v1.0.0-rc1" --notes-file RELEASE_NOTES.md --prerelease --json
gitcode release create v1.0.0 -R owner/repo --title "v1.0.0" --notes-file RELEASE_NOTES.md --draft --json

gitcode release edit v1.0.0 -R owner/repo --title "New title" --json
gitcode release edit v1.0.0 -R owner/repo --notes-file RELEASE_NOTES.md --json
gitcode release edit v1.0.0 -R owner/repo --prerelease false --json
```

`--notes` and `--notes-file` are mutually exclusive.

## Assets

```bash
gitcode release upload v1.0.0 dist/app.zip -R owner/repo --json
gitcode release upload v1.0.0 dist/app.zip dist/checksum.txt -R owner/repo --json

gitcode release download -R owner/repo
gitcode release download v1.0.0 -R owner/repo
gitcode release download v1.0.0 app.zip -R owner/repo -o ./downloads/
```

## Delete

```bash
gitcode release delete v1.0.0 -R owner/repo --dry-run
gitcode release delete v1.0.0 -R owner/repo --yes
```

## Rules

- Confirm the tag, target branch, and release notes before publishing.
- Use `--json` for automation and post-publish verification.
- Use `--dry-run` before deleting a release.
- Do not upload secrets, local config, or unreviewed build artifacts.
- After publishing, verify by viewing the release and, when relevant, downloading an asset.
