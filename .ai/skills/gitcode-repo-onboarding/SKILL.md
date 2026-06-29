---
name: gitcode-repo-onboarding
description: Onboard a contributor to a GitCode repository using GitCode CLI and local repository inspection. Trigger when users ask how to clone, understand, build, test, or contribute to a GitCode repository.
---

# gitcode-repo-onboarding

Help a user understand and contribute to a GitCode repository.

## Command Entry

Use `gitcode` for cross-platform instructions. Linux/macOS may use `gc`; Windows PowerShell should use `gitcode`.

## Workflow

1. Inspect repository metadata.
2. Clone with SSH when local code is needed.
3. Read repository entry documents and build files.
4. Identify build, test, and contribution workflow.
5. Produce a concise onboarding guide.

## Commands

```bash
gitcode repo view owner/repo --json
gitcode repo stats --branch main -R owner/repo --json
ssh -T git@gitcode.com
gitcode repo clone owner/repo --git-protocol ssh
```

After clone:

```bash
git remote -v
git branch --show-current
git status --short --branch
```

Inspect common entry files:

```bash
ls
head -120 README.md 2>/dev/null
head -120 CONTRIBUTING.md 2>/dev/null
head -120 AGENTS.md 2>/dev/null
head -120 CLAUDE.md 2>/dev/null
```

Build-system hints:

```bash
ls go.mod package.json pyproject.toml setup.py Cargo.toml pom.xml build.gradle Makefile 2>/dev/null
```

## Output

```markdown
## Repository Onboarding: owner/repo

### What This Repo Is
- ...

### Local Setup
```bash
gitcode repo clone owner/repo --git-protocol ssh
cd repo
```

### Build and Test
- ...

### Development Workflow
- Branching:
- Issues:
- PR:
- Review:

### First Useful Next Steps
- ...
```

## Rules

- Do not invent build commands; infer them from repository files or state uncertainty.
- Follow target repository instructions when they exist.
- Use SSH for code download.
- Keep target project rules above generic advice.
