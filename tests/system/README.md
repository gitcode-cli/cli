# System Tests

`tests/system/` contains real GitCode CLI command tests.

All repository targets, including read and write tests, must be under
`infra-test/*`. The runner rejects any other repository before running cases.

## Read-Only Suite

```bash
tests/system/run.sh --read
```

This builds `./gc` and runs read-only, dry-run, JSON contract, authentication
guard, and error-code cases against `infra-test/gctest1` by default.

Use a different infra-test repository:

```bash
tests/system/run.sh --read --repo infra-test/another-repo
```

## Write Suite

Write tests are opt-in and must also target `infra-test/*`.

```bash
tests/system/run.sh --write --write-repo infra-test/gctest1
```

The issue write case creates an issue and closes it during cleanup. The PR write
case is skipped unless an existing test branch is provided:

```bash
GC_SYSTEM_PR_HEAD=test-branch tests/system/run.sh --write --write-repo infra-test/gctest1
```

## Safety Rules

- Do not use personal repositories.
- Do not use production repositories.
- Do not use `gitcode-cli/cli`.
- Do not print, read, or pipe real tokens.
- Keep write cases self-cleaning whenever the remote API allows it.
