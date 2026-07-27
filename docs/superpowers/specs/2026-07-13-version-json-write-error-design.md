# Version JSON Write Error Propagation Design

## Context

Issue 416 reports that `gc version --json` discards errors returned by
`cmdutil.WriteJSON`. The current command uses Cobra's `Run` callback, so a
failed stdout writer cannot report its error through `cmd.Execute()`.

The contributor account cannot modify labels in the upstream repository.
Validation evidence has therefore been recorded in the issue comment, while
the maintainer retains ownership of `status/verified` and other workflow
labels. Development and PR delivery use the contributor fork
`zhaogev5_87/cli`.

## Goal

Make `gc version --json` return its stdout write error without changing its
successful JSON payload or the human-readable version output.

## Design

Change `NewCmdVersion` from Cobra `Run` to `RunE` so the command callback can
return an error. In the JSON branch, return `cmdutil.WriteJSON(out, info)`
directly. In the human-readable branch, keep the existing output statements
and return `nil` after they complete.

This is deliberately limited to the reported JSON path. Propagating every
individual `fmt.Fprint*` error in the human-readable path would expand the
issue into a separate output-hardening change.

## Test Strategy

Add a package-local writer whose `Write` method always returns a sentinel
error. A focused test will:

1. construct the version command;
2. enable `--json`;
3. set the failing writer as command output;
4. execute the command; and
5. assert `errors.Is(err, sentinel)`.

The test must first fail against the current implementation because the error
is discarded. After the minimal production change, run the focused package
test, the full test suite, a repository build, scoped pre-commit checks, and a
normal `gc version --json` smoke test.

## Error Handling and Compatibility

- Successful JSON output remains byte-for-byte governed by `cmdutil.WriteJSON`.
- Failed JSON writes become observable through Cobra's returned error.
- Human-readable success output remains unchanged.
- No API, authentication, repository, or remote write path is involved.

## Documentation

No command reference update is required: flags, successful output schema,
examples, and documented behavior remain unchanged. The fix enforces the
existing CLI error-propagation contract.

## Delivery

Commit the focused test and implementation on `bugfix/issue-416`, push the
branch to `zhaogev5_87/cli`, and create a cross-repository PR targeting
`gitcode-cli/cli:main`. The PR will include local test/build/pre-commit
evidence, security and documentation checks, risk classification, and the
maintainer-only label limitation. Its final association line will be
`Closes #416`; other PR text will avoid bare issue references.
