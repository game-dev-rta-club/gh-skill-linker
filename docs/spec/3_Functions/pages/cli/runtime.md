---
title: CLI runtime specification
updated: 2026-07-16
status: implemented
---

# CLI runtime

The [command reference](../../../2_HowToUse/pages/command-reference.md) owns
public syntax and output. This page covers dispatch and preflight behavior.

## CLI-001 Dispatch and parsing

Commands: `install`, `publish`, `status`, `pull`, `push`, and `uninstall`.

- install: `OWNER/REPO [SKILL|PATH | --all] (--branch BRANCH | --tag TAG) [--accept-moved-tag]`
- publish: `OWNER/REPO SKILL --branch BRANCH [--pr]`
- status: no flag, or `--json` only
- pull: exactly one selector
- push: exactly one selector, with optional `--pr`
- uninstall: exactly one selector, with optional `--force`

The CLI validates install and publish repositories as `OWNER/REPO`. A missing
repository, `.git` suffix, or more than two segments is a usage error. Every
two-segment value is treated as a GitHub repository. Install without a selector
performs discovery only.

Help forms are `--help`, `-h`, and `help [command]`. A help flag anywhere in a
command takes precedence. Help and malformed arguments do not initialize
dependencies. The [command reference](../../../2_HowToUse/pages/command-reference.md)
owns public help text.

Other short flags, `--flag=value`, and `--` are not supported.

## CLI-002 Exit mapping

Success is `0`, operation failure is `1`, and usage error is `2`. A pull that
writes conflict markers also exits with `1`.

## CLI-003 Preflight implementation

1. `auth.TokenForHost("github.com")`
2. successful `gh --version`; output is not parsed
3. `git rev-parse --show-toplevel`

Preflight runs after argument validation. Usage errors and help take precedence
over a missing token.

## CLI-004 Status rendering

Merge managed status and provider inventory by normalized path, then sort by
display path. The table shows provider status and managed operation details;
JSON also includes exact paths and preserves existing managed fields. JSON
disables HTML escaping. Inventory warnings use stderr without invalidating JSON.
Replace terminal control characters in table cells and warnings with spaces;
JSON encoding escapes them instead.

## CLI-005 Mutation rendering

Convert results to success or no-op messages. A rendering failure does not
change the exit code.

A conflicting pull leaves stdout empty and writes sorted, project-relative
paths to stderr as `CONFLICT (content)`. It then recommends `status` and, only
for `STATE=push`, `push`, before exiting with `1`.
