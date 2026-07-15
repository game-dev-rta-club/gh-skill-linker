---
title: Linked Skills distribution and support
updated: 2026-07-15
status: archived
---

# Distribution and support

## Extension

The repository, extension, and project names are all `gh-linked-skills`. Release
it as a precompiled Go GitHub CLI extension so users do not need a Go runtime.

```bash
gh extension install game-dev-rta-club/gh-linked-skills
gh extension upgrade game-dev-rta-club/gh-linked-skills
```

Require an update before operation for GitHub CLI versions earlier than 2.96.0.
Do not maintain a compatibility layer for older versions.

## Workflow skill

Embed the workflow skill in the extension binary and install it without an
additional download.

```bash
gh linked-skills skills install --agent codex
gh linked-skills skills install --agent claude-code
```

Provide project scope only. Running the same command after an extension upgrade
updates only managed bundles with no local changes.

## Support

| Target | MVP |
| --- | --- |
| Host | GitHub.com |
| OS | macOS, Linux |
| Git | system Git required |
| Managed skill | Codex-compatible, project scope |
| Workflow adapter | Codex, Claude Code |

GHES, Windows, user scope, and custom destinations are not officially
supported.

## Related

- [Overview](gh-linked-skills.md)
- [Implementation](gh-linked-skills-implementation.md)
