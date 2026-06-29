---
name: gitcode-release-helper
description: Prepare and publish GitCode releases using GitCode CLI, including release notes, tag checks, asset upload/download, release verification, and post-release checks. Trigger for release preparation, release notes, version publishing, or release asset management.
---

# gitcode-release-helper

Plan, publish, and verify a GitCode release.

## Command Entry

Use `gitcode` for cross-platform instructions. Linux/macOS may use `gc`; Windows PowerShell should use `gitcode`.

## Workflow

1. Confirm version, target branch, and previous release.
2. Collect merged PRs and commits.
3. Generate release notes.
4. Build or locate assets.
5. Create release and upload assets.
6. Verify by viewing release and downloading assets when needed.

## Inspect

```bash
gitcode release list -R owner/repo --json
gitcode release view v1.0.0 -R owner/repo --json
gitcode pr list -R owner/repo --state merged --limit 50 --json
gitcode repo stats --branch main -R owner/repo --json
git log <previous_tag>..HEAD --oneline --no-merges
```

## Release Notes

```markdown
## vX.Y.Z

### Added
- ...

### Changed
- ...

### Fixed
- ...

### Security
- ...

### Verification
- ...
```

## Publish

```bash
gitcode release create v1.0.0 -R owner/repo --title "v1.0.0" --notes-file RELEASE_NOTES.md --target main --json
gitcode release upload v1.0.0 dist/app.zip dist/checksums.txt -R owner/repo --json
gitcode release view v1.0.0 -R owner/repo --json
gitcode release download v1.0.0 app.zip -R owner/repo -o ./release-verify/
```

Draft or prerelease:

```bash
gitcode release create v1.0.0-rc1 -R owner/repo --title "v1.0.0-rc1" --notes-file RELEASE_NOTES.md --prerelease --json
gitcode release create v1.0.0 -R owner/repo --title "v1.0.0" --notes-file RELEASE_NOTES.md --draft --json
```

## Rules

- Use `--notes-file` for release notes.
- Verify tags and assets before announcing the release.
- Do not upload unreviewed artifacts or files containing secrets.
- Record exact commands and verification results.
